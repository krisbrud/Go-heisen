package orderprocessor

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// OrderProcessor processes an incoming order from this or other elevators
func OrderProcessor(
	incomingOrdersChan chan order.Order,
	orderTxChan chan order.Order,
	readRequestChan chan readrequest.ReadRequest,
	orderRepoRead chan order.Order,
	orderRepoWriteChan chan order.Order,
	toController chan order.Order,
	toLightManager chan order.Order,
) {
	for {
		select {
		case incomingOrder := <-incomingOrdersChan:
			if incomingOrder.IsValid() {
				// Check OrderRepository for order with same ID:
				readReq := readrequest.ReadRequest{
					OrderID: incomingOrder.OrderID,
					Reader:  readrequest.OrderProcessor}
				readRequestChan <- readReq

				if localOrder := <-orderRepoRead; localOrder.IsValid() {
					switch {
					case localOrder.Completed && !incomingOrder.Completed:
						// Notify other nodes that order is actually completed.
						// Don't update the OrderRepository, local state is newer.
						orderTxChan <- localOrder
					case !localOrder.Completed && incomingOrder.Completed:
						// Overwrite existing order as completed. Update lights.
						orderRepoWriteChan <- incomingOrder

					}
				} else {
					// Incoming order is new. Register to OrderRepository.
					orderRepoWriteChan <- incomingOrder
					if incomingOrder.IsMine() {
						toController <- incomingOrder
					}
					toLightManager <- incomingOrder // Update lights
				}
			}
		}
	}
}
