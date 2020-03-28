package broker

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/monitoring"
)

func TestMonitoringStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newCollector = monitoring.NewHTTPCollector
		ctrl.Finish()
	}()

	collector := monitoring.NewMockHTTPCollector(ctrl)
	newCollector = func(ctx context.Context, target, endpoint string, interval time.Duration) monitoring.HTTPCollector {
		return collector
	}

	sm := NewMonitoringStateMachine(context.TODO(), "endpoint", 10*time.Second)
	collector.EXPECT().Run().AnyTimes()
	collector.EXPECT().Stop().MaxTimes(2)
	sm.Start("target1")
	sm.Start("target1")
	sm.Stop("target1")
	sm.Start("target2")
	sm.StopAll()
	sm.StopAll()
}
