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

package point

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lindb/lindb/pkg/escape"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"

	"github.com/cespare/xxhash"
)

// Points represents a sortable list of points by timestamp.
type Points []Point

func (a Points) Len() int           { return len(a) }
func (a Points) Less(i, j int) bool { return a[i].UnixMilli() < a[j].UnixMilli() }
func (a Points) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Point implements point
type Point struct {
	timestamp          int64          // milliseconds
	key                []byte         // text encoding of metric-name and sorted tags
	fields             []byte         // text encoding of field data
	ts                 []byte         // text encoding of timestamp
	cachedMetricHashID uint64         // cached version of metric hash-id
	cachedTagsHashID   uint64         // cached version of tags hash-id
	cachedName         string         // cached version of parsed name from key
	cachedTags         tag.Tags       // cached version of parsed tags
	cachedFields       field.Fields   // cached version of parsed fields
	itr                *FieldIterator // field iterator
}

func (p *Point) Key() []byte {
	return p.key
}

func (p *Point) name() []byte {
	_, name := scanTo(p.key, 0, ',')
	return name
}

func (p *Point) Name() []byte {
	return escape.Unescape(p.name())
}

func (p *Point) SetName(name string) *Point {
	p.cachedName = ""
	p.cachedMetricHashID = 0
	p.key = MakeKey([]byte(name), p.Tags())
	return p
}

func (p *Point) Tags() tag.Tags {
	if p.cachedTags != nil {
		return p.cachedTags
	}
	p.cachedTags = parseTags(p.key, nil)
	return p.cachedTags
}

func (p *Point) ForEachTag(fn func(k, v []byte) bool) {
	walkTags(p.key, fn)
}

func (p *Point) AddTag(key, value string) *Point {
	p.cachedTagsHashID = 0
	tags := p.Tags()
	tags = append(tags, tag.Tag{Key: []byte(key), Value: []byte(value)})
	sort.Sort(tags)
	p.cachedTags = tags
	p.key = MakeKey(p.Name(), tags)
	return p
}

func (p *Point) AddTags(kv ...string) *Point {
	if len(kv)%2 == 1 {
		return p
	}
	for i := 0; i < len(kv); i += 2 {
		_ = p.AddTag(kv[i], kv[i+1])
	}
	return p
}

func (p *Point) SetTags(tags tag.Tags) *Point {
	p.cachedTagsHashID = 0
	p.key = MakeKey(p.Name(), tags)
	p.cachedTags = tags
	return p
}

func (p *Point) TagsHashKey() []byte {
	return p.Tags().AppendHashKey(nil)
}

func (p *Point) HasTag(tag []byte) bool {
	if len(p.key) == 0 {
		return false
	}

	var exists bool
	walkTags(p.key, func(key, value []byte) bool {
		if bytes.Equal(tag, key) {
			exists = true
			return false
		}
		return true
	})

	return exists
}

func (p *Point) HashID() uint64 {
	if p.cachedMetricHashID != 0 {
		return p.cachedMetricHashID
	}
	p.cachedMetricHashID = xxhash.Sum64(p.Name())
	return p.cachedMetricHashID
}

func (p *Point) TagsHashID() uint64 {
	if p.cachedTagsHashID != 0 {
		return p.cachedTagsHashID
	}
	p.cachedTagsHashID = xxhash.Sum64(p.TagsHashKey())
	return p.cachedTagsHashID
}

func (p *Point) Time() time.Time {
	return time.Unix(0, p.timestamp*1000*1000)
}

func (p *Point) SetTime(t time.Time) *Point {
	p.SetUnixMilli(t.UnixNano() / 1000 / 1000)
	return p
}

func (p *Point) UnixMilli() int64 {
	return p.timestamp
}

func (p *Point) SetUnixMilli(timestamp int64) *Point {
	p.timestamp = MilliSecondOf(timestamp)
	// make ts bytes slice
	p.ts = p.ts[:0]
	p.ts = strconv.AppendInt(p.ts, p.timestamp, 10)
	return p
}

func (p *Point) String() string {
	payload := string(p.Key()) + " " + string(p.fields)
	if len(payload) > 0 && payload[len(payload)-1] != ' ' {
		payload += " "
	}
	if p.timestamp != 0 {
		payload += strconv.FormatInt(p.timestamp, 10)
	}
	return payload
}

func (p *Point) StringSize() int {
	size := len(p.key) + len(p.fields) + 1
	if p.timestamp != 0 {
		digits := 1 // even "0" has one digit
		t := p.timestamp
		if t < 0 {
			// account for negative sign, then negate
			digits++
			t = -t
		}
		for t > 9 { // already accounted for one digit
			digits++
			t /= 10
		}
		size += digits + 1 // digits and a space
	}
	return size
}

func (p *Point) AppendString(buf []byte) []byte {
	buf = append(buf, p.key...)
	buf = append(buf, ' ')
	buf = append(buf, p.fields...)

	if len(buf) > 0 && buf[len(buf)-1] != ' ' {
		buf = append(buf, ' ')
	}
	if p.timestamp != 0 {
		buf = strconv.AppendInt(buf, p.timestamp, 10)
	}
	return buf
}

func (p *Point) Fields() (field.Fields, error) {
	if p.cachedFields != nil {
		return p.cachedFields, nil
	}
	fs, err := p.parseFields()
	if err != nil {
		return nil, err
	}
	p.SetFields(fs)
	return p.cachedFields, nil
}

