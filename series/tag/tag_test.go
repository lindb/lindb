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
	"sync"
	"testing"

	xxhash "github.com/cespare/xxhash/v2"
	protoMetricsV1 "github.com/lindb/common/proto/gen/v1/linmetrics"
	"github.com/stretchr/testify/assert"
)

func TestConcatTagValues(t *testing.T) {
	assert.Equal(t, "", ConcatTagValues(nil))
	assert.Equal(t, "", ConcatKeyValues(KeyValuesFromMap(nil)))
	assert.Equal(t, "", ConcatTagValues([]string{}))
	tags := map[string]string{"t2": "v2", "t1": "v1"}
	assert.Equal(t, "t1=v1,t2=v2", ConcatKeyValues(KeyValuesFromMap(tags)))
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

func Test_XXHashOfKeyValues(t *testing.T) {
	assert.Equal(t, xxhash.Sum64String(""), XXHashOfKeyValues(nil))
}

func Test_getSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := getSlice(10)
		assert.Len(t, *s, 10)
	}
	for i := 0; i < 100; i++ {
		s := getSlice(100)
		assert.Len(t, *s, 100)
		putSlice(s)
	}
	for i := 0; i < 100; i++ {
		s := getSlice(1000)
		assert.Len(t, *s, 1000)
	}
}

var (
	singleKeyValues KeyValues = []*protoMetricsV1.KeyValue{{Key: "env", Value: "prd"}}
	logKeyValues    KeyValues = []*protoMetricsV1.KeyValue{
		{Key: "333339", Value: "22222222222222222211111"},
		{Key: "333338", Value: "22222222222222222211111"},
		{Key: "333337", Value: "22222222222222222211111"},
		{Key: "333336", Value: "22222222222222222211111"},
		{Key: "333335", Value: "22222222222222222211111"},
		{Key: "333334", Value: "22222222222222222211111"},
		{Key: "33333", Value: "22222222222222222211111"},
		{Key: "1", Value: "11111111111111111111111111"},
		{Key: "2222", Value: "22222222222222222211111"},
	}
	commonKeyValues KeyValues = []*protoMetricsV1.KeyValue{
		{Key: "ip", Value: "1.1.1.1"},
		{Key: "host", Value: "alpha-test-machine"},
		{Key: "region", Value: "shanghai"},
		{Key: "1", Value: "2"},
		{Key: "2222", Value: "211111"},
		{Key: "env", Value: "prd"},
	}
)

// on stack concat Sum64, 280ns/op, no escape
// digest with pool,  WriteString..., 687ns/op (too much memmove)
// use bytes buffer, 590ns/op(too much strings.Builder.Grow)
func Benchmark_HashOfConcatTagValues(b *testing.B) {
	kvs := commonKeyValues.Clone()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = XXHashOfKeyValues(kvs)
	}
}

func Benchmark_HashOfConcatTagValues_Sorted(b *testing.B) {
	sorted := commonKeyValues.Clone()
	sort.Sort(sorted)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = XXHashOfKeyValues(sorted)
	}
}

func Benchmark_HashOfConcatTagValues_OnHeap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = XXHashOfKeyValues(logKeyValues)
	}
}

func Benchmark_HashOfConcatTagValues_Single(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = XXHashOfKeyValues(singleKeyValues)
	}
}

func TestXXHashOfKeyValues(t *testing.T) {
	assert.Equal(t, xxhash.Sum64String(""), XXHashOfKeyValues(nil))
	_ = XXHashOfKeyValues(singleKeyValues)
	_ = XXHashOfKeyValues(commonKeyValues)
	_ = XXHashOfKeyValues(logKeyValues)
	long2 := logKeyValues.Clone().Merge(commonKeyValues)
	_ = XXHashOfKeyValues(long2)
	_ = XXHashOfKeyValues(logKeyValues)
	sorted := commonKeyValues.Clone()
	sort.Sort(sorted)
	_ = XXHashOfKeyValues(sorted)
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

func TestTag_Pool(t *testing.T) {
	defer func() {
		slicePool = sync.Pool{}
	}()
	slicePool = sync.Pool{}
	for i := 0; i < 10; i++ {
		s := getSlice(10)
		assert.Len(t, *s, 10)
		putSlice(s)
	}

	assert.Len(t, *(getSlice(100)), 100)
}
