package order

type OrderClass int

const (
	CAB       OrderClass = 0 //iota
	HALL_UP   OrderClass = 1 //iota
	HALL_DOWN OrderClass = 2 //iota
)

type Order struct {
	OrderID    string
	Floor      int
	Class      OrderClass // Defined by iota-"enum"
	RecipentID string
	Completed  bool
}

// TODO:
// isMine() bool
// isCompleted() bool
