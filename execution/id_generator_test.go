package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestIDGenerator(t *testing.T) {
	gen := NewRequestIDGenerator("test-node")
	assert.NotEqual(t, gen.GenerateRequestID(), gen.GenerateRequestID())
}
