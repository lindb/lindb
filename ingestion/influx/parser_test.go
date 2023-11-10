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
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"
	commonseries "github.com/lindb/common/series"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/metric"
)

func Test_tooManyTags(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)
	var tagPair []string
	for i := 0; i <= 10; i++ {
		v := strconv.FormatInt(int64(i), 10)
		tagPair = append(tagPair, fmt.Sprintf("%s=%s", v, v))
	}
	line := fmt.Sprintf("mmm,%s x=1,y=2 1465839830100400200", strings.Join(tagPair, ","))
	err := parseInfluxLine(builder, []byte(line), "ns", -1e6, models.NewDefaultLimits())
	assert.NoError(t, err)
	_, err = builder.Build()
	assert.NoError(t, err)
}

func Test_noTags_noTimestamp(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

	err := parseInfluxLine(builder, []byte("cpu value=1"), "ns2", -1e6, models.NewDefaultLimits())
	assert.Nil(t, err)
	var row metric.BrokerRow
	data, err := builder.Build()
	assert.NoError(t, err)
	(&row).FromBlock(data)
	m := row.Metric()
	assert.NotZero(t, m.Timestamp())
	assert.Equal(t, 0, m.KeyValuesLength())
}

func Test_badTimestamp(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

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
		builder.Reset()
		err := parseInfluxLine(builder, []byte(line), "ns3", 1, models.NewDefaultLimits())
		assert.Equal(t, ErrBadTimestamp, err)
	}
}

func Test_tags(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

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
		builder.Reset()
		err := parseInfluxLine(builder, []byte(example.Line), "ns", 1e6, models.NewDefaultLimits())
		assert.Nil(t, err)
		var br metric.BrokerRow
		data, err := builder.Build()
		assert.NoError(t, err)
		(&br).FromBlock(data)
		m := br.Metric()
		var mp = make(map[string]string)
		var kv flatMetricsV1.KeyValue
		for i := 0; i < m.KeyValuesLength(); i++ {
			m.KeyValues(&kv, i)
			mp[string(kv.Key())] = string(kv.Value())
		}
		assert.EqualValues(t, example.Tags, mp)
	}
}

func Test_InvalidLine(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

	examples := []struct {
		Line string
	}{
		{``},
		{`a`},
		{` a`},
		{`,a=b c=1`},
	}
	for _, example := range examples {
		builder.Reset()
		err := parseInfluxLine(builder, []byte(example.Line), "ns", 1e6, models.NewDefaultLimits())
		if err == nil {
			_, err = builder.Build()
		}
		assert.Error(t, err)
	}
}

func Test_metricName(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

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
		builder.Reset()
		err := parseInfluxLine(builder, []byte(example.Line), "ns", 1e6, models.NewDefaultLimits())
		assert.NoError(t, err)
		var row metric.BrokerRow
		data, err := builder.Build()
		assert.NoError(t, err)
		(&row).FromBlock(data)
		m := row.Metric()
		assert.Equal(t, example.MetricName, string(m.Name()))
	}
}

func Test_missingTagValues(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

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
		builder.Reset()
		err := parseInfluxLine(builder, []byte(example.Line), "ns", -1e6, models.NewDefaultLimits())
		assert.Equal(t, example.Err, err)
	}
}

func Test_missingFieldNames(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

	examples := []struct {
		Line       string
		Err        error
		FieldCount int
	}{
		{`cpu,host=serverA,region=us-west =`, ErrBadFields, 0},
		{`cpu,host=serverA,region=us-west =123i`, ErrBadFields, 0},
		{`cpu,host=serverA,region=us-west a\ =123i`, nil, 2},
		{`cpu,host=serverA,region=us-west value=123i,=456i`, nil, 2},
	}
	for _, example := range examples {
		builder.Reset()
		err := parseInfluxLine(builder, []byte(example.Line), "ns", 1e6, models.NewDefaultLimits())
		assert.Equal(t, example.Err, err)
		if example.FieldCount == 0 {
			assert.Error(t, err)
		} else {
			var row metric.BrokerRow
			data, err := builder.Build()
			assert.NoError(t, err)
			(&row).FromBlock(data)
			m := row.Metric()
			assert.Equalf(t, m.SimpleFieldsLength(), example.FieldCount, example.Line)
		}
	}
}

