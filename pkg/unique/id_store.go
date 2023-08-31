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

package unique

import (
	"bytes"
	"io"

	"github.com/cockroachdb/pebble"

	"github.com/lindb/common/pkg/logger"
)

//go:generate mockgen -source ./id_store.go -destination=./id_store_mock.go -package=unique

// for testing
var (
	pebbleOpenFn = pebble.Open
)

// IDStore represents unique id store for LinDB's metadata.
type IDStore interface {
	io.Closer

	// Get returns the value by given key, if not exist return false.
	Get(key []byte) ([]byte, bool, error)
	// Put puts the value into store by given key.
	Put(key, val []byte) error
	// Merge puts the value into store by given key, if the key exist, will merge value in background.
	Merge(key, val []byte) error
	// Delete deletes the value by given key.
	Delete(key []byte) error
	// IterKeys iterates the key list by given prefix, returns the key list.
	IterKeys(prefix []byte, limit int) (rs [][]byte, err error)
	// Flush flushes the memory table data under pebble db.
	Flush() error
}

// idStore implements IDStore interface.
type idStore struct {
	db     *pebble.DB
	path   string
	logger logger.Logger
}

// NewIDStore creates an IDStore instance.
func NewIDStore(path string) (IDStore, error) {
	// panic when reopen exist db(https://github.com/cockroachdb/pebble/issues/1777)
	db, err := pebbleOpenFn(path, DefaultOptions())
	if err != nil {
		return nil, err
	}
	return &idStore{
		db:     db,
		path:   path,
		logger: logger.GetLogger("PKG", "IDStore"),
	}, nil
}

// Get returns the value by given key, if not exist return false.
func (s *idStore) Get(key []byte) (value []byte, ok bool, err error) {
	val, closer, err := s.db.Get(key)
	defer func() {
		if closer != nil {
			if err0 := closer.Close(); err0 != nil {
				s.logger.Warn("close kv get resource err",
					logger.String("path", s.path),
					logger.Error(err0))
			}
		}
	}()
	if err == pebble.ErrNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	// NOTICE: MUST COPY IT, thread no safe, maybe panic(unexpected fault address x fatal error: fault)
	// unexpected fault address panic cannot be recovery.
	dst := make([]byte, len(val))
	copy(dst, val)
	return dst, true, nil
}

// Put puts the value into store by given key.
func (s *idStore) Put(key, val []byte) error {
	return s.db.Set(key, val, &pebble.WriteOptions{Sync: false})
}

// Merge puts the value into store by given key, if the key exist, will merge value in background.
func (s *idStore) Merge(key, val []byte) error {
	return s.db.Merge(key, val, &pebble.WriteOptions{Sync: false})
}

// Delete deletes the value by given key.
func (s *idStore) Delete(key []byte) error {
	return s.db.Delete(key, &pebble.WriteOptions{Sync: false})
}

// IterKeys iterates the key list by given prefix, returns the key list.
func (s *idStore) IterKeys(prefix []byte, limit int) (rs [][]byte, err error) {
	it := s.db.NewIter(&pebble.IterOptions{
		LowerBound: prefix,
	})
	defer func() {
		if err0 := it.Close(); err0 != nil {
			s.logger.Warn("close kv iterator resource err",
				logger.String("path", s.path),
				logger.Error(err0))
		}
	}()

	for it.First(); it.Valid(); it.Next() {
		if len(rs) >= limit {
			break
		}
		key := it.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		// copy it
		dst := make([]byte, len(key))
		copy(dst, key)
		rs = append(rs, dst)
	}

	// if iterator has err, returns nil result and err
	if err := it.Error(); err != nil {
		return nil, err
	}
	return rs, nil
}

// Flush flushes the memory table data under pebble db.
func (s *idStore) Flush() error {
	return s.db.Flush()
}

// Close closes backend pebble db.
// NOTICE: need flush first
func (s *idStore) Close() error {
	return s.db.Close()
}
