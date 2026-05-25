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
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lindb/common/constants"
	"github.com/lindb/common/field"
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"

	sfield "github.com/lindb/lindb/series/field"
)

var defaultNS = []byte(constants.DefaultNamespace)

// readOnlyRow is an embedded struct used by StorageRow and BrokerRow
type readOnlyRow struct {
	m flatMetricsV1.Metric

	// lazy initialization
	keyValueIterator      KeyValueIterator
	simpleFieldIterator   SimpleFieldIterator
	exemplarIterator      ExemplarIterator
	compoundFieldIterator CompoundFieldIterator
}

func (mr *readOnlyRow) Timestamp() int64 { return mr.m.Timestamp() }

func (mr *readOnlyRow) Name() []byte { return mr.m.Name() }

func (mr *readOnlyRow) NameSpace() []byte {
	ns := mr.m.Namespace()
	if len(ns) == 0 {
		return defaultNS
	}
	return ns
}

func (mr *readOnlyRow) NameHash() uint64 { return mr.m.NameHash() }

func (mr *readOnlyRow) TagsHash() uint64 { return mr.m.KvsHash() }

func (mr *readOnlyRow) TagsLen() int { return mr.m.KeyValuesLength() }

func (mr *readOnlyRow) SimpleFieldsLen() int { return mr.m.SimpleFieldsLength() }

func (mr *readOnlyRow) ExemplarsLen() int { return mr.m.ExemplarsLength() }

func (mr *readOnlyRow) NewKeyValueIterator() *KeyValueIterator {
	mr.keyValueIterator.idx = -1
	mr.keyValueIterator.m = &mr.m
	mr.keyValueIterator.num = mr.m.KeyValuesLength()
	return &mr.keyValueIterator
}

func (mr *readOnlyRow) NewSimpleFieldIterator() *SimpleFieldIterator {
	mr.simpleFieldIterator.idx = -1
	mr.simpleFieldIterator.m = &mr.m
	mr.simpleFieldIterator.num = mr.m.SimpleFieldsLength()
	return &mr.simpleFieldIterator
}

func (mr *readOnlyRow) NewExemplarIterator() *ExemplarIterator {
	mr.exemplarIterator.idx = -1
	mr.exemplarIterator.m = &mr.m
	mr.exemplarIterator.num = mr.m.ExemplarsLength()
	return &mr.exemplarIterator
}

func (mr *readOnlyRow) NewCompoundFieldIterator() (*CompoundFieldIterator, bool) {
	mr.compoundFieldIterator.idx = -1
	mr.compoundFieldIterator.m = &mr.m

	if obj := mr.m.CompoundField(&mr.compoundFieldIterator.f); obj == nil {
		return nil, false
	}
	mr.compoundFieldIterator.num = min(mr.compoundFieldIterator.f.ExplicitBoundsLength(),
		mr.compoundFieldIterator.f.ValuesLength())
	return &mr.compoundFieldIterator, true
}

type TimeSeriesIterator struct{}

type KeyValueIterator struct {
	m   *flatMetricsV1.Metric
	kv  flatMetricsV1.KeyValue
	idx int
	num int
}

func (itr *KeyValueIterator) HasNext() bool {
	itr.idx++
	if itr.idx >= itr.num {
		return false
	}
	return itr.m.KeyValues(&itr.kv, itr.idx)
}

func (itr *KeyValueIterator) Len() int { return itr.num }

func (itr *KeyValueIterator) NextKey() []byte { return itr.kv.Key() }

func (itr *KeyValueIterator) NextValue() []byte { return itr.kv.Value() }

func (itr *KeyValueIterator) Reset() { itr.idx = -1 }

type SimpleFieldIterator struct {
	m   *flatMetricsV1.Metric
	f   flatMetricsV1.SimpleField
	idx int
	num int
}

func (itr *SimpleFieldIterator) HasNext() bool {
	itr.idx++
	if itr.idx >= itr.num {
		return false
	}
	return itr.m.SimpleFields(&itr.f, itr.idx)
}

// Reset iterator for re-iterating simpleFields
func (itr *SimpleFieldIterator) Reset() { itr.idx = -1 }

func (itr *SimpleFieldIterator) Len() int { return itr.num }

func (itr *SimpleFieldIterator) NextName() sfield.Name {
	return sfield.Name(itr.f.Name())
}

func (itr *SimpleFieldIterator) NextRawName() []byte { return itr.f.Name() }

func (itr *SimpleFieldIterator) NextValue() float64 { return itr.f.Value() }

