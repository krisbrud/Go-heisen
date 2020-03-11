package delegator

import (
	"Go-heisen/src/order"
)

type CostRequest struct {
	Order order.Order
}

type CostResponse struct {
	OrderID 	string
	ResponderID	string
	Cost    	int
}

type orderDelegation struct {
	o order.Order
	costs map[string] int
	disallowedRecipent string
}

func makeOrderDelegation(o order) {
	return orderDelegation{
		o: o,
		costs: [],
		disallowedRecipent: ""
	}
}

func makeOrderRedelegation(o order, disallowed string) {
	return orderDelegation{
		o: o,
		costs: [],
		disallowedRecipent: disallowed,
	}
}

func Delegator(
	toDelegate chan order.Order,
	toRedelegate chan order.Order,
	toTransmitter chan order.Order,
	costRequestTx chan CostRequest,
	costResponseRx chan CostResponse,
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
				delegate based on current costs

	*/

	delegations := make(map[string] orderDelegation)

	for {
		select {
		case orderToDelegate := <-toDelegate:
			id := orderToDelegate.OrderID
			if _, currentlyDelegating := delegations[id]; !currentlyDelegating {
				delegations[id] = makeOrderDelegation(orderToDelegate)
			}

		case orderToRedelegate := <-toRedelegate:
			id := orderToDelegate.OrderID
			disallowedRecipent := orderToRedelegate.RecipentID;
			if _, currentlyDelegating := delegations[id]; !currentlyDelegating {
				delegations[id] = makeOrderRedelegation(orderToDelegate, disallowedRecipent)
			}

		case costResponse := <-costResponseRx:
			orderID := costResponse.OrderID
			responderID := costResponse.ResponderID
			cost := costResponse.Cost

			delegation, currentlyDelegating := delegations[orderID]

			if !currentlyDelegating {
				break;
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

			// Check if there are enough costs

						
		}
	}
}
