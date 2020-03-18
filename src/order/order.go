package order

import (
	"Go-heisen/src/elevator"
	"fmt"
	"math/rand"
	"time"
)

type OrderClass int

// TODO: elevio or this so there are no double definitions!
const (
	HALL_UP OrderClass = iota
	HALL_DOWN
	CAB
	INVALID
)

const (
	invalidOrderID  = -1
	invalidFloor    = -1
	invalidClass    = INVALID
	invalidRecipent = ""
)

type Order struct {
	OrderID    int
	Floor      int
	Class      OrderClass // Defined by iota-"enum"
	RecipentID string
	Completed  bool
}

func (o Order) String() string {
	return fmt.Sprintf("Order\n\tOrderID: %v\n\tFloor: %v\n\tClass: %v\n\tRecipentID: %v\n\tCompleted: %v\n", o.OrderID, o.Floor, o.Class, o.RecipentID, o.Completed)
}

func (o Order) Print() { fmt.Println(o.String()) }

type OrderList []Order

func (ol OrderList) Print() {
	if len(ol) == 0 {
		fmt.Println("Orders: []")
	} else {
		fmt.Println("Orders:")
		for _, o := range ol {
			o.Print()
		}
	}
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

var idGeneratorSeeded = false

func GetRandomID() int {
	// TODO: Add mutex
	if !idGeneratorSeeded {
		rand.Seed(time.Now().UTC().UnixNano())
		idGeneratorSeeded = true
	}

	return rand.Int()
}

func (o *Order) SetCompleted() { o.Completed = true }

func (o Order) IsValid() bool {
	return o.OrderID != invalidOrderID || o.Floor != invalidFloor || o.Class != invalidClass
}

func (o Order) IsMine() bool {
	return o.RecipentID == elevator.GetMyElevatorID() // TODO: Update this
}

func (o Order) IsFromHall() bool {
	return o.Class == HALL_UP || o.Class == HALL_DOWN
}

func (o Order) IsFromCab() bool {
	return o.Class == CAB
}

func AreEquivalent(a, b Order) bool {
	return a.Class == b.Class && a.Floor == b.Floor && a.Completed == b.Completed
}

// TODO:
// isMine() bool
