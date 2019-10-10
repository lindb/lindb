package fields

import (
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./field.go -destination=./field_mock.go -package=fields

// Field represents the field series for the time series
type Field interface {
	SetValue(fieldSeries series.Iterator)
	// GetValues returns the values which function call need by given function type
	GetValues(funcType function.FuncType) (result []collections.FloatArray)
	// GetDefaultValues returns the field default values which aggregation need if user not input function type
	GetDefaultValues() (result []collections.FloatArray)
	// Reset resets field's value for reusing
	Reset()
}

// dynamicField represents the dynamic field for storing multi-primitive fields
type dynamicField struct {
	fieldType field.Type
	startTime int64
	interval  int64
	capacity  int

	fields map[uint16]collections.FloatArray
}

// NewDynamicField creates a dynamic field series
func NewDynamicField(fieldType field.Type, startTime int64, interval int64, capacity int) Field {
	return &dynamicField{
		fieldType: fieldType,
		startTime: startTime,
		interval:  interval,
		capacity:  capacity,
		fields:    make(map[uint16]collections.FloatArray),
	}
}

// SetValue sets the field's value by time slot
func (f *dynamicField) SetValue(fieldSeries series.Iterator) {
	if fieldSeries == nil {
		return
	}
	var pField collections.FloatArray
	ok := false
	for fieldSeries.HasNext() {
		startTime, it := fieldSeries.Next()
		if it == nil {
			continue
		}
		for it.HasNext() {
			pIt := it.Next()
			pID := pIt.FieldID()
			pField, ok = f.fields[pID]
			if !ok {
				pField = collections.NewFloatArray(f.capacity)
				f.fields[pID] = pField
			}
			for pIt.HasNext() {
				slot, val := pIt.Next()
				idx := ((int64(slot)*f.interval + startTime) - f.startTime) / f.interval
				pField.SetValue(int(idx), val)
			}
		}
	}
}

// GetValues returns the values which function call need by given function type and field type
func (f *dynamicField) GetValues(funcType function.FuncType) (result []collections.FloatArray) {
	pFields := f.fieldType.GetPrimitiveFields(funcType)
	return f.getFieldValues(pFields)
}

// GetDefaultValues returns the field default values which aggregation need by field type
func (f *dynamicField) GetDefaultValues() []collections.FloatArray {
	return f.getFieldValues(f.fieldType.GetDefaultPrimitiveFields())
}

func (f *dynamicField) Reset() {
	for _, pField := range f.fields {
		pField.Reset()
	}
}

// getFieldValues returns the values by primitive field ids
func (f *dynamicField) getFieldValues(pFields map[uint16]field.AggType) (result []collections.FloatArray) {
	if len(pFields) == 0 {
		return
	}
	for pID := range pFields {
		pField, ok := f.fields[pID]
		if ok {
			result = append(result, pField)
		}
	}
	return
}
