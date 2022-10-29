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

package tsdb

import (
	"sync"
)

var (
	fManager           FamilyManager
	once4FamilyManager sync.Once
)

// GetFamilyManager returns the data family manager singleton instance.
// FIXME: need clean readonly family when no read long term
func GetFamilyManager() FamilyManager {
	once4FamilyManager.Do(func() {
		fManager = newFamilyManager()
	})
	return fManager
}

// FamilyManager represents the data family manager.
type FamilyManager interface {
	// AddFamily adds the family.
	AddFamily(family DataFamily)
	// RemoveFamily removes the family.
	RemoveFamily(family DataFamily)
	// WalkEntry walks each family entry via fn.
	WalkEntry(fn func(family DataFamily))
	// GetFamiliesByShard returns families for spec shard.
	GetFamiliesByShard(shard Shard) []DataFamily
}

// familyManager implements FamilyManager interface.
type familyManager struct {
	families sync.Map
}

// newFamilyManager creates the family manager.
func newFamilyManager() FamilyManager {
	return &familyManager{}
}

// AddFamily adds the family.
func (sm *familyManager) AddFamily(family DataFamily) {
	sm.families.Store(family.Indicator(), family)
}

// RemoveFamily removes the family.
func (sm *familyManager) RemoveFamily(family DataFamily) {
	sm.families.Delete(family.Indicator())
}

// WalkEntry walks each family entry via fn.
func (sm *familyManager) WalkEntry(fn func(family DataFamily)) {
	sm.families.Range(func(_, value interface{}) bool {
		family := value.(DataFamily)
		fn(family)
		return true
	})
}

// GetFamiliesByShard returns families for spec shard.
func (sm *familyManager) GetFamiliesByShard(shard Shard) (rs []DataFamily) {
	sm.families.Range(func(_, value interface{}) bool {
		family := value.(DataFamily)
		if family.Shard().Indicator() == shard.Indicator() {
			rs = append(rs, family)
		}
		return true
	})
	return
}
