package memdb

import (
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

// getFieldType return field type by given field
func getFieldType(f *pb.Field) field.Type {
	switch f.Type {
	case pb.FieldType_Sum:
		return field.SumField
	case pb.FieldType_Max:
		return field.MaxField
	case pb.FieldType_Min:
		return field.MinField
	case pb.FieldType_Gauge:
		return field.GaugeField
	default:
		return field.Unknown
	}
}
