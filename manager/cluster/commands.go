// A list of administration commands
package cluster

type Command struct {
	cmd interface{}
}

type CreateNewSession struct {
	concurrent, delay, host string
}
