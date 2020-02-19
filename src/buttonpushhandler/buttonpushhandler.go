package buttonpushhandler

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// ButtonPushHandler checks if an order alre
func ButtonPushHandler(
	incomingOrderChan chan order.Order,
	orderRepoReadChan chan order.Order,
	readRequestChan chan readrequest.ReadRequest,
	toDelegatorChan chan order.Order,
) {

}
