package config

import (
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
)

// Broker represents a broker configuration
type Broker struct {
	HTTP               HTTP               `toml:"HTTP"`
	Coordinator        state.Config       `toml:"coordinator"`
	User               models.User        `toml:"user"`
	Server             Server             `toml:"server"`
	ReplicationChannel ReplicationChannel `toml:"replicationChannel"`
}

// HTTP represents a HTTP level configuration of broker.
type HTTP struct {
	Port uint16 `toml:"port"`
}

// ReplicationChannel represents config for data replication in broker.
type ReplicationChannel struct {
	Path                       string `toml:"path"`
	BufferSize                 int    `toml:"bufferSize"`
	SegmentFileSize            int    `toml:"segmentFileSize"`
	RemoveTaskIntervalInSecond int    `toml:"remoteTaskIntervalInSecond"`
}

// NewDefaultBrokerCfg creates broker default config
func NewDefaultBrokerCfg() Broker {
	return Broker{
		HTTP: HTTP{
			Port: 9000,
		},
		Coordinator: state.Config{
			Namespace:   "/lindb/broker",
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5,
		},
		User: models.User{
			UserName: "admin",
			Password: "admin123",
		},
		ReplicationChannel: ReplicationChannel{
			Path:                       "/tmp/broker/replication",
			BufferSize:                 32,
			SegmentFileSize:            128 * 1024 * 1024,
			RemoveTaskIntervalInSecond: 60,
		},
	}
}
