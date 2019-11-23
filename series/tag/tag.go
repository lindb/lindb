package tag

import (
	"sort"
	"strings"
)

const emptyStr = ""

var emptyArray []string

// Concat concats map-tags to string
func Concat(tags map[string]string) string {
	if tags == nil {
		return emptyStr
	}
	tagKeys := make([]string, 0, len(tags))
	var b strings.Builder
	b.Grow(128)
	for key := range tags {
		tagKeys = append(tagKeys, key)
	}
	sort.Strings(tagKeys)
	tagKeysLen := len(tagKeys)
	for idx, tagKey := range tagKeys {
		b.WriteString(tagKey)
		b.WriteString("=")
		b.WriteString(tags[tagKey])
		if idx != tagKeysLen-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}

// ConcatTagValues cancats the tag values to string
func ConcatTagValues(tagValues []string) string {
	if len(tagValues) == 0 {
		return emptyStr
	}
	return strings.Join(tagValues, ",")
}

// SplitTagValues splits the string of tag values to array
func SplitTagValues(tags string) []string {
	if tags == "" {
		return emptyArray
	}
	return strings.Split(tags, ",")
}
