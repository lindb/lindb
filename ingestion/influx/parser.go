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

package influx

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/proto/gen/v1/flatMetricsV1"
	"github.com/lindb/lindb/series/metric"
)

var (
	ErrMissingMetricName = errors.New("missing_metric_name")
	ErrMissingWhiteSpace = errors.New("missing_whitespace")
	ErrBadTags           = errors.New("bad_tags")
	ErrBadFields         = errors.New("bad_fields")
	ErrBadTimestamp      = errors.New("bad_timestamp")
)

var (
	influxIngestionScope       = linmetric.NewScope("lindb.ingestion.influx")
	influxCorruptedDataCounter = influxIngestionScope.NewCounter("data_corrupted_count")
	ingestedMetricsCounter     = influxIngestionScope.NewCounter("ingested_metrics")
	ingestedFieldsCounter      = influxIngestionScope.NewCounter("ingested_fields")
	influxReadBytesCounter     = influxIngestionScope.NewCounter("read_bytes")
	droppedMetricsCounter      = influxIngestionScope.NewCounter("dropped_metrics")
	droppedFieldsCounter       = influxIngestionScope.NewCounter("dropped_fields")
)

// Test cases in
// https://github.com/influxdata/influxdb/blob/master/models/points_test.go

func parseInfluxLine(
	builder *metric.RowBuilder,
	content []byte,
	namespace string,
	multiplier int64,
) error {
	// skip comment line
	if bytes.HasPrefix(content, []byte{'#'}) {
		return nil
	}

	escaped := bytes.IndexByte(content, '\\') >= 0
	builder.AddNameSpace(strutil.String2ByteSlice(namespace))
	// parse metric-name
	metricEndAt, err := scanMetricName(content, escaped)
	if err != nil {
		return nil
	}
	builder.AddMetricName(unescapeMetricName(content[:metricEndAt]))

	// parse tags
	tagsEndAt, err := scanTagLine(content, metricEndAt+1, escaped)
	if err != nil {
		return err
	}
	tags, err := parseTags(content, metricEndAt+1, tagsEndAt, escaped)
	if err != nil {
		return err
	}
	for k, v := range tags {
		if err := builder.AddTag(strutil.String2ByteSlice(k), strutil.String2ByteSlice(v)); err != nil {
			return err
		}
	}

	// parse fields
	fieldsEndAt, err := scanFieldLine(content, tagsEndAt+1, escaped)
	if err != nil {
		return err
	}
	fields, err := parseFields(content, tagsEndAt+1, fieldsEndAt, escaped)
	// return error only if fields are empty, just drop fields not supported in lindb like string.
	if err != nil && len(fields) == 0 {
		return err
	}
	for idx := range fields {
		if err := builder.AddSimpleField(fields[idx].Name, fields[idx].Type, fields[idx].Value); err != nil {
			return err
		}
	}

	// parse timestamp
	timestamp, err := parseTimestamp(content, fieldsEndAt+1, multiplier)
	if err != nil {
		return err
	}
	builder.AddTimestamp(timestamp)
	return nil
}

// walkToUnescapedChar returns first position of given unescaped char
// abc\,\,, -> 7
// abc, -> 3
// abc\\\\, -> 7
// \\\, -> -1
func walkToUnescapedChar(buf []byte, char byte, startAt int, isEscaped bool) int {
	if len(buf) <= startAt {
		return -1
	}
	for {
		offset := bytes.IndexByte(buf[startAt:], char)
		if offset < 0 {
			return -1
		}
		if !isEscaped {
			return startAt + offset
		}

		cursor := offset + startAt
		for cursor-1 >= startAt && buf[cursor-1] == '\\' {
			cursor--
		}
		if (offset+startAt-cursor)&1 == 1 {
			// seek right
			startAt += offset + 1
			continue
		}
		return offset + startAt
	}
}

// scanMetricName examines the metric-name part of a Point, and returns the end position
func scanMetricName(buf []byte, isEscaped bool) (endAt int, err error) {
	// unescaped comma;
	commaAt := walkToUnescapedChar(buf, ',', 0, isEscaped)
	switch {
	case commaAt == 0:
		return -1, ErrMissingMetricName
	case commaAt < 0:
		// cpu value=1, no comma
		whiteSpaceAt := walkToUnescapedChar(buf, ' ', 0, isEscaped)
		switch {
		case whiteSpaceAt > 0:
			return whiteSpaceAt, nil
		case whiteSpaceAt < 0:
			return -1, ErrBadFields
		default:
			return -1, ErrMissingMetricName
		}
	default:
		return commaAt, nil
	}
}

