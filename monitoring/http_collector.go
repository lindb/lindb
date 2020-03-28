package monitoring

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./http_collector.go -destination=./http_collector_mock.go -package=monitoring

// for testing
var (
	doGet      = http.Get
	readBody   = ioutil.ReadAll
	newRequest = http.NewRequest
	doRequest  = http.DefaultClient.Do
)

// HTTPCollector represents collect LinDB self metrics via prometheus exporter,
// then writes the metric into internal database.
type HTTPCollector interface {
	// Run runs collect metrics in period
	Run()
	// Stop tops collect metrics
	Stop()
}

// httpCollector implements HTTPCollector
type httpCollector struct {
	target   string // prometheus http exporter
	endpoint string // HTTP endpoint

	ctx      context.Context
	cancel   context.CancelFunc
	interval time.Duration
}

// NewHTTPCollector creates the http collector which collects LinDB self metrics
func NewHTTPCollector(
	ctx context.Context,
	target, endpoint string,
	interval time.Duration,
) HTTPCollector {
	c, cancel := context.WithCancel(ctx)
	return &httpCollector{
		ctx:      c,
		cancel:   cancel,
		target:   target,
		endpoint: endpoint,
		interval: interval,
	}
}

// Run runs collect metrics in period
func (c *httpCollector) Run() {
	c.collect()
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.collect()
		case <-c.ctx.Done():
			log.Info("stop collect monitoring data by http", logger.String("target", c.target))
			return
		}
	}
}

// Stop tops collect metrics
func (c *httpCollector) Stop() {
	c.cancel()
}

// collect collects metrics and writes to internal database
func (c *httpCollector) collect() {
	// 1. collect metric via prometheus exporter
	getResp, err := doGet(c.target)
	if err != nil {
		log.Error("get monitoring data error", logger.String("target", c.target), logger.Error(err))
		return
	}
	defer func() {
		_ = getResp.Body.Close()
	}()
	// 2. read metric data
	body, err := readBody(getResp.Body)
	if err != nil {
		log.Error("reader monitoring data error", logger.String("target", c.target), logger.Error(err))
		return
	}
	// 3. new metric write request
	req, err := newRequest("PUT", c.endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Error("new write monitoring request error", logger.String("target", c.target), logger.Error(err))
		return
	}
	req.Header.Add("Content-Type", "text/plain")
	// 4. send metric data
	writeResp, err := doRequest(req)
	if err != nil {
		log.Error("write monitoring data error", logger.Error(err))
		return
	}
	_ = writeResp.Body.Close()
}
