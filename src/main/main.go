package main

import (
	"fmt"
	"time"

	"Go-heisen/src/arrivedfloorhandler"
	"Go-heisen/src/buttonpushhandler"
	"Go-heisen/src/controller"
	"Go-heisen/src/delegator"
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/orderprocessor"
	"Go-heisen/src/orderrepository"
	"Go-heisen/src/watchdog"

	"github.com/TTK4145/Network-go/network/bcast"
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
	arrivedStateUpdates := make(chan elevator.Elevator)
	// ButtonPushHandler
	buttonPushes := make(chan elevator.ButtonEvent)
	// Controller
	toController := make(chan order.Order)
	// Delegator
	localStateUpdates := make(chan elevator.Elevator)
	toDelegate := make(chan order.Order)
	toRedelegate := make(chan order.Order)
	// OrderRepository
	readSingleRequests := make(chan orderrepository.ReadRequest)
	readAllRequests := make(chan orderrepository.ReadRequest)
	writeRequests := make(chan orderrepository.WriteRequest)
	// OrderProcessor
	toOrderProcessor := make(chan order.Order)
	// Network
	transmitOrder := make(chan order.Order)
	transmitState := make(chan elevator.Elevator)
	receiveState := make(chan elevator.Elevator)

	orderPort := 44232
	go bcast.Transmitter(transmitOrder, orderPort)
	go bcast.Receiver(toOrderProcessor, orderPort)

	statePort := 44233
	go bcast.Transmitter(transmitState, statePort)
	go bcast.Receiver(receiveState, statePort)

	// Start goroutines
	go arrivedfloorhandler.ArrivedFloorHandler(arrivedStateUpdates, readSingleRequests, toOrderProcessor)
	go buttonpushhandler.ButtonPushHandler(buttonPushes, readAllRequests, toDelegate)
	go controller.Controller(toController, buttonPushes, localStateUpdates, arrivedStateUpdates)
	go delegator.Delegator(toDelegate, toRedelegate, transmitOrder, toOrderProcessor, localStateUpdates, transmitState, receiveState)
	go orderrepository.OrderRepository(readSingleRequests, readAllRequests, writeRequests)
	go orderprocessor.OrderProcessor(toOrderProcessor, readSingleRequests, writeRequests, toController, transmitOrder)
	go watchdog.Watchdog(readSingleRequests, toDelegate, transmitOrder)

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
