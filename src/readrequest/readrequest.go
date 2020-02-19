package readrequest

// ReadRequest serves as a request to read an order from OrderRepository
type ReadRequest struct {
	orderID string
	reader  string // TODO: Make own type for this
}
