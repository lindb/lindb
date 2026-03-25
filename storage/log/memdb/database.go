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

package memdb

import (
	"fmt"

	"github.com/lindb/arrow/pkg/constants"
	"github.com/lindb/arrow/pkg/logs"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/log/index"
)

type Database interface {
	Write(namespace []byte, logID uint32, reader *logs.Reader, row int) error
	FindLogIDsByField(fieldID uint32) *roaring.Bitmap
	Flush(flusher kv.Flusher) error
	// MemSize returns an estimated memory footprint of the field index bitmaps in bytes.
	MemSize() int64

	IndexDatabase() index.Database
}

type database struct {
	index        index.Database
	fieldIndexes *imap.IntMap[*roaring.Bitmap] // field index => (global field value id => bitmap(log ids))
}

func NewDatabase(indexDB index.Database) Database {
	return &database{
		index:        indexDB,
		fieldIndexes: imap.NewIntMap[*roaring.Bitmap](),
	}
}

func (md *database) IndexDatabase() index.Database {
	return md.index
}

func (md *database) Write(namespace []byte, logID uint32, reader *logs.Reader, row int) error {
	ns, err := md.index.GetOrCreateNamespaceID(namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		md.indexField(ns, logID)
	}

	traceID := reader.TraceID(row)
	if len(traceID) > 0 {
		traceIDKeyID, _ := md.index.GetOrCreateFieldKeyID(ns, strutil.String2ByteSlice(constants.TraceID))
		traceIDValueID, _ := md.index.GetOrCreateFieldValueID(traceIDKeyID, traceID)

		md.indexField(traceIDKeyID, logID)
		md.indexField(traceIDValueID, logID)
	}

	spanID := reader.SpanID(row)
	if len(spanID) > 0 {
		spanIDKeyID, _ := md.index.GetOrCreateFieldKeyID(ns, strutil.String2ByteSlice(constants.SpanID))
		spanIDValueID, _ := md.index.GetOrCreateFieldValueID(spanIDKeyID, spanID)

		md.indexField(spanIDKeyID, logID)
		md.indexField(spanIDValueID, logID)
	}

	level := reader.Level(row)
	if level != "" {
		levelKeyID, _ := md.index.GetOrCreateFieldKeyID(ns, strutil.String2ByteSlice(constants.Level))
		levelValueID, _ := md.index.GetOrCreateFieldValueID(levelKeyID, strutil.String2ByteSlice(level))

		md.indexField(levelKeyID, logID)
		md.indexField(levelValueID, logID)
	}
	eventName := reader.EventName(row)
	if eventName != "" {
		eventNameKeyID, _ := md.index.GetOrCreateFieldKeyID(ns, strutil.String2ByteSlice(constants.EventName))
		eventNameValueID, _ := md.index.GetOrCreateFieldValueID(eventNameKeyID, strutil.String2ByteSlice(eventName))

		md.indexField(eventNameKeyID, logID)
		md.indexField(eventNameValueID, logID)
	}

	reader.Attributes(row, func(key, value string) {
		keyID, _ := md.index.GetOrCreateFieldKeyID(ns, strutil.String2ByteSlice(key))
		valueID, _ := md.index.GetOrCreateFieldValueID(keyID, strutil.String2ByteSlice(value))

		md.indexField(keyID, logID)
		md.indexField(valueID, logID)
	})

	return nil
}

func (md *database) FindLogIDsByTimeRange(timeRange timeutil.SlotRange) *roaring.Bitmap {
	return nil
}

func (md *database) FindLogIDsByField(fieldID uint32) *roaring.Bitmap {
	value, ok := md.fieldIndexes.Get(fieldID)
	if ok {
		return value
	}
	return nil
}

func (md *database) Flush(flusher kv.Flusher) error {
	return md.fieldIndexes.WalkEntry(func(key uint32, value *roaring.Bitmap) error {
		data, err := value.ToBytes()
		if err != nil {
			return err
		}
		return flusher.Add(key, data)
	})
}

// MemSize returns an estimated memory footprint of all field index bitmaps in bytes.
// It sums GetSizeInBytes() for each roaring bitmap, which reflects the compressed size.
// The actual heap usage may be slightly higher due to Go object overhead.
func (md *database) MemSize() int64 {
	var size int64
	md.fieldIndexes.WalkEntry(func(_ uint32, v *roaring.Bitmap) error { //nolint:errcheck
		size += int64(v.GetSizeInBytes())
		return nil
	})
	return size
}

func (md *database) indexField(fieldID, logID uint32) {
	index, ok := md.fieldIndexes.Get(fieldID)
	if ok {
		index.Add(logID)
	} else {
		md.fieldIndexes.Put(fieldID, roaring.BitmapOf(logID))
	}
}
