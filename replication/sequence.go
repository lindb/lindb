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

package replication

import (
	"io"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

//go:generate mockgen -source=./sequence.go -destination=./sequence_mock.go -package=replication

// for testing
var (
	newPageFactoryFunc = page.NewFactory
)

var sequenceLogger = logger.GetLogger("replication", "Sequence")

const (
	//sequenceMetaSize 8 bytes for int64
	sequenceMetaSize = 8
	metaPageID       = 0
)

// Sequence represents a persistence sequence recorder
// for on storage side when transferring data from broker to storage.
type Sequence interface {
	io.Closer
	// GetHeadSeq returns the head sequence which is the latest sequence of replica received.
	GetHeadSeq() int64
	// SetHeadSeq sets the head sequence which is the latest sequence of replica received.
	SetHeadSeq(seq int64)
	// GetAckSeq returns the ack sequence which is the latest sequence of replica successfully flushed to disk.
	GetAckSeq() int64
	// SetAckSeq sets the ack sequence which is the latest sequence of replica successfully flushed to disk.
	SetAckSeq(seq int64)
	// Sync syncs the Sequence to storage.
	Sync() error

	//TODO need add close method??
}

// sequence implements Sequence.
type sequence struct {
	dirPath     string
	metaPageFct page.Factory
	// meta stores the ackSeq to page cache.
	metaPage page.MappedPage
	// headSeq represents the the max sequence num of replica received.
	headSeq atomic.Int64
	// ackSeq represents the the max sequence num of replica flushed to disk.
	ackSeq atomic.Int64
}

// NewSequence returns a sequence with page cache corresponding to dirPath.
func NewSequence(dirPath string) (Sequence, error) {
	var err error
	metaPageFct, err := newPageFactoryFunc(dirPath, sequenceMetaSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if err1 := metaPageFct.Close(); err1 != nil {
				sequenceLogger.Error("close meta page factory err",
					logger.String("path", dirPath), logger.Error(err1))
			}
		}
	}()

	_, ok := metaPageFct.GetPage(metaPageID)

	metaPage, err := metaPageFct.AcquirePage(metaPageID)
	if err != nil {
		return nil, err
	}
	ackSeq := int64(-1) // for new sequence
	if ok {
		// if exist meta page, need read
		ackSeq = int64(metaPage.ReadUint64(0))
	}

	s := &sequence{
		dirPath:     dirPath,
		metaPageFct: metaPageFct,
		metaPage:    metaPage,
		headSeq:     *atomic.NewInt64(ackSeq),
		ackSeq:      *atomic.NewInt64(ackSeq),
	}
	if err = s.Sync(); err != nil {
		return nil, err
	}
	return s, nil
}

// GetHeadSeq returns the head sequence which is the latest sequence of replica received.
func (s *sequence) GetHeadSeq() int64 {
	return s.headSeq.Load()
}

// SetHeadSeq sets the head sequence which is the latest sequence of replica received.
func (s *sequence) SetHeadSeq(seq int64) {
	s.headSeq.Store(seq)
}

// GetAckSeq returns the ack sequence which is the latest sequence of replica successfully flushed to disk.
func (s *sequence) GetAckSeq() int64 {
	return s.ackSeq.Load()
}

// SetAckSeq sets the ack sequence which is the latest sequence of replica successfully flushed to disk.
func (s *sequence) SetAckSeq(seq int64) {
	s.ackSeq.Store(seq)
}

// Sync syncs the Sequence to storage.
func (s *sequence) Sync() error {
	s.metaPage.PutUint64(uint64(s.GetAckSeq()), 0)
	return s.metaPage.Sync()
}

// Close closes the page factory
func (s *sequence) Close() error {
	return s.metaPageFct.Close()
}
