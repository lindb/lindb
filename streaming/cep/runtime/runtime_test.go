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
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
)

type RPCService struct {
	Interface string
	Timestamp time.Time
	Tags      map[string]string
	Status    string

	TraceID  [16]byte
	SpanID   [8]byte
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
	// add streaming query
	err := runtime.DeployJob("test_app", `
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
	// add streaming query
	err := runtime.DeployJob("test_app", `
	create job test_app
	begin
	  @app(name="test_app")

	  create sink rpc_call with (type="lindb",address="http://localhost:9003",database="_internal");

	  @sink(
		  name="rpc_call",
	    @metric(name="{{.interface}}.rpc_call",tags=["tags_map","interface"],fields=["qps","exemplar"],timestamp="ts")
		)
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

	var rpcs []RPCService

	for range 1 {
		rpcs = append(rpcs, RPCService{
			Interface: "grpc",
			Timestamp: time.Now(),
			Tags:      map[string]string{"host": "1.1.1.1", "app": "order"},
			TraceID:   [16]byte{1, 2, 3},
			SpanID:    [8]byte{4, 5, 6},
			Duration:  time.Duration(200),
		})

		rpcs = append(rpcs, RPCService{
			Interface: "http",
			Timestamp: time.Now(),
			Tags:      map[string]string{"host": "1.1.1.1", "app": "user"},
			TraceID:   [16]byte{1, 2, 3},
			SpanID:    [8]byte{4, 5, 6},
			Duration:  time.Duration(150),
		})

		rpcs = append(rpcs, RPCService{
			Interface: "dubbo",
			Timestamp: time.Now(),
			Tags:      map[string]string{"host": "1.1.1.1", "app": "order"},
			TraceID:   [16]byte{1, 2, 3},
			SpanID:    [8]byte{4, 5, 6},
			Duration:  time.Duration(300),
		})
		rpcs = append(rpcs, RPCService{
			Interface: "dubbo",
			Timestamp: time.Now(),
			Tags:      map[string]string{"host": "1.1.1.1", "app": "order"},
			TraceID:   [16]byte{1, 2, 3},
			SpanID:    [8]byte{4, 5, 6},
			Duration:  time.Duration(300),
		})

		rpcs = append(rpcs, RPCService{
			Interface: "http",
			Timestamp: time.Now(),
			Tags:      map[string]string{"host": "1.1.1.1", "app": "github"},
			TraceID:   [16]byte{1, 2, 3},
			SpanID:    [8]byte{4, 5, 6},
			Duration:  time.Duration(100),
		})
	}

	page, _ := ToRecord(memory.DefaultAllocator, rpcs)

	now := time.Now()
	// var wait sync.WaitGroup
	// wait.Add(5)
	// for range 5 {
	// 	go func() {
	// 		defer wait.Done()
	var total int
	fmt.Println(now)
	for range 100_0000 {

		total += int(page.NumRows())
		input.Send(page)
		// page.Release()
	}
	// 	}()
	// }
	// wait.Wait()
	fmt.Println(time.Since(now))
	fmt.Println(total)
	time.Sleep(10 * time.Second)

	runtime.UndeployJob("test_app")

	time.Sleep(3 * time.Second)

	for range 3 {
		input.Send(page)
	}
	// input.Send(&RPCService{Interface: "grpc", Tags: map[string]string{"host": "1.1.1.1", "app": "order"}})
	// input.Send(&RPCService{Interface: "http", Tags: map[string]string{"host": "1.1.1.1", "app": "user"}})
	// input.Send(&RPCService{Interface: "dubbo", Tags: map[string]string{"host": "1.1.1.1", "app": "order"}})
	// input.Send(&RPCService{Interface: "http", Tags: map[string]string{"host": "1.1.1.1", "app": "github"}})

	time.Sleep(10 * time.Second)
}

// FieldAppender 是一个闭包，负责将 struct 字段值写入 Builder
type FieldAppender func(builder array.Builder, fieldVal reflect.Value)

var (
	// 缓存 Schema 和 Appenders 提升性能
	schemaCache    sync.Map
	appendersCache sync.Map
)

// ToRecord 将任意 struct slice 转为 arrow.Record
func ToRecord(mem memory.Allocator, slice interface{}) (arrow.Record, error) {
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("slice input required")
	}

	if sliceVal.Len() == 0 {
		return nil, fmt.Errorf("empty slice")
	}

	elemType := sliceVal.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// 1. 获取或创建 Schema 和 Appenders
	schema, appenders := getMetadata(elemType)

	// 2. 创建 RecordBuilder
	b := array.NewRecordBuilder(mem, schema)
	defer b.Release()

	b.Reserve(sliceVal.Len())

	// 3. 填充数据
	for i := 0; i < sliceVal.Len(); i++ {
		item := sliceVal.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		for colIdx, appender := range appenders {
			appender(b.Field(colIdx), item.Field(colIdx))
		}
	}

	return b.NewRecord(), nil
}

