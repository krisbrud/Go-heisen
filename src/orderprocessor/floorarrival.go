package orderprocessor

import (
	"Go-heisen/src/elevator"
	"fmt"
)

func clearOrdersOnFloorArrival(
	state elevator.State,
	activeOrders []elevator.Order,
	toOrderProcessor chan elevator.Order,
	transmitOrder chan elevator.Order,
) {
	fmt.Printf("ArrivedFloorHandler! State: %#v\n", state)

	if !state.IsValid() {
		panic("Invalid state in floor arrival handler")
	}

	// Read all active orders from OrderRepository. Set the relevant ones as cleared.
	for _, activeOrder := range activeOrders {
		if activeOrder.Floor == state.Floor {
			if activeOrder.IsFromHall() || (activeOrder.IsFromCab() && activeOrder.IsMine()) {
				fmt.Printf("Active order with floor %#v being set to complete\n", activeOrder.Floor)
				fmt.Println("Clearing order!")
				activeOrder.Print()
				// We have completed this order, make OrderProcessor register it and tell everyone.
				activeOrder.SetCompleted()
				go func() { toOrderProcessor <- activeOrder }()
				go func() { transmitOrder <- activeOrder }()
			}
		}
	}
}
