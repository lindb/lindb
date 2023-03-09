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
	"strings"

	"github.com/lindb/lindb/pkg/ltoml"
)

var (
	// defaultParentDir is the default directory of lindb
	defaultParentDir = filepath.Join(".", "data")
)

// Logging represents a logging configuration
type Logging struct {
	Dir        string     `env:"DIR" toml:"dir"`
	Level      string     `env:"LEVEL" toml:"level"`
	MaxSize    ltoml.Size `env:"MAX_SIZE" toml:"maxsize"`
	MaxBackups uint16     `env:"MAX_BACKUPS" toml:"maxbackups"`
	MaxAge     uint16     `env:"MAX_AGE" toml:"maxage"`
}

// TOML returns Logging's toml config string
func (l *Logging) TOML() string {
	return fmt.Sprintf(`
## logging related configuration.
[logging]
## Dir is the output directory for log-files
## Default: %s
## Env: LOGGING_DIR
dir = "%s"
## Determine which level of logs will be emitted.
## error, warn, info, and debug are available
## Default: %s
## Env: LOGGING_LEVEL
level = "%s"
## MaxSize is the maximum size in megabytes of the log file before it gets rotated. 
## Default: %s
## Env: LOGGING_MAX_SIZE
maxsize = "%s"
## MaxBackups is the maximum number of old log files to retain. The default
## is to retain all old log files (though MaxAge may still cause them to get deleted.)
## Default: %d
## Env: LOGGING_MAX_BACKUPS
maxbackups = %d
## MaxAge is the maximum number of days to retain old log files based on the
## timestamp encoded in their filename.  Note that a day is defined as 24 hours
## and may not exactly correspond to calendar days due to daylight savings, leap seconds, etc.
## The default is not to remove old log files based on age.
## Default: %d
## Env: LOGGING_MAS_AGE
maxage = %d`,
		strings.ReplaceAll(l.Dir, "\\", "\\\\"),
		strings.ReplaceAll(l.Dir, "\\", "\\\\"),
		l.Level,
		l.Level,
		l.MaxSize,
		l.MaxSize,
		l.MaxBackups,
		l.MaxBackups,
		l.MaxAge,
		l.MaxAge,
	)
}

// NewDefaultLogging returns a new default logging config
func NewDefaultLogging() *Logging {
	return &Logging{
		Dir:        filepath.Join(defaultParentDir, "log"),
		Level:      "info",
		MaxSize:    ltoml.Size(100 * 1024 * 1024),
		MaxBackups: 3,
		MaxAge:     7,
	}
}
