package monitoring

import (
	"fmt"
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
	stat := GetMemoryStat()
	assert.NotNil(t, stat)
	assert.True(t, stat.Total > 0)
	assert.True(t, stat.Used > 0)
	assert.True(t, stat.UsedPercent > 0)
}

func TestGetDiskStat(t *testing.T) {
	fmt.Println(filepath.VolumeName("/tmp/test/test11111"))
	stat := GetDiskStat("/tmp/test/test111")
	assert.Nil(t, stat)

	stat = GetDiskStat(fileutil.GetExistPath("/tmp/test/test11111"))
	assert.NotNil(t, stat)
	assert.True(t, stat.Total > 0)
	assert.True(t, stat.Used > 0)
	assert.True(t, stat.UsedPercent > 0)
}

func TestGetCPUStat(t *testing.T) {
	stat := GetCPUStat()
	assert.NotNil(t, stat)
}
