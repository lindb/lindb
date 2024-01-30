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

package prometheus

import (
	"fmt"
	"time"

	"github.com/lindb/common/pkg/logger"
	"github.com/prometheus/prometheus/promql"
	"go.uber.org/zap"
)

type QueryOpts struct {
	opt promql.EngineOpts
}

func (o *QueryOpts) EnablePerStepStats() bool {
	return o.opt.EnablePerStepStats
}

func (o *QueryOpts) LookbackDelta() time.Duration {
	return o.opt.LookbackDelta
}

// Logger implements Prometheus's Logger
type Logger struct {
	logger logger.Logger
}

func (l *Logger) Log(keyvals ...interface{}) error {
	fields := make([]zap.Field, len(keyvals)/2)
	for i := 1; i < len(keyvals); i++ {
		key := fmt.Sprintf("%v", keyvals[i-1])
		fields = append(fields, logger.Any(key, keyvals[i]))
	}
	l.logger.Debug("prometheus", fields...)
	return nil
}

// NewEngine instantiates an instance of the Prometheus Engine.
func NewEngine() *promql.Engine {
	opt := promql.EngineOpts{
		Logger:               &Logger{logger: logger.GetLogger("prometheus", "engine")},
		MaxSamples:           1e5,
		Timeout:              time.Minute,
		LookbackDelta:        time.Minute * 5,
		EnableAtModifier:     true,
		EnableNegativeOffset: true,
		EnablePerStepStats:   true,
	}

	engine := promql.NewEngine(opt)

	return engine
}
