// Package listener provides session REST API handlers
package listener

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/loadcloud/gosiege/state"
)

// Creates a new session
func newSessHandler(w http.ResponseWriter, r *http.Request) {
	event := state.NewSiegeSession{
		Concurrent: reqInt(r, "concurrent"),
		Delay:      reqString(r, "delay"),
		Host:       reqString(r, "target"),
		Port:       reqInt(r, "port"),
	}

	// Write
	writeToState(state.SessionEvent{event})

	w.WriteHeader(http.StatusOK)
}

func updateSessHandler(w http.ResponseWriter, r *http.Request) {
	// Create the event. Id will always be available since
	// it is part of the routing
	siegeCmd := state.SessionEvent{
		Event: state.UpdateSiegeSession{
			SessionId:      mux.Vars(r)["Id"],
			NewTargetUsers: reqInt(r, "concurrent"),
		},
	}

	// Write
	writeToState(siegeCmd)

	// Reply
	w.WriteHeader(http.StatusOK)
}

func stopSessHandler(w http.ResponseWriter, r *http.Request) {
	// Create the event. Id will always be available since
	// it is part of the routing
	id := mux.Vars(r)["Id"]

	stopSession(id)

	w.WriteHeader(http.StatusOK)
}

func stopSession(id string) {
	// Create the event. Id will always be available since
	// it is part of the routing
	siegeCmd := state.SessionEvent{
		Event: state.StopSiegeSession{id},
	}

	log.Println("Stopping session : ", id)
	// Write
	writeToState(siegeCmd)
}

// Safely handle panic handling the user request
func safelyDo(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	wf := func(w http.ResponseWriter, r *http.Request) {
		// Handle Panic
		defer func() {
			if err := recover(); err != nil {
				log.Println("work failed:", err)
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, "Error Processing request")
			}
		}()

		// Call actual function
		f(w, r)
	}
	return wf
}
