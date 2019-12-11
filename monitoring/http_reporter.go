package monitoring

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/rpc/proto/field"

	"github.com/uber-go/tally"
)

// internal metrics will be stored in this database
//const internalDatabase = "_internal"

// httpReporter implements tally.StatsReporter
type httpReporter struct {
	endpoint string          // HTTP endpoint
	metrics  []*field.Metric // buffer
	mux      sync.Mutex      // mutex for metrics
}

func NewHTTPReporter(endpoint string) tally.StatsReporter {
	return &httpReporter{endpoint: endpoint}
}

// Capabilities returns the capabilities description of the reporter.
func (ir *httpReporter) Capabilities() tally.Capabilities {
	return ir
}

func (ir *httpReporter) Reporting() bool {
	return true
}

func (ir *httpReporter) Tagging() bool {
	return false
}

// Flush asks the reporter to flush all reported values.
func (ir *httpReporter) Flush() {
	ir.mux.Lock()
	if len(ir.metrics) == 0 {
		ir.mux.Unlock()
		return
	}

	data := encoding.JSONMarshal(field.MetricList{
		Metrics: ir.metrics,
	})
	ir.metrics = ir.metrics[:0]
	ir.mux.Unlock()

	//TODO set database name
	resp, err := http.Post(ir.endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Error("post error: %s", logger.Error(err))
		return
	}
	resp.Body.Close()
}

// ReportCounter reports a counter value
func (ir *httpReporter) ReportCounter(
	name string,
	tags map[string]string,
	value int64,
) {
	newMetric := &field.Metric{
		Name:      name,
		Timestamp: timeutil.Now(),
		Tags:      tags,
		Fields: []*field.Field{{
			Name: "count",
			Field: &field.Field_Sum{
				Sum: &field.Sum{
					Value: float64(value),
				},
			},
		}}}

	ir.mux.Lock()
	defer ir.mux.Unlock()
	ir.metrics = append(ir.metrics, newMetric)
}

// ReportGauge reports a gauge value
func (ir *httpReporter) ReportGauge(
	name string,
	tags map[string]string,
	value float64,
) {
	newMetric := &field.Metric{
		Name:      name,
		Timestamp: timeutil.Now(),
		Tags:      tags,
		Fields: []*field.Field{{
			Name: "gauge",
			Field: &field.Field_Gauge{
				Gauge: &field.Gauge{
					Value: value,
				},
			},
		}}}

	ir.mux.Lock()
	defer ir.mux.Unlock()
	ir.metrics = append(ir.metrics, newMetric)
}

// ReportTimer reports a timer value
func (ir *httpReporter) ReportTimer(
	name string,
	tags map[string]string,
	interval time.Duration,
) {
	// todo: @codingcrush, implement this later
}

// ReportHistogramValueSamples reports histogram samples for a bucket
func (ir *httpReporter) ReportHistogramValueSamples(
	name string,
	tags map[string]string,
	buckets tally.Buckets,
	bucketLowerBound,
	bucketUpperBound float64,
	samples int64,
) {
	// todo: @codingcrush, implement this later
}

// ReportHistogramDurationSamples reports histogram samples for a bucket
func (ir *httpReporter) ReportHistogramDurationSamples(
	name string,
	tags map[string]string,
	buckets tally.Buckets,
	bucketLowerBound,
	bucketUpperBound time.Duration,
	samples int64,
) {
	// todo: @codingcrush, implement this later
}
