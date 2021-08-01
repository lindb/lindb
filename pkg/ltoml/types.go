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

package ltoml

import (
	"errors"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	jsoniter "github.com/json-iterator/go"
)

// Duration is a TOML wrapper type for time.Duration.
type Duration time.Duration

// String returns the string representation of the duration.
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Duration returns the standard time.Duration
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// UnmarshalText parses a TOML value into a duration value.
// See https://github.com/BurntSushi/toml
func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}

// MarshalText converts a duration to a string for decoding toml
func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// UnmarshalJSON parses a JSON value into a duration value.
func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err := jsoniter.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		var err error
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(duration)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// MarshalJSON converts a duration to a string for decoding json
func (d *Duration) MarshalJSON() (data []byte, err error) {
	return jsoniter.Marshal(d.String())
}

// Size is a TOML wrapper type for size
// k/K -> KB, m/M -> MB, g/G -> GB
type Size uint64

// String returns the string representation of the size.
func (s Size) String() string {
	return humanize.IBytes(uint64(s))
}

// MarshalText converts a size to a string for decoding toml
func (s *Size) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// UnmarshalText parses a byte size from text.
func (s *Size) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return fmt.Errorf("size is empty")
	}
	v, err := humanize.ParseBytes(string(text))
	if err != nil {
		return err
	}
	*s = Size(v)
	return nil
}

// MarshalJSON converts a size to a human readable size
func (s *Size) MarshalJSON() (data []byte, err error) {
	return jsoniter.Marshal(s.String())
}

// UnmarshalJSON parses a JSON value into a size value.
func (s *Size) UnmarshalJSON(data []byte) (err error) {
	var v interface{}
	if err := jsoniter.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*s = Size(uint64(value))
		return nil
	case string:
		if len(data) > 2 && data[0] == '"' && data[len(data)-1] == '"' {
			data = data[1 : len(data)-1]
		}
		size, err := humanize.ParseBytes(string(data))
		if err != nil {
			return err
		}
		*s = Size(size)
		return nil
	default:
		return errors.New("invalid size")
	}
}
