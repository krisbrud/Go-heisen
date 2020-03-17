package ordercost

import (
	"Go-heisen/src/elevator"
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

func getMockElevatorStateFirstFloorUp() elevator.Elevator {
	return elevator.Elevator{
		Floor:       1,
		IntendedDir: elevator.MD_Up,
		Behaviour:   elevator.EB_Moving,
		ElevatorID:  "SomeElevator",
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

func getMockElevatorStateAtThirdFloor() elevator.Elevator {
	return elevator.Elevator{
		Floor:       3,
		IntendedDir: elevator.MD_Stop,
		Behaviour:   elevator.EB_Idle,
		ElevatorID:  "SomeElevator",
	}
}

func makeCostErrorString(o order.Order, es elevator.Elevator, gotCost int, correctCost int) string {
	return fmt.Sprintf("Cost of order %#v while in state %#v should be %v, but was: %v", o, es, correctCost, gotCost)
}