func Test_parseUnescapedMetric(t *testing.T) {
	examples := []struct {
		Line       string
		MetricName string
		Tags       map[string]string
		Fields     []flatSimpleField
	}{
		// commas in metric name
		{`foo\,bar value_total=1i`,
			"foo,bar",
			map[string]string{},
			[]flatSimpleField{
				{
					Name: []byte("value_total_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1,
				},
				{
					Name: []byte("value_total_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
				},
			},
		},
		// comma in metric name with tags
		{`cpu\,main,regions=east value=1.0 1465839830100400200`,
			"cpu,main",
			map[string]string{"regions": "east"},
			[]flatSimpleField{
				{
					Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1,
				},
				{
					Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
				},
			},
		},
		// spaces in metric name
		{`cpu\ load,region=east value_sum=1.0 1465839830100400200`,
			"cpu load",
			map[string]string{"region": "east"},
			[]flatSimpleField{{
				Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1,
			}},
		},
		// equals in metric name, boolean false
		{`cpu\=load,region=east value=false`,
			`cpu\=load`,
			map[string]string{"region": "east"},
			[]flatSimpleField{{
				Name: []byte("value"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 0,
			}},
		},
		// equals in metric name, boolean true
		{`cpu\=load,region=east value=true`,
			`cpu\=load`,
			map[string]string{"region": "east"},
			[]flatSimpleField{{
				Name: []byte("value"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
			}},
		},
		// commas in tag names, boolean true
		{`cpu,region\,zone=east value=t`,
			`cpu`,
			map[string]string{"region,zone": "east"},
			[]flatSimpleField{{
				Name: []byte("value"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
			}},
		},
		// spaces in tag name, boolean false
		{`cpu,region\ zone=east value=f`,
			`cpu`,
			map[string]string{"region zone": "east"},
			[]flatSimpleField{{
				Name: []byte("value"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 0,
			}},
		},
		// backslash with escaped equals in tag name, decimal value
		{`cpu,reg\=ion=east value=1.0`,
			`cpu`,
			map[string]string{"reg=ion": "east"},
			[]flatSimpleField{
				{
					Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1,
				},
				{
					Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
				},
			},
		},
		// space is tag name
		{`cpu,\ =east value=1.0`,
			`cpu`,
			map[string]string{` `: "east"},
			[]flatSimpleField{
				{Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1},
				{Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1},
			},
		},
		// commas in tag values
		{`cpu,regions=east\,west value=1.0`,
			`cpu`,
			map[string]string{"regions": "east,west"},
			[]flatSimpleField{
				{Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1},
				{Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1},
			},
		},
		// backslash literal followed by trailing space
		{`cpu,regions=east\  value=1.0`,
			`cpu`,
			map[string]string{"regions": `east `},
			[]flatSimpleField{
				{Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1},
				{Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1},
			},
		},
		// spaces in tag values
		{`cpu,regions=east\ west value=1.0`,
			`cpu`,
			map[string]string{"regions": `east west`},
			[]flatSimpleField{
				{Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1},
				{Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1},
			},
		},
		// commas in field keys
		{`cpu,regions=east value\,ms_last=1.0`,
			`cpu`,
			map[string]string{"regions": "east"},
			[]flatSimpleField{
				{
					Name: []byte("value,ms_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
				},
			},
		},
		// commas in field keys
		{`cpu,regions=east value\,ms_first=1.0`,
			`cpu`,
			map[string]string{"regions": "east"},
			[]flatSimpleField{
				{
					Name: []byte("value,ms_first"), Type: flatMetricsV1.SimpleFieldTypeFirst, Value: 1,
				},
			},
		},
		// spaces in field keys
		{`cpu,regions=east value\ ms=1.0`,
			`cpu`,
			map[string]string{"regions": "east"},
			[]flatSimpleField{
				{Name: []byte("value ms_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1},
				{Name: []byte("value ms_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1},
			},
		},
		// random character escaped
		{`cpu,regions=eas\t value=1.0`,
			`cpu`,
			map[string]string{"regions": "eas\\t"},
			[]flatSimpleField{
				{
					Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1,
				},
				{
					Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
				},
			},
		},
		// field keys using escape char.
		{`cpu \a=1i`,
			`cpu`,
			map[string]string{},
			[]flatSimpleField{
				{
					Name: []byte("\\a_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1,
				},
				{
					Name: []byte("\\a_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1,
				},
			},
		},
		// measurement, tag and tag value with equals
		{`cpu=load,equals\=foo=tag\=value value=1i,bool=f`,
			`cpu=load`,
			map[string]string{"equals=foo": "tag=value"},
			[]flatSimpleField{
				{Name: []byte("value_sum"), Type: flatMetricsV1.SimpleFieldTypeDeltaSum, Value: 1},
				{Name: []byte("value_last"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 1},
				{Name: []byte("bool"), Type: flatMetricsV1.SimpleFieldTypeLast, Value: 0},
			}},
	}

	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

	for _, example := range examples {
		builder.Reset()
		err := parseInfluxLine(builder, []byte(example.Line), "ns", -1e6, models.NewDefaultLimits())
		assert.Nil(t, err)
		var row metric.BrokerRow
		data, err := builder.Build()
		assert.NoError(t, err)
		(&row).FromBlock(data)
		var m = row.Metric()
		assert.Equal(t, example.MetricName, string(m.Name()))
		var mp = make(map[string]string)
		var kv flatMetricsV1.KeyValue
		for i := 0; i < m.KeyValuesLength(); i++ {
			m.KeyValues(&kv, i)
			mp[string(kv.Key())] = string(kv.Value())
		}
		assert.Equal(t, example.Tags, mp)
		assert.NotZero(t, m.Timestamp)

		var realFields []flatSimpleField
		var sf flatMetricsV1.SimpleField
		for i := 0; i < m.SimpleFieldsLength(); i++ {
			m.SimpleFields(&sf, i)
			realFields = append(realFields, flatSimpleField{
				Name:  sf.Name(),
				Type:  sf.Type(),
				Value: sf.Value(),
			})
		}

		assert.EqualValuesf(t, example.Fields, realFields, example.Line)
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
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)
	for _, line := range lines {
		builder.Reset()
		err := parseInfluxLine(builder, []byte(line), "ns", 1e6, models.NewDefaultLimits())
		assert.Equal(t, ErrBadFields, err)
	}
}

func Test_parseTimestamp(t *testing.T) {
	timestamp := fasttime.UnixMilliseconds()
	assert.Equal(t, timestamp, timestamp2MilliSeconds(timestamp))
	assert.Equal(t, timestamp, timestamp2MilliSeconds(timestamp/1000))
	assert.Equal(t, timestamp, timestamp2MilliSeconds(timestamp*1000))
	assert.Equal(t, timestamp, timestamp2MilliSeconds(timestamp*1000*1000))
	assert.InDelta(t, timestamp, timestamp2MilliSeconds(timestamp/1000/60), float64(1000*60))
	assert.InDelta(t, timestamp, timestamp2MilliSeconds(timestamp/1000/3600), float64(1000*3600))
}

func Test_parseField(t *testing.T) {
	fields, err := parseField(nil, nil)
	assert.Nil(t, fields)
	assert.Equal(t, ErrBadFields, err)

	fields, err = parseField([]byte("test"), nil)
	assert.Nil(t, fields)
	assert.Equal(t, ErrBadFields, err)
}

func Test_limits(t *testing.T) {
	builder, releaseFunc := commonseries.NewRowBuilder()
	defer releaseFunc(builder)

	limits := models.NewDefaultLimits()
	err := parseInfluxLine(builder, []byte("#"), "ns", 0, limits)
	assert.NoError(t, err)

	limits.MaxMetricNameLength = 5
	line := `system,regions=east value=1.0 1465839830100400200`
	// metric name limit
	err = parseInfluxLine(builder, []byte(line), "ns", 0, limits)
	assert.Equal(t, constants.ErrMetricNameTooLong, err)
	limits.MaxMetricNameLength = 0
	limits.MaxTagNameLength = 5
	// tag key limit
	err = parseInfluxLine(builder, []byte(line), "ns", 0, limits)
	assert.Equal(t, constants.ErrTagKeyTooLong, err)
	limits.MaxTagNameLength = 0
	limits.MaxTagValueLength = 3
	// tag value limit
	err = parseInfluxLine(builder, []byte(line), "ns", 0, limits)
	assert.Equal(t, constants.ErrTagValueTooLong, err)
	limits.MaxTagValueLength = 0
	limits.MaxFieldNameLength = 3
	// field nae limit
	err = parseInfluxLine(builder, []byte(line), "ns", 0, limits)
	assert.Equal(t, constants.ErrFieldNameTooLong, err)
	limits.MaxFieldNameLength = 0
	limits.MaxFieldsPerMetric = -1
	// tag value limit
	err = parseInfluxLine(builder, []byte(line), "ns", 0, limits)
	assert.Equal(t, constants.ErrTooManyFields, err)
	limits.MaxFieldsPerMetric = 0
	// tag value limit
	err = parseInfluxLine(builder, []byte(line), "ns", 0, limits)
	assert.NoError(t, err)
}
