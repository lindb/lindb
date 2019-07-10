package discovery

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"

	"go.uber.org/zap"
)

type Register struct {
	key  string
	node models.Node
	ttl  int64
	repo state.Repository
	log  *zap.Logger
}

// NewRegister returns a new register.the key must began with the watch prefix key
func NewRegister(repo state.Repository, key string, node models.Node, ttl int64) *Register {
	return &Register{key: key, node: node, ttl: ttl, repo: repo, log: logger.GetLogger()}
}

// Register registers the node info with prefix and key in the RegisterInfo
func (r *Register) Register(ctx context.Context) error {
	bytes, err := json.Marshal(r.node)
	if err != nil {
		r.log.Error("convert node to byte error", zap.Error(err))
		return err
	}
	closed, err := r.repo.Heartbeat(ctx, r.key, bytes, r.ttl)
	go r.handlerHartBeatChanClosed(ctx, closed)
	return err
}

// UnRegister deletes the node with prefix and key
func (r *Register) UnRegister(ctx context.Context) error {
	return r.repo.Delete(ctx, r.key)
}

// handlerHartBeatChanClosed handlers the event of heartbeat closed
func (r *Register) handlerHartBeatChanClosed(ctx context.Context, closed <-chan state.Closed) {
	select {
	case <-ctx.Done():
		r.log.Warn("the context is canceled,the heartbeat exist")
	case <-closed:
		r.log.Warn("the heartbeat is closed,try again")
		for {
			if err := r.Register(ctx); err != nil {
				r.log.Warn("heartbeat failed,sleep 500ms",
					zap.String("key", r.key))
				time.Sleep(500 * time.Millisecond)
			} else {
				r.log.Info("retry heartbeat success",
					zap.String("key", r.key))
				break
			}
		}
	}
}
