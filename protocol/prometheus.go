package protocol

import (
	"bytes"

	"github.com/cespare/xxhash"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

// PromParse parses prometheus text protocol to LinDB pb protocol.
func PromParse(data []byte) (*pb.MetricList, error) {
	parser := &expfmt.TextParser{}
	out, err := parser.TextToMetricFamilies(bytes.NewBuffer(data))
	if err != nil && len(out) == 0 {
		return nil, err
	}
	metricList := &pb.MetricList{}
	for name, pm := range out {
		metricType := *pm.Type
		if metricType == dto.MetricType_UNTYPED {
			// not support untyped metric type
			continue
		}
		for _, m := range pm.Metric {
			f := getFieldType(metricType, m)
			if f == nil {
				continue
			}

			metric := &pb.Metric{Name: name}
			metric.Fields = []*pb.Field{f}
			if m.TimestampMs != nil {
				metric.Timestamp = *m.TimestampMs
			} else {
				metric.Timestamp = timeutil.Now()
			}
			tagCount := len(m.Label)
			if tagCount > 0 {
				tags := make(map[string]string, tagCount)
				for _, label := range m.Label {
					tags[*label.Name] = *label.Value
				}
				metric.TagsHash = xxhash.Sum64String(tag.Concat(tags))
				metric.Tags = tags
			} else {
				metric.TagsHash = xxhash.Sum64String(metric.Name)
			}

			metricList.Metrics = append(metricList.Metrics, metric)
		}
	}
	return metricList, nil
}

func getFieldType(metricType dto.MetricType, metric *dto.Metric) *pb.Field {
	switch metricType {
	case dto.MetricType_COUNTER:
		if metric.Counter != nil && metric.Counter.Value != nil {
			return &pb.Field{
				Name:   "counter",
				Type:   pb.FieldType_Increase,
				Fields: []*pb.PrimitiveField{{PrimitiveID: int32(field.SimpleFieldPFieldID), Value: *metric.Counter.Value}},
			}
		}
	case dto.MetricType_GAUGE:
		if metric.Gauge != nil && metric.Gauge.Value != nil {
			return &pb.Field{
				Name:   "gauge",
				Type:   pb.FieldType_Gauge,
				Fields: []*pb.PrimitiveField{{PrimitiveID: int32(field.SimpleFieldPFieldID), Value: *metric.Gauge.Value}},
			}
		}
	case dto.MetricType_SUMMARY:
		if metric.Summary == nil || metric.Summary.SampleCount == nil || metric.Summary.SampleSum == nil {
			return nil
		}
		f := &pb.Field{
			Name: "summary",
			Type: pb.FieldType_Summary,
		}
		count := float64(*metric.Summary.SampleCount)
		f.Fields = append(f.Fields, &pb.PrimitiveField{
			PrimitiveID: int32(1),
			Value:       *metric.Summary.SampleSum,
		}, &pb.PrimitiveField{
			PrimitiveID: int32(2),
			Value:       count,
		})
		quantile := metric.Summary.Quantile
		for _, q := range quantile {
			switch *q.Quantile {
			case 0.5:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(50),
					Value:       *q.Value * count,
				})
			case 0.75:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(75),
					Value:       *q.Value * count,
				})
			case 0.90:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(90),
					Value:       *q.Value * count,
				})
			case 0.95:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(95),
					Value:       *q.Value * count,
				})
			case 0.99:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(99),
					Value:       *q.Value * count,
				})
			case 0.999:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(39),
					Value:       *q.Value * count,
				})
			case 0.9999:
				f.Fields = append(f.Fields, &pb.PrimitiveField{
					PrimitiveID: int32(49),
					Value:       *q.Value * count,
				})
			}
		}
		return f
	}
	return nil
}
