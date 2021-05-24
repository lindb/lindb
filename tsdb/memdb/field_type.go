package memdb

import (
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

type FieldKey uint32

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

// buildFieldKey returns field key with family/field id
func buildFieldKey(familyID familyID, fieldID field.ID) FieldKey {
	return FieldKey(uint32(fieldID) | uint32(familyID)<<16)
}
