package kv

// Merger does merge job(compact/rollup etc.)
type Merger interface {

	// Merge values for same key
	Merge(key uint32, value [][]byte)
}
