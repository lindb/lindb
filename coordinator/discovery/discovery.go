package discovery

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

// Listener represents discovery resource event callback interface,
// includes create/delete/cleanup operation
type Listener interface {
	// OnCreate is resource creation callback
	OnCreate(key string, resource []byte)
	// OnDelete is resource deletion callback
	OnDelete(key string)
	// Cleanup cleans all resources
	Cleanup()
}

// Discovery represents discovery resources, through watch resource's prefix
type Discovery interface {
	// Discovery starts discovery resources change, includes create/delete/clean
	Discovery() error
	// Close stops watch, trigger all resource cleanup callback
	Close()
}

// discovery implements discovery interface
type discovery struct {
	prefix   string
	repo     state.Repository
	listener Listener

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewDiscovery returns a Discovery who will watch the changes with the given prefix
func NewDiscovery(repo state.Repository, prefix string, listener Listener) Discovery {
	ctx, cancel := context.WithCancel(context.Background())
	return &discovery{
		prefix:   prefix,
		repo:     repo,
		ctx:      ctx,
		cancel:   cancel,
		listener: listener,
		log:      logger.GetLogger("coordinator/discovery"),
	}
}

// Discovery starts discovery resources change, includes create/delete/clean
func (d *discovery) Discovery() error {
	if len(d.prefix) == 0 {
		return fmt.Errorf("watch prefix is empth for discovery resource")
	}
	watchEventCh := d.repo.WatchPrefix(d.ctx, d.prefix)
	go func() {
		d.handlerResourceChange(watchEventCh)
		d.log.Warn("exit discovery loop")
	}()
	return nil
}

// Cleanup cleans all resources
func (d *discovery) Close() {
	d.cancel()
	d.listener.Cleanup()
	if err := d.repo.Close(); err != nil {
		d.log.Error("close state repo error", logger.String("prefix", d.prefix))
	}
}

// handlerResourceChange handles the changes of event for resources
func (d *discovery) handlerResourceChange(eventCh state.WatchEventChan) {
	for event := range eventCh {
		if event.Err != nil {
			continue
		}
		switch event.Type {
		case state.EventTypeDelete:
			for _, kv := range event.KeyValues {
				d.listener.OnDelete(kv.Key)
			}
		case state.EventTypeAll:
			d.listener.Cleanup()
			fallthrough
		case state.EventTypeModify:
			for _, kv := range event.KeyValues {
				d.listener.OnCreate(kv.Key, kv.Value)
			}
		}
	}
}
