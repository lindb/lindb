package kv

import (
	"fmt"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package kv

// Flusher flushes data into kv store, for big data will be split into many sstable
type Flusher interface {
	// Add puts k/v pair
	Add(key uint32, value []byte) error
	// Commit flushes data and commits metadata
	Commit() error
}

// storeFlusher is family level store flusher
type storeFlusher struct {
	family  *family
	builder table.Builder
	editLog *version.EditLog
}

// newStoreFlusher create family store flusher
func newStoreFlusher(family *family) Flusher {
	return &storeFlusher{
		family:  family,
		editLog: version.NewEditLog(family.option.ID),
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
		sf.builder = builder
	}
	//TODO add file size limit
	return sf.builder.Add(key, value)
}

// Commit flushes data and commits metadata
func (sf *storeFlusher) Commit() error {
	builder := sf.builder
	if builder != nil {
		if err := builder.Close(); err != nil {
			return fmt.Errorf("close table builder error when flush commit, error:%s", err)
		}

		fileMeta := version.NewFileMeta(builder.FileNumber(), builder.MinKey(), builder.MaxKey(), builder.Size())
		sf.editLog.Add(version.CreateNewFile(0, fileMeta))
	}

	if flag := sf.family.commitEditLog(sf.editLog); !flag {
		return fmt.Errorf("commit edit log failure")
	}
	return nil
}
