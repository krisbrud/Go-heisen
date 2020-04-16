package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/orderrepository"
	"fmt"
)

func handleIncomingOrder(
	incomingOrder elevator.Order,
	allOrders *orderrepository.OrderRepository,
	toController chan []elevator.Order,
	toDelegate chan elevator.Order,
	toTransmit chan elevator.Order,
) {
	fmt.Printf("\nProcessor handling incoming order!\n")
	incomingOrder.Print()

	if !incomingOrder.IsValid() {
		fmt.Println("Incoming order not valid!")
		return // Ignore the invalid incoming order
	}

	localOrder, err := allOrders.ReadSingleOrder(incomingOrder.OrderID)
	orderAlreadyExists := err == nil

	if orderAlreadyExists {
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
