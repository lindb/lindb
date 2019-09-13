package series

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	v := NewVersion()
	assert.True(t, v.Elapsed().Seconds() < float64(1))

	v1 := Version(0)
	assert.True(t, v1.IsExpired(time.Hour))

	v2 := Version(1)
	assert.True(t, v1.Before(v2))
	assert.False(t, v1.After(v2))
	assert.False(t, v1.Equal(v2))

	t.Log(v2.String())
}
