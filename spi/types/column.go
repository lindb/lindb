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

package types

import (
	"encoding/json"
	"math"
	"time"

	"github.com/lindb/common/models"
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/pkg/stream"
)

type Column struct {
	Values    []Value
	NumOfRows int
}

func NewColumn() *Column {
	return &Column{}
}

func (c *Column) Marshal(meta ColumnMetadata, w *stream.BufferWriter) {
	switch meta.DataType {
	case DTString:
		for _, v := range c.Values {
			w.PutString(v.(string))
		}
	case DTJSON:
		for _, v := range c.Values {
			data := []byte(v.(json.RawMessage))
			w.PutUvarint32(uint32(len(data)))
			w.PutBytes(data)
		}
	case DTInt:
		for _, v := range c.Values {
			w.PutVarint64(v.(int64))
		}
	case DTFloat:
		for _, v := range c.Values {
			w.PutUvarint64(math.Float64bits(v.(float64)))
		}
	case DTTimestamp:
		for _, v := range c.Values {
			t := v.(time.Time)
			w.PutVarint64(t.UnixMilli())
		}
	case DTDuration:
		for _, v := range c.Values {
			w.PutVarint64(int64(v.(time.Duration)))
		}
	case DTMap:
		for _, v := range c.Values {
			m := v.(map[string]string)
			w.PutUvarint32(uint32(len(m)))
			for key, value := range m {
				w.PutString(key)
				w.PutString(value)
			}
		}
	case DTTimeSeries:
		for _, v := range c.Values {
			if v == nil {
				w.PutUvarint32(uint32(0))
				continue
			}
			ts := v.(*TimeSeries)
			// TODO:need refactor
			data := encoding.JSONMarshal(ts)
			w.PutUvarint32(uint32(len(data)))
			w.PutBytes(data)
		}
	case DTExemplar:
		for _, v := range c.Values {
			if v == nil {
				w.PutUvarint32(uint32(0))
				continue
			}
			ts := v.([]*models.Exemplar)
			// TODO:need refactor
			data := encoding.JSONMarshal(ts)
			w.PutUvarint32(uint32(len(data)))
			w.PutBytes(data)
		}
	}
}

func (c *Column) Unmarshal(meta ColumnMetadata, numOfRows int, r *stream.Reader) {
	if numOfRows == 0 {
		return
	}
	c.NumOfRows = numOfRows
	c.Values = make([]Value, numOfRows)
	for i := range numOfRows {
		switch meta.DataType {
		case DTString:
			c.Values[i] = r.ReadString()
		case DTJSON:
			size := r.ReadUvarint32()
			data := r.ReadBytes(int(size))
			c.Values[i] = json.RawMessage(data)
		case DTInt:
			c.Values[i] = r.ReadVarint64()
		case DTFloat:
			bits := r.ReadUvarint64()
			c.Values[i] = math.Float64frombits(bits)
		case DTTimestamp:
			ms := r.ReadVarint64()
			c.Values[i] = time.UnixMilli(ms)
		case DTDuration:
			dur := r.ReadVarint64()
			c.Values[i] = time.Duration(dur)
		case DTMap:
			size := r.ReadUvarint32()
			m := make(map[string]string, size)
			for range size {
				key := r.ReadString()
				value := r.ReadString()
				m[key] = value
			}
			c.Values[i] = m
		case DTTimeSeries:
			size := r.ReadUvarint32()
			if size == 0 {
				continue
			}
			data := r.ReadBytes(int(size))
			ts := &TimeSeries{}
			err := json.Unmarshal(data, ts)
			if err != nil {
				panic(err)
			}
			c.Values[i] = ts
		case DTExemplar:
			size := r.ReadUvarint32()
			if size == 0 {
				continue
			}
			data := r.ReadBytes(int(size))
			var ts []*models.Exemplar
			err := json.Unmarshal(data, &ts)
			if err != nil {
				panic(err)
			}
			c.Values[i] = ts
		}
	}
}

func (c *Column) Append(val any) {
	c.Values = append(c.Values, val)
	c.NumOfRows++
}

func (c *Column) Reset(row int, val any) {
	c.Values[row] = val
}

func (c *Column) GetString(row int) string {
	if row >= len(c.Values) {
		return ""
	}
	val := c.Values[row]
	if val == nil {
		return ""
	}
	return val.(string)
}

func (c *Column) GetJSON(row int) json.RawMessage {
	if row >= len(c.Values) {
		return nil
	}
	return c.Values[row].(json.RawMessage)
}

func (c *Column) GetInt(row int) int64 {
	if row >= len(c.Values) {
		return 0
	}
	// FIXME:
	return c.Values[row].(int64)
}

func (c *Column) GetFloat(row int) float64 {
	if row >= len(c.Values) {
		return 0
	}
	// FIXME:
	return c.Values[row].(float64)
}

func (c *Column) GetTimestamp(row int) time.Time {
	if row >= len(c.Values) {
		return time.Time{}
	}
	return c.Values[row].(time.Time)
}

func (c *Column) GetDuration(row int) time.Duration {
	if row >= len(c.Values) {
		return time.Duration(0)
	}
	return c.Values[row].(time.Duration)
}

func (c *Column) GetMap(row int) map[string]string {
	if row >= len(c.Values) {
		return nil
	}
	return c.Values[row].(map[string]string)
}

func (c *Column) GetTimeSeries(row int) *TimeSeries {
	if row >= len(c.Values) {
		return nil
	}
	val := c.Values[row]
	if val == nil {
		return nil
	}

	// FIXME:
	return c.Values[row].(*TimeSeries)
}

func (c *Column) Get(row int) any {
	if row >= len(c.Values) {
		return nil
	}
	// FIXME:
	return c.Values[row]
}
