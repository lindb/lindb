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

package tagkeymeta

import (
	"bytes"
	"regexp"
	"sort"
	"strings"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/trie"

	"github.com/lindb/roaring"
)

//go:generate mockgen -source ./meta.go -destination=./meta_mock.go -package tagkeymeta

type TagKeyMeta interface {
	// TagValueIDSeq returns the auto sequence of tag value id under this tag key
	TagValueIDSeq() uint32
	// TagValueIDs returns all tag value ids under this tag key
	TagValueIDs() (*roaring.Bitmap, error)
	// CollectTagValues collects the tag values by tag value ids,
	CollectTagValues(tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) error
	// TrieTree returns the trie tree
	TrieTree() (trie.SuccinctTrie, error)
	// PrefixIterator returns a iterator for prefix iterating the Iterator
	PrefixIterator(tagValuePrefix []byte) (*trie.PrefixIterator, error)
	// FindTagValueID finds tagValueIDs by tagValue
	FindTagValueID(tagValue string) (tagValueIDs []uint32)
	// FindTagValueIDs finds tagValueIDs in tagValue
	FindTagValueIDs(tagValues []string) (tagValueIDs []uint32)
	// FindTagValueIDsByLike finds tagValueIDs like tagValue,
	// 3 cases: *sdb, ts*, *sd*
	FindTagValueIDsByLike(tagValue string) (tagValueIDs []uint32)
	// FindTagValueIDsByRegex finds tagValueIDs by regex pattern,
	FindTagValueIDsByRegex(tagValuePattern string) (tagValueIDs []uint32)
}

const (
	tagFooterSize = 4 + // bitmap position
		4 + // offsets position
		4 + // tag value sequence
		4 // crc32 checksum
)

type TagKeyMetas []TagKeyMeta

// GetTagValueIDs gets all tag value ids under tag-keys meta
func (metas TagKeyMetas) GetTagValueIDs() (*roaring.Bitmap, error) {
	unionIDSet := roaring.New()
	for _, tagMeta := range metas {
		tagValueIDs, err := tagMeta.TagValueIDs()
		if err != nil {
			return nil, err
		}
		unionIDSet.Or(tagValueIDs)
	}
	return unionIDSet, nil
}

// tagKeyMeta implements TagKeyMeta
type tagKeyMeta struct {
	block          []byte
	sr             *stream.Reader
	tree           trie.SuccinctTrie
	unmarshalError error
	offsetsDecoder *encoding.FixedOffsetDecoder
	trieBlock      []byte
	bitmapData     []byte
	offsetsData    []byte
	footerPos      int
	bitmapPos      int
	offsetsPos     int
	tagValueIDSeq  uint32
	crc32CheckSum  uint32
}

func newTagKeyMeta(tagKeyMetaBlock []byte) (TagKeyMeta, error) {
	if len(tagKeyMetaBlock) < tagFooterSize {
		return nil, constants.ErrDataFileCorruption
	}
	meta := &tagKeyMeta{
		block: tagKeyMetaBlock,
		sr:    stream.NewReader(tagKeyMetaBlock),
	}
	// read footer(4+4+4+4+4)
	meta.footerPos = len(tagKeyMetaBlock) - tagFooterSize
	meta.sr.ReadAt(meta.footerPos)
	meta.bitmapPos = int(meta.sr.ReadUint32())
	meta.offsetsPos = int(meta.sr.ReadUint32())
	meta.tagValueIDSeq = meta.sr.ReadUint32()
	meta.crc32CheckSum = meta.sr.ReadUint32()

	expectedOrders := []int{0,
		meta.bitmapPos, meta.bitmapPos + 1,
		meta.offsetsPos, meta.offsetsPos + 1,
		meta.footerPos}
	// data validation
	if !sort.IntsAreSorted(expectedOrders) {
		return nil, constants.ErrDataFileCorruption
	}
	// read trie block data, lazy unmarshal
	meta.sr.SeekStart()
	meta.trieBlock = meta.sr.ReadSlice(meta.bitmapPos) // 0->bitmap pos
	// read bitmap data, lazy unmarshal
	meta.bitmapData = meta.sr.ReadSlice(meta.offsetsPos - meta.bitmapPos)
	// read offsets data, lazy unmarshal
	meta.offsetsData = meta.sr.ReadSlice(meta.footerPos - meta.offsetsPos)
	return meta, nil
}

