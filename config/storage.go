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
	"path/filepath"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// TSDB represents the tsdb configuration
type TSDB struct {
	Dir string `toml:"dir"`
}

func (t *TSDB) TOML() string {
	return fmt.Sprintf(`
    ## where the tsdb data is stored
    dir = "%s"`,
		t.Dir,
	)
}

// StorageBase represents a storage configuration
type StorageBase struct {
	Coordinator RepoState `toml:"coordinator"`
	GRPC        GRPC      `toml:"grpc"`
	TSDB        TSDB      `toml:"tsdb"`
	Query       Query     `toml:"query"`
	Replica     Replica   `toml:"replica"`
}

// TOML returns StorageBase's toml config string
func (s *StorageBase) TOML() string {
	return fmt.Sprintf(`## Config for the Storage Node
[storage]
  [storage.coordinator]%s
  
  [storage.query]%s
  
  [storage.grpc]%s

  [storage.tsdb]%s
`,
		s.Coordinator.TOML(),
		s.Query.TOML(),
		s.GRPC.TOML(),
		s.TSDB.TOML(),
	)
}

// Replica represents config for data replication in storage.
type Replica struct {
	Dir                string         `toml:"dir"`
	DataSizeLimit      int64          `toml:"data-size-limit"`
	RemoveTaskInterval ltoml.Duration `toml:"remove-task-interval"`
	ReportInterval     ltoml.Duration `toml:"report-interval"` // replicator state report interval
	CheckFlushInterval ltoml.Duration `toml:"check-flush-interval"`
	FlushInterval      ltoml.Duration `toml:"flush-interval"`
	BufferSize         int            `toml:"buffer-size"`
}

func (rc *Replica) GetDataSizeLimit() int64 {
	if rc.DataSizeLimit <= 1 {
		return 1024 * 1024 // 1MB
	}
	if rc.DataSizeLimit >= 1024 {
		return 1024 * 1024 * 1024 // 1GB
	}
	return rc.DataSizeLimit * 1024 * 1024
}

// Storage represents a storage configuration with common settings
type Storage struct {
	StorageBase StorageBase `toml:"storage"`
	Monitor     Monitor     `toml:"monitor"`
	Logging     Logging     `toml:"logging"`
}

// NewDefaultStorageBase returns a new default StorageBase struct
func NewDefaultStorageBase() *StorageBase {
	return &StorageBase{
		Coordinator: RepoState{
			Namespace:   "/lindb/storage",
			Endpoints:   []string{"http://localhost:2379"},
			Timeout:     ltoml.Duration(time.Second * 10),
			DialTimeout: ltoml.Duration(time.Second * 5),
		},
		GRPC: GRPC{
			Port: 2891,
			TTL:  ltoml.Duration(time.Second)},
		TSDB: TSDB{
			Dir: filepath.Join(defaultParentDir, "storage/data")},
		Query: *NewDefaultQuery(),
	}
}

// NewDefaultStorageTOML creates storage's default toml config
func NewDefaultStorageTOML() string {
	return fmt.Sprintf(`%s

%s

%s`,
		NewDefaultStorageBase().TOML(),
		NewDefaultMonitor().TOML(),
		NewDefaultLogging().TOML(),
	)
}
