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
			/*
				if at floor of button and not moving
					open/reopendoor
				else
					send command
					reset timer

				set direction otherwise? - should be removed
			*/
			// Print state?
			fmt.Printf("Buttonevent: %#v\n", buttonEvent)
			elev.Print()

			if !elev.IsValid() {
				continue
			}

			switch elev.Behaviour { // Cases are mutually exclusive
			case elevator.EB_DoorOpen:
				if elev.Floor == buttonEvent.Floor {
					doorTimer.Reset(doorDuration)
				} else {
					buttonPushes <- buttonEvent
				}

			case elevator.EB_Moving:
				buttonPushes <- buttonEvent

			case elevator.EB_Idle:
				if elev.Floor == buttonEvent.Floor {
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(doorDuration)
					elev.Behaviour = elevator.EB_DoorOpen // TODO refactor
				} else {
					buttonPushes <- buttonEvent
					nextDir := chooseDirection(elev, activeOrders)
					elev.IntendedDir = nextDir
					elevio.SetMotorDirection(elev.IntendedDir)
					elev.Behaviour = elevator.EB_Moving
				}
			}
			if elev.IsValid() {
				stateUpdates <- elev
			} else {
				fmt.Println("Elevator not valid!")
				elev.Print()
			}
			// printf("\nNew state:\n");
			// elevator_print(elev);

		case newFloor := <-floorUpdates:
			/*
				set new floor in state

				stop if needed
					set in state
					set motor direction

					send floor arrival to handler

				send new state to other nodes


			*/
			fmt.Printf("Floor update: %#v\n", newFloor)
			elev.Print()

			elev.Floor = newFloor
			elevio.SetFloorIndicator(elev.Floor)

			if shouldStop(elev, activeOrders) { // && elev.Behaviour == elevator.EB_Moving
				fmt.Println("Floor reached, elevator should stop.")
				// Stop the elevator
				elevio.SetMotorDirection(elevator.MD_Stop)

				// Clear the orders we have fulfilled
				go func() { toArrivedFloorHandler <- elev }()

				// Open the door
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				elev.Behaviour = elevator.EB_DoorOpen

			}
			if elev.IsValid() {
				stateUpdates <- elev
			}

			fmt.Println("After floor update")
			elev.Print()

		case <-doorTimer.C:
			/*
				find next direction for elevator
					set and execute

			*/

			// Door timer timed out, close door.
			fmt.Println("Door timer!")
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

			elev.Print()

			// Possibly print new state

		case activeOrders = <-activeOrdersUpdates:
			/*
				if door open
					do nothing
				else
					find and set next direction

				set lights
			*/
			fmt.Println("Update of all orders received!")
			elev.Print()
			activeOrders.Print()

			nextDir := chooseDirection(elev, activeOrders)
			fmt.Printf("\nNext intended direction: %v\n", nextDir)
			elev.IntendedDir = nextDir

			// // Send current state to floor updates
			// if elev.Behaviour == elevator.EB_Idle {
			// 	go func() { floorUpdates <- elev.Floor }()
			// }
			// Choose direction and execute if idle
			if elev.Behaviour != elevator.EB_DoorOpen {
				elevio.SetMotorDirection(elev.IntendedDir)
				elev.Behaviour = elevator.EB_Moving
			}

			setAllLights(activeOrders)
		}
	}
}
