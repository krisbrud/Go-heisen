package delegator

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/ordercost"
	"fmt"
)

// Delegator chooses the best recipent for a order to be delegated or redelegated
// based on it's current belief state
func Delegator(
	toDelegate chan order.Order,
	toRedelegate chan order.Order,
	toOrderTransmitter chan order.Order,
	toProcessor chan order.Order,
	transmitState chan elevator.Elevator,
	receiveState chan elevator.Elevator,
) {
	redelegations := make(map[order.OrderIDType]bool)
	elevatorStates := make(map[string]elevator.Elevator)

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
			orderToRedelegate.OrderID = order.GetRandomID() // Give redelegation of order new ID

			recipent, err := bestRecipent(orderToRedelegate, elevatorStates, disallowedRecipent)
			orderToRedelegate.RecipentID = recipent
			if err != nil {
				fmt.Printf("%v\n", err)
				break
			}

			toProcessor <- orderToRedelegate
			toOrderTransmitter <- orderToRedelegate

		case elev := <-receiveState:
			// fmt.Println("Received state from other elevator!")
			// elev.Print()
			if !elev.IsValid() {
				fmt.Printf("Invalid elev incoming!")
				elev.Print()
				break
			}

			// Notify other elevators about own state
			if elev.ElevatorID == elevator.GetMyElevatorID() {
				oldElev, present := elevatorStates[elev.ElevatorID]
				if present {
					if elev != oldElev {
						transmitState <- elev
					}
				} else {
					transmitState <- elev
				}
			}
			// TODO: Possibly add timestamp for elev, only accept states that are
			// recent enough. Then we may also get rid of the "peers variable" for simpler code.
			elevatorStates[elev.ElevatorID] = elev

			// DEBUG: Print all elevator states:
			//fmt.Println("All elevator states after delegator update", elevatorStates)
		}
	}
}

func bestRecipent(o order.Order, states map[string]elevator.Elevator, disallowed string) (string, error) {
	bestElevatorID := ""
	bestCost := 10000 // TODO: Refactor

	// fmt.Printf("Finding best recipent for order %#v\n", o)
	// fmt.Printf("Disallowed: %v\n", disallowed)
	// fmt.Printf("All states: %#v\n", states)

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
		err := fmt.Errorf("Did not find any valid elevator to delegate to! Order %#v\nStates: %#v\nDissallowed: %#v", o, states, disallowed)
		return elevator.GetMyElevatorID(), err
	}

	return bestElevatorID, nil
}
