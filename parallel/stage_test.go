package parallel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStage_String(t *testing.T) {
	assert.Equal(t, "scanner", Scanner.String())
	assert.Equal(t, "filtering", Filtering.String())
	assert.Equal(t, "grouping", Grouping.String())
	assert.Equal(t, "downSampling", DownSampling.String())
	assert.Equal(t, "unknown", Stage(99999).String())
}
