package config

import (
	"path/filepath"
	"testing"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/ltoml"

	"github.com/stretchr/testify/assert"
)

var testPath = "./tmp"

func Test_NewConfig(t *testing.T) {
	_ = fileutil.MkDirIfNotExist(testPath)
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	// validate broker config
	brokerCfgPath := filepath.Join(testPath, "broker.toml")
	assert.Nil(t, ltoml.WriteConfig(brokerCfgPath, NewDefaultBrokerTOML()))
	var brokerCfg Broker
	assert.Nil(t, ltoml.DecodeToml(brokerCfgPath, &brokerCfg))
	assert.Equal(t, brokerCfg.BrokerBase, *NewDefaultBrokerBase())
	assert.Equal(t, brokerCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, brokerCfg.Monitor, *NewDefaultMonitor())

	// validate storage config
	storageCfgPath := filepath.Join(testPath, "storage.toml")
	assert.Nil(t, ltoml.WriteConfig(storageCfgPath, NewDefaultStorageTOML()))
	var storageCfg Storage
	assert.Nil(t, ltoml.DecodeToml(storageCfgPath, &storageCfg))
	assert.Equal(t, storageCfg.StorageBase, *NewDefaultStorageBase())
	assert.Equal(t, storageCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, storageCfg.Monitor, *NewDefaultMonitor())

	// validate standalone config
	standaloneCfgPath := filepath.Join(testPath, "standalone.toml")
	assert.Nil(t, ltoml.WriteConfig(standaloneCfgPath, NewDefaultStandaloneTOML()))
	var standaloneCfg Standalone
	assert.Nil(t, ltoml.DecodeToml(standaloneCfgPath, &standaloneCfg))
	assert.Equal(t, standaloneCfg.BrokerBase, *NewDefaultBrokerBase())
	assert.Equal(t, standaloneCfg.StorageBase, *NewDefaultStorageBase())
	assert.Equal(t, standaloneCfg.Logging, *NewDefaultLogging())
	assert.Equal(t, standaloneCfg.Monitor, *NewDefaultMonitor())
}

func Test_ReplicationChannel_SegmentFileSizeInBytes(t *testing.T) {
	var rc ReplicationChannel
	assert.Equal(t, 1024*1024, rc.SegmentFileSizeInBytes())
	assert.Zero(t, rc.BufferSizeInBytes())
	rc.SegmentFileSize = 10
	assert.Equal(t, 10*1024*1024, rc.SegmentFileSizeInBytes())
	rc.SegmentFileSize = 10000
	assert.Equal(t, 1024*1024*1024, rc.SegmentFileSizeInBytes())
}
