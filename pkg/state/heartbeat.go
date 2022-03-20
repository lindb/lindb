// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package state

import (
	"context"
	"errors"
	"time"

	"github.com/lindb/lindb/pkg/logger"

	etcd "go.etcd.io/etcd/client/v3"
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
	isElect     bool

	ttl    int64
	logger *logger.Logger
}

// newHeartbeat creates heartbeat instance
func newHeartbeat(client *etcd.Client, key string, value []byte, ttl int64, isElect bool) *heartbeat {
	if ttl <= 0 {
		ttl = defaultTTL
	}
	return &heartbeat{
		client:  client,
		isElect: isElect,
		key:     key,
		value:   value,
		ttl:     ttl,
		logger:  logger.GetLogger("etcdState", "HeartBeat"),
	}
}

// withLogger sets a new logger
func (h *heartbeat) withLogger(log *logger.Logger) {
	h.logger = log
}

// grantKeepAliveLease grants ectd lease, if success do keepalive
func (h *heartbeat) grantKeepAliveLease(ctx context.Context) (bool, error) {
	resp, err := h.client.Grant(ctx, h.ttl)
	if err != nil {
		return false, err
	}
	var ops []etcd.Cmp
	if h.isElect {
		ops = append(ops, etcd.Compare(etcd.CreateRevision(h.key), "=", 0))
	}
	txn := h.client.Txn(ctx).If(ops...)
	txn = txn.Then(etcd.OpPut(h.key, string(h.value), etcd.WithLease(resp.ID)))
	txn = txn.Else(etcd.OpGet(h.key))
	response, err := txn.Commit()
	if err != nil {
		return false, err
	}
	if response.Succeeded {
		h.keepaliveCh, err = h.client.KeepAlive(ctx, resp.ID)
	}
	return response.Succeeded, err
}

// keepAlive does keepalive and retry,if the key should be not exist,it should retry
func (h *heartbeat) keepAlive(ctx context.Context) {
	var (
		err error
		gap = 100 * time.Millisecond
	)
	for {
		if err != nil {
			h.logger.Error("do heartbeat keepalive error, retry.", logger.Error(err), logger.String("key", h.key))
			time.Sleep(gap)
			if h.isElect {
				// retry put if not exist. if failed closes the heartbeat
				isSuccess, e := h.grantKeepAliveLease(ctx)
				err = e
				if !isSuccess {
					// put if not exist failed, close heartbeat
					return
				}
			} else {
				// do retry grant keep alive if lease ttl
				_, err = h.grantKeepAliveLease(ctx)
			}
			// if ctx happen err, return stop keepalive
			if ctx.Err() != nil {
				return
			}
		} else {
			err = h.handleAliveResp(ctx)
			// return if keepalive stopped
			if err != nil && err == errKeepaliveStopped {
				return
			}
		}
	}
}

// handleAliveResp handles keepalive response, if keepalive closed or ctx canceled return keep alive stopped error
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
