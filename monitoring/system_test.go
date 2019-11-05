package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCPUs(t *testing.T) {
	cpus := GetCPUs()
	assert.True(t, cpus > 0)
}

func TestGetMemoryStat(t *testing.T) {
	stat := GetMemoryStat()
	assert.NotNil(t, stat)
	assert.True(t, stat.Total > 0)
	assert.True(t, stat.Used > 0)
	assert.True(t, stat.UsedPercent > 0)
}

func TestGetDiskStat(t *testing.T) {
	stat := GetDiskStat("/tmp")
	assert.NotNil(t, stat)
	assert.True(t, stat.Total > 0)
	assert.True(t, stat.Used > 0)
	assert.True(t, stat.UsedPercent > 0)
}

func TestGetCPUStat(t *testing.T) {
	stat := GetCPUStat()
	assert.NotNil(t, stat)
}
