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
		return
		e.requests[e.floor][B_HallDown] ||
			e.requests[e.floor][B_Cab] ||
			!requests_below(e)
	case D_Up:
		return
		e.requests[e.floor][B_HallUp] ||
			e.requests[e.floor][B_Cab] ||
			!requests_above(e)
	}
	return true
}

func anyOrders(orderList []order.Order, predicateFunc func(o order.Order) bool) bool {
	satisfied := false

	for _, o := range orderList {
		satisfied = satisfied | filterFunc(o)
	}

	return satisfied
}

func hasOrdersAbove(elev elevator.Elevator) bool {
	isAbove := func(o order.Order) {
		return o.Floor > elev.Floor
	}
	return anyOrders(elev.ActiveOrders, isAbove)
}

func hasOrdersBelow(elev elevator.Elevator) bool {
	isBelow := func(o order.Order) {
		return o.Floor < elev.Floor
	}
	return anyOrders(elev.ActiveOrders, isBelow)
}

/* func chooseDirection(elev elev) {

Dirn requests_chooseDirection(Elevator e){
    switch(e.dirn){
    case D_Up:
        return  requests_above(e) ? D_Up    :
                requests_below(e) ? D_Down  :
                                    D_Stop  ;
    case D_Down:
    case D_Stop: // there should only be one request in this case. Checking up or down first is arbitrary.
        return  requests_below(e) ? D_Down  :
                requests_above(e) ? D_Up    :
                                    D_Stop  ;
    default:
        return D_Stop;
    }
}
} */
