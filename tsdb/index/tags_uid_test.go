package index

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/eleme/lindb/kv"
	"github.com/eleme/lindb/pkg/stream"
	"github.com/eleme/lindb/pkg/util"

	"github.com/RoaringBitmap/roaring"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/assert"
)

func TestBitmap(t *testing.T) {
	writer := stream.BinaryWriter()

	bitmap1 := roaring.New()
	bitmap1.Add(1)
	bitmap1.Add(2)
	by, _ := bitmap1.ToBytes()
	writer.PutBytes(by)

	bitmap2 := roaring.New()
	bitmap2.Add(100)
	bitmap2.Add(200)
	by, _ = bitmap2.ToBytes()
	offset := writer.Len()
	writer.PutBytes(by)

	data, _ := writer.Bytes()
	newBitmap1 := roaring.New()
	_, _ = newBitmap1.ReadFrom(bytes.NewBuffer(data[0:]))
	assert.Equal(t, true, bitmap1.Equals(newBitmap1))

	newBitmap2 := roaring.New()
	_, _ = newBitmap2.ReadFrom(bytes.NewBuffer(data[offset:]))
	assert.Equal(t, true, bitmap2.Equals(newBitmap2))
}

func TestTagsUid_GetOrCreateTagsId(t *testing.T) {
	tagsUID := initTags()
	defer util.RemoveDir("../test")

	for i := 1; i < 3; i++ {
		for j := 1; j < count; j++ {
			bitmap := tagsUID.GetTagValueBitmap(uint32(i), "a", fmt.Sprintf("%s%d", "value-1-", j))
			if nil != bitmap {
				assert.Equal(t, uint64(10), bitmap.GetCardinality())
				//it := bitmap.Iterator()
				//for ; it.HasNext(); {
				//	fmt.Print(it.Next(), " ")
				//}
			}
		}
	}
	for i := 1; i < 3; i++ {
		for k := 1; k < count; k++ {
			bitmap := tagsUID.GetTagValueBitmap(uint32(i), "b", fmt.Sprintf("%s%d", "value-2-", k))
			if nil != bitmap {
				assert.Equal(t, uint64(10), bitmap.GetCardinality())
				//it := bitmap.Iterator()
				//for ; it.HasNext(); {
				//	fmt.Print(it.Next(), " ")
				//}
			}
		}
	}
	fmt.Println("success")
}

func TestTagsUid_GetTagNames(t *testing.T) {
	defer util.RemoveDir("../test")
	tagsUID := initTags()
	for i := 1; i < 10; i++ {
		tagsUID.GetTagNames(uint32(i), 100)
		fmt.Println("xxxxx")
	}
}

func TestTagsUid_SuggestTagValues(t *testing.T) {
	defer util.RemoveDir("../test")
	tagsUID := initTags()
	for i := 1; i < 10; i++ {
		tagsUID.SuggestTagValues(uint32(i), "a", "value-1-", 100)
		tagsUID.SuggestTagValues(uint32(i), "a", "v", 100)
		tagsUID.SuggestTagValues(uint32(i), "b", "value-2-", 100)
		tagsUID.SuggestTagValues(uint32(i), "b", "v", 100)

		fmt.Println("xxxxx")
	}
}

var count = 11

func initTags() *TagsUID {
	util.RemoveDir("../test")
	tagsUID := NewTagsUID(initTagsFamily())
	for i := 1; i < 10; i++ {
		for j := 1; j < count; j++ {
			for k := 1; k < count; k++ {
				tagsMap := make(map[string]string)
				tagsMap["a"] = fmt.Sprintf("%s%d", "value-1-", j)
				tagsMap["b"] = fmt.Sprintf("%s%d", "value-2-", k)
				tagsUID.GetOrCreateTagsID(uint32(i), MapToString(tagsMap))
			}
		}
	}
	tagsUID.Flush()
	return tagsUID
}

func initTagsFamily() kv.Family {
	option := kv.DefaultStoreOption("../test")
	var indexStore, _ = kv.NewStore("index", option)
	family, _ := indexStore.CreateFamily("tags", kv.FamilyOption{})
	return family
}

func TestTagsMapping(t *testing.T) {
	id := 1
	//writer := stream.BinaryWriter()
	//util.RemoveDir("../m.txt")
	//f, _ := os.Create("../m.txt")
	buf0 := new(bytes.Buffer)
	//writer := snappy.NewBufferedWriter(bufio.NewWriterSize(f, 100*1024*1024))
	writer := snappy.NewBufferedWriter(buf0)
	length := 0

	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			for k := 0; k < 100; k++ {
				tagsMap := make(map[string]string)

				tagsMap["host"] = fmt.Sprintf("%s%d", "host-", i)
				tagsMap["disk"] = fmt.Sprintf("%s%d", "disk-", j)
				tagsMap["partition"] = fmt.Sprintf("%s%d", "partition-", k)
				//writer.PutUInt32(uint32(id))
				//writer.PutKey([]byte(MapToString(tagsMap)))
				_, err := writer.Write([]byte(MapToString(tagsMap)))
				length += len(MapToString(tagsMap))
				//fmt.Println(n)
				if nil != err {
					fmt.Println(err)
				}
				id++
			}
		}
	}

	err := writer.Flush()
	if nil != err {
		fmt.Println(err)
	}
	err = writer.Close()
	if nil != err {
		fmt.Println(err)
	}

	fmt.Println(buf0.Len())
	fmt.Println(length)

	reader := snappy.NewReader(buf0)
	dst := make([]byte, 59)
	reader.Read(dst)
	fmt.Println(string(dst))
	reader.Read(dst)
	fmt.Println(string(dst))

	//by, _ := writer.Bytes()
	//fmt.Println(len(by))
	//dst := snappy.Encode(nil, by)
	//fmt.Println(len(dst))

}
