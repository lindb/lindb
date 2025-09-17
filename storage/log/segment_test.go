package log

import (
	"fmt"
	"testing"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/lindb/common/log"
	"github.com/lindb/common/proto/gen/v1/flatLogV1"
	"github.com/stretchr/testify/assert"
)

func TestSegment_Write(t *testing.T) {
	d, err := NewDatabase("test", Options{})
	assert.NoError(t, err)
	shard, _ := NewShard(1, d)
	segment, err := NewSegment(shard)
	assert.NoError(t, err)
	rb := log.CreateRowBuilder()

	rb.AddMessage([]byte("message10")).
		AddTimestamp(123).
		AddField([]byte("key1"), []byte("value1")).
		AddField([]byte("key2"), []byte("value2"))

	data, _ := rb.Build()

	err = segment.Write(data)
	assert.NoError(t, err)

	rb.AddMessage([]byte("message20")).
		AddTimestamp(123).
		AddField([]byte("key1"), []byte("value1")).
		AddField([]byte("key2"), []byte("value2"))

	data, _ = rb.Build()

	err = segment.Write(data)
	assert.NoError(t, err)
	go segment.indexLog()

	time.Sleep(time.Second)
	indexDB := d.IndexDatabase()
	nsID, _ := indexDB.GetNamespaceID([]byte("ns"))
	logIDs := segment.GetLogIDs(nsID)
	fmt.Println(logIDs)
	nsID, _ = indexDB.GetNamespaceID([]byte("ns"))
	logIDs = segment.GetLogIDs(nsID)
	fmt.Println(logIDs)
	keyID, _ := indexDB.GetFieldKeyID(nsID, []byte("key1"))
	logIDs = segment.GetLogIDs(keyID)
	fmt.Println(logIDs)
	valueID, _ := indexDB.GetFieldValueID(keyID, []byte("value1"))
	logIDs = segment.GetLogIDs(valueID)
	fmt.Println(logIDs)
	it := logIDs.Iterator()
	for it.HasNext() {
		data, err := segment.GetLog(uint32(it.Next()))
		assert.NoError(t, err)
		log := &flatLogV1.Log{}
		log.Init(data, flatbuffers.GetUOffsetT(data))
		// logID := s.logID.Inc()
		fmt.Printf("log message=>%s\n", string(log.Message()))
	}

	fmt.Printf("log total:%v\n", logIDs.GetCardinality())

	segment.Close()
}
