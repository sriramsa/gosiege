// Package clustermanager provides a routine that manages the siege cluster
// Listens to a channel for messages
//	- New Session
//	- Stop Session
package cluster

import (
	"encoding/json"
	"time"

	"github.com/loadcloud/gosiege/logger"
	"github.com/loadcloud/gosiege/manager/session"
)

var Log = logger.NewLogger("ClusterMgr")

// StartClusterManager starts the siege cluster manager. Takes a
// channel for listening to abort signals
func StartClusterManager(l, chan SiegeComman, d chan struct{}) {

	// If there is a panic recover using this function
	defer func() {
		if err := recover(); err != nil {
			Log.Println("StartClusterManager - Failed : ", err)
		}
	}()

	Log.Println("StartingClusterManager")

	listenToIncomingCommands(l, d)
}

// listens to incoming commands on the channel
func listenToIncomingCommands(l chan SiegeCommand, d chan struct{}) {

	var cmd SiegeCommand

	for {
		Log.Println("listening for commands from http listener.")

		select {
		case <-time.After(24 * time.Hour):
			// closing release all listeners of doneCh
			close(d)
		case cmd = <-l:
			parseCommand(cmd)
		case <-d:
			Log.Println("Abort Message Received. Exitting")
			return
		}

		Log.Println("Message Received :")
	}
}

func parseCommand(c SiegeCommand) {
	switch t := c.cmd.(type) {
	case NewSiegeSession:
		Log.Println("Command = ", t)
		createNewSession(c.cmd.(NewSiegeSession))

	case StopSiegeSession:
		Log.Println("Command = ", t)
		stopSession(c.cmd.(StopSiegeSession))

	case UpdateSiegeSession:
		Log.Println("Command = ", t)
		updateSession(c.cmd.(UpdateSiegeSession))

	case EndSiegeSession:
		Log.Println("Command = ", t)
		endSession(c.cmd.(EndSiegeSession))

	default:
		Log.Println("Command = ", t)
	}
}

func createNewSession(c NewSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	Log.Println(marshallOut)

	session.NewSession()
}

func startSession() {

}

func stopSession(c StopSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	Log.Println(marshallOut)
}

func updateSession(c UpdateSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	Log.Println(marshallOut)
}

func endSession(c EndSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	Log.Println(marshallOut)
}
