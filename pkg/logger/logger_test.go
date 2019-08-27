package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/lindb/lindb/config"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func Test_Logger(t *testing.T) {
	logger1 := GetLogger("pkg/logger", "test")
	RunningAtomicLevel.SetLevel(zapcore.DebugLevel)

	fmt.Println(White.Add("white"))
	logger1.Warn("warn for test", String("count", "1"), Reflect("v1", map[string]string{"a": "1"}))
	logger1.Info("info for test", Uint16("value", 1), Int32("v1", 2),
		Int64("v2", 2), Any("v3", 3))
	logger1.Debug("debug for test", Uint32("value", 2))
	logger1.Error("error for test", Error(fmt.Errorf("error")))

	assert.NotNil(t, defaultLogger)

	logger3 := GetLogger("pkg/logger", "")
	logger3.Error("error test")
}

func Test_Logger_Stack(t *testing.T) {
	panicFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				GetLogger("pkg/logger", "test-panic").
					getInitializedOrDefaultLogger().Panic("panic stack", Stack())
			}
		}()
		panic("test-panic")
	}
	assert.Panics(t, panicFunc)
}

func Test_IsTerminal(t *testing.T) {
	assert.False(t, IsTerminal(os.Stdout))
}

func Test_InitLogger(t *testing.T) {
	assert.NotNil(t, GetLogger("test", "test").getInitializedOrDefaultLogger())

	cfg1 := config.Logging{Level: "LLL"}
	assert.NotNil(t, InitLogger(cfg1))

	cfg2 := config.NewDefaultLoggingCfg()
	assert.Nil(t, InitLogger(cfg2))
	thisLogger := GetLogger("test", "test")
	assert.NotNil(t, thisLogger.getInitializedOrDefaultLogger())
	assert.NotNil(t, thisLogger.getInitializedOrDefaultLogger())

	cfg3 := config.Logging{Level: "info"}
	assert.Nil(t, InitLogger(cfg3))

	cfg4 := config.Logging{Level: "debug"}
	assert.Nil(t, InitLogger(cfg4))
}
