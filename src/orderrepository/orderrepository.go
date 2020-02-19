package orderrepository

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// type RepoReader string // TODO: Implement

// OrderRepository is the single source of truth of all known orders in all nodes.
func OrderRepository(
	readRequests chan readrequest.ReadRequest,
	processorWrites chan order.Order,
	processorReads chan order.Order,
	buttonPushReads chan order.Order,
	arrivedFloorReads chan order.Order,
	watchdogReads chan order.Order,
) {
	for {
		select {}
	}
}
