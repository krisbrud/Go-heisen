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

type ElevatorState struct {
	CurrentFloor int
	RelPos       RelativePosition
}
