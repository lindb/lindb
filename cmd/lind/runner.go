package lind

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/pkg/server"
)

// serveStandalone runs the cluster as standalone mode
func run(ctx context.Context, service server.Service) error {
	// start service
	if err := service.Run(); err != nil {
		return fmt.Errorf("run service[%s] error:%s", service.Name(), err)
	}

	// waiting system exit signal
	<-ctx.Done()

	// stop service
	if err := service.Stop(); err != nil {
		return fmt.Errorf("stop service[%s] error:%s", service.Name(), err)
	}

	return nil
}
