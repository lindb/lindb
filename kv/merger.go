package kv

//go:generate mockgen -source ./merger.go -destination=./merger_mock.go -package kv

// MergerType represents the merger type
type MergerType string

// NewMerger represents create merger instance function
type NewMerger func() Merger

var mergers = make(map[MergerType]NewMerger)

// RegisterMerger registers family merger
// NOTICE: must register before create family
func RegisterMerger(name MergerType, merger NewMerger) {
	_, ok := mergers[name]
	if ok {
		panic("merger already register")
	}
	mergers[name] = merger
}

// Merger represents merger values of same key when do compaction job(compact/rollup etc.)
type Merger interface {
	// Init initializes merger params or context, before does merge operation
	Init(params map[string]interface{})
	// Merge merges values for same key, return merged value or err if failure
	Merge(key uint32, values [][]byte) ([]byte, error)
}
