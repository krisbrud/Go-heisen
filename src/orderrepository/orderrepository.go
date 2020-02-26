package orderrepository

import (
	"Go-heisen/src/order"
)

// ReadRequest serves as a request to read an order from OrderRepository
type ReadRequest struct {
	OrderID    string
	ResponseCh chan order.Order
	ErrorCh    chan error
}

// WriteRequest makes it possible for other modules to write to OrderRepository
type WriteRequest struct {
	OrderToWrite order.Order
	ErrorCh      chan error
}

// The InvalidRepoRequestError is returned on ErrorCh in ReadRequest or WriteRequest if something is wrong.
type InvalidRepoRequestError struct {
	why string
}

func (e InvalidRepoRequestError) Error() string { return e.why }

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

			if ok {
				readReq.ErrorCh <- nil
			} else {
				// Order does not exist, inform Reader by sending invalid order back
				storedOrder = order.NewInvalidOrder()
				readReq.ErrorCh <- InvalidRepoRequestError{"Order with ID: " + readReq.OrderID + " already exists."}
			}

			readReq.ResponseCh <- storedOrder
			close(readReq.ResponseCh)
			close(readReq.ErrorCh)

		case readReq := <-readAllActiveRequests:
			// Read back all orders on the request
			for _, storedOrder := range allOrders {
				if !storedOrder.Completed {

				}
			}

		case writeReq := <-writeRequests:
			if writeReq.OrderToWrite.IsValid() {
				allOrders[writeReq.OrderToWrite.OrderID] = writeReq.OrderToWrite
				writeReq.ErrorCh <- nil
			} else {
				writeReq.ErrorCh <- InvalidRepoRequestError{"Trying to write invalid order."}
			}
		}
	}
}
