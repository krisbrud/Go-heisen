package controller

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/elevio"
	"Go-heisen/src/order"
)

func setAllLights(activeOrders order.OrderList) {
	// Make local representation to avoid briefly turning lights off before turning them on again
	numFloors := elevator.GetNumFloors()
	buttonsPerFloor := 3

	// indexed as lights[floor][ButtonType]
	lights := make([][]bool, numFloors, numFloors)
	for i := range lights {
		lights[i] = make([]bool, buttonsPerFloor, buttonsPerFloor)
	}

	for _, o := range activeOrders {
		if !o.Completed && !(o.IsFromCab() && !o.IsMine()) {
			// Found order that is not completed yet, and is not some other
			// elevators cab call. Set the light
			lights[o.Floor][int(o.Class)] = true
		}
	}

	// Iterate through all lights, set
	for floor := range lights {
		for buttonIdx := range lights[floor] {
			button := elevator.ButtonType(buttonIdx)
			if !order.ValidButtonTypeGivenFloor(button, floor) {
				continue
			}
			lightShouldBeOn := lights[floor][buttonIdx]
			elevio.SetButtonLamp(button, floor, lightShouldBeOn)
		}
	}
}
