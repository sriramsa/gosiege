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
	"fmt"
	"log"

	"github.com/loadcloud/gosiege/cluster"
	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/config"
	"github.com/loadcloud/gosiege/listener"
	"github.com/loadcloud/gosiege/session"
	"github.com/loadcloud/gosiege/state"
)

func main() {
	// If there is a panic recover using this function
	defer func() {
		if err := recover(); err != nil {
			log.Println("MAIN: Failed : ", err)
		}
	}()

	// Initialize common resources used across the components
	// logger, channels etc.,
	_ = common.InitResources()
	log.Println("Resources Initialized")

	// Load the configuration
	_ = config.LoadConfig()
	log.Println("Configuration loaded")

	// Initialize Distributed State Engine
	// This also starts a go routine that watches changes and informs
	// the corresponding component of the change
	_ = state.InitGoSiegeState()
	log.Println("InitGoSiegeState Done")

	// Start the State Watcher. This watches for state changes and accepts subscriptions
	// from other components
	go state.StartStateWatcher()
	log.Println("StartStateWatcher Done")

	// Start the cluster manager go routine.
	go cluster.StartClusterManager()
	log.Println("StartClusterManager Done")

	// Start session manager
	go session.StartSessionManager()
	log.Println("StartSessionManager Done")

	// Start the http listener that listens to commands from Admin Web UI and gosiege cli
	go listener.StartHttpCommandListener()
	log.Println("StartHttpCommandListener Done")

	// Wait for a keystroke to exit.
	fmt.Scanln()

	// Closing a channel returns zero value immediately to all waiters.
	// Each goroutine has this wait in their select. This will make them exit.
	close(common.DoneCh)
	log.Println("==================== END ====================")
}
