package field

import (
	"sort"
)

// Meta is the meta-data for field, which contains field-name, fieldID and field-type
type Meta struct {
	ID   uint16 // query not use ID, don't get id in query phase
	Type Type   // query not user type
	Name string
}

// Metas implements sort.Interface, it's sorted by name
type Metas []Meta

func (fms Metas) Len() int           { return len(fms) }
func (fms Metas) Less(i, j int) bool { return fms[i].Name < fms[j].Name }
func (fms Metas) Swap(i, j int)      { fms[i], fms[j] = fms[j], fms[i] }

// GetFromName searches the meta by fieldName, return false when not exist
func (fms Metas) GetFromName(fieldName string) (Meta, bool) {
	idx := sort.Search(len(fms), func(i int) bool { return fms[i].Name >= fieldName })
	if idx >= len(fms) || fms[idx].Name != fieldName {
		return Meta{}, false
	}
	return fms[idx], true
}

// GetFromID searches the meta by fieldID, returns false when not exist
func (fms Metas) GetFromID(fieldID uint16) (Meta, bool) {
	for _, fm := range fms {
		if fm.ID == fieldID {
			return fm, true
		}
	}
	return Meta{}, false
}

// Clone clones a copy of fieldsMetas
func (fms Metas) Clone() (x2 Metas) {
	x2 = make([]Meta, fms.Len())
	for idx, fm := range fms {
		x2[idx] = fm
	}
	return x2
}

// Insert appends a new Meta to the list and sort it.
func (fms Metas) Insert(m Meta) Metas {
	newFms := append(fms, m)
	sort.Sort(fms)
	return newFms
}

// Intersects checks whether each fieldID is in the list,
// and returns the new meta-list corresponding with the fieldID-list.
func (fms Metas) Intersects(fieldIDs []uint16) (x2 Metas, isSubSet bool) {
	isSubSet = true
	for _, fieldID := range fieldIDs {
		fm, ok := fms.GetFromID(fieldID)
		if ok {
			x2 = append(x2, fm)
		} else {
			isSubSet = false
		}
	}
	sort.Sort(x2)
	return x2, isSubSet
}
