package config

import "testing"

func Test_NewConfig(t *testing.T) {
	_ = NewDefaultBrokerCfg()
	_ = NewDefaultStandaloneCfg()
	_ = NewDefaultStorageCfg()
	_ = NewDefaultQueryCfg()
	monitorCfg := NewDefaultMonitorCfg()
	_ = monitorCfg.RuntimeReportInterval()
	_ = monitorCfg.SystemReportInterval()
}
