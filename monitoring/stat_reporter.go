package monitoring

import (
	"context"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

// Reporter represents the monitoring stat reporter
type Reporter interface {
	// Report reports the monitoring stat
	Report(stat interface{})
}

type heartbeatReport struct {
	ctx  context.Context
	repo state.Repository
	path string
}

// NewHeartbeatReporter creates a heartbeat reporter
func NewHeartbeatReporter(ctx context.Context, repo state.Repository, path string) Reporter {
	return &heartbeatReport{
		ctx:  ctx,
		repo: repo,
		path: path,
	}
}

// Report reports the monitoring stat to state repository
func (r *heartbeatReport) Report(stat interface{}) {
	if err := r.repo.Put(r.ctx, r.path, encoding.JSONMarshal(stat)); err != nil {
		log.Error("report stat error", logger.String("path", r.path))
	}
}
