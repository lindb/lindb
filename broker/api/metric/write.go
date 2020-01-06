package metric

import (
	"fmt"
	"net/http"
	"strconv"

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
	c, _ := api.GetParamsFromRequest("c", r, "10", false)
	//count := 40000
	count1, _ := strconv.ParseInt(c, 10, 64)
	count := int(count1)
	var err2 error
	n := 0
	//count := 12500
	for i := 0; i < count; i++ {
		var metrics []*field.Metric
		for j := 0; j < 4; j++ {
			for k := 0; k < 20; k++ {
				metric := &field.Metric{
					Name:      "cpu",
					Timestamp: timeutil.Now(),
					Fields: []*field.Field{
						{Name: "f2", Field: &field.Field_Summary{Summary: &field.Summary{
							Sum:   2.0,
							Count: 1,
						}}},
					},
					Tags: map[string]string{"host": "1.1.1." + strconv.Itoa(i), "disk": "/tmp" + strconv.Itoa(j), "partition": "partition" + strconv.Itoa(k)},
				}
				n++
				metrics = append(metrics, metric)
			}
		}
		//TODO mock data for test
		metricList := &field.MetricList{
			Metrics: metrics,
		}
		if e := m.cm.Write(databaseName, metricList); e != nil {
			err2 = e
		}
	}
	if err2 != nil {
		api.Error(w, err2)
		return
	}
	api.OK(w, fmt.Sprintf("ok,written=%d", n))
}
