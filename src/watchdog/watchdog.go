package watchdog

import "Go-heisen/src/order"

// The watchdog regularly distributes the active orders in the system, and makes sure old orders are redelegated
func Watchdog(
	repoReads chan order.Order,
	toDelegator chan order.Order,
	toTransmitter chan order.Order,
) {
	for {
		select {}
	}
}
