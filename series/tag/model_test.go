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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tags(t *testing.T) {
	var tags = Tags{}
	assert.Len(t, tags.AppendHashKey(nil), 0)
	tags = append(tags, NewTag([]byte("ip"), []byte("1.1.1.1")),
		NewTag([]byte("zone"), []byte("sh")),
		NewTag([]byte("host"), []byte("test")))
	assert.Equal(t, 23, tags.Size())
	assert.False(t, tags.needsEscape())
	assert.Equal(t, ",ip=1.1.1.1,zone=sh,host=test", tags.String())

	tags = append(tags, NewTag([]byte("x x"), []byte("y,y")))
	sort.Sort(tags)
	assert.True(t, tags.needsEscape())
	assert.Equal(t, ",host=test,ip=1.1.1.1,x\\ x=y\\,y,zone=sh", tags.String())
}
