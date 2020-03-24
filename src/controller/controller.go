package controller

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/elevio"
	"Go-heisen/src/order"
	"fmt"
	"math"
	"time"
)

const (
	doorDuration  = 3 * time.Second
	orderCapacity = 100
)

/*
ButtonUpdate:
	if at floor of button and not moving
		open/reopendoor
	else
		send command

FloorUpdate:
	set new floor in state

	stop if needed
		set in state
		set motor direction

		send floor arrival to handler

	send new state to other nodes

CloseDoor:
	find next direction for elevator
		set and execute

ActiveOrdersUpdate:
hensyn: dør, retning, lys
	if should stop at current floor
		open/re-open door

	if door not open
		find and set next direction

	set lights
*/

func Controller(
	activeOrdersUpdates chan order.OrderList,
	buttonPushes chan elevator.ButtonEvent,
	stateUpdates chan elevator.Elevator,
	toArrivedFloorHandler chan elevator.Elevator,
	elevatorPort int,
) {
	elevio.Init(fmt.Sprintf("localhost:%v", elevatorPort), 4)

	buttonUpdates := make(chan elevator.ButtonEvent)
	floorUpdates := make(chan int)

	go elevio.PollButtons(buttonUpdates)
	go elevio.PollFloorSensor(floorUpdates)

	// Initialize internal elevator state
	elev := elevator.UninitializedElevatorBetweenFloors()

	// Run elevator downwards if no state update
	select {
	case newFloor := <-floorUpdates:
		// Send the floor update again on the channel so the normal handler may do it's routine
		fmt.Println("Floor update received, sending state back!")
		go func() {
			fmt.Println("Starting at floor!")
			floorUpdates <- newFloor
		}()
	case <-time.After(200 * time.Millisecond):
		// Elevator initialized between floors, go downwards.
		fmt.Println("Started between floors!")
		elev.IntendedDir = elevator.MD_Down
		elev.Behaviour = elevator.EB_Moving
		elevio.SetMotorDirection(elevator.MD_Down)
	}

	// Initialize timer for doors
	doorTimer := time.NewTimer(math.MaxInt64)
	doorTimer.Stop()

	activeOrders := make(order.OrderList, 0, orderCapacity)

	for {
		select {
		case buttonEvent := <-buttonUpdates:

			// Print state?
			fmt.Printf("Buttonevent: %#v\n", buttonEvent)
			elev.Print()

			if !elev.IsValid() {
				continue
			}

			buttonPushes <- buttonEvent

		case newFloor := <-floorUpdates:
			fmt.Printf("Floor update: %#v\n", newFloor)
			elev.Print()

			elev.Floor = newFloor
			elevio.SetFloorIndicator(elev.Floor)

			if shouldStop(elev, activeOrders) { // && elev.Behaviour == elevator.EB_Moving
				// Stop the elevator
				fmt.Println("Floor reached, elevator should stop.")
				elevio.SetMotorDirection(elevator.MD_Stop)
				// Don't change the IntendedDir, prefer to continue doing orders in same direction

				// Open the door
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				elev.Behaviour = elevator.EB_DoorOpen

				// Make orderprocessor the orders we have fulfilled TODO: OrderManager or processor
				go func() { toArrivedFloorHandler <- elev }()
			}

			stateUpdates <- elev

			// fmt.Println("After floor update")
			// elev.Print()

		case <-doorTimer.C:
			// Door timer timed out, close door.
			elevio.SetDoorOpenLamp(false)

			// Find and set motor direction
			elev.IntendedDir = chooseDirection(elev, activeOrders)
			elevio.SetMotorDirection(elev.IntendedDir)

			// Set the Behaviour accordingly
			if elev.IntendedDir == elevator.MD_Stop {
				elev.Behaviour = elevator.EB_Idle
			} else {
				elev.Behaviour = elevator.EB_Moving
			}

			stateUpdates <- elev

		case activeOrders = <-activeOrdersUpdates:
			fmt.Println("Update of all orders received!")
			elev.Print()
			activeOrders.Print()

			elev.IntendedDir = chooseDirection(elev, activeOrders)
			if ordersAtCurrentFloor(elev, activeOrders) {
				switch elev.Behaviour {
				case elevator.EB_Idle, elevator.EB_DoorOpen:
					// Open/re-open the door
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(doorDuration)
					elev.Behaviour = elevator.EB_DoorOpen

					// Notify ArrivedFloorHandler that we handled the order at our floor (by opening door)
					go func() { toArrivedFloorHandler <- elev }()
				}
			}

			// Execute movement in intended direction if elevator is idle
			if elev.Behaviour == elevator.EB_Idle {
				elevio.SetMotorDirection(elev.IntendedDir)

				if elev.IntendedDir != elevator.MD_Stop {
					elev.Behaviour = elevator.EB_Moving
				}
			}

			setAllLights(activeOrders)
		}
	}
}
