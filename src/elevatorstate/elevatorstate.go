package elevatorstate

// RelativePosition is an enum describing the position of an elevator relative to a floor
type RelativePosition int

// The possible relative positions
const (
	AtFloor    RelativePosition = iota
	OverFloor  RelativePosition = iota
	UnderFloor RelativePosition = iota
	Undefined  RelativePosition = iota
)

const (
	NumFloors   = 4
	BottomFloor = 0
	TopFloor    = BottomFloor + NumFloors - 1 // BottomFloor is a valid floor
)

type ElevatorState struct {
	CurrentFloor int
	RelPos       RelativePosition
}

// TODO: Allow any floors

// IsValid tells us if both fields of ElevatorState are valid given the current configuration
func (es ElevatorState) IsValid() bool {
	return BottomFloor <= es.CurrentFloor && es.CurrentFloor <= TopFloor && es.RelPos != Undefined
}

func (es ElevatorState) IsAtFloor() bool { return es.RelPos == AtFloor }
