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

package kv

import (
	"time"

	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/pkg/timeutil"
)

// FamilyOption defines config items for family level
type FamilyOption struct {
	Name             string `toml:"name"`
	Merger           string `toml:"merger"`
	ID               int    `toml:"id"`
	CompactThreshold int    `toml:"compactThreshold"`
	RollupThreshold  int    `toml:"rollupThreshold"`
	MaxFileSize      uint32 `toml:"maxFileSize"`
}

// StoreOption defines config item for store level
type StoreOption struct {
	Rollup          []timeutil.Interval `toml:"rollup"`
	Levels          int                 `toml:"levels"`
	TTL             ltoml.Duration      `toml:"ttl"`
	CompactInterval ltoml.Duration      `toml:"compactInterval"`
	Source          timeutil.Interval   `toml:"source"`
}

// DefaultStoreOption builds default store option
func DefaultStoreOption() StoreOption {
	return StoreOption{
		Levels: 2,
		TTL:    ltoml.Duration(time.Hour),
	}
}

// storeInfo represents store config option, include all family's option in this kv store
type storeInfo struct {
	Families    map[string]FamilyOption `toml:"families"`
	StoreOption StoreOption             `toml:"store"`
}

// newStoreInfo creates store info instance for saving configs
func newStoreInfo(storeOption StoreOption) *storeInfo {
	return &storeInfo{
		StoreOption: storeOption,
		Families:    make(map[string]FamilyOption),
	}
}
