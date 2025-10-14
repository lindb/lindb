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
	"github.com/lindb/lindb/spi/types"
)

// FuncName represents function name.
type FuncName string

// AggFuncName represents aggregation function name.
type AggFuncName string

const (
	// scalar function names
	Plus  FuncName = "plus"
	Minus FuncName = "minus"
	Div   FuncName = "div"
	Mul   FuncName = "mul"
	Mod   FuncName = "mod"

	// aggregation function names
	Sum   FuncName = "sum"
	Min   FuncName = "min"
	Max   FuncName = "max"
	First FuncName = "first"
	Last  FuncName = "last"
	Count FuncName = "count"

	// time function names
	DateAdd   FuncName = "date_add"
	Now       FuncName = "now"
	StrToDate FuncName = "str_to_date"
)

func GetDefaultFuncReturnType(name FuncName) types.DataType {
	return defaultFuncReturnTypes[name]
}

var defaultFuncReturnTypes = map[FuncName]types.DataType{
	DateAdd:   types.DTTimestamp,
	Now:       types.DTTimestamp,
	StrToDate: types.DTTimestamp,
	Count:     types.DTTimeSeries,
}
