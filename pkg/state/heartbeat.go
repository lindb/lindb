package state

import (
	"context"
	"errors"
	"time"

	"github.com/eleme/lindb/pkg/logger"

	etcd "github.com/coreos/etcd/clientv3"
)

const defaultTTL = 10 // default ttl => 10 seconds

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

//TODO need refactor
func (h *heartbeat) PutIfNotExist(ctx context.Context) (bool, error) {
	resp, err := h.client.Grant(ctx, h.ttl)
	if err != nil {
		return false, err
	}
	txn := h.client.Txn(ctx).If(etcd.Compare(etcd.CreateRevision(h.key), "=", 0))
	txn = txn.Then(etcd.OpPut(h.key, string(h.value), etcd.WithLease(resp.ID)))
	txn = txn.Else(etcd.OpGet(h.key))
	response, err := txn.Commit()
	if err != nil {
		return false, err
	}
	response.Responses[0].GetResponse()
	if response.Succeeded {
		h.keepaliveCh, err = h.client.KeepAlive(ctx, resp.ID)
	}
	return response.Succeeded, err
}

// keepAlive does keepalive and retry,if the key should be not exist,it should retry
func (h *heartbeat) keepAlive(ctx context.Context, keepAliveNotExit bool) {
	log := logger.GetLogger("state/heartbeat")
	var (
		err error
		gap = 100 * time.Millisecond
	)
	for {
		if err != nil {
			log.Error("do heartbeat keepalive error, retry.", logger.Error(err))
			time.Sleep(gap)
			if keepAliveNotExit {
				// retry put if not exist.if failed closes the heartbeat
				isSuccess, e := h.PutIfNotExist(ctx)
				err = e
				if !isSuccess {
					// put if not exist failed ,close heartbeat
					return
				}
			} else {
				// do retry grant keep alive if lease ttl
				err = h.grantKeepAliveLease(ctx)
			}
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
