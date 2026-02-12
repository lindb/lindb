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
	"time"

	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/spi/scalar"
	"github.com/lindb/lindb/sql/tree"
)

var log = logger.GetLogger("Expression", "Eval")

func EvalTime(ctx EvalContext, expression tree.Expression) (time.Time, error) {
	expr := Rewrite(&RewriteContext{}, expression)
	val, err := expr.EvalScalar()
	if err != nil {
		return time.Time{}, err
	}
	return scalar.ToTime(val), nil
}

func EvalString(ctx EvalContext, expression tree.Expression) (string, error) {
	expr := Rewrite(&RewriteContext{}, expression)
	val, err := expr.EvalScalar()
	if err != nil {
		return "", err
	}
	return scalar.ToString(val), nil
}

func EvalProps(ctx EvalContext, props []*tree.Property) (*collections.Properties, error) {
	rs := collections.NewProperties()
	for _, prop := range props {
		expr := prop.Value
		switch propExpr := expr.(type) {
		case *tree.ArrayExpression:
			val := lo.Map(propExpr.Elements, func(item tree.Expression, index int) string {
				val, err := EvalString(ctx, item)
				if err != nil {
					log.Warn("failed to eval property array element", logger.Any("prop", item), logger.Error(err))
					return ""
				}
				return val
			})
			rs.Set(prop.Name.Value, val)
		default:
			val, err := EvalString(ctx, propExpr)
			if err != nil {
				return nil, err
			}
			rs.Set(prop.Name.Value, val)
		}
	}
	return rs, nil
}
