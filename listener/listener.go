// Package main provides listener that listens to the http port for incoming
// messages from Admin NodeJS UI or Command Line
/*
 REST API

*/
package listener

import (
	"log"
	"net/http"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/config"
	"github.com/loadcloud/gosiege/state"
	"github.com/sriramsa/testrument"
)

var urlPrefix string

var tempWriteCh chan state.SessionEvent

var emit *testrument.EventStream

func init() {
	emit = testrument.NewEventStream("listener", true)
}

// Starts a http listener and reports incoming messages to the caller
func StartRESTApiListener(writeCh chan state.SessionEvent) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("FATAL : Listener failed", err)
			emit.Error("Listener start failed", err)
		}

		// Let everybody exit since we can't listen for incoming commands
		close(common.DoneCh)
	}()

	tempWriteCh = writeCh
	port := config.Get("ListeningPort")
	urlPrefix = config.Get("SiegePath")

	regApiRoutes()

	emit.Info("Listening on port : ", port)

	//server := http.Server{ Addr: addr, //ErrorLog:, TODO }
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		emit.Error("Could not listen on port, Exitting :", port, err)

		// Let everybody exit since we can't listen for incoming commands
		close(common.DoneCh)
	}
}

func ShutdownRESTApi() {
	//stopAllSessions()
}

func writeToState(cmd state.SessionEvent) {
	log.Println("sending cmd : ", cmd)
	// this is temp
	tempWriteCh <- cmd
}