func (p *Point) parseFields() (field.Fields, error) {
	itr := p.FieldIterator()
	fields := make(field.Fields, 8)[:0]
	for itr.Next() {
		if len(itr.Name()) == 0 {
			continue
		}
		switch itr.Type() {
		case field.Unknown:
			continue
		case field.SumField, field.MaxField, field.MinField, field.HistogramField, field.SummaryField:
			v, err := itr.Float64Value()
			if err != nil {
				return nil, fmt.Errorf("parse field: %s with error: %s", string(itr.Name()), err)
			}
			fields = append(fields, field.Field{
				Name: itr.Name(), Type: itr.Type(), Value: v})
		}
	}
	if len(fields) == 0 {
		return nil, ErrMissingFields
	}
	return fields, nil
}

func (p *Point) FieldIterator() *FieldIterator {
	if p.itr == nil {
		p.itr = new(FieldIterator)
	}
	p.itr.Reset(p.fields)
	return p.itr
}

func (p *Point) Reset() {
	p.timestamp = 0
	p.key = p.key[:0]
	p.fields = p.fields[:0]
	p.ts = p.ts[:0]
	p.cachedMetricHashID = 0
	p.cachedTagsHashID = 0
	p.cachedName = ""
	p.cachedTags = p.cachedTags[:0]
	p.cachedFields = p.cachedFields[:0]
	p.itr.Reset(nil)
}

func (p *Point) AddField(name string, fType field.Type, value interface{}) *Point {
	p.cachedFields = p.cachedFields.Insert(
		field.Field{Name: []byte(name), Type: fType, Value: value})
	p.fields = p.fields[:0]
	p.fields = MakeFields(p.fields, p.cachedFields)
	return p
}

func (p *Point) SetFields(fs field.Fields) *Point {
	p.cachedFields = fs
	p.fields = p.fields[:0]
	p.fields = MakeFields(p.fields, p.cachedFields)
	return p
}

// ParsePoints returns a slice of *Point from the text representation of lines separated by new lines,
// If any poins fails to parse, a error will be returned.
func ParsePoints(text []byte) ([]*Point, error) {
	points := make([]*Point, bytes.Count(text, []byte{'\n'})+1)[:0]
	var (
		pos    int
		block  []byte
		failed []string
		err    error
	)
	for pos < len(text) {
		pos, block = scanLine(text, pos)
		pos++

		if len(block) == 0 {
			continue
		}

		// lines which start with '#' are comments
		start := skipWhitespace(block, 0)

		// If line is all whitespace, just skip it
		if start >= len(block) {
			continue
		}

		if block[start] == '#' {
			continue
		}

		// strip the newline if one is present
		if block[len(block)-1] == '\n' {
			block = block[:len(block)-1]
		}

		points, err = parsePointsAppend(points, block[start:])
		if err != nil {
			failed = append(failed, fmt.Sprintf("unable to parse '%s': %v", string(block[start:]), err))
		}
	}
	if len(failed) > 0 {
		return points, fmt.Errorf("%s", strings.Join(failed, "\n"))
	}

	return points, nil
}

// ParsePointsFromString is identical to ParsePoints but accepts a string.
func ParsePointsFromString(text string) ([]*Point, error) {
	return ParsePoints([]byte(text))
}

func parsePointsAppend(points []*Point, text []byte) ([]*Point, error) {
	// scan the first block which is measurement[,tag1=value1,tag2=value=2...]
	pos, key, err := scanKey(text, 0)
	if err != nil {
		return nil, err
	}
	// metric-name name is required
	if len(key) == 0 {
		return points, ErrMissingMetricName
	}

	// Since the metric-name is converted to a tag and metric-name & tags have
	// different escaping rules, we need to check if the metric-name needs escaping.
	_, i, _ := scanMetricName(key, 0)
	keyMetricName := key[:i-1]
	if bytes.IndexByte(keyMetricName, '=') != -1 {
		escapedKeyMetricName := bytes.Replace(keyMetricName, []byte("="), []byte(`\=`), -1)

		newKey := make([]byte, len(escapedKeyMetricName)+(len(key)-len(keyMetricName)))
		copy(newKey, escapedKeyMetricName)
		copy(newKey[len(escapedKeyMetricName):], key[len(keyMetricName):])
		key = newKey
	}

	// scan the second block is which is field1=value1[,field2=value2,...]
	// at least one field is required
	pos, fields, err := scanFields(text, pos)
	if err != nil {
		return points, err
	} else if len(fields) == 0 {
		return points, ErrMissingFields
	}

	// scan the last block which is an optional integer timestamp
	pos, ts, err := scanTime(text, pos)
	if err != nil {
		return points, err
	}

	// Build point with timestamp only.
	pt := &Point{ts: ts}

	if len(ts) == 0 {
		pt.timestamp = timeutil.Now()
	} else {
		ts, err := parseInt64Bytes(ts)
		if err != nil {
			return points, err
		}
		pt.timestamp = MilliSecondOf(ts)

		// Determine if there are illegal non-whitespace characters after the
		// timestamp block.
		for pos < len(text) {
			if text[pos] != ' ' {
				return points, ErrInvalidPoint
			}
			pos++
		}
	}

	// validate fields
	if err := walkFields(fields, func(k, v, fieldBuf []byte) bool {
		return true
	}); err != nil {
		return points, err
	}

	pt.key = key
	pt.fields = fields
	points = append(points, pt)
	return points, nil
}
