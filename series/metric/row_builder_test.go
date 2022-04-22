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

package metric

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
)

func Test_NewRowBuilder(t *testing.T) {
	var lastData []byte
	for i := 0; i < 20; i++ {
		rb, releaseFunc := NewRowBuilder()

		assert.NoError(t, rb.AddTag([]byte("a"), []byte("b")))
		assert.NoError(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, 1))
		rb.AddMetricName([]byte("namespace"))
		rb.AddTimestamp(111111)
		thisData, err := rb.Build()
		if i > 0 {
			assert.Equal(t, lastData, thisData)
		}
		lastData = append(lastData[:0], thisData...)
		assert.NoError(t, err)
		releaseFunc(rb)
	}
}

func Test_RowBuilder_ErrorCases(t *testing.T) {
	rb := newRowBuilder()
	// tags validation
	assert.Error(t, rb.AddTag(nil, nil))
	assert.Error(t, rb.AddTag([]byte("tag-key"), nil))

	// simple field validation
	assert.Error(t, rb.AddSimpleField([]byte(""), flatMetricsV1.SimpleFieldTypeDeltaSum, 1))
	assert.Error(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeUnSpecified, 1))
	assert.Error(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, math.Inf(1)))
	assert.Error(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, math.NaN()))
	assert.Zero(t, rb.SimpleFieldsLen())

	// compound field validation
	assert.Error(t, rb.AddCompoundFieldData([]float64{1, 2}, []float64{1}))
	assert.Error(t, rb.AddCompoundFieldData([]float64{1}, []float64{1}))
	// not increasing
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 3, math.Inf(1)},
	))
	// last bound not +Inf
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, 6},
	))
	// first bound < 0
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{-1, 2, 3, 4, 5, math.Inf(1)},
	))
	// value contains Inf
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{math.Inf(1), 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
	// values contains negative float
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{-1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
	// values contains NaN
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{math.NaN(), 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))

	assert.NoError(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
	// mmsc
	assert.Error(t, rb.AddCompoundFieldMMSC(-1, -1, 0, 0))
}

func Test_RowBuilder_BuildError(t *testing.T) {
	rb := newRowBuilder()
	_, err := rb.Build()
	assert.Error(t, err)

	// fields empty
	rb.AddMetricName([]byte("ab"))
	rb.AddMetricName([]byte("a|b"))
	rb.AddNameSpace([]byte("a|b"))
	assert.Equal(t, []byte("a_b"), rb.nameSpace)
	assert.Equal(t, []byte("a_b"), rb.metricName)
	_, err = rb.Build()
	assert.Error(t, err)

	// kvs count too much
	for i := 0; i < 100; i++ {
		_ = rb.AddTag([]byte{uint8(i)}, []byte{2})
	}
	_ = rb.AddSimpleField([]byte("__bucket_1"), flatMetricsV1.SimpleFieldTypeDeltaSum, 1)
	_, err = rb.Build()
	assert.Error(t, err)
	// build to with error
	var row BrokerRow
	assert.Error(t, rb.BuildTo(&row))
}

func Test_RowBuilder_OneSimpleField(t *testing.T) {
	rb := newRowBuilder()
	rb.AddMetricName([]byte("cpu"))
	_ = rb.AddSimpleField([]byte("idle"), flatMetricsV1.SimpleFieldTypeGauge, 1)
	var row BrokerRow
	assert.NoError(t, rb.BuildTo(&row))
	assert.NotZero(t, row.m.Timestamp())
	assert.Equal(t, emptyStringHash, row.m.Hash())
	assert.Equal(t, "cpu", string(row.m.Name()))
}

func Test_RowBuilder_BuildTo(t *testing.T) {
	rb := newRowBuilder()
	assert.NoError(t, rb.AddTag([]byte("ip"), []byte("1.1.1.1")))
	assert.NoError(t, rb.AddTag([]byte("host"), []byte("dev-ecs")))
	rb.AddMetricName([]byte("cpu|load"))
	assert.NoError(t, rb.AddSimpleField([]byte("idle"), flatMetricsV1.SimpleFieldTypeGauge, 1))
	assert.NoError(t, rb.AddCompoundFieldMMSC(1, 1, 1, 1))

	assert.NoError(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))

	var row BrokerRow
	assert.NoError(t, rb.BuildTo(&row))

	assert.Equal(t, 2, row.m.KeyValuesLength())
	assert.Equal(t, "cpu_load", string(row.m.Name()))
	assert.Equal(t, 1, row.m.SimpleFieldsLength())

	var cf flatMetricsV1.CompoundField
	assert.NotNil(t, row.m.CompoundField(&cf))
	assert.InDelta(t, 1, cf.Count(), 1e-6)
	assert.InDelta(t, 1, cf.Sum(), 1e-6)
	assert.InDelta(t, 1, cf.Max(), 1e-6)
	assert.InDelta(t, 1, cf.Min(), 1e-6)

	assert.Equal(t, 6, cf.ValuesLength())
	assert.Equal(t, 6, cf.ExplicitBoundsLength())

	for i := 0; i < 6; i++ {
		assert.InDelta(t, float64(i+1), cf.Values(i), 1e-6)
	}
	assert.Equal(t, float64(1), cf.ExplicitBounds(0))
	assert.True(t, math.IsInf(cf.ExplicitBounds(5), 1))
}

