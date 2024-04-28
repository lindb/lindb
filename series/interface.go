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

package series

import (
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

//go:generate mockgen -source ./interface.go -destination=./interface_mock.go -package=series

// MetricMetaSuggester represents to suggest ability for metricNames and tagKeys.
// default max limit of suggestions is set in constants
type MetricMetaSuggester interface {
	// SuggestNamespace suggests the namespace by namespace's prefix
	SuggestNamespace(prefix string, limit int) (namespaces []string, err error)
	// SuggestMetrics returns suggestions from a given prefix of metricName
	SuggestMetrics(namespace, metricPrefix string, limit int) ([]string, error)
}

// TagValueSuggester represents to suggest ability for tagValues.
// default max limit of suggestions is set in constants
type TagValueSuggester interface {
	// SuggestTagValues returns suggestions from given tag key id and prefix of tagValue
	SuggestTagValues(tagKeyID tag.KeyID, tagValuePrefix string, limit int) ([]string, error)
}

// Filter represents the query ability for filtering seriesIDs by expr from an index of tags.
type Filter interface {
	// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec tag key of metric
	GetSeriesIDsByTagValueIDs(tagKeyID tag.KeyID, tagValueIDs *roaring.Bitmap) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag gets series ids for spec tag key of metric
	GetSeriesIDsForTag(tagKeyID tag.KeyID) (*roaring.Bitmap, error)
	// GetSeriesIDsForMetric gets series ids for spec metric name
	GetSeriesIDsForMetric(metricID metric.ID) (*roaring.Bitmap, error)
}

type FilterTimeRange interface {
	// GetSeriesIDsByTagValueIDs gets series ids by tag value ids for spec tag key of metric
	GetSeriesIDsByTagValueIDs(tagKeyID tag.KeyID, tagValueIDs *roaring.Bitmap, timeRange timeutil.TimeRange) (*roaring.Bitmap, error)
	// GetSeriesIDsForTag gets series ids for spec tag key of metric
	GetSeriesIDsForTag(tagKeyID tag.KeyID, timeRange timeutil.TimeRange) (*roaring.Bitmap, error)
	// GetSeriesIDsForMetric gets series ids for spec metric name
	GetSeriesIDsForMetric(metricID metric.ID, timeRange timeutil.TimeRange) (*roaring.Bitmap, error)
}
