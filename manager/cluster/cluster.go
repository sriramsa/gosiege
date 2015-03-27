// Package clustermanager provides a routine that manages the siege cluster
// Listens to a channel for messages
//	- New Session
//	- Stop Session
package cluster

import (
	"time"

	"github.com/sriramsa/gosiege/logger"
)

// StartClusterManager starts the siege cluster manager
func StartClusterManager(doneCh chan struct{}) {
	var Log = logger.NewLogger("ClusterMgr")

	Log.Println("StartingClusterManager")

	listenChannel := make(chan Command)

	// Start the http listener and pass in a channel
	// for it to report the commands in
	go StartGoSiegeHttpListener(listenChannel, doneCh)

	for {
		Log.Println("listening for messages.")

		select {
		case <-doneCh:
			Log.Println("Abort Message Received. Exitting")
			return
		case <-time.After(12 * time.Hour):
		case <-listenChannel:
		}

		Log.Println("Message Received :")
	}
}

func createNewSession() {

}

func stopSession() {

}
