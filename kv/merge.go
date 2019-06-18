package kv

// Merge iterface do merge job(compact/rollup etc.)
type Merge interface {

	// Merge values for same key
	Merge(key uint32, value [][]byte)
}
