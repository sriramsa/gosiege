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
		case <-sessCmdCh:
			parseEvent(cmd)
		case <-common.DoneCh:
			log.Println("DONE signal received, extting SessionManager")
			return
		}
	}
}

func parseEvent(c state.SessionEvent) {
	switch t := c.Cmd.(type) {
	case state.NewSiegeSession:
		log.Println("Event = ", t)
		sess := createNewSession(c.Cmd.(state.NewSiegeSession))
		// Start the session immediately
		startSession(sess)

	case state.StopSiegeSession:
		log.Println("Event = ", t)
		stopSession(c.Cmd.(state.StopSiegeSession))

	case state.UpdateSiegeSession:
		log.Println("Event = ", t)
		updateSession(c.Cmd.(state.UpdateSiegeSession))

	case state.EndSiegeSession:
		log.Println("Event = ", t)
		endSession(c.Cmd.(state.EndSiegeSession))

	default:
		log.Println("Event = ", t)
	}
}

func createNewSession(c state.NewSiegeSession) state.SiegeSession {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(marshallOut)

	log.Print("Session created...")
	sess := state.SiegeSession{
		Pid:  10,
		Done: make(chan int, 1),
	}

	sess.SetState(state.Ready)
	log.Println("Session State = ", sess.GetState())

	return sess
}

func startSession(sess state.SiegeSession) {

	// spin up the session instance handler
	go StartSessionHandler(sess)
}

func stopSession(c state.StopSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(marshallOut)
}

func updateSession(c state.UpdateSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(marshallOut)
}

func endSession(c state.EndSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(marshallOut)
}
