package orderprocessor

import (
	"Go-heisen/src/elevator"
	"time"
)

// OrderProcessor handles the single source of truth of active and completed orders in the system.
func OrderProcessor(
	incomingOrdersChan chan elevator.Order,
	buttonPushes chan elevator.ButtonEvent,
	floorArrivals chan elevator.State,
	toController chan []elevator.Order,
	toDelegate chan elevator.Order,
	toWatchdog chan []elevator.Order,
	toTransmit chan elevator.Order,
) {
	allOrders := makeEmptyOrderRepository()
	watchdogTicker := time.NewTicker(500 * time.Millisecond) //checking and distributing all orders every 500ms

	for {
		select {
		case elevAtFloor := <-floorArrivals:
			// Clear relevant orders when arriving at floor, send back the completed order(s) to incomingOrdersChan.
			activeOrders := allOrders.readActiveOrders()
			clearOrdersOnFloorArrival(elevAtFloor, activeOrders, incomingOrdersChan, toTransmit)

		case buttonPush := <-buttonPushes:
			// Create orders from button push to be delegated if no equivalent active order exists.
			activeOrders := allOrders.readActiveOrders()
			handleButtonPush(buttonPush, activeOrders, incomingOrdersChan, toDelegate)

		case <-watchdogTicker.C:
			// Static redundancy, resend all active orders to other nodes
			// This solves most issues from packet loss and disconnects/reconnects/restarts
			activeOrders := allOrders.readActiveOrders()
			go func() {
				for _, activeOrder := range activeOrders {
					toTransmit <- activeOrder
				}
			}()

			// Dynamic redundancy, make the watchdog find orders that are too old,
			// and make the Delegator redelegate them
			go func() { toWatchdog <- activeOrders }()

		case incomingOrder := <-incomingOrdersChan:
			// Update the OrderRepository of the incoming order
			// Also notifies other nodes if receiving an order we know is completed
			// Sends all active orders to the controller if the state has changed

			if !incomingOrder.IsValid() {
				continue // Ignore the invalid incoming order
			}

			// Check if this node already has an order with the same ID
			localOrder, err := allOrders.readSingleOrder(incomingOrder.OrderID)
			orderAlreadyExists := err == nil

			if orderAlreadyExists {
				// Check if the status of the orders are different
				switch {
				case localOrder.Completed && !incomingOrder.Completed:
					// Notify other nodes that order is actually completed.
					// Don't update the OrderRepository, local state is newer.
					go func() { toTransmit <- localOrder }()
					continue // Don't resend all active orders to controller

				case !localOrder.Completed && incomingOrder.Completed:
					// Overwrite existing order as completed. Update controller.
					allOrders.writeOrderToRepository(incomingOrder)
				default:
					continue // No changes, don't resend orders to controller
				}
			} else {
				// Incoming order is new. Register to OrderRepository, send to controller and transmitter.
				allOrders.writeOrderToRepository(incomingOrder)
				go func() {
					toTransmit <- incomingOrder
				}()
			}
			// Update the controller about the current active orders
			activeOrders := allOrders.readActiveOrders()
			go func() { toController <- activeOrders }()
		}
	}
}

func resendAllActiveOrders(
	activeOrders []elevator.Order,
	toTransmit chan elevator.Order,
) {
	for _, activeOrder := range activeOrders {
		toTransmit <- activeOrder
	}
}
