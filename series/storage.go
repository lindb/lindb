package series

// Storage represents the time series data storage interface
type Storage interface {
	// Interval returns the interval of storage
	Interval() int64
}
