package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTags(t *testing.T) {
	assert.Nil(t, NewTags(""))
	assert.Len(t, NewTags("host=alpha,ip=1.1.1.1"), 2)
}
