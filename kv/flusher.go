package kv

import (
	"bytes"
	"fmt"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package kv

// Flusher flushes data into kv store, for big data will be split into many sst files
type Flusher interface {
	// Add puts k/v pair
	Add(key uint32, value []byte) error
	// Commit flushes data and commits metadata
	Commit() error
}

// storeFlusher is family level store flusher
type storeFlusher struct {
	family  Family
	builder table.Builder
	editLog version.EditLog
}

// newStoreFlusher create family store flusher
func newStoreFlusher(family Family) Flusher {
	return &storeFlusher{
		family:  family,
		editLog: version.NewEditLog(family.ID()),
	}
}

// Add adds puts k/v pair.
// NOTICE: key must key in sort by desc
func (sf *storeFlusher) Add(key uint32, value []byte) error {
	if sf.builder == nil {
		builder, err := sf.family.newTableBuilder()
		if err != nil {
			return fmt.Errorf("create table build error:%s", err)
		}
		sf.family.addPendingOutput(builder.FileNumber())
		sf.builder = builder
	}
	//TODO add file size limit
	return sf.builder.Add(key, value)
}

// Commit flushes data and commits metadata
func (sf *storeFlusher) Commit() (err error) {
	builder := sf.builder
	defer func() {
		if builder != nil {
			// remove temp file number if fail
			fileNumber := builder.FileNumber()
			sf.family.removePendingOutput(fileNumber)
		}
	}()
	if builder != nil {
		if err := builder.Close(); err != nil {
			err = fmt.Errorf("close table builder error when flush commit, error:%s", err)
			return err
		}

		fileMeta := version.NewFileMeta(builder.FileNumber(), builder.MinKey(), builder.MaxKey(), builder.Size())
		sf.editLog.Add(version.CreateNewFile(0, fileMeta))
	}

	if flag := sf.family.commitEditLog(sf.editLog); !flag {
		err = fmt.Errorf("commit edit log failure")
		return err
	}
	return nil
}

// NopFlusher implements Flusher, but does nothing.
type NopFlusher struct {
	buffer bytes.Buffer
}

// NewNopFlusher returns a new no-operation-flusher
func NewNopFlusher() *NopFlusher {
	return &NopFlusher{}
}

// Bytes returns last-flushed-value
func (nf *NopFlusher) Bytes() []byte { return nf.buffer.Bytes() }

// Add puts value to the buffer.
func (nf *NopFlusher) Add(key uint32, value []byte) error {
	nf.buffer.Reset()
	nf.buffer.Write(value)
	return nil
}

// Commit always return nil
func (nf *NopFlusher) Commit() error { return nil }
