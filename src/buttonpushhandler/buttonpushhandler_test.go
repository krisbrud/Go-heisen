package buttonpushhandler

import (
	"fmt"
	"testing"

	"Go-heisen/src/order"
	"Go-heisen/src/orderprocessor"
	"Go-heisen/src/orderrepository"
	"Go-heisen/src/testutils"
)

func TestButtonPushHandler(t *testing.T) {
	// unusedReadReq := make(chan orderrepository.ReadRequest)
	// readAllRequests := make(chan orderrepository.ReadRequest)
	// incomingOrdersChan := make(chan order.Order)
	// anotherIncomingOrdersChan := make(chan order.Order)
	// toDelegator := make(chan order.Order)
	// writeRequests := make(chan orderrepository.WriteRequest)
	// toTransmitter := make(chan order.Order)
	// toController := make(chan order.Order)

	receiveOrder := make(chan order.Order)
	readAllOrdersRequest := make(chan orderrepository.ReadRequest)
	toDelegator := make(chan order.Order)

	go orderrepository.OrderRepository(unusedReadReq, readAllRequests, writeRequests)
	go orderprocessor.OrderProcessor(incomingOrdersChan, unusedReadReq, writeRequests, toController, toTransmitter)
	go ButtonPushHandler(anotherIncomingOrdersChan, readAllRequests, toDelegator)

	// Send a new, valid order to the OrderProcessor
	newOrder := testutils.GetSomeOrder()
	fmt.Println(newOrder)

	incomingOrdersChan <- newOrder

	// Read the order back from OrderRepository to check that it was written correctly
	readReq := orderrepository.MakeReadAllActiveRequest()
	readAllRequests <- readReq

	anotherOrder := testutils.GetSomeOtherOrder()
	anotherIncomingOrdersChan <- anotherOrder

}
