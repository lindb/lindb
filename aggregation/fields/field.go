package fields

import (
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/collections"
)

//go:generate mockgen -source=./field.go -destination=./field_mock.go -package=fields

// Field represents the field series for the time series
type Field interface {
	// GetValues returns the values which function call need by given function type
	GetValues(funcType function.FuncType) []collections.FloatArray
	// GetDefaultValues returns the field default values which aggregation need if user not input function type
	GetDefaultValues() []collections.FloatArray
}
