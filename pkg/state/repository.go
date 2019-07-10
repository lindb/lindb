package state

import (
	"context"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
)

const (
	// ETCD defines state store using ectd cluster
	ETCD string = "etcd"
)

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
	// PutIfNotExist  puts a key with a value.it will be success
	// if the key does not exist,otherwise it will be failed.When this
	// operation success,it will do keepalive background
	PutIfNotExist(ctx context.Context, key string, value []byte, ttl int64) (bool, <-chan Closed, error)
	// Watch watches on a key. The watched events will be returned through the returned channel.
	Watch(ctx context.Context, key string) WatchEventChan
	// WatchPrefix watches on a prefix.All of the changes who has the prefix
	// will be notified through the WatchEventChan channel.
	WatchPrefix(ctx context.Context, prefixKey string) WatchEventChan
	// Txn returns a etcdcliv3.Txn.
	Txn(ctx context.Context) etcdcliv3.Txn
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

// NewRepo create state repository based on config
func NewRepo(config Config) (Repository, error) {
	return newEtedRepository(config)
}
