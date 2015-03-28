// Maintains and manages distributed state
package state

func InitGoSiegeState() (err error) {

	return nil
}

func ReadValue(key string) (value []byte, err error) {

	// Read from the distributed store
	return nil, nil
}

func WriteValue(key string, value []byte) (err error) {
	return nil
}

// Components subscribe to state using this API
// Returns a channel for the caller to listen to
func SubscribeToTopic(t string) (listen chan struct{}) {

	listen = make(chan struct{})

	return
}
