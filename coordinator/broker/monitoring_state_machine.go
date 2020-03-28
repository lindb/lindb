package broker

import (
	"context"
	"sync"
	"time"

	"github.com/lindb/lindb/monitoring"
)

//go:generate mockgen -source=./monitoring_state_machine.go -destination=./monitoring_state_machine_mock.go -package=broker

// for testing
var (
	newCollector = monitoring.NewHTTPCollector
)

// MonitoringStateMachine represents monitoring state machine which manages monitoring target
type MonitoringStateMachine interface {
	// Start starts monitoring target
	Start(target string)
	// Stop stops monitoring target
	Stop(target string)
	// StopAll stops all monitoring target
	StopAll()
}

// monitoringStateMachine implements MonitoringStateMachine
type monitoringStateMachine struct {
	endpoint string
	interval time.Duration

	nodes map[string]monitoring.HTTPCollector // target monitoring node => http collector

	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.Mutex
}

// NewMonitoringStateMachine creates the monitoring state machine
func NewMonitoringStateMachine(ctx context.Context, endpoint string, interval time.Duration) MonitoringStateMachine {
	c, cancel := context.WithCancel(ctx)
	return &monitoringStateMachine{
		endpoint: endpoint,
		interval: interval,
		nodes:    make(map[string]monitoring.HTTPCollector),
		ctx:      c,
		cancel:   cancel,
	}
}

// Start starts monitoring target
func (m *monitoringStateMachine) Start(target string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_, ok := m.nodes[target]
	if !ok {
		collect := newCollector(m.ctx, target, m.endpoint, m.interval)
		m.nodes[target] = collect
		// start collect data in background goroutine
		go collect.Run()
	}
}

// Stop stops monitoring target
func (m *monitoringStateMachine) Stop(target string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	collect, ok := m.nodes[target]
	if ok {
		// if exist, stop collect monitoring data
		collect.Stop()
		delete(m.nodes, target)
	}
}

// StopAll stops all monitoring target
func (m *monitoringStateMachine) StopAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, collect := range m.nodes {
		collect.Stop()
	}
	// create new nodes cache
	m.nodes = make(map[string]monitoring.HTTPCollector)
}
