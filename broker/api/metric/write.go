package metric

import (
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/rpc/proto/field"
)

type WriteAPI struct {
	cm replication.ChannelManager
}

func NewWriteAPI(cm replication.ChannelManager) *WriteAPI {
	return &WriteAPI{
		cm: cm,
	}
}

func (m *WriteAPI) Sum(w http.ResponseWriter, r *http.Request) {
	databaseName, err := api.GetParamsFromRequest("db", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}

	//TODO mock data for test
	metric := &field.Metric{
		Name:      "cpu",
		Timestamp: timeutil.Now(),
		Fields: []*field.Field{
			{Name: "f1", Field: &field.Field_Sum{Sum: &field.Sum{
				Value: 1.0,
			}}},
		},
	}

	metricList := &field.MetricList{
		Database: databaseName,
		Metrics:  []*field.Metric{metric},
	}

	if err := m.cm.Write(metricList); err != nil {
		api.Error(w, err)
		return
	}

	api.OK(w, "ok")
}
