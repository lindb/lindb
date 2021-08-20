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

package metricchecker

// todo : @codingcrush implement this later
//// ValidateProtoMetricV1 validate and escapes the metric on broker-side
//// If timestamp is zero, current timestamp will be set
//func ValidateProtoMetricV1(metric *protoMetricsV1.Metric) error {
//	if metric == nil {
//		return constants.ErrMetricPBNilMetric
//	}
//	if metric.Timestamp == 0 {
//		metric.Timestamp = fasttime.UnixTimestamp()
//	}
//	if len(metric.Name) == 0 {
//		return constants.ErrMetricPBEmptyMetricName
//	}
//	// empty field
//	if len(metric.SimpleFields) == 0 && metric.CompoundField == nil {
//		return constants.ErrMetricPBEmptyField
//	}
//	timestamp := metric.Timestamp
//	now := fasttime.UnixMilliseconds()
//	// check metric timestamp if in acceptable time range
//	if (timestamp < now-constants.MetricMaxBehindDuration) ||
//		(timestamp > constants.MetricMaxAheadDuration) {
//		return constants.ErrMetricOutOfTimeRange
//	}
//	// check metric timestamp if in acceptable time range
//	// validate empty tags
//	if len(metric.Tags) > 0 {
//		for idx := range metric.Tags {
//			// nil tag
//			if metric.Tags[idx] == nil {
//				return constants.ErrMetricEmptyTagKeyValue
//			}
//			// empty key value
//			if metric.Tags[idx].Key == "" || metric.Tags[idx].Value == "" {
//				return constants.ErrMetricEmptyTagKeyValue
//			}
//		}
//	}
//
//	// check simple fields
//	for idx := range metric.SimpleFields {
//		// nil value
//		if metric.SimpleFields[idx] == nil {
//			return constants.ErrBadMetricPBFormat
//		}
//		// field-name empty
//		if metric.SimpleFields[idx].Name == "" {
//			return constants.ErrMetricEmptyFieldName
//		}
//		// check sanitize
//		if HistogramConverter.NeedToSanitize(metric.SimpleFields[idx].Name) {
//			metric.SimpleFields[idx].Name = HistogramConverter.Sanitize(metric.SimpleFields[idx].Name)
//		}
//		// field type unspecified
//		switch metric.SimpleFields[idx].Type {
//		case protoMetricsV1.SimpleFieldType_SIMPLE_UNSPECIFIED:
//			return constants.ErrBadMetricPBFormat
//		default:
//		}
//		v := metric.SimpleFields[idx].Value
//		if math.IsNaN(v) {
//			return constants.ErrMetricNanField
//		}
//		if math.IsInf(v, 0) {
//			return constants.ErrMetricInfField
//		}
//	}
//	// no more compound field
//	if metric.CompoundField == nil {
//		return nil
//	}
//	// compound field-type unspecified
//	switch metric.CompoundField.Type {
//	case protoMetricsV1.CompoundFieldType_COMPOUND_UNSPECIFIED:
//		return constants.ErrBadMetricPBFormat
//	default:
//	}
//	// value length zero or length not match
//	if len(metric.CompoundField.Values) != len(metric.CompoundField.ExplicitBounds) ||
//		len(metric.CompoundField.Values) <= 2 {
//		return constants.ErrBadMetricPBFormat
//	}
//	// ensure compound field value > 0
//	if (metric.CompoundField.Max < 0) ||
//		metric.CompoundField.Min < 0 ||
//		metric.CompoundField.Sum < 0 ||
//		metric.CompoundField.Count < 0 {
//		return constants.ErrBadMetricPBFormat
//	}
//
//	for idx := 0; idx < len(metric.CompoundField.Values); idx++ {
//		// ensure value > 0
//		if metric.CompoundField.Values[idx] < 0 || metric.CompoundField.ExplicitBounds[idx] < 0 {
//			return constants.ErrBadMetricPBFormat
//		}
//		// ensure explicate bounds increase progressively
//		if idx >= 1 && metric.CompoundField.ExplicitBounds[idx] < metric.CompoundField.ExplicitBounds[idx-1] {
//			return constants.ErrBadMetricPBFormat
//		}
//		// ensure last bound is +Inf
//		if idx == len(metric.CompoundField.ExplicitBounds)-1 && !math.IsInf(metric.CompoundField.ExplicitBounds[idx], 1) {
//			return constants.ErrBadMetricPBFormat
//		}
//	}
//	return nil
//}
//
//// ContainsCumulative checks if the metric contains a cumulative field
//// Called before writing to memdb
//// TSDB provides a simple cache for translating cumulative into delta fields.
//func ContainsCumulative(metric *protoMetricsV1.Metric) bool {
//	for idx := range metric.SimpleFields {
//		switch metric.SimpleFields[idx].Type {
//		case protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM:
//			return true
//		}
//	}
//	if metric.CompoundField == nil {
//		return false
//	}
//	switch metric.CompoundField.Type {
//	case protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM:
//		return true
//	default:
//		return false
//	}
//}
