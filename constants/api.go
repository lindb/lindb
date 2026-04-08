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

package constants

const (
	// APIRoot represents api root path.
	APIRoot = "/api"
	// APIVersion1 represents api version 1 path.
	APIVersion1 = "/v1"
	// APIVersion1CliPath represents api version 1 path for client.
	APIVersion1CliPath = "/api/v1"
	// APIPrometheus represents prometheus api version 1 path
	APIPrometheus = "/prometheus"
	// APIPrometheusPrefix represent prometheus api prefix
	APIPrometheusPrefix = "/prometheus/api/v1"
	// ContentTypeFlat represents flat buffer content type.
	ContentTypeFlat = "application/flatbuffer"
	// ContentTypeProto represents proto buffer content type.
	ContentTypeProto = "application/protobuf"
	// ContentTypeInflux represents influx content type.
	ContentTypeInflux = "application/influx"
	// ContentTypeJSON represents json content type.
	ContentTypeJSON = "application/json"
	// ContentTypeArrow represents arrow content type.
	ContentTypeArrow = "application/vnd.apache.arrow.stream"
	// ContentTypeOTelProto represents the OpenTelemetry protobuf content type.
	// Note: this uses "x-protobuf" (not "protobuf") as required by the OTLP HTTP spec.
	ContentTypeOTelProto = "application/x-protobuf"
)

const (
	// DatabaseHeader represents database header key in http request.
	DatabaseHeader = "X-LinDB-Database"
)

// EncodingType defines the type for content encoding
type EncodingType string

const (
	// EncodingProto represents protobuf encoding type
	EncodingProto EncodingType = "proto"
	// EncodingFlat represents json encoding type
	EncodingJSON EncodingType = "json"
	// EncodingInflux represents arrow encoding type
	EncodingArrow EncodingType = "arrow"
)
