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

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/proto/log"
	"github.com/lindb/lindb/storage/log/index"
)

type Database interface {
	Write(namespace []byte, logID uint32, fields *log.FieldIterator) error
	FindLogIDsByField(fieldID uint32) *roaring.Bitmap
	Flush(flusher kv.Flusher) error

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

func (md *database) Write(namespace []byte, logID uint32, fields *log.FieldIterator) error {
	ns, err := md.index.GetOrCreateNamespaceID(namespace)
	if err != nil {
		fmt.Println(err)
	} else {
		md.indexField(ns, logID)
	}

	for fields.HasNext() {
		keyID, _ := md.index.GetOrCreateFieldKeyID(ns, fields.NextName())
		valueID, _ := md.index.GetOrCreateFieldValueID(keyID, fields.NextValue())

		md.indexField(keyID, logID)
		md.indexField(valueID, logID)
	}

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

func (md *database) indexField(fieldID, logID uint32) {
	index, ok := md.fieldIndexes.Get(fieldID)
	if ok {
		index.Add(logID)
	} else {
		md.fieldIndexes.Put(fieldID, roaring.BitmapOf(logID))
	}
}
