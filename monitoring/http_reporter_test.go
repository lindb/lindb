package monitoring

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockBuckets struct{}

func (b *mockBuckets) String() string               { return "" }
func (b *mockBuckets) Len() int                     { return 0 }
func (b *mockBuckets) Less(i, j int) bool           { return true }
func (b *mockBuckets) Swap(i, j int)                {}
func (b *mockBuckets) AsValues() []float64          { return nil }
func (b *mockBuckets) AsDurations() []time.Duration { return nil }

func Test_NewHTTPReporter(t *testing.T) {
	reporter := NewHTTPReporter("")
	reporter.ReportTimer("", nil, time.Second)
	assert.True(t, reporter.Capabilities().Reporting())
	assert.False(t, reporter.Capabilities().Tagging())

	reporter.ReportHistogramValueSamples("", nil, &mockBuckets{}, 1, 1, 1)
	reporter.ReportHistogramDurationSamples("", nil, &mockBuckets{}, time.Second, time.Second, 1)
}
