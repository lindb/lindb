package point

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnescapeMetricName_EscapeMetricName(t *testing.T) {
	pairs := []struct {
		in  string
		out string
	}{
		{"xxx", "xxx"},
		{"xx x", "xx\\ x"},
		{"xx,x", "xx\\,x"},
	}

	for _, p := range pairs {
		assert.Equal(t, []byte(p.out), EscapeMetricName([]byte(p.in)))
		assert.Equal(t, []byte(p.in), UnescapeMetricName([]byte(p.out)))
	}
}
