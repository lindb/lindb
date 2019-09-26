package lind

import (
	"context"
	"fmt"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"

	"go.uber.org/zap/zapcore"
)

// serveStandalone runs the cluster as standalone mode
func run(ctx context.Context, service server.Service) error {
	printLogoWhenIsTty()

	var mainLogger = logger.GetLogger("cmd", "Main")

	mainLogger.Info(fmt.Sprintf("Lind running as %s with PID: %d (debug: %v)",
		service.Name(), os.Getpid(), debug))
	// enabled debug log level
	if debug {
		logger.RunningAtomicLevel.SetLevel(zapcore.DebugLevel)
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
		mainLogger.Info("pprof listening on 6060")
	}

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
