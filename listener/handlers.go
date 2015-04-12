// Package listener provides session REST API handlers
package listener

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/loadcloud/gosiege/state"
)

// Creates a new session
func newSessHandler(w http.ResponseWriter, r *http.Request) {
	var concurrent, port int
	var host, delay string
	var err error

	err = func() error {
		if concurrent, err = reqInt(r, "concurrent"); err != nil {
			return errors.New("Valid concurrent not found.")
		}

		if port, err = reqInt(r, "port"); err != nil {
			return errors.New("Valid concurrent not found.")
		}
		return nil
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := state.NewSiegeSession{
		Concurrent: concurrent,
		Delay:      delay,
		Host:       host,
		Port:       port,
	}

	// Write
	writeToState(state.SessionEvent{event})

	w.WriteHeader(http.StatusOK)
}

func updateSessHandler(w http.ResponseWriter, r *http.Request) {
	var concurrent int
	var err error

	err = func() error {
		if concurrent, err = reqInt(r, "concurrent"); err != nil {
			return errors.New("Valid concurrent not found.")
		}
		return nil
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the event. Id will always be available since
	// it is part of the routing
	siegeCmd := state.SessionEvent{
		Event: state.UpdateSiegeSession{
			SessionId:      mux.Vars(r)["Id"],
			NewTargetUsers: concurrent,
		},
	}

	// Write
	writeToState(siegeCmd)

	w.WriteHeader(http.StatusOK)
}

func stopSessHandler(w http.ResponseWriter, r *http.Request) {
	// Id will always be available since it is part of the routing
	id := mux.Vars(r)["Id"]

	stopSession(id)

	w.WriteHeader(http.StatusOK)
}

func stopSession(id string) {
	// Create the event.
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
		defer func() {
			if err := recover(); err != nil {
				log.Println("Handler panic:", err)
				http.Error(w, "Error Processing Request", http.StatusBadRequest)
			}
		}()

		f(w, r)
	}
	return wf
}
