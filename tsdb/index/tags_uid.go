package index

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/tree"

	"github.com/RoaringBitmap/roaring"
	"go.uber.org/zap"
)

//TagsUID represents tags unique id under the metric name.
type TagsUID struct {
	metricID  uint32
	tagsIDMap map[uint32]uint32

	tagsMap   map[string]*tree.BTree
	bitmaps   []*roaring.Bitmap
	bitmapSeq uint32 //bitmap sequence Id

	family  kv.Family
	dbField zap.Field
}

//TagsReader represents parses tags byte arrays for reading
type TagsReader struct {
	reader               *stream.ByteBufReader //byte buf reader
	tagNameOffset        map[string]int        //tagName tree start offset
	tagTreePosition      int
	bitmapOffsetPosition int
	bitmapPosition       int
}

//NewTagsUID creation requires kvFamily
func NewTagsUID(f kv.Family) *TagsUID {
	return &TagsUID{
		bitmaps:   make([]*roaring.Bitmap, 8),
		tagsMap:   make(map[string]*tree.BTree),
		tagsIDMap: make(map[uint32]uint32),
		family:    f,
	}
}

//newTagsReader returns TagsReader
func newTagsReader(byteArray []byte) *TagsReader {
	bufReader := stream.NewBufReader(byteArray)

	//header
	tagNames := int(bufReader.ReadUvarint64())
	tagOffset := make(map[string]int, tagNames)
	for i := 0; i < tagNames; i++ {
		_, tagName := bufReader.ReadLenBytes()
		offset := int(bufReader.ReadUvarint64())
		tagOffset[string(tagName)] = offset
	}
	//read tag tree position
	tagTreeLen := int(bufReader.ReadUvarint64())
	tagTreePos := bufReader.GetPosition()
	bufReader.NewPosition(bufReader.GetPosition() + tagTreeLen)
	//read bitmap offset position
	bitmapOffsetLen := int(bufReader.ReadUvarint64())
	bitmapOffsetPos := bufReader.GetPosition()
	bufReader.NewPosition(bufReader.GetPosition() + bitmapOffsetLen)

	//read bitmap and tag tree position
	bitmapPos := bufReader.GetPosition()

	return &TagsReader{
		reader:               bufReader,
		tagNameOffset:        tagOffset,
		tagTreePosition:      tagTreePos,
		bitmapOffsetPosition: bitmapOffsetPos,
		bitmapPosition:       bitmapPos,
	}
}

//getTagValueBitmap returns tag value bitmap from disk
func (tr *TagsReader) getTagValueBitmap(tagName, tagValue string) *roaring.Bitmap {
	offset, ok := tr.tagNameOffset[tagName]
	if ok {
		tr.reader.NewPosition(tr.tagTreePosition + offset)
		treeLen := int(tr.reader.ReadUvarint64())
		treeBytes := tr.reader.ReadBytes(treeLen)
		treeReader := tree.NewReader(treeBytes)
		bitmapIdx, ok := treeReader.Get([]byte(tagValue))
		if ok {
			tr.reader.NewPosition(tr.bitmapOffsetPosition + bitmapIdx*4)
			bitmapPos := int(tr.reader.ReadUint32())
			tr.reader.NewPosition(tr.bitmapPosition + bitmapPos)
			pos := tr.reader.GetPosition()

			bitmap := roaring.New()
			_, err := bitmap.ReadFrom(bytes.NewBuffer(tr.reader.SubArray(pos)))
			if nil != err {
				logger.GetLogger("tsdb/index").Error("decode bitmap error:", zap.String(tagName, tagValue), logger.Error(err))
				return nil
			}
			return bitmap
		}
	}
	return nil
}

//seek returns prefix tag value Iterator
func (tr *TagsReader) seek(tagName string, prefix []byte) tree.Iterator {
	offset, ok := tr.tagNameOffset[tagName]
	if ok {
		tr.reader.NewPosition(tr.tagTreePosition + offset)
		treeLen := int(tr.reader.ReadUvarint64())
		treeBytes := tr.reader.ReadBytes(treeLen)
		treeReader := tree.NewReader(treeBytes)
		return treeReader.Seek(prefix)
	}
	return nil
}

