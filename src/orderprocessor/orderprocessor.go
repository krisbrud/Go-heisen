package orderprocessor

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// OrderReceiver processes an incoming order from this or other elevators
func OrderReceiver(
	incomingOrdersChan chan order.Order,
	orderTxChan chan order.Order,
	readRequestChan chan readrequest.ReadRequest,
	orderRepoRead chan order.Order,
	orderWriteChan chan order.Order,
) {
	for {
		select {}
	}
}
