package watchdog

import (
	"Go-heisen/src/elevator"
	"time"
)

const (
	timeOutDuration = 40 * time.Second
)

// Watchdog gives expired orders to Delegator to be redelegated
func Watchdog(
	activeOrdersUpdate chan []elevator.Order,
	toRedelegate chan elevator.Order,
) {
	timestamps := make(map[elevator.OrderIDType]time.Time)

	for {
		select {
		case activeOrders := <-activeOrdersUpdate:
			// Have orders redelegated if timed out
			for _, activeOrder := range activeOrders {
				// Check if order already has timestamp
				id := activeOrder.OrderID
				orderTimeStamp, alreadyRegistered := timestamps[id]
				now := time.Now()

				if alreadyRegistered {
					// Check if the order has timed out
					if isTimedOut(orderTimeStamp, now) {
						go func() {
							// Order has timed out, have delegator redelegate it
							toRedelegate <- activeOrder
						}()
					}
				} else {
					timestamps[id] = now
				}
			}
		}
	}
}

func isTimedOut(timestamp time.Time, now time.Time) bool {
	return now.Sub(timestamp) > timeOutDuration
}
