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
	"Go-heisen/src/readrequest"
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
	/*
		TODO
		Declare channels
		Make restart-system (with channel)
		Start goroutines
	*/

	/*
		TODO Channels
	*/
	restart := make(chan bool)

	// Declare channels, organized after who reads them

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
	repoReadRequests := make(chan readrequest.ReadRequest)
	processorRepoWrites := make(chan readrequest.ReadRequest)
	// OrderProcessor
	toOrderProcessor := make(chan order.Order)
	processorRepoReads := make(chan order.Order)
	// NetworkReceiver
	toTransmitter := make(chan order.Order)
	// NetworkTransmitter
	fromReceiver := make(chan order.Order)
	// Watchdog

	// Start goroutines
	go arrivedfloorhandler.ArrivedFloorHandler(arrivedStateUpdates, repoReadRequests, arrivedRepoReads, toOrderProcessor)
	go buttonpushhandler.ButtonPushHandler(buttonPushOrders, buttonRepoReads, repoReadRequests, toDelegator)
	go controller.Controller()
	go delegator.Delegator(toDelegator, toTransmitter)
	go orderrepository.OrderRepository(repoReadRequests, processorRepoWrites, processorRepoReads, processorRepoWrites, buttonPushReads, arrivedRepoReads, watchdogRepoReads)
	go networkreceiver.NetworkReceiver(fromReceiver)
	go networktransmitter.NewtorkTransmitter(toTransmitter)

	restartSystem <- restart // Something wrong happened, restart the system
}
