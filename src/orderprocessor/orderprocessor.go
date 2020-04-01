package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
	"time"
)

// OrderProcessor order from this or other elevators
func OrderProcessor(
	incomingOrdersChan chan order.Order,
	buttonPushes chan elevator.ButtonEvent,
	floorArrivals chan elevator.Elevator,
	toController chan order.OrderList,
	toDelegate chan order.Order,
	toWatchdog chan order.OrderList,
	toTransmit chan order.Order,
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

			// Dynamic redund/activeOrders := allOrders.ReadActiveOrders()
			//toWatchdog <- activeOrders
		}
	}
}

func handleIncomingOrder(
	incomingOrder order.Order,
	allOrders *orderrepository.OrderRepository,
	toController chan order.OrderList,
	toDelegate chan order.Order,
	toTransmit chan order.Order,
) {
	fmt.Printf("\nProcessor handling incoming order!\n")
	incomingOrder.Print()

	if !incomingOrder.IsValid() {
		fmt.Println("Incoming order not valid!")
		return // Ignore the incoming order
	}

	localOrder, err := allOrders.ReadSingleOrder(incomingOrder.OrderID)
	exists := err != nil

	if exists {
		fmt.Println("Order already exists!")
		switch {
		case localOrder.Completed && !incomingOrder.Completed:
			// Notify other nodes that order is actually completed.
			// Don't update the OrderRepository, local state is newer.
			go func() { toTransmit <- localOrder }()
			return // Don't resend all active orders to controller

		case !localOrder.Completed && incomingOrder.Completed:
			// Overwrite existing order as completed. Update controller.
			allOrders.WriteOrderToRepository(incomingOrder)
			fmt.Println("Order being marked as completed in processor.")
		default:
			return // No changes, don't resend orders to controller
		}
	} else {
		// Incoming order is new. Register to OrderRepository, send to controller and transmitter.
		fmt.Println("New order incoming in processor")

		allOrders.WriteOrderToRepository(incomingOrder)
		go func() {
			toTransmit <- incomingOrder
		}()
	}
	// Update the controller about the current active orders
	activeOrders := allOrders.ReadActiveOrders()
	go func() { toController <- activeOrders }()
}

func clearOrdersOnFloorArrival(
	elev elevator.Elevator,
	repoptr *orderrepository.OrderRepository,
	allOrders *orderrepository.OrderRepository,
	toController chan order.OrderList,
	transmitOrder chan order.Order,
) {
	fmt.Printf("ArrivedFloorHandler! State: %#v\n", elev)

	if !elev.IsValid() {
		panic("Invalid state in floor arrival handler")
	}

	// Read all active orders from OrderRepository. Set the relevant ones as cleared.
	for _, activeOrder := range repoptr.ReadActiveOrders() {
		if activeOrder.Floor == elev.Floor {
			fmt.Printf("Active order with floor %#v being set to complete\n", activeOrder.Floor)
			if activeOrder.IsFromHall() || (activeOrder.IsFromCab() && activeOrder.IsMine()) {
				// We have completed this order, make OrderProcessor register it and tell everyone.
				activeOrder.SetCompleted()
				allOrders.WriteOrderToRepository(activeOrder)
				activeOrders := allOrders.ReadActiveOrders()
				go func() { toController <- activeOrders }()
				//go func() { handleOrder <- activeOrder }() // New goroutine to avoid deadlock
				go func() { transmitOrder <- activeOrder }()
			}
		}
	}
}

//
func resendAllActiveOrders(
	repoptr *orderrepository.OrderRepository,
	toTransmit chan order.Order,
) {
	for _, activeOrder := range repoptr.ReadActiveOrders() {
		toTransmit <- activeOrder
	}
}
