package orderprocessor

import (
	"Go-heisen/src/elevator"
	"fmt"
	"sync"
)

// OrderRepository is used to store orders in the system.
// It makes sure all orders written and read from it are valid.
type OrderRepository struct {
	orders map[elevator.OrderIDType]elevator.Order
	mtx    sync.Mutex
}

func makeEmptyOrderRepository() OrderRepository {
	return OrderRepository{
		orders: make(map[elevator.OrderIDType]elevator.Order),
		// Mutex mtx implicitly initialized as unlocked
	}
}

// readSingleOrder looks for a single order in the OrderRepository, and returns whether it was found
func (repoptr *OrderRepository) readSingleOrder(id elevator.OrderIDType) (elevator.Order, error) {
	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()
	order, found := repoptr.orders[id]

	var err error = nil
	if !found {
		order = elevator.NewInvalidOrder()
		err = fmt.Errorf("could not find order with id %v in OrderRepository", id)
	} else if !order.IsValid() {
		panic(fmt.Sprintf("invalid order %v inside OrderRepository", order.String()))
	}

	return order, err
}

// ReadActiveOrders returns a slice of all the orders in the OrderRepository which are not marked as completed
func (repoptr *OrderRepository) readActiveOrders() []elevator.Order {
	active := make([]elevator.Order, 0)

	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()
	// Iterate through all the orders, add the ones that are not completed
	for _, order := range repoptr.orders {
		if !order.Completed {
			if order.IsValid() {
				active = append(active, order)
			} else {
				panic(fmt.Sprintf("invalid order %v inside OrderRepository", order.String()))
			}
		}
	}

	return active
}

// WriteOrderToRepository writes the order to the OrderRepository, and panics if the order is invalid
func (repoptr *OrderRepository) writeOrderToRepository(order elevator.Order) {
	if !order.IsValid() {
		panic("trying to write invalid order %v to OrderRepository")
	}

	repoptr.mtx.Lock()
	defer repoptr.mtx.Unlock()
	repoptr.orders[order.OrderID] = order
}
