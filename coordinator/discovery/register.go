package discovery

import (
	"context"
	"encoding/json"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"

	"go.uber.org/zap"
)

type Node struct {
	IP   string
	Port int32
}

type Register struct {
	node   Node
	prefix string
	key    string
	ttl    int64
}

// NewRegister returns a new register
func NewRegister(node Node, prefix string, key string, ttl int64) *Register {
	return &Register{node: node, prefix: prefix, key: key, ttl: ttl}
}

// Register registers the node info with prefix and key in the RegisterInfo
func (r *Register) Register(ctx context.Context) error {
	bytes, err := json.Marshal(r.node)
	if err != nil {
		log := logger.GetLogger()
		log.Error("convert node to byte error", zap.Error(err))
		return err
	}
	_, err = state.GetRepo().Heartbeat(ctx, r.prefix+r.key, bytes, r.ttl)
	return err
}

// UnRegister deletes the node with prefix and key
func (r *Register) UnRegister(ctx context.Context) error {
	return state.GetRepo().Delete(ctx, r.prefix+r.key)
}
