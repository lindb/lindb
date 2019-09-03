package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagsAsString(t *testing.T) {
	assert.Equal(t, "", TagsAsString(nil))
	tags := map[string]string{"t2": "v2", "t1": "v1"}
	tagsStr := TagsAsString(tags)
	assert.Equal(t, "t1=v1,t2=v2", tagsStr)
}
