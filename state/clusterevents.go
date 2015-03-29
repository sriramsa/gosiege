// A list of cluster events
package state

// Holds a generic command. Used for sending across channels
type ClusterEvent struct {
	Cmd interface{}
}

type AddNode struct {
	Concurrent, Delay, Host string
}

type RemoveNode struct {
	SessionId string
}
