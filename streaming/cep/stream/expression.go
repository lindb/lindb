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

package stream

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/sql/tree"
)

type Expr interface {
	Eval(record arrow.RecordBatch) (*roaring.Bitmap, error)
}

// colValue extracts a per-row string value from a record batch column.
// It handles plain string columns and map subscript access (col['key']).
// bind is called once per RecordBatch to resolve the column type;
// after binding, value can be called for each row to get the string value and null status.
type colValue interface {
	bind(record arrow.RecordBatch)
	value(row int) (string, bool)
}

// constCol returns a fixed string value for every row (literal constant).
type constCol struct {
	val string
}

func (c *constCol) bind(_ arrow.RecordBatch) {
}

func (c *constCol) value(_ int) (string, bool) {
	return c.val, true
}

// directCol reads a plain string column by index.
type directCol struct {
	colIndex int

	array arrow.Array
}

func (d *directCol) bind(record arrow.RecordBatch) {
	fmt.Println(record, d.colIndex)
	col := record.Column(d.colIndex)
	d.array = col
}

func (d *directCol) value(row int) (string, bool) {
	if d.array == nil || d.array.IsNull(row) {
		return "", false
	}
	switch c := d.array.(type) {
	case *larray.Generic[string]:
		return c.Value(row), true
	case *array.String:
		return c.Value(row), true
	default:
		return "", false
	}
}

// subscriptCol reads a value from a map column by key.
type subscriptCol struct {
	colIndex int
	key      string

	m       *larray.Map
	keys    *array.String
	values  *array.String
	offsets []int32
}

func (s *subscriptCol) bind(record arrow.RecordBatch) {
	col := record.Column(s.colIndex)

	m, ok := col.(*larray.Map)
	if !ok {
		s.m = nil
		return
	}
	s.m = m
	s.keys = m.Keys().(*array.String)
	s.values = m.Items().(*array.String)
	s.offsets = m.Offsets()
}

func (s *subscriptCol) value(row int) (string, bool) {
	if s.m == nil || s.m.IsNull(row) {
		return "", false
	}
	actualRow := s.m.Row(row)
	start, end := s.offsets[actualRow], s.offsets[actualRow+1]
	for j := int(start); j < int(end); j++ {
		if s.keys.Value(j) == s.key {
			if s.values.IsNull(j) {
				return "", false
			}
			return s.values.Value(j), true
		}
	}
	return "", false
}

type ComparisonExpr struct {
	left  colValue
	right colValue
	// TODO: add operator
}

func (e *ComparisonExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	result := roaring.New()
	e.left.bind(record)
	e.right.bind(record)
	n := int(record.NumRows())
	for i := range n {
		lv, lok := e.left.value(i)
		rv, rok := e.right.value(i)
		if lok && rok && lv == rv {
			result.Add(uint32(i))
		}
	}
	return result, nil
}

type InExpr struct {
	col    colValue
	values []string
}

func (e *InExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	result := roaring.New()
	e.col.bind(record)
	n := int(record.NumRows())
	for i := range n {
		v, ok := e.col.value(i)
		if !ok {
			continue
		}
		if slices.Contains(e.values, v) {
			result.Add(uint32(i))
		}
	}
	return result, nil
}

type NotExpr struct {
	expr Expr
}

func (e *NotExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	inner, err := e.expr.Eval(record)
	if err != nil {
		return nil, err
	}
	n := uint32(record.NumRows())
	all := roaring.New()
	all.AddRange(0, uint64(n))
	all.AndNot(inner)
	return all, nil
}

type LogicalExpr struct {
	op    tree.LogicalOperator
	exprs []Expr
}

func (e *LogicalExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	if len(e.exprs) == 0 {
		return roaring.New(), nil
	}
	result, err := e.exprs[0].Eval(record)
	if err != nil {
		return nil, err
	}
	for _, sub := range e.exprs[1:] {
		r, err := sub.Eval(record)
		if err != nil {
			return nil, err
		}
		switch e.op {
		case tree.LogicalAND:
			result.And(r)
		case tree.LogicalOR:
			result.Or(r)
		}
	}
	return result, nil
}

func likeToRegexp(pattern string) *regexp.Regexp {
	var sb strings.Builder
	sb.WriteByte('^')
	for _, ch := range pattern {
		switch ch {
		case '%':
			sb.WriteString(".*")
		case '_':
			sb.WriteByte('.')
		default:
			sb.WriteString(regexp.QuoteMeta(string(ch)))
		}
	}
	sb.WriteByte('$')
	return regexp.MustCompile(sb.String())
}

// RegexExpr matches rows where the column value matches a regular expression.
type RegexExpr struct {
	col     colValue
	pattern *regexp.Regexp
}

func (e *RegexExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	result := roaring.New()
	e.col.bind(record)
	n := int(record.NumRows())
	for i := range n {
		v, ok := e.col.value(i)
		if ok && e.pattern.MatchString(v) {
			result.Add(uint32(i))
		}
	}
	return result, nil
}

// NullExpr matches rows where the column value is NULL (not == "") when not is false,
// or IS NOT NULL when not is true.
type NullExpr struct {
	col colValue
	not bool
}

func (e *NullExpr) Eval(record arrow.RecordBatch) (*roaring.Bitmap, error) {
	result := roaring.New()
	e.col.bind(record)
	n := int(record.NumRows())
	for i := range n {
		_, ok := e.col.value(i)
		isNull := !ok
		if isNull != e.not {
			result.Add(uint32(i))
		}
	}
	return result, nil
}