func Test_dedupTagsThenXXHash(t *testing.T) {
	rb := newRowBuilder()
	_ = rb.AddTag([]byte("ccc"), []byte("a"))
	_ = rb.AddTag([]byte("d"), []byte("b"))
	_ = rb.AddTag([]byte("a"), []byte("c"))
	_ = rb.AddTag([]byte("ccc"), []byte("d"))
	_ = rb.AddTag([]byte("ccc"), []byte("e"))
	_ = rb.AddTag([]byte("a"), []byte("f"))
	_ = rb.AddTag([]byte("d"), []byte("g"))

	hash1 := rb.dedupTagsThenXXHash()
	assert.Equal(t, "a=f,ccc=e,d=g", rb.hashBuf.String())
	hash2 := rb.dedupTagsThenXXHash()
	assert.Equal(t, "a=f,ccc=e,d=g", rb.hashBuf.String())
	assert.Equal(t, hash2, hash1)
	assert.NotZero(t, hash2)
}

func Test_dedupTags_EmptyKVs(t *testing.T) {
	rb := newRowBuilder()
	hash1 := rb.dedupTagsThenXXHash()
	assert.Equal(t, "", rb.hashBuf.String())
	assert.Equal(t, hash1, emptyStringHash)
}

func Test_dedupTags_SortedKVs(t *testing.T) {
	rb := newRowBuilder()
	_ = rb.AddTag([]byte("a"), []byte("a"))
	_ = rb.AddTag([]byte("c"), []byte("c"))
	_ = rb.dedupTagsThenXXHash()
	assert.Equal(t, "a=a,c=c", rb.hashBuf.String())
}

func Test_dedupTagsThenXXHash_One(t *testing.T) {
	rb := newRowBuilder()
	_ = rb.AddTag([]byte("ccc"), []byte("a"))
	_ = rb.AddTag([]byte("ccc"), []byte("b"))
	_ = rb.AddTag([]byte("ccc"), []byte("c"))
	_ = rb.AddTag([]byte("ccc"), []byte("d"))
	_ = rb.AddTag([]byte("ccc"), []byte("e"))
	_ = rb.AddTag([]byte("ccc"), []byte("f"))
	_ = rb.AddTag([]byte("ccc"), []byte("g"))

	_ = rb.dedupTagsThenXXHash()
	assert.Equal(t, "ccc=g", rb.hashBuf.String())
}
