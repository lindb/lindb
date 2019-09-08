package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagsAsString(t *testing.T) {
	assert.Equal(t, "", TagsAsString(nil))

	assert.Equal(t, "", TagsAsString(map[string]string{}))

	tagsStr := TagsAsString(map[string]string{"t2": "v2", "t1": "v1"})

	assert.Equal(t, "t1v1t2v2", tagsStr)
}
