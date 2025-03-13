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

package tree

import (
	"strconv"
	"strings"
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi/types"
)

type Literal interface {
	Node
}

type StringLiteral struct {
	BaseNode

	Value string `json:"value"`
}

func (n *StringLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type BooleanLiteral struct {
	BaseNode
	Value bool
}

func (n *BooleanLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

func NewBooleanLiteral(id NodeID, location *NodeLocation, value string) *BooleanLiteral {
	return &BooleanLiteral{
		BaseNode: BaseNode{
			ID:       id,
			Location: location,
		},
		Value: strings.EqualFold(value, "true"),
	}
}

type LongLiteral struct {
	BaseNode
	Value int64
}

func NewLongLiteral(id NodeID, location *NodeLocation, value string) *LongLiteral {
	// TODO: check error
	val, _ := strconv.ParseInt(value, 10, 64)
	return &LongLiteral{
		BaseNode: BaseNode{
			ID:       id,
			Location: location,
		},
		Value: val,
	}
}

func (n *LongLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type FloatLiteral struct {
	BaseNode
	Value float64
}

func NewFloatLiteral(id NodeID, location *NodeLocation, value string) *FloatLiteral {
	// TODO: parse value
	val, _ := strconv.ParseFloat(value, 64)
	return &FloatLiteral{
		BaseNode: BaseNode{
			ID:       id,
			Location: location,
		},
		Value: val,
	}
}

func (n *FloatLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type IntervalUnit string

const (
	Second IntervalUnit = "SECOND"
	Minute IntervalUnit = "MINUTE"
	Hour   IntervalUnit = "HOUR"
	Day    IntervalUnit = "DAY"
	Month  IntervalUnit = "MONTH"
	Year   IntervalUnit = "YEAR"
)

type IntervalLiteral struct {
	BaseNode
	Unit  IntervalUnit
	Value time.Duration
}

func NewIntervalLiteral(id NodeID, location *NodeLocation, value string, intervalUnit IntervalUnit) *IntervalLiteral {
	// TODO: parse value
	val, _ := strconv.ParseFloat(value, 64)

	var duration time.Duration
	var result time.Time

	baseTime := time.Now()
	switch intervalUnit {
	case Second:
		duration = time.Duration(val * float64(time.Second))
	case Minute:
		duration = time.Duration(val * float64(time.Minute))
	case Hour:
		duration = time.Duration(val * float64(time.Hour))
	case Day:
		duration = time.Duration(val * float64(time.Hour*24))
	case Month:
		months := int(val)
		days := int((val - float64(months)) * 30)
		result = baseTime.AddDate(0, months, days)
		duration = result.Sub(baseTime)
	case Year:
		years := int(val)
		days := int((val - float64(years)) * 365)
		result = baseTime.AddDate(years, 0, days)
		duration = result.Sub(baseTime)
	}

	return &IntervalLiteral{
		BaseNode: BaseNode{
			ID:       id,
			Location: location,
		},
		Unit:  intervalUnit,
		Value: duration,
	}
}

func (n *IntervalLiteral) Interval() timeutil.Interval {
	return timeutil.Interval(n.Value.Milliseconds())
}

func (n *IntervalLiteral) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type Constant struct {
	BaseNode
	Value any
	Type  types.DataType // TODO:
}

func (n *Constant) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}
