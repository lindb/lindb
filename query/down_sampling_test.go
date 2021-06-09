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

package query

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
)

func Test_downSamplingTimeRange(t *testing.T) {
	timeRange, intervalRatio, interval := downSamplingTimeRange(
		timeutil.Interval(30*timeutil.OneSecond),
		timeutil.Interval(10*timeutil.OneSecond),
		timeutil.TimeRange{
			Start: 35 * timeutil.OneSecond,
			End:   65 * timeutil.OneSecond,
		})
	assert.Equal(t, 3, intervalRatio)
	assert.Equal(t, 30*timeutil.OneSecond, interval.Int64())
	assert.Equal(t, timeutil.TimeRange{
		Start: 30 * timeutil.OneSecond,
		End:   60 * timeutil.OneSecond,
	}, timeRange)
}
