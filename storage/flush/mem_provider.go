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

package flush

import (
	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryUsageProvider provides the current system memory usage information.
type MemoryUsageProvider interface {
	// MemoryUsage returns the current memory usage ratio in [0, 1] and total physical
	// memory in bytes. Returns (0, 0, err) on failure; callers should treat errors as no pressure.
	MemoryUsage() (usageRatio float64, totalBytes uint64, err error)
}

// defaultMemoryUsageProvider is the gopsutil-based implementation.
type defaultMemoryUsageProvider struct{}

// NewDefaultMemoryUsageProvider creates a MemoryUsageProvider backed by gopsutil.
func NewDefaultMemoryUsageProvider() MemoryUsageProvider {
	return &defaultMemoryUsageProvider{}
}

// MemoryUsage implements MemoryUsageProvider using gopsutil VirtualMemory.
// UsedPercent is in [0, 100], so we divide by 100 to normalize to [0, 1].
func (p *defaultMemoryUsageProvider) MemoryUsage() (float64, uint64, error) {
	stat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	return stat.UsedPercent / 100.0, stat.Total, nil
}
