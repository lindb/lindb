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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/logger"
	"github.com/prometheus/common/version"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/textparse"
	"github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/httputil"
	"github.com/prometheus/prometheus/web"
	v1 "github.com/prometheus/prometheus/web/api/v1"

	prometheusIngest "github.com/lindb/lindb/app/broker/api/prometheus/ingest"
	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/prometheus"
)

const (
	metricLabelName = "__name__"
)

// doc https://prometheus.io/docs/prometheus/latest/querying/api/#http-api
const (
	executeQueryPath           = "/query"
	executeQueryRangePath      = "/query_range"
	executeFormatQueryPath     = "/format_query"
	executeStatusBuildInfoPath = "/status/buildinfo"
	executeLabelsPath          = "/labels"
	executeLabelValuesPath     = "/label/:name/values"
	executeMetadataPath        = "/metadata"
	executeSeriesPath          = "/series"
	executeWritePath           = "/write"
)

// ExecuteAPI wraps all Prometheus APIs.
type ExecuteAPI struct {
	prometheusWriter prometheusIngest.Writer
	queryable        storage.Queryable
	logger           logger.Logger
	deps             *depspkg.HTTPDeps
	engine           *promql.Engine
	codecs           []Codec
}

// NewExecuteAPI creates a promql execution api.
func NewExecuteAPI(deps *depspkg.HTTPDeps, prometheusWriter prometheusIngest.Writer) *ExecuteAPI {
	e := &ExecuteAPI{
		deps:             deps,
		prometheusWriter: prometheusWriter,
		logger:           logger.GetLogger("broker", "prometheus.ExecuteAPI"),
		codecs:           []Codec{&JSONCodec{}},
		engine:           prometheus.NewEngine(),
		// queryable:        NewQueryable(deps),
	}
	go e.waitPrometheusWriteErr()
	return e
}

// Register adds promql executor's path.
func (e *ExecuteAPI) Register(route gin.IRoutes) {
	// query
	route.GET(executeQueryPath, e.query)
	route.POST(executeQueryPath, e.query)

	// range query
	route.GET(executeQueryRangePath, e.queryRange)
	route.POST(executeQueryRangePath, e.queryRange)

	// format query
	route.GET(executeFormatQueryPath, e.formatQuery)
	route.POST(executeFormatQueryPath, e.formatQuery)

	// metadata
	route.GET(executeMetadataPath, e.metricMetadata)

	// build info
	route.GET(executeStatusBuildInfoPath, e.buildInfo)

	// label name
	route.GET(executeLabelsPath, e.labelNames)
	route.POST(executeLabelsPath, e.labelNames)

	// label values
	route.GET(executeLabelValuesPath, e.labelValues)

	// series
	route.GET(executeSeriesPath, e.series)
	route.POST(executeSeriesPath, e.series)

	// remote write
	route.POST(executeWritePath, e.remoteWrite)
}

// waitPrometheusWriteErr watch write error
func (e *ExecuteAPI) waitPrometheusWriteErr() {
	errCh := e.prometheusWriter.Errors()
	for {
		err, ok := <-errCh
		e.logger.Error("prometheus remote write", logger.Error(err))
		if !ok {
			break
		}
	}
}

// execute call fn wrap by QueryLimiter.
func (e *ExecuteAPI) execute(c *gin.Context, fn func(c *gin.Context) apiFuncResult) {
	if err := e.deps.QueryLimiter.Do(func() error {
		e.response(c, fn(c))
		return nil
	}); err != nil {
		e.response(c, apiFuncResult{
			err: &apiError{
				typ: errorUnavailable,
				err: err,
			},
		})
	}
}

// series query all label values based on the 'match[]' parameter.
func (e *ExecuteAPI) series(c *gin.Context) {
	e.execute(c, e.querySeries)
}

