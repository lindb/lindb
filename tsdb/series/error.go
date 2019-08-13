package series

import "errors"

// ErrMetaDataNotExist is returned by index-database when the meta-data of metric not exists
var ErrMetaDataNotExist = errors.New("meta data not exist")

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
