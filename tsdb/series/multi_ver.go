package series

import "github.com/RoaringBitmap/roaring"

// MultiVerSeriesIDSet represents a multi version series ids set, can do and/or/and not operator,
// NOTICE: stores the result in the current bitmap, not safe for goroutine concurrent.
// version-> a bitmap of series ids.
type MultiVerSeriesIDSet struct {
	versions map[int64]*roaring.Bitmap
}

// NewMultiVerSeriesIDSet creates a multi-version series id set
func NewMultiVerSeriesIDSet() *MultiVerSeriesIDSet {
	return &MultiVerSeriesIDSet{
		versions: make(map[int64]*roaring.Bitmap),
	}
}

// Add adds a series id set version to map, if version exist ignore new ids
func (mv *MultiVerSeriesIDSet) Add(version int64, ids *roaring.Bitmap) {
	_, ok := mv.versions[version]
	if !ok {
		mv.versions[version] = ids
	}
}

// And computes the intersection between two set and stores the result in the current set
func (mv *MultiVerSeriesIDSet) And(other *MultiVerSeriesIDSet) {
	// 1. computes the intersection between two version
	var notExist []int64
	for version := range mv.versions {
		_, ok := other.versions[version]
		if !ok {
			notExist = append(notExist, version)
		}
	}
	for _, version := range notExist {
		delete(mv.versions, version)
	}

	// 2. computes the intersection between tow bitmap with same version
	for version, ids := range mv.versions {
		ids.And(other.versions[version])
	}
}

// Or computes the union between two set and stores the result in the current set
func (mv *MultiVerSeriesIDSet) Or(other *MultiVerSeriesIDSet) {
	for version, ids := range other.versions {
		existIDs, ok := mv.versions[version]
		if ok {
			existIDs.Or(ids)
		} else {
			mv.versions[version] = ids
		}
	}
}

// AndNot computes the difference between two set and stores the result in the current set
func (mv *MultiVerSeriesIDSet) AndNot(other *MultiVerSeriesIDSet) {
	for version, ids := range other.versions {
		existIDs, ok := mv.versions[version]
		if ok {
			existIDs.AndNot(ids)
		}
	}
}
