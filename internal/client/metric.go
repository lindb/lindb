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

package client

import (
	"net/url"
	"sync"

	resty "github.com/go-resty/resty/v2"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

//go:generate mockgen -source=./metric.go -destination=./metric_mock.go -package=client

// MetricCli represents metric explore client.
type MetricCli interface {
	// FetchMetricData fetches the state metric from each live nodes.
	FetchMetricData(nodes []models.Node, names []string) (interface{}, error)
}

// metricCli implements MetricCli interface.
type metricCli struct {
	logger logger.Logger
}

// NewMetricCli creates a MetricCli instance.
func NewMetricCli() MetricCli {
	return &metricCli{
		logger: logger.GetLogger("Client", "Metric"),
	}
}

// FetchMetricData fetches the state metric from each live nodes.
func (cli *metricCli) FetchMetricData(nodes []models.Node, names []string) (interface{}, error) {
	size := len(nodes)
	if size == 0 {
		return nil, nil
	}
	result := make([]map[string][]*models.StateMetric, size)
	params := make(url.Values)
	for _, name := range names {
		params.Add("names", name)
	}

	var wait sync.WaitGroup
	wait.Add(size)
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			address := node.HTTPAddress()
			metric := make(map[string][]*models.StateMetric)
			_, err := resty.New().R().SetQueryParamsFromValues(params).
				SetHeader("Accept", "application/json").
				SetResult(&metric).
				Get(address + constants.APIVersion1CliPath + "/state/explore/current")
			if err != nil {
				cli.logger.Error("get current metric state from alive node", logger.String("url", address), logger.Error(err))
				return
			}
			result[i] = metric
		}()
	}
	wait.Wait()
	rs := make(map[string][]*models.StateMetric)
	for _, metricList := range result {
		if metricList == nil {
			continue
		}
		for name, list := range metricList {
			if l, ok := rs[name]; ok {
				l = append(l, list...)
				rs[name] = l
			} else {
				rs[name] = list
			}
		}
	}
	return rs, nil
}
