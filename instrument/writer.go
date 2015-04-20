// Package instrument provides EventWriter, the main object that the clients
// use to instrument
package instrument

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

// Object that provides interfaces to write an event
type EventWriter struct {
	pkg    string // Package name
	node   string // Node name
	pretty bool   // Should this be pretty printed

	out io.Writer // MultiWriter, used for writing

	swapOut    io.Writer   // New MultiWriter to be used
	logReaders []io.Writer // List of writers attached. Used for attach/detach op
	log        *log.Logger // Log package logger to log events
}

func NewEventWriter(p string, out io.Writer, pretty bool) *EventWriter {
	h, err := os.Hostname()
	if err != nil {
		log.Println("error getting hostname : ", err)
		h = "<ERROR_READING>"
	}

	w := &EventWriter{
		pkg:    p,
		node:   h,
		pretty: pretty,
	}

	if out != nil {
		w.Attach(out)
	}

	return w
}

// Attach a new writer to the event stream
func (w *EventWriter) Attach(out io.Writer) {
	// If this is the first attatch
	if w.logReaders == nil {
		// Create the channel
		w.logReaders = make([]io.Writer, 1)
		w.logReaders[0] = out

		// Create a pipe to let Log package write to, for us
		// to listen to
		pr, pw := io.Pipe()

		// Create a log instance from the Log package
		w.log = log.New(pw, "", log.Ltime)
		w.out = io.MultiWriter(w.logReaders...)

		// Start a co-routine to write in the background
		go w.writer(pr)
	} else {
		w.logReaders = append(w.logReaders, out)
		w.swapOut = io.MultiWriter(w.logReaders...)
	}
}

// Detach a reader from the list
func (w *EventWriter) Detach(out io.Writer) {
	log.Panic("Detach not implemented yet.")
}

// Write to all the writers attached.
// Swap out the writer to new writer list if there is a new list
func (w *EventWriter) writer(in io.Reader) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if w.swapOut != nil {
			w.out = w.swapOut
			w.swapOut = nil
		}
		if w.out != nil {
			w.out.Write(scanner.Bytes())
		}
	}
}

func (w *EventWriter) Error(msg string, v ...interface{}) {
	if w.out != nil {
		w.emitEvent(Error, msg, v)
	}
}

func (w *EventWriter) Warn(msg string, v ...interface{}) {
	if w.out != nil {
		w.emitEvent(Warning, msg, v)
	}
}

func (w *EventWriter) Metric(msg string, v ...interface{}) {
	if w.out != nil {
		w.emitEvent(Metric, msg, v)
	}
}

func (w *EventWriter) Info(msg string, v ...interface{}) {
	if w.out != nil {
		w.emitEvent(Info, msg, v)
	}
}

func (w *EventWriter) emitEvent(t EventType, msg string, v ...interface{}) {
	e := Event{
		Type:    t.String(),
		Package: w.pkg,
		Time:    time.Now(), // TODO: Format time properly
		Node:    w.node,
		Body: EventBody{
			Message: msg,
			Object:  v,
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
		w.log.Println("Error JSON Marshal: ", err)
	}

	// Use the logger we created to print. This takes care of sync when more
	// than one routine is using this EventWriter
	w.log.Println(string(js))
}
