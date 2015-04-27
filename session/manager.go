// Session manager that manages the set of sessions running on the system.
// Creates a new session
// Stops a session
// Updates a session
// Subscribes to the Session Admin Command events with the State Watcher
// Creates a new session handler go routine for each session.
package session

import (
	"encoding/json"
	"io"
	"log"
	"strconv"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/state"
	"github.com/sriramsa/testrument"
)

// Event Writer for instrumentation
var event *testrument.EventStream

var mw, newMw io.Writer
var pr *io.PipeReader
var pw *io.PipeWriter

func EventHydrantAttach(w io.Writer) *testrument.Attach {
	return nil
	//return event.Attach()
}

func init() {
	//w := io.Writer(os.Stdout)
	//mw = io.MultiWriter(w)

	//pr, pw = io.Pipe()
	//go func() {
	//scanner := bufio.NewScanner(pr)
	//for scanner.Scan() {
	//fmt.Fprintln(mw, "AAATTTTTTT", scanner.Text())
	//if mw != newMw {
	//mw = newMw
	//}
	//mw.Write(scanner.Bytes())
	//}
	//}()

	event = testrument.NewEventStream("session", false)
}

// Session Id is just int, being incremented.
var lastSessionId = 100

func nextSessionId() string {
	lastSessionId++
	return strconv.Itoa(lastSessionId)
}

// List of current sessions running indexed on it's SessionId
var sessList map[string]state.SessionState
var sessCmdCh chan state.SessionEvent

// List of session handlers
var handlerList map[string]*SessionHandler = make(map[string]*SessionHandler, 0)

// Map of all the sessions
//var sessMap = make(map[string]state.SiegeSession)

// Start the session manager. Will be done in a go routine
func StartSessionManager() {

	log.Println("Subscribing to Session events with watcher")
	event.Info("Subscribing to Session events with watcher")

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
			handleEvent(cmd)
		case <-common.DoneCh:
			log.Println("EXIT signal received, exiting SessionManager")
			return
		}
	}
}

func handleEvent(c state.SessionEvent) {

	event.Info("Event Received.", c.Event)

	switch t := c.Event.(type) {
	case state.NewSiegeSession:
		log.Println("NewSiegeSession Event Received", t)
		sessParams, _ := c.Event.(state.NewSiegeSession)
		s := createNewSession(sessParams)

		// Start the session immediately
		startSession(s)

	case state.StopSiegeSession:
		log.Println("StopSiegeSession Event Received", t)
		sessParams, _ := c.Event.(state.StopSiegeSession)

		// Send the Event to the session handler
		if h, found := handlerList[sessParams.SessionId]; found {
			log.Println("Sending event to session handler.")
			h.Stop()
		} else {
			log.Println("Session not found Id : ", sessParams.SessionId)
		}

	case state.UpdateSiegeSession:
		log.Println("UpdateSiegeSession Event Received", t)
		sessParams, _ := c.Event.(state.UpdateSiegeSession)

		// Send the Event to the session handler
		if h, found := handlerList[sessParams.SessionId]; found {
			log.Println("Sending event to session handler.")
			h.Update(sessParams)
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
	event.Info("Create New Session : ", c)

	//if _, found := handlerList[c.SessionId]; found {
	//log.Println("Cannot Create - Session already found: ", c.SessionId)
	//return handlerList[c.SessionId].State
	//}

	handlerCh := make(chan state.SessionEvent, 1)

	// Create a new state for the session
	// TODO: Move this into the creation of session handler
	s := state.SiegeSession{
		SessionId: nextSessionId(),
		Host:      c.Target,
		Delay:     c.Delay,
		Port:      c.Port,

		TargetUsers: c.Concurrent,
		ActiveUsers: 0,

		//HandlerCh: handlerCh,
	}

	handlerList[s.SessionId] = NewSessionHandler(s, handlerCh)

	s.SetState(state.Ready)
	log.Println("New Session Created with sessin id : ", s.SessionId)
	return s
}

func startSession(s state.SiegeSession) {
	h := handlerList[s.SessionId]

	// spin up the session instance handler
	go h.Start()

	//sessMap[sess.SessionId] = sess

	// spin up the session instance handler
	//go StartSessionHandler(sess)
}

func endSession(c state.EndSiegeSession) {
	marshallOut, _ := json.MarshalIndent(c, "", "\t")

	log.Println(string(marshallOut))
}
