// Package session provides session handler that handles the individual session
// This is the main function that
//	- Spin up the siege session
//	- Load balance along peers
package session

import (
	"log"

	"github.com/loadcloud/gosiege/state"
)

func StartSessionHandler(sess state.SiegeSession) {
	// Start the protocol

	// Get current capability
	maxRps := CalculateMaxRpsAvailable()

	log.Println(maxRps)
	// Lock the session in the data store

	// Read current session state from distributed store

	// Calculate proposal

	// Update the session

	// Release the lock on the session

}

func CalculateMaxRpsAvailable() uint {
	return 1000
}

func startSiege(completed chan int) {
	/*
		siegeSession := NewSession()
		_ = siegeSession.Start()

		defer func() {
			if err := recover(); err != nil {
				Log.Println("Command Failed : ", err)

				completed <- -1
			}
		}()

		cmd := exec.Command("siege", "--version")

		marshallOut, _ := json.MarshalIndent(cmd, "", "\t")

		err := cmd.Start()
		if err != nil {

			log.Fatal("ERROR: ", err)
			log.Fatal(string(marshallOut))
		}

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
