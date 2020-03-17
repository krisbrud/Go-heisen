package queue

import (
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/order"
	"Go-heisen/src/ordercost"
	"sort"
)

const (
	queueCapacity = 100
)

type OrderQueue struct {
	orders []order.Order
}

func MakeEmptyQueue() OrderQueue {
	return OrderQueue{orders: make([]order.Order, queueCapacity)}
}

// SortQueue sorts the queue in-place by costs given state
func (q *OrderQueue) SortQueue(state elevatorstate.ElevatorState) {
	lessFunc := func(i, j) bool {
		return ordercost.Cost(q[i], state) < ordercost.Cost(q[j], state)
	}
	sort.Slice(q.orders, lessFunc)
}

// InsertSorted adds a order to the queue, keep it sorted by costs given state
func (q *OrderQueue) InsertSorted(o order.Order, state elevatorstate.ElevatorState) {
	q.orders = append(q.orders, o)
	q.SortQueue(state)
}

func (q OrderQueue) Peek() order.Order {
	return q.orders[0]
}

func (q *OrderQueue) RemoveOrders(completedOrder order.Order) {
	newQueue := MakeEmptyQueue()

	for _, oldOrder := range q.orders {
		if !equivalentOrders(oldOrder, completedOrder) {
			// Keep the orders that the completedOrder does not fulfill
			newQueue.orders = append(newQueue.orders, oldOrder)
		}
	}
	// Assume queue is still sorted, replace old queue by the new one
	*q = newQueue
}

func equivalentOrders(a, b order.Order) bool {
	return a.Class == b.Class && a.Floor == b.Floor
}

// TODO:
/*
	RemoveOrders
*/
