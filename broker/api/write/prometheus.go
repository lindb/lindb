package write

import (
	"io/ioutil"
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/protocol"
	"github.com/lindb/lindb/replication"
)

// for testing
var (
	readAllFunc = ioutil.ReadAll
)

// PrometheusWrite represents support prometheus text protocol
type PrometheusWrite struct {
	cm replication.ChannelManager
}

// NewPrometheusWrite creates prometheus write
func NewPrometheusWrite(cm replication.ChannelManager) *PrometheusWrite {
	return &PrometheusWrite{
		cm: cm,
	}
}

// Write parses prometheus text protocol then writes data into wal
func (m *PrometheusWrite) Write(w http.ResponseWriter, r *http.Request) {
	databaseName, err := api.GetParamsFromRequest("db", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	_, _ = api.GetParamsFromRequest("ns", r, constants.DefaultNamespace, false)
	s, err := readAllFunc(r.Body)
	if err != nil {
		api.Error(w, err)
		return
	}

	metricList, err := protocol.PromParse(s)
	if err != nil {
		api.Error(w, err)
		return
	}

	if err := m.cm.Write(databaseName, metricList); err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, "success")
}