// querySeries is the implementation of series.
func (e *ExecuteAPI) querySeries(c *gin.Context) apiFuncResult {
	r, ctx := c.Request, c.Request.Context()

	if err := r.ParseForm(); err != nil {
		return apiFuncResult{nil, &apiError{errorBadData, fmt.Errorf("error parsing form values: %w", err)}, nil, nil}
	}
	if len(r.Form["match[]"]) == 0 {
		return apiFuncResult{nil, &apiError{errorBadData, errors.New("no match[] parameter provided")}, nil, nil}
	}

	start, err := parseTimeParam(r, "start", MinTime)
	if err != nil {
		return invalidParamError(err, "start")
	}
	end, err := parseTimeParam(r, "end", MaxTime)
	if err != nil {
		return invalidParamError(err, "end")
	}

	matcherSets, err := parseMatchersParam(r.Form["match[]"])
	if err != nil {
		return invalidParamError(err, "match[]")
	}

	q, err := e.queryable.Querier(timestamp.FromTime(start), timestamp.FromTime(end))
	if err != nil {
		return apiFuncResult{nil, returnAPIError(err), nil, nil}
	}

	hints := &storage.SelectHints{
		Start: timestamp.FromTime(start),
		End:   timestamp.FromTime(end),
		// There is no series function, this token is used for lookups that don't need samples.
		Func: "series",
	}

	var set storage.SeriesSet

	if len(matcherSets) > 1 {
		var sets []storage.SeriesSet
		for _, mset := range matcherSets {
			// We need to sort this select results to merge (deduplicate) the series sets later.
			s := q.Select(ctx, true, hints, mset...)
			sets = append(sets, s)
		}
		set = storage.NewMergeSeriesSet(sets, storage.ChainedSeriesMerge)
	} else {
		// At this point at least one match exists.
		set = q.Select(ctx, false, hints, matcherSets[0]...)
	}

	var metrics []labels.Labels
	for set.Next() {
		metrics = append(metrics, set.At().Labels())
	}

	warnings := set.Warnings()
	if set.Err() != nil {
		return apiFuncResult{nil, returnAPIError(set.Err()), warnings, nil}
	}

	return apiFuncResult{metrics, nil, warnings, nil}
}

// metricMetadata query all metric names.
func (e *ExecuteAPI) metricMetadata(c *gin.Context) {
	e.execute(c, e.queryMetricMetadata)
}

// queryMetricMetadata is the implementation of metricMetadata.
func (e *ExecuteAPI) queryMetricMetadata(c *gin.Context) apiFuncResult {
	querier, err := e.queryable.Querier(0, 0)
	if err != nil {
		return unavailableError(err)
	}

	metrics, _, err := querier.LabelValues(c, metricLabelName, nil)
	if err != nil {
		return unavailableError(err)
	}

	res := map[string][]metadata{}
	for _, metric := range metrics {
		res[metric] = []metadata{{Type: textparse.MetricTypeUnknown}}
	}

	return apiFuncResult{
		data: res,
	}
}

// labelNames query label names.
func (e *ExecuteAPI) labelNames(c *gin.Context) {
	e.execute(c, e.queryLabelNames)
}

// queryLabelName is the implementation of labelNames.
func (e *ExecuteAPI) queryLabelNames(c *gin.Context) apiFuncResult {
	querier, err := e.queryable.Querier(0, 0)
	if err != nil {
		return unavailableError(err)
	}

	metrics, _, err := querier.LabelNames(c, nil)
	if err != nil {
		return unavailableError(err)
	}

	ms := make([]string, len(metrics)+1)
	ms[0] = metricLabelName
	copy(ms[1:], metrics)

	return apiFuncResult{
		data: metrics,
	}
}

// labelValues query label values.
// doc see https://prometheus.io/docs/prometheus/latest/querying/api/#querying-label-values
func (e *ExecuteAPI) labelValues(c *gin.Context) {
	e.execute(c, e.queryLabelValues)
}

// queryLabelValues is the implementation of labelValues
func (e *ExecuteAPI) queryLabelValues(c *gin.Context) (result apiFuncResult) {
	labelName := c.Param("name")
	// for now only metric queries are implemented.
	if labelName != metricLabelName {
		return apiFuncResult{
			data: []string{},
		}
	}

	querier, err := e.queryable.Querier(0, 0)
	if err != nil {
		return unavailableError(err)
	}

	metrics, _, err := querier.LabelNames(c, nil)
	if err != nil {
		return unavailableError(err)
	}

	return apiFuncResult{
		data: metrics,
	}
}

// query queries time series data.
func (e *ExecuteAPI) query(c *gin.Context) {
	e.execute(c, e.queryResult)
}

