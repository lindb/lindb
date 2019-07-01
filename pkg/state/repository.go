package state

import (
	"context"
	"fmt"
)

// global repository for state storage
var repository Repository

// Repository stores state data, such as metadata/config/status/task etc.
type Repository interface {
	// Get retrieves value for given key from repository
	Get(ctx context.Context, key string) ([]byte, error)
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
	Watch(ctx context.Context, key string) (WatchEventChan, error)
	// WatchPrefix watches on a prefix.All of the changes who has the prefix
	// will be notified through the WatchEventChan channel.
	WatchPrefix(ctx context.Context, prefixKey string) (WatchEventChan, error)
	// DeleteWithValue deletes the key with the value.it will returns success
	// if the value of the key in the etcd equals the incoming value
	DeleteWithValue(ctx context.Context, key string, value []byte) (bool, error)
	// Close closes repository and release resources
	Close() error
}

// EventType represents a watch event type.
type EventType int

// Event types.
const (
	EventTypeModify EventType = iota
	EventTypeDelete
)

// Event defines repository watch event on key or perfix
type Event struct {
	Type  EventType
	Key   string
	Value []byte

	Err error
}

// Closed represents close status
type Closed interface {
}

// WatchEventChan notify event channel
type WatchEventChan <-chan *Event

// New creates global state reposistory
func New(repoType string, config interface{}) error {
	if repoType == "etcd" {
		repo, err := newEtedRepository(config)
		if err != nil {
			return err
		}
		repository = repo
		return nil
	}
	return fmt.Errorf("repo type not define, type is:%s", repoType)
}

// GetRepo returns state repository
func GetRepo() Repository {
	return repository
}
