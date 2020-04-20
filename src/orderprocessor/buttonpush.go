package orderprocessor

import (
	"Go-heisen/src/config"
	"Go-heisen/src/elevator"
	"fmt"
)

// handleButtonPush creates an order and sends it to be delegated if no equivalent order already exists.
func handleButtonPush(
	pushedButton elevator.ButtonEvent,
	activeOrders []elevator.Order,
	incomingOrdersChan chan elevator.Order,
	toDelegate chan elevator.Order,
) {
	fmt.Println("handleButtonPush")
	if !pushedButton.IsValid() {
		return
	}

	order := makeUnassignedOrder(pushedButton)

	// Check that no existing orders are equivalent with the new order
	for _, activeOrder := range activeOrders {
		if activeOrder.IsEquivalentWith(order) {
			// The order is equivalent with one that already exists, don't create a new one.
			return
		}
	}

	// Send the order to self or to be delegated to the best elevator
	if order.Class == elevator.BT_Cab {
		// My cab call, assign to me
		order.RecipentID = config.GetMyElevatorID()
		go func() { incomingOrdersChan <- order }()
	} else {
		// No active orders are equivalent, have the new order delegated.
		go func() { toDelegate <- order }()
	}
}

func makeUnassignedOrder(pushedButton elevator.ButtonEvent) elevator.Order {
	return elevator.Order{
		OrderID:    elevator.GetRandomID(),
		Floor:      pushedButton.Floor,
		Class:      pushedButton.Button,
		RecipentID: "",
		Completed:  false,
	}
}
