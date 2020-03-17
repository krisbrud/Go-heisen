package ordercost

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"fmt"
)

/*
func cost(o, es, intendedDir)
	if not travelling towards order
		tempState = state after travelling to top/bottom
		return distance to tempstate + cost(o, tempstate)
	else
		return distance to floor
*/

const (
	maxCost = 1000 // TODO find something clever to do here
)

func Cost(o order.Order, es elevator.Elevator) int {
	if !o.IsValid() || !o.IsMine() || !es.IsValid() { // TODO - check what happens if removing ismine
		// TODO panic/restart
		return maxCost
	}

	if atDestinationFloor(o, es) {
		return 0
	}

	if !isTravellingTowardsOrder(o, es) && !es.IsIdle() {
		fmt.Println("Inside recursive cost func")
		// We also need to execute orders before turning around
		// => add distance before turning around and recursively find
		// distance from current state to intermediate + from intermediate state to destination
		intermediateState := getIntermediateState(es)
		return distance(es.Floor, intermediateState.Floor) + Cost(o, intermediateState)
	}

	// Travelling towards order or standing still, return distance to floor
	fmt.Println("At end of cost func")
	return distance(o.Floor, es.Floor)
}

func distance(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func isTravellingTowardsOrder(o order.Order, es elevator.Elevator) bool {
	switch {
	case o.Floor > es.Floor && es.IntendedDir == elevator.MD_Up:
		return true
	case o.Floor < es.Floor && es.IntendedDir == elevator.MD_Down:
		return true
	default:
		return false
	}
}

func atDestinationFloor(o order.Order, es elevator.Elevator) bool {
	if o.Floor == es.Floor {
		switch {
		case o.IsFromCab():
			return true
		case o.Class == order.HALL_UP && es.IntendedDir != elevator.MD_Down:
			return true
		case o.Class == order.HALL_DOWN && es.IntendedDir != elevator.MD_Up:
			return true
		}
	}
	return false
}

func getIntermediateState(es elevator.Elevator) elevator.Elevator {
	// Gets the intermediate state that will take place after changing travel direction
	switch es.IntendedDir {
	case elevator.MD_Up:
		return elevator.Elevator{
			Floor:       elevator.GetTopFloor(),
			IntendedDir: elevator.MD_Down,
			Behaviour:   elevator.EB_Idle,
			ElevatorID:  es.ElevatorID,
		}
	case elevator.MD_Down:
		return elevator.Elevator{
			Floor:       elevator.GetBottomFloor(),
			IntendedDir: elevator.MD_Down,
			Behaviour:   elevator.EB_Idle,
			ElevatorID:  es.ElevatorID,
		}
	default:
		return elevator.MakeInvalidState()
	}
}
