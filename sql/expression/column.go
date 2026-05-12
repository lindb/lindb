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

package expression

import (
	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/spi/scalar"
)

type Column struct {
	name  string
	index int

	rt ResultType
}

func NewColumn(ctx EvalContext, name string, index int, rt ResultType) Expression {
	return &Column{name: name, index: index, rt: rt}
}

func (c *Column) EvalScalar() (scalar.Scalar, error) {
	panic("column is not supported in scalar execution")
}

func (c *Column) Eval(record arrow.RecordBatch) (arrow.Array, error) {
	return record.Column(c.index), nil
}

func (c *Column) ResultType() ResultType {
	return c.rt
}

// String returns the column in string format.
func (c *Column) String() string {
	return c.name
}
