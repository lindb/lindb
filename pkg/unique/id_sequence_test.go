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

package unique

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSequence(t *testing.T) {
	cache := uint32(100)
	cases := []struct {
		name  string
		start uint32
	}{
		{
			name:  "start with 0",
			start: 0,
		},
		{
			name:  "start with value",
			start: 66666,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := NewSequence(tt.start, cache)
			for i := 0; i < 9090; i++ {
				if !s.HasNext() {
					s.Limit(s.Current() + cache)
				}
				assert.Equal(t, uint32(i+1)+tt.start, s.Next())
			}
			assert.Equal(t, 9090+tt.start, s.Current())
		})
	}
}

func TestSaveSequence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockIDStore(ctrl)
	store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	assert.Error(t, SaveSequence(store, []byte("key"), 1))
}
