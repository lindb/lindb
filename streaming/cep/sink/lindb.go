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

package sink

import (
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	lindb "github.com/lindb/client_go"
	"github.com/lindb/client_go/api"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/http"
)

type LinDBSink struct {
	write api.Write
}

func newLinDBSink(props *collections.Properties) Sink {
	address, ok := props.GetString("address")
	if !ok {
		return nil
	}
	database, ok := props.GetString("database")
	if !ok {
		return nil
	}
	// create write client with options
	cli := lindb.NewClientWithOptions(
		http.EnsureHTTP(address),
		lindb.DefaultOptions().SetBatchSize(200).
			SetReqTimeout(60).
			SetRetryBufferLimit(100).
			SetFlushInterval(1000).
			SetMaxRetries(3),
	)
	// get write client
	w := cli.Write(database)
	// get error chan
	errCh := w.Errors()
	go func() {
		for err := range errCh {
			fmt.Printf("got err:%s\n", err)
		}
	}()

	return &LinDBSink{
		write: w,
	}
}

func (s *LinDBSink) Publish(record arrow.RecordBatch) {
	fmt.Println("send lindb", record)
	// if points, ok := event.([]*api.Point); ok {
	// 	for _, point := range points {
	// 		s.write.AddPoint(context.TODO(), point)
	// 	}
	// }
	// // TODO implement the logic to publish event to LinDB
	// fmt.Println("Publishing event to LinDB:", event)
}

func (s *LinDBSink) Close() {
	// TODO implement the logic to close LinDB sink
	s.write.Close()
}
