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

package function

// FuncType is the definition of function type
type FuncType int

const (
	Sum FuncType = iota + 1
	Min
	Max
	Count
	Avg
	LastValue
	Histogram
	Stddev

	Unknown
)

// String return the function's name
func (t FuncType) String() string {
	switch t {
	case Sum:
		return "sum"
	case Min:
		return "min"
	case Max:
		return "max"
	case Count:
		return "count"
	case Avg:
		return "avg"
	case LastValue:
		return "last_value"
	case Histogram:
		return "histogram"
	case Stddev:
		return "stddev"
	default:
		return "unknown"
	}
}
