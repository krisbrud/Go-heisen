package orderrepository

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
	"fmt"
)

// type RepoReader string // TODO: Implement

// OrderRepository is the single source of truth of all known orders in all nodes.
func OrderRepository(
	readRequests chan readrequest.ReadRequest,
	processorWrites chan order.Order,
	processorReads chan order.Order,
	buttonPushReads chan order.Order,
	arrivedFloorReads chan order.Order,
	watchdogReads chan order.Order,
) {
	allOrders := make(map[string]order.Order)

	for {
		select {
		case readReq := <-readRequests:
			fmt.Println(readReq)
			storedOrder, ok := allOrders[readReq.OrderID]

			if !ok {
				// Order does not exist, inform Reader by sending invalid order back
				storedOrder = order.NewInvalidOrder()
			}

			switch readReq.Reader {
			case readrequest.OrderProcessor:
				processorReads <- storedOrder
			case readrequest.ButtonPushHandler:
				buttonPushReads <- storedOrder
			case readrequest.ArrivedFloorHandler:
				arrivedFloorReads <- storedOrder
			case readrequest.Watchdog:
				watchdogReads <- storedOrder
			default:
				fmt.Printf("ERROR! Unknown reader: %v", readReq.Reader)

			}

		case orderToWrite := <-processorWrites:
			if orderToWrite.IsValid() {
				allOrders[orderToWrite.OrderID] = orderToWrite
			} else {
				fmt.Printf("Trying to print invalid order: %v", orderToWrite)
			}
		}
	}
}
