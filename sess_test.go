// Test a session
package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
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

	// Start a test http server to listen
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("HIT HIT HIT HIT")
	}))
	defer ts.Close()
	log.Println("Startng a test http server to recieve hits at : ", ts.URL)

	// Get port number
	re := regexp.MustCompile(`.*:(.*)$`)
	match := re.FindStringSubmatch(ts.URL)
	port := match[1]

	// Start the main go siege server
	log.Println("Starting the Go Siege Server in a separate routine.")
	startServer()
	// TODO: Fix waiting
	time.Sleep(time.Second * 2)

	// Create a session on the server
	log.Println("Sending command to create a session on the server.")

	form := url.Values{}
	form.Add("concurrent", "7")
	form.Add("target", "localhost")
	form.Add("port", port)
	form.Add("delay", "1")
	req, err := http.NewRequest("PUT", "http://localhost:8090/gosiege/sessions/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Error("Creating Request failed with error : ", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Request failed with error : ", err)
	}

	log.Println("Server response : ", resp)

	// Ensure we are getting hits for 5 secs
	log.Println("TBD: Checking if we are getting hits from the server for 5 seconds")

	// Stop the Server
	log.Println("TBD: Stopping the session.")

	// Ensure we are not getting any hits
	log.Println("TBD: Verifying that the hits have stopped.")

	time.Sleep(time.Second * 3)
	//t.Error("Creating a session failed.")
	// Wait for a keystroke to exit.
}
