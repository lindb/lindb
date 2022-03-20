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

package linmetric //nolint:dupl

import (
	"strings"
	"sync"

	"github.com/lindb/lindb/series/tag"
)

type MaxVec struct {
	r          *Registry
	tags       tag.Tags // unique tags
	tagKeys    []string
	metricName string // concated metric name
	fieldName  string
	mu         sync.RWMutex
	maxes      map[string]*BoundMax
}

func newMaxVec(r *Registry, metricName string, fieldName string, tags tag.Tags, tagKey ...string) *MaxVec {
	return &MaxVec{
		r:          r,
		metricName: metricName,
		fieldName:  fieldName,
		tags:       tags,
		tagKeys:    tagKey,
		maxes:      make(map[string]*BoundMax),
	}
}

func (gv *MaxVec) WithTagValues(tagValues ...string) *BoundMax {
	if len(tagValues) != len(gv.tagKeys) {
		panic("count of tagKey and tagValue not match")
	}
	id := strings.Join(tagValues, ",")
	gv.mu.RLock()
	c, ok := gv.maxes[id]
	gv.mu.RUnlock()
	if ok {
		return c
	}

	gv.mu.Lock()
	defer gv.mu.Unlock()

	c, ok = gv.maxes[id]
	if ok {
		return c
	}
	var tagsMap = gv.tags.Map()
	for i := range gv.tagKeys {
		tagsMap[gv.tagKeys[i]] = tagValues[i]
	}
	series := newTaggedSeries(gv.r, gv.metricName, tag.TagsFromMap(tagsMap))
	c = series.NewMax(gv.fieldName)

	gv.maxes[id] = c
	return c
}
