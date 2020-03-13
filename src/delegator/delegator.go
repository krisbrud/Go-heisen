package delegator

import (
	"Go-heisen/src/order"
	"fmt"
)

type CostRequest struct {
	Order order.Order
}

type CostResponse struct {
	OrderID     int
	ResponderID string
	Cost        int
}

type orderDelegation struct {
	o                  order.Order
	costs              map[string]int
	disallowedRecipent string
}

func makeOrderDelegation(o order.Order) orderDelegation {
	return orderDelegation{
		o:                  o,
		costs:              make(map[string]int),
		disallowedRecipent: "",
	}
}

func makeOrderRedelegation(o order.Order, disallowed string) orderDelegation {
	return orderDelegation{
		o:                  o,
		costs:              make(map[string]int),
		disallowedRecipent: disallowed,
	}
}

func Delegator(
	toDelegate chan order.Order,
	toRedelegate chan order.Order,
	toTransmitter chan order.Order,
	toProcessor chan order.Order,
	costRequestTx chan CostRequest,
	costResponseRx chan CostResponse,
	peerUpdates chan []string,
) {
	currentlyRedelegating := make(map[string]bool)
	/*
		initialize delegations

		toDelegate:
			set order in currently delegating
			set costReqTimeout

		toRedelegate:
			if not currently redelegating
				set redelegation
				send to toDelegate

		costReq:
			check cost
			send costresponse

		costResponse:
			if currently delegating order
				if responder not in disallowed recipents for order
					update costs
					if enough costs
						set recipent, send to processor
						remove from currently delegating

		costTimeOut
			if currently delegating


	*/

	delegations := make(map[int]orderDelegation)
	redelegations := make(map[int]bool)
	peers := make([]string, 0)

	for {
		select {
		case orderToDelegate := <-toDelegate:
			id := orderToDelegate.OrderID
			if _, currentlyDelegating := delegations[id]; !currentlyDelegating {
				delegations[id] = makeOrderDelegation(orderToDelegate)
			}

		case orderToRedelegate := <-toRedelegate:
			// Redelegate the order if it isn't redelegated already
			if _, alreadyRedelegating := redelegations[orderToRedelegate.OrderID]; alreadyRedelegating {
				break
			}

			disallowedRecipent := orderToRedelegate.RecipentID
			orderToRedelegate.OrderID = order.GetRandomID()

			delegations[orderToRedelegate.OrderID] = makeOrderRedelegation(orderToRedelegate, disallowedRecipent)

		case costResponse := <-costResponseRx:
			orderID := costResponse.OrderID
			responderID := costResponse.ResponderID
			cost := costResponse.Cost

			delegation, currentlyDelegating := delegations[orderID]

			if !currentlyDelegating {
				break
			}
			if responderID == delegation.disallowedRecipent {
				break
			}

			// Everything seems fine, update the cost if it is not present
			// or worse than the one there
			if oldCost, costExists := delegation.costs[responderID]; costExists {
				if cost > oldCost {
					delegations[orderID].costs[responderID] = cost
				}
			} else {
				delegations[orderID].costs[responderID] = cost
			}

			if !enoughCosts(delegation, peers) {
				break
			}

			// Delegate the order to the elevator with the lowest cost which is still in peers
			bestElevatorID, err := lowestCost(delegation, peers)

			if err != nil {
				fmt.Printf(err.Error())
				break
			}

			delegation.o.RecipentID = bestElevatorID

			if !delegation.o.IsValid() {
				fmt.Println("Order to delegate was not valid!")
				break
			}

		case peerUpdate := <-peerUpdates:
			peers = peerUpdate
		}
	}
}

func lowestCost(delegation orderDelegation, peers []string) (string, error) {
	bestElevatorID := ""
	bestCost := 10000 // TODO: Refactor

	for elevatorID, cost := range delegation.costs {
		if elevatorID != delegation.disallowedRecipent && cost < bestCost {
			bestCost = cost
			bestElevatorID = elevatorID
		}
	}

	if bestElevatorID == "" {
		err := fmt.Errorf("Did not any valid elevator to delegate to! Delegation %#v\n Peers: %#v\n", delegation, peers)
		return "", err
	}

	return bestElevatorID, nil
}

func enoughCosts(delegation orderDelegation, peers []string) bool {
	numValidCosts := len(delegation.costs) // No need to check if disallowed, they are not added to the map
	numNeededCosts := len(peers)

	if delegation.disallowedRecipent != "" {
		// Redelegating order, one elevator disallowed
		numNeededCosts--
	}

	return numValidCosts >= numNeededCosts
}
