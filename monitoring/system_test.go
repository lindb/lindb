package monitoring

import (
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/stretchr/testify/assert"
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
