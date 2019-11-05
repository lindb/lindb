package monitoring

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var log = logger.GetLogger("monitoring", "system")
var cpus = 0

func init() {
	cores, err := cpu.Counts(true)
	if err != nil {
		log.Error("get cpu cores", logger.Error(err))
	}
	cpus = cores
}

// GetCPUs returns the number of logical cores in the system
func GetCPUs() int {
	return cpus
}

// GetCPUStat return the cpu time statistics
func GetCPUStat() *models.CPUStat {
	s, err := cpu.Times(false)
	if err != nil {
		log.Error("get cpu stat", logger.Error(err))
		return nil
	}
	if len(s) == 0 {
		log.Error("cannot get cpu stat")
		return nil
	}
	allStat := s[0]
	return &models.CPUStat{
		User:   allStat.User,
		System: allStat.System,
		Idle:   allStat.Idle,
		Nice:   allStat.Nice,
	}
}

// GetDiskStat returns a file system usage. path is a filesystem path such
// as "/", not device file path like "/dev/vda1".
func GetDiskStat(path string) *models.DiskStat {
	s, err := disk.Usage(path)
	if err != nil {
		log.Error("get disk stat", logger.Error(err))
		return nil
	}
	return &models.DiskStat{
		Total:       s.Total,
		Used:        s.Used,
		UsedPercent: s.UsedPercent,
	}
}

// GetMemoryStat return the memory usage statistics
func GetMemoryStat() *models.MemoryStat {
	v, err := mem.VirtualMemory()
	if err != nil {
		log.Error("get memory stat", logger.Error(err))
		return nil
	}
	return &models.MemoryStat{
		Total:       v.Total,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}
}
