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

package models

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseLimits(t *testing.T) {
	assert.Equal(t, defaultLimits, GetDatabaseLimits("test"))
	limits := NewDefaultLimits()
	limits.MaxMetrics = 10
	SetDatabaseLimits("test", limits)
	assert.Equal(t, limits, GetDatabaseLimits("test"))
}

func TestDefaultLimits(t *testing.T) {
	l := NewDefaultLimits()
	val := l.TOML()
	cfg := &Limits{}
	_, err := toml.Decode(val, cfg)
	assert.NoError(t, err)
	assert.Equal(t, cfg, l)

	l.Metrics["system.cpu"] = 1000
	assert.NotEqual(t, l.TOML(), NewDefaultLimits().TOML())
}

func TestLimits_GetSeriesLimits(t *testing.T) {
	l := NewDefaultLimits()
	ns := "ns"
	name := "name"
	assert.Equal(t, l.MaxSeriesPerMetric, l.GetSeriesLimit(ns, name))
	l.Metrics["ns|name"] = 10
	l.Metrics["name"] = 100
	assert.Equal(t, uint32(10), l.GetSeriesLimit(ns, name))
	assert.Equal(t, uint32(100), l.GetSeriesLimit("default-ns", name))
	assert.Equal(t, l.MaxSeriesPerMetric, l.GetSeriesLimit(ns, "test"))
}

func TestLimits_Disable(t *testing.T) {
	l := NewDefaultLimits()
	assert.True(t, l.EnableNamespaceLengthCheck())
	l.MaxNamespaceLength = 0
	assert.False(t, l.EnableNamespaceLengthCheck())
	assert.False(t, l.EnableNamespacesCheck())
	l.MaxNamespaces = 10
	assert.True(t, l.EnableNamespacesCheck())
	assert.True(t, l.EnableMetricNameLengthCheck())
	l.MaxMetricNameLength = 0
	assert.False(t, l.EnableMetricNameLengthCheck())
	assert.False(t, l.EnableMetricsCheck())
	l.MaxMetrics = 10
	assert.True(t, l.EnableMetricsCheck())
	assert.True(t, l.EnableFieldNameLengthCheck())
	l.MaxFieldNameLength = 0
	assert.False(t, l.EnableFieldNameLengthCheck())
	assert.True(t, l.EnableTagNameLengthCheck())
	l.MaxTagNameLength = 0
	assert.False(t, l.EnableTagNameLengthCheck())
	assert.True(t, l.EnableTagValueLengthCheck())
	l.MaxTagValueLength = 0
	assert.False(t, l.EnableTagValueLengthCheck())
	assert.True(t, l.EnableFieldsCheck())
	l.MaxFieldsPerMetric = 0
	assert.False(t, l.EnableFieldsCheck())
	assert.True(t, l.EnableTagsCheck())
	l.MaxTagsPerMetric = 0
	assert.False(t, l.EnableTagsCheck())

	assert.True(t, l.EnableSeriesCheckForQuery())
	l.MaxSeriesPerQuery = 0
	assert.False(t, l.EnableSeriesCheckForQuery())
}
