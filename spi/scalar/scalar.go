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

package scalar

import (
	"fmt"
	"time"
)

type Scalar any

type String struct {
	Value string
}

func NewStringScalar(value string) *String {
	return &String{Value: value}
}

type Int64 struct {
	Value int64
}

func NewInt64Scalar(value int64) *Int64 {
	return &Int64{Value: value}
}

type Float64 struct {
	Value float64
}

func NewFloat64Scalar(value float64) *Float64 {
	return &Float64{Value: value}
}

type Duration struct {
	Value time.Duration
}

func NewDurationScalar(value time.Duration) *Duration {
	return &Duration{Value: value}
}

type TimeScalar struct {
	Value time.Time
}

func NewTimeScalar(value time.Time) *TimeScalar {
	return &TimeScalar{Value: value}
}

func ToString(scalar Scalar) string {
	strScalar, ok := scalar.(*String)
	if !ok {
		panic(fmt.Sprintf("unexpected scalar type: %T", scalar))
	}
	return strScalar.Value
}

func ToInt64(scalar Scalar) int64 {
	intScalar, ok := scalar.(*Int64)
	if !ok {
		panic(fmt.Sprintf("unexpected scalar type: %T", scalar))
	}
	return intScalar.Value
}

func ToFloat64(scalar Scalar) float64 {
	floatScalar, ok := scalar.(*Float64)
	if !ok {
		panic(fmt.Sprintf("unexpected scalar type: %T", scalar))
	}
	return floatScalar.Value
}

func ToDuration(scalar Scalar) time.Duration {
	durationScalar, ok := scalar.(*Duration)
	if !ok {
		panic(fmt.Sprintf("unexpected scalar type: %T", scalar))
	}
	return durationScalar.Value
}

func ToTime(scalar Scalar) time.Time {
	timeScalar, ok := scalar.(*TimeScalar)
	if !ok {
		panic(fmt.Sprintf("unexpected scalar type: %T", scalar))
	}
	return timeScalar.Value
}
