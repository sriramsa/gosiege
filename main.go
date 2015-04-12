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

	TODO:
		- 12factor service compliance
			- Move config to env vars as per 12factor
			- Logs are streams written directly to the stdout. During local development
			the user can view on screen. In staging/production each process' stream will
			be captured by the execution environment, collated together and routed to
			one or more final destinations.
*/
package main

import (
	"fmt"
	"log"
	"time"

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

	// Temp channel for sending cmds from listener to watcher
	var tempWatcherListenChan chan state.SessionEvent
	tempWatcherListenChan = make(chan state.SessionEvent)

	// Start the State Watcher. This watches for state changes and accepts subscriptions
	// from other components
	go state.StartStateWatcher(tempWatcherListenChan)
	log.Println("StartStateWatcher Done")

	// Start the cluster manager go routine.
	go cluster.StartClusterManager()
	log.Println("StartClusterManager Done")

	// Start session manager
	go session.StartSessionManager()
	log.Println("StartSessionManager Done")

	// Start the http listener that listens to commands from Admin Web UI and gosiege cli
	go listener.StartHttpCommandListener(tempWatcherListenChan)
	log.Println("StartHttpCommandListener Done")

	// Wait for a keystroke to exit.
	fmt.Scanln()

	shutdown()

	// Closing a channel returns zero value immediately to all waiters.
	// Each goroutine has this wait in their select. This will make them exit.
	log.Println("==================== END ====================")
}

func shutdown() {
	listener.ShutdownRESTApiListener()

	// Wait for all the sessions to get kill signal.
	// TODO: Use a channel to signal end
	time.Sleep(time.Second * 2)

	close(common.DoneCh)
}
