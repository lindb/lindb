package metricsmeta

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

//go:generate mockgen -source ./tagkv_entryset.go -destination=./tagkv_entryset_mock.go -package metricsmeta

// for testing
var (
	trieTreeFunc = createTrieTree
)

type TagKVEntrySetINTF interface {
	// TagValueSeq returns the auto sequence of tag value id under this tag key
	TagValueSeq() uint32
	// TagValueIDs returns all tag value ids under this tag key
	TagValueIDs() (*roaring.Bitmap, error)
	// TrieTree builds the trie-tree block for querying
	TrieTree() (trieTreeQuerier, error)
	// GetTagValueID gets tag value id by offset
	GetTagValueID(offset int) uint32
	// CollectTagValues collects the tag values by tag value ids,
	CollectTagValues(tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) error
}

type TagKVEntries []TagKVEntrySetINTF

// GetTagValueIDs gets all tag value ids under tag key entries
func (entries TagKVEntries) GetTagValueIDs() (*roaring.Bitmap, error) {
	unionIDSet := roaring.New()
	for _, entrySet := range entries {
		tagValueIDs, err := entrySet.TagValueIDs()
		if err != nil {
			return nil, err
		}
		unionIDSet.Or(tagValueIDs)
	}
	return unionIDSet, nil
}

// tagKVEntrySet implements tagKVEntrySetINTF
type tagKVEntrySet struct {
	block                []byte
	sr                   *stream.Reader
	tree                 trieTreeQuerier
	crc32CheckSum        uint32
	tagValueSeq          uint32
	tagValueCount        int
	tagValueIDsOffset    int
	tagValueIDWidth      int
	offsetWidth          int
	offsetPos            int
	tagValueBitmapOffset int
}

func newTagKVEntrySet(block []byte) (TagKVEntrySetINTF, error) {
	if len(block) <= tagFooterSize {
		return nil, fmt.Errorf("block length not ok")
	}
	entrySet := &tagKVEntrySet{
		block: block,
		sr:    stream.NewReader(block),
	}
	// read footer(4+4+4+4)
	footerPos := len(block) - tagFooterSize
	entrySet.tagValueSeq = stream.ReadUint32(block, footerPos)
	tagValueIDsStartPos := int(stream.ReadUint32(block, footerPos+4))
	tagValueForwardPos := int(stream.ReadUint32(block, footerPos+8))
	entrySet.crc32CheckSum = stream.ReadUint32(block, footerPos+12)
	// validate offsets
	if !(tagValueIDsStartPos < footerPos) {
		return nil, fmt.Errorf("bad offsets")
	}
	entrySet.tagValueIDWidth = int(block[tagValueIDsStartPos])
	entrySet.tagValueIDsOffset = tagValueIDsStartPos + 1
	entrySet.tagValueCount = (tagValueForwardPos - entrySet.tagValueIDsOffset) / entrySet.tagValueIDWidth
	entrySet.offsetWidth = int(block[tagValueForwardPos])
	entrySet.offsetPos = tagValueForwardPos + 1
	entrySet.tagValueBitmapOffset = tagValueForwardPos + 1 + entrySet.tagValueCount*entrySet.offsetWidth
	return entrySet, nil
}

// TagValueSeq returns the auto sequence of tag value id under this tag key
func (entrySet *tagKVEntrySet) TagValueSeq() uint32 {
	return entrySet.tagValueSeq
}

// TagValueIDs returns all tag value ids under this tag key
func (entrySet *tagKVEntrySet) TagValueIDs() (*roaring.Bitmap, error) {
	tagValueIDs := roaring.New()
	if err := encoding.BitmapUnmarshal(tagValueIDs, entrySet.block[entrySet.tagValueBitmapOffset:]); err != nil {
		return nil, err
	}
	return tagValueIDs, nil
}

// GetTagValueID gets tag value id by index
func (entrySet *tagKVEntrySet) GetTagValueID(index int) uint32 {
	return entrySet.getValue(entrySet.tagValueIDsOffset, entrySet.tagValueIDWidth, index)
}

// CollectTagValues collects the tag values by tag value ids,
func (entrySet *tagKVEntrySet) CollectTagValues(tagValueIDs *roaring.Bitmap, tagValues map[uint32]string) error {
	tagValueIDsInFile, err := entrySet.TagValueIDs()
	if err != nil {
		return err
	}
	needCollectTagValueIDs := roaring.And(tagValueIDs, tagValueIDsInFile)
	// tag value ids not exist in current file
	if needCollectTagValueIDs.IsEmpty() {
		return nil
	}
	// remove found tag value ids
	tagValueIDs.Xor(needCollectTagValueIDs)
	trie, err := trieTreeFunc(entrySet)
	if err != nil {
		return err
	}
	highKeys := tagValueIDsInFile.GetHighKeys()
	idx := 0
	for i, highKey := range highKeys {
		containerInFile := tagValueIDsInFile.GetContainerAtIndex(i)
		collectContainer := needCollectTagValueIDs.GetContainer(highKey)
		if collectContainer != nil {
			it := collectContainer.PeekableIterator()
			for it.HasNext() {
				lowKey := it.Next()
				lowIdx := containerInFile.Rank(lowKey)
				tagValueNodeOffset := entrySet.getValue(entrySet.offsetPos, entrySet.offsetWidth, idx+lowIdx-1)
				tagValue, ok := trie.GetValueByOffset(int(tagValueNodeOffset))
				if ok {
					// flush series id
					hk := uint32(highKey) << 16
					tagValues[encoding.ValueWithHighLowBits(hk, lowKey)] = tagValue
				}
			}
		}

		idx += containerInFile.GetCardinality()
	}
	return nil
}

// createTrieTree creates trie tree
func createTrieTree(entrySet *tagKVEntrySet) (trieTreeQuerier, error) {
	return entrySet.TrieTree()
}

// TrieTree builds the trie-tree block for querying
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

// getTagValueID returns the tag value id by index
func (entrySet *tagKVEntrySet) getValue(offset, width, index int) uint32 {
	offset += index * width
	switch width {
	case 1:
		return uint32(entrySet.block[offset])
	case 2:
		return uint32(entrySet.block[offset]) |
			uint32(entrySet.block[offset+1])<<8
	case 3:
		return uint32(entrySet.block[offset]) |
			uint32(entrySet.block[offset+1])<<8 |
			uint32(entrySet.block[offset+2])<<16
	default:
		return uint32(entrySet.block[offset]) |
			uint32(entrySet.block[offset+1])<<8 |
			uint32(entrySet.block[offset+2])<<16 |
			uint32(entrySet.block[offset+3])<<24
	}
}
