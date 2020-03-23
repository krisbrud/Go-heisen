package controller

import (
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"fmt"
)

func shouldStop(elev elevator.Elevator, activeOrders order.OrderList) bool {
	fmt.Printf("In shouldStop")
	elev.Print()
	activeOrders.Print()

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

func ordersAtCurrentFloor(elev elevator.Elevator, activeOrders []order.Order) bool {
	atCurrentFloor := func(o order.Order) bool {
		return o.Floor == elev.Floor && o.IsMine()
	}
	return anyOrder(activeOrders, atCurrentFloor)
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
