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
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/cespare/xxhash/v2"
	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
)

type rowKV struct {
	key   []byte
	value []byte
}

// rowKVs sorts key values, then computes the hash
type rowKVs struct {
	kvs     []rowKV
	kvCount int
}

func (items rowKVs) Len() int      { return items.kvCount }
func (items rowKVs) Swap(i, j int) { items.kvs[i], items.kvs[j] = items.kvs[j], items.kvs[i] }
func (items rowKVs) Less(i, j int) bool {
	return bytes.Compare(items.kvs[i].key, items.kvs[j].key) < 0
}

type rowSimpleField struct {
	name  []byte
	fType flatMetricsV1.SimpleFieldType
	value float64
}

// RowBuilder builds a flat metric in order.
type RowBuilder struct {
	// metric raw data
	metricName []byte
	nameSpace  []byte
	timestamp  int64

	rowKVs  rowKVs
	hashBuf bytes.Buffer // concat sorted kvs

	simpleFields     []rowSimpleField
	simpleFieldCount int

	compoundFieldValues         []float64
	compoundFieldExplicitValues []float64
	compoundFieldMin            float64
	compoundFieldMax            float64
	compoundFieldSum            float64
	compoundFieldCount          float64

	// context for building flat metrics
	flatBuilder *flatbuffers.Builder
	keys        []flatbuffers.UOffsetT
	values      []flatbuffers.UOffsetT
	kvs         []flatbuffers.UOffsetT
	fieldNames  []flatbuffers.UOffsetT
	fields      []flatbuffers.UOffsetT
}

var rowBuilderPool sync.Pool

// NewRowBuilder picks a row builder from pool for building flat metric
func NewRowBuilder() (
	rb *RowBuilder,
	releaseFunc func(rb *RowBuilder),
) {
	releaseFunc = func(rb *RowBuilder) { rowBuilderPool.Put(rb) }
	item := rowBuilderPool.Get()
	if item != nil {
		builder := item.(*RowBuilder)
		builder.Reset()
	}
	return &RowBuilder{
		flatBuilder: flatbuffers.NewBuilder(1536),
	}, releaseFunc
}

func newRowBuilder() *RowBuilder {
	return &RowBuilder{flatBuilder: flatbuffers.NewBuilder(1536)}
}

// AddTag appends a key-value pair
// Return false if tag is invalid
func (rb *RowBuilder) AddTag(key, value []byte) error {
	if len(key) == 0 || len(value) == 0 {
		return fmt.Errorf("tag[%s: %s] is empty", string(key), string(value))
	}
	rb.rowKVs.kvCount++

	if rb.rowKVs.kvCount > len(rb.rowKVs.kvs) {
		rb.rowKVs.kvs = append(rb.rowKVs.kvs, rowKV{})
	}
	kvIdx := rb.rowKVs.kvCount - 1
	// copy key
	rb.rowKVs.kvs[kvIdx].key = append(rb.rowKVs.kvs[kvIdx].key[:0], key...)
	// copy value
	rb.rowKVs.kvs[kvIdx].value = append(rb.rowKVs.kvs[kvIdx].value[:0], value...)
	return nil
}

// AddSimpleField appends a simple field
// Return false if field is invalid
func (rb *RowBuilder) AddSimpleField(fieldName []byte, fieldType flatMetricsV1.SimpleFieldType, fieldValue float64) error {
	if fieldType == flatMetricsV1.SimpleFieldTypeUnSpecified {
		return fmt.Errorf("flat field type is unspecified")
	}
	if math.IsInf(fieldValue, 0) {
		return fmt.Errorf("fieldValue is Inf :%f", fieldValue)
	}
	if math.IsNaN(fieldValue) {
		return fmt.Errorf("fieldValue is NaN :%f", fieldValue)
	}
	if len(fieldName) == 0 {
		return fmt.Errorf("fieldName is empty")
	}
	if ShouldSanitizeFieldName(fieldName) {
		fieldName = SanitizeFieldName(fieldName)
	}

	rb.simpleFieldCount++

	// add field name
	if rb.simpleFieldCount > len(rb.simpleFields) {
		rb.simpleFields = append(rb.simpleFields, rowSimpleField{})
	}
	sfIdx := rb.simpleFieldCount - 1
	// copy fieldName
	rb.simpleFields[sfIdx].name = append(rb.simpleFields[sfIdx].name[:0], fieldName...)
	// copy field type, field value
	rb.simpleFields[sfIdx].fType = fieldType
	rb.simpleFields[sfIdx].value = fieldValue
	return nil
}

