package delegator

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"testing"
)

// TODO:
// Redelegate the order, check that elev2 should take it
// Send new state updates such that elev3 is best fit to take the order
// Redelegate it again, it should

func TestDelegator(t *testing.T) {
	// Initialize delegator and channels
	toDelegate := make(chan order.Order)
	toRedelegate := make(chan order.Order)
	toTransmitter := make(chan order.Order)
	toProcessor := make(chan order.Order)
	stateUpdates := make(chan elevatorstate.ElevatorState)

	go Delegator(
		toDelegate,
		toRedelegate,
		toTransmitter,
		toProcessor,
		stateUpdates,
	)

	// Send state updates s.t. elev1 should take order. Based on example in spec.
	stateUpdates <- getStateIdleAtFloor(1, "elev1")
	stateUpdates <- getStateIdleAtFloor(0, "elev2")
	stateUpdates <- getStateIdleAtFloor(0, "elev3")

	testOrder := makeTopFloorOrder()
	toDelegate <- testOrder

	// Expect processor to get order with closest recipent
	correctRecipent := "elev1"
	delegated := <-toProcessor
	if delegated.RecipentID != correctRecipent {
		t.Errorf("Order %#v should have been delegated to %#v!", delegated, correctRecipent)
	}
	<-toTransmitter // Ignore

	// Send new state updates so elev2 should take order when redelegated
	// Elev 1 is still closest, but disallowed
	stateUpdates <- getStateIdleAtFloor(2, "elev1")
	stateUpdates <- getStateIdleAtFloor(1, "elev2")
	toRedelegate <- delegated

	correctRecipent = "elev2"
	if redelegated := <-toProcessor; redelegated.RecipentID != correctRecipent {
		t.Errorf("Order %#v should have been redelegated to %#v!", redelegated, correctRecipent)
	}
	<-toTransmitter // Ignore

}

func makeTopFloorOrder() order.Order {
	return order.Order{
		OrderID:    order.GetRandomID(),
		Floor:      3,
		Class:      order.HALL_DOWN,
		RecipentID: "",
		Completed:  false,
	}
}

func getStateIdleAtFloor(floor int, id string) elevatorstate.ElevatorState {
	return elevatorstate.ElevatorState{
		CurrentFloor: floor,
		AtFloor:      true,
		IntendedDir:  elevatorstate.Idle,
		ElevatorID:   id,
	}
}
