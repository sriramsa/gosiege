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
	"strconv"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/state"
)

// Session Id is just int, being incremented.
var lastSessionId = 100

func nextSessionId() string {
	lastSessionId++
	return strconv.Itoa(lastSessionId)
}

// List of current sessions running indexed on it's SessionId
var sessList map[string]state.SessionState
var sessCmdCh chan state.SessionEvent

// Map of all the sessions
var sessMap = make(map[string]state.SiegeSession)

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
	switch t := c.Event.(type) {
	case state.NewSiegeSession:
		log.Println("NewSiegeSession Event Received", t)
		sessParams, _ := c.Event.(state.NewSiegeSession)
		sess := createNewSession(sessParams)

		// Start the session immediately
		startSession(sess)

	case state.StopSiegeSession:
		log.Println("StopSiegeSession Event Received", t)
		sessParams, _ := c.Event.(state.StopSiegeSession)

		// Send the Event to the session handler
		if sess, found := sessMap[sessParams.SessionId]; found {
			log.Println("Sending event to session handler.")
			sess.HandlerCh <- c
		} else {
			log.Println("Session not found Id : ", sessParams.SessionId)
		}
		//stopSession(c.Cmd.(state.StopSiegeSession))

	case state.UpdateSiegeSession:
		log.Println("UpdateSiegeSession Event Received", t)
		sessParams, _ := c.Event.(state.UpdateSiegeSession)

		// Send the Event to the session handler
		if sess, found := sessMap[sessParams.SessionId]; found {
			log.Println("Sending event to session handler.")
			sess.HandlerCh <- c
		} else {
			log.Println("Session not found Id : ", sessParams.SessionId)
		}

	case state.EndSiegeSession:
		log.Println("EndSiegeSession Event Received", t)
		endSession(c.Event.(state.EndSiegeSession))

	default:
		log.Println("Event = ", t)
	}
}

func createNewSession(c state.NewSiegeSession) state.SiegeSession {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println("Event : ", marshallOut)

	sess := state.SiegeSession{
		SessionId: nextSessionId(),
		Host:      c.Target,
		Delay:     c.Delay,
		Port:      c.Port,

		TargetUsers: c.Concurrent,
		ActiveUsers: 0,

		HandlerCh: make(chan state.SessionEvent, 1),
	}

	sess.SetState(state.Ready)
	log.Println("New Session Created with sessin id : ", sess.SessionId)
	return sess
}

func startSession(sess state.SiegeSession) {

	if _, found := sessMap[sess.SessionId]; found {
		log.Println("Session already found: ", sess.SessionId)
		return
	}

	sessMap[sess.SessionId] = sess

	// spin up the session instance handler
	go StartSessionHandler(sess)
}

func endSession(c state.EndSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(string(marshallOut))
}
