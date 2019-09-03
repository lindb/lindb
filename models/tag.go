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
	tagKeyValues := make([]string, 0, len(tags))

	totalLen := 0
	for key, val := range tags {
		keyVal := key + "=" + val
		tagKeyValues = append(tagKeyValues, keyVal)
		totalLen += len(keyVal)
	}

	sort.Strings(tagKeyValues)

	return strings.Join(tagKeyValues, ",")
}
