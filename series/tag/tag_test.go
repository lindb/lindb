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
	tags := map[string]string{"t2": "v2", "t1": "v1"}
	tagsStr := Concat(tags)
	assert.Equal(t, "t1=v1,t2=v2", tagsStr)
}

func TestConcatTagValues(t *testing.T) {
	assert.Equal(t, emptyStr, ConcatTagValues(nil))
	assert.Equal(t, emptyStr, ConcatTagValues([]string{}))
	assert.Equal(t, "a", ConcatTagValues([]string{"a"}))
	assert.Equal(t, "a,b", ConcatTagValues([]string{"a", "b"}))
}

func TestSplitTagValues(t *testing.T) {
	assert.Equal(t, emptyArray, SplitTagValues(""))
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
