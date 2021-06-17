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

package write

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/broker/deps"
	"github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

var (
	SumWritePath = "/metric/sum"
)

type MetricWriteAPI struct {
	deps *deps.HTTPDeps
}

func NewWriteAPI(deps *deps.HTTPDeps) *MetricWriteAPI {
	return &MetricWriteAPI{
		deps: deps,
	}
}

// Register adds metric write url route.
func (m *MetricWriteAPI) Register(route gin.IRoutes) {
	route.PUT(SumWritePath, m.Sum)
}

func (m *MetricWriteAPI) Sum(c *gin.Context) {
	var param struct {
		Database string `form:"db" binding:"required"`
		Count    int    `form:"count"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		http.Error(c, err)
		return
	}
	if param.Count == 0 {
		param.Count = 10
	}
	var err2 error
	n := 0
	//count := 12500
	for i := 0; i < param.Count; i++ {
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
		if e := m.deps.CM.Write(param.Database, metricList); e != nil {
			err2 = e
		}
	}
	if err2 != nil {
		http.Error(c, err2)
		return
	}
	http.OK(c, fmt.Sprintf("ok,written=%d", n))
}
