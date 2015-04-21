// Package instrument provides Attach object that creates a connection to a
// package's published event hydrant
package instrument

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"
)

// Object representing an attach to an event stream
type Attach struct {
	Reader io.Reader
	Writer io.Writer

	// The event writer we are attached to
	eventWriter *EventWriter
}

// Attach a new writer to the event stream
// TODO: Make this an independent object taking EventWriter as a param
func NewAttach(ew *EventWriter) *Attach {
	attach := Attach{
		eventWriter: ew,
	}

	// Create a pipe for the listener to listen to
	attach.Reader, attach.Writer = io.Pipe()

	_ = ew.attach(attach.Writer)

	return &attach
}

// Detach a reader from the list
func (a *Attach) Detach() {
	a.eventWriter.detach(a.Writer)
}

func (a *Attach) WaitForEvent(t EventType, evt string, tout time.Duration) (v interface{}, e error) {
	dec := json.NewDecoder(a.Reader)
	log.Println("Waiting for Event :", evt)

	cont := true
	found := make(chan map[string]interface{})
	go func() {
		for {
			if !cont {
				log.Println("Event Search thread aborted")
				return
			}
			//log.Println("DECODING EVENT")
			var v map[string]interface{}
			if err := dec.Decode(&v); err != nil {
				log.Println("DECODING ERROR EXITTING :", err)
				return
			}
			log.Println(v)
			if v["Type"] == t.String() {
				b := v["Body"].(map[string]interface{})
				if b["Msg"] == evt {
					log.Println("Event Found: ", b["Msg"])
					found <- v
					return
				}
			}
		}
	}()

	select {
	case <-time.After(tout):
		//log.Println("Event Seach Timed Out")
		cont = false
		v = nil
		e = errors.New("Event Search timed out.")

	case v := <-found:
		log.Println("Event Found: ", v)
	}

	return v, e
}
