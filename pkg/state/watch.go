package state

import (
	"context"
	"strings"
	"time"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

const (
	defaultRetryInterval = 100 * time.Millisecond
)

type watcher struct {
	ctx  context.Context
	cli  *etcdRepository
	key  string
	opts []etcdcliv3.OpOption

	EventC WatchEventChan
}

func newWatcher(ctx context.Context, cli *etcdRepository, key string, opts ...etcdcliv3.OpOption) *watcher {
	eventc := make(chan *Event)
	w := &watcher{
		ctx:  ctx,
		cli:  cli,
		key:  key,
		opts: opts,

		EventC: eventc,
	}
	go w.watch(eventc)
	return w
}

func (w *watcher) watch(eventc chan<- *Event) {
	defer close(eventc)

	cli := w.cli.client
	var evtAll *Event
	var resp *etcdcliv3.GetResponse
	// The etcdcliv3.Watch may fail if ErrCompacted or other errors occurs.
	for {
		for {
			var err error
			if resp, err = cli.Get(w.ctx, w.key, w.opts...); err == nil {
				evtAll = w.packAllEvents(resp.Kvs)
				break
			}
			select {
			case <-w.ctx.Done():
				return
			case <-time.After(defaultRetryInterval):
			}
		}
		select {
		case <-w.ctx.Done():
			return
		case eventc <- evtAll:
		}

		opts := append(w.opts, etcdcliv3.WithRev(resp.Header.Revision+1))
		wchc := cli.Watch(w.ctx, w.key, opts...)
		if wchc == nil {
			continue
		}
		for watchResp := range wchc {
			if err := watchResp.Err(); err != nil {
				select {
				case <-w.ctx.Done():
					return
				case eventc <- &Event{Err: err}:
				}
				continue
			}
			for _, event := range watchResp.Events {
				select {
				case <-w.ctx.Done():
					return
				case eventc <- w.packWatchEvent(event):
				}
			}
		}
	}
}

func (w *watcher) parseKey(key string) string {
	if len(w.cli.namespace) == 0 {
		return key
	}
	return strings.Replace(key, w.cli.namespace, "", 1)
}

func (w *watcher) packWatchEvent(watchEvent *etcdcliv3.Event) *Event {
	kv := watchEvent.Kv
	evt := &Event{
		Type: EventTypeModify,
		KeyValues: []EventKeyValue{
			{Key: w.parseKey(string(kv.Key)), Value: kv.Value, Rev: kv.ModRevision},
		},
	}
	if watchEvent.Type == mvccpb.DELETE {
		evt.Type = EventTypeDelete
	}
	return evt
}

func (w *watcher) packAllEvents(kvs []*mvccpb.KeyValue) *Event {
	evt := &Event{Type: EventTypeAll}
	for _, kv := range kvs {
		evt.KeyValues = append(evt.KeyValues, EventKeyValue{
			Key:   w.parseKey(string(kv.Key)),
			Value: kv.Value,
			Rev:   kv.ModRevision,
		})
	}
	return evt
}