func (meta *tagKeyMeta) TagValueIDSeq() uint32 {
	return meta.tagValueIDSeq
}

func (meta *tagKeyMeta) TagValueIDs() (*roaring.Bitmap, error) {
	tagValueIDs := roaring.New()
	if err := encoding.BitmapUnmarshal(tagValueIDs, meta.bitmapData); err != nil {
		return nil, err
	}
	return tagValueIDs, nil
}

func (meta *tagKeyMeta) TrieTree() (trie.SuccinctTrie, error) {
	if meta.tree == nil && meta.unmarshalError == nil {
		meta.tree = trie.NewTrie()
		meta.unmarshalError = meta.tree.UnmarshalBinary(meta.trieBlock)
	}
	return meta.tree, meta.unmarshalError
}

// idRanksOffsets sorts ids slice based on the order in ranks
type idRanksOffsets struct {
	ids     []uint32 // tag-value ids
	keys    []string // tag-values
	ranks   []int
	offsets []uint32
}

func makeIDRankOffsets(size int) idRanksOffsets {
	return idRanksOffsets{
		ids:     make([]uint32, size)[:0],
		keys:    make([]string, size)[:0],
		ranks:   make([]int, size)[:0],
		offsets: make([]uint32, size)[:0],
	}
}
func (ir idRanksOffsets) Len() int           { return len(ir.ranks) }
func (ir idRanksOffsets) Less(i, j int) bool { return ir.ranks[i] < ir.ranks[j] }
func (ir idRanksOffsets) Swap(i, j int) {
	ir.ranks[i], ir.ranks[j] = ir.ranks[j], ir.ranks[i]
	ir.ids[i], ir.ids[j] = ir.ids[j], ir.ids[i]
}

func (meta *tagKeyMeta) CollectTagValues(tagValueIds *roaring.Bitmap, tagValues map[uint32]string) error {
	tagValueIDsInFile, err := meta.TagValueIDs()
	if err != nil {
		return err
	}
	needCollectTagValueIDs := roaring.And(tagValueIds, tagValueIDsInFile)
	// tag value ids not exist in current file
	if needCollectTagValueIDs.IsEmpty() {
		return nil
	}
	// remove found tag value ids
	tagValueIds.Xor(needCollectTagValueIDs)
	// pre-allocate all buffer data
	mappings := makeIDRankOffsets(int(tagValueIds.GetCardinality()))

	highKeys := tagValueIDsInFile.GetHighKeys()
	idx := 0
	for i, highKey := range highKeys {
		hk := uint32(highKey) << 16
		containerInFile := tagValueIDsInFile.GetContainerAtIndex(i)
		collectContainer := needCollectTagValueIDs.GetContainer(highKey)
		if collectContainer != nil {
			it := collectContainer.PeekableIterator()
			for it.HasNext() {
				lowKey := it.Next()
				lowIdx := containerInFile.Rank(lowKey)
				// calculate offset
				offset := idx + lowIdx - 1
				mappings.offsets = append(mappings.offsets, uint32(offset))
				// calculate id
				id := encoding.ValueWithHighLowBits(hk, lowKey)
				mappings.ids = append(mappings.ids, id)
			}
		}
		idx += containerInFile.GetCardinality()
	}
	if err := meta.walkTrieTree(&mappings); err != nil {
		return err
	}
	for i, id := range mappings.ids {
		tagValues[id] = mappings.keys[i]
	}
	return nil
}

