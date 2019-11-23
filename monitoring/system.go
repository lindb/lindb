package monitoring

import (
	"fmt"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var log = logger.GetLogger("monitoring", "System")

var (
	cpuCount      = 0
	once4CpuCount sync.Once
	cpuCountsFunc = cpu.Counts
	cpuTimesFunc  = cpu.Times
	memFunc       = mem.VirtualMemory
)

type (
	MemoryStatGetter func() (*models.MemoryStat, error)
	CPUStatGetter    func() (*models.CPUStat, error)
	DiskStatGetter   func(path string) (*models.DiskStat, error)
)

// GetCPUs returns the number of logical cores in the system
func GetCPUs() int {
	once4CpuCount.Do(
		func() {
			cpuCount = getCPUs()
		})
	return cpuCount
}

// getCPUs returns the number of logical cores in the system
func getCPUs() int {
	count, err := cpuCountsFunc(true)
	if err != nil {
		log.Error("get cpu cores", logger.Error(err))
	}
	return count
}

// GetCPUStat return the cpu time statistics
func GetCPUStat() (*models.CPUStat, error) {
	s, err := cpuTimesFunc(false)
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, fmt.Errorf("cannot get cpu stat")
	}
	allStat := s[0]
	return &models.CPUStat{
		User:   allStat.User,
		System: allStat.System,
		Idle:   allStat.Idle,
		Nice:   allStat.Nice,
	}, nil
}

// GetDiskStat returns a file system usage. path is a filesystem path such
// as "/", not device file path like "/dev/vda1".
func GetDiskStat(path string) (*models.DiskStat, error) {
	s, err := disk.Usage(path)
	if err != nil {
		return nil, err
	}
	return &models.DiskStat{
		Total:       s.Total,
		Used:        s.Used,
		UsedPercent: s.UsedPercent,
	}, nil
}

// GetMemoryStat return the memory usage statistics
func GetMemoryStat() (*models.MemoryStat, error) {
	v, err := memFunc()
	if err != nil {
		return nil, err
	}
	return &models.MemoryStat{
		Total:       v.Total,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}, nil
}
