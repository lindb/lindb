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

package streaming

import (
	contextpkg "context"
	"fmt"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/execution"
	"github.com/lindb/lindb/sql/execution/model"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	printpkg "github.com/lindb/lindb/sql/planner/printer"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/stream"
	"github.com/lindb/lindb/streaming/stream/input"
	"github.com/lindb/lindb/streaming/stream/output"
)

func init() {
	spi.RegisterSourceConnectorProvider(&stream.TableHandle{}, stream.NewSourceConnectorProvider())
}

type Runtime interface {
	AddEventType(eventType any)
	Query(sql string) error
	Startup()
	Shutdown()
}

type runtime struct{}

func NewRuntime() *runtime {
	return &runtime{}
}

func (r *runtime) RegisterStreamByType(eventType any) error {
	return stream.GetManager().GetStreamManager("test").RegisterStreamByType(eventType)
}

func (r *runtime) AddListener(stream string, listener output.Listener) {
	r.GetInputHandler(stream).Subscribe(listener)
}

func (r *runtime) GetInputHandler(stream string) input.InputHandler {
	// return r.inputManager.GetInputHandler(stream)
	return input.GetManager().GetInputHandler("test", stream)
}

func (r *runtime) Query(sql string) error {
	idAllocator := tree.NewNodeIDAllocator()
	stmt, err := tree.GetParser().CreateStatement(sql, idAllocator)
	if err != nil {
		return err
	}

	switch node := stmt.(type) {
	case *tree.StreamingApp:
		fmt.Println(node.Statements)
		fmt.Printf("stat....===%v,%v\n", len(node.Statements), string(encoding.JSONMarshal(node.Annotations)))
		for _, stmt := range node.Statements {
			r.plan(idAllocator, stmt.Statement)
		}
	default:
		r.plan(idAllocator, stmt)
	}
	return nil
}

func (r *runtime) plan(idAllocator *tree.NodeIDAllocator, statement tree.Statement) {
	planner := execution.NewPlanner(analyzer.NewAnalyzerFactory(stream.GetManager().GetStreamManager("test")))
	plan := planner.Plan(&execution.Session{
		Database:        "test", // fixme: set streaming database
		Context:         contextpkg.TODO(),
		NodeIDAllocator: idAllocator,
		Streaming:       true,
	}, statement, execution.StreamingPlanOptimizers())

	printer := printpkg.NewPlanPrinter(printpkg.NewTextRender(0))
	fmt.Printf("final plan:\n%s\n", printer.PrintLogicPlan(plan.Root))

	taskFct := execution.NewTaskExecutionFactory()
	exec := taskFct.Create(&execution.SQLTask{
		ID: model.TaskID{},
		Fragment: &planpkg.PlanFragment{
			Root: plan.Root,
		},
		Streaming: true,
	})
	go exec.Execute(nil)
}