func (rb *RowBuilder) AddTimestamp(ts int64) { rb.timestamp = ts }

func (rb *RowBuilder) AddCompoundFieldData(values, bounds []float64) error {
	if len(values) != len(bounds) {
		return fmt.Errorf("values's length: %d != explicit-bounds's length: %d",
			len(values), len(bounds),
		)
	}
	if len(values) < 2 {
		return fmt.Errorf("compound buckets: %d less than 2", len(values))
	}
	// ensure bounds increasing
	for idx := 1; idx < len(bounds); idx++ {
		if bounds[idx] < bounds[idx-1] {
			return fmt.Errorf("compound explicit bound is not increasing")
		}
	}
	// ensure last bound +Inf
	if !math.IsInf(bounds[len(bounds)-1], 1) {
		return fmt.Errorf("compound last explicit bound: %f is not +Inf", bounds[len(bounds)-1])
	}
	if bounds[0] < 0 {
		return fmt.Errorf("compound first explicit bound: %f < 0", bounds[0])
	}
	for _, v := range values {
		if math.IsInf(v, 0) {
			return fmt.Errorf("compound value contains Inf: %f", v)
		}
		if v < 0 {
			return fmt.Errorf("compound value less than zero: %f", v)
		}
		if math.IsNaN(v) {
			return fmt.Errorf("compound value contains NaN: %f", v)
		}
	}

	rb.compoundFieldValues = append(rb.compoundFieldValues[:0], values...)
	rb.compoundFieldExplicitValues = append(rb.compoundFieldExplicitValues[:0], bounds...)
	return nil
}

func (rb *RowBuilder) AddCompoundFieldMMSC(min, max, sum, count float64) error {
	rb.compoundFieldMin = min
	rb.compoundFieldMax = max
	rb.compoundFieldSum = sum
	rb.compoundFieldCount = count
	if !(min >= 0 && max >= 0 && sum >= 0 && count >= 0) {
		return fmt.Errorf("min: %f, max: %f, sum: %f, count: %f should >= 0",
			min, max, sum, count)
	}
	return nil
}

func (rb *RowBuilder) AddMetricName(metricName []byte) {
	if ShouldSanitizeNamespaceOrMetricName(metricName) {
		metricName = SanitizeNamespaceOrMetricName(metricName)
	}
	rb.metricName = append(rb.metricName[:0], metricName...)
}

var defaultNameSpace = []byte(constants.DefaultNamespace)

func (rb *RowBuilder) AddNameSpace(namespace []byte) {
	if ShouldSanitizeNamespaceOrMetricName(namespace) {
		namespace = SanitizeNamespaceOrMetricName(namespace)
	}
	rb.nameSpace = append(rb.nameSpace[:0], namespace...)
}

func (rb *RowBuilder) Reset() {
	rb.flatBuilder.Reset()
	rb.metricName = rb.metricName[:0]
	rb.nameSpace = rb.nameSpace[:0]
	rb.timestamp = 0

	// reset kvs context
	rb.rowKVs.kvCount = 0

	// reset simple fields context
	rb.simpleFieldCount = 0

	// reset compound context
	rb.compoundFieldValues = rb.compoundFieldValues[:0]
	rb.compoundFieldExplicitValues = rb.compoundFieldExplicitValues[:0]
	rb.compoundFieldMin = 0
	rb.compoundFieldMax = 0
	rb.compoundFieldSum = 0
	rb.compoundFieldCount = 0

	// reset flat builder context
	rb.flatBuilder.Reset()
	rb.keys = rb.keys[:0]
	rb.values = rb.values[:0]
	rb.kvs = rb.kvs[:0]
	rb.fieldNames = rb.fieldNames[:0]
	rb.fields = rb.fields[:0]
}

var (
	emptyStringHash = xxhash.Sum64String("")
)

