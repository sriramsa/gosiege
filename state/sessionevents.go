// A list of administration event
package state

// Holds a generic event. Used for sending across channels
type SessionEvent struct {
	Event interface{}
}

type NewSiegeSession struct {
	Concurrent, Delay, Host string
}

type EndSiegeSession struct {
	SessionId string
}

type StopSiegeSession struct {
	SessionId string
}

type UpdateSiegeSession struct {
	Concurrent, Delay, Host string
}
