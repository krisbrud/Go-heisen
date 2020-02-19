package order

type OrderClass int

const (
	CAB       OrderClass = iota
	HALL_UP   OrderClass = iota
	HALL_DOWN OrderClass = iota
)

type Order struct {
	orderID    string
	floor      int
	class      OrderClass // Defined by iota-"enum"
	recipentID string
	completed  bool
}

// TODO:
// isMine() bool
// isCompleted() bool
