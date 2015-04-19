// Test a session
package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadcloud/gosiege/test"
)

// Starts the main server to test
func startServer() {
	go main()
}

func stat_handler(w http.ResponseWriter, r *http.Request) {
	log.Println("STAT HIT")
}

// Verifies that a session is created and it sends requests.
//
func TestCreatingASession(t *testing.T) {
	log.Println("TEST: Creating a Session")

	var hitCounter uint32 = 0
	checkHit := true
	hitRecd := make(chan struct{})
	noHitRecd := make(chan struct{})
	var verifyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment hit counter
		atomic.AddUint32(&hitCounter, 1)
		if hitCounter == 10 && checkHit {
			log.Println("More than 10 hits received...")
			close(hitRecd)
			// No need to check hits now. Now check for no hits
			checkHit = false

			// Spin up a timer to set the counter to 0 every second
			// and verify that no hits were obtained for a second
			var tr *time.Timer
			tr = time.AfterFunc(time.Second*1, func() {
				if hitCounter == 0 {
					close(noHitRecd)
				} else {
					atomic.StoreUint32(&hitCounter, 0)
					tr.Reset(time.Second * 1)
				}
			})
		}
	})

	ts, port := test.StartTestServer(verifyHandler)
	defer ts.Close()
	log.Println("Started a test http server to recieve hits : ", ts.URL)

	// Start the main go siege server
	log.Println("Starting the Go Siege Server in a separate routine.")
	startServer()
	// TODO: Fix waiting
	time.Sleep(time.Second * 1)

	// Create a session on the server
	log.Println("Sending command to create a session on the server.")

	req := test.NewSessionReq{
		Url:        ts.URL,
		Target:     "localhost",
		Port:       port,
		Concurrent: "11",
		Delay:      "1",
	}

	resp, err := req.Send()
	if err != nil {
		t.Error("Request failed with error : ", err)
	}

	log.Println("Server response : ", resp)

	// Wait till we get some hits
	log.Println("Waiting for hits on test server.")
	<-hitRecd

	// Stop the Server
	log.Println("Stopping the session.")
	sr := test.StopSessionReq{
		SessionId: "101",
	}

	resp, err = sr.Send()
	if err != nil {
		t.Error("Stop Session request failed with error : ", err)
	}
	log.Println("Server response : ", resp)

	// Ensure we are not getting any hits
	log.Println("Verifying that the hits have stopped.")
	select {
	case <-noHitRecd:
		log.Println("Verified no hits.")
	case <-time.After(time.Second * 3):
		t.Error("Hits didn't stop for 3 seconds after sending stop.")
	}

	//t.Error("Creating a session failed.")
	// Wait for a keystroke to exit.
}
