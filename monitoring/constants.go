package monitoring

// DefaultHistogramBuckets represents default prometheus histogram buckets in LinDB
var DefaultHistogramBuckets = []float64{
	0.0, 10.0, 25.0, 50.0, 75.0,
	100.0, 200.0, 300.0, 400.0, 500.0, 600.0, 800.0,
	1000.0, 2000.0, 5000.0,
}
