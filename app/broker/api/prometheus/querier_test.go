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
	stmtpkg "github.com/lindb/lindb/sql/stmt"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateMatcher(t *testing.T) {
	a := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "a",
		Value: "1",
	}
	b := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "b",
		Value: "2",
	}
	c := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "c",
		Value: "3",
	}
	d := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "d",
		Value: "4",
	}

	expected1 := &stmtpkg.BinaryExpr{
		Left: &stmtpkg.EqualsExpr{
			Key:   a.Name,
			Value: a.Value,
		},
		Right: &stmtpkg.EqualsExpr{
			Key:   b.Name,
			Value: b.Value,
		},
		Operator: stmtpkg.ADD,
	}
	expected2 := &stmtpkg.BinaryExpr{
		Left: &stmtpkg.EqualsExpr{
			Key:   a.Name,
			Value: a.Value,
		},
		Right: &stmtpkg.BinaryExpr{
			Left: &stmtpkg.EqualsExpr{
				Key:   b.Name,
				Value: b.Value,
			},
			Right: &stmtpkg.EqualsExpr{
				Key:   c.Name,
				Value: c.Value,
			},
			Operator: stmtpkg.ADD,
		},
		Operator: stmtpkg.ADD,
	}
	expected3 := &stmtpkg.BinaryExpr{
		Left: &stmtpkg.EqualsExpr{
			Key:   a.Name,
			Value: a.Value,
		},
		Right: &stmtpkg.BinaryExpr{
			Left: &stmtpkg.EqualsExpr{
				Key:   b.Name,
				Value: b.Value,
			},
			Right: &stmtpkg.BinaryExpr{
				Left: &stmtpkg.EqualsExpr{
					Key:   c.Name,
					Value: c.Value,
				},
				Right: &stmtpkg.EqualsExpr{
					Key:   d.Name,
					Value: d.Value,
				},
				Operator: stmtpkg.ADD,
			},
			Operator: stmtpkg.ADD,
		},
		Operator: stmtpkg.ADD,
	}

	param1 := []*labels.Matcher{a, b}
	param2 := []*labels.Matcher{a, b, c}
	param3 := []*labels.Matcher{a, b, c, d}

	root1 := &stmtpkg.BinaryExpr{Operator: stmtpkg.ADD}
	root2 := &stmtpkg.BinaryExpr{Operator: stmtpkg.ADD}
	root3 := &stmtpkg.BinaryExpr{Operator: stmtpkg.ADD}

	walkMatcher(root1, param1)
	assert.Equal(t, root1, expected1)

	walkMatcher(root2, param2)
	assert.Equal(t, root2, expected2)

	walkMatcher(root3, param3)
	assert.Equal(t, root3, expected3)
}

func TestMakeCondition(t *testing.T) {
	a := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "a",
		Value: "1",
	}
	b := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "b",
		Value: "2",
	}
	c := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "c",
		Value: "3",
	}
	d := &labels.Matcher{
		Type:  labels.MatchEqual,
		Name:  "d",
		Value: "4",
	}

	expected1 := &stmtpkg.EqualsExpr{
		Key:   a.Name,
		Value: a.Value,
	}
	expected2 := &stmtpkg.BinaryExpr{
		Left: &stmtpkg.EqualsExpr{
			Key:   a.Name,
			Value: a.Value,
		},
		Right: &stmtpkg.EqualsExpr{
			Key:   b.Name,
			Value: b.Value,
		},
		Operator: stmtpkg.ADD,
	}
	expected3 := &stmtpkg.BinaryExpr{
		Left: &stmtpkg.EqualsExpr{
			Key:   a.Name,
			Value: a.Value,
		},
		Right: &stmtpkg.BinaryExpr{
			Left: &stmtpkg.EqualsExpr{
				Key:   b.Name,
				Value: b.Value,
			},
			Right: &stmtpkg.EqualsExpr{
				Key:   c.Name,
				Value: c.Value,
			},
			Operator: stmtpkg.ADD,
		},
		Operator: stmtpkg.ADD,
	}
	expected4 := &stmtpkg.BinaryExpr{
		Left: &stmtpkg.EqualsExpr{
			Key:   a.Name,
			Value: a.Value,
		},
		Right: &stmtpkg.BinaryExpr{
			Left: &stmtpkg.EqualsExpr{
				Key:   b.Name,
				Value: b.Value,
			},
			Right: &stmtpkg.BinaryExpr{
				Left: &stmtpkg.EqualsExpr{
					Key:   c.Name,
					Value: c.Value,
				},
				Right: &stmtpkg.EqualsExpr{
					Key:   d.Name,
					Value: d.Value,
				},
				Operator: stmtpkg.ADD,
			},
			Operator: stmtpkg.ADD,
		},
		Operator: stmtpkg.ADD,
	}

	metric := &labels.Matcher{
		Name:  metricLabelName,
		Value: "http_requests_total",
	}

	param1 := []*labels.Matcher{a, metric}
	param2 := []*labels.Matcher{a, b, metric}
	param3 := []*labels.Matcher{a, b, c, metric}
	param4 := []*labels.Matcher{a, b, c, d, metric}

	me, expr := makeCondition(param1...)
	assert.Equal(t, me, metric.Value)
	assert.Equal(t, expr, expected1)

	me, expr = makeCondition(param2...)
	assert.Equal(t, me, metric.Value)
	assert.Equal(t, expr, expected2)

	me, expr = makeCondition(param3...)
	assert.Equal(t, me, metric.Value)
	assert.Equal(t, expr, expected3)

	me, expr = makeCondition(param4...)
	assert.Equal(t, me, metric.Value)
	assert.Equal(t, expr, expected4)
}
