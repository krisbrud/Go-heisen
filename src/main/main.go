package main

import (
	//"Go-heisen/src/controller"
	"fmt"

	"../controller"
	"../elevatorio"
	"../elevatorstate"
	"../order"
	//"Go-heisen/src/readrequest"
)

func main() {

	/*	restartSystem := make(chan bool)

		go startSystem(restartSystem)
S
		for {
			select {
			case <-restartSystem:
				go startSystem(restartSystem)
			}
		}*/
	toButtonPushHandler := make(chan order.Order)
	toArrivedFloorHandler := make(chan elevatorstate.ElevatorState)
	toDelegator := make(chan elevatorstate.ElevatorState)
	readState := make(chan elevatorstate.ElevatorState)
	readQueue := make(chan order.Order)
	readButtonPush := make(chan order.Order)
	readCurrentFloor := make(chan elevatorio.ButtonEvent)

	go controller.Controller(toButtonPushHandler,
		toArrivedFloorHandler,
		toDelegator,
		readState,
		readQueue,
		readButtonPush)

	for {
		select {
		case btn := <-toButtonPushHandler:
			fmt.Println(btn)
		}

	}

}

//func startSystem(restartSystem chan bool) {
/*
	TODO
	Declare channels
	Make restart-system (with channel)
	Start goroutines
*/

/*
	TODO Channels
*/
//restart := make(chan bool)

// Declare channels, organized after who reads them

// ArrivedFloorHandler
/*
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
*/
// Start goroutines
/*
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
*/
