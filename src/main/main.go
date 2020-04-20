package main

import (
	"Go-heisen/src/Network-go/network/bcast"
	"Go-heisen/src/controller"
	"Go-heisen/src/delegator"
	"Go-heisen/src/elevator"
	"Go-heisen/src/orderprocessor"
	"Go-heisen/src/watchdog"
	"runtime"
)

func main() {
	// Parse command line flags
	elevator.ParseConfigFlags()

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Declare channels, organized after who reads them
	// Controller
	activeOrders := make(chan []elevator.Order)
	// Delegator
	toDelegate := make(chan elevator.Order)
	toRedelegate := make(chan elevator.Order)
	receiveState := make(chan elevator.State)
	// OrderProcessor
	buttonPushes := make(chan elevator.ButtonEvent)
	floorArrivals := make(chan elevator.State)
	toOrderProcessor := make(chan elevator.Order)
	// Network
	transmitOrder := make(chan elevator.Order)
	transmitState := make(chan elevator.State)
	// Watchdog
	toWatchdog := make(chan []elevator.Order)

	orderPort := config.GetOrderPort()
	go bcast.Transmitter(orderPort, transmitOrder)
	go bcast.Receiver(orderPort, toOrderProcessor)

	statePort := config.GetStatePort()
	go bcast.Transmitter(statePort, transmitState)
	go bcast.Receiver(statePort, receiveState)

	// Start goroutines
	go controller.Controller(activeOrders, buttonPushes, receiveState, floorArrivals)
	go delegator.Delegator(toDelegate, toRedelegate, transmitOrder, toOrderProcessor, transmitState, receiveState)
	go orderprocessor.OrderProcessor(toOrderProcessor, buttonPushes, floorArrivals, activeOrders, toDelegate, toWatchdog, transmitOrder)
	go watchdog.Watchdog(toWatchdog, toRedelegate)

	// Block such that main goroutine does not exit
	select {}
}
