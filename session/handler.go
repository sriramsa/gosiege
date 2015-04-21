// Package session provides session handler that handles the individual session
// This is the main function that
//	- Spin up the siege session
//	- Load balance along peers
package session

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/state"
)

// Current Siege Session being managed by this handler
//var sess state.SiegeSession

type siegeProcInfo struct {
	Concurrent int
	Proc       exec.Cmd
}

type SessionHandler struct {
	State     state.SiegeSession
	StartTime time.Time // When was it started
	EndTime   time.Time // When it should end

	// List of siege processes in this session
	Procs []siegeProcInfo

	// Session Handler listens on this channel admin commands
	ListenCh chan state.SessionEvent
}

// List of processes running siege instances. Indexed by session id.
//var siegeProcs = make([]exec.Cmd, 0)

// Create a new SessionHandler and return it
func NewSessionHandler(state state.SiegeSession, listen chan state.SessionEvent) *SessionHandler {
	h := SessionHandler{
		State:     state,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Procs:     make([]siegeProcInfo, 0),

		ListenCh: listen,
	}

	return &h
}

func (h *SessionHandler) Start() {

	event.Info("Starting Session Handler", h)

	// Start the protocol

	// Get current capability
	//maxRps := CalculateMaxRpsAvailable()

	//log.Println("Max RPS : ", maxRps)
	//log.Println("host : ", h.State.Host)
	//log.Println("concurrent : ", h.State.TargetUsers)
	//log.Println("delay : ", h.State.Delay)

	h.startOrUpdateSiege()
	// Lock the session in the data store

	// Read current session state from distributed store

	// Calculate proposal

	// Update the session

	// Release the lock on the session

	// listen for commands
	for {
		log.Println("SessionHandler waiting for Done")
		select {
		//case cmd := <-h.ListenCh:
		//h.handleCommand(cmd)
		case <-common.DoneCh:
			log.Println("EXIT signal received, stopping procs and exiting SessionHandler : ", h.State.SessionId)
			h.Stop()
			return
		}
	}
}

/*
func StartSessionHandler(session state.SiegeSession) {

	if jSess, err := json.MarshalIndent(session, "", "\t"); err != nil {
		log.Println("Starting a new session handler for session : ", string(jSess))
	} else {
		log.Println("Error JSON MarshalIndent :", err)
	}

	// Start the protocol

	// Get current capability
	maxRps := CalculateMaxRpsAvailable()

	log.Println("Max RPS : ", maxRps)
	log.Println("host : ", session.Host)
	log.Println("concurrent : ", session.TargetUsers)
	log.Println("delay : ", session.Delay)

	sess = session

	startOrUpdateSiege()
	// Lock the session in the data store

	// Read current session state from distributed store

	// Calculate proposal

	// Update the session

	// Release the lock on the session

	// listen for commands
	for {
		select {
		case cmd := <-session.HandlerCh:
			parseCommand(cmd)
		}
	}
}
*/
func (h *SessionHandler) handleCommand(e state.SessionEvent) {

	event.Info("Handler: Event Received", e)
	switch e.Event.(type) {
	case state.UpdateSiegeSession:
		log.Println("Update Siege Session")
		//h.updateSiege(e.Event.(state.UpdateSiegeSession))

	case state.StopSiegeSession:
		log.Println("Stop Siege Session received")
		//h.stopAllSiege()
	}
}

func CalculateMaxRpsAvailable() uint {
	return 1000
}

func (h *SessionHandler) Update(e state.UpdateSiegeSession) {
	log.Println("Updating Siege - New target : ", e.Concurrent, " Old :", h.State.TargetUsers)
	h.State.TargetUsers = e.Concurrent

	h.startOrUpdateSiege()
}

// Stop all siege processes in this session
func (h *SessionHandler) Stop() {
	log.Println("Stopping ", len(h.Procs), " Siege processes : ", h.State.SessionId)

	for i := range h.Procs {
		cmd := h.Procs[i].Proc

		// go
		func() {
			log.Print("... Killing Process : ", cmd.Process.Pid)
			if err := cmd.Process.Kill(); err != nil {
				log.Println("Could not Kill : ", err)
			} else {
				log.Println("...killed")
			}
		}()

	}
	// FIX with sync
	h.State.ActiveUsers = 0

	// Is there a better way to empty this?
	h.Procs = nil
}

func (h *SessionHandler) startOrUpdateSiege() {
	/*
		defer func() {
			if err := recover(); err != nil {
				log.Println("Command Failed : ", err)
				h.State.SetState(state.Error)

				h.Stop()
			}
		}()
	*/
	h.spinUpSiege(10)
	// Number of siege procs we need
	//numSigeProcs := (sess.Concurrent / users) + ((sess.Concurrent % users) / (sess.Concurrent % users))

	// Create the slice to hold these procs
	//siegeProcs = make([]exec.Cmd, numSigeProcs)
}

// Spins up one siege per usersPerSiege
func (h *SessionHandler) spinUpSiege(usersPerSiege int) {
	// Build the parameters
	delayParam := fmt.Sprintf("--delay=%s", h.State.Delay)
	URLParam := fmt.Sprintf("http://%s:%d", h.State.Host, h.State.Port)

	for h.State.TargetUsers > h.State.ActiveUsers {
		users := h.State.TargetUsers - h.State.ActiveUsers
		if users > usersPerSiege {
			users = usersPerSiege
		}

		userParam := fmt.Sprintf("--concurrent=%d", users)

		// Construct the command
		cmd := exec.Command("siege",
			"--quiet",
			userParam,
			delayParam,
			URLParam)

		// Start the Process
		if err := cmd.Start(); err != nil {
			log.Panic("ERROR: ", err)
			return
		}

		h.State.ActiveUsers += users

		h.addProc(users, *cmd)
		log.Println("Siege process created - PID : ", cmd.Process.Pid, " ", userParam)
		event.Info("New Siege Proc", cmd)
	}

	event.Info("Session information", h.State)
}

func (h *SessionHandler) addProc(u int, cmd exec.Cmd) {
	p := siegeProcInfo{
		Concurrent: u,
		Proc:       cmd,
	}

	h.Procs = append(h.Procs, p)
}
