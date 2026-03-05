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
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
)

type newRule func(mapper *StringMapper) Rule

var rules = map[arrow.Type]newRule{
	(&arrow.MapType{}).ID():                 newMapRule,
	arrow.BinaryTypes.String.ID():           newStringRule,
	arrow.FixedWidthTypes.Timestamp_ns.ID(): newTimestampRule,
}

func CreateRule(dt arrow.DataType, mapper *StringMapper) (Rule, error) {
	newFn, ok := rules[dt.ID()]
	if !ok {
		return nil, fmt.Errorf("rule not support data type: %s", dt)
	}
	return newFn(mapper), nil
}
