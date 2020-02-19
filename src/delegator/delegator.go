package delegator

import (
	"Go-heisen/src/order"
)

func Delegator(
	toDelegator chan order.Order,
	toTransmitter chan order.Order,
) {
	for {
		select {}
	}
}
