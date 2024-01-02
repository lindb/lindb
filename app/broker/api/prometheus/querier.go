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

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"reflect"
// 	"sort"
//
// 	"github.com/lindb/lindb/app/broker/api/exec/command"
// 	depspkg "github.com/lindb/lindb/app/broker/deps"
// 	"github.com/lindb/lindb/models"
// 	"github.com/lindb/lindb/pkg/timeutil"
// 	stmtpkg "github.com/lindb/lindb/sql/stmt"
//
// 	commonmodels "github.com/lindb/common/models"
// 	"github.com/prometheus/prometheus/model/labels"
// 	"github.com/prometheus/prometheus/promql"
// 	"github.com/prometheus/prometheus/storage"
// 	"github.com/prometheus/prometheus/util/annotations"
// )
//
// // Queryable is implementation of storage.Queryable of Prometheus.
// type Queryable struct {
// 	deps *depspkg.HTTPDeps
// }
//
// func NewQueryable(deps *depspkg.HTTPDeps) storage.Queryable {
// 	return &Queryable{deps: deps}
// }
//
// func (q *Queryable) Querier(mint, maxt int64) (storage.Querier, error) {
// 	return newQuerier(mint, maxt, q), nil
// }
//
// // Querier is implementation of storage.Querier of Prometheus.
// type Querier struct {
// 	queryable *Queryable
// 	mint      int64
// 	maxt      int64
// }
//
// func newQuerier(mint, maxt int64, queryable *Queryable) *Querier {
// 	return &Querier{
// 		mint:      mint,
// 		maxt:      maxt,
// 		queryable: queryable,
// 	}
// }
//
// // metadata query metrics or tag keys.
// func (q *Querier) metadata(ctx context.Context, metric string, metadataType stmtpkg.MetricMetadataType) ([]string, error) {
// 	var sql string
//
// 	switch metadataType {
// 	case stmtpkg.Metric:
// 		sql = "show metrics"
// 	case stmtpkg.TagKey:
// 		sql = fmt.Sprintf("show tag keys from %s", metric)
// 	default:
// 		return nil, fmt.Errorf("not supported type %d", metadataType)
// 	}
//
// 	param := &models.ExecuteParam{
// 		Database: q.queryable.deps.BrokerCfg.Prometheus.Database,
// 		SQL:      sql,
// 	}
// 	stmt := &stmtpkg.MetricMetadata{
// 		Namespace:  q.queryable.deps.BrokerCfg.Prometheus.Namespace,
// 		Type:       metadataType,
// 		MetricName: metric,
// 	}
//
// 	result, err := command.MetricMetadataCommand(ctx, q.queryable.deps, param, stmt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rs, ok := result.(*commonmodels.Metadata)
// 	if !ok {
// 		return nil, fmt.Errorf("expected type Metadata got %s", reflect.TypeOf(result))
// 	}
//
// 	switch v := rs.Values.(type) {
// 	case []string:
// 		return v, nil
// 	case []commonmodels.Field:
// 		var metrics = make([]string, len(v))
// 		for idx, field := range v {
// 			metrics[idx] = field.Name
// 		}
// 		return metrics, nil
// 	default:
// 		return nil, fmt.Errorf("expected type []string or []Field, got %s", reflect.TypeOf(rs.Values))
// 	}
// }
//
// // LabelValues query label values.
// func (q *Querier) LabelValues(ctx context.Context, name string, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
// 	if name != metricLabelName {
// 		return nil, nil, fmt.Errorf("not supported name %s", name)
// 	}
// 	metrics, err := q.metadata(ctx, "", stmtpkg.Metric)
// 	return metrics, nil, err
// }
//
// // LabelNames query label names.
// func (q *Querier) LabelNames(ctx context.Context, matchers ...*labels.Matcher) ([]string, annotations.Annotations, error) {
// 	metrics, err := q.metadata(ctx, "", stmtpkg.Metric)
// 	return metrics, nil, err
// }
//
// func (q *Querier) Close() error {
// 	return nil
// }
//
// // query time series data using the method in LinDB.
// func (q *Querier) query(ctx context.Context, hints *storage.SelectHints, matchers ...*labels.Matcher) (result any, err error) {
// 	metric, condition := makeCondition(matchers...)
// 	if metric == "" {
// 		return nil, errors.New("metric name does not exist")
// 	}
//
// 	tagKeys, err := q.metadata(ctx, metric, stmtpkg.TagKey)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	param := &models.ExecuteParam{
// 		Database: q.queryable.deps.BrokerCfg.Prometheus.Database,
// 	}
//
// 	stmt := &stmtpkg.Query{
// 		Namespace:  q.queryable.deps.BrokerCfg.Prometheus.Namespace,
// 		MetricName: metric,
// 		AllFields:  true,
// 		Condition:  condition,
// 		TimeRange:  timeutil.TimeRange{Start: hints.Start, End: hints.End},
// 		// let LinDB calculate the interval
// 		Interval: 0,
// 		// the logic of LinDB is that if any tag does not exist, the result is empty.
// 		// if some data has already been persisted and new tag value pair are added later,
// 		// there are some issues with this approach.
// 		// Prometheus, on the other hand, aims to query results even if not all labels exist simultaneously
// 		GroupBy: tagKeys,
// 		Limit:   1e5,
// 	}
//
// 	return command.QueryCommand(ctx, q.queryable.deps, param, stmt)
// }
//
// // Select calls the query method to retrieve data and transforms the results into the SeriesSet.
// func (q *Querier) Select(ctx context.Context, sortSeries bool, hints *storage.SelectHints, matchers ...*labels.Matcher) storage.SeriesSet {
// 	var set = newSeriesSet()
// 	result, err := q.query(ctx, hints, matchers...)
// 	if err != nil {
// 		set.setErr(err)
// 		return set
// 	}
//
// 	rs, ok := result.(*commonmodels.ResultSet)
// 	if !ok {
// 		set.setErr(fmt.Errorf("expected ResultSet type got %s type", reflect.TypeOf(result).String()))
// 		return set
// 	}
//
// 	if len(rs.Fields) == 0 {
// 		return set
// 	}
//
// 	field := rs.Fields[0]
// 	seriesSlice := make([]storage.Series, 0)
// 	for _, s := range rs.Series {
// 		lbs := []string{metricLabelName, rs.MetricName}
// 		for k, v := range s.Tags {
// 			lbs = append(lbs, k, v)
// 		}
// 		var points []promql.FPoint
// 		if hints.Func != "series" {
// 			points = make([]promql.FPoint, 0, len(s.Fields[field]))
// 			for t, f := range s.Fields[field] {
// 				fp := promql.FPoint{
// 					T: t,
// 					F: f,
// 				}
// 				points = append(points, fp)
// 			}
// 			sort.Slice(points, func(i, j int) bool {
// 				return points[i].T < points[j].T
// 			})
// 			if sortSeries {
// 				sort.Strings(lbs)
// 			}
// 		}
// 		one := promql.NewStorageSeries(promql.Series{
// 			Metric:     labels.FromStrings(lbs...),
// 			Floats:     points,
// 			Histograms: nil,
// 		})
// 		seriesSlice = append(seriesSlice, one)
// 	}
//
// 	set.setSeries(seriesSlice)
//
// 	return set
// }
