// Session manager that manages the set of sessions running on the system.
// Creates a new session
// Stops a session
// Updates a session
// Subscribes to the Session Admin Command events with the State Watcher
// Creates a new session handler go routine for each session.
package session

import (
	"encoding/json"
	"log"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/state"
)

type SessionId string

// List of current sessions running indexed on it's SessionId
var sessList map[SessionId]state.SessionState
var sessCmdCh chan state.SessionEvent

// Start the session manager. Will be done in a go routine
func StartSessionManager() {

	log.Println("Subscribing to Session events with watcher")

	// Subscribe to the StateWatcher for Session Events
	sessCmdCh = state.SubscribeToSessionEvents()

	listenToSessionEvents()
}

func listenToSessionEvents() {

	var cmd state.SessionEvent

	log.Println("Waiting for Session events from watcher.")
	for {
		select {
		case cmd = <-sessCmdCh:
			log.Println("Event received from watcher.")
			parseEvent(cmd)
		case <-common.DoneCh:
			log.Println("DONE signal received, exiting SessionManager")
			return
		}
	}
}

func parseEvent(c state.SessionEvent) {
	switch t := c.Cmd.(type) {
	case state.NewSiegeSession:
		log.Println("NewSiegeSession Command Received", t)
		sessParams, _ := c.Cmd.(state.NewSiegeSession)
		sess := createNewSession(sessParams)
		// Start the session immediately
		startSession(sess)

	case state.StopSiegeSession:
		log.Println("StopSiegeSession Command Received", t)
		log.Println("Event = ", t)
		stopSession(c.Cmd.(state.StopSiegeSession))

	case state.UpdateSiegeSession:
		log.Println("UpdateSiegeSession Command Received", t)
		updateSession(c.Cmd.(state.UpdateSiegeSession))

	case state.EndSiegeSession:
		log.Println("EndSiegeSession Command Received", t)
		endSession(c.Cmd.(state.EndSiegeSession))

	default:
		log.Println("Event = ", t)
	}
}

func createNewSession(c state.NewSiegeSession) state.SiegeSession {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println("Command : ", marshallOut)

	sess := state.SiegeSession{
		SessionId:  "1234",
		Concurrent: c.Concurrent,
		Host:       c.Host,
		Delay:      c.Delay,
		Done:       make(chan int, 1),
	}

	log.Print("Session created...")

	sess.SetState(state.Ready)
	log.Println("New Session Created with sessin id : ", sess.SessionId)
	return sess
}

func startSession(sess state.SiegeSession) {

	// spin up the session instance handler
	go StartSessionHandler(sess)
}

func stopSession(c state.StopSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(string(marshallOut))
}

func updateSession(c state.UpdateSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(string(marshallOut))
}

func endSession(c state.EndSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(string(marshallOut))
}
