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

//go:build integration
// +build integration

package standalone

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/standalone"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/series/metric"
)

func TestMain(m *testing.M) {
	cfg := config.NewDefaultStandalone()
	if err := logger.InitLogger(cfg.Logging, "standalone.log"); err != nil {
		panic(fmt.Errorf("init logging err: %s", err))
	}

	gin.SetMode(gin.ReleaseMode)
	// run cluster as standalone mode
	runtime := standalone.NewStandaloneRuntime(config.Version, &cfg)
	if err := runtime.Run(); err != nil {
		panic(fmt.Errorf("run standalone cluster err: %s", err))
	}
	defer func() {
		fmt.Println("all tests run completed, now stop standalone cluster")
		runtime.Stop()
	}()
	// wait server start complete
	time.Sleep(10 * time.Second)
	fmt.Println("standalone cluster start successfully")
	mockMetricData()
	// wait data write complete
	time.Sleep(10 * time.Second)
	fmt.Println("write mock metric completed")

	m.Run()
}

func TestQuery_Group_by(t *testing.T) {
	cli := client.NewExecuteCli("http://localhost:9000/api")
	rs, err := cli.ExecuteAsResult(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f1 from cpu_data where host='host1' and time>now()-1h group by host,app",
	}, &models.ResultSet{})
	assert.NoError(t, err)
	fmt.Println(rs)

	// no data found
	rs, err = cli.ExecuteAsResult(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f1 from cpu_data where host='host3434' and time>now()-1h group by host,app",
	}, &models.ResultSet{})
	assert.NoError(t, err)
	fmt.Println(rs)
}

func TestTagValueNotFound(t *testing.T) {
	cli := client.NewExecuteCli("http://localhost:9000/api")
	rs, err := cli.ExecuteAsResult(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f1 from cpu_data where host='host' and time>now()-1h group by host,app",
	}, &models.ResultSet{})
	assert.NoError(t, err)
	fmt.Println(rs)
}

func TestMetaNotFound(t *testing.T) {
	cli := client.NewExecuteCli("http://localhost:9000/api")
	err := cli.Execute(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f4 from cpu_data where host='host' and time>now()-1h group by host,app",
	}, &models.ResultSet{})
	assert.Equal(t, err, errors.New(`"field not found, field: f4"`))

	err = cli.Execute(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f1 from cpu_data2 where host='host' and time>now()-1h group by host,app",
	}, &models.ResultSet{})
	assert.Equal(t, err, errors.New(`"metric not found, metric: cpu_data2"`))

	err = cli.Execute(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f1 from cpu_data where host2='host' and time>now()-1h group by host,app",
	}, &models.ResultSet{})
	assert.Equal(t, err, errors.New(`"tag key not found, tag key: host2"`))
	err = cli.Execute(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select f1 from cpu_data where host='host' and time>now()-1h group by host,app2",
	}, &models.ResultSet{})
	assert.Equal(t, err, errors.New(`"tag key not found, tag key: app2"`))
}

func mockMetricData() {
	timestamp := timeutil.Now()
	var buf bytes.Buffer
	for i := 0; i < 50; i++ {
		var brokerRow metric.BrokerRow
		converter := metric.NewProtoConverter()
		err := converter.ConvertTo(&protoMetricsV1.Metric{
			Name:      "cpu_data",
			Timestamp: timestamp,
			Tags: []*protoMetricsV1.KeyValue{
				{Key: "host", Value: "host" + strconv.Itoa(i)},
				{Key: "app", Value: "app" + strconv.Itoa(i)},
			},
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1},
				{Name: "f2", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 2},
			},
		}, &brokerRow)
		_, _ = brokerRow.WriteTo(&buf)
		if err != nil {
			panic(err)
		}
	}
	body := buf.Bytes()
	_, err := resty.New().R().SetBody(body).Put("http://127.0.0.1:9000/api/flat/write?db=_internal")
	if err != nil {
		panic(err)
	}
}
