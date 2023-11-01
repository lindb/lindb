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

	prompt "github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
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
				cli.EXPECT().Execute(gomock.Any(), gomock.Any()).
					DoAndReturn(func(param models.ExecuteParam, rs interface{}) error {
						m := rs.(*models.Master)
						m.Node = &models.StatelessNode{}
						return nil
					})

				runPromptFn = func(p *prompt.Prompt) {
				}
				newPrompt = func(executor prompt.Executor, completer prompt.Completer,
					opts ...prompt.Option) *prompt.Prompt {
					return nil
				}
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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
				exit = func(code int) {}
			},
		},
		{
			name: "use database",
			in:   "use database;",
		},
		{
			name: "show master",
			in:   "show master;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "show storages",
			in:   "show storages;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "show broker alive",
			in:   "show broker alive;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "show databases",
			in:   "show databases;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "show schemas",
			in:   "show schemas;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "show namespaces",
			in:   "show namespaces;",
			prepare: func() {
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "show metrics, but not use database",
			in:   "show metrics;",
			prepare: func() {
				inputC.db = ""
			},
		},
		{
			name: "select query, but not use database",
			in:   "select f from cpu;",
			prepare: func() {
				inputC.db = ""
			},
		},
		{
			name: "select query, but execute failure",
			in:   "select f from cpu;",
			prepare: func() {
				inputC.db = "test"
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any()).Return("", fmt.Errorf("err"))
			},
		},
		{
			name: "select query successfully",
			in:   "select f from cpu;",
			prepare: func() {
				inputC.db = "test"
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "from query successfully",
			in:   "from cpu select f;",
			prepare: func() {
				inputC.db = "test_2"
				mockCli.EXPECT().ExecuteAsResult(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "parse query sql failure",
			in:   "select f;",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				executor(";")
				exit = os.Exit
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			executor(tt.in)
		})
	}
}

func Test_completer(t *testing.T) {
	assert.Nil(t, completer(""))
	assert.NotEmpty(t, completer("s"))
}
