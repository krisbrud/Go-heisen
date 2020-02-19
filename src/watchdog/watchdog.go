package watchdog

import (
	"Go-heisen/src/order"
	"Go-heisen/src/readrequest"
)

// The watchdog regularly distributes the active orders in the system, and makes sure old orders are redelegated
func Watchdog(
	repoReads chan readrequest.ReadRequest,
	toDelegator chan order.Order,
	toTransmitter chan order.Order,
) {
	for {
		select {}
	}
}
