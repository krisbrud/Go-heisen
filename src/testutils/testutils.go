package testutils

import "Go-heisen/src/order"

func GetSomeOrder() order.Order {
	return order.Order{
		"Some ID",
		1,
		order.CAB,
		"Some recipent",
		false,
	}
}

func GetSomeOtherOrder() order.Order {
	return order.Order{
		"Some other ID",
		2,
		order.CAB,
		"Some other recipent",
		false,
	}
}
