package monitoring

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestGetCPUs(t *testing.T) {
	cpus := GetCPUs()
	assert.True(t, cpus > 0)
}

func TestGetMemoryStat(t *testing.T) {
	stat, err := GetMemoryStat()
	assert.Nil(t, err)
	assert.True(t, stat.Total > 0)
	assert.True(t, stat.Used > 0)
	assert.True(t, stat.UsedPercent > 0)
}

func TestGetDiskStat(t *testing.T) {
	t.Log(filepath.VolumeName("/tmp/test/test11111"))
	_, err := GetDiskStat("/tmp/test/test111")
	assert.NotNil(t, err)

	stat, err := GetDiskStat(fileutil.GetExistPath("/tmp/test/test11111"))
	assert.Nil(t, err)
	assert.True(t, stat.Total > 0)
	assert.True(t, stat.Used > 0)
	assert.True(t, stat.UsedPercent > 0)
}

func TestGetCPUStat(t *testing.T) {
	_, err := GetCPUStat()
	assert.Nil(t, err)
}
