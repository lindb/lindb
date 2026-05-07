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

package client

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	resty "github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

//go:generate mockgen -source=./execute.go -destination=./execute_mock.go -package=client

// ExecuteCli represents lin query language execute client.
type ExecuteCli interface {
	// ExecuteAsRecord executes lin query language, returns result as Arrow RecordBatch.
	// The caller is responsible for calling Release() on the returned RecordBatch.
	ExecuteAsRecord(param models.ExecuteParam) (arrow.RecordBatch, error)
}

// executeCli implements ExecuteCli interface.
type executeCli struct {
	Base
}

// NewExecuteCli creates a lin query language execute client instance.
func NewExecuteCli(endpoint string) ExecuteCli {
	cli := resty.New()
	cli.SetBaseURL(endpoint)
	return &executeCli{
		Base{
			cli: cli,
		},
	}
}

// ExecuteAsRecord executes lin query language, returns the result as an Arrow RecordBatch
// by requesting the Arrow IPC stream format from the server.
// The caller must call Release() on the returned RecordBatch when done.
func (cli *executeCli) ExecuteAsRecord(param models.ExecuteParam) (arrow.RecordBatch, error) {
	resp, err := cli.cli.R().
		SetBody(&param).
		SetHeader("Accept", constants.ContentTypeArrow).
		Put("/exec")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, parseErrorResponse(resp.StatusCode(), resp.Body())
	}
	data := resp.Body()
	// Empty body means the server produced no rows (e.g. query matched nothing).
	// This is a successful empty result set, not an error — same as MySQL's "Empty set".
	if len(data) == 0 {
		return nil, nil
	}
	reader, err := ipc.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Release()

	if !reader.Next() {
		if err := reader.Err(); err != nil {
			return nil, err
		}
		// Arrow IPC stream with no record batches — treat as empty result set.
		return nil, nil
	}
	record := reader.RecordBatch()
	record.Retain()
	return record, nil
}

// parseErrorResponse converts a non-200 HTTP response into a meaningful error.
// When the server returns a JSON null body (e.g. 404 Not Found from gin's NotFound helper),
// the body string is "null" which is useless — fall back to a generic status-based message.
func parseErrorResponse(statusCode int, body []byte) error {
	msg := strings.TrimSpace(string(body))
	// "null" or empty body means the server sent no real message; build one from the status code.
	if msg == "" || msg == "null" {
		msg = fmt.Sprintf("server returned %s", http.StatusText(statusCode))
	}
	return fmt.Errorf("%d %s", statusCode, msg)
}
