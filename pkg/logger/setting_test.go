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

package logger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/lindb/common/pkg/logger"
)

func TestLogger_InitAccessLogger(t *testing.T) {
	defer func() {
		initLoggerFn = logger.InitLogger
	}()

	t.Run("init ok", func(t *testing.T) {
		assert.NoError(t, InitAccessLogger(logger.Setting{}, AccessLogFileName))
	})
	t.Run("init fail", func(t *testing.T) {
		initLoggerFn = func(_ string, _ logger.Setting, _ *zapcore.EncoderConfig, _ ...zap.Option) (*zap.Logger, error) {
			return nil, fmt.Errorf("err")
		}
		assert.Error(t, InitAccessLogger(logger.Setting{}, AccessLogFileName))
	})
}

func TestLogger_InitSlowSQLLogger(t *testing.T) {
	defer func() {
		initLoggerFn = logger.InitLogger
	}()

	t.Run("init ok", func(t *testing.T) {
		assert.NoError(t, InitSlowSQLLogger(logger.Setting{}, SlowSQLLogFileName))
	})
	t.Run("init fail", func(t *testing.T) {
		initLoggerFn = func(_ string, _ logger.Setting, _ *zapcore.EncoderConfig, _ ...zap.Option) (*zap.Logger, error) {
			return nil, fmt.Errorf("err")
		}
		assert.Error(t, InitSlowSQLLogger(logger.Setting{}, SlowSQLLogFileName))
	})
}

func TestLogger_InitLogger(t *testing.T) {
	defer func() {
		initLoggerFn = logger.InitLogger
	}()

	t.Run("init ok", func(t *testing.T) {
		assert.NoError(t, InitLogger(logger.Setting{}, SlowSQLLogFileName))
	})
	t.Run("init fail", func(t *testing.T) {
		initLoggerFn = func(_ string, _ logger.Setting, _ *zapcore.EncoderConfig, _ ...zap.Option) (*zap.Logger, error) {
			return nil, fmt.Errorf("err")
		}
		assert.Error(t, InitLogger(logger.Setting{}, SlowSQLLogFileName))
	})
}
