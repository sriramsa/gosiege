// Package main provides listener that listens to the http port for incoming
// messages from Admin NodeJS UI or Command Line
/*
 REST API

*/
package listener

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/config"
	"github.com/loadcloud/gosiege/state"
)

var urlPrefix string

var tempWriteCh chan state.SessionEvent

// Starts a http listener and reports incoming messages to the caller
func StartHttpCommandListener(writeCh chan state.SessionEvent) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("FATAL : Listener failed", err)
		}

		// Let everybody exit since we can't listen for incoming commands
		close(common.DoneCh)
	}()

	tempWriteCh = writeCh
	port := config.Get("ListeningPort")
	urlPrefix = config.Get("SiegePath")

	regApiRoutes()

	log.Println("Listening on port : ", port)

	//server := http.Server{ Addr: addr, //ErrorLog:, TODO }
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Could not listen on port : ", port, err)

		// Let everybody exit since we can't listen for incoming commands
		close(common.DoneCh)
	}
}

func ShutdownRESTApiListener() {
	stopSession("1234")
}

func reqString(r *http.Request, s string) (val string, err error) {

	val = r.FormValue(s)

	if val == "" {
		return val, errors.New(s + " could not be read. error :")
	}

	return val, nil
}

func reqInt(r *http.Request, s string) (val int, err error) {
	var sv string
	err = func() error {
		if sv, err = reqString(r, s); err != nil {
			return err
		}

		val, err = strconv.Atoi(sv)
		return err
	}()

	if err != nil {
		return 0, errors.New(s + " could not be found.")
	}
	return val, nil
}

func writeToState(cmd state.SessionEvent) {
	// this is temp
	tempWriteCh <- cmd
}
