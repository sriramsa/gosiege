// Package instrumentation
package instrument

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type EventType int

const (
	Info EventType = 1 + iota
	Warning
	Error
)

func (i EventType) String() string {
	s := [...]string{
		"INFO",
		"WARNING",
		"ERROR",
	}

	return s[i]
}

type EventBody struct {
	Message    string      `json:"Msg,omitempty"`
	JSONObject interface{} `json:"Obj,omitempty"`
}

type Event struct {
	Type    string    `json:"Type,omitempty"`
	Time    time.Time `json:"Time,omitempty"`
	Node    string    `json:"Node,omitempty"`
	Package string    `json:"Package,omitempty"`

	Body EventBody `json:"Body,omitempty"`
}

type EventWriter struct {
	l      *log.Logger
	pkg    string
	node   string
	pretty bool
}

func (w *EventWriter) Info(msg string, v ...interface{}) {

	e := Event{
		// TODO: Format time properly
		Time:    time.Now(),
		Package: w.pkg,
		Type:    Info.String(),
		Node:    w.node,
		Body: EventBody{
			Message:    msg,
			JSONObject: v,
		},
	}

	var js []byte
	var err error
	if w.pretty {
		js, err = json.MarshalIndent(e, "", "\t")
	} else {
		js, err = json.Marshal(e)
	}
	if err != nil {
		log.Println("Error JSON Marshal: ", err)
	} else {
		w.l.Println("EVENT : ", string(js))
	}
}

// Create a new logger with the prefix given
func NewEventWriter(p string, w io.Writer, pr bool) *EventWriter {
	h, err := os.Hostname()
	if err != nil {
		log.Println("error getting hostname : ", err)
	}
	return &EventWriter{
		l:      log.New(os.Stdout, "", log.Ltime),
		pkg:    p,
		node:   h,
		pretty: pr,
	}
}
