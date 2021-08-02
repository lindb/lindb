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
)

var (
	// defaultParentDir is the default directory of lindb
	defaultParentDir = "/tmp/lindb"
)

// Logging represents a logging configuration
type Logging struct {
	Dir        string `toml:"dir"`
	Level      string `toml:"level"`
	MaxSize    uint16 `toml:"maxsize"`
	MaxBackups uint16 `toml:"maxbackups"`
	MaxAge     uint16 `toml:"maxage"`
}

// TOML returns Logging's toml config string
func (l *Logging) TOML() string {
	return fmt.Sprintf(`
[logging]
  ## Dir is the output directory for log-files
  dir = "%s"

  ## Determine which level of logs will be emitted.
  ## error, warn, info, and debug are available
  level = "%s"

  ## MaxSize is the maximum size in megabytes of the log file before it gets
  ## rotated. It defaults to 100 megabytes.
  maxsize = %d

  ## MaxBackups is the maximum number of old log files to retain.  The default
  ## is to retain all old log files (though MaxAge may still cause them to get
  ## deleted.)
  maxbackups = %d

  ## MaxAge is the maximum number of days to retain old log files based on the
  ## timestamp encoded in their filename.  Note that a day is defined as 24
  ## hours and may not exactly correspond to calendar days due to daylight
  ## savings, leap seconds, etc. The default is not to remove old log files
  ## based on age.
  maxage = %d`,
		l.Dir,
		l.Level,
		l.MaxSize,
		l.MaxBackups,
		l.MaxAge)
}

// NewDefaultLogging returns a new default logging config
func NewDefaultLogging() *Logging {
	return &Logging{
		Dir:        filepath.Join(defaultParentDir, "log"),
		Level:      "info",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30}
}
