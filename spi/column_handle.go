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

package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/sql/tree"
)

func init() {
	// register json encoder/decoder for column handle
	jsoniter.RegisterTypeEncoder("spi.ColumnHandle", &encoding.JSONEncoder[ColumnHandle]{})
	jsoniter.RegisterTypeDecoder("spi.ColumnHandle", &encoding.JSONDecoder[ColumnHandle]{})
}

type ColumnHandle interface{}

type ColumnAssignment struct {
	Handler ColumnHandle `json:"handler"`
	Column  string       `json:"column"`
}

type ColumnAggregation struct {
	Column      string
	AggFuncName tree.FuncName
}
