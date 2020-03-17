package controller

import (
	"Go-heisen/src/elevatorio"
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"math"
	"time"
)

const (
	doorDuration = 3 * time.Second
)

func Controller(
	incomingOrders chan order.Order,
	createOrder chan elevatorio.ButtonEvent,
	stateUpdates chan elevatorstate.ElevatorState,
	toArrivedFloorHandler chan elevatorstate.ElevatorState,
	// TODO config
) {
	elevatorio.Init("localhost:15657", 4)

	buttonUpdates := make(chan elevatorio.ButtonEvent)
	floorUpdates := make(chan int)

	go elevatorio.PollButtons(buttonUpdates)
	go elevatorio.PollFloorSensor(floorUpdates)

	/* void fsm_onInitBetweenFloors(void){
		outputDevice.motorDirection(D_Down);
		elev.dirn = D_Down;
		elev.behaviour = EB_Moving;
	} */

	// elev := initialize elev between floors

	// Initialize timer for doors
	doorTimer := time.NewTimer(math.MaxInt64)
	doorTimer.Stop()

	for {
		select {
		case buttonEvent := <-buttonUpdates:
			// Print state?

			switch { // Cases are mutually exclusive
			case elev.IsDoorOpen():
				if elev.Floor == buttonEvent.Floor {
					doorTimer.Reset(doorDuration)
					// timer_start(elev.config.doorOpenDuration_s);
				} else {
					createOrder <- buttonEvent
				}

			case elev.IsMoving():
				createOrder <- buttonEvent

			case elev.IsIdle():
				if elev.Floor == btn_floor {
					elevatorio.SetDoorOpenLamp(true)
					doorTimer.Reset(doorDuration)
					// elev.behaviour = EB_DoorOpen; // TODO refactor
				} else {
					createOrder <- buttonEvent
					elev.dirn = chooseDirection(elev)
					elevatorio.SetMotorDirection(elev.dirn)
					elev.behaviour = EB_Moving
				}
			}

			stateUpdates <- elev

			// setAllLights(elev); // Set on incomingorder and completed order instead

			// printf("\nNew state:\n");
			// elevator_print(elev);

		case newFloor := <-floorUpdates:
			// TODO maybe print something

			elev.Floor = newFloor
			elevatorio.SetFloorIndicator(elev.Floor)

			if elev.IsMoving() && shouldStop(elev) {
				// Clear the orders we have fulfilled
				toArrivedFloorHandler <- elev

				// Stop the elevator
				elevatorio.SetMotorDirection(elevatorio.MD_Stop)

				// Open the door
				elevatorio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				elev.behaviour = EB_DoorOpen

				// setAllLights(elev);
			}

		case <-doorTimer.C:
			// Door timer timed out, close door.
			elevatorio.SetDoorOpenLamp(false)

			// Find and set motor direction
			elev.dirn = chooseDirection(elev)
			elevatorio.SetMotorDirection(elev.dirn)

			// Set the behaviour accordingly
			if elev.IsMotorStopped() {
				elev.SetBehaviourIdle()
			} else {
				elev.SetBehaviourMoving()
			}

			// Possibly print new state

		}
	}
}

func shouldStop(elev elevator.Elevator) bool {
	switch e.dirn {
	case D_Down:
		shouldStopAtOrder := func(o order.Order) {
			atSameFloor := elev.Floor == o.Floor
			notOppositeDirection := o.IsFromCab() || o.Class == order.HALL_UP

			return (atSameFloor && notOppositeDirection) || !isBelow(o, elev.Floor)
		}
		return anyOrder(elev.ActiveOrders, shouldStopAtOrder)

	case D_Up:
		shouldStopAtOrder := func(o order.Order) {
			atSameFloor := elev.Floor == o.Floor
			notOppositeDirection := o.IsFromCab() || o.Class == order.HALL_DOWN

			return (atSameFloor && notOppositeDirection) || !isAbove(o, elev.Floor)
		}
		return anyOrder(elev.ActiveOrders, shouldStopAtOrder)
	}

	return true // Default
}

func anyOrder(orderList []order.Order, predicateFunc func(o order.Order) bool) bool {
	satisfied := false

	for _, o := range orderList {
		satisfied = satisfied | filterFunc(o)
	}

	return satisfied
}

func isAbove(o order.Order, floor int) bool {
	return o.Floor > elev.Floor
}

func isBelow(o order.Order) bool {
	return o.Floor < elev.Floor
}

func chooseDirection(elev elevator.Elevator) {
	switch {
	case elev.dirn == D_Up:
		switch {
		case anyOrder(elev.ActiveOrders, isAbove):
			return D_Up
		case anyOrder(elev.ActiveOrders, isBelow):
			return D_Down
		}

	case elev.Dirn == D_Up || elev.Dirn == D_Down:
		switch {
		case anyOrder(elev.ActiveOrders, isBelow):
			return D_Down
		case anyOrder(elev.ActiveOrders, isAbove):
			return D_Up
		}
	default:
		return D_Stop // Default case
	}
}
