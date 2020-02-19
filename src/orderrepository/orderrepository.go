package orderrepositry

import "order"

// type RepoReader string // TODO: Implement

// ReadRequest serves as a request to read an order
type ReadRequest struct {
	orderID string
	reader  string // TODO: Make own type for this
}

// OrderRepository is the single source of truth of all known orders in all nodes.
func OrderRepository(readRequestChan chan ReadRequest,
	buttonPushChan chan order.Order) {

}
