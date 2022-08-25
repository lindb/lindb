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

package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
)

func TestGetValue(t *testing.T) {
	ctx := context.WithValue(context.TODO(), constants.ContextKeySQL, "value")
	val, ok := GetStringFromContext(ctx, constants.ContextKeySQL)
	assert.True(t, ok)
	assert.Equal(t, "value", val)
	assert.NotNil(t, GetFromContext(ctx, constants.ContextKeySQL))

	val, ok = GetStringFromContext(ctx, constants.ContextKey("no"))
	assert.False(t, ok)
	assert.Empty(t, val)
	assert.Nil(t, GetFromContext(ctx, constants.ContextKey("no")))
}
