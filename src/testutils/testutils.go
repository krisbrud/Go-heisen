package testutils

import "Go-heisen/src/order"

func GetSomeOrder() order.Order {
	return order.Order{
		1234,
		1,
		order.CAB,
		"Some recipent",
		false,
	}
}

func GetSomeOtherOrder() order.Order {
	return order.Order{
		4321,
		2,
		order.CAB,
		"Some other recipent",
		false,
	}
}
