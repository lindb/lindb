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
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/lindb/lindb/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	isTerminal = IsTerminal(os.Stdout)
	// max length of all modules
	maxModuleNameLen uint32
	lindLogger       atomic.Value
	accessLogger     atomic.Value
	// uninitialized logger for default usage
	defaultLogger = newDefaultLogger()
	// RunningAtomicLevel supports changing level on the fly
	RunningAtomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
)

func init() {
	// get log level from evn
	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		var zapLevel zapcore.Level
		if err := zapLevel.Set(level); err == nil {
			RunningAtomicLevel.SetLevel(zapLevel)
		}
	}
}

const (
	accessLogFileName = "access.log"
)

// GetLogger return logger with module name
func GetLogger(module, role string) *Logger {
	length := len(module)
	for {
		currentMaxModuleLen := atomic.LoadUint32(&maxModuleNameLen)
		if uint32(length) <= currentMaxModuleLen {
			break
		}
		if atomic.CompareAndSwapUint32(&maxModuleNameLen, currentMaxModuleLen, uint32(length)) {
			break
		}
	}
	return &Logger{
		module: module,
		role:   role,
	}
}

// newDefaultLogger creates a default logger for uninitialized usage
func newDefaultLogger() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = SimpleTimeEncoder
	encoderConfig.EncodeLevel = SimpleLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		os.Stdout,
		RunningAtomicLevel)
	return zap.New(core)
}

// InitLogger initializes a zap logger from user config
func InitLogger(cfg config.Logging, fileName string) error {
	if err := initLogger(fileName, cfg); err != nil {
		return err
	}
	if err := initLogger(accessLogFileName, cfg); err != nil {
		return err
	}
	return nil
}

// initLogger initializes a zap logger for different module
func initLogger(logFilename string, cfg config.Logging) error {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, logFilename),
		MaxSize:    int(cfg.MaxSize),
		MaxBackups: int(cfg.MaxBackups),
		MaxAge:     int(cfg.MaxAge),
	})
	// check if it is terminal
	if isTerminal {
		w = os.Stdout
	}
	// parse logging level
	if err := RunningAtomicLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		return err
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = SimpleTimeEncoder
	switch {
	case logFilename == accessLogFileName:
		encoderConfig.EncodeLevel = SimpleAccessLevelEncoder
	default:
		encoderConfig.EncodeLevel = SimpleLevelEncoder
	}
	// check format
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		RunningAtomicLevel)
	switch {
	case logFilename == accessLogFileName:
		accessLogger.Store(zap.New(core))
	default:
		lindLogger.Store(zap.New(core))
	}
	return nil
}
