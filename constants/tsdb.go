package constants

const (
	// use this limit of metric-store when maxTagsLimit is not set
	DefaultMStoreMaxTagsCount = 100000
	// max tag keys limitation of a metric-store
	MStoreMaxTagKeysCount = 512
	// max fields limitation of a tsStore.
	TStoreMaxFieldsCount = 1024
	// the max number of suggestions count
	MaxSuggestions = 10000
)