//getSortTagNames returns get the sorted tag names
func (t *TagsUID) getSortTagNames() []string {
	tagNames := make([]string, len(t.tagsMap))
	var idx int
	for tagName := range t.tagsMap {
		tagNames[idx] = tagName
		idx++
	}
	sort.Strings(tagNames)
	return tagNames
}

//GetOrCreateTagsID returns find the tags ID associated with given tags or create it.
func (t *TagsUID) GetOrCreateTagsID(metricID uint32, tags string) (uint32, error) {
	tagsMap := StringToMap(tags)
	//load from kv-store
	tagsID := t.getTagsIDFromDisk(metricID, tagsMap)
	if tagsID == NotFoundTagsID {
		if t.metricID != metricID {
			err := t.Flush()
			if nil != err {
				logger.GetLogger("tsdb/index").Error("flush metric field error!", t.dbField, logger.Error(err))
				return NotFoundTagsID, err
			}
			t.metricID = metricID
			t.clear()
		}

		tagsID = t.tagsIDMap[metricID] + 1
		//tagsMap -> tagsID
		for tagName, tagValue := range tagsMap {
			tagTree, ok := t.tagsMap[tagName]
			if !ok {
				tagTree = tree.NewBTree()
				t.tagsMap[tagName] = tagTree
			}
			var bitmapIdx uint32
			v, ok := tagTree.Get([]byte(tagValue))
			if !ok {
				bitmapIdx = t.bitmapSeq
				if bitmapIdx >= uint32(len(t.bitmaps)) {
					target := make([]*roaring.Bitmap, len(t.bitmaps)*2)
					copy(target, t.bitmaps)
					t.bitmaps = target
				}
				t.bitmaps[bitmapIdx] = roaring.New()
				tagTree.Put([]byte(tagValue), int(bitmapIdx))
				t.bitmapSeq++
			} else {
				bitmapIdx = uint32(v)
			}
			bitmap := t.bitmaps[bitmapIdx]
			bitmap.Add(tagsID)
		}
		t.tagsIDMap[metricID]++
		return tagsID, nil
	}
	return tagsID, nil
}

func (t *TagsUID) clear() {
	//clear all
	for k := range t.tagsMap {
		delete(t.tagsMap, k)
	}
	t.bitmapSeq = 0
	t.bitmaps = make([]*roaring.Bitmap, 8)
}

//getTagsIDFromDisk returns find the tags ID associated with given tags
func (t *TagsUID) getTagsIDFromDisk(metricID uint32, tags map[string]string) uint32 {
	var result *roaring.Bitmap
	t.family.Lookup(metricID, func(byteArray []byte) bool {
		tagsReader := newTagsReader(byteArray)
		for tagName, tagValue := range tags {
			bitmap := tagsReader.getTagValueBitmap(tagName, tagValue)
			if nil == bitmap {
				return true
			}
			if nil == result {
				result = bitmap
			} else {
				result.And(bitmap)
			}
		}
		if nil != result && !result.IsEmpty() {
			return true
		}
		return false
	})
	if nil != result && !result.IsEmpty() {
		//first tags ID
		return result.Iterator().Next()
	}
	return NotFoundTagsID
}

//GetTagValueBitmap returns find bitmap associated with a given tag value
func (t *TagsUID) GetTagValueBitmap(metricID uint32, tagName string, tagValue string) *roaring.Bitmap {
	var result *roaring.Bitmap
	t.family.Lookup(metricID, func(byteArray []byte) bool {
		tagsReader := newTagsReader(byteArray)
		bitmap := tagsReader.getTagValueBitmap(tagName, tagValue)
		if nil != bitmap {
			result = bitmap
			return true
		}
		return false
	})
	return result
}

