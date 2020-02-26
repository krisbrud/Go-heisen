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
				go func() { toArrivedFloorHandler := <- drv_floors }() // må sende beskjed til arrived floor handler 
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