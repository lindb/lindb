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

package field

import (
	"sort"
	"strings"
)

// Meta is the meta-data for field, which contains field-name, fieldID and field-type
type Meta struct {
	ID   ID   `json:"id"`   // query not use id, don't get id in query phase
	Type Type `json:"type"` // query not use type
	Name Name `json:"name"`
}

// Metas implements sort.Interface, it's sorted by name
type Metas []Meta

func (fms Metas) Len() int           { return len(fms) }
func (fms Metas) Less(i, j int) bool { return fms[i].Name < fms[j].Name }
func (fms Metas) Swap(i, j int)      { fms[i], fms[j] = fms[j], fms[i] }

// GetFromName searches the meta by fieldName, return false when not exist
func (fms Metas) GetFromName(fieldName Name) (Meta, bool) {
	idx := sort.Search(len(fms), func(i int) bool { return fms[i].Name >= fieldName })
	if idx >= len(fms) || fms[idx].Name != fieldName {
		return Meta{}, false
	}
	return fms[idx], true
}

// GetFromID searches the meta by fieldID, returns false when not exist
func (fms Metas) GetFromID(fieldID ID) (Meta, bool) {
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
	copy(x2, fms)
	return x2
}

// Insert appends a new Meta to the list and sort it.
func (fms Metas) Insert(m Meta) Metas {
	newFms := append(fms, m)
	return newFms
}

// Intersects checks whether each fieldID is in the list,
// and returns the new meta-list corresponding with the fieldID-list.
func (fms Metas) Intersects(fields Metas) (x2 Metas, isSubSet bool) {
	isSubSet = true
	for _, f := range fields {
		fm, ok := fms.GetFromID(f.ID)
		if ok {
			x2 = append(x2, fm)
		} else {
			isSubSet = false
		}
	}
	sort.Sort(x2)
	return x2, isSubSet
}

// Stringer returns the fields in string
func (fms Metas) String() string {
	switch len(fms) {
	case 0:
		return ""
	case 1:
		return string(fms[0].Name)
	case 2:
		return string(fms[0].Name) + "," + string(fms[1].Name)
	default:
		b := strings.Builder{}
		for i := range fms {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(string(fms[i].Name))
		}
		return b.String()
	}
}
