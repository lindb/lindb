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
	rsStr := `{"metricName":"lindb.runtime.mem","groupBy":["node"],"fields":["heap_inuse"],"startTime":1674831650000,"endTime":1674835250000,"interval":10000,"series":[{"tags":{"node":"192.168.0.103:9001"}},{"tags":{"node":"192.168.0.108:9001"},"fields":{"heap_inuse":{"1674831650000":34938880,"1674831660000":29769728,"1674831670000":30056448,"1674831680000":30359552,"1674831690000":30785536,"1674831700000":31211520,"1674831710000":31670272,"1674831720000":32153600,"1674831730000":32669696,"1674831740000":33341440,"1674831750000":33980416,"1674831760000":34775040,"1674831770000":35430400,"1674831780000":29802496,"1674831790000":30130176,"1674831800000":30425088,"1674831810000":30801920,"1674831820000":31162368,"1674831830000":31571968,"1674831840000":32055296,"1674831850000":32530432,"1674831860000":33218560,"1674831870000":34209792,"1674831880000":34897920,"1674831890000":35741696,"1674831900000":29818880,"1674831910000":30105600,"1674831920000":30375936,"1674831930000":30777344,"1674831940000":31326208,"1674831950000":31850496,"1674831960000":32514048,"1674831970000":32980992,"1674831980000":33595392,"1674831990000":34250752,"1674832000000":34930688,"1674832010000":35790848,"1674832020000":29827072,"1674832030000":30097408,"1674832040000":30416896,"1674832050000":30826496,"1674832060000":31383552,"1674832070000":31834112,"1674832080000":32325632,"1674832090000":32874496,"1674832100000":33406976,"1674832110000":34127872,"1674832120000":34914304,"1674832130000":35627008,"1674832140000":29851648,"1674832150000":30269440,"1674832160000":30605312,"1674832170000":30957568,"1674832180000":31358976,"1674832190000":31809536,"1674832200000":32374784,"1674832210000":32890880,"1674832220000":33546240,"1674832230000":34193408,"1674832240000":35078144,"1674832250000":35987456,"1674832260000":29769728,"1674832270000":30195712,"1674832280000":30613504,"1674832290000":30982144,"1674832300000":31383552,"1674832310000":31825920,"1674832320000":32342016,"1674832330000":33095680,"1674832340000":33669120,"1674832350000":34381824,"1674832360000":35102720,"1674832370000":35823616,"1674832380000":29990912,"1674832390000":30294016,"1674832400000":30556160,"1674832410000":30867456,"1674832420000":31432704,"1674832430000":31850496,"1674832440000":32399360,"1674832450000":33136640,"1674832460000":33652736,"1674832470000":34373632,"1674832480000":35045376,"1674832490000":35840000,"1674832500000":30220288,"1674832510000":30490624,"1674832520000":30760960,"1674832530000":31113216,"1674832540000":31457280,"1674832550000":31916032,"1674832560000":32391168,"1674832570000":32923648,"1674832580000":33513472,"1674832590000":34168832,"1674832600000":34824192,"1674832610000":35651584,"1674832620000":29925376,"1674832630000":30228480,"1674832640000":30523392,"1674832650000":31096832,"1674832660000":31399936,"1674832670000":31793152,"1674832680000":32325632,"1674832690000":32882688,"1674832700000":33423360,"1674832710000":34144256,"1674832720000":34791424,"1674832730000":35561472,"1674832740000":29769728,"1674832750000":30146560,"1674832760000":30523392,"1674832770000":31113216,"1674832780000":31514624,"1674832790000":32006144,"1674832800000":32563200,"1674832810000":33062912,"1674832820000":33693696,"1674832830000":34381824,"1674832840000":35151872,"1674832850000":35921920,"1674832860000":29745152,"1674832870000":30072832,"1674832880000":30441472,"1674832890000":30924800,"1674832900000":31244288,"1674832910000":31817728,"1674832920000":32342016,"1674833080000":23486464,"1674833090000":27992064,"1674833100000":29523968,"1674833110000":30171136,"1674833120000":31342592,"1674833130000":32145408,"1674833140000":32972800,"1674833150000":28434432,"1674833160000":28753920,"1674833170000":29048832,"1674833200000":23371776,"1674833210000":27377664}}}],"stats":{"node":"192.168.0.108:9000","waitCost":9859000,"waitStart":1674835250462420000,"waitEnd":1674835250472279000,"totalCost":10743335,"start":1674835250461979000,"end":1674835250472723000,"stages":[{"identifier":"Physical Plan","start":1674835250462017000,"end":1674835250462575000,"cost":557706,"state":"Complete","errMsg":"","async":false,"operators":[{"identifier":"Physical Plan","start":1674835250462026000,"end":1674835250462406000,"cost":380451}],"children":[{"identifier":"TaskSend","start":1674835250462417000,"end":1674835250462574000,"cost":157495,"state":"Complete","errMsg":"","async":false,"operators":[{"identifier":"Task Sender","start":1674835250462418000,"end":1674835250462572000,"cost":154388}],"children":null}]},{"identifier":"Expression","start":1674835250472681000,"end":1674835250472723000,"cost":41739,"state":"Complete","errMsg":"","async":false,"children":null}],"children":[{"node":"192.168.0.108:9001","waitCost":7795000,"waitStart":1674835250463798000,"waitEnd":1674835250471593000,"netPayload":5694,"totalCost":8435793,"start":1674835250463450000,"end":1674835250471886000,"stages":[{"identifier":"Physical Plan","start":1674835250463484000,"end":1674835250463830000,"cost":346120,"state":"Complete","errMsg":"","async":false,"operators":[{"identifier":"Physical Plan","start":1674835250463491000,"end":1674835250463787000,"cost":295661}],"children":[{"identifier":"TaskSend","start":1674835250463793000,"end":1674835250463829000,"cost":36279,"state":"Complete","errMsg":"","async":false,"operators":[{"identifier":"Task Sender","start":1674835250463796000,"end":1674835250463827000,"cost":30700}],"children":null}]}],"children":[{"node":"192.168.0.108:2891","netPayload":4933,"totalCost":7010289,"start":1674835250463939000,"end":1674835250470949000,"stages":[{"identifier":"Metadata Lookup","start":1674835250463992000,"end":1674835250466590000,"cost":2598552,"state":"Complete","errMsg":"","async":false,"operators":[{"identifier":"Metadata Lookup","start":1674835250463998000,"end":1674835250464012000,"cost":14036},{"identifier":"Tag Value Lookup","start":1674835250464016000,"end":1674835250466178000,"cost":2161505}],"children":[{"identifier":"Shard Scan[Shard(0)]","start":1674835250466194000,"end":1674835250470061000,"cost":3867256,"state":"Complete","errMsg":"","async":true,"operators":[{"identifier":"Series Filtering","start":1674835250466609000,"end":1674835250467583000,"cost":974087,"stats":{"numOfSeries":2}},{"identifier":"Data Family Read","start":1674835250467588000,"end":1674835250469328000,"cost":1740724},{"identifier":"Data Family Read","start":1674835250469332000,"end":1674835250469343000,"cost":11376},{"identifier":"Grouping Context Build","start":1674835250469347000,"end":1674835250470019000,"cost":671754}],"children":[{"identifier":"Grouping[Shard(0)]","start":1674835250470036000,"end":1674835250470171000,"cost":134793,"state":"Complete","errMsg":"","async":true,"operators":[{"identifier":"Grouping Tags Lookup","start":1674835250470087000,"end":1674835250470140000,"cost":52415}],"children":[{"identifier":"Data Load[2023-01-27 23:00:00]","start":1674835250470153000,"end":1674835250470868000,"cost":714618,"state":"Complete","errMsg":"","async":true,"operators":[{"identifier":"Data Load[/day/20230127/23/000046.sst]","start":1674835250470180000,"end":1674835250470829000,"cost":649384,"stats":{"numOfSeries":1}},{"identifier":"Reduce","start":1674835250470832000,"end":1674835250470833000,"cost":182},{"identifier":"Data Load[/day/20230127/23/000050.sst]","start":1674835250470833000,"end":1674835250470843000,"cost":10410,"stats":{"numOfSeries":1}},{"identifier":"Reduce","start":1674835250470845000,"end":1674835250470845000,"cost":123},{"identifier":"Data Load[/day/20230127/23/000052.sst]","start":1674835250470845000,"end":1674835250470855000,"cost":10226,"stats":{"numOfSeries":1}},{"identifier":"Reduce","start":1674835250470856000,"end":1674835250470867000,"cost":10790}],"children":null}]},{"identifier":"Grouping[Shard(0)]","start":1674835250470059000,"end":1674835250470198000,"cost":139152,"state":"Complete","errMsg":"","async":true,"operators":[{"identifier":"Grouping Tags Lookup","start":1674835250470107000,"end":1674835250470172000,"cost":65277}],"children":[{"identifier":"Data Load[2023-01-27 23:00:00]","start":1674835250470188000,"end":1674835250470290000,"cost":102538,"state":"Complete","errMsg":"","async":true,"operators":[{"identifier":"Data Load[/day/20230127/23/000046.sst]","start":1674835250470214000,"end":1674835250470227000,"cost":13080,"stats":{"numOfSeries":0}},{"identifier":"Reduce","start":1674835250470253000,"end":1674835250470254000,"cost":416},{"identifier":"Data Load[/day/20230127/23/000050.sst]","start":1674835250470260000,"end":1674835250470262000,"cost":2396,"stats":{"numOfSeries":0}},{"identifier":"Reduce","start":1674835250470265000,"end":1674835250470265000,"cost":134},{"identifier":"Data Load[/day/20230127/23/000052.sst]","start":1674835250470265000,"end":1674835250470271000,"cost":5373,"stats":{"numOfSeries":0}},{"identifier":"Reduce","start":1674835250470272000,"end":1674835250470283000,"cost":11409}],"children":null}]}]}]},{"identifier":"Grouping Collect","start":1674835250466185000,"end":1674835250470869000,"cost":4684000,"state":"Complete","errMsg":"","async":false,"children":null}]}]}]}}`
	rs := &ResultSet{}
	err := encoding.JSONUnmarshal([]byte(rsStr), rs)
	assert.NoError(t, err)
	rows, table := rs.ToTable()
	fmt.Println(rows)
	fmt.Println(table)
}
