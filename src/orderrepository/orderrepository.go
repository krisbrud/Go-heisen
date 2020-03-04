package orderrepository

import (
	"Go-heisen/src/order"
	"fmt"
)

// ReadRequest serves as a request to read an order from OrderRepository
type ReadRequest struct {
	OrderID    string
	ResponseCh chan order.Order
}

// MakeReadRequest returns a ReadRequest with a new response channel
func MakeReadRequest(OrderID string) ReadRequest {
	return ReadRequest{OrderID, make(chan order.Order)}
}

// MakeWriteRequest returns a WriteRequest with a success response channel
func MakeWriteRequest(orderToWrite order.Order) WriteRequest {
	return WriteRequest{orderToWrite, make(chan bool)}
}

// WriteRequest makes it possible for other modules to write to OrderRepository
type WriteRequest struct {
	OrderToWrite order.Order
	SuccessCh    chan bool
}

// OrderRepository is the single source of truth of all known orders in all nodes.
func OrderRepository(
	readSingleRequests chan ReadRequest,
	readAllActiveRequests chan ReadRequest,
	writeRequests chan WriteRequest,
) {
	allOrders := make(map[string]order.Order)

	for {
		select {
		case readReq := <-readSingleRequests:
			storedOrder, ok := allOrders[readReq.OrderID]
			if !ok {
				// Order does not exist, inform Reader by sending invalid order back
				storedOrder = order.NewInvalidOrder()
			}

			readReq.ResponseCh <- storedOrder // Send order back to requester
			close(readReq.ResponseCh)

		case readReq := <-readAllActiveRequests:
			// Read back all orders on the request and close channel afterwards
			for _, storedOrder := range allOrders {
				if !storedOrder.Completed {
					if storedOrder.IsValid() {
						readReq.ResponseCh <- storedOrder
					} else {
						// Invalid Order in OrderRepository for some reason. Restart.
						// TODO restart mechanics
						fmt.Println("Invalid order in repository!")
					}
				}
			}
			close(readReq.ResponseCh)

		case writeReq := <-writeRequests:
			if writeReq.OrderToWrite.IsValid() {
				allOrders[writeReq.OrderToWrite.OrderID] = writeReq.OrderToWrite
				go func() { writeReq.SuccessCh <- true }() // Don't wait for SuccessCh to be read
			} else {
				go func() { writeReq.SuccessCh <- false }()
			}
		}
	}
}
