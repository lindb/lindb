package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSegmentStore(t *testing.T) {
	sStore := newSegmentStore(0)
	assert.NotNil(t, sStore)
}
