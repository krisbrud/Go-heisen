package testutils

func GetSomeOrder() elevator.Order {
	return elevator.Order{
		1234,
		1,
		elevator.CAB,
		"Some recipent",
		false,
	}
}

func GetSomeOtherOrder() elevator.Order {
	return elevator.Order{
		4321,
		2,
		elevator.CAB,
		"Some other recipent",
		false,
	}
}
