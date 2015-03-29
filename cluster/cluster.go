// Package clustermanager provides a routine that manages the siege cluster
// Listens to a channel for messages
//	- New Session
//	- Stop Session
package cluster

import (
	"log"
	"time"

	"github.com/loadcloud/gosiege/common"
	"github.com/loadcloud/gosiege/state"
)

// StartClusterManager starts the siege cluster manager. Takes a
// channel for listening to abort signals
func StartClusterManager() {

	// If there is a panic recover using this function
	defer func() {
		if err := recover(); err != nil {
			log.Println("StartClusterManager - Failed : ", err)
		}
	}()

	log.Println("StartingClusterManager")

	// Subscribe to the State Watcher for Cluster Administration Events
	listenCh := state.SubscribeToClusterEvents()

	listenToIncomingEvents(listenCh)
}

// listens to incoming events on the channel
func listenToIncomingEvents(listen chan state.ClusterEvent) {

	var cmd state.ClusterEvent

	for {
		log.Println("listening for Incoming events.")

		select {
		case <-time.After(24 * time.Hour):
			// closing release all listeners of doneCh
			close(common.DoneCh)
		case cmd = <-listen:
			parseEvent(cmd)
		case <-common.DoneCh:
			log.Println("Abort Message Received. Exitting")
			return
		}

		log.Println("Message Received :")
	}
}

func parseEvent(c state.ClusterEvent) {
	switch t := c.Cmd.(type) {
	case state.AddNode:
		log.Println("Event = ", t)
		addNode(c.Cmd.(state.AddNode))

	case state.RemoveNode:
		log.Println("Event = ", t)
		removeNode(c.Cmd.(state.RemoveNode))

	default:
		log.Println("Event = ", t)
	}
}

func addNode(cmd state.AddNode) {
}

func removeNode(cmd state.RemoveNode) {
}