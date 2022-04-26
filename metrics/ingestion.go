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

package metrics

import (
	"time"

	"github.com/lindb/lindb/internal/linmetric"
)

// FlatIngestionStatistics represents flot ingestion statistics.
type FlatIngestionStatistics struct {
	CorruptedData   *linmetric.BoundCounter // corrupted when parse
	DroppedMetric   *linmetric.BoundCounter // drop when append
	IngestedMetrics *linmetric.BoundCounter // ingested metrics
	ReadBytes       *linmetric.BoundCounter // read data bytes
	LT10KiBCounter  *linmetric.BoundCounter // <=10k count
	LT100KiBCounter *linmetric.BoundCounter // <=100k count
	LT1MiBCounter   *linmetric.BoundCounter // <=1mb count
	LT10MiBCounter  *linmetric.BoundCounter // <=10mb count
	GT10MiBCounter  *linmetric.BoundCounter // >=10mb count
}

// InfluxIngestionStatistics represents influx ingestion statistics.
type InfluxIngestionStatistics struct {
	CorruptedData   *linmetric.BoundCounter // corrupted when parse
	IngestedMetrics *linmetric.BoundCounter // ingested metrics
	IngestedFields  *linmetric.BoundCounter // ingested fields
	ReadBytes       *linmetric.BoundCounter // read data bytes
	DroppedMetrics  *linmetric.BoundCounter // drop metric when append
	DroppedFields   *linmetric.BoundCounter // drop field when append
}

// NativeIngestionStatistics represents native ingestion statistics.
type NativeIngestionStatistics struct {
	CorruptedData   *linmetric.BoundCounter // corrupted when parse
	IngestedMetrics *linmetric.BoundCounter // ingested metrics
	ReadBytes       *linmetric.BoundCounter // read data bytes
	DroppedMetrics  *linmetric.BoundCounter // drop metric when append
}

// CommonIngestionStatistics represents ingestion common statistics.
type CommonIngestionStatistics struct {
	Duration *linmetric.DeltaHistogramVec // ingest duration(include count)
}

// NewNativeIngestionStatistics creates a native ingestion statistics.
func NewNativeIngestionStatistics() *NativeIngestionStatistics {
	influxIngestionScope := linmetric.BrokerRegistry.NewScope("lindb.ingestion.influx")
	return &NativeIngestionStatistics{
		CorruptedData:   influxIngestionScope.NewCounter("data_corrupted"),
		IngestedMetrics: influxIngestionScope.NewCounter("ingested_metrics"),
		ReadBytes:       influxIngestionScope.NewCounter("read_bytes"),
		DroppedMetrics:  influxIngestionScope.NewCounter("dropped_metrics"),
	}
}

// NewInfluxIngestionStatistics creates an influx ingestion statistics.
func NewInfluxIngestionStatistics() *InfluxIngestionStatistics {
	influxIngestionScope := linmetric.BrokerRegistry.NewScope("lindb.ingestion.influx")
	return &InfluxIngestionStatistics{
		CorruptedData:   influxIngestionScope.NewCounter("data_corrupted"),
		IngestedMetrics: influxIngestionScope.NewCounter("ingested_metrics"),
		IngestedFields:  influxIngestionScope.NewCounter("ingested_fields"),
		ReadBytes:       influxIngestionScope.NewCounter("read_bytes"),
		DroppedMetrics:  influxIngestionScope.NewCounter("dropped_metrics"),
		DroppedFields:   influxIngestionScope.NewCounter("dropped_fields"),
	}
}

// NewCommonIngestionStatistics creates an ingestion common statistics.
func NewCommonIngestionStatistics() *CommonIngestionStatistics {
	return &CommonIngestionStatistics{
		Duration: linmetric.BrokerRegistry.
			NewScope(
				"lindb.http.ingest_duration",
			).
			NewHistogramVec("path").
			WithExponentBuckets(time.Millisecond, time.Second*5, 20),
	}
}

// NewFlatIngestionStatistics creates a flat ingestion statistics.
func NewFlatIngestionStatistics() *FlatIngestionStatistics {
	scope := linmetric.BrokerRegistry.NewScope("lindb.ingestion.flat")
	flatIngestionBlockScope := scope.NewCounterVec("block", "size")
	return &FlatIngestionStatistics{
		CorruptedData:   scope.NewCounter("data_corrupted"),
		DroppedMetric:   scope.NewCounter("dropped_metrics"),
		IngestedMetrics: scope.NewCounter("ingested_metrics"),
		ReadBytes:       scope.NewCounter("read_bytes"),
		LT10KiBCounter:  flatIngestionBlockScope.WithTagValues("<10KiB"),
		LT100KiBCounter: flatIngestionBlockScope.WithTagValues("<100KiB"),
		LT1MiBCounter:   flatIngestionBlockScope.WithTagValues("<1MiB"),
		LT10MiBCounter:  flatIngestionBlockScope.WithTagValues("<10MiB"),
		GT10MiBCounter:  flatIngestionBlockScope.WithTagValues(">=10MiB"),
	}
}
