// A list of session related events
package state

// Holds a generic event. Used for sending across channels
type SessionEvent struct {
	Event interface{}
}

type NewSiegeSession struct {
	Target     string `schema:"target"`
	Port       int    `schema:"port"`
	Concurrent int    `schema:"concurrent"`
	Delay      string `schema:"delay"`
}

type EndSiegeSession struct {
	SessionId string
}

type StopSiegeSession struct {
	SessionId string
}

type UpdateSiegeSession struct {
	SessionId string `schema:"Id"`

	Target     string `schema:"target"`
	Port       int    `schema:"port"`
	Concurrent int    `schema:"concurrent"`
	Delay      string `schema:"delay"`
}
