package config

// BrokerConfig represents a broker configuration
type BrokerConfig struct {
	HTTP HTTP `toml:"HTTP"`
}

// HTTP represents an HTTP level configuration of broker.
type HTTP struct {
	Port   uint16 `toml:"port"`
	Static string `toml:"static"`
}
