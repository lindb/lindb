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

package lind

import (
	"context"
	"fmt"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/server"

	"github.com/gin-gonic/gin"
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
		gin.SetMode(gin.DebugMode)
		go func() {
			if err := http.ListenAndServe(":6060", nil); err != nil {
				mainLogger.Error("close debug http listener with err", logger.Error(err))
			}
		}()
		mainLogger.Info("pprof listening on 6060")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// start service
	if err := service.Run(); err != nil {
		return fmt.Errorf("run service[%s] error:%s", service.Name(), err)
	}

	// waiting system exit signal
	<-ctx.Done()

	// stop service
	service.Stop()

	return nil
}
