package elevator

import (
	"fmt"
)

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
	return fmt.Sprintf("Elevator:\n\tFloor: %v\n\tIntendedDir: %v\n\tBehaviour: %v\n\tElevatorID: %v\n", elev.Floor, elev.IntendedDir, elev.Behaviour, elev.ElevatorID)
}

func (elev Elevator) Print() {
	fmt.Printf(elev.String() + "\n")
}

var numFloors int = 4
var bottomFloor int = 0
var topFloor int = bottomFloor + numFloors - 1 // bottomFloor is a valid floor

func GetTopFloor() int    { return topFloor }
func GetBottomFloor() int { return bottomFloor }
func GetNumFloors() int   { return numFloors }

// TODO setconfiguration?

// IsValid tells us if both fields of Elevator are valid given the current configuration
func (elev Elevator) IsValid() bool {
	return bottomFloor <= elev.Floor && elev.Floor <= topFloor
}

func (be ButtonEvent) IsValid() bool {
	return bottomFloor <= be.Floor && be.Floor <= topFloor &&
		be.Button == BT_HallUp || be.Button == BT_HallDown || be.Button == BT_Cab
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
		Behaviour:   EB_Moving,
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

var myElevatorID string

func SetMyElevatorID(id string) {
	myElevatorID = id
	fmt.Println("Set my ID to", id)
}

func GetMyElevatorID() string {
	if myElevatorID == "" {
		// ElevatorID not initialized, set to parents process ID
		panic("Trying to get elevator ID before setting it!")
	}
	return myElevatorID // TODO refactor - maybe a "config" module?
}
