package state

// Config represents state repository config
type Config struct {
	Namespace   string   `toml:"namespace" json:"namespace"`
	Endpoints   []string `toml:"endpoints" json:"endpoints"`
	DialTimeout int64    `toml:"dialTimeout" json:"dialTimeout"`
}
