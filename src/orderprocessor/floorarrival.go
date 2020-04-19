package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/orderrepository"
	"fmt"
	"time"
)

func clearOrdersOnFloorArrival(
	state elevator.State,
	repoptr *orderrepository.OrderRepository,
	allOrders *orderrepository.OrderRepository,
	toController chan []elevator.Order,
	transmitOrder chan elevator.Order,
) {
	fmt.Printf("ArrivedFloorHandler! State: %#v\n", state)

	if !state.IsValid() {
		panic("Invalid state in floor arrival handler")
	}

	// Read all active orders from OrderRepository. Set the relevant ones as cleared.
	for _, activeOrder := range repoptr.ReadActiveOrders() {
		if activeOrder.Floor == state.Floor && state.IsDoorOpen() {
			fmt.Printf("Active order with floor %#v being set to complete\n", activeOrder.Floor)
			if activeOrder.IsFromHall() || (activeOrder.IsFromCab() && activeOrder.IsMine()) {
				fmt.Println("Clearing order!")
				activeOrder.Print()
				// We have completed this order, make OrderProcessor register it and tell everyone.
				activeOrder.SetCompleted()
				allOrders.WriteOrderToRepository(activeOrder)
				go func() { transmitOrder <- activeOrder }()

				activeOrders := allOrders.ReadActiveOrders()
				time.Sleep(10 * time.Millisecond) //this fixes a race condition causing lights on elevators to not properly switch off
				go func() { toController <- activeOrders }()

			}
		}
	}
}
