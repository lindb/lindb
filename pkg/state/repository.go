// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package state

import (
	"context"
	"errors"

	"github.com/lindb/lindb/config"
)

//go:generate mockgen -source=./repository.go -destination=./repository_mock.go -package=state

var (
	// ErrNotExist represents key not exist
	ErrNotExist = errors.New("not exist")
)

// RepositoryFactory represents the repository create factory
type RepositoryFactory interface {
	// CreateRepo creates state repository based on config
	CreateRepo(repoState config.RepoState) (Repository, error)
}

// Repository stores state data, such as metadata/config/status/task etc.
type Repository interface {
	// Get retrieves value for given key from repository
	Get(ctx context.Context, key string) ([]byte, error)
	// List retrieves list for given prefix from repository
	List(ctx context.Context, prefix string) ([]KeyValue, error)
	// Put puts a key-value pair into repository
	Put(ctx context.Context, key string, val []byte) error
	// Delete deletes value for given key from repository
	Delete(ctx context.Context, key string) error
	// Heartbeat does heartbeat on the key with a value and ttl
	Heartbeat(ctx context.Context, key string, value []byte, ttl int64) (<-chan Closed, error)
	// Elect puts a key with a value,
	// 1) returns success if the key does not exist and puts success
	// 2) returns failure if key exist
	// When this operation success, it will do keepalive background for keep session
	Elect(ctx context.Context, key string, value []byte, ttl int64) (bool, <-chan Closed, error)
	// Watch watches on a key. The watched events will be returned through the returned channel.
	// fetchVal: if fetch prefix key's values for init
	Watch(ctx context.Context, key string, fetchVal bool) WatchEventChan
	// WatchPrefix watches on a prefix.All of the changes who has the prefix
	// will be notified through the WatchEventChan channel.
	// fetchVal: if fetch prefix key's values for init
	WatchPrefix(ctx context.Context, prefixKey string, fetchVal bool) WatchEventChan
	// Batch puts k/v list, this operation is atomic
	Batch(ctx context.Context, batch Batch) (bool, error)
	// NextSequence returns next sequence number.
	NextSequence(ctx context.Context, key string) (int64, error)
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

// String returns event type string value
func (e EventType) String() string {
	switch e {
	case EventTypeModify:
		return "modify"
	case EventTypeDelete:
		return "delete"
	case EventTypeAll:
		return "all"
	default:
		return "unknown"
	}
}

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
	owner string
}

// NewRepositoryFactory creates a repository factory by owner
func NewRepositoryFactory(owner string) RepositoryFactory {
	return &repositoryFactory{owner: owner}
}

// CreateRepo creates state repository based on config
func (f *repositoryFactory) CreateRepo(repoState config.RepoState) (Repository, error) {
	return newEtcdRepository(repoState, f.owner)
}

type Transaction interface {
	ModRevisionCmp(key, op string, v interface{})
	Put(key string, value []byte)
	Delete(key string)
}
