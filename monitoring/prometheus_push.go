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

	"github.com/OneOfOne/xxhash"
	"github.com/gogo/protobuf/proto"
	"github.com/klauspost/compress/gzip"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"

	"github.com/lindb/lindb/pkg/logger"
)

// for testing
var (
	newRequest = http.NewRequest
	doRequest  = http.DefaultClient.Do
)

var pushLogger = logger.GetLogger("monitoring", "Pusher")

var separatorByteSlice = []byte{model.SeparatorByte} // For convenient use with xxhash.
const (
	contentTypeHeader     = "Content-Type"
	contentEncodingHeader = "Content-Encoding"
)

// PrometheusPusher represents a pusher,
// collects metrics from prometheus registry, then pushes metrics data via http.
// Counter/Summary/Histogram need calc delta value with previous.
type PrometheusPusher interface {
	// Start starts push metrics data in period
	Start()
	// Stop stops push metrics data
	Stop()
}

// metricFamily represents the metric family data that stores series data
type metricFamily struct {
	metrics map[uint64]*dto.Metric
}

// newMetricFamily creates a new metric family
func newMetricFamily() *metricFamily {
	return &metricFamily{
		metrics: make(map[uint64]*dto.Metric),
	}
}

// prometheusPusher implements PrometheusPusher interface
type prometheusPusher struct {
	ctx          context.Context
	cancel       context.CancelFunc
	interval     time.Duration
	endpoint     string // HTTP endpoint
	gatherers    prometheus.Gatherers
	globalLabels []*dto.LabelPair
	expfmt       expfmt.Format

	newRequest func(method, url string, body io.Reader) (*http.Request, error)
	doRequest  func(req *http.Request) (*http.Response, error)

	hash       *xxhash.XXHash64
	gatherFunc func(gatherers prometheus.Gatherers) ([]*dto.MetricFamily, error)
	encodeFunc func(enc expfmt.Encoder, mf *dto.MetricFamily) error

	metricFamilies map[string]*metricFamily
}

// NewPrometheusPusher creates a new prometheus pusher
func NewPrometheusPusher(
	ctx context.Context,
	endpoint string,
	interval time.Duration,
	gatherers prometheus.Gatherers,
	globalLabels []*dto.LabelPair,
) PrometheusPusher {
	c, cancel := context.WithCancel(ctx)
	return &prometheusPusher{
		ctx:            c,
		cancel:         cancel,
		endpoint:       endpoint,
		interval:       interval,
		gatherers:      gatherers,
		globalLabels:   globalLabels,
		doRequest:      doRequest,
		newRequest:     newRequest,
		expfmt:         expfmt.FmtText,
		hash:           xxhash.New64(),
		metricFamilies: make(map[string]*metricFamily),
		encodeFunc:     encode,
		gatherFunc:     gather,
	}
}

// Start starts push metrics data in period
func (p *prometheusPusher) Start() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.run()
		case <-p.ctx.Done():
			pushLogger.Info("stop prometheus pusher")
			return
		}
	}
}

// Stop stops push metrics data
func (p *prometheusPusher) Stop() {
	p.cancel()
}

// run collects the metrics from prometheus registry, calc delta value with previous value for counter/summary/histogram,
// finally cache previous data and sends current metrics data.
func (p *prometheusPusher) run() {
	// 1. get prometheus metrics
	mfs, err := p.gatherFunc(p.gatherers)
	if err != nil {
		pushLogger.Error("get metric fail from prometheus registry", logger.Error(err))
		return
	}
	// 2. encode metric, calc delta value if need calc delta
	buf := &bytes.Buffer{}
	var enc expfmt.Encoder

	gzipWriter := gzip.NewWriter(buf)
	enc = expfmt.NewEncoder(gzipWriter, p.expfmt)

	for _, mf := range mfs {
		if p.needCalcDelta(mf.GetType()) {
			p.metricFamilies[mf.GetName()] = p.calcDelta(mf)
		}
		// add global labels
		for _, m := range mf.GetMetric() {
			m.Label = append(m.Label, p.globalLabels...)
		}
		if err = p.encodeFunc(enc, mf); err != nil {
			pushLogger.Error("encode prometheus metric error", logger.Error(err))
		}
	}
	_ = gzipWriter.Close()

	// 3. new metric write request
	req, err := p.newRequest("PUT", p.endpoint, buf)
	if err != nil {
		pushLogger.Error("new write monitoring request error", logger.Error(err))
		return
	}
	req.Header.Set(contentEncodingHeader, "gzip")
	req.Header.Add(contentTypeHeader, string(p.expfmt))

	// 4. send metric data
	writeResp, err := p.doRequest(req)
	if err != nil {
		pushLogger.Error("write monitoring data error", logger.Error(err))
		return
	}
	_ = writeResp.Body.Close()
}

