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

package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"

	"go.uber.org/automaxprocs/maxprocs"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/server"
	"github.com/lindb/lindb/pkg/logger"
)

// serveStandalone runs the cluster as standalone mode
func run(ctx context.Context, service server.Service, reloadConfigFunc func() error) error {
	printLogoWhenIsTty()

	var mainLogger = logger.GetLogger("CMD", "Main")

	mainLogger.Info(fmt.Sprintf("Lind running as %s with PID: %d (pprof: %v)",
		service.Name(), os.Getpid(), pprof))

	config.Profile = pprof
	config.Doc = doc

	// auto set maxprocs
	_, _ = maxprocs.Set(maxprocs.Logger(func(s string, i ...interface{}) {
		mainLogger.Info(fmt.Sprintf(s, i))
	}))
	// start service
	if err := service.Run(); err != nil {
		return fmt.Errorf("run service[%s] error:%s", service.Name(), err)
	}

	signUpCh := newSigHupCh()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-signUpCh:
				mainLogger.Info("received SIGHUP signal, reloading config...")
				if err := reloadConfigFunc(); err != nil {
					mainLogger.Error("failed to reload config", logger.Error(err))
				} else {
					mainLogger.Info("reload config successfully")
				}
			}
		}
	}()

	// waiting system exit signal
	<-ctx.Done()

	// stop service
	service.Stop()

	return nil
}
