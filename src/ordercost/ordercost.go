package ordercost

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"fmt"
)

/*
func cost(o, elev, intendedDir)
	if not travelling towards order
		tempState = state after travelling to top/bottom
		return distance to tempstate + cost(o, tempstate)
	else
		return distance to floor
*/

const (
	maxCost = 1000 // TODO find something clever to do here
)

func Cost(o order.Order, elev elevator.Elevator) int {
	if !o.IsValid() || !elev.IsValid() { // TODO - check what happens if removing ismine
		// TODO panic/restart
		return maxCost
	}

	if atDestinationFloor(o, elev) {
		return 0
	}

	if !isTravellingTowardsOrder(o, elev) && !elev.IsIdle() {
		fmt.Println("Inside recursive cost func")
		// We also need to execute orders before turning around
		// => add distance before turning around and recursively find
		// distance from current state to intermediate + from intermediate state to destination
		intermediateState := getIntermediateState(elev)
		return distance(elev.Floor, intermediateState.Floor) + Cost(o, intermediateState)
	}

	// Travelling towards order or standing still, return distance to floor
	fmt.Println("At end of cost func")
	return distance(o.Floor, elev.Floor)
}

func distance(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func isTravellingTowardsOrder(o order.Order, elev elevator.Elevator) bool {
	switch {
	case o.Floor > elev.Floor && elev.IntendedDir == elevator.MD_Up:
		return true
	case o.Floor < elev.Floor && elev.IntendedDir == elevator.MD_Down:
		return true
	default:
		return false
	}
}

func atDestinationFloor(o order.Order, elev elevator.Elevator) bool {
	if o.Floor == elev.Floor {
		switch {
		case o.IsFromCab():
			return true
		case o.Class == elevator.BT_HallUp && elev.IntendedDir != elevator.MD_Down:
			return true
		case o.Class == elevator.BT_HallDown && elev.IntendedDir != elevator.MD_Up:
			return true
		}
	}
	return false
}

func getIntermediateState(elev elevator.Elevator) elevator.Elevator {
	// Gets the intermediate state that will take place after changing travel direction
	switch elev.IntendedDir {
	case elevator.MD_Up:
		return elevator.Elevator{
			Floor:       elevator.GetTopFloor(),
			IntendedDir: elevator.MD_Down,
			Behaviour:   elevator.EB_Idle,
			ElevatorID:  elev.ElevatorID,
		}
	case elevator.MD_Down:
		return elevator.Elevator{
			Floor:       elevator.GetBottomFloor(),
			IntendedDir: elevator.MD_Down,
			Behaviour:   elevator.EB_Idle,
			ElevatorID:  elev.ElevatorID,
		}
	default:
		return elevator.MakeInvalidState()
	}
}
