package config

import (
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
)

// Broker represents a broker configuration
type Broker struct {
	HTTP        HTTP         `toml:"HTTP"`
	Coordinator state.Config `toml:"coordinator"`
	User        models.User  `toml:"user"`
}

// HTTP represents an HTTP level configuration of broker.
type HTTP struct {
	Port uint16 `toml:"port"`
}

// NewDefaultBrokerCfg creates broker default config
func NewDefaultBrokerCfg() Broker {
	return Broker{
		HTTP: HTTP{
			Port: 9000,
		},
		Coordinator: state.Config{
			Namespace:   "/lindb/broker",
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 5,
		},
	}
}
