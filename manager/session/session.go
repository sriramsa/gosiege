package session

import (
	"encoding/json"
	"os/exec"
	"time"

	"github.com/loadcloud/gosiege/logger"
)

var Log = logger.NewLogger("SessionManager")

type sessionState int

const (
	Ready sessionState = iota
	Running
	Stopping
	Stopped
	Aborted
	Error
)

// SiegeSession struct
type SiegeSession struct {
	Guid string
	Pid  int
	Done chan int

	state sessionState
}

func (s SiegeSession) GetState() sessionState {
	return s.state
}

func (s SiegeSession) SetState(st sessionState) error {

	s.state = st

	return nil
}

// Creates a new SiegeSession Struct and returns a pointer to it
func NewSession() *SiegeSession {

	Log.Print("Session created...")
	s := SiegeSession{
		Pid:  10,
		Done: make(chan int, 1),
	}

	Log.Println("Session State = ", s.GetState())

	return &s
}

func (s SiegeSession) Start() int {
	Log.Print("Started")
	return -1
}

func (s SiegeSession) Stop() int {
	Log.Print("Stopped")
	return -1
}

func runSiege(completed chan int) (retOutput string) {

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
