// Package main provides listener that listens to the http port for incoming
// messages from Admin NodeJS UI or Command Line
package listener

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

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

	//http.HandleFunc(urlPrefix, newSessHandler)
	http.HandleFunc("/gosiege/new", newSessHandler)
	http.HandleFunc("/gosiege/stop", stopSessHandler)
	//addr := "127.0.0.1:" + port

	log.Println("Listening on port : ", port)

	//server := http.Server{ Addr: addr, //ErrorLog:, TODO }
	//if err := server.ListenAndServe(); err != nil {
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Could not listen on port : ", port, err)

		// Let everybody exit since we can't listen for incoming commands
		close(common.DoneCh)
	}
}

func stopSessHandler(w http.ResponseWriter, r *http.Request) {
	// error handling separated from code flow for easy readability
	// Follow this idiom
	err := func() error {
		if r.Method != "GET" {
			return errors.New("expected GET")
		}

		//if input := parseInput(r); input != "command" {
		//return errors.New("malformed command")
		//}
		return nil
	}()

	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "error found")
		return
	}

	// Log the request
	rjson, _ := json.MarshalIndent(r.URL, "", "\t")

	fmt.Fprintf(w, string(rjson))
	log.Printf("Request %s", rjson)

	query_params := r.URL.Query()
	log.Println(query_params)

	// Parse the REST API command
	//siegeCmd := parseCommand(r.URL.Query())

	var check = func(s string) string {
		val, ok := query_params[s]
		if ok {
			return val[0]
		}
		log.Panic(s, " could not be read. error :")
		return ""
	}

	sessId := check("sessId")

	cmd := state.StopSiegeSession{
		SessionId: sessId,
	}

	siegeCmd := state.SessionEvent{
		Event: cmd,
	}

	// Write
	writeToState(siegeCmd)
}

// http handler
func newSessHandler(w http.ResponseWriter, r *http.Request) {
	// error handling separated from code flow for easy readability
	// Follow this idiom
	err := func() error {
		if r.Method != "GET" {
			return errors.New("expected GET")
		}

		//if input := parseInput(r); input != "command" {
		//return errors.New("malformed command")
		//}
		return nil
	}()

	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "error found")
		return
	}

	// Log the request
	rjson, _ := json.MarshalIndent(r.URL, "", "\t")

	admin := r.URL.Path[1:]
	cmd := r.URL.Path[2:]
	fmt.Fprintf(w, "Hi there, thanks for the command %s/%s!\n", admin, cmd)
	fmt.Fprintf(w, string(rjson))

	log.Printf("Request %s", string(rjson))

	query_params := r.URL.Query()
	log.Println(query_params)

	// Parse the REST API command
	siegeCmd := parseCommand(r.URL.Query())

	// Write
	writeToState(siegeCmd)

	//out, err := exec.Command("/usr/bin/siege", "--delay="+delay[0], "--concurrent="+concurrent[0], "http://"+host[0]+":"+port[0]+"/"+r.URL.Path[1:]).Output()

	//if err != nil {
	//log.Println("Error occurred")
	//log.Printf("%s", err)
	//}

	//log.Printf("%s", out)

}

func parseCommand(q url.Values) state.SessionEvent {

	var check = func(s string) string {
		val, ok := q[s]
		if ok {
			return val[0]
		}
		log.Fatal(s, " could not be read. error :")
		return ""
	}

	concurrent := check("concurrent")
	delay := check("delay")
	host := check("target")

	cmd := state.NewSiegeSession{
		Concurrent: concurrent,
		Delay:      delay,
		Host:       host,
	}

	return state.SessionEvent{
		Event: cmd,
	}
}

func writeToState(cmd state.SessionEvent) {
	// this is temp
	tempWriteCh <- cmd
}
