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

package indexdb

import (
	"io"

	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/metric"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=indexdb

var indexLogger = logger.GetLogger("TSDB", "IndexDB")

// FileIndexDatabase represents a database of index files, it is shard-level
// it provides the abilities to filter seriesID from the index.
// See `tsdb/doc` for index file layout.
type FileIndexDatabase interface {
	series.Filter
	series.TagValueSuggester
}

// IndexDatabase represents a index database includes memory/file storage, it is shard level.
// index database will generate series id if tags hash not exist in mapping storage, and
// builds inverted index for tags => series id
type IndexDatabase interface {
	io.Closer
	series.TagValueSuggester
	series.Filter
	// GetOrCreateSeriesID gets series by tags hash, if not exist generate new series id in memory,
	// if generate a new series id returns isCreate is true
	// if generate fail return err
	GetOrCreateSeriesID(metricID metric.ID, tagsHash uint64) (seriesID uint32, isCreated bool, err error)
	// BuildInvertIndex builds the inverted index for tag value => series ids,
	// the tags is considered as an empty key-value pair while tags is nil.
	BuildInvertIndex(namespace, metricName string, tagIterator *metric.KeyValueIterator, seriesID uint32)
	// Flush flushes index data to disk
	Flush() error
}
