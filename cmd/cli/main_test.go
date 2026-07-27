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

package main

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	prompt "github.com/elk-language/go-prompt"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/internal/client"
)

func Test_main(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cli := client.NewMockExecuteCli(ctrl)

	endpoint = "localhost/url"
	cases := []struct {
		name    string
		prepare func()
	}{
		{
			name: "url parse failure",
			prepare: func() {
				urlParse = func(rawURL string) (*url.URL, error) {
					return nil, fmt.Errorf("err")
				}
			},
		},
		{
			name: "get master failure",
		},
		{
			name: "get master successfully",
			prepare: func() {
				newExecuteCli = func(endpoint string) client.ExecuteCli {
					return cli
				}
				cli.EXPECT().ExecuteAsRecord(gomock.Any()).DoAndReturn(func(_ any) (arrow.RecordBatch, error) {
					// Return a RecordBatch with one string column "version" = "1.0.0"
					schema := arrow.NewSchema([]arrow.Field{
						{Name: "version", Type: arrow.BinaryTypes.String},
					}, nil)
					rb := array.NewRecordBuilder(memory.NewGoAllocator(), schema)
					defer rb.Release()
					rb.Field(0).(*array.StringBuilder).Append("1.0.0")
					return rb.NewRecordBatch(), nil
				})

				runPromptFn = func(p *prompt.Prompt) {
				}
				newPrompt = func(executor prompt.Executor,
					opts ...prompt.Option,
				) *prompt.Prompt {
					return nil
				}
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(_ *testing.T) {
			defer func() {
				urlParse = url.Parse
				newExecuteCli = client.NewExecuteCli
				runPromptFn = runPrompt
				newPrompt = prompt.New
			}()
			if tt.prepare != nil {
				tt.prepare()
			}

			main()
		})
	}
}

func Test_executor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCli := client.NewMockExecuteCli(ctrl)
	cli = mockCli

	cases := []struct {
		name    string
		in      string
		prepare func()
	}{
		{
			name: "exit",
			in:   "exit;",
			prepare: func() {
				exitFn = func(code int) {}
			},
		},
		{
			name: "use empty database",
			in:   "use;",
		},
		{
			name: "use database",
			in:   "use database;",
		},
		{
			name: "history",
			in:   "history;",
		},
		{
			name: "show master",
			in:   "select * from master;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsRecord(gomock.Any())
			},
		},
		{
			name: "show brokers",
			in:   "select * from brokers;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsRecord(gomock.Any())
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				executor(";")
				exitFn = os.Exit
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			executor(tt.in)
		})
	}
}
