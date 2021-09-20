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
	"os"
	"os/signal"
	"syscall"
)

// newCtxWithSignals returns a context which will can be canceled by sending signal.
func newCtxWithSignals() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case sig := <-c:
				// preventing exiting when receiving SigHup Signal
				if sig == syscall.SIGHUP {
					continue
				}
				return
			}
		}
	}()
	return ctx
}

// newSigHupCh returns a channel for handling sigHup signal,
// which is used for config reloading.
func newSigHupCh() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	return ch
}
