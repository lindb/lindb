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
	"context"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/utils"
)

// EvalContext is the context for evaluating expression.
type EvalContext interface {
	// CurrentTime returns the current time.
	CurrentTime() time.Time
}

// evalContext implements EvalContext interface.
type evalContext struct {
	ctx context.Context
}

// NewEvalContext creates an EvalContext.
func NewEvalContext(ctx context.Context) EvalContext {
	return &evalContext{
		ctx: ctx,
	}
}

// CurrentTime returns the current time.
func (e *evalContext) CurrentTime() time.Time {
	return utils.GetTimeFromContext(e.ctx, constants.ContextKeyCurrentTime)
}
