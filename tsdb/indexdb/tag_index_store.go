package indexdb

import "github.com/lindb/roaring"

// tagIndexStore represents all tag indexes storage
type tagIndexStore struct {
	tagKeyIDs *roaring.Bitmap // store all tag key ids
	indexes   [][]TagIndex    // store all tag indexes by high/low key
}

// newTagIndexStore creates a tag index stare
func newTagIndexStore() *tagIndexStore {
	return &tagIndexStore{
		tagKeyIDs: roaring.New(),
	}
}

// get returns tag index by tag key id, if exist returns it, else returns nil, false
func (m *tagIndexStore) get(tagKeyIDs uint32) (TagIndex, bool) {
	if len(m.indexes) == 0 {
		return nil, false
	}
	found, highIdx, lowIdx := m.tagKeyIDs.ContainsAndRank(tagKeyIDs)
	if !found {
		return nil, false
	}
	return m.indexes[highIdx][lowIdx-1], true
}

// put puts the tag index by tag key id
func (m *tagIndexStore) put(tagKeyID uint32, tagIndex TagIndex) {
	if len(m.indexes) == 0 {
		// if indexes is empty, append new low container directly
		m.tagKeyIDs.Add(tagKeyID)
		m.indexes = append(m.indexes, []TagIndex{tagIndex})
		return
	}

	// try find tag key id if exist
	found, highIdx, lowIdx := m.tagKeyIDs.ContainsAndRank(tagKeyID)
	if !found {
		// not found
		m.tagKeyIDs.Add(tagKeyID)
		if highIdx >= 0 {
			// high container exist
			stores := m.indexes[highIdx]
			// insert operation
			stores = append(stores, nil)
			copy(stores[lowIdx+1:], stores[lowIdx:len(stores)-1])
			stores[lowIdx] = tagIndex
			m.indexes[highIdx] = stores
		} else {
			// high container not exist, append operation
			m.indexes = append(m.indexes, []TagIndex{tagIndex})
		}
	}
}

// getAllTagKeyIDs returns the all tag key ids
func (m *tagIndexStore) getAllTagKeyIDs() *roaring.Bitmap {
	return m.tagKeyIDs
}
