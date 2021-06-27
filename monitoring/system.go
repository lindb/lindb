// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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

var collectorLogger = logger.GetLogger("monitoring", "Collector")

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
		collectorLogger.Error("get cpu cores", logger.Error(err))
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
