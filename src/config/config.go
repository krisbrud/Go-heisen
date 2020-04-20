package config

import (
	"flag"
	"os"
	"strconv"
	"sync"
)

// Default values for config
const (
	defaultDriverPort = 15657
	defaultOrderPort  = 44232
	defaultStatePort  = 44233
	defaultNumFloors  = 4
	// Default ElevatorID dynamically set to make them unique
)

type config struct {
	ElevatorDriverPort int
	OrderPort          int
	Stateport          int
	ElevatorID         string
	NumFloors          int
	BottomFloor        int
	TopFloor           int
}

var mtx sync.Mutex
var initialized bool = false
var myConfig config

// Reuse the getElevatorConfig code for parsing values automatically
func GetElevatorDriverPort() int { return getElevatorConfig().ElevatorDriverPort }
func GetOrderPort() int          { return getElevatorConfig().OrderPort }
func GetStatePort() int          { return getElevatorConfig().Stateport }
func GetMyElevatorID() string    { return getElevatorConfig().ElevatorID }
func GetNumFloors() int          { return getElevatorConfig().NumFloors }
func GetBottomFloor() int        { return getElevatorConfig().BottomFloor }
func GetTopFloor() int           { return getElevatorConfig().TopFloor }

// ParseConfigFlags parses command line flags for the configuration of the system, and sets the parameters to logical values if not
func ParseConfigFlags() {
	if initialized {
		return // Don't parse flags multiple times, they can't change.
	}

	mtx.Lock()
	defer mtx.Unlock()
	// Ensure unique elevator ids if not provided
	// Set defaultID to "elev-" + local network ip address + parent process id
	parentProcessID := os.Getppid()
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "" // In case the host name is not available
	}
	defaultElevatorID := "elev-" + hostName + strconv.Itoa(parentProcessID)

	// Define the flags to parse
	flag.IntVar(&myConfig.ElevatorDriverPort, "port", defaultDriverPort, "Port for connection to elevator")
	flag.IntVar(&myConfig.NumFloors, "floors", defaultNumFloors, "Number of floors in each elevator")
	flag.IntVar(&myConfig.OrderPort, "orderport", defaultOrderPort, "Port to broadcast and receive order updates")
	flag.IntVar(&myConfig.Stateport, "stateport", defaultStatePort, "Port to broadcast and receive state updates")

	flag.StringVar(&myConfig.ElevatorID, "id", defaultElevatorID, "ID of this elevator")

	// Parse the flags, set the variables
	flag.Parse()

	// Set bottom and top floors
	myConfig.BottomFloor = 0
	myConfig.TopFloor = myConfig.BottomFloor + myConfig.NumFloors - 1

	initialized = true
}

func getElevatorConfig() config {
	// Make sure the config is initialized
	if !initialized {
		ParseConfigFlags()
	}
	mtx.Lock()
	defer mtx.Unlock()
	return myConfig
}
