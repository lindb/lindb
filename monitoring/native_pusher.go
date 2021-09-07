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
	"io"
	"net/http"
	"time"

	"github.com/klauspost/compress/gzip"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/logger"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"
)

var nativePushLogger = logger.GetLogger("monitoring", "Pusher")

var (
	monitorScope       = linmetric.NewScope("lindb.monitor")
	nativePusherScope  = monitorScope.Scope("native_pusher")
	pushBytesCounter   = nativePusherScope.NewCounter("push_bytes")
	pushMetricsCounter = nativePusherScope.NewCounter("push_metrics_count")
	pushErrorCounter   = nativePusherScope.NewCounter("push_error_count")
)

const (
	ProtoType     = `application/protobuf`
	ProtoProtocol = `io.lindb.proto.Metric`
	ProtoFmt      = ProtoType + "; proto=" + ProtoProtocol + ";"
)

// NativePusher collects metrics from internal lin-metric registry,
// then pushes metrics data via http.
type NativePusher interface {
	// Start starts push metrics data in period
	Start()
	// Stop stops push metrics data
	Stop()
}

// nativeProtoPusher writes native protobuf data to ingestion endpoint.
type nativeProtoPusher struct {
	ctx             context.Context
	cancel          context.CancelFunc
	interval        time.Duration
	endpoint        string // HTTP endpoint
	globalKeyValues tag.KeyValues
	gather          linmetric.Gather
	client          *http.Client
}

// NewNativeProtoPusher creates a new native pusher
func NewNativeProtoPusher(
	ctx context.Context,
	endpoint string,
	interval time.Duration,
	pushTimeout time.Duration,
	globalKeyValues tag.KeyValues,
) NativePusher {
	c, cancel := context.WithCancel(ctx)
	return &nativeProtoPusher{
		ctx:             c,
		cancel:          cancel,
		endpoint:        endpoint,
		interval:        interval,
		globalKeyValues: globalKeyValues,
		gather: linmetric.NewGather(
			linmetric.WithReadRuntimeOption(),
			linmetric.WithGlobalKeyValueOption(globalKeyValues),
		),
		client: &http.Client{Timeout: pushTimeout},
	}
}

func (np *nativeProtoPusher) Start() {
	nativePushLogger.Info("native proto pusher starting...")
	ticker := time.NewTicker(np.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			buf := np.gatherAndMarshal()
			np.push(buf)
		case <-np.ctx.Done():
			nativePushLogger.Info("native proto pusher stopped")
			return
		}
	}
}

func (np *nativeProtoPusher) Stop() {
	np.cancel()
}

func (np *nativeProtoPusher) gatherAndMarshal() *bytes.Buffer {
	metrics := np.gather.Gather()
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	ml := protoMetricsV1.MetricList{Metrics: metrics}
	data, err := ml.Marshal()
	if err != nil {
		pushErrorCounter.Add(float64(len(metrics)))
		nativePushLogger.Error("failed to marshal metric", logger.Error(err))
		return nil
	}
	_, _ = gzipWriter.Write(data)
	pushMetricsCounter.Add(float64(len(metrics)))
	_ = gzipWriter.Close()
	pushBytesCounter.Add(float64(buf.Len()))
	return &buf
}

func (np *nativeProtoPusher) push(r io.Reader) {
	if r == nil {
		return
	}
	req, _ := http.NewRequest(http.MethodPut, np.endpoint, r)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", ProtoFmt)

	resp, err := np.client.Do(req)
	if err != nil {
		pushErrorCounter.Incr()
		nativePushLogger.Error("failed to post request", logger.Error(err))
		return
	}
	_ = resp.Body.Close()
}
