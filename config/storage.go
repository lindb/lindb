package config

import (
	"path/filepath"
)

// Storage represents a storage configuration
type StorageKernel struct {
	Coordinator RepoState   `toml:"coordinator"`
	GRPC        GRPC        `toml:"grpc"`
	Engine      Engine      `toml:"engine"`
	Replication Replication `toml:"replication"`
	Query       Query       `toml:"query"`
}

// Storage represents a storage configuration with common settings
type Storage struct {
	StorageKernel
	Logging Logging `toml:"logging"`
}

// GRPC represents grpc server config
type GRPC struct {
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
			GRPC: GRPC{
				Port: 2891,
				TTL:  1},
			Engine: Engine{
				Dir: filepath.Join(defaultParentDir, "storage/data")},
			Replication: Replication{
				Dir: filepath.Join(defaultParentDir, "storage/replication")},
			Query: NewDefaultQueryCfg(),
		},
		Logging: NewDefaultLoggingCfg(),
	}
}
