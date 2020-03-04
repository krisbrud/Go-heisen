package delegator

import (
	"Go-heisen/src/order"
)

type CostRequest struct {
	Order order.Order
}

type CostResponse struct {
	OrderID string
	Cost    int
}

type orderDelegation struct {
	o order.Order
	costs map[string] int
	disallowedRecipent string
}

func makeOrderDelegation(o order, disallowed)

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
				delegations[id] = 
			}

		case orderToRedelegate := <-toRedelegate:
			

		case costResponse := <-costResponseRx
			
		}
	}
}
