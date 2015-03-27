// Maintains and manages distributed state
package state

import "github.com/sriramsa/gosiege/config"

func InitGoSiegeState() (err error) {
	// Read configuration
	_ = config.LoadConfig()

	return nil
	// load plugins
}

func ReadValue(key string) (value []byte, err error) {

	// Read from the distributed store
	return nil, nil
}

func WriteValue(key string, value []byte) (err error) {
	return nil
}
