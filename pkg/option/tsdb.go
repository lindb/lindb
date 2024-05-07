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

package option

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/timeutil"
)

// Intervals represents the list of Interval.
type Intervals []Interval

func (m Intervals) Len() int { return len(m) }

func (m Intervals) Less(i, j int) bool { return m[i].Interval < m[j].Interval }

func (m Intervals) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

// String returns the string representation of the Intervals.
func (m Intervals) String() string {
	rs := make([]string, len(m))
	for idx, i := range m {
		rs[idx] = i.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(rs, ","))
}

// IsValid checks if intervals is valid, if invalid return error.
func (m Intervals) IsValid() error {
	intervalMap := make(map[timeutil.IntervalType]Interval)
	for _, i := range m {
		intervalType := i.Interval.Type()
		exist, ok := intervalMap[intervalType]
		if ok {
			return fmt.Errorf("duplicate interval type,[%s(%s),%s(%s)]",
				exist.String(), intervalType.String(), i.String(), intervalType.String())
		}
		intervalMap[intervalType] = i
	}
	return nil
}

// Interval represents the database's interval option, include interval and data retention.
type Interval struct {
	Interval  timeutil.Interval `toml:"interval" json:"interval,omitempty" validate:"required"`
	Retention timeutil.Interval `toml:"retention" json:"retention,omitempty" validate:"required"`
}

// String returns the string representation of the Interval.
func (m Interval) String() string {
	return fmt.Sprintf("%s->%s", m.Interval, m.Retention)
}

// FlusherOption represents a flusher configuration for index and memory db
type FlusherOption struct {
	TimeThreshold int64 `toml:"timeThreshold" json:"timeThreshold"` // time level flush threshold
	SizeThreshold int64 `toml:"sizeThreshold" json:"sizeThreshold"` // size level flush threshold, unit(MB)
}

// DatabaseOption represents a database option include shard ids and shard's option
type DatabaseOption struct {
	Behind string `toml:"behind" json:"behind,omitempty"`
	Ahead  string `toml:"ahead" json:"ahead,omitempty"`

	Intervals Intervals     `toml:"intervals" json:"intervals,omitempty"  validate:"required"`
	Index     FlusherOption `toml:"index" json:"index,omitempty"`
	Data      FlusherOption `toml:"data" json:"data,omitempty"`

	ahead  int64
	behind int64

	AutoCreateNS bool `toml:"autoCreateNS" json:"autoCreateNS,omitempty"`
}

// FindMatchSmallestInterval returns the smallest interval which match query interval.
func (e *DatabaseOption) FindMatchSmallestInterval(interval timeutil.Interval) timeutil.Interval {
	storageIntervals := make([]timeutil.Interval, len(e.Intervals))
	idx := 0
	for k := range e.Intervals {
		storageIntervals[idx] = e.Intervals[k].Interval
		idx++
	}
	// desc order
	sort.Slice(storageIntervals, func(i, j int) bool {
		return storageIntervals[i] > storageIntervals[j]
	})

	storageInterval := e.Intervals[0].Interval // init using the smallest interval
	for _, sInterval := range storageIntervals {
		if interval >= sInterval {
			storageInterval = sInterval
			break
		}
	}
	return storageInterval
}

// Validate validates engine option if valid
func (e *DatabaseOption) Validate() error {
	if len(e.Intervals) == 0 {
		return errors.New("intervals cannot be empty")
	}
	if err := e.Intervals.IsValid(); err != nil {
		return err
	}
	// TODO: need remove
	if err := validateInterval(e.Ahead, false); err != nil {
		return err
	}
	if err := validateInterval(e.Behind, false); err != nil {
		return err
	}
	return nil
}

// GetAcceptWritableRange returns accept writable time range.
func (e *DatabaseOption) GetAcceptWritableRange() (ahead, behind int64) {
	if e.ahead <= 0 {
		e.ahead = e.getIntervalVal(e.Ahead)
	}
	if e.behind <= 0 {
		e.behind = e.getIntervalVal(e.Behind)
	}
	return e.ahead, e.behind
}

// getIntervalVal returns interval value.
func (e *DatabaseOption) getIntervalVal(interval string) int64 {
	var intervalVal timeutil.Interval
	_ = intervalVal.ValueOf(interval)
	return intervalVal.Int64()
}

// Default sets default value if some configuration item not set.
func (e *DatabaseOption) Default() {
	if e.Ahead == "" {
		e.Ahead = constants.MetricMaxAheadDurationStr
	}
	if e.Behind == "" {
		e.Behind = constants.MetricMaxBehindDurationStr
	}
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
