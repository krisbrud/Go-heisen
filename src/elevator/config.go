package elevator

import (
	"flag"
	"os"
	"strconv"
	"sync"
)

// Default values for config
const (
	defaultElevatorPort = 15657
	defaultNumFloors    = 4
	// Default ElevatorID dynamically set to make them unique
)

type config struct {
	ElevatorDriverPort int
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
func GetMyElevatorID() string    { return getElevatorConfig().ElevatorID }
func GetNumFloors() int          { return getElevatorConfig().NumFloors }
func GetBottomFloor() int        { return getElevatorConfig().BottomFloor }
func GetTopFloor() int           { return getElevatorConfig().TopFloor }

// ParseConfigFlags parses command line flags for the configuration of the system, and sets the parameters to logical values if not
func ParseConfigFlags() {
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
	flag.IntVar(&myConfig.ElevatorDriverPort, "port", defaultElevatorPort, "Port for connection to elevator")
	flag.IntVar(&myConfig.NumFloors, "floors", defaultNumFloors, "Number of floors in each elevator")
	flag.StringVar(&myConfig.ElevatorID, "id", defaultElevatorID, "ID of this elevator")

	// Parse the flags, set the variables
	flag.Parse()

	// Set bottom and top floors
	myConfig.BottomFloor = 0
	myConfig.TopFloor = myConfig.BottomFloor + myConfig.NumFloors - 1

	initialized = true
}

func getElevatorConfig() config {
	// Make config initialize by reading flags and setting them to default if not presenteEE
	if !initialized {
		ParseConfigFlags()
	}
	mtx.Lock()
	defer mtx.Unlock()
	return myConfig
}
