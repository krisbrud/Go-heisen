package orderrepositry

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// type RepoReader string // TODO: Implement

// OrderRepository is the single source of truth of all known orders in all nodes.
func OrderRepository(readRequestChan chan readrequest.ReadRequest,
	orderReceiverChan chan order.Order,
	buttonPushChan chan order.Order,
	arrivedFloorChan chan order.Order,
	watchdogChan chan order.Order) {
}
