package arrivedfloorhandler

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
)

// ArrivedFloorHandler handles clearing of orders when arriving at a destination floor
func ArrivedFloorHandler(
	stateUpdates chan elevatorstate.ElevatorState,
	repoReadRequests chan orderrepository.ReadRequest,
	fromOrderRepo chan order.Order,
	toOrderProcessor chan order.Order,
) {
	// Remember to differentiate between cab and hall orders when clearing!

	for {
		select {}
	}
}
