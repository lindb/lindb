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

package option

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
)

func TestDatabaseOption_Validate(t *testing.T) {
	databaseOption := DatabaseOption{Interval: "ad"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "-10s"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s"}
	assert.Nil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "aa"}}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"1s", "1m", "1h"}}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"20s", "1m", "1h"}}
	assert.Nil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}, Ahead: "aa"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"10s", "1m", "1h"}, Behind: "aa"}
	assert.NotNil(t, databaseOption.Validate())
	databaseOption = DatabaseOption{Interval: "10s", Rollup: []string{"20s", "1m", "1h"}, Behind: "10h", Ahead: "1h"}
	assert.Nil(t, databaseOption.Validate())
}

func TestDatabaseOption_Default(t *testing.T) {
	databaseOption := DatabaseOption{Interval: "1m"}
	databaseOption.Default()
	assert.Equal(t, databaseOption.Ahead, constants.MetricMaxAheadDurationStr)
	assert.Equal(t, databaseOption.Behind, constants.MetricMaxBehindDurationStr)
	_, _ = databaseOption.GetAcceptWritableRange()
	ahead, behind := databaseOption.GetAcceptWritableRange()
	assert.Equal(t, ahead, int64(constants.MetricMaxAheadDuration))
	assert.Equal(t, behind, int64(constants.MetricMaxBehindDuration))
}