func (meta *tagKeyMeta) walkTrieTree(mappings *idRanksOffsets) error {
	tree, err := meta.TrieTree()
	if err != nil {
		return err
	}
	if meta.offsetsDecoder == nil {
		meta.offsetsDecoder = encoding.NewFixedOffsetDecoder(meta.offsetsData)
	}

	for _, offset := range mappings.offsets {
		rank, ok := meta.offsetsDecoder.Get(int(offset))
		if !ok {
			return constants.ErrDataFileCorruption
		}
		mappings.ranks = append(mappings.ranks, rank)
	}
	sort.Sort(mappings)

	itr := tree.NewIterator()
	itr.SeekToFirst()

	expectedRankIdx := 0 // pop left from ranks
	walkedRankAt := 0
	for itr.Valid() {
		if expectedRankIdx >= len(mappings.ranks) {
			break
		}
		if mappings.ranks[expectedRankIdx] == walkedRankAt {
			mappings.keys = append(mappings.keys, string(itr.Key()))
			expectedRankIdx++
		}
		itr.Next()
		walkedRankAt++
	}
	if len(mappings.keys) != len(mappings.ranks) {
		return constants.ErrDataFileCorruption
	}
	return nil
}

func (meta *tagKeyMeta) FindTagValueID(tagValue string) (tagValueIDs []uint32) {
	tree, err := meta.TrieTree()
	if err != nil {
		return nil
	}
	slice, ok := tree.Get([]byte(tagValue))
	if !ok {
		return nil
	}
	return []uint32{encoding.ByteSlice2Uint32(slice)}
}

func (meta *tagKeyMeta) FindTagValueIDs(tagValues []string) (tagValueIDs []uint32) {
	for _, tagValue := range tagValues {
		tagValueIDs = append(tagValueIDs, meta.FindTagValueID(tagValue)...)
	}
	return tagValueIDs
}

func (meta *tagKeyMeta) PrefixIterator(tagValuePrefix []byte) (*trie.PrefixIterator, error) {
	tree, err := meta.TrieTree()
	if err != nil {
		return nil, err
	}
	return tree.NewPrefixIterator(tagValuePrefix), nil
}

func (meta *tagKeyMeta) FindTagValueIDsByLike(tagValue string) (tagValueIDs []uint32) {
	hashPrefix := strings.HasPrefix(tagValue, "*")
	hasSuffix := strings.HasSuffix(tagValue, "*")
	tagValueSlice := strutil.String2ByteSlice(tagValue)
	switch {
	case tagValue == "":
		break
	// only endswith *
	case !hashPrefix && hasSuffix:
		itr, err := meta.PrefixIterator(tagValueSlice[:len(tagValueSlice)-1])
		if err != nil {
			return nil
		}
		for itr.Valid() {
			tagValueIDs = append(tagValueIDs, encoding.ByteSlice2Uint32(itr.Value()))
			itr.Next()
		}
	// only startswith *
	case hashPrefix && !hasSuffix:
		suffix := tagValueSlice[1:]
		itr, err := meta.PrefixIterator(nil)
		if err != nil {
			return nil
		}
		for itr.Valid() {
			if bytes.HasSuffix(itr.Key(), suffix) {
				tagValueIDs = append(tagValueIDs, encoding.ByteSlice2Uint32(itr.Value()))
			}
			itr.Next()
		}
	// startswith and endswith *
	case hashPrefix && hasSuffix:
		middle := tagValueSlice[1 : len(tagValueSlice)-1]
		itr, err := meta.PrefixIterator(nil)
		if err != nil {
			return nil
		}
		for itr.Valid() {
			if bytes.Contains(itr.Key(), middle) {
				tagValueIDs = append(tagValueIDs, encoding.ByteSlice2Uint32(itr.Value()))
			}
			itr.Next()
		}
	default:
		return meta.FindTagValueID(tagValue)
	}
	return tagValueIDs
}

func (meta *tagKeyMeta) FindTagValueIDsByRegex(tagValuePattern string) (tagValueIDs []uint32) {
	rp, err := regexp.Compile(tagValuePattern)
	if err != nil {
		return nil
	}
	literalPrefix, _ := rp.LiteralPrefix()
	literalPrefixByte := strutil.String2ByteSlice(literalPrefix)
	itr, err := meta.PrefixIterator(literalPrefixByte)
	if err != nil {
		return nil
	}
	for itr.Valid() {
		if rp.Match(itr.Key()) {
			tagValueIDs = append(tagValueIDs, encoding.ByteSlice2Uint32(itr.Value()))
		}
		itr.Next()
	}
	return tagValueIDs
}
