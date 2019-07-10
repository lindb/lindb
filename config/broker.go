package config

import (
	"github.com/eleme/lindb/pkg/state"
)

// Broker represents a broker configuration
type Broker struct {
	HTTP        HTTP         `toml:"HTTP"`
	Coordinator state.Config `toml:"coordinator"`
}

// HTTP represents an HTTP level configuration of broker.
type HTTP struct {
	Port uint16 `toml:"port"`
}
