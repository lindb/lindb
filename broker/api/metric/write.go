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

package metric

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/replication"
	pb "github.com/lindb/lindb/rpc/proto/field"
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
		var metrics []*pb.Metric
		for j := 0; j < 4; j++ {
			for k := 0; k < 20; k++ {
				metric := &pb.Metric{
					Name:      "cpu",
					Timestamp: timeutil.Now() + 10*timeutil.OneSecond*int64(n),
					Fields: []*pb.Field{
						{Name: "f2", Type: pb.FieldType_Sum, Value: 1.0},
					},
					Tags: map[string]string{"host": "1.1.1." + strconv.Itoa(i), "disk": "/tmp" + strconv.Itoa(j), "partition": "partition" + strconv.Itoa(k)},
				}
				n++
				metrics = append(metrics, metric)
			}
		}
		//TODO mock data for test
		metricList := &pb.MetricList{
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
