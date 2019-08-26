package config

// RepoState represents state repository config
type RepoState struct {
	Namespace   string   `toml:"namespace" json:"namespace"`
	Endpoints   []string `toml:"endpoints" json:"endpoints"`
	DialTimeout int64    `toml:"dialTimeout" json:"dialTimeout"`
}

// StorageCluster represents config of storage cluster
type StorageCluster struct {
	Name   string    `json:"name"`
	Config RepoState `json:"config"`
}
