// Package main provides listener that listens to the http port for incoming
// messages from Admin NodeJS UI or Command Line
package cluster

import (
	"fmt"
	"log"
	"net/http"
)

var writeCh chan Command

// Starts a http listener and reports incoming messages to the caller
func StartGoSiegeHttpListener(writeCh chan Command, doneCh chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Listener failed", err)
		}
	}()

	writeCh = writeCh
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	fmt.Println("hit")

	query_params := r.URL.Query()
	concurrent, _ := query_params["concurrent"]
	delay, _ := query_params["delay"]
	host, _ := query_params["target"]

	cmd := CreateNewSession{
		concurrent: concurrent[0],
		delay:      delay[0],
		host:       host[0],
	}

	// Write the command into the channel
	writeCh <- Command{
		cmd: cmd,
	}

	fmt.Println(query_params)
	//out, err := exec.Command("/usr/bin/siege", "--delay="+delay[0], "--concurrent="+concurrent[0], "http://"+host[0]+":"+port[0]+"/"+r.URL.Path[1:]).Output()

	//if err != nil {
	//fmt.Println("Error occurred")
	//fmt.Printf("%s", err)
	//}

	//fmt.Printf("%s", out)

}
