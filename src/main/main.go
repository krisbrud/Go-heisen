package main

import (
	"Go-heisen/src/arrivedfloorhandler"
	"Go-heisen/src/buttonpushhandler"
	"Go-heisen/src/controller"
	"Go-heisen/src/delegator"
	"Go-heisen/src/elevatorstate"
	"Go-heisen/src/networkreceiver"
	"Go-heisen/src/networktransmitter"
	"Go-heisen/src/order"
	"Go-heisen/src/orderrepository"
	"Go-heisen/src/watchdog"
	"fmt"
	"time"
)

func main() {
	restartSystem := make(chan bool)

	go startSystem(restartSystem)

	for {
		select {
		case <-restartSystem:
			go startSystem(restartSystem)
		}
	}
}

func startSystem(restartSystem chan bool) {

	// Declare channels, organized after who reads them
	restart := make(chan bool)

	// ArrivedFloorHandler
	arrivedStateUpdates := make(chan elevatorstate.ElevatorState)
	arrivedRepoReads := make(chan order.Order)
	// ButtonPushHandler
	buttonPushOrders := make(chan order.Order)
	buttonRepoReads := make(chan order.Order)
	// Controller - None yet!
	// Delegator
	toDelegator := make(chan order.Order)
	// OrderRepository
	repoReadRequests := make(chan orderrepository.ReadRequest)
	processorRepoWrites := make(chan order.Order)
	// OrderProcessor
	toOrderProcessor := make(chan order.Order)
	processorRepoReads := make(chan order.Order)
	// NetworkReceiver
	toTransmitter := make(chan order.Order)
	// NetworkTransmitter
	fromReceiver := make(chan order.Order)
	// Watchdog
	watchdogRepoReads := make(chan order.Order)

	// Start goroutines
	go arrivedfloorhandler.ArrivedFloorHandler(arrivedStateUpdates, repoReadRequests, arrivedRepoReads, toOrderProcessor)
	go buttonpushhandler.ButtonPushHandler(buttonPushOrders, buttonRepoReads, repoReadRequests, toDelegator)
	go controller.Controller()
	go delegator.Delegator(toDelegator, toTransmitter)
	go orderrepository.OrderRepository(repoReadRequests, processorRepoWrites, processorRepoReads, buttonRepoReads, arrivedRepoReads, watchdogRepoReads)
	go networkreceiver.NetworkReceiver(fromReceiver)
	go networktransmitter.NewtorkTransmitter(toTransmitter)
	go watchdog.Watchdog(repoReadRequests, toDelegator, toTransmitter)

	tick := time.Tick(1000 * time.Millisecond) // 1 second

	for {
		select {
		case <-tick:
			fmt.Println("Tick!") // Needed currently to prevent deadlock...
		case <-restart:
			// Something wrong happened, restart the system
			break
		}
	}

	restartSystem <- true
}
