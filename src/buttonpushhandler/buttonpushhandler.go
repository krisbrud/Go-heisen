package buttonpushhandler

import (
	//"../order"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
)


// ButtonPushHandler checks if an order already exist
func ButtonPushHandler(
	receiveOrder chan order.Order,
	readAllOrdersRequest chan orderrepository.ReadRequest,
	toDelegator chan order.Order,
) {
	var orderExist = false;

	for {
		select {
		case o := <- receiveOrder:
			fmt.Println("inside case")
			if o.IsValid(){	
				fmt.Println("order is valid")
				readReq :=  orderrepository.MakeReadAllActiveRequest()
				for order := range readReq.ResponseCh{
					fmt.Println("inside for loop")
					if order == o{
						orderExist = true; 
						fmt.Println("Order already exists")
						readReq.ResponseCh <- order
					}
				}
				close(readReq.ResponseCh)
				if !orderExist{
					toDelegator <- o
					fmt.Println("Order sent to delegator")
				}
			}
		}
	}
}