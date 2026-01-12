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
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
)

type mapValuesFunc struct {
	baseFunc

	keys []string
}

func newMapValuesFunc(ctx EvalContext, args []Expression) Func {
	keys := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		key, _, _ := args[i].EvalString(types.EmptyRow)
		keys[i-1] = key
	}
	return &mapValuesFunc{
		keys: keys,
		baseFunc: baseFunc{
			ctx:  ctx,
			args: args,
		},
	}
}

func (n *mapValuesFunc) EvalMap(row types.Row) (val map[string]string, isNull bool, err error) {
	// TODO: check error/now func
	tsStr, _, _ := n.args[0].EvalMap(row)
	return lo.PickByKeys(tsStr, n.keys), true, nil
}
