package models

import "github.com/eleme/lindb/pkg/field"

//go:generate mockgen -source ./point.go -destination=./point_mock.go -package models

// Point contains the methods for accessing a point.
type Point interface {
	Name() string
	Timestamp() int64
	// Tags sorts keys in ascii ascending order, then concat each key and value.
	// example: ezone=nj(delim)host=alpha-1.vm(delim)ip=1.1.1.1(delim)
	Tags() string
	Fields() map[string]Field
	TagsMap() map[string]string
}

// Field is the numerical key-value pair of metric.
type Field interface {
	Type() field.Type
	IsComplex() bool
}

// SimpleField is simple field for single value
type SimpleField interface {
	Field
	ValueType() field.ValueType
	AggType() field.AggType
	Value() interface{}
}
