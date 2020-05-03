package config

import (
	"fmt"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

var (
	// defaultPusherURL is the default push target url of LinDB
	defaultPusherURL = "http://127.0.0.1:9000/metric/prometheus?db=_internal"
)

// Monitor represents a configuration for the internal monitor
type Monitor struct {
	SystemReportInterval  ltoml.Duration `toml:"system-report-interval"`
	RuntimeReportInterval ltoml.Duration `toml:"runtime-report-interval"`
	URL                   string         `toml:"url"`
}

// TOML returns Monitor's toml config
func (m *Monitor) TOML() string {
	return fmt.Sprintf(`
[monitor]
  ## Config for the Internal Monitor
  ## monitor won't start when interval is sets to 0
  
  ## system-monitor collects the system metrics, 
  ## such as cpu, memory, and disk
  system-report-interval = "%s"
  
  ## runtime-monitor collects the golang runtime memory metrics,
  ## such as stack, heap, off-heap, and gc
  runtime-report-interval = "%s"

  ## URL is the target of prometheus pusher 
  url = "%s"`,
		m.SystemReportInterval.String(),
		m.RuntimeReportInterval.String(),
		m.URL,
	)
}

// NewDefaultMonitor returns a new default monitor config
func NewDefaultMonitor() *Monitor {
	return &Monitor{
		SystemReportInterval:  ltoml.Duration(30 * time.Second),
		RuntimeReportInterval: ltoml.Duration(10 * time.Second),
		URL:                   defaultPusherURL,
	}
}
