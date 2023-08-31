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
	"github.com/lindb/common/pkg/fileutil"
	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/tsdb/indexdb"
	"github.com/lindb/lindb/tsdb/memdb"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

// for testing
var (
	mkDirIfNotExist        = fileutil.MkDirIfNotExist
	listDir                = fileutil.GetDirectoryList
	removeDir              = fileutil.RemoveDir
	fileExist              = fileutil.Exist
	decodeToml             = ltoml.DecodeToml
	newDatabaseFunc        = newDatabase
	newSegmentFunc         = newSegment
	newMetadataFunc        = metadb.NewMetadata
	newShardFunc           = newShard
	encodeToml             = ltoml.EncodeToml
	newReaderFunc          = metricsdata.NewReader
	newFilterFunc          = metricsdata.NewFilter
	newIntervalSegmentFunc = newIntervalSegment
	newIndexDBFunc         = indexdb.NewIndexDatabase
	newMemoryDBFunc        = memdb.NewMemoryDatabase
	newDataFamilyFunc      = newDataFamily
	newMetricDataFlusher   = metricsdata.NewFlusher
	closeFamilyFunc        = closeFamily
	writeConfigFn          = ltoml.WriteConfig
)
