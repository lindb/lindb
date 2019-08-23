package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	isatty "github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// max length of all modules
	maxModuleNameLen uint32
	logger           *zap.Logger
	once             sync.Once
)

// SimpleTimeEncoder serializes a time.Time to a simplified format without timezone
func SimpleTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// SimpleColorLevelEncoder serializes a Level to a lowercase string. For example,
// InfoLevel is serialized to "info".
func SimpleLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(LevelString(l))
}

// LevelString returns a upper-case ASCII representation of the log level.
func LevelString(l zapcore.Level) string {
	switch l {
	case zapcore.DebugLevel:
		return Magenta.Add("DEBUG")
	case zapcore.InfoLevel:
		return Green.Add("INFO")
	case zapcore.WarnLevel:
		return Yellow.Add("WARN")
	case zapcore.ErrorLevel:
		return Red.Add("ERROR")
	default:
		return Red.Add(fmt.Sprintf("LEVEL(%d)", l))
	}
}

// Logger is wrapper for zap logger with module, it is singleton.
type Logger struct {
	module string
	role   string
	log    *zap.Logger
}

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
		log:    getLogger(),
	}
}

// getLogger returns the zap logger
func getLogger() *zap.Logger {
	once.Do(func() {
		logger = New()
	})
	return logger
}

// New creates a zap logger
func New() *zap.Logger {
	config := NewConfig()
	l, _ := config.New()
	return l
}

// New creates a zap logger based on user config
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
	encoderConfig.EncodeTime = SimpleTimeEncoder
	encoderConfig.EncodeLevel = SimpleLevelEncoder
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
	if runtime.GOOS == "windows" {
		return false
	}
	if f, ok := w.(interface {
		Fd() uintptr
	}); ok {
		return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
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
	moduleName := Cyan.Add(fmt.Sprintf("[%*s]", atomic.LoadUint32(&maxModuleNameLen), l.module))
	if l.role == "" {
		return fmt.Sprintf("%s: %s",
			moduleName, msg)
	}
	return fmt.Sprintf("%s [%s]: %s",
		moduleName, l.role, msg)
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
