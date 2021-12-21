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

package indexdb

import (
	"encoding/binary"
	"fmt"
	"io"
	"path"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/unique"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source ./id_mapping_backend.go -destination=./id_mapping_backend_mock.go -package=indexdb

// for testing
var (
	mkDir        = fileutil.MkDirIfNotExist
	newIDStoreFn = unique.NewIDStore
)

const MappingDB = "mapping"

// IDMappingBackend represents the id mapping backend storage,
// save series data(tags hash => series id) under metric
type IDMappingBackend interface {
	io.Closer
	// loadMetricIDMapping loads metric id mapping include id sequence
	loadMetricIDMapping(metricID metric.ID) (idMapping MetricIDMapping, err error)
	// getSeriesID gets series id by metric id/tags hash, if not exist return constants.ErrNotFount
	getSeriesID(metricID metric.ID, tagsHash uint64) (seriesID uint32, err error)
	// genSeries generates series id by metric id/tags hash.
	genSeriesID(metricID metric.ID, tagsHash uint64, seriesID uint32) error
	// sync the backend memory data into persist storage.
	sync() error
}

// idMappingBackend implements IDMappingBackend interface
type idMappingBackend struct {
	db unique.IDStore
}

// newIDMappingBackend creates new id mapping backend storage
func newIDMappingBackend(parent string) (IDMappingBackend, error) {
	if err := mkDir(parent); err != nil {
		return nil, err
	}
	db, err := newIDStoreFn(path.Join(parent, MappingDB))
	if err != nil {
		return nil, err
	}
	return &idMappingBackend{
		db: db,
	}, nil
}

func (imb *idMappingBackend) loadMetricIDMapping(metricID metric.ID) (idMapping MetricIDMapping, err error) {
	mID := metricID.MarshalBinary()
	val, exist, err := imb.db.Get(mID)
	if err != nil {
		return nil, err
	}
	if !exist {
		return newMetricIDMapping(metricID, 0), nil
	}
	sequence := binary.LittleEndian.Uint32(val)
	return newMetricIDMapping(metricID, sequence), nil
}

// getSeriesID gets series id by metric id/tags hash, if not exist return constants.ErrNotFount
func (imb *idMappingBackend) getSeriesID(metricID metric.ID, tagsHash uint64) (seriesID uint32, err error) {
	mID := metricID.MarshalBinary()
	mIDLen := len(mID)
	key := make([]byte, mIDLen+8)
	copy(key, mID)
	binary.LittleEndian.PutUint64(key[mIDLen:], tagsHash)
	val, exist, err := imb.db.Get(key)
	if err != nil {
		return series.EmptySeriesID, err
	}
	if !exist {
		return series.EmptySeriesID, fmt.Errorf("%w, metricID: %d, tagsHash: %d",
			constants.ErrSeriesIDNotFound, metricID, tagsHash)
	}
	seriesID = binary.LittleEndian.Uint32(val)
	return
}

// genSeriesID gets series id by metric id/tags hash, if not exist return constants.ErrNotFount
func (imb *idMappingBackend) genSeriesID(metricID metric.ID, tagsHash uint64, seriesID uint32) error {
	mID := metricID.MarshalBinary()
	mIDLen := len(mID)
	key := make([]byte, mIDLen+8)
	copy(key, mID)
	binary.LittleEndian.PutUint64(key[mIDLen:], tagsHash)

	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], seriesID)
	return imb.db.Put(key, scratch[:])
}

// Close closes the backend storage resource.
func (imb *idMappingBackend) Close() error {
	return imb.db.Close()
}

// sync the backend memory data into persist storage.
func (imb *idMappingBackend) sync() error {
	return imb.db.Flush()
}
