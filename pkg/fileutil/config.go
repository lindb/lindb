package fileutil

import (
	"fmt"
)

// LoadConfig loads config from file, if fail return err
func LoadConfig(cfgPath, defaultCfgPath string, v interface{}) error {
	if cfgPath == "" {
		cfgPath = defaultCfgPath
	}
	if !Exist(cfgPath) {
		return fmt.Errorf("config file doesn't exist`")
	}

	if err := DecodeToml(cfgPath, v); err != nil {
		return fmt.Errorf("decode config file error:%s", err)
	}
	return nil
}
