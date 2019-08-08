package standalone

import (
	"net/url"
	"time"

	"github.com/coreos/etcd/embed"

	"github.com/lindb/lindb/broker"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/storage"
)

var log = logger.GetLogger("standalone")

// runtime represents the runtime dependency of standalone mode
type runtime struct {
	state   server.State
	cfg     config.Standalone
	etcd    *embed.Etcd
	broker  server.Service
	storage server.Service
}

// NewStandaloneRuntime creates the runtime
func NewStandaloneRuntime(cfg config.Standalone) server.Service {
	return &runtime{
		state: server.New,
		cfg:   cfg,
	}
}

// Name returns the cluster mode
func (r *runtime) Name() string {
	return "standalone"
}

// Run runs the cluster as standalone mode
func (r *runtime) Run() error {
	if err := r.startETCD(); err != nil {
		r.state = server.Failed
		return err
	}
	// start storage server
	storageRuntime := storage.NewStorageRuntime(r.cfg.Storage)
	if err := storageRuntime.Run(); err != nil {
		r.state = server.Failed
		return err
	}
	// start broker server
	brokerRuntime := broker.NewBrokerRuntime(r.cfg.Broker)
	if err := brokerRuntime.Run(); err != nil {
		r.state = server.Failed
		return err
	}
	r.state = server.Running
	return nil
}

// State returns the state of cluster
func (r *runtime) State() server.State {
	return r.state
}

// Stop stops the cluster
func (r *runtime) Stop() error {
	if r.etcd != nil {
		r.etcd.Close()
		log.Info("etcd server stopped")
	}
	if r.broker != nil {
		if err := r.broker.Stop(); err != nil {
			log.Error("stop broker server", logger.Error(err))
		}
		log.Info("broker server stopped")
	}
	if r.storage != nil {
		if err := r.storage.Stop(); err != nil {
			log.Error("stop storage server", logger.Error(err))
		}
		log.Info("storage server stopped")
	}
	r.state = server.Terminated
	return nil
}

// startETCD starts embed etcd server
func (r *runtime) startETCD() error {
	cfg := embed.NewConfig()
	lcurl, _ := url.Parse(r.cfg.ETCD.URL)
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.Dir = r.cfg.ETCD.Dir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		r.etcd = e
		return err
	}
	select {
	case <-e.Server.ReadyNotify():
		log.Info("etcd server is ready")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Error("etcd server took too long to start")
	case err := <-e.Err():
		log.Error("etcd server error", logger.Error(err))
	}
	return nil
}
