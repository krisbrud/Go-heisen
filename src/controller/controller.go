package controller

import (
	
)


func Controller(
	toButtonPushHandler chan order.Order, 
	toArrivedFloorHandler chan elevatorstate.ElevatorState, 
	toDelegator chan elevatorstate.ElevatorState, 
	readState chan elevatorstate.ElevatorState,
	readQueue chan order.Order
) {
	
	// TODO Setup/init code here
	// TODO declare channels here
	
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
					if handleIncomingOrder.Class == CAB { //legg inn at dette kun skal skje for den heisen det gjelder 
						SetButtonLamp(2, handleIncomingOrder.floor, true)  //BT_Cab, 2 eller CAB??
					}
					if handleIncomingOrder.Class == HALL_UP {
						SetButtonLamp(0, handleIncomingOrder.floor, true)  //BT_hallup, 0 eller HALL_UP??
					}
					if handleIncomingOrder.Class == HALL_DOWN {	 //legg inn at dette skal skje på alle heisene 
						SetButtonLamp(1, handleIncomingOrder.floor, true)  //BT_hallup/down, 0,1 eller HALL_UP/DOWN??
					}
				}	
				if handleIncomingOrder.Completed == true {
					for f := 0; f < 3; f++ {
						SetButtonLamp(f, handleIncomingOrder.floor, false) 
					}
				}
			}
		case floorIndicator:
			floor = getFloor()
			SetFloorIndicator(floor)
		}
		case openDoor: 	// sett åpen dør knapp 
			//settes viss motordir = stop og timeren ikke har rent ut 

	}
}