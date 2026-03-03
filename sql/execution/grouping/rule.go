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

package grouping

import (
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/common/pkg/timeutil"
)

type Rule interface {
	Map(column arrow.Array, row int, buf *Buffer)
	Unmap(builder array.Builder, buf *Buffer)
}

type MapRule struct {
	mapper *StringMapper
}

func newMapRule(mapper *StringMapper) Rule {
	return &MapRule{mapper: mapper}
}

func (r *MapRule) Map(column arrow.Array, row int, buf *Buffer) {
	values, ok := column.(*larray.Map)
	if !ok || column.IsNull(row) {
		buf.Write(0)
		fmt.Println("kkk..", row)
		return
	}
	// FIXME: need sort keys of map
	rawRow := values.Row(row)
	keys := values.Keys().(*array.String)
	items := values.Items().(*array.String)
	offsets := values.Offsets()
	start, end := offsets[rawRow], offsets[rawRow+1]
	buf.Write(uint32(end - start))
	for i := int(start); i < int(end); i++ {
		buf.Write(r.mapper.GetID(keys.Value(i)))
		buf.Write(r.mapper.GetID(items.Value(i)))
	}
}

func (r *MapRule) Unmap(builder array.Builder, buf *Buffer) {
	count := buf.Read()
	if count == 0 {
		// for empty map, just append null value
		builder.AppendNull()
		return
	}
	mb := builder.(*array.MapBuilder)

	keysBuilder := mb.KeyBuilder().(*array.StringBuilder)
	valuesBuilder := mb.ItemBuilder().(*array.StringBuilder)

	mb.Append(true)
	mb.Reserve(int(count))

	for range count {
		keysBuilder.Append(r.mapper.GetValue(buf.Read()))
		valuesBuilder.Append(r.mapper.GetValue(buf.Read()))
	}
}

type StringRule struct {
	mapper *StringMapper
}

func newStringRule(mapper *StringMapper) Rule {
	return &StringRule{mapper: mapper}
}

func (r *StringRule) Map(column arrow.Array, row int, buf *Buffer) {
	value, ok := column.(*larray.Generic[string])
	if !ok || column.IsNull(row) {
		buf.Write(0)
		return
	}
	buf.Write(r.mapper.GetID(value.Value(row)))
}

func (r *StringRule) Unmap(builder array.Builder, buf *Buffer) {
	v := buf.Read()
	if v == 0 {
		builder.AppendNull()
		return
	}

	sb := builder.(*array.StringBuilder)
	sb.Append(r.mapper.GetValue(v))
}

type TimestampRule struct {
	mapper *StringMapper
}

func newTimestampRule(mapper *StringMapper) Rule {
	return &TimestampRule{
		mapper: mapper,
	}
}

func (r *TimestampRule) Map(column arrow.Array, row int, buf *Buffer) {
	value, ok := column.(*larray.Generic[arrow.Timestamp])
	if !ok || column.IsNull(row) {
		buf.Write(0)
		return
	}
	// TODO: check timetamp unit???
	ts := timeutil.FormatTimestamp(int64(value.Value(row)), timeutil.DataTimeFormat4)
	buf.Write(r.mapper.GetID(ts))
}

func (r *TimestampRule) Unmap(builder array.Builder, buf *Buffer) {
	val := buf.Read()
	if val == 0 {
		builder.AppendNull()
		return
	}
	tsStr := r.mapper.GetValue(val)
	t, err := time.ParseInLocation(timeutil.DataTimeFormat4, tsStr, time.Local)
	if err != nil {
		panic("parse timestamp string error:" + tsStr)
	}
	tb := builder.(*array.TimestampBuilder)
	tb.Append(arrow.Timestamp(t.UnixMilli()))
}
