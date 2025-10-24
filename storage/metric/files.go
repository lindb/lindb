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

package metric

import (
	"path/filepath"
	"strconv"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/store"
)

// define database storage structure.
// directory tree for database[xx]:
//
//	xx/OPTIONS => config file
//	xx/meta/namespace => namespace metadata
//	xx/meta/metric => metrics' name metadata
//	xx/meta/field => metrics' field metadata
//	xx/meta/tagkey => metrics' tag key metadata
//	xx/meta/tagvalue => metrics' tag value metadata
//	xx/shard/1/(path)
//	xx/shard/1/buffer/123213123131 // time of ns
//	xx/shard/1/index
//	xx/shard/1/partition/day/20191012/
//	xx/shard/1/partition/month/201910/
const (
	options          = "OPTIONS"
	shardDir         = "shard"
	metaDir          = "meta"
	tagValueMetaDir  = "tagvalue"
	tagValueDir      = "tag_value"
	indexParentDir   = "index"
	forwardIndexDir  = "forward"
	invertedIndexDir = "inverted"
	bufferDir        = "buffer"
	limits           = "limits.toml"
)

// metricsMetaPath returns metrics' metadata storage path.
func metricsMetaPath(database string) string {
	return filepath.Join(config.GlobalStorageConfig().TSDB.Dir, database, metaDir)
}

// shardIndicator returns shard indicator information.
func shardIndicator(database string, shardID models.ShardID) string {
	return filepath.Join(database, shardDir, strconv.Itoa(int(shardID)))
}

// shardTempBufferPath returns temp buffer path for write data.
func shardTempBufferPath(database string, shardID models.ShardID) string {
	return filepath.Join(store.ShardPath(database, shardID), bufferDir)
}

// shardIndexPath returns shard level index index path.
func shardIndexPath(database string, shardID models.ShardID) string {
	return filepath.Join(store.ShardPath(database, shardID), indexParentDir)
}
