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
	"time"

	"github.com/lindb/common/pkg/timeutil"
)

type Rule interface {
	Map(value any, buf *Buffer)
	Unmap(buf *Buffer) any
}

type MapRule struct {
	mapper *StringMapper
}

func newMapRule(mapper *StringMapper) Rule {
	return &MapRule{mapper: mapper}
}

func (r *MapRule) Map(value any, buf *Buffer) {
	values, ok := value.(map[string]string)
	if !ok {
		return
	}
	// TODO: need sort keys of map
	buf.Write(uint32(len(values)))
	for k, v := range values {
		buf.Write(r.mapper.GetID(k))
		buf.Write(r.mapper.GetID(v))
	}
}

func (r *MapRule) Unmap(buf *Buffer) any {
	count := buf.Read()
	if count == 0 {
		return nil
	}
	values := make(map[string]string, count)
	for range count {
		values[r.mapper.GetValue(buf.Read())] = r.mapper.GetValue(buf.Read())
	}

	return values
}

type StringRule struct {
	mapper *StringMapper
}

func newStringRule(mapper *StringMapper) Rule {
	return &StringRule{mapper: mapper}
}

func (r *StringRule) Map(value any, buf *Buffer) {
	buf.Write(r.mapper.GetID(value.(string)))
}

func (r *StringRule) Unmap(buf *Buffer) any {
	return r.mapper.GetValue(buf.Read())
}

type TimestampRule struct {
	mapper *StringMapper
}

func newTimestampRule(mapper *StringMapper) Rule {
	return &TimestampRule{
		mapper: mapper,
	}
}

func (r *TimestampRule) Map(value any, buf *Buffer) {
	// OPT: need refactor time mapping logic
	switch t := value.(type) {
	case *time.Time:
		ts := timeutil.FormatTimestamp(t.UnixMilli(), timeutil.DataTimeFormat4)
		buf.Write(r.mapper.GetID(ts))
	case time.Time:
		ts := timeutil.FormatTimestamp(t.UnixMilli(), timeutil.DataTimeFormat4)
		buf.Write(r.mapper.GetID(ts))
	default:
		buf.Write(0)
	}
}

func (r *TimestampRule) Unmap(buf *Buffer) any {
	val := buf.Read()
	if val == 0 {
		return nil
	}
	tsStr := r.mapper.GetValue(val)
	t, err := time.ParseInLocation(timeutil.DataTimeFormat4, tsStr, time.Local)
	if err != nil {
		panic("parse timestamp string error:" + tsStr)
	}
	return t
}
