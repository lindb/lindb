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

package model

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/pkg/trie"
)

func TestTrieBucketBuilder_Error(t *testing.T) {
	t.Run("io write error", func(t *testing.T) {
		w := &mockWriter{}
		b := NewTrieBucketBuilder(math.MaxUint16, w)
		assert.Error(t, b.Write([][]byte{{1}}, []uint32{1}))
	})

	t.Run("trie write error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer func() {
			ctrl.Finish()
			newTrieBuilder = trie.NewBuilder
		}()
		bm := trie.NewMockBuilder(ctrl)
		newTrieBuilder = func() trie.Builder {
			return bm
		}
		w := bytes.NewBuffer([]byte{})
		gomock.InOrder(
			bm.EXPECT().Reset(),
			bm.EXPECT().Build(gomock.Any(), gomock.Any()),
			bm.EXPECT().MarshalSize().Return(10),
			bm.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("err")),
		)
		b := NewTrieBucketBuilder(math.MaxUint16, w)
		assert.Error(t, b.Write([][]byte{{1}}, []uint32{1}))
	})
}
