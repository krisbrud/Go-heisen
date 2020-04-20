package delegator

import (
	"Go-heisen/src/elevator"
	"fmt"
	"testing"
)

func TestCost(t *testing.T) {
	// Cost from floor 1 to 3 while going up should be 2
	order := getMockThirdFloorCabCall()
	es := getMockElevatorStateFirstFloorUp()
	correctCost := 2
	if foundcost := cost(order, es); foundcost != correctCost {
		t.Errorf(makeCostErrorString(order, es, foundcost, correctCost))
	}

	// Cost to third floor while standing still at the third floor should be zero
	es = getMockElevatorStateAtThirdFloor()
	correctCost = 0
	if foundcost := cost(order, es); foundcost != correctCost {
		t.Errorf(makeCostErrorString(order, es, foundcost, correctCost))
	}

	// Cost to third floor while standing still at the third floor should be zero
	order = getMockFloorZeroCabUpOrder()
	es = getMockElevatorStateFirstFloorUp()
	correctCost = 5
	if foundcost := cost(order, es); foundcost != correctCost {
		t.Errorf(makeCostErrorString(order, es, foundcost, correctCost))
	}

}

func getMockFloorZeroCabUpOrder() elevator.Order {
	return elevator.Order{
		OrderID:    12345,
		Floor:      0,
		Class:      elevator.BT_HallUp,
		RecipentID: "",
		Completed:  false,
	}
}

func getMockElevatorStateFirstFloorUp() elevator.State {
	return elevator.State{
		Floor:       1,
		IntendedDir: elevator.MD_Up,
		Behaviour:   elevator.EB_Moving,
		ElevatorID:  "SomeElevator",
	}
}

func getMockThirdFloorCabCall() elevator.Order {
	return elevator.Order{
		OrderID:    234678,
		Floor:      3,
		Class:      elevator.BT_Cab,
		RecipentID: "SomeElevator",
		Completed:  false,
	}
}

func getMockElevatorStateAtThirdFloor() elevator.State {
	return elevator.State{
		Floor:       3,
		IntendedDir: elevator.MD_Stop,
		Behaviour:   elevator.EB_Idle,
		ElevatorID:  "SomeElevator",
	}
}

func makeCostErrorString(order elevator.Order, es elevator.State, gotCost int, correctCost int) string {
	return fmt.Sprintf("Cost of order %#v while in state %#v should be %v, but was: %v", order, es, correctCost, gotCost)
}
