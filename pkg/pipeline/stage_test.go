package pipeline

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStage_Shutdown(t *testing.T) {
	stage := NewStage(&ConfigTest{
		taskSize: 2,
	})
	stage.Shutdown()
	time.Sleep(time.Second)
	for i := 0; i < len(stage.runs); i++ {
		assert.Equal(t, int32(1), atomic.LoadInt32(&stage.runs[i].closed))
	}
}
