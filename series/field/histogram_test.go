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

package field

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HistogramConverter(t *testing.T) {
	v := math.MaxFloat64 + 1
	assert.Equal(t, "__bucket_+Inf", HistogramConverter.BucketName(v))
	assert.Equal(t, "__bucket_0.025", HistogramConverter.BucketName(0.025))
	assert.Equal(t, "__bucket_0.5", HistogramConverter.BucketName(0.500))
	assert.Equal(t, "__bucket_5000", HistogramConverter.BucketName(5000))

	f, err := HistogramConverter.UpperBound("__bucket_5000")
	assert.Nil(t, err)
	assert.Equal(t, float64(5000), f)

	f, err = HistogramConverter.UpperBound("__bucket_0.025")
	assert.Nil(t, err)
	assert.Equal(t, 0.025, f)

	_, err = HistogramConverter.UpperBound("_bucket_0.025")
	assert.NotNil(t, err)

	_, err = HistogramConverter.UpperBound("__bucket_x")
	assert.NotNil(t, err)

	assert.Equal(t, "_Histogram", HistogramConverter.Sanitize("_Histogram"))
	assert.Equal(t, "_Histogram", HistogramConverter.Sanitize("Histogram"))
	assert.Equal(t, "_HistogramCount", HistogramConverter.Sanitize("HistogramCount"))
	assert.Equal(t, "_bucket_32", HistogramConverter.Sanitize("__bucket_32"))

}
