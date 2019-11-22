package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// Replication represents replication config
type Replication struct {
	Dir string `toml:"dir"`
}

func (r *Replication) TOML() string {
	return fmt.Sprintf(`
    ## Where the WAL log is stored
    dir = "%s"`,
		r.Dir)
}

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
	Coordinator RepoState   `toml:"coordinator"`
	GRPC        GRPC        `toml:"grpc"`
	TSDB        TSDB        `toml:"tsdb"`
	Replication Replication `toml:"replication"`
	Query       Query       `toml:"query"`
}

// TOML returns StorageBase's toml config string
func (s *StorageBase) TOML() string {
	return fmt.Sprintf(`## Config for the Storage Node
[storage]
  [storage.coordinator]%s
  
  [storage.query]%s
  
  [storage.grpc]%s

  [storage.tsdb]%s
	
  [storage.replication]%s
`,
		s.Coordinator.TOML(),
		s.Query.TOML(),
		s.GRPC.TOML(),
		s.TSDB.TOML(),
		s.Replication.TOML(),
	)
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
			DialTimeout: ltoml.Duration(time.Second * 5)},
		GRPC: GRPC{
			Port: 2891,
			TTL:  ltoml.Duration(time.Second)},
		TSDB: TSDB{
			Dir: filepath.Join(defaultParentDir, "storage/data")},
		Replication: Replication{
			Dir: filepath.Join(defaultParentDir, "storage/replication")},
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
