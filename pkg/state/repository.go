package state

import (
	"context"
	"fmt"
	"strings"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
)

// global repository for state storage
var repository Repository

const (
	ETCD string = "etcd"
)

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
	Watch(ctx context.Context, key string) WatchEventChan
	// WatchPrefix watches on a prefix.All of the changes who has the prefix
	// will be notified through the WatchEventChan channel.
	WatchPrefix(ctx context.Context, prefixKey string) WatchEventChan
	// DeleteWithValue deletes the key with the value. It will returns nil
	// if the value of the key in the etcd equals the incoming value
	DeleteWithValue(ctx context.Context, key string, value []byte) error
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

// New creates global state repository
func New(repoType string, config interface{}) error {
	if repoType == ETCD {
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

// the custom repository config
type RepositoryConfig struct {
	RepositoryType string
	URL            string
}

// convert custom config to real repository config
func NewRepositoryConfig(config RepositoryConfig) (interface{}, error) {
	switch config.RepositoryType {
	case ETCD:
		return &etcdcliv3.Config{
			Endpoints: strings.Split(config.URL, ","),
		}, nil
	default:
		return nil, fmt.Errorf("repo type not support, type is:%s", config.RepositoryType)
	}

}