func (rb *RowBuilder) _xxHashOfKVs() uint64 {
	if rb.rowKVs.kvCount == 0 {
		return emptyStringHash
	}
	rb.hashBuf.Reset()
	for idx := 0; idx < rb.rowKVs.kvCount; idx++ {
		if idx >= 1 {
			_ = rb.hashBuf.WriteByte(',')
		}
		_, _ = rb.hashBuf.Write(rb.rowKVs.kvs[idx].key)
		_ = rb.hashBuf.WriteByte('=')
		_, _ = rb.hashBuf.Write(rb.rowKVs.kvs[idx].value)
	}
	return xxhash.Sum64(rb.hashBuf.Bytes())
}

// dedupTags removes duplicated tags
func (rb *RowBuilder) dedupTagsThenXXHash() uint64 {
	if rb.rowKVs.kvCount < 2 {
		return rb._xxHashOfKVs()
	}
	if !sort.IsSorted(rb.rowKVs) {
		sort.Sort(rb.rowKVs)
	}
	// fast path
	shouldDeDup := false
	for cursor := 1; cursor < rb.rowKVs.kvCount; cursor++ {
		if bytes.Equal(rb.rowKVs.kvs[cursor].key, rb.rowKVs.kvs[cursor-1].key) {
			shouldDeDup = true
			break
		}
	}
	if !shouldDeDup {
		return rb._xxHashOfKVs()
	}

	// tags with same key will keep order as they are appended after sorting
	// high index key has higher priority
	// use 2-pointer algorithm
	var slow = 0
	for high := 1; high < rb.rowKVs.kvCount; high++ {
		if !bytes.Equal(rb.rowKVs.kvs[slow].key, rb.rowKVs.kvs[high].key) {
			slow++
		}
		rb.rowKVs.kvs[slow].value = append(rb.rowKVs.kvs[slow].value[:0], rb.rowKVs.kvs[high].value...)
		rb.rowKVs.kvs[slow].key = append(rb.rowKVs.kvs[slow].key[:0], rb.rowKVs.kvs[high].key...)
	}
	rb.rowKVs.kvCount = slow + 1
	return rb._xxHashOfKVs()
}

