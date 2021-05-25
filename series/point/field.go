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

package point

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/lindb/lindb/pkg/escape"
	"github.com/lindb/lindb/series/field"
)

const fieldSuffixLength = 4

var (
	SumFieldSuffix       = []byte("_SUM")
	MinFieldSuffix       = []byte("_MIN")
	MaxFieldSuffix       = []byte("_MAX")
	SummaryFieldSuffix   = []byte("_SMY")
	HistogramFieldSuffix = []byte("_HGM")
)

// FieldIterator is used to traverse the fields of a point
// without constructing the in-memory map.
type FieldIterator struct {
	data        []byte // fields data
	start, end  int
	key, keybuf []byte
	valueBuf    []byte
	fieldType   field.Type
}

func (fi *FieldIterator) Next() bool {
	fi.start = fi.end
	if fi.start >= len(fi.data) {
		return false
	}

	fi.end, fi.key = scanTo(fi.data, fi.start, '=')
	if escape.IsEscaped(fi.key) {
		fi.keybuf = escape.AppendUnescaped(fi.keybuf[:0], fi.key)
		fi.key = fi.keybuf
	}

	fi.end, fi.valueBuf = scanFieldValue(fi.data, fi.end+1)
	fi.end++

	switch {
	case (len(fi.key) <= fieldSuffixLength) || (len(fi.valueBuf) == 0):
		fi.fieldType = field.Unknown
	case bytes.HasSuffix(fi.key, SumFieldSuffix):
		fi.fieldType = field.SumField
	case bytes.HasSuffix(fi.key, MinFieldSuffix):
		fi.fieldType = field.MinField
	case bytes.HasSuffix(fi.key, MaxFieldSuffix):
		fi.fieldType = field.MaxField
	case bytes.HasSuffix(fi.key, SummaryFieldSuffix):
		fi.fieldType = field.SummaryField
	case bytes.HasSuffix(fi.key, HistogramFieldSuffix):
		fi.fieldType = field.HistogramField
	}
	return true
}

func (fi *FieldIterator) Name() []byte {
	if len(fi.key) <= fieldSuffixLength {
		return fi.key
	}
	return fi.key[:len(fi.key)-fieldSuffixLength]
}

func (fi *FieldIterator) Type() field.Type {
	return fi.fieldType
}

func (fi *FieldIterator) Reset(data []byte) {
	fi.fieldType = field.Unknown
	fi.key = nil
	fi.valueBuf = nil
	fi.start = 0
	fi.end = 0
	fi.data = data
}

func (fi *FieldIterator) Int64Value() (int64, error) {
	n, err := parseInt64Bytes(fi.valueBuf)
	if err != nil {
		return 0, fmt.Errorf("parse int64 field value %q with error: %v", fi.valueBuf, err)
	}
	return n, nil
}

func (fi *FieldIterator) Float64Value() (float64, error) {
	f, err := parseFloat64Bytes(fi.valueBuf)
	if err != nil {
		return 0, fmt.Errorf("parse float64 field value %q with error: %v", fi.valueBuf, err)
	}
	return f, nil
}

func scanFieldValue(buf []byte, i int) (int, []byte) {
	start := i
	for i < len(buf) {
		if buf[i] == '\\' && i+1 < len(buf) && buf[i+1] == '\\' {
			i += 2
			continue
		}

		if buf[i] == ',' {
			break
		}
		i++
	}
	return i, buf[start:i]
}

// parseInt64Bytes is a zero-alloc wrapper around strconv.ParseInt.
func parseInt64Bytes(b []byte) (i int64, err error) {
	s := unsafeBytesToString(b)
	return strconv.ParseInt(s, 10, 64)
}

// parseFloat64Bytes is a zero-alloc wrapper around strconv.ParseFloat.
func parseFloat64Bytes(b []byte) (float64, error) {
	s := unsafeBytesToString(b)
	return strconv.ParseFloat(s, 64)
}

// parseUintBytes is a zero-alloc wrapper around strconv.ParseUint.
func parseUint64Bytes(b []byte) (i uint64, err error) {
	s := unsafeBytesToString(b)
	return strconv.ParseUint(s, 10, 64)
}

// unsafeBytesToString converts a []byte to a string without a heap allocation.
//
// It is unsafe, and is intended to prepare input to short-lived functions
// that require strings.
func unsafeBytesToString(in []byte) string {
	src := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	dst := reflect.StringHeader{
		Data: src.Data,
		Len:  src.Len,
	}
	s := *(*string)(unsafe.Pointer(&dst))
	return s
}

// MakeFields creates a key for a set of fields
func MakeFields(dst []byte, fs field.Fields) []byte {
	if len(fs) == 0 {
		return dst
	}
	for idx, f := range fs {
		var suffix []byte
		switch f.Type {
		case field.SumField:
			suffix = SumFieldSuffix
		case field.MinField:
			suffix = MinFieldSuffix
		case field.MaxField:
			suffix = MaxFieldSuffix
		case field.SummaryField:
			suffix = SummaryFieldSuffix
		case field.HistogramField:
			suffix = HistogramFieldSuffix
		default:
			continue
		}
		// int64 field is temporarily not existed
		var thisValue float64
		switch v := f.Value.(type) {
		case float64:
			thisValue = v
		case float32:
			thisValue = float64(v)
		case int:
			thisValue = float64(v)
		case int64:
			thisValue = float64(v)
		default:
			continue
		}
		dst = append(dst, escape.Bytes(f.Name)...)
		dst = append(dst, suffix...)
		dst = append(dst, '=')
		dst = strconv.AppendFloat(dst, thisValue, 'E', -1, 64)
		if idx != len(fs)-1 {
			dst = append(dst, ',')
		}
	}
	return dst
}
