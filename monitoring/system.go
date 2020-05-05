package monitoring

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var log = logger.GetLogger("monitoring", "System")

var (
	cpuCount      = 0
	once4CpuCount sync.Once
	cpuCountsFunc = cpu.Counts
	cpuTimesFunc  = cpu.Times
)

type (
	MemoryStatGetter    func() (*mem.VirtualMemoryStat, error)
	CPUStatGetter       func() (*models.CPUStat, error)
	DiskUsageStatGetter func(ctx context.Context, path string) (*disk.UsageStat, error)
	NetStatGetter       func(ctx context.Context) ([]net.IOCountersStat, error)
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
	total := allStat.Total()
	return &models.CPUStat{
		User:    allStat.User / total,
		System:  allStat.System / total,
		Idle:    allStat.Idle / total,
		Nice:    allStat.Nice / total,
		Iowait:  allStat.Iowait / total,
		Irq:     allStat.Irq / total,
		Softirq: allStat.Softirq / total,
		Steal:   allStat.Steal / total,
	}, nil
}

// GetNetStat return the network usage statistics
func GetNetStat(ctx context.Context) ([]net.IOCountersStat, error) {
	stats, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return nil, err
	}
	var availableStats []net.IOCountersStat
	for _, stat := range stats {
		switch {
		// OS X
		case strings.HasPrefix(stat.Name, "en"):
		// Linux
		case strings.HasPrefix(stat.Name, "eth"):
		default:
			continue
		}
		// ignore empty interface
		if stat.BytesRecv == 0 || stat.BytesSent == 0 {
			continue
		}
		availableStats = append(availableStats, stat)
	}
	return availableStats, nil
}