func (rb *RowBuilder) Build() ([]byte, error) {
	if len(rb.metricName) == 0 {
		return nil, fmt.Errorf("metric-name is empty")
	}
	if rb.simpleFieldCount == 0 && len(rb.compoundFieldValues) == 0 {
		return nil, fmt.Errorf("simple field and compound field are both empty")
	}
	if rb.rowKVs.kvCount > config.GlobalStorageConfig().TSDB.MaxTagKeysNumber {
		return nil, fmt.Errorf("too many tag pairs: %d", rb.rowKVs.kvCount)
	}
	hash := rb.dedupTagsThenXXHash()
	for i := 0; i < rb.rowKVs.kvCount; i++ {
		rb.keys = append(rb.keys, rb.flatBuilder.CreateByteString(rb.rowKVs.kvs[i].key))
		rb.values = append(rb.values, rb.flatBuilder.CreateByteString(rb.rowKVs.kvs[i].value))
	}
	// building key values vector
	for i := 0; i < len(rb.keys); i++ {
		flatMetricsV1.KeyValueStart(rb.flatBuilder)
		flatMetricsV1.KeyValueAddKey(rb.flatBuilder, rb.keys[i])
		flatMetricsV1.KeyValueAddValue(rb.flatBuilder, rb.values[i])
		rb.kvs = append(rb.kvs, flatMetricsV1.KeyValueEnd(rb.flatBuilder))
	}
	// building field names
	for i := 0; i < rb.simpleFieldCount; i++ {
		rb.fieldNames = append(rb.fieldNames, rb.flatBuilder.CreateByteString(rb.simpleFields[i].name))
	}

	for i := 0; i < rb.simpleFieldCount; i++ {
		flatMetricsV1.SimpleFieldStart(rb.flatBuilder)
		flatMetricsV1.SimpleFieldAddName(rb.flatBuilder, rb.fieldNames[i])
		flatMetricsV1.SimpleFieldAddType(rb.flatBuilder, rb.simpleFields[i].fType)
		flatMetricsV1.SimpleFieldAddValue(rb.flatBuilder, rb.simpleFields[i].value)
		rb.fields = append(rb.fields, flatMetricsV1.SimpleFieldEnd(rb.flatBuilder))
	}
	flatMetricsV1.MetricStartKeyValuesVector(rb.flatBuilder, rb.rowKVs.kvCount)
	for i := rb.rowKVs.kvCount - 1; i >= 0; i-- {
		rb.flatBuilder.PrependUOffsetT(rb.kvs[i])
	}
	kvs := rb.flatBuilder.EndVector(rb.rowKVs.kvCount)
	// serialize fields
	flatMetricsV1.MetricStartSimpleFieldsVector(rb.flatBuilder, rb.simpleFieldCount)
	for i := rb.simpleFieldCount - 1; i >= 0; i-- {
		rb.flatBuilder.PrependUOffsetT(rb.fields[i])
	}
	fields := rb.flatBuilder.EndVector(rb.simpleFieldCount)

	var (
		compoundFieldBounds flatbuffers.UOffsetT
		compoundFieldValues flatbuffers.UOffsetT
		compoundField       flatbuffers.UOffsetT
	)
	if len(rb.compoundFieldValues) == 0 {
		goto Serialize
	}
	// serialize compound fields
	// add compound buckets explicit bounds
	flatMetricsV1.CompoundFieldStartValuesVector(rb.flatBuilder, len(rb.compoundFieldValues))
	for i := len(rb.compoundFieldValues) - 1; i >= 0; i-- {
		rb.flatBuilder.PrependFloat64(rb.compoundFieldValues[i])
	}
	compoundFieldValues = rb.flatBuilder.EndVector(len(rb.compoundFieldValues))
	// add compound buckets values
	flatMetricsV1.CompoundFieldStartExplicitBoundsVector(rb.flatBuilder, len(rb.compoundFieldExplicitValues))
	for i := len(rb.compoundFieldExplicitValues) - 1; i >= 0; i-- {
		rb.flatBuilder.PrependFloat64(rb.compoundFieldExplicitValues[i])
	}
	compoundFieldBounds = rb.flatBuilder.EndVector(len(rb.compoundFieldExplicitValues))
	// add count sum min max
	flatMetricsV1.CompoundFieldStart(rb.flatBuilder)
	flatMetricsV1.CompoundFieldAddCount(rb.flatBuilder, rb.compoundFieldCount)
	flatMetricsV1.CompoundFieldAddSum(rb.flatBuilder, rb.compoundFieldSum)
	flatMetricsV1.CompoundFieldAddMin(rb.flatBuilder, rb.compoundFieldMin)
	flatMetricsV1.CompoundFieldAddMax(rb.flatBuilder, rb.compoundFieldMax)
	flatMetricsV1.CompoundFieldAddValues(rb.flatBuilder, compoundFieldValues)
	flatMetricsV1.CompoundFieldAddExplicitBounds(rb.flatBuilder, compoundFieldBounds)
	compoundField = flatMetricsV1.CompoundFieldEnd(rb.flatBuilder)

Serialize:
	metricName := rb.flatBuilder.CreateByteString(rb.metricName)
	if len(rb.nameSpace) == 0 {
		rb.nameSpace = defaultNameSpace
	}
	namespace := rb.flatBuilder.CreateByteString(rb.nameSpace)
	flatMetricsV1.MetricStart(rb.flatBuilder)
	flatMetricsV1.MetricAddNamespace(rb.flatBuilder, namespace)
	flatMetricsV1.MetricAddName(rb.flatBuilder, metricName)
	if rb.timestamp == 0 {
		rb.timestamp = fasttime.UnixMilliseconds()
	}
	flatMetricsV1.MetricAddTimestamp(rb.flatBuilder, rb.timestamp)
	flatMetricsV1.MetricAddKeyValues(rb.flatBuilder, kvs)
	flatMetricsV1.MetricAddHash(rb.flatBuilder, hash)
	flatMetricsV1.MetricAddSimpleFields(rb.flatBuilder, fields)
	if compoundField != 0 {
		flatMetricsV1.MetricAddCompoundField(rb.flatBuilder, compoundField)
	}
	end := flatMetricsV1.MetricEnd(rb.flatBuilder)
	// size prefix encoding
	rb.flatBuilder.FinishSizePrefixed(end)

	return rb.flatBuilder.FinishedBytes(), nil
}

func (rb *RowBuilder) BuildTo(row *BrokerRow) error {
	generatedBlock, err := rb.Build()
	if err != nil {
		return err
	}
	row.FromBlock(generatedBlock)
	return nil
}

func (rb *RowBuilder) SimpleFieldsLen() int { return rb.simpleFieldCount }
