package config

// Storage represents a storeage configuration
type Storage struct {
	Engine `toml:"engine"`
}

// Engine represetns an engine level configuration
type Engine struct {
	Path string `toml:"path"`
}
