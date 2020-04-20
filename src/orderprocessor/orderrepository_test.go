package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/testutils"
	"testing"
)

func TestOrderRepository(t *testing.T) {

	repo := makeEmptyOrderRepository()

	nonExistingID := elevator.OrderIDType(12364)

	// Test writing some order and reading it back
	someOrder := testutils.GetSomeOrder()

	someWriteReq := WriteRequest{
		someOrder,
		make(chan bool),
	}
	writeRequests <- someWriteReq
	<-someWriteReq.SuccessCh

	someReadRequest := MakeReadRequest(someOrder.OrderID)
	readSingleRequests <- someReadRequest

	if result := <-someReadRequest.ResponseCh; result != someOrder {
		t.Errorf(" <- someOrderDid not get expected order when reading!: %v", result)
	}

	// Write another order to same ID, should overwrite
	someOtherOrder := testutils.GetSomeOtherOrder()

	someOtherWriteReq := WriteRequest{
		someOtherOrder,
		make(chan bool),
	}

	writeRequests <- someOtherWriteReq
	<-someOtherWriteReq.SuccessCh

	someOtherReadReq := ReadRequest{
		someOtherOrder.OrderID,
		make(chan elevator.Order),
	}
	readSingleRequests <- someOtherReadReq

	if result := <-someOtherReadReq.ResponseCh; result != someOtherOrder {
		t.Errorf("Did not get expected order when reading after overwrite!: %v", result)
	}
}
