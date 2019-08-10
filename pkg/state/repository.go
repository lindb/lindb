package state

import (
	"context"
	"errors"
)

//go:generate mockgen -source=./repository.go -destination=./repository_mock.go -package=state

var (
	// ErrNotExist represents key not exist
	ErrNotExist = errors.New("not exist")
)

// RepositoryFactory represents the repository create factory
type RepositoryFactory interface {
	// CreateRepo creates state repository based on config
	CreateRepo(config Config) (Repository, error)
}

// Repository stores state data, such as metadata/config/status/task etc.
type Repository interface {
	// Get retrieves value for given key from repository
	Get(ctx context.Context, key string) ([]byte, error)
	// List retrieves list for given prefix from repository
	List(ctx context.Context, prefix string) ([][]byte, error)
	// Put puts a key-value pair into repository
	Put(ctx context.Context, key string, val []byte) error
	// Delete deletes value for given key from repository
	Delete(ctx context.Context, key string) error
	// Heartbeat does heartbeat on the key with a value and ttl
	Heartbeat(ctx context.Context, key string, value []byte, ttl int64) (<-chan Closed, error)
	// PutIfNotExist puts a key with a value,
	// 1) returns success if the key does not exist and puts success
	// 2) returns failure if key exist
	// When this operation success, it will do keepalive background for keep session
	PutIfNotExist(ctx context.Context, key string, value []byte, ttl int64) (bool, <-chan Closed, error)
	// Watch watches on a key. The watched events will be returned through the returned channel.
	Watch(ctx context.Context, key string) WatchEventChan
	// WatchPrefix watches on a prefix.All of the changes who has the prefix
	// will be notified through the WatchEventChan channel.
	WatchPrefix(ctx context.Context, prefixKey string) WatchEventChan
	// Batch puts k/v list, this operation is atomic
	Batch(ctx context.Context, batch Batch) (bool, error)
	// NewTransaction creates a new transaction
	NewTransaction() Transaction
	// Commit commits the transaction, if fail return err
	Commit(ctx context.Context, txn Transaction) error
	// Close closes repository and release resources
	Close() error
}

// EventType represents a watch event type.
type EventType int

// Event types.
const (
	EventTypeModify EventType = iota
	EventTypeDelete
	EventTypeAll
)

// KeyValue represents key/value pair
type KeyValue struct {
	Key   string
	Value []byte
}

// Batch represents put list for batch operation
type Batch struct {
	KVs []KeyValue
}

// EventKeyValue represents task event
type EventKeyValue struct {
	Key   string
	Value []byte
	Rev   int64
}

// Event defines repository watch event on key or perfix
type Event struct {
	Type      EventType
	KeyValues []EventKeyValue

	Err error
}

// Closed represents close status
type Closed struct{}

// WatchEventChan notify event channel
type WatchEventChan <-chan *Event

// repositoryFactory represents a repository create factory
type repositoryFactory struct {
}

// NewRepositoryFactory creates a repository factory
func NewRepositoryFactory() RepositoryFactory {
	return &repositoryFactory{}
}

// CreateRepo creates state repository based on config
func (f *repositoryFactory) CreateRepo(config Config) (Repository, error) {
	return newEtedRepository(config)
}

type Transaction interface {
	ModRevisionCmp(key, op string, v interface{})
	Put(key string, value []byte)
	Delete(key string)
}
