package series

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	v := NewVersion()
	assert.True(t, v.Elapsed().Seconds() < float64(1))
}
