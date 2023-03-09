// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package config

import (
	"fmt"
	"testing"

	"github.com/caarlos0/env/v7"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/ltoml"
)

func TestSetGlobalConfig(t *testing.T) {
	b := &BrokerBase{}
	SetGlobalBrokerConfig(b)
	assert.Equal(t, b, GlobalBrokerConfig())

	s := &StorageBase{}
	SetGlobalStorageConfig(s)
	assert.Equal(t, s, GlobalStorageConfig())
}

func TestLoadAndSetBrokerConfig(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(cfg *Broker)
		wantErr bool
	}{
		{
			name: "load config failure",
			prepare: func(_ *Broker) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "load env failure",
			prepare: func(_ *Broker) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				envParseFn = func(v interface{}, opts ...env.Options) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "valid coordinator failure",
			prepare: func(cfg *Broker) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.Coordinator.Namespace = ""
			},
			wantErr: true,
		},
		{
			name: "valid broker failure",
			prepare: func(cfg *Broker) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.BrokerBase.HTTP.Port = 0
			},
			wantErr: true,
		},
		{
			name: "load and set cfg success",
			prepare: func(_ *Broker) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				loadConfigFn = ltoml.LoadConfig
				envParseFn = env.Parse
			}()
			cfg := &Broker{
				Coordinator: *NewDefaultCoordinator(),
				Query:       *NewDefaultQuery(),
				BrokerBase:  *NewDefaultBrokerBase(),
			}
			if tt.prepare != nil {
				tt.prepare(cfg)
			}
			err := LoadAndSetBrokerConfig("test", "broker.toml", cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAndSetBrokerConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAdnSetRootConfig(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(cfg *Root)
		wantErr bool
	}{
		{
			name: "load config failure",
			prepare: func(_ *Root) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "load env failure",
			prepare: func(_ *Root) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				envParseFn = func(v interface{}, opts ...env.Options) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "valid coordinator failure",
			prepare: func(cfg *Root) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.Coordinator.Namespace = ""
			},
			wantErr: true,
		},
		{
			name: "load and set cfg success",
			prepare: func(_ *Root) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				loadConfigFn = ltoml.LoadConfig
				envParseFn = env.Parse
			}()
			cfg := &Root{
				Coordinator: *NewDefaultCoordinator(),
				Query:       *NewDefaultQuery(),
			}
			if tt.prepare != nil {
				tt.prepare(cfg)
			}
			err := LoadAndSetRootConfig("test", "storage.toml", cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAndSetRootConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndSetStorageConfig(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(cfg *Storage)
		wantErr bool
	}{
		{
			name: "load config failure",
			prepare: func(_ *Storage) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "load env failure",
			prepare: func(_ *Storage) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				envParseFn = func(v interface{}, opts ...env.Options) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "valid coordinator failure",
			prepare: func(cfg *Storage) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.Coordinator.Namespace = ""
			},
			wantErr: true,
		},
		{
			name: "valid storage failure",
			prepare: func(cfg *Storage) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.StorageBase.GRPC.Port = 0
			},
			wantErr: true,
		},
		{
			name: "load and set cfg success",
			prepare: func(_ *Storage) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				loadConfigFn = ltoml.LoadConfig
				envParseFn = env.Parse
			}()
			cfg := &Storage{
				Coordinator: *NewDefaultCoordinator(),
				Query:       *NewDefaultQuery(),
				StorageBase: *NewDefaultStorageBase(),
			}
			if tt.prepare != nil {
				tt.prepare(cfg)
			}
			err := LoadAndSetStorageConfig("test", "storage.toml", cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAndSetStorageConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadAndSetStandaloneConfig(t *testing.T) {
	cases := []struct {
		name    string
		prepare func(cfg *Standalone)
		wantErr bool
	}{
		{
			name: "load config failure",
			prepare: func(_ *Standalone) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "load env failure",
			prepare: func(_ *Standalone) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				envParseFn = func(v interface{}, opts ...env.Options) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "valid coordinator failure",
			prepare: func(cfg *Standalone) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.Coordinator.Namespace = ""
			},
			wantErr: true,
		},
		{
			name: "valid broker failure",
			prepare: func(cfg *Standalone) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.BrokerBase.HTTP.Port = 0
			},
			wantErr: true,
		},
		{
			name: "valid storage failure",
			prepare: func(cfg *Standalone) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
				cfg.StorageBase.GRPC.Port = 0
			},
			wantErr: true,
		},
		{
			name: "load and set cfg success",
			prepare: func(_ *Standalone) {
				loadConfigFn = func(cfgPath, defaultCfgPath string, v interface{}) error {
					return nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				loadConfigFn = ltoml.LoadConfig
				envParseFn = env.Parse
			}()
			cfg := &Standalone{
				Coordinator: *NewDefaultCoordinator(),
				Query:       *NewDefaultQuery(),
				StorageBase: *NewDefaultStorageBase(),
				BrokerBase:  *NewDefaultBrokerBase(),
			}
			if tt.prepare != nil {
				tt.prepare(cfg)
			}
			err := LoadAndSetStandAloneConfig("test", "standalone.toml", cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAndSetStandAloneConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
