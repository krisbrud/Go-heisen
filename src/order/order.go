package order

import (
	"math/rand"
	"time"
)

type OrderClass int

const (
	CAB OrderClass = iota
	HALL_UP
	HALL_DOWN
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
	return true // TODO: Update this
}

func (o Order) IsFromHall() bool {
	return o.Class == HALL_UP || o.Class == HALL_DOWN
}

func (o Order) IsFromCab() bool {
	return o.Class == CAB
}

// TODO:
// isMine() bool
