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

package kv

import (
	"bytes"
	"fmt"
	"hash"
	"hash/crc32"
	"time"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/metrics"
)

//go:generate mockgen -source ./flusher.go -destination=./flusher_mock.go -package kv

// Flusher flushes data into kv store, for big data will be split into many sst files.
type Flusher interface {
	// StreamWriter creates a stream writer for flushing in stream.
	StreamWriter() (table.StreamWriter, error)
	// Add puts k/v pair
	Add(key uint32, value []byte) error
	// Sequence sets write sequence number.
	Sequence(leader int32, seq int64)
	// Commit flushes data and commits metadata.
	Commit() error
	// Release releases the resource of flusher.
	// NOTE: MUST invoke Release() after new fluster instance.
	Release()
}

// storeFlusher is family level store flusher.
type storeFlusher struct {
	start     time.Time
	family    Family
	builder   table.Builder
	editLog   version.EditLog
	sequences map[int32]int64
	releaseFn func()
	outputs   []table.FileNumber
}

// newStoreFlusher create family store flusher.
func newStoreFlusher(family Family, releaseFn func()) Flusher {
	metrics.FlushStatistics.Flushing.Incr()
	return &storeFlusher{
		family:    family,
		editLog:   version.NewEditLog(family.ID()),
		sequences: make(map[int32]int64),
		releaseFn: releaseFn,
		start:     time.Now(),
	}
}

func (sf *storeFlusher) checkBuilder() error {
	if sf.builder == nil {
		builder, err := sf.family.newTableBuilder()
		if err != nil {
			return fmt.Errorf("create table build error:%s", err)
		}
		sf.outputs = append(sf.outputs, builder.FileNumber())
		sf.builder = builder
	}
	return nil
}

// Add adds puts k/v pair.
// NOTICE: key must key in sort by desc
func (sf *storeFlusher) Add(key uint32, value []byte) error {
	if err := sf.checkBuilder(); err != nil {
		metrics.FlushStatistics.Failure.Incr()
		return err
	}
	// TODO add file size limit
	return sf.builder.Add(key, value)
}

// Sequence sets write sequence number.
func (sf *storeFlusher) Sequence(leader int32, seq int64) {
	sf.sequences[leader] = seq
}

func (sf *storeFlusher) StreamWriter() (table.StreamWriter, error) {
	if err := sf.checkBuilder(); err != nil {
		metrics.FlushStatistics.Failure.Incr()
		return nil, err
	}
	return sf.builder.StreamWriter(), nil
}

// Commit flushes data and commits metadata.
func (sf *storeFlusher) Commit() (err error) {
	builder := sf.builder
	defer func() {
		if err != nil {
			metrics.FlushStatistics.Failure.Incr()
		}
		if builder != nil {
			// remove temp file number if fail
			fileNumber := builder.FileNumber()
			sf.family.removePendingOutput(fileNumber)
		}
	}()
	if builder != nil {
		if builder.Size() > 0 {
			err = builder.Close()
			if err != nil {
				return fmt.Errorf("close table builder error when flush commit, error:%s", err)
			}

			fileMeta := version.NewFileMeta(builder.FileNumber(), builder.MinKey(), builder.MaxKey(), builder.Size())
			sf.editLog.Add(version.CreateNewFile(0, fileMeta))
		} else {
			// if no keys, abandon kv builder
			err = builder.Abandon()
			if err != nil {
				return fmt.Errorf("abandon table builder error when flush commit, error:%s", err)
			}
		}
	}
	for leader, seq := range sf.sequences {
		// add sequence for each leader
		sf.editLog.Add(version.CreateSequence(leader, seq))
	}

	// check if it needs add rollup log to target store
	if len(sf.outputs) > 0 {
		store := sf.family.getStore()
		rollupTargetStores := store.Option().Rollup
		for _, interval := range rollupTargetStores {
			// add rollup files edit log in source version
			for _, output := range sf.outputs {
				sf.editLog.Add(version.CreateNewRollupFile(output, interval))
			}
		}
	}

	// TODO: check if edit log is empty?
	if flag := sf.family.commitEditLog(sf.editLog); !flag {
		err = fmt.Errorf("commit edit log failure")
		return err
	}
	return nil
}

// Release releases the resource of flusher.
func (sf *storeFlusher) Release() {
	metrics.FlushStatistics.Flushing.Decr()
	metrics.FlushStatistics.Duration.UpdateSince(sf.start)
	sf.releaseFn()
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

func (nf *NopFlusher) StreamWriter() (table.StreamWriter, error) {
	return &nopStreamWriter{
		buffer: &nf.buffer,
		crc32:  crc32.New(crc32.IEEETable),
	}, nil
}

// Add puts value to the buffer.
func (nf *NopFlusher) Add(_ uint32, value []byte) error {
	nf.buffer.Reset()
	nf.buffer.Write(value)
	return nil
}

func (nf *NopFlusher) Sequence(_ int32, _ int64) {}

// Commit always return nil
func (nf *NopFlusher) Commit() error {
	nf.buffer.Reset()
	return nil
}

func (nf *NopFlusher) Release() {}

type nopStreamWriter struct {
	crc32  hash.Hash32
	buffer *bytes.Buffer
	size   uint32
}

func (nw *nopStreamWriter) Prepare(_ uint32) {
	nw.size = 0
	nw.buffer.Reset()
	nw.crc32.Reset()
}

func (nw *nopStreamWriter) Write(data []byte) (int, error) {
	_, _ = nw.buffer.Write(data)
	nw.size += uint32(len(data))
	_, _ = nw.crc32.Write(data)
	return len(data), nil
}

func (nw *nopStreamWriter) CRC32CheckSum() uint32 { return nw.crc32.Sum32() }

func (nw *nopStreamWriter) Size() uint32 { return nw.size }

func (nw *nopStreamWriter) Commit() error { return nil }

func (nw *nopStreamWriter) Release() {}
