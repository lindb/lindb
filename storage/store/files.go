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

package store

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

// for testing
var (
	mkDirIfNotExist = fileutil.MkDirIfNotExist
)

// define database storage structure.
// directory tree for database[xx]:
//
//	xx/OPTIONS => config file
//	xx/shard/1/(path)
//	xx/shard/1/segment/day/20191012/
//	xx/shard/1/segment/month/201910/
const (
	options        = "OPTIONS"
	shardDir       = "shard"
	metaDir        = "meta"
	PartitionDir   = "partition"
	segmentDir     = "segment"
	indexParentDir = "index"
	bufferDir      = "buffer"
	limits         = "limits.toml"
)

// CreateDatabasePath creates database's root path if existed.
func CreateDatabasePath(database string) (string, error) {
	dbPath := filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database)
	if err := mkDirIfNotExist(dbPath); err != nil {
		return "", fmt.Errorf("create database[%s]'s path with error: %s", database, err)
	}
	return dbPath, nil
}

// LimitsPath returns database's limits file path.
func LimitsPath(database string) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database, limits)
}

// OptionsPath returns database's options file path.
func OptionsPath(database string) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database, options)
}

// metricsMetaPath returns metrics' metadata storage path.
func metricsMetaPath(database string) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database, metaDir)
}

// shardIndicator returns shard indicator information.
func shardIndicator(database string, shardID models.ShardID) string {
	return filepath.Join(database, shardDir, strconv.Itoa(int(shardID)))
}

// ShardPath returns shard's storage path.
func ShardPath(database string, shardID models.ShardID) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, shardIndicator(database, shardID))
}

// shardTempBufferPath returns temp buffer path for write data.
func shardTempBufferPath(database string, shardID models.ShardID) string {
	return filepath.Join(ShardPath(database, shardID), bufferDir)
}

// shardIndexPath returns shard level index index path.
func shardIndexPath(database string, shardID models.ShardID) string {
	return filepath.Join(ShardPath(database, shardID), indexParentDir)
}

// FIXME: new
func ShardIntervalSegmentPath(database string, shardID models.ShardID, interval timeutil.Interval) string {
	return filepath.Join(ShardPath(database, shardID), segmentDir, interval.Type().String())
}

// ShardSegmentPath returns segment path in shard dir.
func ShardSegmentPath(database string, shardID models.ShardID, interval timeutil.Interval, name string) string {
	return filepath.Join(ShardPath(database, shardID), segmentDir, interval.Type().String(), name)
}

// PartitionPath returns partition path in shard dir.
func PartitionPath(database string, shardID models.ShardID, interval timeutil.Interval, name string) string {
	return filepath.Join(ShardPath(database, shardID), PartitionDir, interval.Type().String(), name)
}
