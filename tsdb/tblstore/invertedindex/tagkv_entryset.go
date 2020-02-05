package invertedindex

import (
	"encoding/binary"
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./tagkv_entryset.go -destination=./tagkv_entryset_mock.go -package invertedindex

type TagKVEntrySetINTF interface {
	// TagValuesCount returns the count of tag values under this tag key
	TagValuesCount() int
	// TagValueIDs returns all tag value ids under this tag key
	TagValueIDs() *roaring.Bitmap
	// TrieTree builds the trie-tree block for querying
	TrieTree() (trieTreeQuerier, error)
	// GetTagValueID gets tag value id by offset
	GetTagValueID(offset int) uint32
}

type TagKVEntries []TagKVEntrySetINTF

func (entries TagKVEntries) TagValuesCount() (count int) {
	for _, entry := range entries {
		count += entry.TagValuesCount()
	}
	return
}

// GetTagValueIDs gets all tag value ids under tag key entries
func (entries TagKVEntries) GetTagValueIDs() *roaring.Bitmap {
	unionIDSet := roaring.New()
	for _, entrySet := range entries {
		unionIDSet.Or(entrySet.TagValueIDs())
	}
	return unionIDSet
}

// tagKVEntrySet implements tagKVEntrySetINTF
type tagKVEntrySet struct {
	sr            *stream.Reader
	tree          trieTreeQuerier
	crc32CheckSum uint32
	tagValueIDs   *encoding.FixedOffsetDecoder
}

func newTagKVEntrySet(block []byte) (TagKVEntrySetINTF, error) {
	if len(block) <= tagFooterSize {
		return nil, fmt.Errorf("block length no ok")
	}
	entrySet := &tagKVEntrySet{
		sr: stream.NewReader(block),
	}
	// read footer(4+4)
	footerPos := len(block) - tagFooterSize
	tagValueIDsStartPos := int(binary.LittleEndian.Uint32(block[footerPos : footerPos+4]))
	entrySet.crc32CheckSum = binary.LittleEndian.Uint32(block[footerPos+4 : footerPos+8])
	// validate offsets
	if !(tagValueIDsStartPos < footerPos) {
		return nil, fmt.Errorf("bad offsets")
	}
	entrySet.tagValueIDs = encoding.NewFixedOffsetDecoder(block[tagValueIDsStartPos:footerPos])
	return entrySet, nil
}

// TagValuesCount returns the count of tag values under this tag key
func (entrySet *tagKVEntrySet) TagValuesCount() int {
	return entrySet.tagValueIDs.Size()
}

// TagValueIDs returns all tag value ids under this tag key
func (entrySet *tagKVEntrySet) TagValueIDs() *roaring.Bitmap {
	size := entrySet.tagValueIDs.Size()
	tagValueIDs := roaring.New()
	for i := 0; i < size; i++ {
		tagValueIDs.Add(uint32(entrySet.tagValueIDs.Get(i)))
	}
	return tagValueIDs
}

// GetTagValueID gets tag value id by offset
func (entrySet *tagKVEntrySet) GetTagValueID(offset int) uint32 {
	return uint32(entrySet.tagValueIDs.Get(offset))
}

func (entrySet *tagKVEntrySet) TrieTree() (trieTreeQuerier, error) {
	var tree trieTreeBlock
	////////////////////////////////
	// Block: LOUDS Trie-Tree
	////////////////////////////////
	// read trie-tree length
	expectedTrieTreeLen := entrySet.sr.ReadUvarint64()
	startPosOfTree := entrySet.sr.Position()
	// read label length
	labelsLen := entrySet.sr.ReadUvarint64()
	// read labels block
	tree.labels = entrySet.sr.ReadSlice(int(labelsLen))
	// read isPrefix length
	isPrefixKeyLen := entrySet.sr.ReadUvarint64()
	// read isPrefixKey bitmap
	isPrefixBlock := entrySet.sr.ReadSlice(int(isPrefixKeyLen))
	// read LOUDS length
	loudsLen := entrySet.sr.ReadUvarint64()
	// read LOUDS block
	LOUDSBlock := entrySet.sr.ReadSlice(int(loudsLen))
	// validation of stream error
	if entrySet.sr.Error() != nil {
		return nil, entrySet.sr.Error()
	}
	// validation of length
	if entrySet.sr.Position()-startPosOfTree != int(expectedTrieTreeLen) {
		return nil, fmt.Errorf("failed validation of trie-tree")
	}
	// unmarshal LOUDS block to rank-select
	tree.LOUDS = NewRankSelect()
	if err := tree.LOUDS.UnmarshalBinary(LOUDSBlock); err != nil {
		return nil, err
	}
	// unmarshal isPrefixKey block to rank-select
	tree.isPrefixKey = NewRankSelect()
	if err := tree.isPrefixKey.UnmarshalBinary(isPrefixBlock); err != nil {
		return nil, err
	}
	entrySet.tree = &tree
	return entrySet.tree, nil
}
