package state

import (
	"context"
	"errors"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
)

const defaultTTL = 10 // defalut ttl => 5 seconds

// define errors
var errKeepaliveStopped = errors.New("heartbeat keepalive stopped")

type heartbeat struct {
	client *etcd.Client
	key    string
	value  []byte

	keepaliveCh <-chan *etcd.LeaseKeepAliveResponse

	ttl int64
}

func newHeartbeat(client *etcd.Client, key string, value []byte, ttl int64) *heartbeat {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	return &heartbeat{
		client: client,
		key:    key,
		value:  value,
		ttl:    ttl,
	}
}

func (h *heartbeat) grantKeepAliveLease(ctx context.Context) error {
	resp, err := h.client.Grant(ctx, h.ttl)
	if err != nil {
		return err
	}
	_, err = h.client.Put(ctx, h.key, string(h.value), etcd.WithLease(resp.ID))
	if err != nil {
		return err
	}
	h.keepaliveCh, err = h.client.KeepAlive(ctx, resp.ID)
	return err
}

func (h *heartbeat) keepAlive(ctx context.Context) {
	log := logger.GetLogger()
	var (
		err error
		gap = 100 * time.Millisecond
	)

	for {
		if err != nil {
			log.Error("do heartbeat keepalive error, retry.", zap.Error(err))
			time.Sleep(gap)
			err = h.grantKeepAliveLease(ctx)
		} else {
			err = h.handleAliveResp(ctx)
			if err != nil && err != errKeepaliveStopped {
				return
			}
		}
	}

}
func (h *heartbeat) handleAliveResp(ctx context.Context) error {
	select {
	case aliveResp := <-h.keepaliveCh:
		if aliveResp == nil {
			return errKeepaliveStopped
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
