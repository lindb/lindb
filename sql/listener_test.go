package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListener_build(t *testing.T) {
	l := &listener{}
	s, err := l.statement()
	assert.NoError(t, err)
	assert.Nil(t, s)
}