// getMetadata 解析并缓存 struct 的反射信息
func getMetadata(t reflect.Type) (*arrow.Schema, []FieldAppender) {
	if s, ok := schemaCache.Load(t); ok {
		a, _ := appendersCache.Load(t)
		return s.(*arrow.Schema), a.([]FieldAppender)
	}

	fields := make([]arrow.Field, t.NumField())
	appenders := make([]FieldAppender, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		arrowType := goTypeToArrow(f.Type)
		fields[i] = arrow.Field{Name: lo.SnakeCase(f.Name), Type: arrowType}
		appenders[i] = makeAppender(f.Type)
	}

	schema := arrow.NewSchema(fields, nil)
	schemaCache.Store(t, schema)
	appendersCache.Store(t, appenders)
	return schema, appenders
}

// goTypeToArrow 处理 Go 类型到 Arrow 类型的映射
func goTypeToArrow(t reflect.Type) arrow.DataType {
	if t == reflect.TypeFor[time.Duration]() {
		return arrow.FixedWidthTypes.Duration_ns
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int64:
		return arrow.PrimitiveTypes.Int64
	case reflect.String:
		return arrow.BinaryTypes.String
	case reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return &arrow.FixedSizeBinaryType{ByteWidth: int(t.Len())}
		}
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return arrow.BinaryTypes.Binary
		}
	case reflect.Float64:
		return arrow.PrimitiveTypes.Float64
	case reflect.Map:
		if t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.String {
			return arrow.MapOf(arrow.BinaryTypes.String, arrow.BinaryTypes.String)
		}
	case reflect.TypeFor[time.Time]().Kind():
		return arrow.FixedWidthTypes.Timestamp_ms
	}
	return arrow.BinaryTypes.String // 默认回退
}

// makeAppender 预编译每个字段的写入逻辑
func makeAppender(t reflect.Type) FieldAppender {
	if t == reflect.TypeFor[time.Duration]() {
		return func(b array.Builder, v reflect.Value) {
			b.(*array.DurationBuilder).Append(arrow.Duration(v.Int()))
		}
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int64:
		return func(b array.Builder, v reflect.Value) {
			b.(*array.Int64Builder).Append(v.Int())
		}
	case reflect.String:
		return func(b array.Builder, v reflect.Value) {
			b.(*array.StringBuilder).Append(v.String())
		}
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return func(b array.Builder, v reflect.Value) {
				b.(*array.BinaryBuilder).Append(v.Bytes())
			}
		}
	case reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return func(b array.Builder, v reflect.Value) {
				b.(*array.FixedSizeBinaryBuilder).Append(v.Bytes())
			}
		}
	case reflect.TypeFor[time.Time]().Kind():
		return func(b array.Builder, v reflect.Value) {
			b.(*array.TimestampBuilder).Append(arrow.Timestamp(v.Interface().(time.Time).UnixMilli()))
		}
	case reflect.Map:
		if t.Key().Kind() == reflect.String && t.Elem().Kind() == reflect.String {
			return func(b array.Builder, v reflect.Value) {
				mapBuilder := b.(*array.MapBuilder)
				keyBuilder := mapBuilder.KeyBuilder().(*array.StringBuilder)
				valueBuilder := mapBuilder.ItemBuilder().(*array.StringBuilder)

				mapBuilder.Append(true)
				keys := v.MapKeys()
				for _, key := range keys {
					value := v.MapIndex(key)
					keyBuilder.Append(key.String())
					valueBuilder.Append(value.String())
				}
			}
		}
	}
	return func(b array.Builder, v reflect.Value) { b.AppendNull() }
}
