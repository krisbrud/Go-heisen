package orderprocessor

import (
	"Go-heisen/src/elevator"
)

func clearOrdersOnFloorArrival(
	state elevator.State,
	activeOrders []elevator.Order,
	toOrderProcessor chan elevator.Order,
	transmitOrder chan elevator.Order,
) {
	if !state.IsValid() {
		panic("Invalid state in floor arrival handler")
	}

	// Read all active orders from OrderRepository. Set the relevant ones as cleared.
	for _, activeOrder := range activeOrders {
		if activeOrder.Floor == state.Floor {
			if activeOrder.IsFromHall() || (activeOrder.IsFromCab() && activeOrder.IsMine()) {
				// We have completed this order, make OrderProcessor register it and tell everyone.
				activeOrder.SetCompleted()
				go func() { toOrderProcessor <- activeOrder }()
				go func() { transmitOrder <- activeOrder }()
			}
		}
	}
}
