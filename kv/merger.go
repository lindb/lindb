package kv

//go:generate mockgen -source ./merger.go -destination=./merger_mock.go -package kv

// Merger represents merger values of same key when do compaction job(compact/rollup etc.)
type Merger interface {
	// Merge merges values for same key, return merged value or err if failure
	Merge(key uint32, value [][]byte) ([]byte, error)
}
