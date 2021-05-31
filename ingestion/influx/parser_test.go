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
	"github.com/lindb/lindb/constants"
	pb "github.com/lindb/lindb/rpc/proto/field"

	"github.com/stretchr/testify/assert"

	"fmt"
	"strconv"
	"strings"
	"testing"
)

func Test_tooManyTags(t *testing.T) {
	var tagPair []string
	for i := 0; i <= constants.DefaultMaxTagKeysCount+1; i++ {
		v := strconv.FormatInt(int64(i), 10)
		tagPair = append(tagPair, fmt.Sprintf("%s=%s", v, v))
	}
	line := fmt.Sprintf("mmm,%s x=1,y=2 1465839830100400200", strings.Join(tagPair, ","))
	_, err := parseInfluxLine([]byte(line), "ns", -1e6)
	assert.Equal(t, ErrTooManyTags, err)
}

func Test_noTags_noTimestamp(t *testing.T) {
	m, err := parseInfluxLine([]byte("cpu value=1"), "ns2", -1e6)
	assert.Nil(t, err)
	assert.NotZero(t, m.Timestamp)
	assert.Empty(t, m.Tags)
}

func Test_badTimestamp(t *testing.T) {
	lines := []string{
		"cpu value=1 9223372036854775808",
		"cpu value=1 -92233720368547758078",
		"cpu value=1 -",
		"cpu value=1 -/",
		"cpu value=1 -1?",
		"cpu value=1 1-",
		"cpu value=1 9223372036854775807 12",
	}
	for _, line := range lines {
		m, err := parseInfluxLine([]byte(line), "ns3", 1)
		assert.Equal(t, ErrBadTimestamp, err)
		assert.Nil(t, m)
	}
}

func Test_tags(t *testing.T) {
	examples := []struct {
		Line string
		Tags map[string]string
	}{
		{`cpu value=1`, map[string]string{}},
		{"cpu,tag0=v0 value=1", map[string]string{"tag0": "v0"}},
		{"cpu,tag0=v0,tag1=v0 value=1", map[string]string{"tag0": "v0", "tag1": "v0"}},
		{`cpu,tag0=v\ 0 value=1`, map[string]string{"tag0": "v 0"}},
		{`cpu,tag0=v\ 0\ 1,tag1=v2 value=1`, map[string]string{"tag0": "v 0 1", "tag1": "v2"}},
		{`cpu,tag0=\, value=1`, map[string]string{"tag0": ","}},
		{`cpu,ta\ g0=\, value=1`, map[string]string{"ta g0": ","}},
		{`cpu,tag0=\,1 value=1`, map[string]string{"tag0": ",1"}},
		{`cpu,tag0=1\"\",t=k value=1`, map[string]string{"tag0": `1\"\"`, "t": "k"}},
	}
	for _, example := range examples {
		m, err := parseInfluxLine([]byte(example.Line), "ns", 1e6)
		assert.NotNil(t, m)
		assert.Nil(t, err)
		assert.Equal(t, example.Tags, m.Tags)
	}
}

func Test_InvalidLine(t *testing.T) {
	examples := []struct {
		Line string
		Err  error
	}{
		{``, ErrBadFields},
		{`a`, ErrBadFields},
		{` a`, ErrMissingMetricName},
		{`,a=b c=1`, ErrMissingMetricName},
		{`# `, nil},
	}
	for _, example := range examples {
		_, err := parseInfluxLine([]byte(example.Line), "ns", 1e6)
		assert.Equal(t, example.Err, err)
	}
}

func Test_metricName(t *testing.T) {
	examples := []struct {
		Line       string
		MetricName string
	}{
		{`cpu,tag0=v0 value=1 111`, "cpu"},
		{`cpu value=1 222`, "cpu"},
		{`cpu\  value=1`, "cpu "},
		{`cpu\ a,    tag0=v0 value=1`, "cpu a"},
		{`cpu\,a, tag0=v0 value=1`, "cpu,a"},
		{`cpu\,\ a, tag0=v0 value=1`, "cpu, a"},
		{`cpu\\\,\ a, tag0=v0 value=1`, "cpu\\\\, a"},
	}
	for _, example := range examples {
		m, err := parseInfluxLine([]byte(example.Line), "ns", 1e6)
		assert.NotNil(t, m)
		assert.Nil(t, err)
		assert.Equal(t, example.MetricName, m.Name)
	}
}

func Test_missingTagValues(t *testing.T) {
	examples := []struct {
		Line string
		Err  error
	}{
		{`cpu,host`, ErrMissingWhiteSpace},
		{`cpu,host,`, ErrMissingWhiteSpace},
		{`cpu,host=`, ErrMissingWhiteSpace},
		{`cpu,host value=1i`, ErrBadTags},
		{`cpu,host=serverA,region value=1i`, ErrBadTags},
		{`cpu,host=serverA,region= value=1i`, ErrBadTags},
		{`cpu,host=serverA,region=,zone=us-west value=1i`, ErrBadTags},
		{`cpu,host=f=o,`, ErrMissingWhiteSpace},
		{`cpu,host=f\==o,`, ErrMissingWhiteSpace},
	}
	for _, example := range examples {
		m, err := parseInfluxLine([]byte(example.Line), "ns", -1e6)
		assert.Equal(t, example.Err, err)
		assert.Nil(t, m)
	}
}

