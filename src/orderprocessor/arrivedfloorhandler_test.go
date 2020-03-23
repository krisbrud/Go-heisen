package orderprocessor

// import (
// 	"Go-heisen/src/elevator"
// 	"Go-heisen/src/order"
// 	"Go-heisen/src/orderrepository"
// 	"Go-heisen/src/testutils"
// 	"testing"
// )

// func TestOrderArrivedFloorHandler(t *testing.T) {
// 	stateUpdates := make(chan elevator.Elevator)
// 	readAllActiveRequests := make(chan orderrepository.ReadRequest)
// 	toOrderProcessor := make(chan order.Order)

// 	go ArrivedFloorHandler(stateUpdates, readAllActiveRequests, toOrderProcessor)

// 	readSingleRequests := make(chan orderrepository.ReadRequest)
// 	writeRequests := make(chan orderrepository.WriteRequest)

// 	go orderrepository.OrderRepository(readSingleRequests, readAllActiveRequests, writeRequests)

// 	// Insert order to be cleared on arrive into orderrepository
// 	someOrder := testutils.GetSomeOrder()
// 	writeReq := orderrepository.MakeWriteRequest(someOrder)
// 	writeRequests <- writeReq

// 	if writeSuccess := <-writeReq.SuccessCh; !writeSuccess {
// 		t.Errorf("Couldn't write order to OrderRepo in ArrivedFloorHandler test!: %v", someOrder)
// 	}

// 	incomingState := elevator.Elevator{
// 		Floor:       someOrder.Floor,
// 		IntendedDir: elevator.MD_Up,
// 		Behaviour:   elevator.EB_Idle,
// 	}
// 	stateUpdates <- incomingState

// 	// Expect the order to be sent to OrderProcessor, and that it is cleared
// 	someClearedOrder := <-toOrderProcessor

// 	if !someClearedOrder.IsValid() || !someClearedOrder.Completed {
// 		t.Errorf("Error in cleared floor in ArrivedFloorHandler test!: %v", someOrder)
// 	}
// }
