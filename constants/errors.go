package constants

import "errors"

var (
	ErrDatabaseNotFound = errors.New("database not found")
	ErrShardNotFound    = errors.New("shard not found")

	// ErrNotFound represents the data not found
	ErrNotFound = errors.New("not found")
)
