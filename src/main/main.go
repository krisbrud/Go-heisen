package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"Go-heisen/src/Network-go/network/bcast"
	"Go-heisen/src/controller"
	"Go-heisen/src/delegator"
	"Go-heisen/src/elevator"
	"Go-heisen/src/order"
	"Go-heisen/src/orderprocessor"
)

func main() {
	restartSystem := make(chan bool)

	var elevatorPort int = 15657
	var elevatorID string
	flag.IntVar(&elevatorPort, "port", 15657, "Port for connection to elevator")
	flag.StringVar(&elevatorID, "id", "elev"+strconv.Itoa(os.Getppid()), "ID of this elevator")
	flag.Parse()

	// Set ID of this elevator
	elevator.SetMyElevatorID(elevatorID)

	fmt.Printf("ElevatorPort %v\n", elevatorPort)

	go startSystem(restartSystem, elevatorPort)

	for {
		strconv.Itoa(os.Getppid())
		select {
		case <-restartSystem:
			go startSystem(restartSystem, elevatorPort)
		}
	}
}

func startSystem(restartSystem chan bool, elevatorPort int) {

	// Declare channels, organized after who reads them
	// ArrivedFloorHandler
	floorArrivals := make(chan elevator.Elevator)
	// ButtonPushHandler
	buttonPushes := make(chan elevator.ButtonEvent)
	// Controller
	activeOrders := make(chan order.OrderList)
	// Delegator
	toDelegate := make(chan order.Order)
	toRedelegate := make(chan order.Order)
	// OrderProcessor
	toOrderProcessor := make(chan order.Order)
	// Network
	transmitOrder := make(chan order.Order)
	transmitState := make(chan elevator.Elevator)
	receiveState := make(chan elevator.Elevator)

	orderPort := 44232
	go bcast.Transmitter(orderPort, transmitOrder)
	go bcast.Receiver(orderPort, toOrderProcessor)

	statePort := 44233
	go bcast.Transmitter(statePort, transmitState)
	go bcast.Receiver(statePort, receiveState)

	// Start goroutines
	go controller.Controller(activeOrders, buttonPushes, receiveState, floorArrivals, elevatorPort)
	go delegator.Delegator(toDelegate, toRedelegate, transmitOrder, toOrderProcessor, transmitState, receiveState)
	go orderprocessor.OrderProcessor(toOrderProcessor, buttonPushes, floorArrivals, activeOrders, toDelegate, transmitOrder)
	// go watchdog.Watchdog(readSingleRequests, toDelegate, transmitOrder)

	// Block such that goroutine does not exit
	<-restartSystem
	restartSystem <- true
}
