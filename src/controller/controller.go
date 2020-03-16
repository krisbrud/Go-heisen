package controller

import (
	"Go-heisen/src/elevatorio"
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"fmt"
)

func Controller(
	toButtonPushHandler chan order.Order,
	toArrivedFloorHandler chan elevatorstate.ElevatorState,
	toDelegator chan elevatorstate.ElevatorState,
	incomingOrders chan order.Order,
	// readState chan elevatorstate.ElevatorState,
	// readQueue chan order.Order,
) {
	// Initialize driver
	numFloors := 4 // TODO: Refactor
	elevatorio.Init("localhost:15657", numFloors)

	drvButtons := make(chan elevatorio.ButtonEvent)
	drvFloors := make(chan int)

	go elevatorio.PollButtons(drvButtons)
	go elevatorio.PollFloorSensor(drvFloors)

	// Initialize belief state
	var state elevatorstate.ElevatorState
	var destination int

	for {
		select {
		case buttonPushed := <-drvButtons:
			// TODO: Make order based on pushed button
			// go func() { toButtonPushHandler <- buttonPushed }() // må sende dette til button pushed handler

		// case elevatorStateChanged := <-readState:
		// 	if elevatorStateChanged != state {
		// 		state := elevatorStateChanged
		// 		go func() { toDelegator <- elevatorStateChanged }() // må lese elevator state og sende den til delegatoren
		// 	}

		case arrivedFloor := <-drvFloors:
			if arrivedFloor == destination {
				elevatorio.SetMotorDirection(elevatorio.MD_Stop)
				// TODO: send state instead
				// go func() { toArrivedFloorHandler <- drvFloors }() // må sende beskjed til arrived floor handler
			}

			// case orderToExecute := <- readQueue:
			// 	destination = orderToExecute.Floor
			// 	if floor < destination {
			// 		elevatorio.SetMotorDirection(elevatorio.MD_Up)
			// 	} else if floor > destination {
			// 		elevatorio.SetMotorDirection(elevatorio.MD_Down)
			// 	}
		}
	}
}

func setOrderLights(incomingOrder order.Order) {
	if !incomingOrder.IsValid() {
		fmt.Printf("Trying to set lights of invalid order! Order: %#v", incomingOrder)
		return
	}

	if incomingOrder.IsFromCab() && !incomingOrder.IsMine() {
		// Cab call but not mine, don't change any lights
		return
	}

	// All is good, set the lights
	elevatorio.SetButtonLamp(0, incomingOrder.Floor, incomingOrder.Completed) //BT_hallup, 0 eller HALL_UP??
}
