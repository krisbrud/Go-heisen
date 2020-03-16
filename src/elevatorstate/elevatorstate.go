package elevatorstate

// RelativePosition is an enum describing the position of an elevator relative to a floor
type RelativePosition int

// The possible relative positions
const (
	AtFloor RelativePosition = iota
	OverFloor
	UnderFloor
	Undefined
)

const (
	NumFloors   = 4
	BottomFloor = 0
	TopFloor    = BottomFloor + NumFloors - 1 // BottomFloor is a valid floor
)

type Direction int

const (
	Up Direction = iota // TODO: get rid off
	Down
	Idle
)

// TODO: Refactor to elevator.State

type ElevatorState struct {
	CurrentFloor int
	AtFloor      bool
	IntendedDir  Direction
	ElevatorID   string
}

// TODO: Allow any floors

// IsValid tells us if both fields of ElevatorState are valid given the current configuration
func (es ElevatorState) IsValid() bool {
	return BottomFloor <= es.CurrentFloor && es.CurrentFloor <= TopFloor
}

func (es ElevatorState) IsAtFloor() bool { return es.AtFloor }

func (dir Direction) Opposite() Direction {
	switch dir {
	case Up:
		return Down
	case Down:
		return Up
	default:
		return Idle // TODO Maybe invalid?
	}
}

func MakeInvalidState() ElevatorState {
	return ElevatorState{BottomFloor - 1, false, Idle, ""}
}

func GetMyElevatorID() string {
	return "My ElevatorID123" // TODO refactor - maybe a "config" module?
}
