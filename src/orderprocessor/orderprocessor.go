package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
)

// OrderManager order from this or other elevators
func OrderManager(
	incomingOrdersChan chan order.Order,
	buttonPushes chan elevator.ButtonEvent,
	toController chan order.Order,
	toDelegate chan order.Order,
	toRedelegate chan order.Order,
	toTransmit chan order.Order,
) {
	allOrders := orderrepository.MakeEmptyOrderRepository()


	for {
		select {
		case incomingOrder := <-incomingOrdersChan:
			handleIncomingOrder(incomingOrder, &allOrders, toController, toDelegate, toTransmit)
		case buttonPush := <-buttonPushes:
			handleButtonPush(buttonPush, &allOrders, toDelegate)
		case 
		}
	}
}

func handleIncomingOrder(
	incomingOrder order.Order,
	allOrders *orderrepository.OrderRepository,
	toController chan order.Order,
	toDelegate chan order.Order,
	toTransmit chan order.Order,
) {
	// TODO: Comment here
	fmt.Printf("\nHandling incoming order!\n")
	incomingOrder.Print()

	if !incomingOrder.IsValid() {
		return // Ignore the incoming order
	}

	localOrder, err := allOrders.ReadSingleOrder(incomingOrder.OrderID)
	exists := err != nil

	if exists {
		switch {
		case localOrder.Completed && !incomingOrder.Completed:
			// Notify other nodes that order is actually completed.
			// Don't update the OrderRepository, local state is newer.
			go func() {
				toTransmit <- localOrder
			}()
		case !localOrder.Completed && incomingOrder.Completed:
			// Overwrite existing order as completed. Update controller.
			allOrders.WriteOrderToRepository(incomingOrder)
			go func() {
				toController <- incomingOrder
			}()
		}
	} else {
		// Incoming order is new. Register to OrderRepository, send to controller and transmitter.
		allOrders.WriteOrderToRepository(incomingOrder)
		go func() {
			toController <- incomingOrder
			toTransmit <- incomingOrder
		}()
	}
}

// handleButtonPush creates an order and sends it to be delegated if no equivalent order already exists.
func handleButtonPush(
	pushedButton elevator.ButtonEvent,
	repoptr *orderrepository.OrderRepository,
	toDelegate chan order.Order,
) {
	if !pushedButton.IsValid() {
		return
	}

	o := order.MakeUnassignedOrder(pushedButton)

	if !repoptr.HasEquivalentOrders(o) {
		go func() {
			// No active orders are equivalent, have the new order delegated.
			toDelegate <- o
		}()
	}
}
