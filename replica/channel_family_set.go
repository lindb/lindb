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

package replica

import (
	"sync/atomic"
)

type channelEntry struct {
	familyTime int64
	channel    FamilyChannel
}

type familyChannelSet struct {
	value atomic.Value // channelEntry
}

func newFamilyChannelSet() *familyChannelSet {
	set := &familyChannelSet{}
	set.value.Store(make([]channelEntry, 0))
	return set
}

func (ss *familyChannelSet) InsertFamily(familyTime int64, channel FamilyChannel) {
	set := ss.value.Load().([]channelEntry)

	newSet := make([]channelEntry, 0, len(set)+1)

	newSet = append(newSet, set...)
	newSet = append(newSet, channelEntry{
		familyTime: familyTime,
		channel:    channel,
	})

	ss.value.Store(newSet)
}

func (ss *familyChannelSet) GetFamilyChannel(familyTime int64) (FamilyChannel, bool) {
	set := ss.value.Load().([]channelEntry)
	for idx := range set {
		if set[idx].familyTime == familyTime {
			return set[idx].channel, true
		}
	}
	return nil, false
}

func (ss *familyChannelSet) Entries() []FamilyChannel {
	set := ss.value.Load().([]channelEntry)
	dst := make([]FamilyChannel, len(set))
	for idx := range set {
		dst[idx] = set[idx].channel
	}
	return dst
}

// RemoveFamilies removes given families from set.
func (ss *familyChannelSet) RemoveFamilies(needRemoveFamilies map[int64]struct{}) {
	if len(needRemoveFamilies) == 0 {
		return
	}
	set := ss.value.Load().([]channelEntry)
	dst := make([]channelEntry, 0)

	for _, family := range set {
		_, ok := needRemoveFamilies[family.familyTime]
		if !ok {
			dst = append(dst, family)
		}
	}
	ss.value.Store(dst)
}
