package order

import (
	"Go-heisen/src/elevator"
	"fmt"
	"math/rand"
	"time"
)

type OrderClass int
type OrderIDType int

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
	OrderID    OrderIDType
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

func MakeUnassignedOrder(pushedButton elevator.ButtonEvent) Order {
	return Order{
		OrderID:    GetRandomID(),
		Floor:      pushedButton.Floor,
		Class:      OrderClass(pushedButton.Button), // TODO Verify that definitions are the same
		RecipentID: "",
		Completed:  false,
	}
}

var idGeneratorSeeded = false

func GetRandomID() OrderIDType {
	// TODO: Add mutex
	if !idGeneratorSeeded {
		rand.Seed(time.Now().UTC().UnixNano())
		idGeneratorSeeded = true
	}

	return OrderIDType(rand.Int())
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

// AreEquivalent returns true if orders have the same class, floor and completion status
func AreEquivalent(a, b Order) bool {
	return a.Class == b.Class && a.Floor == b.Floor && a.Completed == b.Completed
}
