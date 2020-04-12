package delegator

import (
	"Go-heisen/src/elevator"
	"fmt"
)

// Delegator chooses the best recipent for a order to be delegated or redelegated
// based on it's current belief state
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

	// TODO stateupdates
	for {
		select {
		case orderToDelegate := <-toDelegate:
			// fmt.Println("Received order to delegate!")
			// orderToDelegate.Print()

			// Find best recipent for order based on current belief state
			recipent, err := bestRecipent(orderToDelegate, elevatorStates, "")
			orderToDelegate.RecipentID = recipent

			// Doing order myself, but warn user? TODO fix this
			if err != nil {
				fmt.Printf("%v\n", err)
				break
			}
			toProcessor <- orderToDelegate
			toOrderTransmitter <- orderToDelegate

		case orderToRedelegate := <-toRedelegate:
			// Redelegate the order if it isn't redelegated already
			oldID := orderToRedelegate.OrderID
			if _, alreadyRedelegating := redelegations[oldID]; alreadyRedelegating {
				break
			}
			// Set the order as being redelegated
			redelegations[oldID] = true

			disallowedRecipent := orderToRedelegate.RecipentID
			orderToRedelegate.OrderID = elevator.GetRandomID() // Give redelegation of order new ID

			recipent, err := bestRecipent(orderToRedelegate, elevatorStates, disallowedRecipent)
			orderToRedelegate.RecipentID = recipent
			if err != nil {
				fmt.Printf("%v\n", err)
				break
			}

			toProcessor <- orderToRedelegate
			toOrderTransmitter <- orderToRedelegate

		case state := <-receiveState:
			// fmt.Println("Received state from other elevator!")
			// state.Print()
			if !state.IsValid() {
				fmt.Printf("Invalid state incoming!")
				state.Print()
				break
			}

			// Notify other elevators about own state
			if state.ElevatorID == elevator.GetElevatorID() {
				oldElev, present := elevatorStates[state.ElevatorID]
				if present {
					if state != oldElev {
						transmitState <- state
					}
				} else {
					transmitState <- state
				}
			}
			// TODO: Possibly add timestamp for state, only accept states that are
			// recent enough. Then we may also get rid of the "peers variable" for simpler code.
			elevatorStates[state.ElevatorID] = state

			// DEBUG: Print all elevator states:
			//fmt.Println("All elevator states after delegator update", elevatorStates)
		}
	}
}

func bestRecipent(order elevator.Order, states map[string]elevator.State, disallowed string) (string, error) {
	bestElevatorID := ""
	bestCost := 10000 // TODO: Refactor

	// fmt.Printf("Finding best recipent for order %#v\n", order)
	// fmt.Printf("Disallowed: %v\n", disallowed)
	// fmt.Printf("All states: %#v\n", states)

	for elevatorID, state := range states {
		stateCost := cost(order, state)
		fmt.Printf("Cost for %v: %v", elevatorID, stateCost)
		if elevatorID != disallowed && stateCost < bestCost {
			bestCost = stateCost
			bestElevatorID = elevatorID
		}
	}
	fmt.Println("")

	if bestElevatorID == "" {
		err := fmt.Errorf("Did not find any valid elevator to delegate to! Order %#v\nStates: %#v\nDissallowed: %#v", order, states, disallowed)
		return elevator.GetElevatorID(), err
	}

	return bestElevatorID, nil
}
