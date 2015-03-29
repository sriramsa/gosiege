// Package config provides configuration needed for initializing the cluster
package config

import (
	"fmt"

	"github.com/loadcloud/gosiege/logger"
)

// local variable with configuration loaded
var cfg *ClusterConfig

var log = logger.NewLogger("Config")

type Machine struct {
	HostName  string
	IpAddress string
}

// ClusterConfig struct that contains the cluster configuration information
type ClusterConfig struct {
	HeartbeatIntervalInSeconds uint
	ListeningPort              uint

	// Machines in the cluster
	ClusterMachines []Machine

	// Maximum number of resources to use
	MaxProcsToUse  int
	MaxMemoryToUse int

	DistributedStateProvider string

	StressEngine string
}

// Return the key value store option from the configuration
func KeyValueStoreOption() string {
	return cfg.DistributedStateProvider
}

// Gets the config from the environment and returns the same
func LoadConfig() error {
	cfg = loadConfig()

	if err := validateConfig(); err != nil {
		return fmt.Errorf("Couldn't load config:", err)
	}

	return nil
}

func validateConfig() error {
	return nil
}
func loadConfig() (config *ClusterConfig) {
	config = &ClusterConfig{
		HeartbeatIntervalInSeconds: 5,
		ListeningPort:              0,
		MaxProcsToUse:              1,

		DistributedStateProvider: "etcd",
		StressEngine:             "siege",
	}

	log.Println("Config loaded :", config)
	return config
}
