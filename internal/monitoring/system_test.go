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
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/v3/cpu"
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
