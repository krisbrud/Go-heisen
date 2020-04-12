package controller

import (
	"Go-heisen/src/elevator"
	"fmt"
)

func shouldStop(state elevator.State, activeOrders elevator.OrderList) bool {
	fmt.Printf("In shouldStop")
	state.Print()
	activeOrders.Print()
	if len(activeOrders) == 0 { //Checking whether there are any active orders. Stop if no
		fmt.Printf("No active orders, stopping\n")
		return true
	}

	for _, activeOrder := range activeOrders {
		if activeOrder.Floor == state.Floor && activeOrder.Class == elevator.BT_HallUp && state.IntendedDir == elevator.MD_Up {
			fmt.Printf("ShouldStop found a floor to stop atwhile going up\n")
			return true

		}
		if activeOrder.Floor == state.Floor && activeOrder.Class == elevator.BT_HallDown && state.IntendedDir == elevator.MD_Down {
			fmt.Printf("ShouldStop found a floor to stop at while going down\n")
			return true

		}
		if activeOrder.Class == elevator.BT_Cab && activeOrder.Floor == state.Floor {
			fmt.Printf("ShouldStop found a cab call at this floor and stopped\n")
			return true
		}
		if (state.IntendedDir == elevator.MD_Down && !ordersBelow(state, activeOrders)) || (state.IntendedDir == elevator.MD_Up && !ordersAbove(state, activeOrders)) {
			fmt.Printf("ShouldStop foud no orders below this one and stopped\n")
			return true
		}

	}
	return false // Default
}

// Generic helper function, returns true if predicateFunc returns true for any order in orderList
func anyOrder(orderList []elevator.Order, predicateFunc func(order elevator.Order) bool) bool {
	for _, order := range orderList {
		if predicateFunc(order) {
			return true
		}
	}
	return false
}

// TODO kanskje fjerne
func ordersAtCurrentFloor(state elevator.State, activeOrders []elevator.Order) bool {
	atCurrentFloor := func(order elevator.Order) bool {
		return order.Floor == state.Floor && order.IsMine()
	}
	return anyOrder(activeOrders, atCurrentFloor)
}

func ordersAbove(state elevator.State, activeOrders []elevator.Order) bool {
	if len(activeOrders) == 0 {
		return false
	}

	isAbove := func(order elevator.Order) bool {
		return order.Floor > state.Floor && order.IsMine()
	}

	return anyOrder(activeOrders, isAbove)
}

func ordersBelow(state elevator.State, activeOrders []elevator.Order) bool {
	if len(activeOrders) == 0 {
		return false
	}

	isBelow := func(order elevator.Order) bool {
		return order.Floor < state.Floor && order.IsMine()
	}

	return anyOrder(activeOrders, isBelow)
}

func chooseDirection(state elevator.State, activeOrders []elevator.Order) elevator.MotorDirection {
	switch state.IntendedDir {
	case elevator.MD_Up:
		switch {
		case ordersAbove(state, activeOrders):
			return elevator.MD_Up
		case ordersBelow(state, activeOrders):
			return elevator.MD_Down
		}

	case elevator.MD_Down, elevator.MD_Stop:
		switch {
		case ordersBelow(state, activeOrders):
			return elevator.MD_Down
		case ordersAbove(state, activeOrders):
			return elevator.MD_Up
		}
	}
	return elevator.MD_Stop // Default case
}
