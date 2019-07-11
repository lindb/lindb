package config

import "github.com/eleme/lindb/pkg/state"

// Storage represents a storage configuration
type Storage struct {
	Coordinator state.Config `toml:"coordinator"`
	Server      Server       `toml:"server"`

	Engine `toml:"engine"`
}

// Server represents tcp server config
type Server struct {
	Port uint16 `toml:"port"`
	TTL  int64  `toml:"ttl"`
}

// Engine represents an engine level configuration
type Engine struct {
	Path string `toml:"path"`
	Name string `toml:"name"`
}
