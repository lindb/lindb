package standalone

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/server"
)

func TestRuntime_Run(t *testing.T) {
	standalone := NewStandaloneRuntime(config.NewDefaultStandaloneCfg())
	err := standalone.Run()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, server.Running, standalone.State())
	err = standalone.Stop()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, server.Terminated, standalone.State())
	assert.Equal(t, "standalone", standalone.Name())
}
