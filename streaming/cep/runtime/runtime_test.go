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

package runtime

import (
	"fmt"
	"testing"
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/streaming/cep/stream/output"
)

type RPCService struct {
	Interface string
	Timestamp time.Time
	Tags      map[string]string
	Status    string

	TraceID  string
	SpanID   string
	Duration time.Duration
}

type Result struct{}

/*
Output[columnNames = []]
│ Layout: []
└─ Insert[]
   │ Layout: [tags_map:map, interface:string, qps:int@sum, ts:timestamp]
   └─ Projection[]
      │ Layout: [tags_map:map, interface:string, qps:int@sum, ts:timestamp]
      │ tags_map := map_values
      │ qps := count
      │ ts := time_trunc
      └─ Aggregate[step = SINGLE, keys = [map_values:map, time_trunc:timestamp, interface:string]]
         │ Layout: [map_values:map, time_trunc:timestamp, interface:string, count:int@sum]
         │ count := count("expr")
         └─ ScanFilterProjection[table = test:RPCService, filterPredicate = (interface IN ('grpc','http'))]
              Layout: [interface:string, map_values:map, time_trunc:timestamp]
              map_values := map_values(tags,app)
              time_trunc := time_trunc(timestamp,10s)
              Interval: 10s

Output[columnNames = [tags_map, interface, qps, ts]]
│ Layout: [map_values:map, interface:string, count:int@sum, time_trunc:timestamp]
│ tags_map := map_values
│ qps := count
│ ts := time_trunc
└─ Aggregate[keys = [map_values:map, time_trunc:timestamp, interface:string], step = SINGLE]
   │ Layout: [map_values:map, time_trunc:timestamp, interface:string, count:int@sum]
   │ count := count("expr")
   └─ ScanFilterProjection[table = test:RPCService, filterPredicate = (interface IN ('grpc','http'))]
        Layout: [interface:string, map_values:map, time_trunc:timestamp]
        map_values := map_values(tags,app)
        time_trunc := time_trunc(timestamp,10s)
        Interval: 10s
*/

func Test_Runtime_Insert(t *testing.T) {
	runtime := NewRuntime("test")
	runtime.RegisterStreamByType(RPCService{})
	runtime.RegisterStreamByType(Result{})
	// add result listener
	runtime.AddListener("Result", output.NewConsoleOutput())
	// add streaming query
	err := runtime.Query(`
	@app(name="test_app")
	@metric(name="count_rpc",tags=["tags_map","interface"],fields=["qps"],timestamp="ts")
	insert into Result
	select map_values(tags,'app') as tags_map,interface,count(1) as qps,time_trunc(timestamp,interval 10 second) as ts
	from RPCService
	where interface in('grpc','http')
	group by tags_map,interface,ts;
		`)
	fmt.Println(err)
	//
	// insert into Result
	// select interface,'rpc_call',count(1)
	// from RPCService where interface in('grpc','http')
	// group by interface;
	//
	//
	// @name(name="select result")
	// select * from Result
	//
	// @name(name="count_rpc",@header(user="test_user",pwd="pwd"))
	// insert into Result
	// select map_value(tags,'app') as tags_map,interface,count(1)
	// from RPCService where interface in('grpc','http')
	// group by map_value(tags,'app'),interfaec;

	// insert into Result
	// select interface,count(1)
	// from RPCService where interface in('grpc','http')
	// group by interfaec;
	//
	// @name(name="select result")
	//   select * from Result

	// send event
	input := runtime.GetInputHandler("RPCService")

	page := types.NewPage()
	interfaceColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "interface"}, interfaceColumn)
	timestampColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimestamp, Name: "timestamp"}, timestampColumn)
	tagsColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTMap, Name: "tags"}, tagsColumn)
	statusColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "status"}, statusColumn)
	traceColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "trace_id"}, traceColumn)
	spanColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "span_id"}, spanColumn)
	durationColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTInt, Name: "duration"}, durationColumn)

	interfaceColumn.Append("grpc")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "order"})
	traceColumn.Append("trace_grpc_order")
	spanColumn.Append("span_grpc_order")
	durationColumn.Append(int64(200))

	interfaceColumn.Append("http")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "user"})
	traceColumn.Append("trace_http_user")
	spanColumn.Append("span_http_user")
	durationColumn.Append(int64(150))

	interfaceColumn.Append("dubbo")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "order"})
	traceColumn.Append("trace_dubbo_order")
	spanColumn.Append("span_dubbo_order")
	durationColumn.Append(int64(300))

	interfaceColumn.Append("http")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "github"})
	traceColumn.Append("trace_http_github")
	spanColumn.Append("span_http_github")
	durationColumn.Append(int64(100))

	now := time.Now()
	// var wait sync.WaitGroup
	// wait.Add(5)
	// for range 5 {
	// 	go func() {
	// 		defer wait.Done()
	for range 3 {
		input.Send(page)
	}
	// 	}()
	// }
	// wait.Wait()
	fmt.Println(time.Since(now))

	// input.Send(&RPCService{Interface: "grpc", Tags: map[string]string{"host": "1.1.1.1", "app": "order"}})
	// input.Send(&RPCService{Interface: "http", Tags: map[string]string{"host": "1.1.1.1", "app": "user"}})
	// input.Send(&RPCService{Interface: "dubbo", Tags: map[string]string{"host": "1.1.1.1", "app": "order"}})
	// input.Send(&RPCService{Interface: "http", Tags: map[string]string{"host": "1.1.1.1", "app": "github"}})

	time.Sleep(10 * time.Second)
}

