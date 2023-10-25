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
	"go.uber.org/zap"

	"github.com/lindb/common/pkg/logger"
)

var (
	initLoggerFn = logger.InitLogger
)

const (
	AccessLogFileName  = "access.log"
	SlowSQLLogFileName = "show_sql.log"

	SlowSQLModule = "SlowSQL"
)

// InitAccessLogger initializes a zap logger for access log.
func InitAccessLogger(cfg logger.Setting, fileName string) error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = logger.SimpleTimeEncoder
	encoderConfig.EncodeLevel = logger.SimpleAccessLevelEncoder

	log, err := initLoggerFn(fileName, cfg, &encoderConfig)
	if err != nil {
		return err
	}
	logger.RegisterLogger(logger.AccessLogModule, log, true)
	return nil
}

// InitSlowSQLLogger initializes a zap logger for slow sql log.
func InitSlowSQLLogger(cfg logger.Setting, fileName string) error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = logger.SimpleTimeEncoder
	encoderConfig.EncodeLevel = logger.SimpleLevelEncoder
	encoderConfig.LevelKey = ""
	encoderConfig.TimeKey = ""

	log, err := initLoggerFn(fileName, cfg, &encoderConfig)
	if err != nil {
		return err
	}
	logger.RegisterLogger(SlowSQLModule, log, true)
	return nil
}

// InitLogger initializes a zap logger for default log.
func InitLogger(cfg logger.Setting, fileName string) error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = logger.SimpleTimeEncoder
	encoderConfig.EncodeLevel = logger.SimpleLevelEncoder

	log, err := initLoggerFn(fileName, cfg, &encoderConfig, zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	logger.DefaultLogger.Store(log)
	return nil
}
