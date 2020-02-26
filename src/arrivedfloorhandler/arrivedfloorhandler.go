package arrivedfloorhandler

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"fmt"
)

// ArrivedFloorHandler handles clearing of orders when arriving at a destination CurrentFloor
func ArrivedFloorHandler(
	stateUpdates chan elevatorstate.ElevatorState,
	readAllActiveRequests chan orderrepository.ReadRequest,
	toOrderProcessor chan order.Order,
) {
	for {
		select {
		case newState := <-stateUpdates:
			go func() {
				if !newState.IsValid() {
					fmt.Println("New state not valid!")
					// TODO restart
				}

				if !newState.IsAtFloor() {
					return // Not at floor, no orders to clear
				}

				// Read all active orders from OrderRepository. Set the relevant ones as cleared.
				readAllReq := orderrepository.MakeReadAllActiveRequest()
				readAllActiveRequests <- readAllReq

				for activeOrder := range readAllReq.ResponseCh {
					if activeOrder.Floor == newState.CurrentFloor {
						if activeOrder.IsFromHall() || (activeOrder.IsFromCab() && activeOrder.IsMine()) {
							// We have completed this order, make OrderProcessor register it and tell everyone.
							activeOrder.SetCompleted()
							toOrderProcessor <- activeOrder
						}
					}
				}
			}()
		}
	}
}
