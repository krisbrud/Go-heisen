package orderprocessor

import (
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"Go-heisen/src/testutils"
	"testing"
)

func TestOrderProcessor(t *testing.T) {
	unusedReadReq := make(chan orderrepository.ReadRequest)
	singleReadRequests := make(chan orderrepository.ReadRequest)
	writeRequests := make(chan orderrepository.WriteRequest)
	incomingOrdersChan := make(chan order.Order)
	toTransmitter := make(chan order.Order)
	toController := make(chan order.Order)

	go orderrepository.OrderRepository(singleReadRequests, unusedReadReq, writeRequests)
	go OrderProcessor(incomingOrdersChan, singleReadRequests, writeRequests, toController, toTransmitter)

	// Send a new, valid order to the OrderProcessor
	newOrder := testutils.GetSomeOrder()
	incomingOrdersChan <- newOrder

	// Read the order back from OrderRepository to check that it was written correctly
	readReq := orderrepository.MakeReadRequest(newOrder.OrderID)
	singleReadRequests <- readReq

	if fromOrderRepo := <-readReq.ResponseCh; fromOrderRepo != newOrder {
		t.Errorf("Order not read back correctly!: %v", fromOrderRepo)
	}
	// Check that the order was sent to the right channels
	if controllerMsg := <-toController; controllerMsg != newOrder {
		t.Errorf("Did not send order to controller!: %v", controllerMsg)
	}

	if transmitterMsg := <-toTransmitter; transmitterMsg != newOrder {
		t.Errorf("Did not send order to transmitter!: %v", transmitterMsg)
	}

	// Send the first order to OrderProcessor again. Nothing should happen.
	incomingOrdersChan <- newOrder

	// Set the order to completed and send to OrderProcessor. It should tell Controller.
	newOrderCompleted := newOrder
	newOrderCompleted.SetCompleted()
	incomingOrdersChan <- newOrderCompleted

	if controllerMsg := <-toController; controllerMsg != newOrderCompleted {
		t.Errorf("Did not send order to controller!: %v", controllerMsg)
	}

	// Try to send the first order to OrderProcessor.
	// It should notify the other nodes that it actually is completed by sending it as completed to transmitter.
	incomingOrdersChan <- newOrder

	if transmitterMsg := <-toTransmitter; transmitterMsg != newOrderCompleted {
		t.Errorf("OrderProcessor did not notify transmitter that order is actually completed!: %v", transmitterMsg)
	}

	// Try to send an invalid order to the OrderProcessor. It should not write it to OrderRepository
	invalidOrder := order.NewInvalidOrder()
	incomingOrdersChan <- invalidOrder

	invalidReadReq := orderrepository.MakeReadRequest(invalidOrder.OrderID)
	singleReadRequests <- invalidReadReq

	if fromOrderRepo := <-invalidReadReq.ResponseCh; fromOrderRepo.IsValid() {
		t.Errorf("The invalid order was written to OrderRepository by OrderProcessor!: %v", fromOrderRepo)
	}

}
