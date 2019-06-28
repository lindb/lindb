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

// heartbeat represents a heartbeat with etcd, it will start a goroutine does keepalive in background
type heartbeat struct {
	client *etcd.Client
	key    string
	value  []byte

	keepaliveCh <-chan *etcd.LeaseKeepAliveResponse

	ttl int64
}

// newHeartbeat creates heartbeat instance
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

// grantKeepAliveLease grants ectd lease, if success do keepalive
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

// keepAlive does keepalive and retry
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
			// do retry grant keep alive if lease ttl
			err = h.grantKeepAliveLease(ctx)
			// if ctx happen err, return stop keepalive
			if ctx.Err() != nil {
				return
			}
		} else {
			err = h.handleAliveResp(ctx)
			// return if keepalive stopped
			if err != nil && err != errKeepaliveStopped {
				return
			}
		}
	}

}

// handleAliveResp handles keepalive response, if keepalive closed or ctx canceled return keep liave stopped error
func (h *heartbeat) handleAliveResp(ctx context.Context) error {
	select {
	case aliveResp := <-h.keepaliveCh:
		if aliveResp == nil {
			return errKeepaliveStopped
		}
	case <-ctx.Done():
		return errKeepaliveStopped
	}
	return nil
}
