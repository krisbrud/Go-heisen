package buttonpushhandler

import (
	"../order"
	"../readrequest"
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
			//TODO sørg for at alle active ordre blir lest, oppdater for løkken når du har fått til riktig syntax på dette. 
			activeOrders := <- readAllActiveRequests; //her var tanken at active ORders skulle være en liste med alle aktive ordre
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


