package client

import (
	"net/url"
	"sync"

	resty "github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./metric.go -destination=./metric_mock.go -package=client

// MetricCli represents metric explore client.
type MetricCli interface {
	// FetchMetricData fetches the state metric from each live nodes.
	FetchMetricData(nodes []models.Node, names []string) (interface{}, error)
}

// metricCli implements MetricCli interface.
type metricCli struct {
	logger *logger.Logger
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
