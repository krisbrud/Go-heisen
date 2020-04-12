package controller

import (
	"Go-heisen/src/elevator"
	"fmt"
)

func shouldStop(state elevator.State, activeOrders elevator.OrderList) bool {
	fmt.Printf("In shouldStop")
	state.Print()
	activeOrders.Print()

	if len(activeOrders) == 0 {
		fmt.Println("No active orders, stopping")
		return true
	}

	for _, activeOrder := range activeOrders {
		// Check if the order is at our floor and not in the opposite direction
		if activeOrder.Floor == state.Floor {
			switch activeOrder.Class {
			case elevator.BT_HallUp:
				if state.IntendedDir == elevator.MD_Up {
					fmt.Println("ShouldStop found a floor to stop at while going up")
					return true
				}
			case elevator.BT_HallDown:
				if state.IntendedDir == elevator.MD_Down {
					fmt.Println("ShouldStop found a floor to stop at while going down")
					return true
				}
			case elevator.BT_Cab:
				fmt.Println("ShouldStop found a cab call at this floor and stopped")
				return true
			}
		}
	}
	
	// Handle the cases where we wish to stop at an order in the opposite direction
	// E.g. travellling up, an elevator should stop at a hall call going down if there are no orders above it
	if (state.IntendedDir == elevator.MD_Down && !ordersBelow(state, activeOrders)) 
		|| (state.IntendedDir == elevator.MD_Up && !ordersAbove(state, activeOrders)) {
		fmt.Println("ShouldStop found no orders in the direction of travel and stopped")
		return true
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
