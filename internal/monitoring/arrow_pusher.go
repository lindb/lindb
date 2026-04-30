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

package monitoring

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/apache/arrow-go/v18/arrow/memory"
	lmetrics "github.com/lindb/arrow/pkg/metrics"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/series/tag"
)

var arrowPushLogger = logger.GetLogger("Monitoring", "ArrowPusher")

// arrowNativePusher implements NativePusher, pushing metrics as Arrow IPC bytes.
type arrowNativePusher struct {
	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
	endpoint string
	database string
	gather   linmetric.ArrowGather
	builder  *lmetrics.MetricBuilder
	client   *http.Client

	statistics struct {
		pushBytesCounter   *linmetric.BoundCounter
		pushMetricsCounter *linmetric.BoundCounter
		pushErrorCounter   *linmetric.BoundCounter
	}
}

// NewArrowNativePusher creates an ArrowNativePusher that collects metrics from r
// and periodically sends them as Arrow IPC bytes to endpoint.
// The database name is sent via the X-LinDB-Database HTTP header on each push.
func NewArrowNativePusher(
	ctx context.Context,
	endpoint string,
	database string,
	interval time.Duration,
	pushTimeout time.Duration,
	r *linmetric.Registry,
	globalKeyValues tag.Tags,
) NativePusher {
	c, cancel := context.WithCancel(ctx)
	p := &arrowNativePusher{
		ctx:      c,
		cancel:   cancel,
		endpoint: endpoint,
		database: database,
		interval: interval,
		gather: r.NewArrowGather(
			linmetric.WithReadRuntimeOption(newRuntimeObserver(r)),
			linmetric.WithGlobalKeyValueOption(globalKeyValues),
		),
		builder: lmetrics.NewMetricBuilder(memory.NewGoAllocator()),
		client:  &http.Client{Timeout: pushTimeout},
	}

	monitorScope := r.NewScope("lindb.monitor")
	arrowPusherScope := monitorScope.Scope("arrow_pusher")
	p.statistics.pushBytesCounter = arrowPusherScope.NewCounter("push_bytes")
	p.statistics.pushMetricsCounter = arrowPusherScope.NewCounter("push_metrics_count")
	p.statistics.pushErrorCounter = arrowPusherScope.NewCounter("push_error_count")

	return p
}

// Start runs the push loop until Stop() is called.
func (p *arrowNativePusher) Start() {
	arrowPushLogger.Info("arrow native pusher starting...")
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.gatherAndBuild()
		case <-p.ctx.Done():
			arrowPushLogger.Info("arrow native pusher stopped")
			return
		}
	}
}

// Stop cancels the push loop.
func (p *arrowNativePusher) Stop() {
	p.cancel()
}

// gatherAndBuild collects metrics, appends them to the MetricBuilder,
// serializes to Arrow IPC bytes, then pushes over HTTP.
func (p *arrowNativePusher) gatherAndBuild() {
	metrics := p.gather.GatherMetrics()
	if len(metrics) == 0 {
		return
	}

	for _, m := range metrics {
		p.builder.Append(m)
	}

	data, err := p.builder.Bytes()
	if err != nil {
		p.statistics.pushErrorCounter.Incr()
		arrowPushLogger.Error("failed to serialize Arrow IPC bytes", logger.Error(err))
		return
	}

	p.statistics.pushMetricsCounter.Add(float64(len(metrics)))
	p.statistics.pushBytesCounter.Add(float64(len(data)))
	p.push(data)
}

// push sends the Arrow IPC bytes to the configured endpoint via HTTP PUT.
// The target database is specified via the X-LinDB-Database header.
func (p *arrowNativePusher) push(data []byte) {
	if len(data) == 0 {
		return
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPut, p.endpoint, bytes.NewReader(data))
	if err != nil {
		p.statistics.pushErrorCounter.Incr()
		arrowPushLogger.Error("failed to create HTTP request", logger.Error(err))
		return
	}
	req.Header.Set("Content-Type", constants.ContentTypeArrow)
	req.Header.Set(constants.DatabaseHeader, p.database)

	resp, err := p.client.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		p.statistics.pushErrorCounter.Incr()
		arrowPushLogger.Error("failed to push Arrow metrics", logger.Error(err))
	}
}
