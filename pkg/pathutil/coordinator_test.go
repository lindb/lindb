package pathutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetName(t *testing.T) {
	assert.Equal(t, "name", GetName("/test/name"))
	assert.Equal(t, "name", GetName("name"))
}
