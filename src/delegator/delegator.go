package delegator

import (
	"Go-heisen/src/config"
	"Go-heisen/src/elevator"
	"math"
	"time"
)

const (
	maxTimeSinceStateUpdate     = 5 * time.Second // Max time since we got an update from an elevator in order to delegate an order to it
	stateRedistributionInterval = 500 * time.Millisecond
)

// Delegator chooses the best recipent for a order to be delegated or redelegated
// based on it's current belief states
func Delegator(
	toDelegate chan elevator.Order,
	toRedelegate chan elevator.Order,
	toOrderTransmitter chan elevator.Order,
	toProcessor chan elevator.Order,
	transmitState chan elevator.State,
	receiveState chan elevator.State,
) {
	redelegations := make(map[elevator.OrderIDType]bool)
	elevatorStates := make(map[string]elevator.State)

	stateRedistributionTimer := time.NewTimer(stateRedistributionInterval)

	for {
		select {
		case orderToDelegate := <-toDelegate:
			// Find best recipent for order based on current belief state
			recipent := bestRecipent(orderToDelegate, elevatorStates, "")
			orderToDelegate.RecipentID = recipent

			toProcessor <- orderToDelegate
			toOrderTransmitter <- orderToDelegate

		case orderToRedelegate := <-toRedelegate:
			// Redelegate the order if it isn't redelegated already
			oldID := orderToRedelegate.OrderID
			if _, isAlreadyRedelegated := redelegations[oldID]; isAlreadyRedelegated {
				break // Don't redelegate the order, it already has been
			}
			// Set the order as being redelegated
			redelegations[oldID] = true

			disallowedRecipent := orderToRedelegate.RecipentID
			orderToRedelegate.OrderID = elevator.GetRandomID() // Give redelegation of order new ID

			recipent := bestRecipent(orderToRedelegate, elevatorStates, disallowedRecipent)
			orderToRedelegate.RecipentID = recipent

			toProcessor <- orderToRedelegate
			toOrderTransmitter <- orderToRedelegate

		case state := <-receiveState:
			if !state.IsValid() {
				break // Ignore invalid incoming state
			}

			// Make sure states are synced to local time.
			state.Timestamp = time.Now()

			elevatorStates[state.ElevatorID] = state

		case <-stateRedistributionTimer.C:
			// Redistribute the state regularly, to combat lost packets with state updates
			if state, ok := elevatorStates[config.GetMyElevatorID()]; ok {
				go func() { transmitState <- state }()
			}
			stateRedistributionTimer.Reset(stateRedistributionInterval)
		}
	}
}

func bestRecipent(order elevator.Order, states map[string]elevator.State, disallowed string) string {
	bestElevatorID := ""
	bestCost := math.MaxInt64

	for elevatorID, state := range states {
		stateCost := cost(order, state)

		// Check that the state update is recent enough
		if time.Since(state.Timestamp) > maxTimeSinceStateUpdate {
			continue // The state of this elevator is too old. Don't delegate to it.
		}
		if elevatorID != disallowed && stateCost < bestCost {
			bestCost = stateCost
			bestElevatorID = elevatorID
		}
	}

	if bestElevatorID == "" {
		// Did for some reason not find any valid recipents. Set self as best recipent.
		return config.GetMyElevatorID()
	}

	return bestElevatorID
}
