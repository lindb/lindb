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

package tag

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Concat(t *testing.T) {
	assert.Equal(t, "", Concat(nil))
	assert.Equal(t, "", ConcatKeyValues(KeyValuesFromMap(nil)))
	tags := map[string]string{"t2": "v2", "t1": "v1"}
	assert.Equal(t, "t1=v1,t2=v2", ConcatKeyValues(KeyValuesFromMap(tags)))
	assert.Equal(t, "t1=v1,t2=v2", Concat(tags))
}

func TestConcatTagValues(t *testing.T) {
	assert.Equal(t, "", ConcatTagValues(nil))
	assert.Equal(t, "", ConcatTagValues([]string{}))
	assert.Equal(t, "a", ConcatTagValues([]string{"a"}))
	assert.Equal(t, "a,b", ConcatTagValues([]string{"a", "b"}))
}

func TestSplitTagValues(t *testing.T) {
	assert.Len(t, SplitTagValues(""), 0)
	assert.Equal(t, []string{"a"}, SplitTagValues("a"))
	assert.Equal(t, []string{"a", "b"}, SplitTagValues("a,b"))
	assert.Equal(t, []string{"a", "b", ""}, SplitTagValues("a,b,"))
}

var _testTags = map[string]string{
	"a": "aaaaaaaaa",
	"b": "bbb",
	"c": "ccccc",
	"d": "ddddd",
}

func Benchmark_TagsAsString_old(b *testing.B) {
	tagsAsString := func(tags map[string]string) string {
		if tags == nil {
			return ""
		}
		tagKeyValues := make([]string, 0, len(tags))

		totalLen := 0
		for key, val := range tags {
			keyVal := key + "=" + val
			tagKeyValues = append(tagKeyValues, keyVal)
			totalLen += len(keyVal)
		}
		sort.Strings(tagKeyValues)
		return strings.Join(tagKeyValues, ",")
	}

	for i := 0; i < b.N; i++ {
		tagsAsString(_testTags)
	}
}

func Benchmark_TagsAsString_new(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concat(_testTags)
	}
}

func Test_KeyValues(t *testing.T) {
	// merge
	originMap := map[string]string{
		"host": "alpha",
		"ip":   "1.1.1.1",
		"z":    "32",
	}
	origin := KeyValuesFromMap(originMap)
	assert.Equal(t, origin.Merge(nil), origin)
	other := KeyValuesFromMap(map[string]string{
		"host": "beta",
		"a":    "222",
		"z":    "33",
	})
	merged := origin.Merge(other)
	assert.Equal(t, originMap, origin.Map())
	assert.Equal(t, KeyValues{
		{Key: "a", Value: "222"},
		{Key: "host", Value: "beta"},
		{Key: "ip", Value: "1.1.1.1"},
		{Key: "z", Value: "33"},
	}, merged)

	// append
	origin = KeyValuesFromMap(map[string]string{
		"host": "alpha",
		"ip":   "1.1.1.1",
		"z":    "32",
	})
	other = KeyValuesFromMap(map[string]string{
		"b": "222",
		"c": "33",
		"d": "313",
		"e": "323",
		"f": "333",
	})
	assert.Equal(t, KeyValues{
		{Key: "b", Value: "222"},
		{Key: "c", Value: "33"},
		{Key: "d", Value: "313"},
		{Key: "e", Value: "323"},
		{Key: "f", Value: "333"},
		{Key: "host", Value: "alpha"},
		{Key: "ip", Value: "1.1.1.1"},
		{Key: "z", Value: "32"},
	}, origin.Merge(other))
}

func Test_KeyValuesDeDup(t *testing.T) {
	assert.Len(t, KeyValuesFromMap(map[string]string{"a": "1"}).DeDup(), 1)

	assert.Equal(t, KeyValues{
		{Key: "1", Value: "2"},
		{Key: "2", Value: "4"},
		{Key: "3", Value: "6"},
	}, KeyValues{
		{Key: "2", Value: "4"},
		{Key: "2", Value: "4"},
		{Key: "1", Value: "2"},
		{Key: "3", Value: "6"},
		{Key: "3", Value: "6"},
	}.DeDup())

	assert.Equal(t, KeyValues{
		{Key: "1", Value: "2"},
		{Key: "2", Value: "4"},
		{Key: "3", Value: "6"},
	}, KeyValues{
		{Key: "2", Value: "4"},
		{Key: "1", Value: "2"},
		{Key: "3", Value: "6"},
		{Key: "3", Value: "6"},
	}.DeDup())

	assert.Equal(t, KeyValues{
		{Key: "1", Value: "2"},
		{Key: "2", Value: "4"},
		{Key: "3", Value: "6"},
	}, KeyValues{
		{Key: "2", Value: "4"},
		{Key: "1", Value: "2"},
		{Key: "3", Value: "6"},
	}.DeDup())
}
