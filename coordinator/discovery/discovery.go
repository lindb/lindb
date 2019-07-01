package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"

	"go.uber.org/zap"
)

// Discovery defines a discovery of a list of node.it will watch the node
// online and offline and update the the lived node list
type Discovery struct {
	serverMap sync.Map
}

// NewDiscovery returns a Discovery who will watch the changes of the key with
// the given prefix
func NewDiscovery(ctx context.Context, prefix string) (*Discovery, error) {
	if len(prefix) == 0 {
		return nil, fmt.Errorf("the key must not be null")
	}
	discovery := &Discovery{}
	repo := state.GetRepo()
	watchEventChan, err := repo.WatchPrefix(ctx, prefix)
	if err != nil {
		return nil, err
	}
	go discovery.handlerNodeChangeEvent(watchEventChan)
	return discovery, nil
}

// NodeList returns the current lived nod array
func (d *Discovery) NodeList() []*Node {
	nodeList := make([]*Node, 0)
	d.serverMap.Range(func(key, value interface{}) bool {
		nodeList = append(nodeList, value.(*Node))
		return true
	})
	return nodeList
}

// handlerServerChange handles the changes of the node and update the
// node map
func (d *Discovery) handlerNodeChangeEvent(eventChan state.WatchEventChan) {
	for event := range eventChan {
		switch event.Type {
		case state.EventTypeDelete:
			d.serverMap.Delete(event.Key)
		case state.EventTypeModify:
			node := &Node{}
			err := json.Unmarshal(event.Value, node)
			if err != nil {
				log := logger.GetLogger()
				log.Error(" deserialize error", zap.Error(err))
			} else {
				d.serverMap.Store(event.Key, node)
			}
		}
	}
}
