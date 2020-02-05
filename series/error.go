package series

import "errors"

// ErrTooManyTags is the error returned by tsdb when
// writes exceed the max limit of tag identifiers.
var ErrTooManyTags = errors.New("too many tags")

// ErrTooManyTagKeys is the error returned by tsdb when
// writes exceed the max limit of tag keys.
var ErrTooManyTagKeys = errors.New("too many tag keys")

// ErrTooManyFields is the error returned by tsdb when
// writes exceed the max limit of fields.
var ErrTooManyFields = errors.New("too many fields")

// ErrWrongFieldType is the error returned by tsdb when
// field-type of new point is different from the type before.
var ErrWrongFieldType = errors.New("field type is wrong")

// ErrResetVersionUnavailable is the error returned by tsdb when
// the immutable tagIndex has not been flushed yet.
var ErrResetVersionUnavailable = errors.New("reset version unavailable")
