package delegator

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/ordercost"
	"fmt"
)

// Delegator chooses the best recipent for a order to be delegated or redelegated
// based on it's current belief state
func Delegator(
	toDelegate chan order.Order,
	toRedelegate chan order.Order,
	toTransmitter chan order.Order,
	toProcessor chan order.Order,
	stateUpdates chan elevatorstate.ElevatorState,
	//peerUpdates chan []string,
) {

	redelegations := make(map[int]bool)
	// peers := make([]string, 0)
	elevatorStates := make(map[string]elevatorstate.ElevatorState)

	for {
		select {
		case orderToDelegate := <-toDelegate:
			// Find best recipent for order based on current belief state
			recipent, err := bestRecipent(orderToDelegate, elevatorStates, "")
			orderToDelegate.RecipentID = recipent
			// Doing order myself, but warn user? TODO fix this
			if err != nil {
				fmt.Printf("%v\n", err)
				break
			}
			toProcessor <- orderToDelegate
			toTransmitter <- orderToDelegate

		case orderToRedelegate := <-toRedelegate:
			// Redelegate the order if it isn't redelegated already
			oldID := orderToRedelegate.OrderID
			if _, alreadyRedelegating := redelegations[oldID]; alreadyRedelegating {
				break
			}
			// Set the order as being redelegated
			redelegations[oldID] = true

			disallowedRecipent := orderToRedelegate.RecipentID
			orderToRedelegate.OrderID = order.GetRandomID() // Give redelegation of order new ID

			recipent, err := bestRecipent(orderToRedelegate, elevatorStates, disallowedRecipent)
			orderToRedelegate.RecipentID = recipent
			if err != nil {
				fmt.Printf("%v\n", err)
				break
			}

			toProcessor <- orderToRedelegate
			toTransmitter <- orderToRedelegate

		case state := <-stateUpdates:
			if !state.IsValid() {
				fmt.Printf("Invalid state incoming! state: %#v", state)
				break
			}
			// TODO: Possibly add timestamp for state, only accept states that are
			// recent enough. Then we may also get rid of the "peers variable" for simpler code.
			elevatorStates[state.ElevatorID] = state

			// case peerUpdate := <-peerUpdates:
			// 	peers = peerUpdate
		}
	}
}

func bestRecipent(o order.Order, states map[string]elevatorstate.ElevatorState, disallowed string) (string, error) {
	bestElevatorID := ""
	bestCost := 10000 // TODO: Refactor

	fmt.Printf("Finding best recipent for order %#v\n", o)
	fmt.Printf("Disallowed: %v\n", disallowed)
	for elevatorID, state := range states {
		cost := ordercost.Cost(o, state)
		fmt.Printf("Cost for %v: %v", elevatorID, cost)
		if elevatorID != disallowed && cost < bestCost {
			bestCost = cost
			bestElevatorID = elevatorID
		}
	}
	fmt.Println("")

	if bestElevatorID == "" {
		err := fmt.Errorf("Did not any valid elevator to delegate to! Order %#v\nStates: %#v\nDissallowed: %#v", o, states, disallowed)
		return elevatorstate.GetMyElevatorID(), err
	}

	return bestElevatorID, nil
}
