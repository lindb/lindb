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

	"github.com/lindb/lindb/spi/types"
)

type Expression interface {
	EvalInt(ctx EvalContext, row types.Row) (val int64, isNull bool, err error)
	EvalString(ctx EvalContext, row types.Row) (val string, isNull bool, err error)
	EvalFloat(ctx EvalContext, row types.Row) (val float64, isNull bool, err error)
	EvalDuration(ctx EvalContext, row types.Row) (val time.Duration, isNull bool, err error)
	EvalTimeSeries(ctx EvalContext, row types.Row) (val *types.TimeSeries, isNull bool, err error)
	EvalTime(ctx EvalContext, row types.Row) (val time.Time, isNull bool, err error)
	EvalMap(ctx EvalContext, row types.Row) (val map[string]string, isNull bool, err error)
	// Getype returns the data type of the expression returns.
	GetType() types.DataType
	// String returns the expression in string format.
	String() string
}
