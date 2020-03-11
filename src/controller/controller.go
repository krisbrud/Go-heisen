package controller

import (
	//"time"
	"fmt"

	"../elevatorio"
	"../elevatorstate"
	"../order"
	//"github.com/kjk/betterguid"
)

//type OrderClass int
//const (
//	CAB       OrderClass = iota
//	HALL_UP   OrderClass = iota
//	HALL_DOWN OrderClass = iota
//)

// type Order struct {
// 	orderID    string
// 	floor      int
// 	class      OrderClass // Defined by iota-"enum"
// 	recipentID string
// 	completed  bool
// }

/*Controller blabla */
func Controller(toButtonPushHandler chan order.Order,
	toArrivedFloorHandler chan elevatorstate.ElevatorState,
	toDelegator chan elevatorstate.ElevatorState,
	readState chan elevatorstate.ElevatorState,
	readQueue chan order.Order,
	readButtonPush chan elevatorio.ButtonEvent) {

	var state elevatorstate.ElevatorState
	var floor int
	var destination int
	var newOrder order.Order

	for {
		select {

		case btn := <-readButtonPush:

			orderClass := order.OrderClass(btn.Button)
			newOrder = order.Order{OrderID: "dette er en test", Floor: btn.Floor, Class: orderClass, RecipentID: "", Completed: false}

			//newOrder := order.Order{}
			go func() { toButtonPushHandler <- newOrder }() // må sende dette til button pushed handler

		case s := <-readState:
			if s != state {
				state = s
				go func() { toDelegator <- s }() // må lese elevator state og sende den til delegatoren

			}
			if s.CurrentFloor == destination {
				elevatorio.SetMotorDirection(elevatorio.MD_Stop)
				go func() { toArrivedFloorHandler <- state }() // må sende beskjed til arrived floor handler
			}

		case order := <-readQueue:
			floor = elevatorio.GetFloor() //må skrives om til å bruke PollFloorSensor viss ikke vi kan skrive om elevatorio.
			//TODO floor = state.CurrentFloor, change?
			if floor == -1 {
				fmt.Println("HUGE ERRROR, START DEBUGGING NOW")
				panic("^")
			}
			destination = order.Floor
			if floor < destination {
				elevatorio.SetMotorDirection(elevatorio.MD_Up)
			} else if floor > destination {
				elevatorio.SetMotorDirection(elevatorio.MD_Down)
			} else {
				//TODO elevator already on destionation floor //IS THIS code needed?
				elevatorio.SetMotorDirection(elevatorio.MD_Stop)
				go func() { toArrivedFloorHandler <- state }()

			}
		}
	}
}

func lightManager(
	incomingOrders chan order.Order,
	toElevatorNetwork chan order.Order,
	readCurrentFloor chan int) {

	for {
		select {
		case handleIncomingOrder := <-incomingOrders:
			go func() {
				if handleIncomingOrder.Completed == false {
					if handleIncomingOrder.Class == 2 /*&& handleIncomingOrder.isMine()*/ {
						elevatorio.SetButtonLamp(2, handleIncomingOrder.Floor, true) //BT_Cab, 2 eller CAB??
					}
					if handleIncomingOrder.Class == 0 {
						elevatorio.SetButtonLamp(0, handleIncomingOrder.Floor, true) //BT_hallup, 0 eller HALL_UP??
					}
					if handleIncomingOrder.Class == 1 {
						elevatorio.SetButtonLamp(1, handleIncomingOrder.Floor, true) //BT_hallup/down, 0,1 eller HALL_DOWN??
					}
				}
				if handleIncomingOrder.Completed == true {
					for btn := 0; btn < 3; btn++ {
						elevatorio.SetButtonLamp(elevatorio.ButtonType(btn), handleIncomingOrder.Floor, false)
						/*if handleIncomingOrder.isMine() {
							SetDoorOpenLamp(true)
							time.NewTimer(3 * time.Second) //meeeeget usikker på denne
						}*/
						elevatorio.SetDoorOpenLamp(false)
					}
				}
			}()
		case floorIndicator := <-readCurrentFloor:
			//TODO  trenger man en kanal og egen state for å lese hvilken floor vi er på?
			//floor = getFloor()
			elevatorio.SetFloorIndicator(floorIndicator)
		}

	}
}
