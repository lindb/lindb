package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTags(t *testing.T) {
	assert.Nil(t, NewTags(""))
	assert.Len(t, NewTags("host=alpha,ip=1.1.1.1"), 2)
}

func TestTagsAsString(t *testing.T) {
	assert.Equal(t, "", TagsAsString(map[string]string{}))

	tagsStr := TagsAsString(map[string]string{"t2": "v2", "t1": "v1"})

	assert.Equal(t, "t1v1t2v2", tagsStr)
}
