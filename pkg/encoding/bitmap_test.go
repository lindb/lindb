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

package encoding

import (
	"fmt"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestBitmapMarshal(t *testing.T) {
	defer func() {
		BitmapMarshal = bitmapMarshal
		BitmapUnmarshal = bitmapUnmarshal
	}()
	data, err := BitmapMarshal(roaring.BitmapOf(1))
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)
	BitmapMarshal = func(bitmap *roaring.Bitmap) (bytes []byte, err error) {
		return nil, fmt.Errorf("err")
	}
	_, err = BitmapMarshal(roaring.BitmapOf(1))
	assert.Error(t, err)

	bitmap := roaring.New()
	err = BitmapUnmarshal(bitmap, data)
	assert.NoError(t, err)
	assert.EqualValues(t, roaring.BitmapOf(1).ToArray(), bitmap.ToArray())

	BitmapUnmarshal = func(bitmap *roaring.Bitmap, data []byte) error {
		return fmt.Errorf("err")
	}
	err = BitmapUnmarshal(bitmap, data)
	assert.Error(t, err)
}
