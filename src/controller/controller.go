package controller

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/elevio"
	"Go-heisen/src/order"
	"fmt"
	"math"
	"time"
)

const (
	doorDuration  = 3 * time.Second
	orderCapacity = 100
)

func Controller(
	incomingOrders chan order.Order,
	buttonPushes chan elevator.ButtonEvent,
	stateUpdates chan elevator.Elevator,
	toArrivedFloorHandler chan elevator.Elevator,
	elevatorPort int,
) {
	elevio.Init(fmt.Sprintf("localhost:%v", elevatorPort), 4)

	buttonUpdates := make(chan elevator.ButtonEvent)
	floorUpdates := make(chan int)

	go elevio.PollButtons(buttonUpdates)
	go elevio.PollFloorSensor(floorUpdates)

	// Initialize internal elevator state
	elev := elevator.UninitializedElevatorBetweenFloors()

	// Run elevator downwards if no state update
	select {
	case newFloor := <-floorUpdates:
		// Send the floor update again on the channel so the normal handler may do it's routine
		fmt.Println("Floor update received, sending state back!")
		go func() { floorUpdates <- newFloor }()
	case <-time.After(200 * time.Millisecond):
		// Elevator initialized between floors, go downwards.
		fmt.Println("Started between floors!")
		elev.IntendedDir = elevator.MD_Down
		elev.Behaviour = elevator.EB_Moving
		elevio.SetMotorDirection(elevator.MD_Down)
	}

	// Initialize timer for doors
	doorTimer := time.NewTimer(math.MaxInt64)
	doorTimer.Stop()

	activeOrders := make(order.OrderList, 0, orderCapacity)

	for {
		select {
		case buttonEvent := <-buttonUpdates:
			// Print state?
			fmt.Printf("Buttonevent: %#v\n", buttonEvent)
			elev.Print()
			activeOrders.Print()

			if !elev.IsValid() {
				continue
			}

			switch elev.Behaviour { // Cases are mutually exclusive
			case elevator.EB_DoorOpen:
				if elev.Floor == buttonEvent.Floor {
					doorTimer.Reset(doorDuration)
					// timer_start(elev.config.doorOpenDuration_s);
				} else {
					buttonPushes <- buttonEvent
				}

			case elevator.EB_Moving:
				buttonPushes <- buttonEvent

			case elevator.EB_Idle:
				if elev.Floor == buttonEvent.Floor {
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(doorDuration)
					// elev.Behaviour = EB_DoorOpen; // TODO refactor
				} else {
					buttonPushes <- buttonEvent
					nextDir := chooseDirection(elev, activeOrders)
					elev.IntendedDir = nextDir
					elevio.SetMotorDirection(elev.IntendedDir)
					elev.Behaviour = elevator.EB_Moving
				}
			}
			if elev.IsValid() {
				stateUpdates <- elev
			}

			setAllLights(activeOrders) // Set on incomingorder and completed order instead

			// printf("\nNew state:\n");
			// elevator_print(elev);

		case newFloor := <-floorUpdates:
			fmt.Printf("Floor update: %#v\n", newFloor)
			elev.Print()
			activeOrders.Print()

			elev.Floor = newFloor
			elevio.SetFloorIndicator(elev.Floor)

			if elev.Behaviour == elevator.EB_Moving && shouldStop(elev, activeOrders) {
				// Stop the elevator
				elevio.SetMotorDirection(elevator.MD_Stop)

				// Clear the orders we have fulfilled
				go func() { toArrivedFloorHandler <- elev }()

				// Open the door
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(doorDuration)
				elev.Behaviour = elevator.EB_DoorOpen

			}
			if elev.IsValid() {
				stateUpdates <- elev
			}

			setAllLights(activeOrders)
			fmt.Println("After floor update")
			elev.Print()
			activeOrders.Print()

		case <-doorTimer.C:
			// Door timer timed out, close door.
			elevio.SetDoorOpenLamp(false)

			// Find and set motor direction
			elev.IntendedDir = chooseDirection(elev, activeOrders)
			elevio.SetMotorDirection(elev.IntendedDir)

			// Set the Behaviour accordingly
			if elev.IntendedDir == elevator.MD_Stop {
				elev.Behaviour = elevator.EB_Idle
			} else {
				elev.Behaviour = elevator.EB_Moving
			}

			// Possibly print new state

		case newOrder := <-incomingOrders:
			fmt.Printf("\nNew order in controller: %v\nElevator: %v", newOrder.String(), elev.String())
			if !newOrder.IsValid() {
				fmt.Println("Controller received invalid order", newOrder)
			}

			if !newOrder.Completed {
				activeOrders = append(activeOrders, newOrder)

				// Choose direction and execute
				nextDir := chooseDirection(elev, activeOrders)
				fmt.Printf("Next intended direction: %v\n", nextDir)
				elev.IntendedDir = nextDir
				elevio.SetMotorDirection(elev.IntendedDir)
				elev.Behaviour = elevator.EB_Moving
			} else {
				// Remove from queue
				fmt.Println("Removing order from queue!")
				newOrder.Print()
				activeOrders = removeEquivalentOrders(activeOrders, newOrder)
			}

			fmt.Printf("\nNew order in controller handled: %s\nElevator: %s", newOrder, elev)
			activeOrders.Print()
			setAllLights(activeOrders)

		}
	}
}

