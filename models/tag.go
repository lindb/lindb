package models

import (
	"sort"
	"strings"
)

// Tag is the key/value tag pair of a metric point.
type Tag struct {
	// tag-key
	Key string
	// tag-value
	Value string
}

// convert tags to string
func TagsAsString(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	tagKeyValues := make([]string, 0, len(tags))

	totalLen := 0
	for key, val := range tags {
		keyVal := key + val
		tagKeyValues = append(tagKeyValues, keyVal)
		totalLen += len(keyVal)
	}

	sort.Strings(tagKeyValues)

	var builder strings.Builder
	builder.Grow(totalLen)

	for _, tagKeyValue := range tagKeyValues {
		builder.WriteString(tagKeyValue)
	}

	return builder.String()
}
