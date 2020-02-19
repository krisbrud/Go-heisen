package arrivedfloorhandler

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// ArrivedFloorHandler handles clearing of orders when arriving at a destination floor
func ArrivedFloorHandler(
	stateUpdates chan elevatorstate.ElevatorState,
	readReqChan chan readrequest.ReadRequest,
	fromOrderRepo chan order.Order,
	toOrderReceiver chan order.Order,
) {

}
