package buttonpushhandler

import (
	"Go-heisen/src/elevatorio"
	"fmt"

	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
)

// ButtonPushHandler checks if an order already exist
func ButtonPushHandler(
	buttonPush chan elevatorio.ButtonEvent,
	readAllOrdersRequests chan orderrepository.ReadRequest,
	toDelegator chan order.Order,
) {
	for {
		select {
		case buttonEvent := <-buttonPush:
			// Make order based on button push
			o := makeUnassignedOrder(buttonEvent)

			if !o.IsValid() {
				fmt.Println("Invalid order in ButtonPushHandler!")
				break
			}

			// Check that equivalent orders don't exist already 
			readReq := orderrepository.MakeReadAllActiveRequest()
			readAllOrdersRequests <- readReq

			orderExists := false
			for existingOrder :|= range readReq.ResponseCh {
				if existingOrder == o {
					orderExists = true
					readReq.ResponseCh <- order
				}
			}
			close(readReq.ResponseCh)

			if !orderExists {
				toDelegator <- o
				fmt.Println("Order sent to delegator")
			}
		}
	}
}

func makeUnassignedOrder(pushedButton elevatorio.ButtonEvent) order.Order {
	return order.Order{
		OrderID:    order.GetRandomID(),
		Floor:      pushedButton.Floor,
		Class:      order.OrderClass(pushedButton.Button), // TODO Verify that definitions are the same
		RecipentID: "",
		Completed:  false,
	}
}
