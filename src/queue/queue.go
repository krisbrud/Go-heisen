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

type OrderQueue []order.Order

func MakeEmptyQueue() OrderQueue {
	return make([]order.Order, queueCapacity)
}

// SortQueue sorts the queue in-place by costs given state
func (q *OrderQueue) SortQueue(state elevatorstate.ElevatorState) {
	lessFunc := func(i, j int) bool {
		return ordercost.Cost((*q)[i], state) < ordercost.Cost((*q)[j], state)
	}
	sort.Slice(q, lessFunc)
}

// InsertSorted adds a order to the queue, keep it sorted by costs given state
func (q *OrderQueue) InsertSorted(o order.Order, state elevatorstate.ElevatorState) {
	*q = append(*q, o)
	q.SortQueue(state)
}

func (q OrderQueue) Peek() order.Order {
	return q[0]
}

func (q *OrderQueue) RemoveOrders(completedOrder order.Order) {
	newQueue := MakeEmptyQueue()

	for _, oldOrder := range *q {
		if !order.AreEquivalent(oldOrder, completedOrder) {
			// Keep the orders that the completedOrder does not fulfill
			newQueue = append(newQueue, oldOrder)
		}
	}
	// Assume queue is still sorted, replace old queue by the new one
	*q = newQueue
}
