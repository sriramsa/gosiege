// Package config provides configuration needed for initializing the cluster
package config

import "github.com/sriramsa/gosiege/logger"

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

var config *ClusterConfig

// Gets the config from the environment and returns the same
func LoadConfig() (config *ClusterConfig) {
	config = loadConfig()

	return config
}

func loadConfig() (config *ClusterConfig) {
	config = &ClusterConfig{
		HeartbeatIntervalInSeconds: 5,
		ListeningPort:              0,

		DistributedStateProvider: "etcd",
		StressEngine:             "siege",
	}

	log.Println("Config loaded :", config)
	return config
}
