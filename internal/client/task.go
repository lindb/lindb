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
	"errors"
	"net/http"

	resty "github.com/go-resty/resty/v2"
	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/sql/execution/model"
)

type TaskCli interface {
	Submit(task *model.TaskRequest) error
}

type taskCli struct {
	Base
}

func NewTaskCli(endpoint string) TaskCli {
	cli := resty.New()
	cli.SetBaseURL(endpoint)
	return &taskCli{
		Base{
			cli: cli,
		},
	}
}

func (cli *taskCli) Submit(task *model.TaskRequest) error {
	resp, err := cli.cli.R().
		SetBody(encoding.JSONMarshal(task)).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Post("/task")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		return nil
	}
	return errors.New(string(resp.Body()))
}