// needCalcDelta checks if need calc delta value
func (p *prometheusPusher) needCalcDelta(metricType dto.MetricType) bool {
	return metricType == dto.MetricType_COUNTER ||
		metricType == dto.MetricType_SUMMARY ||
		metricType == dto.MetricType_HISTOGRAM
}

// calcDelta calculates delta value with previous
func (p *prometheusPusher) calcDelta(mf *dto.MetricFamily) *metricFamily {
	metricName := mf.GetName()
	newFamily := newMetricFamily()
	family, familyExist := p.metricFamilies[metricName]
	if !familyExist {
		p.metricFamilies[metricName] = family
	}
	for _, m := range mf.GetMetric() {
		p.hash.Reset()
		for _, l := range m.GetLabel() {
			_, _ = p.hash.WriteString(l.GetName())
			_, _ = p.hash.Write(separatorByteSlice)
			_, _ = p.hash.WriteString(l.GetValue())
			_, _ = p.hash.Write(separatorByteSlice)
		}
		hash := p.hash.Sum64()
		previous := m
		// if family exist in previous cache
		if familyExist {
			oldMetric, ok := family.metrics[hash]
			if ok {
				switch mf.GetType() {
				case dto.MetricType_COUNTER:
					previous = p.calcDeltaCounter(m, oldMetric)
				case dto.MetricType_SUMMARY:
					previous = p.calcDeltaSummary(m, oldMetric)
				case dto.MetricType_HISTOGRAM:
					previous = p.calcDeltaHistogram(m, oldMetric)
				}
			}
		}
		newFamily.metrics[hash] = previous
	}
	return newFamily
}

// calcDeltaCounter calculates delta counter
func (p *prometheusPusher) calcDeltaCounter(new, old *dto.Metric) *dto.Metric {
	metric := &dto.Metric{
		Label: new.Label,
		Counter: &dto.Counter{
			Value: new.Counter.Value,
		},
		TimestampMs: new.TimestampMs,
	}

	new.Counter.Value = proto.Float64(new.Counter.GetValue() - old.Counter.GetValue())
	return metric
}

// calcDeltaSummary calculates delta summary
func (p *prometheusPusher) calcDeltaSummary(new, old *dto.Metric) *dto.Metric {
	metric := &dto.Metric{
		Label: new.Label,
		Summary: &dto.Summary{
			SampleCount: new.Summary.SampleCount,
			SampleSum:   new.Summary.SampleSum,
			Quantile:    new.Summary.Quantile,
		},
		TimestampMs: new.TimestampMs,
	}

	new.Summary.SampleSum = proto.Float64(new.Summary.GetSampleSum() - old.Summary.GetSampleSum())
	new.Summary.SampleCount = proto.Uint64(new.Summary.GetSampleCount() - old.Summary.GetSampleCount())

	return metric
}

// calcDeltaHistogram calculates delta histogram
func (p *prometheusPusher) calcDeltaHistogram(new, old *dto.Metric) *dto.Metric {
	// clone from new metric
	metric := &dto.Metric{
		Label: new.Label,
		Histogram: &dto.Histogram{
			SampleCount: new.Histogram.SampleCount,
			SampleSum:   new.Histogram.SampleSum,
			Bucket:      make([]*dto.Bucket, len(new.Histogram.Bucket)),
		},
		TimestampMs: new.TimestampMs,
	}

	new.Histogram.SampleSum = proto.Float64(new.Histogram.GetSampleSum() - old.Histogram.GetSampleSum())
	new.Histogram.SampleCount = proto.Uint64(new.Histogram.GetSampleCount() - old.Histogram.GetSampleCount())

	buckets := make(map[float64]uint64)

	for _, b := range old.Histogram.Bucket {
		buckets[b.GetUpperBound()] = b.GetCumulativeCount()
	}
	// calc delta for bucket
	for idx, b := range new.Histogram.Bucket {
		metric.Histogram.Bucket[idx] = b
		bucket, ok := buckets[b.GetUpperBound()]
		if ok {
			b.CumulativeCount = proto.Uint64(b.GetCumulativeCount() - bucket)
		}
	}
	return metric
}

// Gather collects prometheus metrics from registry
func gather(gatherers prometheus.Gatherers) ([]*dto.MetricFamily, error) {
	return gatherers.Gather()
}

// encode encodes metric families into an underlying wire prometheus.
func encode(enc expfmt.Encoder, mf *dto.MetricFamily) error {
	return enc.Encode(mf)
}
