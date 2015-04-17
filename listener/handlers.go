// Package listener provides session REST API handlers
package listener

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/loadcloud/gosiege/state"
)

// Creates a new session
func newSessHandler(w http.ResponseWriter, r *http.Request) {
	event := new(state.NewSiegeSession)

	err := func() error {
		if err := decodeFormIntoStruct(r, event); err != nil {
			return err
		}

		log.Println("event : ", event)
		// Ensure all values are in the structure
		if err := validateFields(*event); err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error : ", err.Error())
		return
	}

	writeToState(state.SessionEvent{*event})

	w.WriteHeader(http.StatusOK)
}

func updateSessHandler(w http.ResponseWriter, r *http.Request) {
	// Id comes in the URL path
	event := new(state.UpdateSiegeSession)

	err := func() error {
		if err := decodeFormIntoStruct(r, event); err != nil {
			return err
		}

		event.SessionId = mux.Vars(r)["Id"]
		log.Println("event : ", *event)
		// Ensure all values are in the structure
		if err := validateFields(*event); err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error : ", err.Error())
		return
	}

	writeToState(state.SessionEvent{*event})

	w.WriteHeader(http.StatusOK)
}

func stopSessHandler(w http.ResponseWriter, r *http.Request) {
	// Id will always be available since it is part of the routing
	id := mux.Vars(r)["Id"]

	stopSession(id)

	w.WriteHeader(http.StatusOK)
}

/* Handler Helpers */

// Decode the form values into the event struct passed
func decodeFormIntoStruct(r *http.Request, e interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Gorilla decoder to decode form into struct
	decoder := schema.NewDecoder()

	log.Println("r.PostForm :", r.PostForm)
	// r.PostForm now has a map of our POST form values
	if err = decoder.Decode(e, r.PostForm); err != nil {
		log.Println("Error decoding.")
		return err
	}

	return nil
}

// Verifies that all the structure fields are filled.
func validateFields(event interface{}) error {
	s := reflect.ValueOf(event)

	for i := 0; i < s.NumField(); i++ {
		name := s.Type().Field(i).Name
		val := s.Field(i).Interface()

		switch val.(type) {
		case int:
			if val == 0 {
				return errors.New(name + " is 0")
			}
		case string:
			if val == "" {
				return errors.New(name + " is empty.")
			}
		default:
			panic("Verification to be implemented for type in handler.")
		}
	}
	return nil
}

func stopSession(id string) {
	siegeCmd := state.SessionEvent{
		Event: state.StopSiegeSession{id},
	}

	log.Println("Sending message to stop session : ", id)

	writeToState(siegeCmd)
}

// Safely handle panic handling the user request
func safelyDo(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Handler Panic:", err)
				http.Error(w, "Error Processing Request", http.StatusBadRequest)
			}
		}()

		f(w, r)
	}
}
