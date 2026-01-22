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
	"fmt"
	"reflect"
	"time"

	"github.com/samber/lo"
)

type Datum struct {
	val any
}

func MakeDatums(vals ...any) []*Datum {
	return lo.Map(vals, func(v any, _ int) *Datum {
		return &Datum{val: v}
	})
}

func (d *Datum) String() string {
	if str, ok := d.val.(string); ok {
		return str
	}
	return fmt.Sprint(d.val)
}

func (d *Datum) Float() float64 {
	v := reflect.ValueOf(d.val)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	default:
		return 0 // TODO:set 0?
	}
}

func (d *Datum) Int() int64 {
	v := reflect.ValueOf(d.val)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Int()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint())
	default:
		return 0 // TODO:set 0?
	}
}

func (d *Datum) Duration() time.Duration {
	return d.val.(time.Duration)
}
