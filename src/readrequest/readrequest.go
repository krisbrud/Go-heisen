package readrequest

type ReadRequester int

// Possible readers of OrderRepository
const (
	OrderProcessor ReadRequester = iota
	ButtonPushHandler
	ArrivedFloorHandler
	Watchdog
)

// ReadRequest serves as a request to read an order from OrderRepository
type ReadRequest struct {
	OrderID string
	Reader  ReadRequester // TODO: Make own type for this
}
