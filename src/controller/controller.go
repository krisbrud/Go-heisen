package controller

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/elevio"
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
	activeOrdersUpdates chan elevator.OrderList,
	buttonPushes chan elevator.ButtonEvent,
	stateUpdates chan elevator.State,
	toArrivedFloorHandler chan elevator.State,
) {
	// Initialize driver for ElevatorServer
	elevio.Init(fmt.Sprintf("localhost:%v", elevator.GetElevatorDriverPort()), elevator.GetNumFloors())

	buttonUpdates := make(chan elevator.ButtonEvent)
	floorUpdates := make(chan int)

	go elevio.PollButtons(buttonUpdates)
	go elevio.PollFloorSensor(floorUpdates)

	// Initialize internal elevator state
	state := elevator.UninitializedElevatorBetweenFloors()

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
		state.IntendedDir = elevator.MD_Down
		state.Behaviour = elevator.EB_Moving
		elevio.SetMotorDirection(elevator.MD_Down)
	}

	// Initialize timer for doors
	doorTimer := time.NewTimer(math.MaxInt64)
	doorTimer.Stop()

	activeOrders := make(elevator.OrderList, 0, orderCapacity)

	for {
		select {
		case buttonEvent := <-buttonUpdates:

			// Print state?
			fmt.Printf("Buttonevent: %#v\n", buttonEvent)
			state.Print()

			if !state.IsValid() {
				continue
			}

			buttonPushes <- buttonEvent

		case newFloor := <-floorUpdates:
			fmt.Printf("Floor update: %#v\n", newFloor)
			state.Floor = newFloor
			state.Print()
			elevio.SetFloorIndicator(state.Floor)

			if shouldStop(state, activeOrders) { // && state.Behaviour == elevator.EB_Moving
				elevio.SetMotorDirection(elevator.MD_Stop)
				// Don't change the IntendedDir to MD_Stop,
				// so we may continue in same direction when door closes

				// Open the door
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				state.Behaviour = elevator.EB_DoorOpen

				// Make orderprocessor the orders we have fulfilled TODO: OrderManager or processor
				go func() { toArrivedFloorHandler <- state }()
			}

			stateUpdates <- state

			// fmt.Println("After floor update")
			// state.Print()

		case <-doorTimer.C:
			// Door timer timed out, close door.
			elevio.SetDoorOpenLamp(false)

			// Find and set motor direction
			state.IntendedDir = chooseDirection(state, activeOrders)
			elevio.SetMotorDirection(state.IntendedDir)

			// Set the Behaviour accordingly
			if state.IntendedDir == elevator.MD_Stop {
				state.Behaviour = elevator.EB_Idle
			} else {
				state.Behaviour = elevator.EB_Moving
			}

			stateUpdates <- state

		case activeOrders = <-activeOrdersUpdates:
			fmt.Println("Update of all orders received!")
			//state.Print()
			activeOrders.Print()

			state.IntendedDir = chooseDirection(state, activeOrders)
			if ordersAtCurrentFloor(state, activeOrders) {
				switch state.Behaviour {
				case elevator.EB_Idle, elevator.EB_DoorOpen:
					// Open/re-open the door
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(doorDuration)
					state.Behaviour = elevator.EB_DoorOpen

					// Notify ArrivedFloorHandler that we handled the order at our floor (by opening door)
					go func() { toArrivedFloorHandler <- state }()
				}
			}

			// Execute movement in intended direction if elevator is idle
			if state.Behaviour == elevator.EB_Idle {
				elevio.SetMotorDirection(state.IntendedDir)

				if state.IntendedDir != elevator.MD_Stop {
					state.Behaviour = elevator.EB_Moving
				}
			}

			setAllLights(activeOrders)
		}
	}
}
