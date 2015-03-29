// Package main provides listener that listens to the http port for incoming
// messages from Admin NodeJS UI or Command Line
package listener

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Starts a http listener and reports incoming messages to the caller
func StartHttpCommandListener() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Listener failed", err)
		}
	}()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
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

	fmt.Fprintf(w, "Hi there, thanks for the command %s!", r.URL.Path[1:])
	fmt.Println("hit by %s", r.URL.Path[1:])

	query_params := r.URL.Query()
	//concurrent, _ := query_params["concurrent"]
	//delay, _ := query_params["delay"]
	//host, _ := query_params["target"]

	fmt.Println(query_params)

	// Parse the REST API command
	parseCommand()

	// Write
	writeToState()

	//out, err := exec.Command("/usr/bin/siege", "--delay="+delay[0], "--concurrent="+concurrent[0], "http://"+host[0]+":"+port[0]+"/"+r.URL.Path[1:]).Output()

	//if err != nil {
	//fmt.Println("Error occurred")
	//fmt.Printf("%s", err)
	//}

	//fmt.Printf("%s", out)

}

func parseCommand() {

}

func writeToState() {

}
