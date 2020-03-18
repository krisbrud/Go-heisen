package elevator

import "fmt"

const (
	orderCapacity = 50
)

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	Floor       int
	IntendedDir MotorDirection
	Behaviour   ElevatorBehaviour
	ElevatorID  string
	// ActiveOrders []order.Order
}

func (elev Elevator) String() string {
	return fmt.Sprintf("Order:\n\tFloor: %v\n\tIntendedDir: %v\n\tBehaviour: %v\n\tElevatorID: %v\n")
}

var numFloors int = 4
var bottomFloor int = 0
var topFloor int = bottomFloor + numFloors - 1 // bottomFloor is a valid floor

func GetTopFloor() int    { return topFloor }
func GetBottomFloor() int { return bottomFloor }
func GetNumFloors() int   { return numFloors }

// IsValid tells us if both fields of Elevator are valid given the current configuration
func (elev Elevator) IsValid() bool {
	return bottomFloor <= elev.Floor && elev.Floor <= topFloor
}

func (elev Elevator) IsIdle() bool { return elev.Behaviour == EB_Idle }

func (dir MotorDirection) Opposite() MotorDirection {
	switch dir {
	case MD_Up:
		return MD_Down
	case MD_Down:
		return MD_Up
	default:
		return MD_Stop // TODO Maybe invalid?
	}
}

func (elev Elevator) IsDoorOpen() bool { return elev.Behaviour == EB_DoorOpen }

func UninitializedElevatorBetweenFloors() Elevator {
	return Elevator{
		Floor:       bottomFloor - 1,
		IntendedDir: MD_Down,
		Behaviour:   EB_Idle,
		ElevatorID:  GetMyElevatorID(),
	}
}

func MakeInvalidState() Elevator {
	return Elevator{
		Floor:       bottomFloor - 1,
		IntendedDir: MD_Stop,
		Behaviour:   EB_Idle,
	}
}

func GetMyElevatorID() string {
	return "My ElevatorID123" // TODO refactor - maybe a "config" module?
}
