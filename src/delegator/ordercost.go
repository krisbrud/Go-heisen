package delegator

import (
	"Go-heisen/src/config"
	"Go-heisen/src/elevator"
)

func cost(order elevator.Order, state elevator.State) int {
	if !order.IsValid() || !state.IsValid() {
		return config.GetNumFloors() * 10 // Higher than any possible cost, but avoids overflow
	}

	if atDestinationFloor(order, state) {
		return 0
	}

	if !isTravellingTowardsOrder(order, state) && state.IntendedDir != elevator.MD_Stop {
		// The elevator needs to execute orders before turning around
		// => add distance before turning around and recursively find
		// distance from current state to intermediate + from intermediate state to destination
		intermediateState := getIntermediateState(state)
		return distance(state.Floor, intermediateState.Floor) + cost(order, intermediateState)
	}

	// Travelling towards order or standing still, return distance to floor
	return distance(order.Floor, state.Floor)
}

func distance(a, b int) int {
	// Absolute value of distance between floor a and b
	if a < b {
		return b - a
	}
	return a - b
}

func isTravellingTowardsOrder(order elevator.Order, state elevator.State) bool {
	switch {
	case order.Floor > state.Floor && state.IntendedDir == elevator.MD_Up:
		return true
	case order.Floor < state.Floor && state.IntendedDir == elevator.MD_Down:
		return true
	default:
		return false
	}
}

func atDestinationFloor(order elevator.Order, state elevator.State) bool {
	if order.Floor == state.Floor {
		switch {
		case order.IsFromCab():
			return true
		case order.Class == elevator.BT_HallUp && state.IntendedDir != elevator.MD_Down:
			return true
		case order.Class == elevator.BT_HallDown && state.IntendedDir != elevator.MD_Up:
			return true
		}
	}
	return false
}

func getIntermediateState(state elevator.State) elevator.State {
	// Helper function, make an intermediate state that will take place after changing travel direction
	// Assume that the elevator seeks to travel to the top/bottom floor (depending on direction) before
	// turning around.
	switch state.IntendedDir {
	case elevator.MD_Up:
		return elevator.State{
			Floor:       config.GetTopFloor(),
			IntendedDir: elevator.MD_Down,
			Behaviour:   elevator.EB_Idle,
			ElevatorID:  state.ElevatorID,
		}
	case elevator.MD_Down:
		return elevator.State{
			Floor:       config.GetBottomFloor(),
			IntendedDir: elevator.MD_Up,
			Behaviour:   elevator.EB_Idle,
			ElevatorID:  state.ElevatorID,
		}
	default:
		return elevator.MakeInvalidState()
	}
}
