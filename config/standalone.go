package config

import (
	"path/filepath"
)

// Standalone represents the configuration of standalone mode
type Standalone struct {
	ETCD    ETCD          `toml:"etcd"`
	Broker  BrokerKernel  `toml:"broker"`
	Storage StorageKernel `toml:"storage"`
	Logging Logging       `toml:"logging"`
}

// ETCD represents embed etcd's configuration
type ETCD struct {
	Dir string `toml:"dir"`
	URL string `toml:"url"`
}

// NewDefaultStandaloneCfg creates define config of standalone mode
func NewDefaultStandaloneCfg() Standalone {
	return Standalone{
		Broker:  NewDefaultBrokerCfg().BrokerKernel,
		Storage: NewDefaultStorageCfg().StorageKernel,
		Logging: NewDefaultLoggingCfg(),
		ETCD: ETCD{
			Dir: filepath.Join(defaultParentDir, "standalone"),
			URL: "http://localhost:2379"},
	}
}
