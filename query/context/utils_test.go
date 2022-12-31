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

package context

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
)

func Test_calcTimeRangeAndInterval(t *testing.T) {
	cfg := models.Database{
		Option: &option.DatabaseOption{
			Intervals: option.Intervals{
				{Interval: timeutil.Interval(timeutil.OneSecond)},
				{Interval: timeutil.Interval(timeutil.OneMinute)},
			},
		},
	}
	statement := &stmt.Query{}
	calcTimeRangeAndInterval(statement, cfg)
	assert.Equal(t, timeutil.Interval(timeutil.OneSecond), statement.Interval)

	statement.Interval = timeutil.Interval(timeutil.OneHour)
	statement.TimeRange = timeutil.TimeRange{Start: timeutil.Now(), End: timeutil.Now() + 6*timeutil.OneHour}
	calcTimeRangeAndInterval(statement, cfg)
	assert.Equal(t, timeutil.Interval(timeutil.OneHour), statement.Interval)
}
