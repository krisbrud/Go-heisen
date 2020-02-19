package orderrepository

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
	"testing"
)

func TestOrderRepository(t *testing.T) {
	unused := make(chan order.Order)
	orderProcessorReads := make(chan order.Order)
	orderProcessorWrites := make(chan order.Order)
	readRequests := make(chan readrequest.ReadRequest)

	go OrderRepository(readRequests, orderProcessorWrites, orderProcessorReads, unused, unused, unused)

	nonExistingID := "Non-existent;)"
	myReadReq := readrequest.ReadRequest{
		nonExistingID,
		readrequest.OrderProcessor,
	}

	readRequests <- myReadReq

	if result := <-orderProcessorReads; result.IsValid() {
		t.Errorf("Order that should not exist exists!: %v", result)
	}

	// Test writing some order and reading it back

	someOrder := order.Order{
		"Some ID",
		1,
		order.CAB,
		"Some recipent",
		false,
	}

	orderProcessorWrites <- someOrder

	readRequests <- readrequest.ReadRequest{
		OrderID: someOrder.OrderID,
		Reader:  readrequest.OrderProcessor,
	}

	if result := <-orderProcessorReads; result != someOrder {
		t.Errorf("Did not get expected order when reading!: %v", result)
	}

	// Write another order to same ID, see if it is overwritten
	someOtherOrder := order.Order{
		"Some ID",
		2,
		order.CAB,
		"Some other recipent",
		false,
	}

	orderProcessorWrites <- someOtherOrder

	readRequests <- readrequest.ReadRequest{someOrder.OrderID, readrequest.OrderProcessor}

	if result := <-orderProcessorReads; result == someOrder {
		t.Errorf("Did not get expected order when reading after rewrite!: %v", result)
	}
}
