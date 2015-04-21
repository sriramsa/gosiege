// Package instrumentation
// Lets one generate events within code and provides Attach/Detach
// functions to attach to or detatch from the event stream.
// Will be used for testing too with tests attaching to each
// package stream
package instrument

import "time"

type EventType int

const (
	Info EventType = 0 + iota
	Metric
	Warning
	Error
)

func (i EventType) String() string {
	return [...]string{
		"INFO",
		"METRIC",
		"WARNING",
		"ERROR",
	}[i]
}

type EventBody struct {
	Message string      `json:"Msg,omitempty"`
	Object  interface{} `json:"Obj,omitempty"`
}

type Event struct {
	Type    string    `json:"Type,omitempty"`
	Package string    `json:"Package,omitempty"`
	Time    time.Time `json:"Time,omitempty"`
	Node    string    `json:"Node,omitempty"`

	Body EventBody `json:"Body,omitempty"`
}
