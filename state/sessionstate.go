// Session state
package state

import "github.com/loadcloud/gosiege/logger"

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
	SessionId  string
	Pid        int
	HandlerCh  chan SessionEvent
	Delay      string
	Concurrent string
	Host       string

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
	Log.Print("Started")
	return -1
}

func (s SiegeSession) Stop() int {
	Log.Print("Stopped")
	return -1
}