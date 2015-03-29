// Provides functions to CRUD on state elements
package state

import "github.com/loadcloud/gosiege/config"

// Initialize the Gosiege State
func InitGoSiegeState() (err error) {

	// Read the config and initialize the appropriate plugin
	// If etcd, use etcd plugin
	// If only one in the cluster, use local file?
	_ = config.KeyValueStoreOption()
	return nil
}

func ReadValue(key string) (value []byte, err error) {

	// Read from the distributed store
	return nil, nil
}

func WriteValue(key string, value []byte) (err error) {
	return nil
}
