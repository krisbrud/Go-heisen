package orderprocessor

import (
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
)

// OrderProcessor processes an incoming order from this or other elevators
func OrderProcessor(
	incomingOrdersChan chan order.Order,
	singleReadRequests chan orderrepository.ReadRequest,
	repoWriteRequests chan orderrepository.WriteRequest,
	toController chan order.Order,
	toTransmitter chan order.Order,
) {
	for {
		select {
		case incomingOrder := <-incomingOrdersChan:
			go func() {
				fmt.Printf("\nIncoming order in processor! %v\n", incomingOrder)

				if !incomingOrder.IsValid() {
					return // Ignore the incoming order
				}

				// Check OrderRepository for order with same ID:
				readReq := orderrepository.MakeReadRequest(incomingOrder.OrderID)
				singleReadRequests <- readReq

				if localOrder := <-readReq.ResponseCh; localOrder.IsValid() {
					// Order exists in OrderRepository.
					switch {
					case localOrder.Completed && !incomingOrder.Completed:
						// Notify other nodes that order is actually completed.
						// Don't update the OrderRepository, local state is newer.
						toTransmitter <- localOrder
					case !localOrder.Completed && incomingOrder.Completed:
						// Overwrite existing order as completed. Update controller.
						writeReq := orderrepository.MakeWriteRequest(incomingOrder)
						repoWriteRequests <- writeReq
						if !<-writeReq.SuccessCh {
							fmt.Println("Error: Could not overwrite order to repo for some reason!")
						}
						toController <- incomingOrder
					}
				} else {
					// Incoming order is new. Register to OrderRepository, send to controller and transmitter.
					writeReq := orderrepository.MakeWriteRequest(incomingOrder)
					repoWriteRequests <- writeReq
					fmt.Println("Trying to write new order to processor")
					if <-writeReq.SuccessCh {
						toController <- incomingOrder
						toTransmitter <- incomingOrder
					} else {
						fmt.Println("Error: Could not write new order to repo for some reason!")
						// TODO maybe add restart
					}
				}
			}()
		}
	}
}
