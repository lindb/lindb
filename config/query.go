package config

// Query represents query rpc config
type Query struct {
	NumOfTasks    int   `toml:"numOfTasks"`    // task pool size
	QueueCapacity int   `toml:"queueCapacity"` // task pool queue's capacity
	Timeout       int64 `toml:"timeout"`       // task process timeout, the number of second
}

// NewDefaultQueryCfg creates the default query config
func NewDefaultQueryCfg() Query {
	return Query{
		NumOfTasks:    30,
		QueueCapacity: 30,
		Timeout:       30,
	}
}
