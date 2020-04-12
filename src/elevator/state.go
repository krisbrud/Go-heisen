package elevator

import (
	"fmt"
	"time"
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

type State struct {
	Floor       int
	IntendedDir MotorDirection
	Behaviour   ElevatorBehaviour
	ElevatorID  string
	Timestamp   time.Time
}

func (eb ElevatorBehaviour) String() string {
	switch eb {
	case EB_Idle:
		return "Idle"
	case EB_DoorOpen:
		return "DoorOpen"
	case EB_Moving:
		return "Moving"
	default:
		return "invalid"
	}

}

func (state State) String() string {
	return fmt.Sprintf("State:\n\tFloor: %v\n\tIntendedDir: %v\n\tBehaviour: %v\n\tElevatorID: %v\n", state.Floor, state.IntendedDir, state.Behaviour.String(), state.ElevatorID)
}

func (state State) Print() {
	fmt.Printf(state.String() + "\n")
}

// IsValid tells us if both fields of State are valid given the current configuration
func (state State) IsValid() bool {
	return GetBottomFloor() <= state.Floor && state.Floor <= GetTopFloor()
}

func (be ButtonEvent) IsValid() bool {
	return GetBottomFloor() <= be.Floor && be.Floor <= GetTopFloor() &&
		(be.Button == BT_HallUp || be.Button == BT_HallDown || be.Button == BT_Cab)
}

func (state State) IsIdle() bool { return state.Behaviour == EB_Idle }

func (dir MotorDirection) Opposite() MotorDirection {
	switch dir {
	case MD_Up:
		return MD_Down
	case MD_Down:
		return MD_Up
	default:
		return MD_Stop
	}
}

func (state State) IsDoorOpen() bool { return state.Behaviour == EB_DoorOpen }

func UninitializedElevatorBetweenFloors() State {
	return State{
		Floor:       GetBottomFloor() - 1,
		IntendedDir: MD_Down,
		Behaviour:   EB_Moving,
		ElevatorID:  GetElevatorID(),
	}
}

func MakeInvalidState() State {
	return State{
		Floor:       GetBottomFloor() - 1,
		IntendedDir: MD_Stop,
		Behaviour:   EB_Idle,
	}
}
