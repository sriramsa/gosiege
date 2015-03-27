/*
Provides the main function to start the siege host process.
This process keeps running that
	- Initializes
	- Session Life-cycle Management
		- Create
		- Start/Stop
		- Monitor for issues

	- General pattern
		- Main initializes and starts the components
		- Each component subscribes for notifications from the StateEngine
		- StateEngine:
*/
package main

import (
	"github.com/sriramsa/gosiege/config"
	"github.com/sriramsa/gosiege/logger"
	"github.com/sriramsa/gosiege/manager/cluster"
	"github.com/sriramsa/gosiege/state"
)

// For graceful shutdown of the service. When this is closed, all goroutines exit.
var DoneCh chan struct{}

func main() {
	var Log = logger.NewLogger("Main")

	Log.Println("==================== BEGIN ====================")

	// If there is a panic recover using this function
	defer func() {
		if err := recover(); err != nil {
			Log.Println("Failed : ", err)
		}
	}()

	// Load the configuration
	_ = config.LoadConfig()

	// Initialize Distributed State Engine
	// This also starts a go routine that watches changes and informs
	// the corresponding component of the change
	_ = state.InitGoSiegeState()

	// Start the cluster manager go routine
	go cluster.StartClusterManager(DoneCh)

	// admin channel is used to listen to incoming messages from
	// admin UI or command line
	adminCh := make(chan string)

	// Start the http command listener

	// Wait for the command
	Log.Println("Waiting for command...")
	select {
	case cmd := <-adminCh:
		Log.Println("Command Received :", cmd)
	}

	// Closing a channel returns zero value immediately to all waiters.
	// Each goroutine has this wait in their select. This will make them exit.
	close(DoneCh)
	Log.Println("==================== END ====================")
}
