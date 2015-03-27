package session

import (
	"encoding/json"
	"os/exec"
	"time"

	"github.com/loadcloud/gosiege/logger"
	"github.com/loadcloud/gosiege/session"
)

var Log = logger.NewLogger("SessionManager")

type SiegeSession struct {
	Guid string
	Pid  int
	Done chan int
}

// Creates a new SiegeSession Struct and returns a pointer to it
func CreateSiegeSession() *SiegeSession {
	Log.Print("Session created...")
	return &SiegeSession{
		Pid:  10,
		Done: make(chan int, 1),
	}
}

func (session SiegeSession) Start() int {
	Log.Print("Started")
	return -1
}

func (session SiegeSession) Stop() int {
	Log.Print("Stopped")
	return -1
}

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
