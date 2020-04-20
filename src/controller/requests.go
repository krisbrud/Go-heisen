package controller

import (
	"Go-heisen/src/elevator"
)

func shouldStop(state elevator.State, activeOrders []elevator.Order) bool {
	if len(activeOrders) == 0 {
		return true
	}

	for _, activeOrder := range activeOrders {
		// Check if the order is at our floor and not in the opposite direction
		if activeOrder.Floor == state.Floor {
			switch activeOrder.Class {
			case elevator.BT_HallUp:
				if state.IntendedDir == elevator.MD_Up {
					return true
				}
			case elevator.BT_HallDown:
				if state.IntendedDir == elevator.MD_Down {
					return true
				}
			case elevator.BT_Cab:
				return true
			}
		}
	}

	// Handle the cases where we wish to stop at an order in the opposite direction
	// E.g. travellling up, an elevator should stop at a hall call going down if there are no orders above it
	if (state.IntendedDir == elevator.MD_Down && !ordersBelow(state, activeOrders)) ||
		(state.IntendedDir == elevator.MD_Up && !ordersAbove(state, activeOrders)) {
		return true
	}

	return false // Default
}

// Generic helper function, returns true if predicateFunc returns true for any order in orderList, returns false otherwise
func anyOrder(orderList []elevator.Order, predicateFunc func(order elevator.Order) bool) bool {
	for _, order := range orderList {
		if predicateFunc(order) {
			return true
		}
	}
	return false
}

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
