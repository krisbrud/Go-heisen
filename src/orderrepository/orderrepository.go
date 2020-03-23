package orderrepository

import (
	"Go-heisen/src/order"
	"fmt"
	"sync"
)

const (
	defaultCapacity = 100
)

type OrderRepository struct {
	orders map[order.OrderIDType]order.Order
	mtx    sync.Mutex
}

func MakeEmptyOrderRepository() OrderRepository {
	return OrderRepository{
		orders: make(map[order.OrderIDType]order.Order),
		// Mutex mtx implicitly initialized
	}
}

// ReadSingleOrder looks for a single order in the OrderRepository, and returns an error if it isn't found
func (repoptr *OrderRepository) ReadSingleOrder(id order.OrderIDType) (order.Order, error) {
	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()
	o, found := repoptr.orders[id]

	var err error = nil
	if found {
		o = order.NewInvalidOrder()
		err = fmt.Errorf("could not find order with id %v in OrderRepository", id)
	} else if !o.IsValid() {
		panic(fmt.Sprintf("invalid order %v inside OrderRepository", o.String()))
	}

	return o, err
}

// ReadActiveOrders returns a slice of all the orders in the OrderRepository which are not marked as completed
func (repoptr *OrderRepository) ReadActiveOrders() order.OrderList {
	active := make(order.OrderList, 0)

	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()
	// Iterate through all the orders, add the ones that are not completed
	for _, o := range repoptr.orders {
		if !o.Completed {
			if o.IsValid() {
				active = append(active, o)
			} else {
				panic(fmt.Sprintf("invalid order %v inside OrderRepository", o.String()))
			}
		}
	}

	return active
}

// WriteOrderToRepository writes the order to the OrderRepository, and panics if the order is invalid
func (repoptr *OrderRepository) WriteOrderToRepository(o order.Order) {
	if !o.IsValid() {
		panic("trying to write invalid order %v to OrderRepository")
	}

	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()
	repoptr.orders[o.OrderID] = o
}

// HasEquivalentOrders returns true if the OrderRepository has at least one order that is equivalent
func (repoptr *OrderRepository) HasEquivalentOrders(o order.Order) bool {
	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()

	for _, storedOrder := range repoptr.orders {
		if order.AreEquivalent(o, storedOrder) {
			return true
		}
	}
	return false
}
