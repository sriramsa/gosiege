// Package session provides session handler that handles the individual session
// This is the main function that
//	- Spin up the siege session
//	- Load balance along peers
package session

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/loadcloud/gosiege/state"
)

// Current Siege Session being managed by this handler
var sess state.SiegeSession

// List of processes running siege instances
var siegeProcs = make([]exec.Cmd, 0)

func StartSessionHandler(session state.SiegeSession) {

	if jSess, err := json.MarshalIndent(session, "", "\t"); err != nil {
		log.Println("Starting a new SESSION HANDLER for : ", string(jSess))
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

func parseCommand(e state.SessionEvent) {

	switch e.Event.(type) {
	case state.UpdateSiegeSession:
		log.Println("Update Siege Session")
		updateSiege(e.Event.(state.UpdateSiegeSession))

	case state.StopSiegeSession:
		log.Println("Stop Siege Session received")
		stopAllSiege()
	}
}

func CalculateMaxRpsAvailable() uint {
	return 1000
}

func updateSiege(e state.UpdateSiegeSession) {
	log.Println("Updating Siege - New target : ", e.NewTargetUsers, " Old :", sess.TargetUsers)
	sess.TargetUsers = e.NewTargetUsers

	startOrUpdateSiege()
}

// Stop all siege processes
func stopAllSiege() {
	log.Println("Stopping all Siege processes")

	for i := range siegeProcs {
		cmd := siegeProcs[i]
		// go
		func() {
			log.Print(". Killing Process : ", cmd.Process.Pid)
			if err := cmd.Process.Kill(); err != nil {
				log.Println("Could not Kill : ", err)
			} else {
				log.Println("...killed")
			}
		}()
		//marshallOut, err = json.MarshalIndent(cmd, "after :", "\t")
		//log.Println(string(marshallOut))

		// FIX with sync
		sess.ActiveUsers = 0

		// Is there a better way?
		siegeProcs = make([]exec.Cmd, 0)
	}
}

func startOrUpdateSiege() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Command Failed : ", err)
			sess.SetState(state.Error)

			stopAllSiege()
		}
	}()

	spinUpSiege(10)
	// Number of siege procs we need
	//numSigeProcs := (sess.Concurrent / users) + ((sess.Concurrent % users) / (sess.Concurrent % users))

	// Create the slice to hold these procs
	//siegeProcs = make([]exec.Cmd, numSigeProcs)

	//marshallOut, _ := json.MarshalIndent(cmd, "", "\t")

	//log.Println("Cmd : ", string(marshallOut))
}

// Spins up one siege per usersPerSiege
func spinUpSiege(usersPerSiege int) {
	// Build the parameters
	delayParam := fmt.Sprintf("--delay=%s", sess.Delay)
	URLParam := fmt.Sprintf("http://%s:%d", sess.Host, sess.Port)

	for sess.TargetUsers > sess.ActiveUsers {
		users := sess.TargetUsers - sess.ActiveUsers
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

		// Starting siege
		if err := cmd.Start(); err != nil {
			log.Panic("ERROR: ", err)
			return
		}

		sess.ActiveUsers += users

		addToProcList(*cmd)
		log.Println("Siege process created - PID : ", cmd.Process.Pid, " ", userParam)
	}

	if jSess, err := json.MarshalIndent(sess, "", "\t"); err != nil {
		log.Println("Current session : ", string(jSess))
	} else {
		log.Println("Error JSON MarshalIndent :", err)
	}
}

func addToProcList(cmd exec.Cmd) {
	siegeProcs = append(siegeProcs, cmd)
}
