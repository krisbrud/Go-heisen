package controller

import (
	"Go-heisen/src/config"
	"Go-heisen/src/elevator"
	"Go-heisen/src/elevio"
)

func setAllLights(activeOrders []elevator.Order) {
	numFloors := config.GetNumFloors()
	buttonsPerFloor := 3

	// Make local representation to avoid briefly turning lights off before turning them on again
	// The slice of slices indexed as lights[floor][ButtonType]
	lights := make([][]bool, numFloors, numFloors)
	for i := range lights {
		lights[i] = make([]bool, buttonsPerFloor, buttonsPerFloor)
	}

	for _, order := range activeOrders {
		if !order.Completed && !(order.IsFromCab() && !order.IsMine()) {
			// Found order that is not completed yet, and is not some other
			// elevators cab call. Set the light
			lights[order.Floor][int(order.Class)] = true
		}
	}

	// Iterate through lights, set and clear as needed
	for floor := range lights {
		for buttonIdx := range lights[floor] {
			button := elevator.ButtonType(buttonIdx)
			if !elevator.ValidButtonTypeGivenFloor(button, floor) {
				continue
			}
			lightShouldBeOn := lights[floor][buttonIdx]
			elevio.SetButtonLamp(button, floor, lightShouldBeOn)
		}
	}
}
