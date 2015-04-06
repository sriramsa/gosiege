// Package common provides command resources accessed by all components
package common

import "log"

// For graceful shutdown of the service. When this is closed, all goroutines exit.
var DoneCh chan struct{}

// Channel for writing to and listening for Admin Commands
// Cluster Manager listens to this.
//var AdminCmdCh chan cluster.SiegeCommand

func InitResources() error {

	DoneCh = make(chan struct{})

	// Print with file and line numbers
	log.SetFlags(log.Lshortfile)

	log.Println("Initialized common resources")

	return nil
}
