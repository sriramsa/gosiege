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
	"github.com/loadcloud/gosiege/instrument"
	"github.com/loadcloud/gosiege/listener"
	"github.com/loadcloud/gosiege/session"
	"github.com/loadcloud/gosiege/state"
)

// Event Writer for instrumentation
var emit *instrument.EventWriter

func EventHydrantAttach() *instrument.Attach {
	if emit == nil {
		// Race condition with main just started and somebody trying to
		// attach to the event stream
		emit = instrument.NewEventWriter("main", nil, true)
	}
	//return emit.Attach()
	return nil
}

func GetEventHydrant() *instrument.EventWriter {
	if emit == nil {
		// Race condition with main just started and somebody trying to
		// attach to the event stream
		emit = instrument.NewEventWriter("main", nil, true)
	}
	return emit
}

func main() {
	// If there is a panic recover using this function
	defer func() {
		if err := recover(); err != nil {
			if emit != nil {
				emit.Error("MAIN: Failed : ", err)
			} else {
				log.Println("ERROR - MAIN : Panic : ", err)
			}
		}
	}()

	//event = instrument.NewEventWriter("main", nil, true)
	if emit == nil {
		emit = instrument.NewEventWriter("main", nil, true)
	}

	// Initialize common resources used across the components
	// logger, channels etc.,
	_ = common.InitResources()
	//log.Println("Resources Initialized")
	emit.Info("Resources Initialized")

	// Load the configuration
	_ = config.LoadConfig()
	//log.Println("Configuration loaded")
	emit.Info("Configuration loaded")

	// Initialize Distributed State Engine
	// This also starts a go routine that watches changes and informs
	// the corresponding component of the change
	_ = state.InitGoSiegeState()
	emit.Info("InitGoSiegeState Done")

	// Temp channel for sending cmds from listener to watcher
	var tempWatcherListenChan chan state.SessionEvent
	tempWatcherListenChan = make(chan state.SessionEvent)

	// Start the State Watcher. This watches for state changes and accepts subscriptions
	// from other components
	go state.StartStateWatcher(tempWatcherListenChan)
	emit.Info("StartStateWatcher Done")

	// Start the cluster manager go routine.
	go cluster.StartClusterManager()
	emit.Info("StartClusterManager Done")

	// Start session manager
	go session.StartSessionManager()
	emit.Info("StartSessionManager Done")

	// Start the http listener that listens to commands from Admin Web UI and gosiege cli
	go listener.StartRESTApiListener(tempWatcherListenChan)
	emit.Info("StartHttpCommandListener Done")

	emit.Info("ServerInitDone")

	// Wait for a keystroke to exit.
	// TODO:
	runningFromCmdline := false
	if runningFromCmdline {
		fmt.Scanln()
	} else {
		<-common.DoneCh // will wait forever
	}

	shutdown()

	// Closing a channel returns zero value immediately to all waiters.
	// Each goroutine has this wait in their select. This will make them exit.
	emit.Info("Shutting down")
}

func shutdown() {
	listener.ShutdownRESTApi()

	close(common.DoneCh)

	// Wait for all the sessions to get kill signal.
	// TODO: Use a channel to signal end
	time.Sleep(time.Second * 1)

}
