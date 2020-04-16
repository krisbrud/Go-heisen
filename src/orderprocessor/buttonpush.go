package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/orderrepository"
	"fmt"
)

// handleButtonPush creates an order and sends it to be delegated if no equivalent order already exists.
func handleButtonPush(
	pushedButton elevator.ButtonEvent,
	repoptr *orderrepository.OrderRepository,
	incomingOrdersChan chan elevator.Order,
	toDelegate chan elevator.Order,
) {
	fmt.Println("handleButtonPush")
	if !pushedButton.IsValid() {
		return
	}

	order := makeUnassignedOrder(pushedButton)

	if !repoptr.HasEquivalentOrders(order) {
		if order.Class == elevator.BT_Cab {
			// My cab call, assign to me
			order.RecipentID = elevator.GetMyElevatorID()
			go func() { incomingOrdersChan <- order }()
		} else {
			// No active orders are equivalent, have the new order delegated.
			go func() { toDelegate <- order }()
		}
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
