package models

//go:generate mockgen -source ./point.go -destination=./point_mock.go -package models

// Point contains the methods for accessing a point.
type Point interface {
	Name() string
	Timestamp() int64
	Tags() map[string]string
	Fields() map[string]Field

	// sort tags by key in ascii ascending order, then concat each key and value.
	TagsID() string
	TsID() uint32
}

// Field is the numerical key-value pair of metric.
type Field interface {
}
