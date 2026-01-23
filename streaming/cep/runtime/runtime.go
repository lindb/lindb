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
	contextpkg "context"
	"fmt"
	"strings"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/analyzer"
	"github.com/lindb/lindb/sql/execution"
	"github.com/lindb/lindb/sql/execution/model"
	"github.com/lindb/lindb/sql/expression"
	planpkg "github.com/lindb/lindb/sql/planner/plan"
	printpkg "github.com/lindb/lindb/sql/planner/printer"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/cep/annotation"
	"github.com/lindb/lindb/streaming/cep/sink"
	"github.com/lindb/lindb/streaming/cep/stream"
	"github.com/lindb/lindb/streaming/cep/stream/input"
	"github.com/lindb/lindb/streaming/cep/stream/output"
)

func init() {
	spi.RegisterSourceConnectorProvider(&stream.TableHandle{}, stream.NewSourceConnectorProvider())
}

type Runtime interface {
	AddEventType(eventType any)
	Query(sql string) error
	RegisterStreamByType(eventType any) error
	RegisterStreamBySchema(name string, schema *types.TableSchema) error
	AddListener(stream string, listener output.Listener)
	GetInputHandler(stream string) input.InputHandler
	Startup()
	Shutdown()
}

type runtime struct {
	database string

	sinks map[string]sink.Sink
}

// AddEventType implements [Runtime].
func (r *runtime) AddEventType(eventType any) {
	panic("unimplemented")
}

// Shutdown implements [Runtime].
func (r *runtime) Shutdown() {
	panic("unimplemented")
}

// Startup implements [Runtime].
func (r *runtime) Startup() {
	panic("unimplemented")
}

func NewRuntime(database string) Runtime {
	return &runtime{
		database: database,
		sinks:    make(map[string]sink.Sink),
	}
}

func (r *runtime) RegisterStreamByType(eventType any) error {
	return stream.GetManager().GetStreamManager(r.database).RegisterStreamByType(eventType)
}

func (r *runtime) RegisterStreamBySchema(name string, schema *types.TableSchema) error {
	return stream.GetManager().GetStreamManager(r.database).RegisterStreamBySchema(name, schema)
}

func (r *runtime) AddListener(stream string, listener output.Listener) {
	r.GetInputHandler(stream).Subscribe(listener)
}

func (r *runtime) GetInputHandler(stream string) input.InputHandler {
	return input.GetManager().GetInputHandler(r.database, stream)
}

func (r *runtime) Query(sql string) error {
	idAllocator := tree.NewNodeIDAllocator()
	stmt, err := tree.GetParser().CreateStatement(sql, idAllocator)
	if err != nil {
		return err
	}
	fmt.Printf("parsed statement: %T\n", stmt)

	switch node := stmt.(type) {
	case *tree.StreamingApp:
		r.deployStreaming(node, idAllocator)
	case *tree.CreateJob:
		if node.Streaming != nil {
			r.deployStreaming(node.Streaming, idAllocator)
		}
	default:
		r.deploy(stmt, idAllocator, nil)
	}
	return nil
}

func (r *runtime) deployStreaming(stmt *tree.StreamingApp, idAllocator *tree.NodeIDAllocator) {
	// create data sinks
	for _, stmt := range stmt.CreateSinks {
		r.creaetSink(stmt)
	}
	// dploy streaming query statements
	for _, stmt := range stmt.Statements {
		r.deploy(stmt.Statement, idAllocator, stmt.Annotations)
	}
}

func (r *runtime) deploy(statement tree.Statement, idAllocator *tree.NodeIDAllocator, annotations []*tree.Annotation) {
	// TODO: generate stream name for statement
	planner := execution.NewPlanner(analyzer.NewAnalyzerFactory(stream.GetManager().GetStreamManager(r.database)))
	plan := planner.Plan(&execution.Session{
		Database:        r.database,
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
		Database:   r.database,
		StreamName: statement.String(),
	})
	sinkBridges, err := r.parseSinkBridges(annotations)
	if err != nil {
		panic(err)
	}
	// add listener to input handler for this statement
	output := input.GetManager().GetInputHandler(r.database, statement.String())
	output.Subscribe(NewListener(sinkBridges))

	// run streaming execution
	go exec.Execute(nil)
}

func (r *runtime) creaetSink(statment *tree.CreateSink) {
	name := statment.Name
	props, err := expression.EvalProps(expression.NewEvalContext(contextpkg.TODO()), statment.Props)
	if err != nil {
		panic(err)
	}
	sink := sink.CreateSink(props)
	if sink == nil {
		panic(fmt.Errorf("unsupported sink type"))
	}
	r.sinks[name] = sink
}

func (r *runtime) parseSinkBridges(annotations []*tree.Annotation) ([]*sink.SinkBridge, error) {
	var sinkBridges []*sink.SinkBridge
	for _, ann := range annotations {
		annotationMeta := annotation.ParseAnnotation(ann)

		if !strings.EqualFold(annotationMeta.Name, "sink") {
			continue
		}
		sinkName, ok := annotationMeta.Props.GetString("name")
		if !ok {
			return nil, fmt.Errorf("sink annotation must contain 'name' property")
		}

		s, ok := r.sinks[sinkName]
		if !ok {
			return nil, fmt.Errorf("sink '%s' not found", sinkName)
		}

		sinkBridges = append(sinkBridges, sink.NewSinkBridge(s, annotationMeta))
	}

	return sinkBridges, nil
}
