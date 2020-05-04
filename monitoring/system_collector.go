package monitoring

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

var (
	memStatGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_mem_stat",
		Help: "System mem stats",
	}, []string{"type"})
	cpuStatGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_cpu_stat",
		Help: "System cpu stats",
	}, []string{"type"})
	diskStatGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_disk_stat",
		Help: "System disk stats",
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(memStatGauge)
	prometheus.MustRegister(cpuStatGauge)
	prometheus.MustRegister(diskStatGauge)
}

// SystemCollector collects the system stat
type SystemCollector struct {
	ctx        context.Context
	interval   time.Duration
	storage    string
	repository state.Repository // data will be putted to this
	path       string           // repository key
	systemStat *models.SystemStat
	nodeStat   *models.NodeStat
	// used for mock
	MemoryStatGetter MemoryStatGetter
	CPUStatGetter    CPUStatGetter
	DiskStatGetter   DiskStatGetter
}

// NewSystemCollector creates a new system stat collector
func NewSystemCollector(
	ctx context.Context,
	interval time.Duration,
	storage string,
	repository state.Repository,
	path string,
	node models.ActiveNode,
) *SystemCollector {
	r := &SystemCollector{
		interval:   interval,
		storage:    fileutil.GetExistPath(storage),
		repository: repository,
		path:       path,
		systemStat: &models.SystemStat{},
		nodeStat: &models.NodeStat{
			Node: node,
		},
		ctx:              ctx,
		MemoryStatGetter: GetMemoryStat,
		CPUStatGetter:    GetCPUStat,
		DiskStatGetter:   GetDiskStat,
	}
	return r
}

// Run starts a background goroutine that collects the monitoring stat
func (r *SystemCollector) Run() {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	// collect system status
	r.collect()

	for {
		select {
		case <-ticker.C:
			// collect system status
			r.collect()
		case <-r.ctx.Done():
			return
		}
	}
}

// collect collects the monitoring stat
func (r *SystemCollector) collect() {
	var err error
	r.systemStat.CPUs = GetCPUs()

	if r.systemStat.MemoryStat, err = r.MemoryStatGetter(); err != nil {
		log.Error("get memory stat", logger.Error(err))
	}
	if r.systemStat.CPUStat, err = r.CPUStatGetter(); err != nil {
		log.Error("get cpu stat", logger.Error(err))
	}
	if r.systemStat.DiskStat, err = r.DiskStatGetter(r.storage); err != nil {
		log.Error("get disk stat", logger.Error(err))
	}

	r.nodeStat.System = *r.systemStat

	if r.systemStat.MemoryStat != nil {
		memStat := r.systemStat.MemoryStat
		memStatGauge.WithLabelValues("total").Set(float64(memStat.Total))
		memStatGauge.WithLabelValues("used").Set(float64(memStat.Used))
		memStatGauge.WithLabelValues("used_percent").Set(memStat.UsedPercent)
	}

	if r.systemStat.DiskStat != nil {
		diskStat := r.systemStat.DiskStat
		diskStatGauge.WithLabelValues("total").Set(float64(diskStat.Total))
		diskStatGauge.WithLabelValues("used").Set(float64(diskStat.Used))
		diskStatGauge.WithLabelValues("used_percent").Set(diskStat.UsedPercent)
	}

	if r.systemStat.CPUStat != nil {
		cpuStat := r.systemStat.CPUStat
		cpuStatGauge.WithLabelValues("idle").Set(cpuStat.Idle)
		cpuStatGauge.WithLabelValues("nice").Set(cpuStat.Nice)
		cpuStatGauge.WithLabelValues("system").Set(cpuStat.System)
		cpuStatGauge.WithLabelValues("user").Set(cpuStat.User)
		cpuStatGauge.WithLabelValues("irq").Set(cpuStat.Irq)
		cpuStatGauge.WithLabelValues("steal").Set(cpuStat.Steal)
		cpuStatGauge.WithLabelValues("softirq").Set(cpuStat.Softirq)
		cpuStatGauge.WithLabelValues("iowait").Set(cpuStat.Iowait)
	}

	if err := r.repository.Put(r.ctx, r.path, encoding.JSONMarshal(r.nodeStat)); err != nil {
		log.Error("report stat error", logger.String("path", r.path))
	}
}
