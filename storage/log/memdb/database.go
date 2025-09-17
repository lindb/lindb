package memdb

import (
	"fmt"
	"sync"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/proto/log"
	"github.com/lindb/lindb/storage/log/index"
)

type Database interface {
	Write(namespace []byte, logID uint32, timestamp int64, fields *log.FieldIterator) error
	GetLogIDs(fieldID uint32) *roaring.Bitmap
	Flush(flusher kv.Flusher) error

	IndexDatabase() index.Database
}

type database struct {
	index            index.Database
	timestampIndexes sync.Map                      // timestamp index => (timestamp(second) => bitmap(log ids) )
	fieldIndexes     *imap.IntMap[*roaring.Bitmap] // field index => (global field value id => bitmap(log ids))
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

func (md *database) Write(namespace []byte, logID uint32, timestamp int64, fields *log.FieldIterator) error {
	md.indexTimestamp(timestamp, logID)

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

func (md *database) GetLogIDs(fieldID uint32) *roaring.Bitmap {
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

func (md *database) indexTimestamp(timestamp int64, logID uint32) {
	// TODO: truncate timestamp
	index, ok := md.timestampIndexes.Load(timestamp)
	if ok {
		(index.(*roaring.Bitmap)).Add(logID)
	} else {
		md.timestampIndexes.Store(timestamp, roaring.BitmapOf(logID))
	}
}
