package ltoml

import (
	"os"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"

	"github.com/stretchr/testify/assert"
)

type TestCfg struct {
	Path string `toml:"path"`
}

var cfgFile = "./test.test"
var defaultCfgFile = "./test.test"

func TestLoadConfig(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(cfgFile)
	}()
	assert.NotNil(t, LoadConfig(cfgFile, defaultCfgFile, &TestCfg{}))

	f, err := os.Create(cfgFile)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("Hello World")
	assert.NotNil(t, LoadConfig(cfgFile, defaultCfgFile, &TestCfg{}))

	_ = EncodeToml(cfgFile, &TestCfg{Path: "/data/path"})
	cfg := TestCfg{}
	err = LoadConfig(cfgFile, defaultCfgFile, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)

	err = LoadConfig("", defaultCfgFile, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)
}
