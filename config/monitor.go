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

package config

import (
	"fmt"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

var (
	// defaultPusherURL is the default push target url of LinDB
	defaultPusherURL = "http://127.0.0.1:9000/api/v1/metric/prometheus?db=_internal"
)

// Monitor represents a configuration for the internal monitor
type Monitor struct {
	ReportInterval ltoml.Duration `toml:"report-interval"`
	URL            string         `toml:"url"`
}

// TOML returns Monitor's toml config
func (m *Monitor) TOML() string {
	return fmt.Sprintf(`
[monitor]
  ## Config for the Internal Monitor
  ## monitor won't start when interval is sets to 0
  ## such as cpu, memory, and disk, process and go runtime
  report-interval = "%s"

  ## URL is the target of prometheus pusher 
  url = "%s"`,
		m.ReportInterval.String(),
		m.URL,
	)
}

// NewDefaultMonitor returns a new default monitor config
func NewDefaultMonitor() *Monitor {
	return &Monitor{
		ReportInterval: ltoml.Duration(10 * time.Second),
		URL:            defaultPusherURL,
	}
}
