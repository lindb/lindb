package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// Logger is wrapper for zap logger with module, it is singleton.
type Logger struct {
	module string
	log    *zap.Logger
}

// GetLogger return logger with module name
func GetLogger(module string) *Logger {
	return &Logger{
		module: module,
		log:    getLogger(),
	}
}

func getLogger() *zap.Logger {
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

func (c *Config) New() (*zap.Logger, error) {
	//TODO ?????
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

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(l.formatMsg(msg), fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(l.formatMsg(msg), fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(l.formatMsg(msg), fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(l.formatMsg(msg), fields...)
}

// formatMsg formats msg using module name
func (l *Logger) formatMsg(msg string) string {
	return fmt.Sprintf("[%s]:%s", l.module, msg)
}

// String constructs a field with the given key and value.
func String(key string, val string) zap.Field {
	return zap.Field{Key: key, Type: zapcore.StringType, String: val}
}

// Error is shorthand for the common idiom NamedError("error", err).
func Error(err error) zap.Field {
	return zap.NamedError("error", err)
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Uint16Type, Integer: int64(val)}
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Uint32Type, Integer: int64(val)}
}

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack() zap.Field {
	return zap.Stack("stack")
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) zap.Field {
	return zap.Reflect(key, val)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Int32Type, Integer: int64(val)}
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) zap.Field {
	return zap.Field{Key: key, Type: zapcore.Int64Type, Integer: val}
}
