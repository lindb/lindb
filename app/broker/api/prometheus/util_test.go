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

package prometheus

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
)

func TestParseMatchersParam(t *testing.T) {
	result, err := parseMatchersParam([]string{`http_requests_total{idc="sh"}`, `cpu_load{ip="1.1.1.1"}`})
	assert.Nil(t, err)
	expected := []*labels.Matcher{
		{
			Type:  labels.MatchEqual,
			Name:  "idc",
			Value: "sh",
		},
		{
			Type:  labels.MatchEqual,
			Name:  "__name__",
			Value: "http_requests_total",
		},
		{
			Type:  labels.MatchEqual,
			Name:  "ip",
			Value: "1.1.1.1",
		},
		{
			Type:  labels.MatchEqual,
			Name:  "__name__",
			Value: "cpu_load",
		},
	}

	got := make([]*labels.Matcher, 0)
	for _, slice := range result {
		got = append(got, slice...)
	}

	assert.Equal(t, expected, got)
}
