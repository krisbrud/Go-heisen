package arrivedfloorhandler

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"Go-heisen/src/testutils"
	"testing"
)

func TestOrderArrivedFloorHandler(t *testing.T) {
	stateUpdates := make(chan elevatorstate.ElevatorState)
	readAllActiveRequests := make(chan orderrepository.ReadRequest)
	toOrderProcessor := make(chan order.Order)

	go ArrivedFloorHandler(stateUpdates, readAllActiveRequests, toOrderProcessor)

	readSingleRequests := make(chan orderrepository.ReadRequest)
	writeRequests := make(chan orderrepository.WriteRequest)

	go orderrepository.OrderRepository(readSingleRequests, readAllActiveRequests, writeRequests)

	// Insert order to be cleared on arrive into orderrepository
	someOrder := testutils.GetSomeOrder()
	writeReq := orderrepository.MakeWriteRequest(someOrder)
	writeRequests <- writeReq

	if writeSuccess := <-writeReq.SuccessCh; !writeSuccess {
		t.Errorf("Couldn't write order to OrderRepo in ArrivedFloorHandler test!: %v", someOrder)
	}

	incomingState := elevatorstate.ElevatorState{
		CurrentFloor: someOrder.Floor,
		RelPos:       elevatorstate.AtFloor,
	}
	stateUpdates <- incomingState

	// Expect the order to be sent to OrderProcessor, and that it is cleared
	someClearedOrder := <-toOrderProcessor

	if !someClearedOrder.IsValid() || !someClearedOrder.Completed {
		t.Errorf("Error in cleared floor in ArrivedFloorHandler test!: %v", someOrder)
	}
}
