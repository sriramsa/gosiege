package session

import (
	"encoding/json"
	"os/exec"
	"time"

	"github.com/sriramsa/gosiege/logger"
	"github.com/sriramsa/gosiege/session"
)

var Log = logger.NewLogger("SessionManager")

func runSiege(completed chan int) (retOutput string) {

	siegeSession := session.CreateSiegeSession()
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

		Log.Fatal("ERROR: ", err)
		Log.Fatal(string(marshallOut))
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
}
