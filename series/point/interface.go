package point

import (
	"fmt"
	"time"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

// point abstracts the abilities of the metric Point
// it's the interface that Point should implements
type point interface {
	// Key returns the key of the point
	// metric-name + sorted tags
	Key() []byte

	// Name returns the metric-name
	Name() []byte

	// SetName updates the metric-name for this point
	SetName(name string) *Point

	// Tags returns the tag list of the point
	Tags() tag.Tags

	// ForEachTag iterates over each tag pair,
	// if f returns false, iteration stops
	ForEachTag(f func(key, value []byte) bool)

	// AddTag adds or replaces a tag value for a point.
	AddTag(key, value string) *Point

	// AddTags adds or replaces multi tag value for a point.
	// Each data will be dropped if given an odd number of arguments.
	AddTags(kv ...string) *Point

	// SetTags replaces the tags for the point.
	SetTags(tags tag.Tags) *Point

	// TagsHashKey returns the concated tags
	// hash-key: host=test,ip=1.1.1.1,zone=sh
	TagsHashKey() []byte

	// HasTag returns true if the tag exists for the point.
	HasTag(tag []byte) bool

	// HashID returns the xxhash value for metric-name
	HashID() uint64

	// TagsHashID returns the xxhash for concated tags
	TagsHashID() uint64

	// Time returns the timestamp for the metric
	Time() time.Time

	// SetTime updates the timestamp for the point.
	SetTime(t time.Time) *Point

	// UnixMilli returns the timestamp as millisecond
	UnixMilli() int64

	// SetUnixMilli updates the timestamp from milliseconds for the point
	SetUnixMilli(timestamp int64) *Point

	fmt.Stringer

	// StringSize returns the length of String()
	StringSize() int

	// AppendString appends the string representation of the point to buf.
	AppendString(buf []byte) []byte

	// Fields returns the fields for the point.
	Fields() (field.Fields, error)

	// FieldIterator returns a FieldIterator that can be used to traverse the
	// fields of a point without constructing the in-memory map.
	FieldIterator() *FieldIterator

	// Reset clears all internal data,
	Reset()

	// AddField adds or replaces a field pair for a point.
	AddField(name string, fType field.Type, value interface{}) *Point

	// SetFields replaces the fields for the point.
	SetFields(fs field.Fields) *Point
}

// fieldIterator provides a low-allocation interface to iterate through a point's fields.
type fieldIterator interface {
	// Next indicates whether there any fields remaining.
	Next() bool

	// Name returns the key of the current field.
	Name() []byte

	// Type returns the FieldType of the current field.
	Type() field.Type

	// Int64Value returns the integer value of the current field.
	Int64Value() (int64, error)

	// Float64Value returns the float value of the current field.
	Float64Value() (float64, error)

	// Reset resets the iterator to its initial state.
	Reset(data []byte)
}

var (
	// static type assertion
	_ point         = (*Point)(nil)
	_ fieldIterator = (*FieldIterator)(nil)
)
