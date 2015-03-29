// Package common provides command resources accessed by all components
package common

import (
	"log"
	"os"
)

// For graceful shutdown of the service. When this is closed, all goroutines exit.
var DoneCh chan struct{}

// Logger
var Log *log.Logger

// Channel for writing to and listening for Admin Commands
// Cluster Manager listens to this.
//var AdminCmdCh chan cluster.SiegeCommand

func InitResources() error {

	Log = log.New(os.Stdout, "", log.Ltime)

	// Create a channel for receiving commands
	//listen := make(chan cluster.SiegeCommand)

	// Print with file and line numbers
	log.SetFlags(log.Llongfile)

	Log.Println("Initialized common resources")

	return nil
}
