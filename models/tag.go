package models

import (
	"sort"
	"strings"
)

// convert tags to string
func TagsAsString(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	tagKeys := make([]string, 0, len(tags))
	var b strings.Builder
	b.Grow(128)
	for key := range tags {
		tagKeys = append(tagKeys, key)
	}
	sort.Strings(tagKeys)
	for idx, tagKey := range tagKeys {
		b.WriteString(tagKey)
		b.WriteString("=")
		b.WriteString(tags[tagKey])
		if idx != len(tagKeys)-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}
