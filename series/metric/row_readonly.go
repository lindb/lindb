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
	"bytes"
	"strings"

	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
	"github.com/lindb/lindb/series/field"
)

// readOnlyRow is an embedded struct used by StorageRow and BrokerRow
type readOnlyRow struct {
	m flatMetricsV1.Metric

	// lazy initialization
	keyValueIterator      readOnlyKeyValueIterator
	simpleFieldIterator   readOnlySimpleFieldIterator
	compoundFieldIterator readOnlyCompoundFieldIterator
}

// ShouldSanitizeName checks if metric-name is in necessary of sanitizing
func (mr *readOnlyRow) ShouldSanitizeName() bool {
	return bytes.IndexByte(mr.Name(), '|') >= 0
}

// ShouldSanitizeNameSpace checks if namespace is in necessary of sanitizing
func (mr *readOnlyRow) ShouldSanitizeNameSpace() bool {
	return bytes.IndexByte(mr.NameSpace(), '|') >= 0
}

func (mr *readOnlyRow) SanitizedName() string {
	return strings.Replace(strutil.ByteSlice2String(mr.Name()), "|", "_", -1)
}

func (mr *readOnlyRow) SanitizedNamespace() string {
	return strings.Replace(strutil.ByteSlice2String(mr.NameSpace()), "|", "_", -1)
}

func (mr *readOnlyRow) Timestamp() int64  { return mr.m.Timestamp() }
func (mr *readOnlyRow) Name() []byte      { return mr.m.Name() }
func (mr *readOnlyRow) NameSpace() []byte { return mr.m.Namespace() }
func (mr *readOnlyRow) TagsHash() uint64  { return mr.m.Hash() }

func (mr *readOnlyRow) NewKeyValueIterator() *readOnlyKeyValueIterator {
	mr.keyValueIterator.idx = -1
	mr.keyValueIterator.m = &mr.m
	mr.keyValueIterator.num = mr.m.KeyValuesLength()
	return &mr.keyValueIterator
}
func (mr *readOnlyRow) NewSimpleFieldIterator() *readOnlySimpleFieldIterator {
	mr.simpleFieldIterator.idx = -1
	mr.simpleFieldIterator.m = &mr.m
	mr.simpleFieldIterator.num = mr.m.SimpleFieldsLength()
	return &mr.simpleFieldIterator
}
func (mr *readOnlyRow) NewCompoundFieldIterator() (*readOnlyCompoundFieldIterator, bool) {
	mr.compoundFieldIterator.idx = -1
	mr.compoundFieldIterator.m = &mr.m

	if obj := mr.m.CompoundField(&mr.compoundFieldIterator.f); obj == nil {
		return nil, false
	}
	mr.compoundFieldIterator.num = mr.compoundFieldIterator.f.ExplicitBoundsLength()
	return &mr.compoundFieldIterator, true
}

type readOnlyKeyValueIterator struct {
	m   *flatMetricsV1.Metric
	kv  flatMetricsV1.KeyValue
	idx int
	num int
}

func (itr *readOnlyKeyValueIterator) HasNext() bool {
	itr.idx++
	if itr.idx >= itr.num {
		return false
	}
	return itr.m.KeyValues(&itr.kv, itr.idx)
}
func (itr *readOnlyKeyValueIterator) Len() int          { return itr.num }
func (itr *readOnlyKeyValueIterator) NextKey() []byte   { return itr.kv.Key() }
func (itr *readOnlyKeyValueIterator) NextValue() []byte { return itr.kv.Value() }
func (itr *readOnlyKeyValueIterator) Reset()            { itr.idx = -1 }

type readOnlySimpleFieldIterator struct {
	m   *flatMetricsV1.Metric
	f   flatMetricsV1.SimpleField
	idx int
	num int
}

func (itr *readOnlySimpleFieldIterator) HasNext() bool {
	itr.idx++
	if !(itr.idx < itr.num) {
		return false
	}
	return itr.m.SimpleFields(&itr.f, itr.idx)
}

// Reset iterator for re-iterating simpleFields
func (itr *readOnlySimpleFieldIterator) Reset()             { itr.idx = -1 }
func (itr *readOnlySimpleFieldIterator) Len() int           { return itr.num }
func (itr *readOnlySimpleFieldIterator) NextName() []byte   { return itr.f.Name() }
func (itr *readOnlySimpleFieldIterator) NextValue() float64 { return itr.f.Value() }
func (itr *readOnlySimpleFieldIterator) NextType() field.Type {
	switch itr.f.Type() {
	// assertion: cumulative should be converted before writing into memdb
	case flatMetricsV1.SimpleFieldTypeDeltaSum:
		return field.SumField
	case flatMetricsV1.SimpleFieldTypeGauge:
		return field.GaugeField
	case flatMetricsV1.SimpleFieldTypeMax:
		return field.MaxField
	case flatMetricsV1.SimpleFieldTypeMin:
		return field.MinField
	default:
		return field.Unknown
	}
}

func (itr *readOnlySimpleFieldIterator) ShouldSanitizeNextName() bool {
	v := itr.NextName()
	// internal histogram field
	return bytes.HasPrefix(v, []byte("Histogram")) ||
		bytes.HasPrefix(v, []byte("__bucket_")) // bucket field
}

// SanitizeNextName escapes the illegal field name,
// if reserved field-name is used, the input will be escaped with underline.
// HistogramSum-> _HistogramSum
// __bucket_ -> _bucket_
func (itr *readOnlySimpleFieldIterator) SanitizeNextName() string {
	v := itr.NextName()
	switch {
	case bytes.HasPrefix(v, []byte("Histogram")):
		return "_" + string(v)
	case bytes.HasPrefix(v, []byte("__bucket_")):
		return string(v[1:])
	default:
		return string(v)
	}
}

type readOnlyCompoundFieldIterator struct {
	m   *flatMetricsV1.Metric
	f   flatMetricsV1.CompoundField
	idx int
	num int
}

func (itr *readOnlyCompoundFieldIterator) HasNextBucket() bool {
	itr.idx++
	return itr.idx < itr.num
}
func (itr *readOnlyCompoundFieldIterator) NextExplicitBound() float64 {
	return itr.f.ExplicitBounds(itr.idx)
}
func (itr *readOnlyCompoundFieldIterator) BucketLen() int     { return itr.num }
func (itr *readOnlyCompoundFieldIterator) NextValue() float64 { return itr.f.Values(itr.idx) }
func (itr *readOnlyCompoundFieldIterator) Reset()             { itr.idx = -1 }
func (itr *readOnlyCompoundFieldIterator) Min() float64       { return itr.f.Min() }
func (itr *readOnlyCompoundFieldIterator) Max() float64       { return itr.f.Max() }
func (itr *readOnlyCompoundFieldIterator) Sum() float64       { return itr.f.Sum() }
func (itr *readOnlyCompoundFieldIterator) Count() float64     { return itr.f.Count() }
func (itr *readOnlyCompoundFieldIterator) BucketName() string {
	return HistogramConverter.BucketName(itr.NextExplicitBound())
}
