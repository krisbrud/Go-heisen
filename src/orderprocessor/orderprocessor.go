package orderprocessor

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
)

// OrderProcessor order from this or other elevators
func OrderProcessor(
	incomingOrdersChan chan order.Order,
	buttonPushes chan elevator.ButtonEvent,
	floorArrivals chan elevator.Elevator,
	toController chan order.Order,
	toDelegate chan order.Order,
	toTransmit chan order.Order,
) {
	allOrders := orderrepository.MakeEmptyOrderRepository()

	for {
		select {
		case incomingOrder := <-incomingOrdersChan:
			handleIncomingOrder(incomingOrder, &allOrders, toController, toDelegate, toTransmit)
		case buttonPush := <-buttonPushes:
			handleButtonPush(buttonPush, &allOrders, incomingOrdersChan, toDelegate)
		case elevAtFloor := <-floorArrivals:
			clearOrdersOnFloorArrival(elevAtFloor, &allOrders, incomingOrdersChan)
		}
	}
}

func handleIncomingOrder(
	incomingOrder order.Order,
	allOrders *orderrepository.OrderRepository,
	toController chan order.Order,
	toDelegate chan order.Order,
	toTransmit chan order.Order,
) {
	fmt.Printf("\nProcessor handling incoming order!\n")
	incomingOrder.Print()

	if !incomingOrder.IsValid() {
		fmt.Println("Incoming order not valid!")
		return // Ignore the incoming order
	}

	localOrder, err := allOrders.ReadSingleOrder(incomingOrder.OrderID)
	exists := err != nil

	if exists {
		fmt.Println("Order already exists!")
		switch {
		case localOrder.Completed && !incomingOrder.Completed:
			// Notify other nodes that order is actually completed.
			// Don't update the OrderRepository, local state is newer.
			go func() { toTransmit <- localOrder }()
		case !localOrder.Completed && incomingOrder.Completed:
			// Overwrite existing order as completed. Update controller.
			allOrders.WriteOrderToRepository(incomingOrder)
			fmt.Println("Order being marked as completed in processor.")

			go func() {
				toController <- incomingOrder
			}()
		}
	} else {
		// Incoming order is new. Register to OrderRepository, send to controller and transmitter.
		fmt.Println("New order incoming in processor")
		incomingOrder.Print()

		allOrders.WriteOrderToRepository(incomingOrder)
		go func() {
			toController <- incomingOrder
			toTransmit <- incomingOrder
		}()
	}
}

func clearOrdersOnFloorArrival(
	elev elevator.Elevator,
	repoptr *orderrepository.OrderRepository,
	handleOrder chan order.Order,
) {
	fmt.Printf("ArrivedFloorHandler! State: %#v", elev)

	if !elev.IsValid() {
		fmt.Println("New state not valid!")
		// TODO restart
	}

	// if elev.Behaviour == elevator.EB_Moving {
	// 	return // Elevator is moving, no orders to clear
	// }

	// Read all active orders from OrderRepository. Set the relevant ones as cleared.
	for _, activeOrder := range repoptr.ReadActiveOrders() {
		if activeOrder.Floor == elev.Floor {
			if activeOrder.IsFromHall() || (activeOrder.IsFromCab() && activeOrder.IsMine()) {
				// We have completed this order, make OrderProcessor register it and tell everyone.
				activeOrder.SetCompleted()
				go func() { handleOrder <- activeOrder }() // New goroutine to avoid deadlock
			}
		}
	}
}

//
func resendAllActiveOrders(
	repoptr *orderrepository.OrderRepository,
	toTransmit chan order.Order,
) {
	for _, activeOrder := range repoptr.ReadActiveOrders() {
		toTransmit <- activeOrder
	}
}

// handleButtonPush creates an order and sends it to be delegated if no equivalent order already exists.
func handleButtonPush(
	pushedButton elevator.ButtonEvent,
	repoptr *orderrepository.OrderRepository,
	incomingOrdersChan chan order.Order,
	toDelegate chan order.Order,
) {
	fmt.Println("handleButtonPush")
	if !pushedButton.IsValid() {
		return
	}

	o := order.MakeUnassignedOrder(pushedButton)

	if !repoptr.HasEquivalentOrders(o) {
		if o.Class == elevator.BT_Cab {
			// My cab call, assign to me
			o.RecipentID = elevator.GetMyElevatorID()
			go func() { incomingOrdersChan <- o }()
		} else {
			// No active orders are equivalent, have the new order delegated.
			go func() { toDelegate <- o }()
		}
	}
}
