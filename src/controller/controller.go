package controller

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/elevio"
	"Go-heisen/src/order"
	"math"
	"time"
)

const (
	doorDuration  = 3 * time.Second
	orderCapacity = 100
)

func Controller(
	incomingOrders chan order.Order,
	buttonPushes chan elevator.ButtonEvent,
	stateUpdates chan elevator.Elevator,
	toArrivedFloorHandler chan elevator.Elevator,
	// TODO config
) {
	elevio.Init("localhost:15657", 4)

	buttonUpdates := make(chan elevator.ButtonEvent)
	floorUpdates := make(chan int)

	go elevio.PollButtons(buttonUpdates)
	go elevio.PollFloorSensor(floorUpdates)

	/* void fsm_onInitBetweenFloors(void){
		outputDevice.motorDirection(elevator.MD_Down);
		elev.IntendedDir = elevator.MD_Down;
		elev.Behaviour = EB_Moving;
	} */

	// elev := initialize elev between floors

	elev := elevator.MakeInvalidState()

	// Initialize timer for doors
	doorTimer := time.NewTimer(math.MaxInt64)
	doorTimer.Stop()

	activeOrders := make([]order.Order, 0, orderCapacity)

	for {
		select {
		case buttonEvent := <-buttonUpdates:
			// Print state?

			switch elev.Behaviour { // Cases are mutually exclusive
			case elevator.EB_DoorOpen:
				if elev.Floor == buttonEvent.Floor {
					doorTimer.Reset(doorDuration)
					// timer_start(elev.config.doorOpenDuration_s);
				} else {
					buttonPushes <- buttonEvent
				}

			case elevator.EB_Moving:
				buttonPushes <- buttonEvent

			case elevator.EB_Idle:
				if elev.Floor == buttonEvent.Floor {
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(doorDuration)
					// elev.Behaviour = EB_DoorOpen; // TODO refactor
				} else {
					buttonPushes <- buttonEvent
					nextDir := chooseDirection(elev, activeOrders)
					elev.IntendedDir = nextDir
					elevio.SetMotorDirection(elev.IntendedDir)
					elev.Behaviour = elevator.EB_Moving
				}
			}

			stateUpdates <- elev

			// setAllLights(elev); // Set on incomingorder and completed order instead

			// printf("\nNew state:\n");
			// elevator_print(elev);

		case newFloor := <-floorUpdates:
			// TODO maybe print something

			elev.Floor = newFloor
			elevio.SetFloorIndicator(elev.Floor)

			if elev.Behaviour == elevator.EB_Moving && shouldStop(elev, activeOrders) {
				// Clear the orders we have fulfilled
				toArrivedFloorHandler <- elev

				// Stop the elevator
				elevio.SetMotorDirection(elevator.MD_Stop)

				// Open the door
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				elev.Behaviour = elevator.EB_DoorOpen

				// setAllLights(elev);
			}

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

			// Possibly print new state

		}
	}
}

func shouldStop(elev elevator.Elevator, activeOrders []order.Order) bool {
	switch elev.IntendedDir {
	case elevator.MD_Down:
		shouldStopAtOrder := func(o order.Order) bool {
			atSameFloor := elev.Floor == o.Floor
			notOppositeDirection := o.IsFromCab() || o.Class == order.HALL_UP
			isOrderBelow := o.Floor < elev.Floor

			return (atSameFloor && notOppositeDirection) || !isOrderBelow
		}
		return anyOrder(activeOrders, shouldStopAtOrder)

	case elevator.MD_Up:
		shouldStopAtOrder := func(o order.Order) bool {
			atSameFloor := elev.Floor == o.Floor
			notOppositeDirection := o.IsFromCab() || o.Class == order.HALL_DOWN
			isOrderAbove := o.Floor > elev.Floor

			return (atSameFloor && notOppositeDirection) || !isOrderAbove
		}
		return anyOrder(activeOrders, shouldStopAtOrder)
	}

	return true // Default
}

func anyOrder(orderList []order.Order, predicateFunc func(o order.Order) bool) bool {
	satisfied := false

	for _, o := range orderList {
		satisfied = satisfied || predicateFunc(o)
	}

	return satisfied
}

func ordersAbove(elev elevator.Elevator, activeOrders []order.Order) bool {
	isAbove := func(o order.Order) bool {
		return o.Floor > elev.Floor
	}

	return anyOrder(activeOrders, isAbove)
}

func ordersBelow(elev elevator.Elevator, activeOrders []order.Order) bool {
	isBelow := func(o order.Order) bool {
		return o.Floor < elev.Floor
	}

	return anyOrder(activeOrders, isBelow)
}

func chooseDirection(elev elevator.Elevator, activeOrders []order.Order) elevator.MotorDirection {
	switch elev.IntendedDir {
	case elevator.MD_Up:
		switch {
		case ordersAbove(elev, activeOrders):
			return elevator.MD_Up
		case ordersBelow(elev, activeOrders):
			return elevator.MD_Down
		}

	case elevator.MD_Down, elevator.MD_Stop:
		switch {
		case ordersBelow(elev, activeOrders):
			return elevator.MD_Down
		case ordersAbove(elev, activeOrders):
			return elevator.MD_Up
		}
	}
	return elevator.MD_Stop // Default case
}
