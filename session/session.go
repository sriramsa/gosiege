// Package  provides a session that shares load across
// the siege nodes
package session

import (
	"github.com/sriramsa/gosiege/logger"
)

type SiegeSession struct {
	Guid string
	Pid  int
	Done chan int
}

// Log to console
var Log = logger.NewLogger("Session")

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
