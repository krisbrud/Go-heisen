package ordercost

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
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
	maxCost = elevatorstate.NumFloors * 10
)

func Cost(o order.Order, es elevatorstate.ElevatorState) int {
	if !o.IsValid() || !o.IsMine() || !es.IsValid() { // TODO - check what happens if removing ismine
		// TODO panic/restart
		return maxCost
	}

	if atDestinationFloor(o, es) {
		return 0
	}

	if !isTravellingTowardsOrder(o, es) {
		// We also need to execute orders before turning around
		// => add distance before turning around and recursively find
		// distance from current state to intermediate + from intermediate state to destination
		intermediateState := getIntermediateState(es)
		return distance(es.CurrentFloor, intermediateState.CurrentFloor) + Cost(o, intermediateState)
	}

	// Travelling towards order, return distance to floor
	return distance(o.Floor, es.CurrentFloor)
}

func distance(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func isTravellingTowardsOrder(o order.Order, es elevatorstate.ElevatorState) bool {
	switch {
	case o.Floor > es.CurrentFloor && es.IntendedDir == elevatorstate.Up:
		return true
	case o.Floor < es.CurrentFloor && es.IntendedDir == elevatorstate.Down:
		return true
	default:
		return false
	}
}

func atDestinationFloor(o order.Order, es elevatorstate.ElevatorState) bool {
	if o.Floor == es.CurrentFloor {
		switch {
		case o.IsFromCab():
			return true
		case o.Class == order.HALL_UP && es.IntendedDir != elevatorstate.Down:
			return true
		case o.Class == order.HALL_DOWN && es.IntendedDir != elevatorstate.Up:
			return true
		}
	}
	return false
}

func getIntermediateState(es elevatorstate.ElevatorState) elevatorstate.ElevatorState {
	// Gets the intermediate state that will take place after changing travel direction
	switch es.IntendedDir {
	case elevatorstate.Up:
		return elevatorstate.ElevatorState{
			CurrentFloor: elevatorstate.TopFloor,
			AtFloor:      true,
			IntendedDir:  elevatorstate.Down,
		}
	case elevatorstate.Down:
		return elevatorstate.ElevatorState{
			CurrentFloor: elevatorstate.BottomFloor,
			AtFloor:      true,
			IntendedDir:  elevatorstate.Down,
		}
	default:
		return elevatorstate.MakeInvalidState()
	}
}
