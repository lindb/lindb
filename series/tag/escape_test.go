package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnescapeTag_EscapeTag(t *testing.T) {
	pairs := []struct {
		in  string
		out string
	}{
		{"xxx", "xxx"},
		{"xx x", "xx\\ x"},
		{"xx,x", "xx\\,x"},
		{"xx=x", "xx\\=x"},
	}

	for _, p := range pairs {
		assert.Equal(t, []byte(p.out), EscapeTag([]byte(p.in)))
		assert.Equal(t, []byte(p.in), UnescapeTag([]byte(p.out)))
	}
}
