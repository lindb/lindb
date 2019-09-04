package config

import (
	"path/filepath"
)

// BrokerKernel represents a broker configuration
type BrokerKernel struct {
	HTTP               HTTP               `toml:"http"`
	Coordinator        RepoState          `toml:"coordinator"`
	User               User               `toml:"user"`
	GRPC               GRPC               `toml:"grpc"`
	TCP                TCP                `toml:"tcp"`
	ReplicationChannel ReplicationChannel `toml:"replicationChannel"`
}

// Broker represents a broker configuration with common settings
type Broker struct {
	BrokerKernel
	Logging Logging `toml:"logging"`
}

// HTTP represents a HTTP level configuration of broker.
type HTTP struct {
	Port uint16 `toml:"port"`
}

// User represents user model
type User struct {
	UserName string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

type TCP struct {
	Port uint16 `toml:"port"`
}

// ReplicationChannel represents config for data replication in broker.
type ReplicationChannel struct {
	Dir                        string `toml:"path"`
	BufferSize                 int    `toml:"bufferSize"`
	SegmentFileSize            int    `toml:"segmentFileSize"`
	RemoveTaskIntervalInSecond int    `toml:"remoteTaskIntervalInSecond"`
	ReportInterval             int    `toml:"reportInterval"` // replicator state report interval
	// interval for check flush
	CheckFlushIntervalInSecond int `toml:"checkFlushIntervalInSecond"`
	// interval for flush
	FlushIntervalInSecond int `toml:"flushIntervalInSecond"`
	//buffer size limit for batch bytes before append to queue
	BufferSizeLimit int `toml:"bufferSizeLimit"`
}

// NewDefaultBrokerCfg creates broker default config
func NewDefaultBrokerCfg() Broker {
	return Broker{
		BrokerKernel: BrokerKernel{
			HTTP: HTTP{
				Port: 9000,
			},
			GRPC: GRPC{
				Port: 9001,
			},
			TCP: TCP{
				Port: 9002,
			},
			Coordinator: RepoState{
				Namespace:   "/lindb/broker",
				Endpoints:   []string{"http://localhost:2379"},
				DialTimeout: 5,
			},
			User: User{
				UserName: "admin",
				Password: "admin123",
			},
			ReplicationChannel: ReplicationChannel{
				Dir:                        filepath.Join(defaultParentDir, "broker/replication"),
				BufferSize:                 32,
				SegmentFileSize:            128 * 1024 * 1024,
				RemoveTaskIntervalInSecond: 60,
				CheckFlushIntervalInSecond: 1,
				FlushIntervalInSecond:      5,
				BufferSizeLimit:            128 * 1024,
			}},
		Logging: NewDefaultLoggingCfg(),
	}
}
