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

package index

import (
	"os"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/stream"
)

// for testing
var (
	openFileFn = os.OpenFile
	rwMapFn    = fileutil.RWMap
	syncFn     = fileutil.Sync
	unmapFn    = fileutil.Unmap
)

const (
	SeqSize          = 4 * 4 // ns/name/tag key/tag value
	NamespaceOffset  = 0
	MetricNameOffset = 4
	TagKeyOffset     = 8
	TagValueOffset   = 12
)

// Sequence represents the sequence allocate of metadata.
type Sequence struct {
	buf []byte // mmap buf
	f   *os.File

	ns       *atomic.Uint32
	metric   *atomic.Uint32
	tagKey   *atomic.Uint32
	tagValue *atomic.Uint32
}

// NewSequence creates a Sequence.
func NewSequence(fileName string) (*Sequence, error) {
	f, err := openFileFn(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	buf, err := rwMapFn(f, SeqSize)
	if err != nil {
		// need close file, if map file failure
		_ = f.Close()
		return nil, err
	}
	return &Sequence{
		f:   f,
		buf: buf,
		// init sequence value
		ns:       atomic.NewUint32(stream.ReadUint32(buf, NamespaceOffset)),
		metric:   atomic.NewUint32(stream.ReadUint32(buf, MetricNameOffset)),
		tagKey:   atomic.NewUint32(stream.ReadUint32(buf, TagKeyOffset)),
		tagValue: atomic.NewUint32(stream.ReadUint32(buf, TagValueOffset)),
	}, nil
}

// GetNamespaceSeq retruns the sequence for namespace.
func (s *Sequence) GetNamespaceSeq() uint32 {
	return s.ns.Inc() - 1
}

// GetMetricNameSeq returns the sequence for metric name.
func (s *Sequence) GetMetricNameSeq() uint32 {
	return s.metric.Inc() - 1
}

// GetTagKeySeq returns the sequence for tag key.
func (s *Sequence) GetTagKeySeq() uint32 {
	return s.tagKey.Inc() - 1
}

// GetTagValueSeq returns the sequence for tag value.
func (s *Sequence) GetTagValueSeq() uint32 {
	return s.tagValue.Inc() - 1
}

// Sync persists the sequence data.
func (s *Sequence) Sync() error {
	stream.PutUint32(s.buf, NamespaceOffset, s.ns.Load())
	stream.PutUint32(s.buf, MetricNameOffset, s.metric.Load())
	stream.PutUint32(s.buf, TagKeyOffset, s.tagKey.Load())
	stream.PutUint32(s.buf, TagValueOffset, s.tagValue.Load())
	return syncFn(s.buf)
}

// Close flushs sequence data then closes it.
func (s *Sequence) Close() error {
	defer func() {
		_ = s.f.Close()
	}()
	return unmapFn(s.f, s.buf)
}