// queryResult is the implementation of query.
func (e *ExecuteAPI) queryResult(c *gin.Context) (result apiFuncResult) {
	r := c.Request
	ts, err := parseTimeParam(r, "time", time.Now())
	if err != nil {
		return invalidParamError(err, "time")
	}
	ctx := r.Context()
	if to := r.FormValue("timeout"); to != "" {
		var cancel context.CancelFunc
		timeout, err0 := parseDuration(to)
		if err0 != nil {
			return invalidParamError(err0, "timeout")
		}

		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(timeout))
		defer cancel()
	}

	opts, err := extractQueryOpts(r)
	if err != nil {
		return apiFuncResult{nil, &apiError{errorBadData, err}, nil, nil}
	}

	qry, err := e.engine.NewInstantQuery(ctx, e.queryable, opts, r.FormValue("query"), ts)
	if err != nil {
		return invalidParamError(err, "query")
	}

	param := models.ExecuteParam{SQL: r.FormValue("query")}
	c.Set(constants.CurrentSQL, &param)

	ctx = httputil.ContextFromRequest(ctx, r)

	res := qry.Exec(ctx)
	if res.Err != nil {
		return apiFuncResult{nil, returnAPIError(res.Err), res.Warnings, qry.Close}
	}

	return apiFuncResult{&v1.QueryData{
		ResultType: res.Value.Type(),
		Result:     res.Value,
	}, nil, res.Warnings, qry.Close}
}

// queryRange queries time series data within a time range.
func (e *ExecuteAPI) queryRange(c *gin.Context) {
	e.execute(c, e.queryRangeResult)
}

// queryRangeResult is the implementation of queryRange.
func (e *ExecuteAPI) queryRangeResult(c *gin.Context) (result apiFuncResult) {
	r := c.Request
	start, err := parseTime(r.FormValue("start"))
	if err != nil {
		return invalidParamError(err, "start")
	}
	end, err := parseTime(r.FormValue("end"))
	if err != nil {
		return invalidParamError(err, "end")
	}
	if end.Before(start) {
		return invalidParamError(errors.New("end timestamp must not be before start time"), "end")
	}

	step, err := parseDuration(r.FormValue("step"))
	if err != nil {
		return invalidParamError(err, "step")
	}

	if step <= 0 {
		return invalidParamError(errors.New("zero or negative query resolution step widths are not accepted. Try a positive integer"), "step")
	}

	// For safety, limit the number of returned points per timeseries.
	// This is sufficient for 60s resolution for a week or 1h resolution for a year.
	if end.Sub(start)/step > 11000 {
		err0 := errors.New("exceeded maximum resolution of 11,000 points per timeseries. Try decreasing the query resolution (?step=XX)")
		return apiFuncResult{nil, &apiError{errorBadData, err0}, nil, nil}
	}

	ctx := r.Context()
	if to := r.FormValue("timeout"); to != "" {
		var cancel context.CancelFunc
		timeout, err0 := parseDuration(to)
		if err0 != nil {
			return invalidParamError(err0, "timeout")
		}

		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, err = extractQueryOpts(r)
	if err != nil {
		return apiFuncResult{nil, &apiError{errorBadData, err}, nil, nil}
	}

	// queryable := NewQueryable(e.deps)
	//
	// qry, err := e.engine.NewRangeQuery(ctx, queryable, opts, r.FormValue("query"), start, end, step)
	// if err != nil {
	// 	return invalidParamError(err, "query")
	// }
	//
	// param := models.ExecuteParam{SQL: r.FormValue("query")}
	// c.Set(constants.CurrentSQL, &param)
	//
	// ctx = httputil.ContextFromRequest(ctx, r)
	//
	// res := qry.Exec(ctx)
	// if res.Err != nil {
	// 	return apiFuncResult{nil, returnAPIError(res.Err), res.Warnings, qry.Close}
	// }
	//
	return apiFuncResult{&v1.QueryData{
		// 	ResultType: res.Value.Type(),
		// 	Result:     res.Value,
	}, nil, nil, nil}
}

// formatQuery formats the query statement.
func (e *ExecuteAPI) formatQuery(c *gin.Context) {
	r := c.Request
	expr, err := parser.ParseExpr(r.FormValue("query"))
	if err != nil {
		e.response(c, invalidParamError(err, "query"))
		return
	}
	e.response(c, apiFuncResult{data: expr.Pretty(0)})
}

// buildInfo returns build information.
func (e *ExecuteAPI) buildInfo(c *gin.Context) {
	info := &web.PrometheusVersion{
		Version:   version.Version,
		Revision:  version.Revision,
		Branch:    version.Branch,
		BuildUser: version.BuildUser,
		BuildDate: version.BuildDate,
		GoVersion: version.GoVersion,
	}
	e.response(c, apiFuncResult{data: info})
}
