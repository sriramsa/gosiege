// A list of administration commands
package cluster

// Holds a generic command. Used for sending across channels
type Command struct {
	cmd interface{}
}

type NewSiegeSession struct {
	concurrent, delay, host string
}

type EndSiegeSession struct {
	sessionId string
}

type StopSiegeSession struct {
	sessionId string
}

type UpdateSiegeSession struct {
	concurrent, delay, host string
}
