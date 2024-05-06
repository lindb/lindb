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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/logger"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"

	"github.com/lindb/lindb/app/broker/api/prometheus/ingest"
)

// remoteWrite implements a remote write interface similar to Prometheus.
func (e *ExecuteAPI) remoteWrite(c *gin.Context) {
	r, w := c.Request, c.Writer
	req, err := remote.DecodeWriteRequest(c.Request.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.write(r.Context(), req)

	w.WriteHeader(http.StatusNoContent)
}

// write asynchronously writes data to LinDB.
func (e *ExecuteAPI) write(ctx context.Context, req *prompb.WriteRequest) {
	for _, ts := range req.Timeseries { //nolint:gocritic
		lbs := labelProtosToLabels(ts.Labels)
		if !lbs.IsValid() {
			e.logger.Warn("Invalid metric names or labels", logger.String("labels", lbs.String()))
			continue
		}
		for _, sample := range ts.Samples {
			point := new(ingest.Point)
			for _, label := range lbs {
				if label.Name == metricLabelName {
					point.SetMetricName(label.Value)
				} else {
					point.AddTag(label.Name, label.Value)
				}
			}
			var timestamp time.Time
			if sample.Timestamp <= 0 {
				timestamp = time.Now()
			} else {
				timestamp = time.UnixMilli(sample.Timestamp)
			}
			point.SetNamespace(e.deps.BrokerCfg.Prometheus.Namespace)
			point.SetTimestamp(timestamp)
			point.AddField(ingest.NewLast(e.deps.BrokerCfg.Prometheus.Field, sample.Value))
			e.prometheusWriter.AddPoint(ctx, point)
		}
	}
}
