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
func StartClusterManager(d chan struct{}) {

	// If there is a panic recover using this function
	defer func() {
		if err := recover(); err != nil {
			Log.Println("StartClusterManager - Failed : ", err)
		}
	}()

	Log.Println("StartingClusterManager")

	listenChannel := make(chan Command)

	// Start the http listener and pass in a channel
	// for it to report the commands in
	go StartHttpCommandListener(listenChannel, d)

	listenToIncomingCommands(listenChannel, d)
}

// listens to incoming commands on the channel
func listenToIncomingCommands(l chan Command, d chan struct{}) {

	var cmd Command

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

func parseCommand(c Command) {
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
