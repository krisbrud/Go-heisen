package buttonpushhandler

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// ButtonPushHandler checks if an order alre
func ButtonPushHandler(
	buttonPushOrders chan order.Order,
	readAllActiveRequests chan orderrepository.ReadRequest,
	toDelegatorChan chan order.Order,
) {
	var activeOrders order.Order; 
	for {
		select {
		case readButtonPush := <- buttonPushOrders:
			// Lag readreq. for alle aktive ordre
			// Vent på resultat
			// Sjekk at ingen aktive ordre er like
			// 		send i så fall videre
			activeOrders := <- readAllActiveRequests;
			var orderExist = false;
			if readButtonPush.IsValid() {
				for order in activeOrders{
					if order == readButtonPush{
						orderExist = true; 
					}
				
				}
				if !orderExist{
					toDelegatorChan <- readButtonPush
				}
			}
		}
	}
}


