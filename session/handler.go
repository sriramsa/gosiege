// Package session provides session handler that handles the individual session
// This is the main function that
//	- Spin up the siege session
//	- Load balance along peers
package session

import (
	"encoding/json"
	"log"
	"os/exec"

	"github.com/loadcloud/gosiege/state"
)

var siegeProcessList = make([]exec.Cmd, 0)

func StartSessionHandler(sess state.SiegeSession) {

	if jSess, err := json.MarshalIndent(sess, "", "\t"); err != nil {
		log.Println("Starting a new SESSION HANDLER for : ", string(jSess))
	} else {
		log.Println("Error JSON MarshalIndent :", err)
	}

	// Start the protocol

	// Get current capability
	maxRps := CalculateMaxRpsAvailable()

	log.Println("Max RPS : ", maxRps)
	log.Println("host : ", sess.Host)
	log.Println("concurrent : ", sess.Concurrent)
	log.Println("delay : ", sess.Delay)

	startSiege(sess)
	// Lock the session in the data store

	// Read current session state from distributed store

	// Calculate proposal

	// Update the session

	// Release the lock on the session

	// listen for commands
	for {
		select {
		case cmd := <-sess.HandlerCh:
			parseCommand(cmd)
		}
	}
}

func parseCommand(e state.SessionEvent) {
	switch e.Event.(type) {
	case state.UpdateSiegeSession:
		log.Println("Update Siege Session")
	case state.StopSiegeSession:
		log.Println("Stop Siege Session received")
	}
}

func CalculateMaxRpsAvailable() uint {
	return 1000
}

func startSiege(sess state.SiegeSession) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Command Failed : ", err)
		}
	}()

	cmd := exec.Command("siege", "--quiet", "-c", sess.Concurrent, "-d", sess.Delay, "http://localhost:8888")

	marshallOut, _ := json.MarshalIndent(cmd, "", "\t")

	log.Println("Cmd : ", string(marshallOut))
	// Starting siege
	err := cmd.Start()
	if err != nil {

		log.Fatal("ERROR: ", err)
		log.Fatal(string(marshallOut))
	}

	/*
		Log.Println("Waiting for 5 secs")
		time.Sleep(5 * time.Second)
		Log.Println("Sending kill signal")
		if err := cmd.Process.Kill(); err != nil {
			Log.Println("Could not Kill : ", err)
		} else {
			Log.Println("Killed Process")
		}

		process := cmd.Process

		retOutput = "abcd"

		time.Sleep(2 * time.Second)

		marshallOut, err = json.MarshalIndent(cmd, "after :", "\t")
		Log.Println("\nAFTER:")
		Log.Println(string(marshallOut))

		completed <- process.Pid

		return retOutput
	*/
}
