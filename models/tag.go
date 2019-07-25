package models

import "strings"

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
