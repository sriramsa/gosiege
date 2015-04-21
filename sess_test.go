// Test a session.
package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/loadcloud/gosiege/instrument"
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
func TestCreateAndStopASession(t *testing.T) {
	t.Log("TEST: Creating a Session")

	var hitCounter uint32 = 0
	checkHit := true
	hitRecd := make(chan struct{})
	noHitRecd := make(chan struct{})
	var verifyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment hit counter
		atomic.AddUint32(&hitCounter, 1)
		if hitCounter == 10 && checkHit {
			t.Log("More than 10 hits received...")
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
	t.Log("Started a test http server to recieve hits : ", ts.URL)

	// Start the main go siege server
	t.Log("Starting the Go Siege Server in a separate routine.")
	startServer()

	// Attach to the event stream and wait for ServerInitDone event
	mainEvents := instrument.NewAttach(GetEventHydrant())
	jsonObj, e := mainEvents.WaitForEvent(instrument.Info, "ServerInitDone", time.Second*1)
	if e != nil {
		t.Error("Could not get event: ", e)
	} else {
		t.Log("JSON: ", jsonObj)
	}

	// Got to detach for now since log write is synchronous and blocks if
	// nobody is listening on it.
	mainEvents.Detach()

	time.Sleep(time.Second * 1)

	// Create a session on the server
	t.Log("Sending command to create a session on the server.")

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

	t.Log("Server response : ", resp)

	// Wait till we get some hits
	t.Log("Waiting for hits on test server.")
	<-hitRecd

	// Stop the Server
	t.Log("Stopping the session.")
	sr := test.StopSessionReq{
		SessionId: "101",
	}

	resp, err = sr.Send()
	if err != nil {
		t.Error("Stop Session request failed with error : ", err)
	}
	t.Log("Server response : ", resp)

	// Ensure we are not getting any hits
	t.Log("Verifying that the hits have stopped.")
	select {
	case <-noHitRecd:
		t.Log("Verified no hits.")
	case <-time.After(time.Second * 3):
		t.Error("Hits didn't stop for 3 seconds after sending stop.")
	}
	time.Sleep(time.Second * 2)
}
