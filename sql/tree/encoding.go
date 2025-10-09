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
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
)

func init() {
	jsoniter.RegisterTypeEncoder("tree.Expression", &encoding.JSONEncoder[Expression]{})
	jsoniter.RegisterTypeDecoder("tree.Expression", &encoding.JSONDecoder[Expression]{})

	encoding.RegisterNodeType(ComparisonExpression{})
	encoding.RegisterNodeType(NotExpression{})
	encoding.RegisterNodeType(InPredicate{})
	encoding.RegisterNodeType(LikePredicate{})
	encoding.RegisterNodeType(RegexPredicate{})
	encoding.RegisterNodeType(NullPredicate{})
	encoding.RegisterNodeType(InListExpression{})
	encoding.RegisterNodeType(LogicalExpression{})
	encoding.RegisterNodeType(TimePredicate{})
	encoding.RegisterNodeType(StringLiteral{})
	encoding.RegisterNodeType(LongLiteral{})
	encoding.RegisterNodeType(Identifier{})
	encoding.RegisterNodeType(Cast{})
	encoding.RegisterNodeType(FunctionCall{})
	encoding.RegisterNodeType(SymbolReference{})
	encoding.RegisterNodeType(Constant{})
}