func Test_Runtime_Query(t *testing.T) {
	runtime := NewRuntime("test")
	runtime.RegisterStreamByType(RPCService{})
	runtime.RegisterStreamByType(Result{})
	// add result listener
	runtime.AddListener("Result", output.NewConsoleOutput())
	// add streaming query
	err := runtime.Query(`
	create job test_app
	begin
	  @app(name="test_app")

	  create sink rpc_call with (type="lindb",address="http://localhost:9003",database="_internal");

	  @sink(name="rpc_call")
	  @metric(name="{{.interface}}.rpc_call",tags=["tags_map","interface"],fields=["qps","exemplar"],timestamp="ts")
	  select map_values(tags,'app') as tags_map,
	  	interface,
	  	count(1) as qps,
	  	sampling(trace_id,span_id,duration) as exemplar,
	  	time_trunc(timestamp,interval 10 second) as ts 
	  from RPCService
	  where interface in('grpc','http')
	  group by tags_map,interface,ts;
	end
		`)
	fmt.Println(err)
	//
	// insert into Result
	// select interface,'rpc_call',count(1)
	// from RPCService where interface in('grpc','http')
	// group by interface;
	//
	//
	// @name(name="select result")
	// select * from Result
	//
	// @name(name="count_rpc",@header(user="test_user",pwd="pwd"))
	// insert into Result
	// select map_value(tags,'app') as tags_map,interface,count(1)
	// from RPCService where interface in('grpc','http')
	// group by map_value(tags,'app'),interfaec;

	// insert into Result
	// select interface,count(1)
	// from RPCService where interface in('grpc','http')
	// group by interfaec;
	//
	// @name(name="select result")
	//   select * from Result

	// send event
	input := runtime.GetInputHandler("RPCService")

	page := types.NewPage()
	interfaceColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "interface"}, interfaceColumn)
	timestampColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTTimestamp, Name: "timestamp"}, timestampColumn)
	tagsColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTMap, Name: "tags"}, tagsColumn)
	statusColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "status"}, statusColumn)
	traceColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "trace_id"}, traceColumn)
	spanColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTString, Name: "span_id"}, spanColumn)
	durationColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTDuration, Name: "duration"}, durationColumn)

	interfaceColumn.Append("grpc")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "order"})
	traceColumn.Append("trace_grpc_order")
	spanColumn.Append("span_grpc_order")
	durationColumn.Append(time.Duration(200))
	statusColumn.Append(nil)

	interfaceColumn.Append("http")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "user"})
	traceColumn.Append("trace_http_user")
	spanColumn.Append("span_http_user")
	durationColumn.Append(time.Duration(150))
	statusColumn.Append(nil)

	interfaceColumn.Append("dubbo")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "order"})
	traceColumn.Append("trace_dubbo_order")
	spanColumn.Append("span_dubbo_order")
	durationColumn.Append(time.Duration(300))
	statusColumn.Append(nil)

	interfaceColumn.Append("http")
	timestampColumn.Append(time.Now())
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "github"})
	traceColumn.Append("trace_http_github")
	spanColumn.Append("span_http_github")
	durationColumn.Append(time.Duration(100))
	statusColumn.Append(nil)

	now := time.Now()
	// var wait sync.WaitGroup
	// wait.Add(5)
	// for range 5 {
	// 	go func() {
	// 		defer wait.Done()
	for range 3 {
		input.Send(page)
	}
	// 	}()
	// }
	// wait.Wait()
	fmt.Println(time.Since(now))

	// input.Send(&RPCService{Interface: "grpc", Tags: map[string]string{"host": "1.1.1.1", "app": "order"}})
	// input.Send(&RPCService{Interface: "http", Tags: map[string]string{"host": "1.1.1.1", "app": "user"}})
	// input.Send(&RPCService{Interface: "dubbo", Tags: map[string]string{"host": "1.1.1.1", "app": "order"}})
	// input.Send(&RPCService{Interface: "http", Tags: map[string]string{"host": "1.1.1.1", "app": "github"}})

	time.Sleep(30 * time.Second)
}
