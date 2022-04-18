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

package tsdb

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

// define database storage structure.
// directory tree for database[xx]:
//    xx/OPTIONS => config file
//    xx/meta/namespace => namespace metadata
//    xx/meta/metric => metrics' name metadata
//    xx/meta/field => metrics' field metadata
//    xx/meta/tagkey => metrics' tag key metadata
//    xx/meta/tagvalue => metrics' tag value metadata
//    xx/shard/1/(path)
//    xx/shard/1/buffer/123213123131 // time of ns
//    xx/shard/1/meta/
//    xx/shard/1/index/inverted/
//    xx/shard/1/segment/day/20191012/
//    xx/shard/1/segment/month/201910/
const (
	options          = "OPTIONS"
	shardDir         = "shard"
	metaDir          = "meta"
	tagValueMetaDir  = "tagvalue"
	tagValueDir      = "tag_value"
	segmentDir       = "segment"
	indexParentDir   = "index"
	forwardIndexDir  = "forward"
	invertedIndexDir = "inverted"
	bufferDir        = "buffer"
)

// createDatabasePath creates database's root path if existed.
func createDatabasePath(database string) (string, error) {
	dbPath := filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database)
	if err := mkDirIfNotExist(dbPath); err != nil {
		return "", fmt.Errorf("create database[%s]'s path with error: %s", database, err)
	}
	return dbPath, nil
}

// optionsPath returns database's options file path.
func optionsPath(database string) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database, options)
}

// metricsMetaPath returns metrics' metadata storage path.
func metricsMetaPath(database string) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database, metaDir)
}

// tagMetaIndicator returns database's tag metadata indicator information.
func tagMetaIndicator(database string) string {
	return filepath.Join(database, metaDir, tagValueMetaDir)
}

// shardIndicator returns shard indicator information.
func shardIndicator(database string, shardID models.ShardID) string {
	return filepath.Join(database, shardDir, strconv.Itoa(int(shardID)))
}

// shardPath returns shard's storage path.
func shardPath(database string, shardID models.ShardID) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, shardIndicator(database, shardID))
}

// shardTempBufferPath returns temp buffer path for write data.
func shardTempBufferPath(database string, shardID models.ShardID) string {
	return filepath.Join(shardPath(database, shardID), bufferDir)
}

// shardIndexIndicator returns shard level index indicator information.
func shardIndexIndicator(database string, shardID models.ShardID) string {
	return filepath.Join(shardIndicator(database, shardID), indexParentDir)
}

// shardMetaPath returns shard level metadata path.
func shardMetaPath(database string, shardID models.ShardID) string {
	return filepath.Join(shardPath(database, shardID), metaDir)
}

// ShardSegmentIndicator returns the segment name indicator information.
func ShardSegmentIndicator(database string, shardID models.ShardID, interval timeutil.Interval, name string) string {
	return filepath.Join(shardIndicator(database, shardID), segmentDir, interval.Type().String(), name)
}

// ShardSegmentPath returns segment path in shard dir.
func ShardSegmentPath(database string, shardID models.ShardID, interval timeutil.Interval) string {
	return filepath.Join(shardPath(database, shardID), segmentDir, interval.Type().String())
}
