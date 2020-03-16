package controller

import (
	"time"
)

func Controller(
	toButtonPushHandler chan order.Order, 
	toArrivedFloorHandler chan elevatorstate.ElevatorState, 
	toDelegator chan elevatorstate.ElevatorState, 
	readState chan elevatorstate.ElevatorState,
	readQueue chan order.Order,
) {
	
	var state elevatorstate.ElevatorState  
	var floor int
	var destination int 

	for {
		select {
			
		case buttonPushed := <- drv_buttons:
			go func() { toButtonPushHandler <- buttonPushed }() // må sende dette til button pushed handler 

		case elevatorStateChanged := <-readState:
			if elevatorStateChanged != state{
				go func() { toDelegator <- elevatorStateChanged }() // må lese elevator state og sende den til delegatoren 
				state := elevatorStateChanged
			}
			
		case arrivedAtFloor := <- drv_floors: 
			if arrivedFloor == destination{
				elevio.SetMotorDirection(elevio.MD_Stop)
				go func() { toArrivedFloorHandler <- drv_floors }() // må sende beskjed til arrived floor handler 
			}
			
		case orderToExecute := <- reaadQueue:
			floor = elevio.GetFloor()
			destination = orderToExecute.Floor
			if floor < destination {
				elevio.SetMotorDirection(elevio.MD_Up)
			} else if floor > destination {
				elevio.SetMotorDirection(elevio.MD_Down)
			}
		}
	}
}


func lightManager(
	incomingOrders chan order.Order,
	toElevatorNetwork chan order.Order
) {
	
	for {
		select {
		case handleIncomingOrder := <- incomingOrders: 
			go func() {
				if handleIncomingOrder.Completed == false {
					if handleIncomingOrder.Class == CAB && handleIncomingOrder.isMine(){  
						SetButtonLamp(2, handleIncomingOrder.floor, true)  //BT_Cab, 2 eller CAB??
					}
					if handleIncomingOrder.Class == HALL_UP {
						SetButtonLamp(0, handleIncomingOrder.floor, true)  //BT_hallup, 0 eller HALL_UP??
					}
					if handleIncomingOrder.Class == HALL_DOWN {	 
						SetButtonLamp(1, handleIncomingOrder.floor, true)  //BT_hallup/down, 0,1 eller HALL_DOWN??
					}
				}	
				if handleIncomingOrder.Completed == true {
					for f := 0; f < 3; f++ {
						SetButtonLamp(f, handleIncomingOrder.floor, false)
						if handleIncomingOrder.isMine() {
							SetDoorOpenLamp(true)
							time.Sleep(3* time.Second) //meeeeget usikker på denne
						} 
						SetDoorOpenLamp(false)
					}
				}
			}
		case floorIndicator:
			floor = getFloor()
			SetFloorIndicator(floor)
		}

	}
}