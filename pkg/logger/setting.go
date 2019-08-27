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
	logger           atomic.Value
	// uninitialized logger for default usage
	defaultLogger = newDefaultLogger()
	// RunningAtomicLevel supports changing level on the fly
	RunningAtomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
)

const (
	lindLogFilename = "lind.log"
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
func InitLogger(cfg config.Logging) error {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, lindLogFilename),
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
	encoderConfig.EncodeLevel = SimpleLevelEncoder
	// check format
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		RunningAtomicLevel)
	logger.Store(zap.New(core))
	return nil
}
