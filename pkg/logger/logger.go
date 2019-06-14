package logger

import (
	"io"
	"github.com/mattn/go-isatty"
	"sync"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	"os"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		logger = New()
	})
	return logger
}

func New() *zap.Logger {
	config := NewConfig()
	l, _ := config.New()
	return l
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func (c *Config) New() (*zap.Logger, error) {
	//w := zapcore.AddSync(&lumberjack.Logger{
	//	Filename:   "/var/log/myapp/foo.log",
	//	MaxSize:    500, // megabytes
	//	MaxBackups: 3,
	//	MaxAge:     28, // days
	//})
	//core := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
	//	w,
	//	zap.InfoLevel,
	//)
	//logger := zap.New(core)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		os.Stdout,
		c.Level,
	)

	logger := zap.New(core)
	logger.Info("logger init success")
	return logger, nil
}

// IsTerminal checks if w is a file and whether it is an interactive terminal session.
func IsTerminal(w io.Writer) bool {
	if f, ok := w.(interface {
		Fd() uintptr
	}); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}