// scanTagLine returns the end position of tags
// weather,location=us-midwest,season=summer temperature=82,humidity=71 1465839830100400200
func scanTagLine(buf []byte, startAt int, isEscaped bool) (endAt int, err error) {
	// no tags
	if startAt > 0 && buf[startAt-1] == ' ' {
		return startAt - 1, nil
	}
	// skip prefix whitespace
	for startAt <= len(buf)-1 && buf[startAt] == ' ' {
		startAt++
	}
	endAt = walkToUnescapedChar(buf, ' ', startAt, isEscaped)
	switch {
	case endAt < 0:
		return -1, ErrMissingWhiteSpace
	default:
		// if endAt = 0, tags are empty
		return endAt, nil
	}
}

func parseTags(buf []byte, startAt int, endAt int, isEscaped bool) (map[string]string, error) {
	// empty
	tags := make(map[string]string)

WalkBeforeComma:
	{
		if startAt >= endAt-1 {
			return tags, nil
		}
		commaAt := walkToUnescapedChar(buf, ',', startAt, isEscaped)
		// '=' does not exist
		equalAt := walkToUnescapedChar(buf, '=', startAt, isEscaped)
		if equalAt <= startAt || equalAt+1 >= endAt {
			return tags, ErrBadTags
		}
		boundaryAt := endAt
		if commaAt > 0 && commaAt <= endAt {
			boundaryAt = commaAt
		}
		// move to next tag pair
		if equalAt+1 >= boundaryAt {
			return tags, ErrBadTags
		}
		// move to next tag pair
		tagKey, tagValue := buf[startAt:equalAt], buf[equalAt+1:boundaryAt]
		tags[string(unescapeTag(tagKey))] = string(unescapeTag(tagValue))
		startAt = boundaryAt + 1
		goto WalkBeforeComma
	}
}

// scanFieldLine returns the end position of fields
func scanFieldLine(buf []byte, startAt int, isEscaped bool) (endAt int, err error) {
	endAt = walkToUnescapedChar(buf, ' ', startAt, isEscaped)
	switch {
	case endAt < 0:
		// case: no timestamp
		endAt = len(buf)
		// but field line is empty
		if startAt == endAt {
			return -1, ErrBadFields
		}
		return endAt, nil
	case endAt == 0:
		return -1, ErrBadFields
	default:
		return endAt, nil
	}
}

type flatSimpleField struct {
	Name  []byte
	Type  flatMetricsV1.SimpleFieldType
	Value float64
}

func parseFields(
	buf []byte,
	startAt int,
	endAt int,
	isEscaped bool,
) (fields []flatSimpleField, err error) {
WalkBeforeComma:
	{
		if startAt >= endAt-1 {
			if len(fields) == 0 {
				return fields, ErrBadFields
			}
			return fields, nil
		}
		commaAt := walkToUnescapedChar(buf, ',', startAt, isEscaped)
		// '=' does not exist
		equalAt := walkToUnescapedChar(buf, '=', startAt, isEscaped)
		if equalAt <= startAt || equalAt+1 >= endAt {
			return fields, ErrBadFields
		}
		boundaryAt := endAt
		if commaAt > 0 && commaAt <= endAt {
			boundaryAt = commaAt
		}
		// move to next field pair
		if equalAt+1 >= boundaryAt {
			return fields, ErrBadFields
		}
		// move to next field pair
		var (
			parsedFields []flatSimpleField
		)
		parsedFields, err = parseField(buf[startAt:equalAt], buf[equalAt+1:boundaryAt])
		if err == nil {
			fields = append(fields, parsedFields...)
		} else {
			droppedFieldsCounter.Incr()
		}
		startAt = boundaryAt + 1
		goto WalkBeforeComma
	}
}

func parseField(key, value []byte) ([]flatSimpleField, error) {
	if len(value) == 0 {
		return nil, ErrBadFields
	}
	unescapedKey := unescapeTag(key)
	if len(unescapedKey) == 0 {
		return nil, ErrBadFields
	}
	if len(bytes.TrimSpace(unescapedKey)) == 0 {
		return nil, ErrBadFields
	}
	tail := value[len(value)-1]
	switch tail {
	case 'i', 'I', 'u', 'U': // is int or unsigned
		v, err := strconv.ParseInt(strutil.ByteSlice2String(value[0:len(value)-1]), 10, 64)
		if err != nil {
			return nil, ErrBadFields
		}
		return toLinGaugeAndSumField(unescapedKey, float64(v)), nil
	case 't', 'T': // boolean true
		if len(value) == 1 {
			return []flatSimpleField{{
				Name:  unescapedKey,
				Type:  flatMetricsV1.SimpleFieldTypeGauge,
				Value: float64(1),
			}}, nil
		}
		return nil, ErrBadFields
	case 'f', 'F': // boolean false
		if len(value) == 1 {
			return []flatSimpleField{{
				Name:  unescapedKey,
				Type:  flatMetricsV1.SimpleFieldTypeGauge,
				Value: float64(0),
			}}, nil
		}
		return nil, ErrBadFields
	default:
		// boolean, always gauge
		lf := strutil.ByteSlice2String(value)
		// still boolean
		switch lf {
		case "false", "False", "FALSE":
			return []flatSimpleField{{
				Name:  unescapedKey,
				Type:  flatMetricsV1.SimpleFieldTypeGauge,
				Value: float64(0),
			}}, nil
		case "true", "True", "TRUE":
			return []flatSimpleField{{
				Name:  unescapedKey,
				Type:  flatMetricsV1.SimpleFieldTypeGauge,
				Value: float64(1),
			}}, nil
		default:
			v, err := strconv.ParseFloat(lf, 64)
			if err != nil {
				return nil, ErrBadFields
			}
			return toLinGaugeAndSumField(unescapedKey, v), nil
		}
	}
}