func setAllLights(activeOrders order.OrderList) {
	// Make local representation to avoid briefly turning lights off before turning them on again
	if elevator.GetBottomFloor() != 0 {
		panic("routine setAllLights assumes the bottom floor is zero!")
	}

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

func removeEquivalentOrders(activeOrders order.OrderList, completedOrder order.Order) order.OrderList {
	fmt.Println("Remove equivalent orders")
	newActiveOrders := make(order.OrderList, 0, orderCapacity)

	for _, existingOrder := range activeOrders {
		if !order.AreEquivalent(existingOrder, completedOrder) {
			newActiveOrders = append(newActiveOrders, existingOrder)
		}
	}

	return newActiveOrders
}

func shouldStop(elev elevator.Elevator, activeOrders []order.Order) bool {
	switch elev.IntendedDir {
	case elevator.MD_Down:
		shouldStopAtOrder := func(o order.Order) bool {
			atSameFloor := elev.Floor == o.Floor
			notOppositeDirection := o.IsFromCab() || o.Class == elevator.BT_HallUp

			return atSameFloor && notOppositeDirection && o.IsMine()
		}
		return anyOrder(activeOrders, shouldStopAtOrder) || !ordersBelow(elev, activeOrders)

	case elevator.MD_Up:
		shouldStopAtOrder := func(o order.Order) bool {
			atSameFloor := elev.Floor == o.Floor
			notOppositeDirection := o.IsFromCab() || o.Class == elevator.BT_HallDown

			return atSameFloor && notOppositeDirection && o.IsMine()
		}
		return anyOrder(activeOrders, shouldStopAtOrder) || !ordersAbove(elev, activeOrders)
	}

	return true // Default
}

func anyOrder(orderList []order.Order, predicateFunc func(o order.Order) bool) bool {
	for _, o := range orderList {
		if predicateFunc(o) {
			return true
		}
	}
	return false
}

func ordersAbove(elev elevator.Elevator, activeOrders []order.Order) bool {
	if len(activeOrders) == 0 {
		return false
	}

	isAbove := func(o order.Order) bool {
		return o.Floor > elev.Floor && o.IsMine()
	}

	return anyOrder(activeOrders, isAbove)
}

func ordersBelow(elev elevator.Elevator, activeOrders []order.Order) bool {
	if len(activeOrders) == 0 {
		return false
	}

	isBelow := func(o order.Order) bool {
		return o.Floor < elev.Floor && o.IsMine()
	}

	return anyOrder(activeOrders, isBelow)
}

func chooseDirection(elev elevator.Elevator, activeOrders []order.Order) elevator.MotorDirection {
	switch elev.IntendedDir {
	case elevator.MD_Up:
		switch {
		case ordersAbove(elev, activeOrders):
			return elevator.MD_Up
		case ordersBelow(elev, activeOrders):
			return elevator.MD_Down
		}

	case elevator.MD_Down, elevator.MD_Stop:
		switch {
		case ordersBelow(elev, activeOrders):
			return elevator.MD_Down
		case ordersAbove(elev, activeOrders):
			return elevator.MD_Up
		}
	}
	return elevator.MD_Stop // Default case
}
