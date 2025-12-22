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
	"github.com/lindb/lindb/spi/types"
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

func (m *MapRule) Map(value any, buf *Buffer) {
	values, ok := value.(map[string]string)
	if !ok {
		return
	}
	// TODO: need sort keys of map
	buf.Write(uint32(len(values)))
	for k, v := range values {
		buf.Write(m.mapper.GetID(k))
		buf.Write(m.mapper.GetID(v))
	}
}

func (m *MapRule) Unmap(buf *Buffer) any {
	count := buf.Read()
	if count == 0 {
		return nil
	}
	values := make(map[string]string, count)
	for range count {
		values[m.mapper.GetValue(buf.Read())] = m.mapper.GetValue(buf.Read())
	}

	return values
}

type StringRule struct {
	mapper *StringMapper
}

func newStringRule(mapper *StringMapper) Rule {
	return &StringRule{mapper: mapper}
}

func (s *StringRule) Map(value any, buf *Buffer) {
	switch t := value.(type) {
	case string:
		buf.Write(s.mapper.GetID(t))
	case *types.String:
		buf.Write(s.mapper.GetID(string(*t)))
	}
}

func (s *StringRule) Unmap(buf *Buffer) any {
	return s.mapper.GetValue(buf.Read())
}
