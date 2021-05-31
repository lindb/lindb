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
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
)

var (
	ErrMissingMetricName = errors.New("missing_metric_name")
	ErrMissingWhiteSpace = errors.New("missing_whitespace")
	ErrBadTags           = errors.New("bad_tags")
	ErrTooManyTags       = errors.New("too_many_tags")
	ErrBadFields         = errors.New("bad_fields")
	ErrBadTimestamp      = errors.New("bad_timestamp")
)

// Test cases in
// https://github.com/influxdata/influxdb/blob/master/models/points_test.go

func parseInfluxLine(content []byte, namespace string, multiplier int64) (*pb.Metric, error) {
	// skip comment line
	if bytes.HasPrefix(content, []byte{'#'}) {
		return nil, nil
	}

	escaped := bytes.IndexByte(content, '\\') >= 0
	var (
		m pb.Metric
	)
	m.Namespace = namespace
	// parse metric-name
	metricEndAt, err := scanMetricName(content, escaped)
	if err != nil {
		return nil, err
	}
	m.Name = string(unescapeMetricName(content[:metricEndAt]))

	// parse tags
	tagsEndAt, err := scanTagLine(content, metricEndAt+1, escaped)
	if err != nil {
		return nil, err
	}
	if m.Tags, err = parseTags(content, metricEndAt+1, tagsEndAt, escaped); err != nil {
		return nil, err
	}
	if len(m.Tags) >= constants.DefaultMaxTagKeysCount {
		return nil, ErrTooManyTags
	}

	// parse fields
	fieldsEndAt, err := scanFieldLine(content, tagsEndAt+1, escaped)
	if err != nil {
		return nil, err
	}
	if m.Fields, err = parseFields(content, tagsEndAt+1, fieldsEndAt, escaped); err != nil {
		return nil, err
	}

	// parse timestamp
	if m.Timestamp, err = parseTimestamp(content, fieldsEndAt+1, multiplier); err != nil {
		return nil, err
	}
	return &m, nil
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

func parseFields(buf []byte, startAt int, endAt int, isEscaped bool) ([]*pb.Field, error) {
	var fields []*pb.Field
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
		f, err := parseField(buf[startAt:equalAt], buf[equalAt+1:boundaryAt])
		if err != nil {
			return fields, err
		}
		fields = append(fields, f)
		startAt = boundaryAt + 1
		goto WalkBeforeComma
	}
}

func parseField(key, value []byte) (*pb.Field, error) {
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
		return &pb.Field{Name: string(unescapedKey), Type: guessFieldType(key), Value: float64(v)}, nil
	case 't', 'T': // boolean true
		if len(value) == 1 {
			return &pb.Field{Name: string(unescapedKey), Type: pb.FieldType_Gauge, Value: 1}, nil
		}
		return nil, ErrBadFields
	case 'f', 'F': // boolean false
		if len(value) == 1 {
			return &pb.Field{Name: string(unescapedKey), Type: pb.FieldType_Gauge, Value: 0}, nil
		}
		return nil, ErrBadFields
	default:
		// boolean, always gauge
		lf := strutil.ByteSlice2String(value)
		// still boolean
		switch lf {
		case "false", "False", "FALSE":
			return &pb.Field{Name: string(unescapedKey), Type: pb.FieldType_Gauge, Value: 0}, nil
		case "true", "True", "TRUE":
			return &pb.Field{Name: string(unescapedKey), Type: pb.FieldType_Gauge, Value: 1}, nil
		default:
			// todo, decimal, such like 1e-3, 1e20, -3e23
			v, err := strconv.ParseFloat(lf, 64)
			if err != nil {
				return nil, ErrBadFields
			}
			return &pb.Field{Name: string(unescapedKey), Type: guessFieldType(key), Value: v}, nil
		}
	}
}

func guessFieldType(key []byte) pb.FieldType {
	switch {
	case bytes.HasSuffix(key, []byte("total")):
		return pb.FieldType_Sum
	case bytes.HasSuffix(key, []byte("sum")):
		return pb.FieldType_Sum
	default:
		return pb.FieldType_Gauge
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
	if multiplier > 0 {
		return f * multiplier, nil
	}
	return -1 * f / multiplier, nil
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
