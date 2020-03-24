package order

import (
	"Go-heisen/src/elevator"
	"fmt"
	"math/rand"
	"time"
)

type OrderIDType int

const (
	invalidOrderID  OrderIDType = -1
	invalidRecipent             = ""
)

type Order struct {
	OrderID    OrderIDType
	Floor      int
	Class      elevator.ButtonType // Defined by iota-"enum"
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
		fmt.Println("Orders: [")
		for _, o := range ol {
			o.Print()
		}
		fmt.Println("]")
	}
}

func NewInvalidOrder() Order {
	return Order{
		invalidOrderID,
		-1,
		elevator.ButtonType(-1),
		"no recipent",
		false,
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

func ValidButtonTypeGivenFloor(bt elevator.ButtonType, floor int) bool {
	switch bt {
	case elevator.BT_Cab:
		return elevator.GetBottomFloor() <= floor && floor <= elevator.GetTopFloor()
	case elevator.BT_HallDown:
		return elevator.GetBottomFloor()+1 <= floor && floor <= elevator.GetTopFloor()
	case elevator.BT_HallUp:
		return elevator.GetBottomFloor() <= floor && floor <= elevator.GetTopFloor()-1
	default:
		// Invalid ButtonType
		return false
	}
}

func (o Order) IsValid() bool {
	return ValidButtonTypeGivenFloor(o.Class, o.Floor)
}

func (o Order) IsMine() bool {
	return o.RecipentID == elevator.GetMyElevatorID() // TODO: Update this
}

func (o Order) IsFromHall() bool {
	return o.Class == elevator.BT_HallUp || o.Class == elevator.BT_HallDown
}

func (o Order) IsFromCab() bool {
	return o.Class == elevator.BT_Cab
}

// AreEquivalent returns true if orders have the same class, floor and completion status
func AreEquivalent(a, b Order) bool {
	switch a.Class {
	case elevator.BT_Cab:
		// Cab calls from different elevators are not equivalent.
		return a.Class == b.Class && a.Floor == b.Floor && a.Completed == b.Completed && a.RecipentID == b.RecipentID
	default:
		return a.Class == b.Class && a.Floor == b.Floor && a.Completed == b.Completed
	}
}