//GetTagNames return get all tag names within the metric name
func (t *TagsUID) GetTagNames(metricID uint32, limit int16) map[string]struct{} {
	//todo
	t.family.Lookup(metricID, func(byteArray []byte) bool {
		tagsReader := newTagsReader(byteArray)
		for tagName := range tagsReader.tagNameOffset {
			fmt.Println("xx====", tagName)
		}
		return false
	})
	return nil
}

//SuggestTagValues returns suggestions of tag values given a search prefix
func (t *TagsUID) SuggestTagValues(metricID uint32, tagName string, tagValuePrefix string,
	limit uint16) map[string]struct{} {
	//todo
	t.family.Lookup(metricID, func(byteArray []byte) bool {
		tagsReader := newTagsReader(byteArray)
		_ = tagsReader.seek(tagName, []byte(tagValuePrefix))
		return false
	})
	return nil
}

//Flush represents forces a flush of in-memory data, and clear it
func (t *TagsUID) Flush() error {
	if len(t.tagsIDMap) > 0 {
		writer := stream.BinaryWriter()
		tagNameOffset := stream.BinaryWriter()
		tagTreeWriter := stream.BinaryWriter()
		bitmapOffset := stream.BinaryWriter()
		bitmapWriter := stream.BinaryWriter()

		//sort tag names
		tagNames := t.getSortTagNames()
		for _, tagName := range tagNames {
			tagTree := t.tagsMap[tagName]

			tagNameOffset.PutLenBytes([]byte(tagName))
			tagNameOffset.PutUvarint64(uint64(tagTreeWriter.Len()))

			by, err := tree.NewWriter(tagTree).Encode()
			if nil != err {
				logger.GetLogger("tsdb/index").Error("encode tag tree error:", t.dbField, logger.Error(err))
				return err
			}
			tagTreeWriter.PutUvarint64(uint64(len(by)))
			tagTreeWriter.PutBytes(by)
		}

		for _, bitmap := range t.bitmaps {
			if nil != bitmap {
				bitmapOffset.PutUint32(uint32(bitmapWriter.Len()))

				bitmap.RunOptimize()
				bitmapBytes, err := bitmap.ToBytes()
				if nil != err {
					logger.GetLogger("tsdb/index").Error("encode tags bitmap error:", t.dbField, logger.Error(err))
					return err
				}
				bitmapWriter.PutBytes(bitmapBytes)
			}
		}
		//
		writer.PutUvarint64(uint64(len(tagNames)))
		by, err := tagNameOffset.Bytes()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("encode tag names error: ", t.dbField, logger.Error(err))
			return err
		}
		writer.PutBytes(by)

		//write tagTree
		writer.PutUvarint64(uint64(tagTreeWriter.Len()))
		by, err = tagTreeWriter.Bytes()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("encode tag tree writer error:", t.dbField, logger.Error(err))
			return err
		}
		writer.PutBytes(by)

		//write bitmap offset
		writer.PutUvarint64(uint64(bitmapOffset.Len()))
		by, err = bitmapOffset.Bytes()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("encode bitmap offset error:", t.dbField, logger.Error(err))
			return err
		}
		writer.PutBytes(by)

		//write bitmap
		by, err = bitmapWriter.Bytes()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("encode bitmap writer error:", t.dbField, logger.Error(err))
			return err
		}
		writer.PutBytes(by)

		//flush tagsID to kv-store
		flusher := t.family.NewFlusher()
		by, err = writer.Bytes()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("encode tags data error:", t.dbField, logger.Error(err))
			return err
		}
		err = flusher.Add(t.metricID, by)
		if nil != err {
			logger.GetLogger("tsdb/index").Error("write metric tags error!",
				t.dbField, zap.String("metricID", string(t.metricID)), logger.Error(err))
			return err
		}
		err = flusher.Commit()
		if nil != err {
			logger.GetLogger("tsdb/index").Error("flush metric tags error!", t.dbField, logger.Error(err))
			return err
		}
		t.clear()
	}
	return nil
}
