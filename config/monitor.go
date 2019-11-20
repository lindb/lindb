package config

import (
	"time"
)

// Monitor represents a configuration for the internal monitor
type Monitor struct {
	SystemReportIntervalInSeconds  int `toml:"systemReportIntervalInSeconds"`
	RuntimeReportIntervalInSeconds int `toml:"runtimeReportIntervalInSeconds"`
}

// SystemReportInterval returns a duration value
func (m *Monitor) SystemReportInterval() time.Duration {
	return time.Duration(m.SystemReportIntervalInSeconds) * time.Second
}

// RuntimeReportInterval returns a duration value
func (m *Monitor) RuntimeReportInterval() time.Duration {
	return time.Duration(m.RuntimeReportIntervalInSeconds) * time.Second
}

// NewDefaultMonitorCfg returns a new default monitor config
// zero disables the monitor
func NewDefaultMonitorCfg() Monitor {
	return Monitor{
		SystemReportIntervalInSeconds:  30,
		RuntimeReportIntervalInSeconds: 10}
}
