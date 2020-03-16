package buttonpushhandler

import (
	"fmt"

	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
)

// ButtonPushHandler checks if an order already exist
func ButtonPushHandler(
	receiveOrder chan order.Order,
	readAllOrdersRequests chan orderrepository.ReadRequest,
	toDelegator chan order.Order,
) {
	for {
		select {
		case o := <-receiveOrder:
			if !o.IsValid() {
				fmt.Println("Invalid order in ButtonPushHandler!")
				break
			}

			readReq := orderrepository.MakeReadAllActiveRequest()
			readAllOrdersRequests <- readReq

			orderExists := false
			for order := range readReq.ResponseCh {
				if order == o {
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
