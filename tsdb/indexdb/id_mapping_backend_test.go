package indexdb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
)

const testPath = "test"

func TestIdMappingBackend_new(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		seriesBucketName = []byte("s")
		closeFunc = closeDB
		mkDir = fileutil.MkDirIfNotExist
	}()
	// case 1: new backend
	backend, err := newIDMappingBackend("test", testPath)
	assert.NoError(t, err)
	assert.NotNil(t, backend)
	// case 2: cannot reopen
	backend2, err := newIDMappingBackend("test", testPath)
	assert.Error(t, err)
	assert.Nil(t, backend2)
	err = backend.Close()
	assert.NoError(t, err)

	// case 3: mock create root bucket
	seriesBucketName = []byte("")
	backend2, err = newIDMappingBackend("test", testPath)
	assert.Error(t, err)
	assert.Nil(t, backend2)
	closeFunc = func(db *bbolt.DB) error {
		return fmt.Errorf("err")
	}
	seriesBucketName = []byte("")
	backend2, err = newIDMappingBackend("test", testPath)
	assert.Error(t, err)
	assert.Nil(t, backend2)
	// case 4: create parent err
	mkDir = func(path string) error {
		return fmt.Errorf("err")
	}
	backend, err = newIDMappingBackend("test", testPath)
	assert.Error(t, err)
	assert.Nil(t, backend)
}

func TestIdMappingBackend_mapping(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	backend, err := newIDMappingBackend("test", filepath.Join(testPath, "test"))
	assert.NoError(t, err)
	event := newMappingEvent()
	event.addSeriesID(1, 20, 200)
	event.addSeriesID(2, 10, 100)
	event.addSeriesID(2, 30, 300)
	err = backend.saveMapping(event)
	assert.NoError(t, err)

	// case 1: get series
	seriesID, err := backend.getSeriesID(2, 30)
	assert.NoError(t, err)
	assert.Equal(t, uint32(300), seriesID)
	// case 2: metric id not exist
	seriesID, err = backend.getSeriesID(4, 30)
	assert.Equal(t, err, constants.ErrNotFound)
	assert.Equal(t, uint32(0), seriesID)
	// case 3: series id not exist
	seriesID, err = backend.getSeriesID(2, 300)
	assert.Equal(t, err, constants.ErrNotFound)
	assert.Equal(t, uint32(0), seriesID)
	// case 4: load mapping not exist
	mapping, err := backend.loadMetricIDMapping(30)
	assert.Equal(t, err, constants.ErrNotFound)
	assert.Nil(t, mapping)
	// case 5: load mapping not exist
	mapping, err = backend.loadMetricIDMapping(2)
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), mapping.GetMetricID())
	mapping1 := mapping.(*metricIDMapping)
	assert.Equal(t, uint32(300), mapping1.idSequence.Load())

	err = backend.Close()
	assert.NoError(t, err)

	//reopen
	backend, _ = newIDMappingBackend("test", filepath.Join(testPath, "test"))
	mapping, err = backend.loadMetricIDMapping(2)
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), mapping.GetMetricID())
	mapping = mapping.(*metricIDMapping)
	assert.Equal(t, uint32(300), mapping1.idSequence.Load())
}

func TestIdMappingBackend_save_err(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
		setSequenceFunc = setSequence
		createBucketFunc = createBucket
		putFunc = put
	}()
	backend, err := newIDMappingBackend("test", testPath)
	assert.NoError(t, err)
	event := newMappingEvent()
	event.addSeriesID(1, 20, 200)
	event.addSeriesID(2, 10, 100)
	event.addSeriesID(2, 30, 300)
	setSequenceFunc = func(bucket *bbolt.Bucket, seq uint64) error {
		return fmt.Errorf("err")
	}
	err = backend.saveMapping(event)
	assert.Error(t, err)

	setSequenceFunc = setSequence
	createBucketFunc = func(parentBucket *bbolt.Bucket, name []byte) (bucket *bbolt.Bucket, err error) {
		return nil, fmt.Errorf("err")
	}
	err = backend.saveMapping(event)
	assert.Error(t, err)
	createBucketFunc = createBucket
	putFunc = func(bucket *bbolt.Bucket, key, value []byte) error {
		return fmt.Errorf("err")
	}
	err = backend.saveMapping(event)
	assert.Error(t, err)
}
