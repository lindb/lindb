package option

import (
	"fmt"

	"github.com/lindb/lindb/pkg/timeutil"
)

// DatabaseOption represents a database option include shard ids and shard's option
type DatabaseOption struct {
	Interval string `toml:"interval" json:"interval,omitempty"` // write interval(the number of second)
	// rollup intervals(like seconds->minute->hour->day)
	Rollup []string `toml:"rollup" json:"rollup,omitempty"`

	// auto create namespace
	AutoCreateNS bool `toml:"autoCreateNS" json:"autoCreateNS,omitempty"`

	Behind string `toml:"behind" json:"behind,omitempty"` // allowed timestamp write behind
	Ahead  string `toml:"ahead" json:"ahead,omitempty"`   // allowed timestamp write ahead

	Index FlusherOption `toml:"index" json:"index,omitempty"` // index flusher option
	Data  FlusherOption `toml:"data" json:"data,omitempty"`   // data flusher data
}

// FlusherOption represents a flusher configuration for index and memory db
type FlusherOption struct {
	TimeThreshold int64 `toml:"timeThreshold" json:"timeThreshold"` // time level flush threshold
	SizeThreshold int64 `toml:"sizeThreshold" json:"sizeThreshold"` // size level flush threshold, unit(MB)
}

// Validate validates engine option if valid
func (e DatabaseOption) Validate() error {
	if err := validateInterval(e.Interval, true); err != nil {
		return err
	}
	for _, interval := range e.Rollup {
		if err := validateInterval(interval, true); err != nil {
			return err
		}
	}
	if err := validateInterval(e.Ahead, false); err != nil {
		return err
	}
	if err := validateInterval(e.Behind, false); err != nil {
		return err
	}
	var interval timeutil.Interval
	_ = interval.ValueOf(e.Interval)
	for _, intervalStr := range e.Rollup {
		var rollupInterval timeutil.Interval
		_ = rollupInterval.ValueOf(intervalStr)
		if interval.Int64() >= rollupInterval.Int64() {
			return fmt.Errorf("rollup interval must be large than write interval")
		}
	}
	return nil
}

// validateInterval checks interval string if valid
func validateInterval(intervalStr string, require bool) error {
	if !require && intervalStr == "" {
		return nil
	}
	var interval timeutil.Interval
	if err := interval.ValueOf(intervalStr); err != nil {
		return err
	}
	if interval <= 0 {
		return fmt.Errorf("interval cannot be negative")
	}
	return nil
}
