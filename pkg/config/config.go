package config

import (
	"github.com/BurntSushi/toml"
)

func Parse(configFile string, v interface{}) {
	if _, err := toml.DecodeFile(configFile, v); err != nil {
		panic(err)
	}
}
