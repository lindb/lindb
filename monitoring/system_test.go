package monitoring

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func Test_getCPUCounts(t *testing.T) {
	defer func() {
		cpuCountsFunc = cpu.Counts
	}()

	cpuCountsFunc = func(logical bool) (i int, e error) {
		return 0, fmt.Errorf("err")
	}
	core := getCPUs()
	assert.Equal(t, 0, core)
}

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

func TestGetMemoryStat2(t *testing.T) {
	defer func() {
		memFunc = mem.VirtualMemory
	}()
	memFunc = func() (stat *mem.VirtualMemoryStat, e error) {
		return nil, fmt.Errorf("err")
	}
	stat, err := GetMemoryStat()
	assert.Nil(t, stat)
	assert.Error(t, err)
}

func TestGetDiskStat(t *testing.T) {
	fmt.Println(filepath.VolumeName("/tmp/test/test11111"))
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

func TestGetCPUStat2(t *testing.T) {
	defer func() {
		cpuTimesFunc = cpu.Times
	}()
	cpuTimesFunc = func(perCPU bool) (stats []cpu.TimesStat, e error) {
		return nil, fmt.Errorf("err")
	}
	stat, err := GetCPUStat()
	assert.Nil(t, stat)
	assert.Error(t, err)

	cpuTimesFunc = func(perCPU bool) (stats []cpu.TimesStat, e error) {
		return nil, nil
	}
	stat, err = GetCPUStat()
	assert.Nil(t, stat)
	assert.Error(t, err)
}
