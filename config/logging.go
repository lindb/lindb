package config

import "path/filepath"

var (
	// defaultParentDir is the default directory of lindb
	defaultParentDir = "/tmp/lindb"
)

// Logging represents a logging configuration
type Logging struct {
	// Dir is the output directory for log-files
	Dir string `toml:"dir"`
	// Determine which level of logs will be emitted.
	// error, warn, info, and debug are available
	Level string `toml:"level"`
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize uint16 `toml:"maxsize"`
	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups uint16 `toml:"maxbackups"`
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge uint16 `toml:"maxage"`
}

// NewDefaultLoggingCfg returns a new default logging config
func NewDefaultLoggingCfg() Logging {
	return Logging{
		Dir:        filepath.Join(defaultParentDir, "log"),
		Level:      "info",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30}
}
