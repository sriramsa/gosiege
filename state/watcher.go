// Provides a watcher to watch for state changes and inform
// subscribers.
// State is organized as a set of events that can be subscribed to
package state

import (
	"log"
	"time"

	"github.com/loadcloud/gosiege/common"
)

type EventType int

// Events that a command can subscribe to
const (
	EventClusterAdminCmd = iota
	EventSessionAdminCmd
	EventSessionState
)

// Subscriber for a particular session type
type SessionInstanceSubscriber struct {
	subs      chan struct{}
	sessionId string
}

// Subscriber list holders
var subsList map[string]SessionInstanceSubscriber

// Cluster Manager is the only known subscriber
var clusterEventSubscriber = make([]chan ClusterEvent, 0)

// Start with 5, slice will grow as needed.
var sessionEventSubscriber = make([]chan SessionEvent, 0)

// Start the State Watcher
func StartStateWatcher() {

	log.Println("Started")
	// poll every second and look for state changes
	for {
		select {
		case <-time.After(1 * time.Second):
			// Read the etcd store and see if there are any changes

		case <-common.DoneCh:
		}
	}
}

// Components subscribe to Cluster Related Events
// Returns a channel for the caller to listen to
func SubscribeToClusterEvents() (listen chan ClusterEvent) {

	log.Println("Subscription Requested for ClusterEvents")
	listen = make(chan ClusterEvent)

	// Add to the subscriber list for cluster events
	clusterEventSubscriber = append(clusterEventSubscriber, listen)

	return listen
}

// Subscribe to Session Administration related events
// Events include
//	Create a new Session
//	Stop a session
//	Update a session
func SubscribeToSessionEvents() (listen chan SessionEvent) {
	log.Println("Subscription Requested for SessionEvents")
	listen = make(chan SessionEvent)

	// Add to the subscriber list for session events
	sessionEventSubscriber = append(sessionEventSubscriber, listen)

	return make(chan SessionEvent)
}

// Subscribe to the session load channel for a particular session
// This is the channel that session handlers talk on.
func SubscribeToSessionLoadBalanceEvents(sessId string) (listen chan struct{}) {
	log.Println("Subscription Requested for LoadBalanceEvents")
	return make(chan struct{})
}
