package config

// Standalone represents the configuration of standalone mode
type Standalone struct {
	ETCD    ETCD    `toml:"etcd"`
	Broker  Broker  `toml:"broker"`
	Storage Storage `toml:"storage"`
}

// ETCD represents embed etcd's configuration
type ETCD struct {
	Dir string `toml:"dir"`
	URL string `toml:"url"`
}

// NewDefaultStandaloneCfg creates define config of standalone mode
func NewDefaultStandaloneCfg() Standalone {
	return Standalone{
		Broker:  NewDefaultBrokerCfg(),
		Storage: NewDefaultStorageCfg(),
		ETCD:    ETCD{Dir: "/tmp/lindb/standalone", URL: "http://localhost:2379"},
	}
}
