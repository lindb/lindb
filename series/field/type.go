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

package field

import (
	"github.com/lindb/common/models"
)

// EmptyFieldID represents empty value for field id.
const EmptyFieldID = ID(0)

// ID represents field id.
type ID uint8

// Name represents field name.
type Name string

func (n Name) String() string {
	return string(n)
}

// ExemplarAggregate merges two exemplars, returning the one with longer trace duration.
// This is LinDB-specific semantics and is kept here rather than in the common package.
func ExemplarAggregate(a, b *models.Exemplar) *models.Exemplar {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.Duration >= b.Duration {
		return a
	}
	return b
}
