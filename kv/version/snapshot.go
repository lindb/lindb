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

package version

import (
	"errors"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/kv/table"
)

//go:generate mockgen -source ./snapshot.go -destination=./snapshot_mock.go -package version

// Snapshot represents a current family version for reading data.
// NOTICE: current version will retain like ref count, so snapshot must close.
type Snapshot interface {
	// GetCurrent returns current mutable version
	GetCurrent() Version
	// FindReaders finds all files include key
	FindReaders(key uint32) ([]table.Reader, error)
	// Load loads value by key, if exist invoke loader function.
	Load(key uint32, loader func(value []byte) error) error
	// GetReader returns file reader
	GetReader(fileNumber table.FileNumber) (table.Reader, error)
	// Close releases related resources
	Close()
}

// snapshot implements Snapshot interface
type snapshot struct {
	familyName string
	cache      table.Cache

	readers []table.Reader // current read table.Reader list
	version Version
	closed  atomic.Bool
}

// newSnapshot new snapshot instance
func newSnapshot(familyName string, version Version, cache table.Cache) Snapshot {
	version.Retain()
	return &snapshot{
		version:    version,
		familyName: familyName,
		cache:      cache,
	}
}

// GetCurrent returns current mutable version
func (s *snapshot) GetCurrent() Version {
	return s.version
}

// FindReaders finds all files include key
func (s *snapshot) FindReaders(key uint32) ([]table.Reader, error) {
	// find files related given key
	// current version is readonly, if modify version will clone a new one, so needn't lock here.
	files := s.version.FindFiles(key)
	var readers []table.Reader
	for _, fileMeta := range files {
		// get store reader from cache
		reader, err := s.cache.GetReader(s.familyName, Table(fileMeta.GetFileNumber()))
		if err != nil {
			return nil, err
		}
		if reader != nil {
			s.readers = append(s.readers, reader)
			readers = append(readers, reader)
		}
	}
	return readers, nil
}

// Load loads value by key, if exist invoke loader function.
func (s *snapshot) Load(key uint32, loader func(value []byte) error) error {
	files := s.version.FindFiles(key)
	for _, fileMeta := range files {
		reader, err := s.cache.GetReader(s.familyName, Table(fileMeta.GetFileNumber()))
		if err != nil {
			return err
		}
		if reader != nil {
			value, err := reader.Get(key)
			if errors.Is(err, table.ErrKeyNotExist) {
				// not exist
				continue
			}
			if err != nil {
				return err
			}
			if err := loader(value); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetReader returns the file reader
func (s *snapshot) GetReader(fileNumber table.FileNumber) (table.Reader, error) {
	reader, err := s.cache.GetReader(s.familyName, Table(fileNumber))
	if reader != nil {
		s.readers = append(s.readers, reader)
	}
	return reader, err
}

// Close releases related resources
func (s *snapshot) Close() {
	// atomic set closed status, make sure only release once
	if s.closed.CompareAndSwap(false, true) {
		s.version.Release()
		s.cache.ReleaseReaders(s.readers)
	}
}
