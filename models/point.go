package models

type Point interface {
	Name() string
	Timestamp() int64
	Tags() string
	Fields() map[string]Field
}

type Field interface {
}