func Test_missingFieldNames(t *testing.T) {
	examples := []struct {
		Line string
		Err  error
	}{
		{`cpu,host=serverA,region=us-west =`, ErrBadFields},
		{`cpu,host=serverA,region=us-west =123i`, ErrBadFields},
		{`cpu,host=serverA,region=us-west a\ =123i`, nil},
		{`cpu,host=serverA,region=us-west value=123i,=456i`, ErrBadFields},
	}
	for _, example := range examples {
		_, err := parseInfluxLine([]byte(example.Line), "ns", 1e6)
		assert.Equal(t, example.Err, err)
	}
}

func Test_parseUnescapedMetric(t *testing.T) {
	examples := []struct {
		Line       string
		MetricName string
		Tags       map[string]string
		Fields     []*pb.Field
	}{
		// commas in metric name
		{`foo\,bar value_total=1i`,
			"foo,bar",
			map[string]string{},
			[]*pb.Field{{Name: "value_total", Type: pb.FieldType_Sum, Value: 1}}},
		// comma in metric name with tags
		{`cpu\,main,regions=east value=1.0 1465839830100400200`,
			"cpu,main",
			map[string]string{"regions": "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// spaces in metric name
		{`cpu\ load,region=east value_sum=1.0 1465839830100400200`,
			"cpu load",
			map[string]string{"region": "east"},
			[]*pb.Field{{Name: "value_sum", Type: pb.FieldType_Sum, Value: 1}}},
		// equals in metric name, boolean false
		{`cpu\=load,region=east value=false`,
			`cpu\=load`,
			map[string]string{"region": "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 0}}},
		// equals in metric name, boolean true
		{`cpu\=load,region=east value=true`,
			`cpu\=load`,
			map[string]string{"region": "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// commas in tag names, boolean true
		{`cpu,region\,zone=east value=t`,
			`cpu`,
			map[string]string{"region,zone": "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// spaces in tag name, boolean false
		{`cpu,region\ zone=east value=f`,
			`cpu`,
			map[string]string{"region zone": "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 0}}},
		// backslash with escaped equals in tag name, decimal value
		{`cpu,reg\=ion=east value=1.0`,
			`cpu`,
			map[string]string{"reg=ion": "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// space is tag name
		{`cpu,\ =east value=1.0`,
			`cpu`,
			map[string]string{` `: "east"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// commas in tag values
		{`cpu,regions=east\,west value=1.0`,
			`cpu`,
			map[string]string{"regions": "east,west"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// backslash literal followed by trailing space
		{`cpu,regions=east\  value=1.0`,
			`cpu`,
			map[string]string{"regions": `east `},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// spaces in tag values
		{`cpu,regions=east\ west value=1.0`,
			`cpu`,
			map[string]string{"regions": `east west`},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// commas in field keys
		{`cpu,regions=east value\,ms=1.0`,
			`cpu`,
			map[string]string{"regions": "east"},
			[]*pb.Field{{Name: "value,ms", Type: pb.FieldType_Gauge, Value: 1}}},
		// spaces in field keys
		{`cpu,regions=east value\ ms=1.0`,
			`cpu`,
			map[string]string{"regions": "east"},
			[]*pb.Field{{Name: "value ms", Type: pb.FieldType_Gauge, Value: 1}}},
		// random character escaped
		{`cpu,regions=eas\t value=1.0`,
			`cpu`,
			map[string]string{"regions": "eas\\t"},
			[]*pb.Field{{Name: "value", Type: pb.FieldType_Gauge, Value: 1}}},
		// field keys using escape char.
		{`cpu \a=1i`,
			`cpu`,
			map[string]string{},
			[]*pb.Field{{Name: "\\a", Type: pb.FieldType_Gauge, Value: 1}}},
		// measurement, tag and tag value with equals
		{`cpu=load,equals\=foo=tag\=value value=1i,bool=f`,
			`cpu=load`,
			map[string]string{"equals=foo": "tag=value"},
			[]*pb.Field{
				{Name: "value", Type: pb.FieldType_Gauge, Value: 1},
				{Name: "bool", Type: pb.FieldType_Gauge, Value: 0},
			}},
	}

	for _, example := range examples {
		m, err := parseInfluxLine([]byte(example.Line), "ns", -1e6)
		assert.Nil(t, err)
		assert.Equal(t, example.MetricName, m.Name)
		assert.Equal(t, example.Tags, m.Tags)
		assert.NotZero(t, m.Timestamp)
		assert.EqualValues(t, example.Fields, m.Fields)
	}
}

func Test_parseBadFields(t *testing.T) {
	lines := []string{
		`cpu,regions=east value="a1i"`,
		`cpu,regions=east value=a1i`,
		`cpu,regions=east value=`,
		`cpu,regions=east  value=2`,
		`cpu,regions=east =1`,
		`cpu,regions=east \ =1`,
		`cpu,regions=east value=a2e3`,
		`cpu,regions=east value=-a2e3`,
		`cpu,regions=east value=1t`,
		`cpu,regions=east value=2f`,
	}
	for _, line := range lines {
		_, err := parseInfluxLine([]byte(line), "ns", 1e6)
		assert.Equal(t, ErrBadFields, err)
	}
}
