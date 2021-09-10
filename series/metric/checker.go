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
	"bytes"
	"strings"
)

// SanitizeMetricName checks if metric-name is in necessary of sanitizing
func SanitizeMetricName(metricName string) string {
	if !strings.Contains(metricName, "|") {
		return metricName
	}
	return strings.Replace(metricName, "|", "_", -1)
}

// SanitizeNamespace checks if namespace is in necessary of sanitizing
func SanitizeNamespace(namespace string) string {
	if !strings.Contains(namespace, "|") {
		return namespace
	}
	return strings.Replace(namespace, "|", "_", -1)
}

func ShouldSanitizeNamespaceOrMetricName(name []byte) bool {
	return bytes.IndexByte(name, '|') >= 0
}

func SanitizeNamespaceOrMetricName(name []byte) []byte {
	for idx := range name {
		if name[idx] == '|' {
			name[idx] = '_'
		}
	}
	return name
}

func ShouldSanitizeFieldName(fieldName []byte) bool {
	return bytes.HasPrefix(fieldName, []byte("Histogram")) ||
		bytes.HasPrefix(fieldName, []byte("__bucket_")) // bucket field
}

func SanitizeFieldName(fieldName []byte) []byte {
	switch {
	case bytes.HasPrefix(fieldName, []byte("Histogram")):
		var dst = make([]byte, len(fieldName)+1)
		dst[0] = byte('_')
		copy(dst[1:], fieldName)
		return dst
	case bytes.HasPrefix(fieldName, []byte("__bucket_")):
		return fieldName[1:]
	default:
		return fieldName
	}
}

// JoinNamespaceMetric concat namespace and metric-name for storage with a delimiter
func JoinNamespaceMetric(namespace, metricName string) string {
	return namespace + "|" + metricName
}
