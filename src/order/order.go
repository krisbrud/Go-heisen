package order

type OrderClass int

const (
	CAB OrderClass = iota
	HALL_UP
	HALL_DOWN
	INVALID
)

const (
	invalidOrderID  = ""
	invalidFloor    = -1
	invalidClass    = INVALID
	invalidRecipent = ""
)

type Order struct {
	OrderID    string
	Floor      int
	Class      OrderClass // Defined by iota-"enum"
	RecipentID string
	Completed  bool
}

func NewInvalidOrder() Order {
	return Order{
		invalidOrderID,
		-1,
		INVALID,
		"no recipent",
		false,
	}
}

func (o *Order) SetCompleted() { o.Completed = true }

func (o Order) IsValid() bool {
	return o.OrderID != invalidOrderID || o.Floor != invalidFloor || o.Class != invalidClass
}

// TODO:
// isMine() bool
