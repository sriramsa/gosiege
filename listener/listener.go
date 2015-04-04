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

// Starts a http listener and reports incoming messages to the caller
func StartHttpCommandListener() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("FATAL : Listener failed", err)
		}

		// Let everybody exit since we can't listen for incoming commands
		close(common.DoneCh)
	}()

	port := config.Get("ListeningPort")
	urlPrefix = config.Get("SiegePath")

	http.HandleFunc(urlPrefix, handler)
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

// http handler
func handler(w http.ResponseWriter, r *http.Request) {
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
	log.Println(string(rjson))

	admin := r.URL.Path[1:]
	cmd := r.URL.Path[2:]
	fmt.Fprintf(w, "Hi there, thanks for the command %s/%s!\n", admin, cmd)
	fmt.Fprintf(w, string(rjson))
	log.Printf("Request %s", rjson)

	query_params := r.URL.Query()
	//concurrent, _ := query_params["concurrent"]
	//delay, _ := query_params["delay"]
	//host, _ := query_params["target"]

	log.Println(query_params)

	// Parse the REST API command
	siegeCmd := parseCommand(r.URL.Query())

	// Write
	writeToState()

	//out, err := exec.Command("/usr/bin/siege", "--delay="+delay[0], "--concurrent="+concurrent[0], "http://"+host[0]+":"+port[0]+"/"+r.URL.Path[1:]).Output()

	//if err != nil {
	//log.Println("Error occurred")
	//log.Printf("%s", err)
	//}

	//log.Printf("%s", out)

}

func parseCommand(q url.Values) interface{} {

	cmd := state.NewSiegeSession{
		Concurrent, _: q["concurrent"],
		Delay, _: q["delay"],
		Host: q["target"],
	}

	return cmd
}

func writeToState() {

}
