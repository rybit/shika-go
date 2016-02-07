package components

import (
	"encoding/json"
)

// Configuration a container for the possible config values
type Configuration struct {
	DefaultPartitionCount int32  `json:"default_partition_count"`
	DataDirectory         string `json:"data_directory"`
	NodePort              int    `json:"node_port"`
	// TODO define routing startegies????
}

func fromString(jsonString string) (*Configuration, error) {
	var config Configuration
	err := json.Unmarshal([]byte(jsonString), &config)
	return &config, err
}
