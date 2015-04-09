// Session state
package state

import "log"

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
	SessionId string
	Pid       uint32
	Delay     string
	Host      string
	Proto     string // http or https
	Port      int

	TargetUsers int // Target users for the run
	ActiveUsers int // Active users

	HandlerCh chan SessionEvent

	state sessionState
}

func (s SiegeSession) GetState() sessionState {
	return s.state
}

func (s SiegeSession) SetState(st sessionState) error {
	s.state = st
	return nil
}

func (s SiegeSession) Start() int {
	log.Print("Started")
	return -1
}

func (s SiegeSession) Stop() int {
	log.Print("Stopped")
	return -1
}
