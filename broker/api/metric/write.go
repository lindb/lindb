package metric

import (
	"context"
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
			{Name: "f1", Field: &field.Field_Sum{Sum: 1.0}},
		},
	}

	data, _ := metric.Marshal()

	ch, err := m.cm.GetChannel(databaseName, 0)
	if err != nil {
		api.Error(w, err)
		return
	}

	if err := ch.Write(context.TODO(), data); err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, "ok")
}
