package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/orderrepository"
	"fmt"
	"time"
)

// OrderProcessor order from this or other elevators
func OrderProcessor(
	incomingOrdersChan chan elevator.Order,
	buttonPushes chan elevator.ButtonEvent,
	floorArrivals chan elevator.State,
	toController chan []elevator.Order,
	toDelegate chan elevator.Order,
	toWatchdog chan []elevator.Order,
	toTransmit chan elevator.Order,
) {
	allOrders := orderrepository.MakeEmptyOrderRepository()
	watchdogTicker := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case elevAtFloor := <-floorArrivals:
			// Clear relevant orders when arriving at floor, notify OrderProcessor and other nodes.
			clearOrdersOnFloorArrival(elevAtFloor, &allOrders, &allOrders, toController, toTransmit)
		case incomingOrder := <-incomingOrdersChan:
			// Update the OrderRepository of the incoming order
			// Also notifies other nodes if receiving an order we know is completed
			// Sends all active orders to the controller if the state has changed
			handleIncomingOrder(incomingOrder, &allOrders, toController, toDelegate, toTransmit)
		case buttonPush := <-buttonPushes:
			// Create orders from button push to be delegated if needed.
			handleButtonPush(buttonPush, &allOrders, incomingOrdersChan, toDelegate)
		case <-watchdogTicker.C:
			// Static redundancy, resend all active orders to other nodes
			// This solves most issues from packet loss and disconnects/reconnects/restarts
			resendAllActiveOrders(&allOrders, toTransmit)

			// Dynamic redundancy
			activeOrders := allOrders.ReadActiveOrders()
			fmt.Println("Resending all active orders!")
			elevator.PrintOrders(activeOrders)
			toWatchdog <- activeOrders
		}
	}
}

func resendAllActiveOrders(
	repoptr *orderrepository.OrderRepository,
	toTransmit chan elevator.Order,
) {
	for _, activeOrder := range repoptr.ReadActiveOrders() {
		toTransmit <- activeOrder
	}
}
