package controller

import (
	"Go-heisen/src/elevatorio"
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/queue"
	"fmt"
)

/*
ArriveFloor:
	update floor in state

	if should stop:
		update atfloor in state
		send to ArrivedFloorHandler
		remove from orders
		set lights for completed order
		if door closed:
			open door

	send state to delegator

ButtonPush:
	if atfloor and idle and same floor:
		open door
		update door state
	else:
		Send to buttonpushhandler (ordercreator)

IncomingOrder:
	if not valid:
		break

	if completed:
		remove all equivalent from queue
	else:
		add to queue

	set lights for order

DoorTimer:
	turn off door light
	set door as closed
	if more orders, execute order
*/

func Controller(
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
	queue := queue.MakeEmptyQueue()

	for {
		select {
		case newFloor := <-drvFloors:
			state.CurrentFloor = newFloor

		}
	}
}

func shouldStop(q queue.OrderQueue, state elevatorstate.ElevatorState) bool {

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
