package config

import (
	"path/filepath"
)

// Storage represents a storage configuration
type StorageKernel struct {
	Coordinator RepoState   `toml:"coordinator"`
	Server      Server      `toml:"server"`
	Engine      Engine      `toml:"engine"`
	Replication Replication `toml:"replication"`
}

// Storage represents a storage configuration with common settings
type Storage struct {
	StorageKernel
	Logging Logging `toml:"logging"`
}

// Server represents tcp server config
type Server struct {
	Port uint16 `toml:"port"`
	TTL  int64  `toml:"ttl"`
}

// Replication represents replication config
type Replication struct {
	Dir string `toml:"path"`
}

// Engine represents a tsdb engine level configuration
type Engine struct {
	Dir string `toml:"path"`
}

// NewDefaultStorageCfg creates storage define config
func NewDefaultStorageCfg() Storage {
	return Storage{
		StorageKernel: StorageKernel{
			Coordinator: RepoState{
				Namespace:   "/lindb/storage",
				Endpoints:   []string{"http://localhost:2379"},
				DialTimeout: 5},
			Server: Server{
				Port: 2891,
				TTL:  1},
			Engine: Engine{
				Dir: filepath.Join(defaultParentDir, "storage/data")},
			Replication: Replication{
				Dir: filepath.Join(defaultParentDir, "storage/replication")}},
		Logging: NewDefaultLoggingCfg(),
	}
}
