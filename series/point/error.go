package point

import "fmt"

var (
	ErrInvalidPoint      = fmt.Errorf("point is invalid")
	ErrMissingMetricName = fmt.Errorf("metric name is missing")
	ErrDuplicateTags     = fmt.Errorf("duplicat tags")
	ErrMissingTagValue   = fmt.Errorf("tag value is missing")
	ErrMissingFields     = fmt.Errorf("fields is missing")
	ErrMissingFieldName  = fmt.Errorf("field name is missing")
	ErrMissingFieldValue = fmt.Errorf("field value is missing")
	ErrInvalidNumber     = fmt.Errorf("invalid number")
)
