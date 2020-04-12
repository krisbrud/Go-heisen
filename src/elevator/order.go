package elevator

import (
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
	Class      ButtonType // Defined by iota-"enum"
	RecipentID string
	Completed  bool
}

func (order Order) String() string {
	return fmt.Sprintf("Order\n\tOrderID: %v\n\tFloor: %v\n\tClass: %v\n\tRecipentID: %v\n\tCompleted: %v\n", order.OrderID, order.Floor, order.Class, order.RecipentID, order.Completed)
}

func (order Order) Print() { fmt.Println(order.String()) }

type OrderList []Order

func (ol OrderList) Print() {
	if len(ol) == 0 {
		fmt.Println("Orders: []")
	} else {
		fmt.Println("Orders: [")
		for _, order := range ol {
			order.Print()
		}
		fmt.Println("]")
	}
}

func NewInvalidOrder() Order {
	return Order{
		invalidOrderID,
		-1,
		ButtonType(-1),
		"no recipent",
		false,
	}
}

var idGeneratorSeeded = false

func GetRandomID() OrderIDType {
	if !idGeneratorSeeded {
		rand.Seed(time.Now().UTC().UnixNano())
		idGeneratorSeeded = true
	}

	return OrderIDType(rand.Int())
}

func (order *Order) SetCompleted() { order.Completed = true }

func ValidButtonTypeGivenFloor(bt ButtonType, floor int) bool {
	switch bt {
	case BT_Cab:
		return GetBottomFloor() <= floor && floor <= GetTopFloor()
	case BT_HallDown:
		return GetBottomFloor()+1 <= floor && floor <= GetTopFloor()
	case BT_HallUp:
		return GetBottomFloor() <= floor && floor <= GetTopFloor()-1
	default:
		// Invalid ButtonType
		return false
	}
}

func (order Order) IsValid() bool {
	return ValidButtonTypeGivenFloor(order.Class, order.Floor)
}

func (order Order) IsMine() bool {
	return order.RecipentID == GetElevatorID()
}

func (order Order) IsFromHall() bool {
	return order.Class == BT_HallUp || order.Class == BT_HallDown
}

func (order Order) IsFromCab() bool {
	return order.Class == BT_Cab
}

// IsEquivalentWith returns true if orders have the same class, floor and completion status
func (a Order) IsEquivalentWith(b Order) bool {
	switch a.Class {
	case BT_Cab:
		// Cab calls from different elevators are not equivalent.
		return a.Class == b.Class && a.Floor == b.Floor && a.Completed == b.Completed && a.RecipentID == b.RecipentID
	default:
		return a.Class == b.Class && a.Floor == b.Floor && a.Completed == b.Completed
	}
}
