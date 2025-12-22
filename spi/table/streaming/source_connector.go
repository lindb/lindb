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
	"context"
	"fmt"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/stream/input"
)

type sourceConnectorProvider struct{}

func NewSourceConnectorProvider() spi.SourceConnectorProvider {
	return &sourceConnectorProvider{}
}

func (s *sourceConnectorProvider) CreateSourceConnector(ctx context.Context,
	table spi.TableHandle, partitions []int, columnMapping map[string]string,
	predicate tree.Expression,
	outputColumns []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
) spi.SourceConnector {
	fmt.Printf("create log source connector,table=%v,partitions=%v\n", outputColumns, assignments)

	tableHandle := table.(*TableHandle)
	inputHandle := input.GetManager().GetInputHandler(tableHandle.App, tableHandle.Stream)

	// stream.GetManager().GetStreamManager(tableHandle.App).GetTableMetadata("", "", tableHandle.Stream)

	connector := &sourceConnector{
		ctx:       ctx,
		predicate: predicate,
		inbound:   operator.NewQueue(make(chan *types.Page)),
	}
	inputHandle.Subscribe(connector)
	fmt.Println("create streaming connector")
	return connector
}

type sourceConnector struct {
	ctx context.Context

	predicate tree.Expression

	inbound *operator.Queue
}

func (sc *sourceConnector) Receive(event any) {
	fmt.Printf("stream h connector receiver, receive event:%v\n", event)
	if page, ok := event.(*types.Page); ok {
		sc.inbound.Produce(page)
	}
}

func (sc *sourceConnector) Run(output chan<- *types.Page) {
	v := &visitor{}
	for {
		source, ok := sc.inbound.Consume(sc.ctx)
		if !ok {
			break
		}
		if sc.predicate == nil {
			output <- source
			continue
		}
		fmt.Println("get event")
		if val, ok := v.Visit(source, sc.predicate).(bool); val && ok {
			output <- source
		}
	}
}

type visitor struct{}

func (v *visitor) Visit(event any, n tree.Node) any {
	// 	switch node := n.(type) {
	// 	case *tree.ComparisonExpression:
	// 		v, err := GetFieldValue(event, getValue(node.Left))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return false
	// 		}
	// 		return v == getValue(node.Right)
	// 	case *tree.InPredicate:
	// 		var values []string
	// 		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
	// 			values = lo.Map(inListExpression.Values, func(item tree.Expression, index int) string {
	// 				return getValue(item)
	// 			})
	// 		}
	// 		v, err := GetFieldValue(event, getValue(node.Value))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return false
	// 		}
	// 		return lo.Contains(values, v.(string))
	// 	case *tree.LogicalExpression:
	// 		for _, term := range node.Terms {
	// 			val, ok := term.Accept(event, v).(bool)
	// 			if !ok {
	// 				return false
	// 			}
	// 			if node.Operator == tree.LogicalOR && val {
	// 				return true
	// 			} else if !val {
	// 				return false
	// 			}
	// 		}
	// 		return true
	// 	default:
	// 		panic(fmt.Errorf("not support,%T", n))
	// 	}
	return true
}
