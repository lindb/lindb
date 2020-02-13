package memdb

import (
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
)

// getFieldType return field type by given field
func getFieldType(f *pb.Field) field.Type {
	switch f.Field.(type) {
	case *pb.Field_Sum:
		return field.SumField
	case *pb.Field_Max:
		return field.MaxField
	case *pb.Field_Min:
		return field.MinField
	case *pb.Field_Gauge:
		return field.GaugeField
	case *pb.Field_Summary:
		return field.SummaryField
	default:
		return field.Unknown
	}
}

// buildFieldKeys returns field key with family/field/primitive
func buildFieldKey(familyID familyID, fieldID field.ID, primitiveID field.PrimitiveID) uint32 {
	return uint32(primitiveID) | uint32(fieldID)<<8 | uint32(familyID)<<16
}
