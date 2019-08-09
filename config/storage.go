package config

import "github.com/lindb/lindb/pkg/state"

// Storage represents a storage configuration
type Storage struct {
	Coordinator state.Config `toml:"coordinator"`
	Server      Server       `toml:"server"`

	Engine      Engine      `toml:"engine"`
	Replication Replication `toml:"replication"`
}

// Server represents tcp server config
type Server struct {
	Port uint16 `toml:"port"`
	TTL  int64  `toml:"ttl"`
}

// Replication represents replication config
type Replication struct {
	Path string `toml:"path"`
}

// Engine represents a tsdb engine level configuration
type Engine struct {
	Path string `toml:"path"`
}

// NewDefaultStorageCfg creates storage define config
func NewDefaultStorageCfg() Storage {
	return Storage{
		Coordinator: state.Config{
			Namespace:   "/lindb/storage",
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5,
		},
		Server: Server{
			Port: 2891,
			TTL:  1,
		},
		Engine: Engine{
			Path: "/tmp/storage/data",
		},
		Replication: Replication{
			Path: "/tmp/storage/replication",
		},
	}
}
