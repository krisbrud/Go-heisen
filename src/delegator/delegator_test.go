package delegator

import (
	"Go-heisen/src/elevator"
	"testing"
)

func TestDelegator(t *testing.T) {
	// Initialize delegator and channels
	toDelegate := make(chan elevator.Order)
	toRedelegate := make(chan elevator.Order)
	toTransmitter := make(chan elevator.Order)
	toProcessor := make(chan elevator.Order)
	stateUpdates := make(chan elevator.State)

	go Delegator(
		toDelegate,
		toRedelegate,
		toTransmitter,
		toProcessor,
		stateUpdates,
	)

	// Send state updates s.t. elev1 should take elevator. Based on example in spec.
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

func makeTopFloorOrder() elevator.Order {
	return elevator.Order{
		OrderID:    config.GetRandomID(),
		Floor:      3,
		Class:      elevator.HALL_DOWN,
		RecipentID: "",
		Completed:  false,
	}
}

func getStateIdleAtFloor(floor int, id string) elevator.State {
	return elevator.State{
		Floor:       floor,
		IntendedDir: elevator.MD_Stop,
		Behaviour:   elevator.EB_Idle,
		ElevatorID:  id,
	}
}
