package standalone

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/server"
)

var testPath = "./test_data"

func TestRuntime_Run(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	cfg := config.NewDefaultStandaloneCfg()
	cfg.Storage.Engine.Dir = testPath
	standalone := NewStandaloneRuntime(cfg)
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
