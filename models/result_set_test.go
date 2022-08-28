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

package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestResultSet(t *testing.T) {
	rs := NewResultSet()
	series := NewSeries(map[string]string{"key": "value"}, "value")
	rs.AddSeries(series)
	points := NewPoints()
	points.AddPoint(int64(10), 10.0)
	series.AddField("f1", points)
	points = NewPoints()
	points.AddPoint(int64(20), 10.0)
	series.AddField("f1", points)

	assert.Equal(t, 1, len(rs.Series))
	s := rs.Series[0]
	assert.Equal(t, map[string]string{"key": "value"}, s.Tags)
	assert.Equal(t, map[int64]float64{
		int64(10): 10.0,
		int64(20): 10.0},
		s.Fields["f1"])
}

func TestResultSet_ToTable(t *testing.T) {
	rows, rs := NewResultSet().ToTable()
	assert.Zero(t, rows)
	assert.Empty(t, rs)

	rows, rs = (&ResultSet{
		MetricName: "cpu",
		GroupBy:    []string{"host", "ip"},
		Fields:     []string{"usage", "load"},
		Series: []*Series{{
			Tags:   map[string]string{"host": "host1", "ip": "1.1.1.1"},
			Fields: map[string]map[int64]float64{"usage": {timeutil.Now(): 1.1}, "load": {timeutil.Now(): 1.1}},
		}, {
			Tags:   map[string]string{"host": "host2", "ip": "1.1.1.1"},
			Fields: map[string]map[int64]float64{"usage": {timeutil.Now(): 1.1}, "load": {timeutil.Now(): 1.1}},
		}},
	}).ToTable()
	assert.Equal(t, rows, 2)
	assert.NotEmpty(t, rs)
}

func TestResultSet_Stats_ToTable(t *testing.T) {
	//nolint:lll
	rsStr := `{"metricName":"lindb.monitor.system.cpu_stat","fields":["idle"],"startTime":1659407350000,"endTime":1659410950000,"interval":10000,"series":[{"fields":{"idle":{"1659407840000":0.883618673490009,"1659407850000":0.8836069701859954,"1659407860000":0.8835990069256352,"1659407870000":0.8835915848206877,"1659407880000":0.8835812405425646,"1659407890000":0.8835690023036138,"1659409710000":0.8827340629060619,"1659409720000":0.8827288247011711,"1659410740000":0.8821966656343508,"1659410750000":0.882196074380958,"1659410760000":0.8821958493943506,"1659410770000":0.8821959788548275,"1659410780000":0.8821964560173532,"1659410790000":0.8821965958177824,"1659410800000":0.8821969801618673,"1659410810000":0.882197131228225,"1659410820000":0.8821973572553974,"1659410830000":0.8821977823366373,"1659410840000":0.8821983713876528,"1659410850000":0.8821989610709579,"1659410860000":0.882199573300995,"1659410870000":0.8821998618931616,"1659410880000":0.8821996656853625,"1659410890000":0.8822003545934965,"1659410900000":0.8822007902309552,"1659410910000":0.8822010526563662,"1659410920000":0.8822014856853968,"1659410930000":0.8822010475746209,"1659410940000":0.8822011836026491,"1659410950000":0.8822019006037599}}}],"stats":{"root":"localhost:9001","leafNodes":{"localhost:2891":{"netPayload":2323,"totalCost":3773922,"start":1659410959349,"end":1659410959353,"stages":[{"identifier":"Metadata lookup","start":1659410959350,"end":1659410959350,"cost":237992,"state":"Complete","errMsg":"","operators":[{"identifier":"Metadata lookup","start":1659410959350,"end":1659410959350,"cost":9969}],"children":[{"identifier":"Shard scan, shard(0)","start":1659410959350,"end":1659410959353,"cost":3004838,"state":"Complete","errMsg":"","operators":[{"identifier":"All series","start":1659410959350,"end":1659410959351,"cost":1082569,"stats":{"numOfSeries":5}},{"identifier":"Data family read","start":1659410959351,"end":1659410959352,"cost":858153},{"identifier":"Data family read","start":1659410959352,"end":1659410959353,"cost":815789}],"children":[{"identifier":"Grouping, shard(0)","start":1659410959353,"end":1659410959353,"cost":68503,"state":"Complete","errMsg":"","operators":[{"identifier":"Grouping tags lookup","start":1659410959353,"end":1659410959353,"cost":30815}],"children":[{"identifier":"Data load[2022-08-02 10:00:00]","start":1659410959353,"end":1659410959353,"cost":76013,"state":"Complete","errMsg":"","operators":[{"identifier":"Data load[/day/20220802/10/000002.sst]","start":1659410959353,"end":1659410959353,"cost":57904,"stats":{"numOfSeries":2}},{"identifier":"Reduce","start":1659410959353,"end":1659410959353,"cost":179}],"children":null},{"identifier":"Data load[2022-08-02 11:00:00]","start":1659410959353,"end":1659410959353,"cost":95115,"state":"Complete","errMsg":"","operators":[{"identifier":"Data load[2022-08-02 11:00:00/memory/readwrite]","start":1659410959353,"end":1659410959353,"cost":42833,"stats":{"numOfSeries":4}},{"identifier":"Reduce","start":1659410959353,"end":1659410959353,"cost":220},{"identifier":"Data load[/day/20220802/11/000004.sst]","start":1659410959353,"end":1659410959353,"cost":4096,"stats":{"numOfSeries":2}},{"identifier":"Reduce","start":1659410959353,"end":1659410959353,"cost":16270}],"children":null}]}]}]}]}},"netPayload":2323,"planCost":18375,"waitCost":5335903,"expressCost":22307,"totalCost":5376585}}`
	rs := &ResultSet{}
	err := encoding.JSONUnmarshal([]byte(rsStr), rs)
	assert.NoError(t, err)
	rows, table := rs.ToTable()
	fmt.Println(rows)
	fmt.Println(table)
}
