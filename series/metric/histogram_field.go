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

import "strings"

const (
	// histoFieldSep is the ".__" separator between the histogram name and its physical sub-field.
	// Using ".__" avoids collision with user-defined fields (e.g. "duration_sum").
	histoFieldSep = ".__"

	// HistoBucketSubField is the sub-field name prefix for bucket columns, without the leading dot.
	// Full bucket field: <histoName> + "." + HistoBucketSubField + <bound>
	// e.g. "duration.__bucket_0.5"
	HistoBucketSubField = "__bucket_"

	// Stat sub-field name tokens — passed as the second argument to HistoStatFieldName.
	HistoStatSum   = "sum"
	HistoStatCount = "count"
	HistoStatMin   = "min"
	HistoStatMax   = "max"
)

// histoStatSuffixes are the full suffixes for stat fields, derived from the constants above.
var histoStatSuffixes = []string{
	histoFieldSep + HistoStatSum,
	histoFieldSep + HistoStatCount,
	histoFieldSep + HistoStatMin,
	histoFieldSep + HistoStatMax,
}

// IsBucketField reports whether name follows the histogram bucket field naming convention:
//
//	<histoName>.__bucket_<bound>
//
// The separator is always the FIRST dot in the name; everything after it must start with
// "__bucket_" to avoid false-positives from the decimal point in the bound value itself.
func IsBucketField(name string) bool {
	dot := strings.Index(name, ".")
	return dot >= 0 && strings.HasPrefix(name[dot+1:], HistoBucketSubField)
}

// IsHistoStatField reports whether name is a histogram scalar statistic field,
// i.e. ends with one of the suffixes: .__sum, .__count, .__min, .__max.
func IsHistoStatField(name string) bool {
	for _, s := range histoStatSuffixes {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

// IsHistogramRelatedField reports whether name is any histogram physical field —
// either a bucket field (<histoName>.__bucket_<bound>) or a statistic field.
func IsHistogramRelatedField(name string) bool {
	return IsBucketField(name) || IsHistoStatField(name)
}

// HistoNameFromField extracts the logical histogram name from a physical field name.
//   - Bucket fields (<histoName>.__bucket_<bound>): returns the part before the first dot.
//   - Stat fields (<histoName>.__sum|.__count|.__min|.__max): returns the part before the first dot.
//
// Both naming conventions use a ".<double-underscore>" separator, so the first-dot
// check handles both cases correctly.
// Returns "" if name does not match any histogram naming pattern.
func HistoNameFromField(name string) string {
	if dot := strings.Index(name, "."); dot >= 0 {
		return name[:dot]
	}
	return ""
}

// HistoStatFieldName constructs the physical storage field name for a histogram stat:
//
//	<histoName>.__<stat>   e.g. "duration.__sum", "duration.__count"
//
// stat should be one of HistoStatSum, HistoStatCount, HistoStatMin, HistoStatMax.
func HistoStatFieldName(histoName, stat string) string {
	return histoName + histoFieldSep + stat
}

// ParseBucketBound extracts the upper-bound float64 from a histogram bucket field name
// following the <histoName>.__bucket_<bound> convention.
// Returns (0, false) if fieldName is not a valid bucket field.
func ParseBucketBound(fieldName string) (float64, bool) {
	dot := strings.Index(fieldName, ".")
	if dot < 0 {
		return 0, false
	}
	bound, err := UpperBound(fieldName[dot+1:])
	if err != nil {
		return 0, false
	}
	return bound, true
}

// BucketFieldName constructs the full physical storage field name for a named histogram bucket:
//
//	<histoName>.<HistoBucketSubField><bound>
//
// e.g. BucketFieldName("latency", 0.5) → "latency.__bucket_0.5"
func BucketFieldName(histoName string, upperBound float64) string {
	return histoName + "." + BucketNameOfHistogramExplicitBound(upperBound)
}
