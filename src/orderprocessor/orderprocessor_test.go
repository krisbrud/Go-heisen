package orderprocessor

import (
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"Go-heisen/src/readrequest"
	"testing"
)

func TestOrderProcessor(t *testing.T) {
	//
	unused := make(chan order.Order)
	orderProcessorReads := make(chan order.Order)
	orderProcessorWrites := make(chan order.Order)
	readRequests := make(chan readrequest.ReadRequest)
	incomingOrdersChan := make(chan order.Order)

	go orderrepository.OrderRepository(readRequests, orderProcessorWrites, orderProcessorReads, unused, unused, unused)
	go OrderProcessor(incomingOrdersChan)

	nonExistingID := "Non-existent;)"
	myReadReq := readrequest.ReadRequest{
		nonExistingID,
		readrequest.OrderProcessor,
	}

}