func (itr *SimpleFieldIterator) NextRawType() flatMetricsV1.SimpleFieldType { return itr.f.Type() }

func (itr *SimpleFieldIterator) NextType() field.Type {
	switch itr.f.Type() {
	// assertion: cumulative should be converted before writing into memdb
	case flatMetricsV1.SimpleFieldTypeDeltaSum:
		return field.Sum
	case flatMetricsV1.SimpleFieldTypeLast:
		return field.Last
	case flatMetricsV1.SimpleFieldTypeMax:
		return field.Max
	case flatMetricsV1.SimpleFieldTypeMin:
		return field.Min
	case flatMetricsV1.SimpleFieldTypeFirst:
		return field.First
	default:
		return field.Unknown
	}
}

// ExemplarIterator iterates exemplar fields
type ExemplarIterator struct {
	m   *flatMetricsV1.Metric
	f   flatMetricsV1.Exemplar
	idx int
	num int
}

func (itr *ExemplarIterator) HasNext() bool {
	itr.idx++
	if itr.idx >= itr.num {
		return false
	}
	return itr.m.Exemplars(&itr.f, itr.idx)
}

// Reset iterator for re-iterating exemplarFields
func (itr *ExemplarIterator) Reset() { itr.idx = -1 }

func (itr *ExemplarIterator) Len() int { return itr.num }

func (itr *ExemplarIterator) NextName() sfield.Name { return sfield.Name(itr.f.Name()) }

func (itr *ExemplarIterator) NextRawName() []byte { return itr.f.Name() }

func (itr *ExemplarIterator) NextTraceID() []byte { return itr.f.TraceId() }

func (itr *ExemplarIterator) NextSpanID() []byte { return itr.f.SpanId() }

func (itr *ExemplarIterator) NextDuration() int64 { return itr.f.Duration() }

type CompoundFieldIterator struct {
	m   *flatMetricsV1.Metric
	f   flatMetricsV1.CompoundField
	idx int
	num int
}

func (itr *CompoundFieldIterator) HasNextBucket() bool {
	itr.idx++
	return itr.idx < itr.num
}

func (itr *CompoundFieldIterator) NextExplicitBound() float64 {
	return itr.f.ExplicitBounds(itr.idx)
}

func (itr *CompoundFieldIterator) BucketLen() int { return itr.num }

func (itr *CompoundFieldIterator) NextValue() float64 { return itr.f.Values(itr.idx) }

func (itr *CompoundFieldIterator) Reset() { itr.idx = -1 }

func (itr *CompoundFieldIterator) Min() float64 { return itr.f.Min() }

func (itr *CompoundFieldIterator) Max() float64 { return itr.f.Max() }

func (itr *CompoundFieldIterator) Sum() float64 { return itr.f.Sum() }

func (itr *CompoundFieldIterator) Count() float64 { return itr.f.Count() }

func (itr *CompoundFieldIterator) BucketName() sfield.Name {
	return sfield.Name(BucketNameOfHistogramExplicitBound(itr.NextExplicitBound()))
}

const (
	histogramSum   = sfield.Name("HistogramSum")
	histogramCount = sfield.Name("HistogramCount")
	histogramMax   = sfield.Name("HistogramMax")
	histogramMin   = sfield.Name("HistogramMin")
)

func (itr *CompoundFieldIterator) HistogramSumName() sfield.Name   { return histogramSum }
func (itr *CompoundFieldIterator) HistogramCountFieldName() sfield.Name { return histogramCount }
func (itr *CompoundFieldIterator) HistogramMaxName() sfield.Name   { return histogramMax }
func (itr *CompoundFieldIterator) HistogramMinName() sfield.Name   { return histogramMin }

// BucketNameOfHistogramExplicitBound converts reserved field-name for histogram buckets.
func BucketNameOfHistogramExplicitBound(upperBound float64) string {
	if math.IsInf(upperBound, 1) {
		return HistoBucketSubField + "+Inf"
	}
	return HistoBucketSubField + strconv.FormatFloat(upperBound, 'f', -1, 32)
}

// UpperBound extracts the upper-bound from bucketName
func UpperBound(bucketName string) (float64, error) {
	// make sure it has prefix with __bucket_
	if !strings.HasPrefix(bucketName, HistoBucketSubField) {
		return 0, fmt.Errorf("bucketName:%s not startswith '%s'", bucketName, HistoBucketSubField)
	}
	raw := bucketName[len(HistoBucketSubField):]
	return strconv.ParseFloat(raw, 64)
}
