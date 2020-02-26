package orderrepository

import (
	"Go-heisen/src/order"
	"testing"
)

func TestOrderRepository(t *testing.T) {
	readSingleRequests := make(chan ReadRequest)
	readAllRequests := make(chan ReadRequest)
	writeRequests := make(chan WriteRequest)

	go OrderRepository(readSingleRequests, readAllRequests, writeRequests)

	nonExistingID := "Non-existent;)"
	myReadReq := MakeReadRequest(nonExistingID)

	readSingleRequests <- myReadReq

	if result := <-myReadReq.ResponseCh; result.IsValid() {
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

	someWriteReq := WriteRequest{
		someOrder,
		make(chan bool),
	}
	writeRequests <- someWriteReq
	<-someWriteReq.successCh

	someReadRequest := MakeReadRequest(someOrder.OrderID)
	readSingleRequests <- someReadRequest

	if result := <-someReadRequest.ResponseCh; result != someOrder {
		t.Errorf(" <- someOrderDid not get expected order when reading!: %v", result)
	}

	// Write another order to same ID, should overwrite
	someOtherOrder := order.Order{
		"Some ID",
		2,
		order.CAB,
		"Some other recipent",
		false,
	}

	someOtherWriteReq := WriteRequest{
		someOtherOrder,
		make(chan bool),
	}

	writeRequests <- someOtherWriteReq
	<-someOtherWriteReq.successCh

	someOtherReadReq := ReadRequest{
		someOtherOrder.OrderID,
		make(chan order.Order),
	}
	readSingleRequests <- someOtherReadReq

	if result := <-someOtherReadReq.ResponseCh; result != someOtherOrder {
		t.Errorf("Did not get expected order when reading after overwrite!: %v", result)
	}
}
