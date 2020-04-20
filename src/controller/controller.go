package controller

import (
	"Go-heisen/src/config"
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

// Controller controls the elevator motor and lights, and
func Controller(
	activeOrdersUpdates chan []elevator.Order,
	buttonPushes chan elevator.ButtonEvent,
	stateUpdates chan elevator.State,
	toArrivedFloorHandler chan elevator.State,
) {
	// Initialize driver for ElevatorServer
	elevio.Init(fmt.Sprintf("localhost:%v", config.GetElevatorDriverPort()), config.GetNumFloors())

	buttonUpdates := make(chan elevator.ButtonEvent)
	floorUpdates := make(chan int)

	go elevio.PollButtons(buttonUpdates)
	go elevio.PollFloorSensor(floorUpdates)

	// Initialize internal elevator state
	state := elevator.UninitializedElevatorBetweenFloors()

	// Run elevator downwards if no floor update is received
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

	// Initialize timer for door
	doorTimer := time.NewTimer(math.MaxInt64)
	doorTimer.Stop()

	activeOrders := make([]elevator.Order, 0, orderCapacity)
	setAllLights(activeOrders) // Clear all the lights, in case some of them were turned on when we started

	for {
		select {
		case newFloor := <-floorUpdates:
			fmt.Printf("Floor update: %#v\n", newFloor)
			state.Floor = newFloor
			state.Print()
			elevio.SetFloorIndicator(state.Floor)

			if shouldStop(state, activeOrders) { // && state.Behaviour == elevator.EB_Moving
				elevio.SetMotorDirection(elevator.MD_Stop)
				state.IntendedDir = chooseDirection(state, activeOrders)

				// Open the door
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				state.Behaviour = elevator.EB_DoorOpen

				// Make OrderProcessor clear fulfilled orders (if any)
				go func() { toArrivedFloorHandler <- state }()
			}
			// Sending new state to delegator
			go func() { stateUpdates <- state }()

		case buttonEvent := <-buttonUpdates:
			// Print state?
			fmt.Printf("Buttonevent: %#v\n", buttonEvent)
			state.Print()

			if !state.IsValid() {
				continue // Don't take orders if we have not reached a valid floor yet
			}

			go func() {
				buttonPushes <- buttonEvent
			}()

		case <-doorTimer.C:
			// Door timer timed out, close door.
			fmt.Println("Door timer!")
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
			// Sending new state to delegator
			go func() { stateUpdates <- state }()

		//Received new active order list from order processor
		case activeOrders = <-activeOrdersUpdates:
			fmt.Println("Update of all orders received!")
			elevator.PrintOrders(activeOrders)

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
