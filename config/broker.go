package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/lindb/lindb/pkg/ltoml"
)

// HTTP represents a HTTP level configuration of broker.
type HTTP struct {
	Port uint16 `toml:"port"`
}

func (h *HTTP) TOML() string {
	return fmt.Sprintf(`
    ## which port broker's HTTP Server is listening on 
    port = %d`,
		h.Port,
	)
}

// User represents user model
type User struct {
	UserName string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

func (u *User) TOML() string {
	return fmt.Sprintf(`
    ## admin user setting
    username = "%s"
    password = "%s"`,
		u.UserName,
		u.Password)
}

type TCP struct {
	Port uint16 `toml:"port"`
}

func (t *TCP) TOML() string {
	return fmt.Sprintf(`
    port = %d`,
		t.Port)
}

// ReplicationChannel represents config for data replication in broker.
type ReplicationChannel struct {
	Dir                string         `toml:"dir"`
	SegmentFileSize    uint16         `toml:"segment-file-size"`
	RemoveTaskInterval ltoml.Duration `toml:"remove-task-interval"`
	ReportInterval     ltoml.Duration `toml:"report-interval"` // replicator state report interval
	CheckFlushInterval ltoml.Duration `toml:"check-flush-interval"`
	FlushInterval      ltoml.Duration `toml:"flush-interval"`
	BufferSize         uint16         `toml:"buffer-size"`
}

func (rc *ReplicationChannel) SegmentFileSizeInBytes() int {
	if rc.SegmentFileSize <= 1 {
		return 1024 * 1024 // 1MB
	}
	if rc.SegmentFileSize >= 1024 {
		return 1024 * 1024 * 1024 // 1GB
	}
	return int(rc.SegmentFileSize) * 1024 * 1024
}

func (rc *ReplicationChannel) BufferSizeInBytes() int {
	return int(rc.BufferSize) * 1024
}

func (rc *ReplicationChannel) TOML() string {
	return fmt.Sprintf(`
    ## WAL mmaped log directory
    dir = "%s"
    
    ## segment-file-size is the maximum size in megabytes of the segment file before a new
    ## file is created. It defaults to 128 megabytes, available size is in [1MB, 1GB]
    segment-file-size = %d
	
    ## interval for how often a new segment will be created
    remove-task-interval = "%s"

    ## replicator state report interval
    report-interval = "%s"

    ## interval for how often buffer will be checked if it's available to flush
    check-flush-interval = "%s"

    ## interval for how often data will be flushed if data not exceeds the buffer-size
    flush-interval = "%s"

    ## will flush if this size of data in kegabytes get buffered
    buffer-size = %d`,
		rc.Dir,
		rc.SegmentFileSize,
		rc.RemoveTaskInterval.String(),
		rc.ReportInterval.String(),
		rc.CheckFlushInterval.String(),
		rc.FlushInterval.String(),
		rc.BufferSize,
	)
}

// BrokerBase represents a broker configuration
type BrokerBase struct {
	Coordinator        RepoState          `toml:"coordinator"`
	Query              Query              `toml:"query"`
	HTTP               HTTP               `toml:"http"`
	User               User               `toml:"user"`
	GRPC               GRPC               `toml:"grpc"`
	TCP                TCP                `toml:"tcp"`
	ReplicationChannel ReplicationChannel `toml:"replication_channel"`
}

func (bb *BrokerBase) TOML() string {
	return fmt.Sprintf(`## Config for the Broker Node
[broker]
  [broker.coordinator]%s
  
  [broker.query]%s

  [broker.http]%s
	
  [broker.user]%s

  [broker.grpc]%s

  [broker.tcp]%s

  [broker.replication_channel]%s`,
		bb.Coordinator.TOML(),
		bb.Query.TOML(),
		bb.HTTP.TOML(),
		bb.User.TOML(),
		bb.GRPC.TOML(),
		bb.TCP.TOML(),
		bb.ReplicationChannel.TOML(),
	)
}

func NewDefaultBrokerBase() *BrokerBase {
	return &BrokerBase{
		HTTP: HTTP{
			Port: 9000,
		},
		GRPC: GRPC{
			Port: 9001,
		},
		TCP: TCP{
			Port: 9002,
		},
		Coordinator: RepoState{
			Namespace:   "/lindb/broker",
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: ltoml.Duration(time.Second * 5),
		},
		User: User{
			UserName: "admin",
			Password: "admin123",
		},
		ReplicationChannel: ReplicationChannel{
			Dir:                filepath.Join(defaultParentDir, "broker/replication"),
			SegmentFileSize:    128,
			RemoveTaskInterval: ltoml.Duration(time.Minute),
			CheckFlushInterval: ltoml.Duration(time.Second),
			FlushInterval:      ltoml.Duration(5 * time.Second),
			BufferSize:         128,
		},
		Query: *NewDefaultQuery(),
	}
}

// Broker represents a broker configuration with common settings
type Broker struct {
	BrokerBase BrokerBase `toml:"broker"`
	Monitor    Monitor    `toml:"monitor"`
	Logging    Logging    `toml:"logging"`
}

// NewDefaultBrokerTOML creates broker default toml config
func NewDefaultBrokerTOML() string {
	return fmt.Sprintf(`%s

%s

%s`,
		NewDefaultBrokerBase().TOML(),
		NewDefaultMonitor().TOML(),
		NewDefaultLogging().TOML(),
	)
}
