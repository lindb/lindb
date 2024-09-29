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

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/execution/model"
)

//go:generate mockgen -source=./execute.go -destination=./execute_mock.go -package=client

// ExecuteCli represents lin query language execute client.
type ExecuteCli interface {
	// Execute executes lin query language, then returns execute result.
	Execute(param models.ExecuteParam) (*model.ResultSet, error)
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

// Execute executes lin query language, then returns execute result.
func (cli *executeCli) Execute(param models.ExecuteParam) (*model.ResultSet, error) {
	// send request
	resp, err := cli.cli.R().
		SetBody(&param).
		SetHeader("Accept", "application/json").
		Put("/exec")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusOK {
		// if success, unmarshal response
		data := resp.Body()
		if len(data) > 0 {
			rs := &model.ResultSet{}
			if err := encoding.JSONUnmarshal(data, rs); err != nil {
				return nil, err
			}
			return rs, nil
		}
		return nil, errors.New("no data found")
	}
	return nil, errors.New(string(resp.Body()))
}
