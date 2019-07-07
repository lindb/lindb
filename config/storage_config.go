package config

import "github.com/eleme/lindb/pkg/state"

// StorageConfig represents a storage configuration
type StorageConfig struct {
	StorageRepositoryConfig `toml:"StorageRepositoryConfig"`
	Engine                  `toml:"engine"`
	StoragePort             uint16 `toml:"Port"`
}

// Engine represents an engine level configuration
type Engine struct {
	Path string `toml:"path"`
	Name string `toml:"name"`
}

// RepositoryConfig represents the repository config
type StorageRepositoryConfig struct {
	state.RepositoryConfig `toml:"RepositoryConfig"`
	HeartBeatTTL           int64  `toml:"HeartBeatTTL"`
	HeartBeatPrefix        string `toml:"HeartBeatPrefix"`
}
