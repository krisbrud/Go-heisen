package elevator

// import (
// 	"fmt"
// 	"time"
// )

// /*
// func (eb ElevatorBehaviour) String() string {
// 	switch eb {
// 	case EB_Idle:
// 		return "Idle"
// 	case EB_DoorOpen:
// 		return "DoorOpen"
// 	case EB_Moving:
// 		return "Moving"
// 	default:
// 		return "invalid"
// 	}

// }
// */
// func (elev Elevator) String() string {
// 	return fmt.Sprintf("Elevator:\n\tFloor: %v\n\tIntendedDir: %v\n\tBehaviour: %v\n\tElevatorID: %v\n", elev.Floor, elev.IntendedDir, elev.Behaviour.String(), elev.ElevatorID)
// }

// func (elev Elevator) Print() {
// 	fmt.Printf(elev.String() + "\n")
// }

// var numFloors int = 4
// var bottomFloor int = 0
// var topFloor int = bottomFloor + numFloors - 1 // bottomFloor is a valid floor

// /*func GetTopFloor() int    { return topFloor }
// func GetBottomFloor() int { return bottomFloor }
// func GetNumFloors() int   { return numFloors }
// */
// // TODO setconfiguration?

// // IsValid tells us if both fields of Elevator are valid given the current configuration
// func (elev Elevator) IsValid() bool {
// 	return bottomFloor <= elev.Floor && elev.Floor <= topFloor
// }

// /*
// func (be ButtonEvent) IsValid() bool {
// 	return bottomFloor <= be.Floor && be.Floor <= topFloor &&
// 		(be.Button == BT_HallUp || be.Button == BT_HallDown || be.Button == BT_Cab)
// }
// */
// func (elev Elevator) IsIdle() bool { return elev.Behaviour == EB_Idle }

// /*
// func (dir MotorDirection) Opposite() MotorDirection {
// 	switch dir {
// 	case MD_Up:
// 		return MD_Down
// 	case MD_Down:
// 		return MD_Up
// 	default:
// 		return MD_Stop // TODO Maybe invalid?
// 	}
// }
// */
// func (elev Elevator) IsDoorOpen() bool { return elev.Behaviour == EB_DoorOpen }

// /*
// func UninitializedElevatorBetweenFloors() Elevator {
// 	return Elevator{
// 		Floor:       bottomFloor - 1,
// 		IntendedDir: MD_Down,
// 		Behaviour:   EB_Moving,
// 		ElevatorID:  GetMyElevatorID(),
// 		Timestamp:   time.Now(),
// 	}
// }

// func MakeInvalidState() Elevator {
// 	return Elevator{
// 		Floor:       bottomFloor - 1,
// 		IntendedDir: MD_Stop,
// 		Behaviour:   EB_Idle,
// 	}
// }
// */
// var myElevatorID string

// func SetMyElevatorID(id string) {
// 	myElevatorID = id
// 	fmt.Println("Set my ID to", id)
// }

// func GetMyElevatorID() string {
// 	if myElevatorID == "" {
// 		// ElevatorID not initialized, set to parents process ID
// 		panic("Trying to get elevator ID before setting it!")
// 	}
// 	return myElevatorID // TODO refactor - maybe a "config" module?
// }
