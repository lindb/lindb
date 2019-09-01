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

// NewTags returns a Tag list from string.
func NewTags(tagStr string) (theTags []Tag) {
	pairs := strings.Split(tagStr, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}
		theTags = append(theTags, Tag{Key: kv[0], Value: kv[1]})
	}
	return
}

// convert tags to string
func TagsAsString(tags map[string]string) string {
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
