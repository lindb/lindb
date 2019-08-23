package logger

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Logger(t *testing.T) {
	logger1 := GetLogger("pkg/logger", "test")
	logger1.Warn("warn for test", String("count", "1"), Reflect("v1", map[string]string{"a": "1"}))
	logger1.Info("info for test", Uint16("value", 1), Int32("v1", 2),
		Int64("v2", 2), Any("v3", 3))
	logger1.Debug("debug for test", Uint32("value", 2))
	logger1.Error("error for test", Error(fmt.Errorf("error")))

	logger2 := New()
	assert.NotNil(t, logger2)

	logger3 := GetLogger("pkg/logger", "")
	logger3.Error("error test")
}

func Test_Logger_Stack(t *testing.T) {
	panicFunc := func() {
		defer func() {
			if r := recover(); r != nil {
				GetLogger("pkg/logger", "test-panic").log.Panic("panic stack", Stack())
			}
		}()
		panic("test-panic")
	}
	assert.Panics(t, panicFunc)
}

func Test_IsTerminal(t *testing.T) {
	assert.False(t, IsTerminal(os.Stdout))

	w := bytes.NewBuffer(nil)
	assert.False(t, IsTerminal(w))
}
