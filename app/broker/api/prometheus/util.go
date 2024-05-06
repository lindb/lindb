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
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"

	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

func parseTimeParam(r *http.Request, paramName string, defaultValue time.Time) (time.Time, error) {
	val := r.FormValue(paramName)
	if val == "" {
		return defaultValue, nil
	}
	result, err := parseTime(val)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time value for '%s', err:%w", paramName, err)
	}
	return result, nil
}

func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		ns = math.Round(ns*1000) / 1000
		return time.Unix(int64(s), int64(ns*float64(time.Second))).UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}

	// Stdlib's time parser can only handle 4 digit years. As a workaround until
	// that is fixed we want to at least support our own boundary times.
	// Context: https://github.com/prometheus/client_golang/issues/614
	// Upstream issue: https://github.com/golang/go/issues/20555
	switch s {
	case minTimeFormatted:
		return MinTime, nil
	case maxTimeFormatted:
		return MaxTime, nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q to a valid timestamp", s)
}

func parseDuration(s string) (time.Duration, error) {
	if d, err := strconv.ParseFloat(s, 64); err == nil {
		ts := d * float64(time.Second)
		if ts > float64(math.MaxInt64) || ts < float64(math.MinInt64) {
			return 0, fmt.Errorf("cannot parse %q to a valid duration. It overflows int64", s)
		}
		return time.Duration(ts), nil
	}
	if d, err := model.ParseDuration(s); err == nil {
		return time.Duration(d), nil
	}
	return 0, fmt.Errorf("cannot parse %q to a valid duration", s)
}

func invalidParamError(err error, parameter string) apiFuncResult {
	return apiFuncResult{nil, &apiError{
		errorBadData, fmt.Errorf("invalid parameter %q, err: %w", parameter, err),
	}, nil, nil}
}

func unavailableError(err error) apiFuncResult {
	return apiFuncResult{
		err: &apiError{
			typ: errorUnavailable,
			err: err,
		},
	}
}

func extractQueryOpts(r *http.Request) (promql.QueryOpts, error) {
	var duration time.Duration

	if strDuration := r.FormValue("lookback_delta"); strDuration != "" {
		parsedDuration, err := parseDuration(strDuration)
		if err != nil {
			return nil, fmt.Errorf("error parsing lookback delta duration: %w", err)
		}
		duration = parsedDuration
	}

	return promql.NewPrometheusQueryOpts(r.FormValue("stats") == "all", duration), nil
}

func parseMatchersParam(matchers []string) ([][]*labels.Matcher, error) {
	var matcherSets [][]*labels.Matcher
	for _, s := range matchers {
		matchers, err := parser.ParseMetricSelector(s)
		if err != nil {
			return nil, err
		}
		matcherSets = append(matcherSets, matchers)
	}

OUTER:
	for _, ms := range matcherSets {
		for _, lm := range ms {
			if lm != nil && !lm.Matches("") {
				continue OUTER
			}
		}
		return nil, errors.New("match[] must contain at least one non-empty matcher")
	}
	return matcherSets, nil
}

// walkMatcher iterates matchers and make binary tree.
func walkMatcher(root *stmtpkg.BinaryExpr, matchers []*labels.Matcher) {
	if root == nil || len(matchers) == 0 {
		return
	}
	if root.Left == nil {
		root.Left = &stmtpkg.EqualsExpr{
			Key:   matchers[0].Name,
			Value: matchers[0].Value,
		}
	} else if root.Right == nil {
		if len(matchers) > 1 {
			expr := &stmtpkg.BinaryExpr{
				Left: &stmtpkg.EqualsExpr{
					Key:   matchers[0].Name,
					Value: matchers[0].Value,
				},
				Operator: stmtpkg.ADD,
			}
			root.Right = expr
			root = expr
		} else {
			root.Right = &stmtpkg.EqualsExpr{
				Key:   matchers[0].Name,
				Value: matchers[0].Value,
			}
		}
	}

	matchers = matchers[1:]
	walkMatcher(root, matchers)
}

// makeCondition extracts metric name and condition from matchers.
func makeCondition(matchers ...*labels.Matcher) (metricName string, expr stmtpkg.Expr) {
	pureMatchers := make([]*labels.Matcher, 0, len(matchers)-1)
	for index := range matchers {
		matcher := matchers[index]
		if matcher.Name == metricLabelName {
			metricName = matcher.Value
		} else {
			pureMatchers = append(pureMatchers, matcher)
		}
	}

	switch len(pureMatchers) {
	case 0:
		return metricName, nil
	case 1:
		return metricName, &stmtpkg.EqualsExpr{
			Key:   pureMatchers[0].Name,
			Value: pureMatchers[0].Value,
		}
	default:
		e := &stmtpkg.BinaryExpr{Operator: stmtpkg.ADD}
		walkMatcher(e, pureMatchers)
		return metricName, e
	}
}

func labelProtosToLabels(labelPairs []prompb.Label) labels.Labels {
	b := labels.ScratchBuilder{}
	for _, l := range labelPairs {
		b.Add(l.Name, l.Value)
	}
	b.Sort()
	return b.Labels()
}
