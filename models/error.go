package models

import (
	"errors"
)

// ErrTooManyTags is the error returned by memory-database when
// writes exceed the max limit of tag identifiers.
var ErrTooManyTags = errors.New("too many tags")
