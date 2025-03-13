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

package flow

import (
	"context"
	"time"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

// TaskContext represents task execute context.
type TaskContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	Start time.Time
}

// NewTaskContextWithTimeout creates a task context with timeout.
func NewTaskContextWithTimeout(ctx context.Context, timeout time.Duration) *TaskContext {
	c, cancel := context.WithTimeout(ctx, timeout)
	return &TaskContext{
		Ctx:    c,
		Cancel: cancel,
		Start:  time.Now(),
	}
}

// Release releases context's resource after query.
func (ctx *TaskContext) Release() {
	ctx.Cancel()
}

// TagFilterResult represents the tag filter result, include tag key id and tag value ids.
type TagFilterResult struct {
	TagKeyID    tag.KeyID
	TagValueIDs *roaring.Bitmap
}

type MetricScanContext struct {
	MetricID  metric.ID
	SeriesIDs *roaring.Bitmap
	Fields    field.Metas
	TimeRange timeutil.TimeRange
}
