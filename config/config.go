// Package config provides configuration needed for initializing the cluster
package config

import (
	"fmt"
	"log"
)

// local variable with configuration loaded
var cfg *ClusterConfig

type Machine struct {
	HostName  string
	IpAddress string
}

// ClusterConfig struct that contains the cluster configuration information
type ClusterConfig struct {
	HeartbeatIntervalInSeconds string
	ListeningPort              string
	SiegePath                  string

	// Machines in the cluster
	ClusterMachines []Machine

	// Maximum number of resources to use
	MaxProcsToUse  string
	MaxMemoryToUse string

	DistributedStateProvider string

	StressEngine string
}

type SessionConfig struct {
	UsersPerSiegeProc int
}

func Get(s string) string {
	switch s {
	case "ListeningPort":
		return cfg.ListeningPort
	case "SiegePath":
		return cfg.SiegePath
	default:
		// TODO: fix
		return ""
	}
}

// Return the key value store option from the configuration
func KeyValueStoreOption() string {
	return cfg.DistributedStateProvider
}

var event *testrument.EventWriter

// Gets the config from the environment and returns the same
func LoadConfig() error {
	event = testrument.NewEventWriter("config", nil, true)
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
		HeartbeatIntervalInSeconds: "5",
		ListeningPort:              "8090",
		MaxProcsToUse:              "1",

		DistributedStateProvider: "etcd",
		StressEngine:             "siege",
		SiegePath:                "/gosiege/",
	}

	log.Println("Config loaded :", config)
	event.Info("Configuration Loaded : ", config)
	return config
}
