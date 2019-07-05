package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
)

// Discovery defines a discovery of a list of node.it will watch the node
// online and offline and update the the lived node list
type Discovery struct {
	serverMap atomic.Value
}

// NewDiscovery returns a Discovery who will watch the changes of the key with
// the given prefix
func NewDiscovery(ctx context.Context, prefix string) (*Discovery, error) {
	if len(prefix) == 0 {
		return nil, fmt.Errorf("the key must not be null")
	}
	discovery := &Discovery{}
	discovery.serverMap.Store(&sync.Map{})
	repo := state.GetRepo()
	watchEventChan := repo.WatchPrefix(ctx, prefix)
	go discovery.handlerNodeChangeEvent(watchEventChan)
	return discovery, nil
}

// NodeList returns the current lived nod array
func (d *Discovery) NodeList() []*Node {
	nodeList := make([]*Node, 0)
	d.serverMap.Load().(*sync.Map).Range(func(key, value interface{}) bool {
		nodeList = append(nodeList, value.(*Node))
		return true
	})
	return nodeList
}

// handlerServerChange handles the changes of the node and update the
// node map
func (d *Discovery) handlerNodeChangeEvent(eventChan state.WatchEventChan) {
	log := logger.GetLogger()
	for event := range eventChan {
		if event.Err != nil {
			continue
		}
		switch event.Type {
		case state.EventTypeDelete:
			m := d.serverMap.Load().(*sync.Map)
			for _, kv := range event.KeyValues {
				m.Delete(kv.Key)
			}
		case state.EventTypeAll:
			d.serverMap.Store(&sync.Map{})
			fallthrough
		case state.EventTypeModify:
			m := d.serverMap.Load().(*sync.Map)
			for _, kv := range event.KeyValues {
				node := &Node{}
				if err := json.Unmarshal(kv.Value, node); err != nil {
					log.Error(" deserialize error", zap.Error(err))
				} else {
					m.Store(kv.Key, node)
				}
			}
		}
	}
}
