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

package cep

import (
	"fmt"
	"testing"
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/streaming/cep/stream/output"
)

type RPCService struct {
	Interface string
	Tags      map[string]string
	Status    string
}

type Result struct{}

func Test_Runtime(t *testing.T) {
	runtime := NewRuntime("test")
	runtime.RegisterStreamByType(RPCService{})
	runtime.RegisterStreamByType(Result{})
	// add result listener
	runtime.AddListener("Result", output.NewConsoleOutput())
	// add streaming query
	err := runtime.Query(`
	@app(name="test_app")
	@name(name="count_rpc",@header(user="test_user",pwd="pwd"))
	insert into Result
	select map_values(tags,'app') as tags_map,interface,count(1) as qps
	from RPCService
	where interface in('grpc','http')
	group by map_values(tags,'app'),interface;
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
	tagsColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTMap, Name: "tags"}, tagsColumn)
	statusColumn := types.NewColumn()
	page.AppendColumn(types.ColumnMetadata{DataType: types.DTMap, Name: "status"}, statusColumn)

	interfaceColumn.AppendString("grpc")
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "order"})

	interfaceColumn.AppendString("http")
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "user"})

	interfaceColumn.AppendString("dubbo")
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "order"})

	interfaceColumn.AppendString("http")
	tagsColumn.Append(map[string]string{"host": "1.1.1.1", "app": "github"})

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
