package monitoring

import (
	"context"
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

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
	if err := r.repository.Put(r.ctx, r.path, encoding.JSONMarshal(r.nodeStat)); err != nil {
		log.Error("report stat error", logger.String("path", r.path))
	}
}
