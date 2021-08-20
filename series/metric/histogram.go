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

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var HistogramConverter = histogramConverter{
	SumFieldName:   "HistogramSum",
	CountFieldName: "HistogramCount",
	MaxFieldName:   "HistogramMax",
	MinFieldName:   "HistogramMin",
}

type histogramConverter struct {
	SumFieldName   string
	CountFieldName string
	MaxFieldName   string
	MinFieldName   string
}

// UpperBound extracts the upper-bound from bucketName
func (hc histogramConverter) UpperBound(bucketName string) (float64, error) {
	// make sure it has prefix with __bucket_
	if !strings.HasPrefix(bucketName, "__bucket_") {
		return 0, fmt.Errorf("bucketName:%s not startswith '__bucket_", bucketName)
	}
	raw := bucketName[len("__bucket_"):]
	return strconv.ParseFloat(raw, 64)
}

// BucketName converts reserved field-name for histogram buckets.
func (hc histogramConverter) BucketName(upperBound float64) string {
	if math.IsInf(upperBound, 1) {
		return "__bucket_+Inf"
	}
	return "__bucket_" + strconv.FormatFloat(upperBound, 'f', -1, 32)
}

// Sanitize escapes the illegal field,
// if reserved field-name is used, the input will be escaped with underline.
// HistogramSum-> _HistogramSum
// __bucket_ -> _bucket_
func (hc histogramConverter) Sanitize(fieldName string) string {
	if strings.HasPrefix(fieldName, "Histogram") {
		return "_" + fieldName
	} else if strings.HasPrefix(fieldName, "__bucket_") {
		// truncate leading underline
		return fieldName[1:]
	}
	return fieldName
}
func (hc histogramConverter) NeedToSanitize(fieldName string) bool {
	if strings.HasPrefix(fieldName, "Histogram") {
		return true
	}
	if strings.HasPrefix(fieldName, "__bucket_") {
		return true
	}
	return false
}