func toLinGaugeAndSumField(key []byte, value float64) []flatSimpleField {
	switch {
	case bytes.HasSuffix(key, []byte("gauge")):
		return []flatSimpleField{{
			Name:  key,
			Type:  flatMetricsV1.SimpleFieldTypeGauge,
			Value: value,
		}}
	case bytes.HasSuffix(key, []byte("sum")):
		return []flatSimpleField{{
			Name:  key,
			Type:  flatMetricsV1.SimpleFieldTypeDeltaSum,
			Value: value,
		}}
	default:
		return []flatSimpleField{
			{
				Name:  []byte(string(key) + "_sum"),
				Type:  flatMetricsV1.SimpleFieldTypeDeltaSum,
				Value: value,
			},
			{
				Name:  []byte(string(key) + "_gauge"),
				Type:  flatMetricsV1.SimpleFieldTypeGauge,
				Value: value,
			},
		}
	}
}

func parseTimestamp(buf []byte, startAt int, multiplier int64) (int64, error) {
	// no timestamp
	if startAt >= len(buf) {
		return timeutil.Now(), nil
	}
	f, err := strconv.ParseInt(string(buf[startAt:]), 10, 64)
	if err != nil {
		return 0, ErrBadTimestamp
	}
	switch {
	// precision query not exist, unknown multiplier
	case multiplier == 0:
		return timestamp2MilliSeconds(f), nil
	case multiplier > 0:
		return f * multiplier, nil
	default:
		return -1 * f / multiplier, nil
	}
}

// timestamp2MilliSeconds guesses the real timestamp precision,
// then converts it into milliseconds
func timestamp2MilliSeconds(timestamp int64) int64 {
	min := fasttime.UnixMilliseconds() - int64(constants.MetricMaxBehindDuration)
	max := fasttime.UnixMilliseconds() + int64(constants.MetricMaxAheadDuration)

	switch {
	// ms
	case min < timestamp && timestamp < max:
		return timestamp
	// ns
	case min < (timestamp/1e6) && (timestamp/1e6) < max:
		return timestamp / 1e6
	// us
	case min < (timestamp/1e3) && (timestamp/1e3) < max:
		return timestamp / 1e3
	// m
	case min < (timestamp*60*1000) && (timestamp*60*1000) < max:
		return timestamp * 1000 * 60
	// h
	case min < (timestamp*1000*3600) && (timestamp*1000*3600) < max:
		return timestamp * 1000 * 3600
	// unknown precision, use milliseconds
	default:
		return fasttime.UnixMilliseconds()
	}
}

type escapeSet struct {
	k   [1]byte
	esc [2]byte
}

var (
	metricNameEscapeCodes = [...]escapeSet{
		{k: [1]byte{','}, esc: [2]byte{'\\', ','}},
		{k: [1]byte{' '}, esc: [2]byte{'\\', ' '}},
	}

	tagEscapeCodes = [...]escapeSet{
		{k: [1]byte{','}, esc: [2]byte{'\\', ','}},
		{k: [1]byte{' '}, esc: [2]byte{'\\', ' '}},
		{k: [1]byte{'='}, esc: [2]byte{'\\', '='}},
	}
)

func unescapeMetricName(in []byte) []byte {
	if bytes.IndexByte(in, '\\') == -1 {
		return in
	}
	for i := range metricNameEscapeCodes {
		c := &metricNameEscapeCodes[i]
		if bytes.IndexByte(in, c.k[0]) != -1 {
			in = bytes.Replace(in, c.esc[:], c.k[:], -1)
		}
	}
	return in
}

func unescapeTag(in []byte) []byte {
	if bytes.IndexByte(in, '\\') == -1 {
		return in
	}

	for i := range tagEscapeCodes {
		c := &tagEscapeCodes[i]
		if bytes.IndexByte(in, c.k[0]) != -1 {
			in = bytes.Replace(in, c.esc[:], c.k[:], -1)
		}
	}
	return in
}
