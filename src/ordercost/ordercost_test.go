package ordercost

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"fmt"
	"testing"
)

func TestCost(t *testing.T) {
	// Cost from floor 1 to 3 while going up should be 2
	o := getMockThirdFloorCabCall()
	es := getMockElevatorStateFirstFloorUp()
	correctCost := 2
	if cost := Cost(o, es); cost != correctCost {
		t.Errorf(makeCostErrorString(o, es, cost, correctCost))
	}

	// Cost to third floor while standing still at the third floor should be zero
	es = getMockElevatorStateAtThirdFloor()
	correctCost = 0
	if cost := Cost(o, es); cost != correctCost {
		t.Errorf(makeCostErrorString(o, es, cost, correctCost))
	}

	// Cost to third floor while standing still at the third floor should be zero
	o = getMockFloorZeroCabUpOrder()
	es = getMockElevatorStateFirstFloorUp()
	correctCost = 5
	if cost := Cost(o, es); cost != correctCost {
		t.Errorf(makeCostErrorString(o, es, cost, correctCost))
	}

}

func getMockFloorZeroCabUpOrder() order.Order {
	return order.Order{
		OrderID:    12345,
		Floor:      0,
		Class:      order.HALL_UP,
		RecipentID: "SomeElevator",
		Completed:  false,
	}
}

func getMockElevatorStateFirstFloorUp() elevatorstate.ElevatorState {
	return elevatorstate.ElevatorState{
		CurrentFloor: 1,
		AtFloor:      false,
		IntendedDir:  elevatorstate.Up,
	}
}

func getMockThirdFloorCabCall() order.Order {
	return order.Order{
		OrderID:    234678,
		Floor:      3,
		Class:      order.CAB,
		RecipentID: "SomeElevator",
		Completed:  false,
	}
}

func getMockElevatorStateAtThirdFloor() elevatorstate.ElevatorState {
	return elevatorstate.ElevatorState{
		CurrentFloor: 3,
		AtFloor:      true,
		IntendedDir:  elevatorstate.Idle,
	}
}

func makeCostErrorString(o order.Order, es elevatorstate.ElevatorState, gotCost int, correctCost int) string {
	return fmt.Sprintf("Cost of order %#v while in state %#v should be %v, but was: %v", o, es, correctCost, gotCost)
}
