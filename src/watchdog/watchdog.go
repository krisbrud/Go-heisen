package watchdog

import (
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"time"
)

const (
	millisBetweenTicks = 200
	timeOutDuration    = 40 * time.Second
)

// Watchdog regularly distributes the active orders in the system, and gives expired order to Delegator to be redelegated
func Watchdog(
	repoptr *orderrepository.OrderRepository,
	toRedelegate chan order.Order,
	toTransmitter chan order.Order,
) {
	// Initialize monotonic clock
	initTime := time.Now()
	getTimeSinceStartup := func() time.Duration {
		return time.Since(initTime)
	}

	timestamps := make(map[order.OrderIDType]time.Duration)

	ticker := time.NewTicker(millisBetweenTicks * time.Millisecond)

	for {
		select {
		case <-ticker.C: // New tick

			activeOrders := repoptr.ReadActiveOrders()
			for _, activeOrder := range activeOrders {
				// Check if order already has timestamp
				id := activeOrder.OrderID
				orderTimeStamp, alreadyRegistered := timestamps[id]
				now := getTimeSinceStartup()

				if alreadyRegistered {
					// Check if the order has timed out
					if isTimedOut(orderTimeStamp, now) {
						go func() {
							// Order has timed out, have delegator redelegate it
							toRedelegate <- activeOrder
						}()
					} else {
						// Static redundancy, broadcast the active order to other nodes
						go func() {
							toTransmitter <- activeOrder
						}()
					}
				} else {
					timestamps[id] = now
					// Static redundancy, broadcast the active order to other nodes
					go func() {
						toTransmitter <- activeOrder
					}()
				}

			}
		}
	}
}

func isTimedOut(timestamp time.Duration, now time.Duration) bool {
	return now-timestamp > timeOutDuration
}
