package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
)

// handleButtonPush creates an order and sends it to be delegated if no equivalent order already exists.
func handleButtonPush(
	pushedButton elevator.ButtonEvent,
	repoptr *orderrepository.OrderRepository,
	incomingOrdersChan chan order.Order,
	toDelegate chan order.Order,
) {
	fmt.Println("handleButtonPush")
	if !pushedButton.IsValid() {
		return
	}

	o := order.MakeUnassignedOrder(pushedButton)

	if !repoptr.HasEquivalentOrders(o) {
		if o.Class == elevator.BT_Cab {
			// My cab call, assign to me
			o.RecipentID = elevator.GetMyElevatorID()
			go func() { incomingOrdersChan <- o }()
		} else {
			// No active orders are equivalent, have the new order delegated.
			go func() { toDelegate <- o }()
		}
	}
}

func makeUnassignedOrder(pushedButton elevator.ButtonEvent) order.Order {
	return order.Order{
		OrderID:    order.GetRandomID(),
		Floor:      pushedButton.Floor,
		Class:      pushedButton.Button, // TODO Verify that definitions are the same
		RecipentID: "",
		Completed:  false,
	}
}
