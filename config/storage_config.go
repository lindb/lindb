package config

// StorageConfig represents a storage configuration
type StorageConfig struct {
	Engine `toml:"engine"`
}

// Engine represents an engine level configuration
type Engine struct {
	Path string `toml:"path"`
	Name string `toml:"name"`
}
