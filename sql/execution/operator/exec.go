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

package operator

import (
	"context"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/pkg/logger"
)

// RunAsync runs op.Run in a background goroutine.
// Any panic is recovered and sent to errCh so callers can propagate it back to
// the HTTP layer. The output channel is always closed when the goroutine exits.
func RunAsync(ctx context.Context, op Operator, output chan<- arrow.RecordBatch, errCh chan<- error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("run operator panic", logger.Any("error", r), logger.Stack())
				// Forward the panic as an error so the pipeline/task-execution layer
				// can surface it to the caller instead of silently swallowing it.
				select {
				case errCh <- fmt.Errorf("%v", r):
				default:
				}
			}
			close(output)
		}()

		op.Run(ctx, output)
	}()
}
