package buttonpushhandler

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// ButtonPushHandler checks if an order alre
func ButtonPushHandler(
	buttonPushOrders chan order.Order,
	buttonRepoReads chan order.Order,
	repoReadRequests chan readrequest.ReadRequest,
	toDelegatorChan chan order.Order,
) {
	for {
		select {}
	}
}
