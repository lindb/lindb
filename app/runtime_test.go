package app

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/monitoring"
	"github.com/lindb/lindb/series/tag"
)

func TestBaseRuntime_SystemCollector(t *testing.T) {
	r := NewBaseRuntime(context.TODO(), config.Monitor{}, linmetric.RootRegistry, tag.Tags{})
	r.SystemCollector()
}

func TestBaseRuntime_NativePusher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newNativeProtoPusher = monitoring.NewNativeProtoPusher
		ctrl.Finish()
	}()

	r := NewBaseRuntime(context.TODO(), config.Monitor{}, linmetric.RootRegistry, tag.Tags{})
	r.NativePusher()
	assert.Nil(t, r.pusher)

	pusher := monitoring.NewMockNativePusher(ctrl)
	newNativeProtoPusher = func(_ context.Context, _ string, _, _ time.Duration,
		_ *linmetric.Registry, _ tag.Tags) monitoring.NativePusher {
		return pusher
	}
	r = NewBaseRuntime(context.TODO(), config.Monitor{ReportInterval: 1000}, linmetric.RootRegistry, tag.Tags{})
	ch := make(chan struct{})
	pusher.EXPECT().Start().Do(func() {
		close(ch)
	})
	pusher.EXPECT().Stop()
	r.NativePusher()
	assert.NotNil(t, r.pusher)
	<-ch
	r.Shutdown()
}
