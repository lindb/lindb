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
	"time"

	"go.etcd.io/bbolt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source ./id_mapping_backend.go -destination=./id_mapping_backend_mock.go -package=indexdb

const MappingDB = "mapping.db"

// for testing
var (
	mkDir            = fileutil.MkDirIfNotExist
	closeFunc        = closeDB
	setSequenceFunc  = setSequence
	createBucketFunc = createBucket
	putFunc          = put
)

var (
	seriesBucketName = []byte("s")
)

// IDMappingBackend represents the id mapping backend storage,
// save series data(tags hash => series id) under metric
type IDMappingBackend interface {
	io.Closer
	// loadMetricIDMapping loads metric id mapping include id sequence
	loadMetricIDMapping(metricID uint32) (idMapping MetricIDMapping, err error)
	// getSeriesID gets series id by metric id/tags hash, if not exist return constants.ErrNotFount
	getSeriesID(metricID uint32, tagsHash uint64) (seriesID uint32, err error)
	// saveMapping saves the id mapping event
	saveMapping(event *mappingEvent) (err error)
}

// idMappingBackend implements IDMappingBackend interface
type idMappingBackend struct {
	db *bbolt.DB
}

// newIDMappingBackend creates new id mapping backend storage
func newIDMappingBackend(parent string) (IDMappingBackend, error) {
	if err := mkDir(parent); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path.Join(parent, MappingDB), 0600, &bbolt.Options{Timeout: 1 * time.Second, NoSync: true})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		// create series root bucket for save metric's id mapping
		_, err := tx.CreateBucketIfNotExists(seriesBucketName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		// close bbolt.DB if init mapping backend err
		if e := closeFunc(db); e != nil {
			indexLogger.Error("close bbolt.db err when create mapping backend fail",
				logger.String("db", parent), logger.Error(e))
		}
		return nil, err
	}
	return &idMappingBackend{
		db: db,
	}, nil
}

// loadMetricIDMapping loads metric id mapping include id sequence
func (imb *idMappingBackend) loadMetricIDMapping(metricID uint32) (idMapping MetricIDMapping, err error) {
	var sequence uint32
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = imb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(seriesBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return fmt.Errorf("%w, metricID: %d", constants.ErrMetricBucketNotFound, metricID)
		}
		sequence = uint32(metricBucket.Sequence())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w, metricID: %d, loadMetricIDMapping with error: %s",
			constants.ErrMetricBucketNotFound, metricID, err)
	}
	return newMetricIDMapping(metricID, sequence), nil
}

// getSeriesID gets series id by metric id/tags hash, if not exist return constants.ErrNotFount
func (imb *idMappingBackend) getSeriesID(metricID uint32, tagsHash uint64) (seriesID uint32, err error) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], metricID)
	err = imb.db.View(func(tx *bbolt.Tx) error {
		metricBucket := tx.Bucket(seriesBucketName).Bucket(scratch[:])
		if metricBucket == nil {
			return fmt.Errorf("%w, metricID: %d, tagsHash: %d",
				constants.ErrMetricBucketNotFound, metricID, tagsHash)
		}
		var hash [8]byte
		binary.LittleEndian.PutUint64(hash[:], tagsHash)
		value := metricBucket.Get(hash[:])
		if len(value) == 0 {
			return fmt.Errorf("%w, metricID: %d, tagsHash: %d",
				constants.ErrSeriesIDNotFound, metricID, tagsHash)
		}
		seriesID = binary.LittleEndian.Uint32(value)
		return nil
	})
	return
}

// saveMapping saves the id mapping event
func (imb *idMappingBackend) saveMapping(event *mappingEvent) (err error) {
	err = imb.db.Update(func(tx *bbolt.Tx) error {
		for metricID, metricEvent := range event.events {
			var scratch [4]byte
			binary.LittleEndian.PutUint32(scratch[:], metricID)
			id := scratch[:]
			root := tx.Bucket(seriesBucketName)
			metricBucket := root.Bucket(id)
			if metricBucket == nil {
				// create metric bucket if metric id not exist
				metricBucket, err = createBucketFunc(root, id)
				if err != nil {
					return err
				}
			}
			// save series data
			for _, seriesEvent := range metricEvent.events {
				var seriesID [4]byte
				var hash [8]byte
				binary.LittleEndian.PutUint64(hash[:], seriesEvent.tagsHash)
				binary.LittleEndian.PutUint32(seriesID[:], seriesEvent.seriesID)
				if err = putFunc(metricBucket, hash[:], seriesID[:]); err != nil {
					return err
				}
			}
			// save metric id sequence
			if err = setSequenceFunc(metricBucket, uint64(metricEvent.metricIDSeq)); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Close closes the bbolt.DB
func (imb *idMappingBackend) Close() error {
	return imb.db.Close()
}

// closeDB closes the bbolt.DB
func closeDB(db *bbolt.DB) error {
	return db.Close()
}

// setSequence sets the bucket's sequence
func setSequence(bucket *bbolt.Bucket, seq uint64) error {
	if bucket.Sequence() < seq {
		return bucket.SetSequence(seq)
	}
	return nil
}

// createBucket creates the bucket with name
func createBucket(parentBucket *bbolt.Bucket, name []byte) (*bbolt.Bucket, error) {
	return parentBucket.CreateBucket(name)
}

// put puts the key/value
func put(bucket *bbolt.Bucket, key, value []byte) error {
	return bucket.Put(key, value)
}
