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
	"fmt"
	"net/http"
	"time"

	resty "github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
)

//go:generate mockgen -source=./execute.go -destination=./execute_mock.go -package=client

// ExecuteCli represents lin query language execute client.
type ExecuteCli interface {
	// Execute executes lin query language, then returns execute result.
	Execute(param models.ExecuteParam, rs interface{}) error
	// ExecuteAsResult executes lin query language, then returns terminal result.
	ExecuteAsResult(param models.ExecuteParam, rs interface{}) (string, error)
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
		}}
}

// Execute executes lin query language, then returns execute result.
func (cli *executeCli) Execute(param models.ExecuteParam, rs interface{}) error {
	// send request
	resp, err := cli.cli.R().
		SetBody(&param).
		SetHeader("Accept", "application/json").
		Put("/exec")
	if err != nil {
		return err
	}
	if resp.StatusCode() == http.StatusOK {
		// if success, unmarshal response
		data := resp.Body()
		if rs != nil && len(data) > 0 {
			if err := encoding.JSONUnmarshal(data, &rs); err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New(string(resp.Body()))
}

// ExecuteAsResult executes lin query language, then returns terminal result.
func (cli *executeCli) ExecuteAsResult(param models.ExecuteParam, rs interface{}) (string, error) {
	n := time.Now()
	err := cli.Execute(param, rs)
	cost := time.Since(n)
	if err != nil {
		return "", err
	}
	result := ""
	rows := 0
	if formatter, ok := rs.(models.TableFormatter); ok {
		rows, result = formatter.ToTable()
	}
	if rows == 0 {
		return fmt.Sprintf("Query OK, 0 rows affected (%s)", ltoml.Duration(cost)), nil
	}
	if result != "" {
		result += "\n"
	}
	return fmt.Sprintf("%s%s", result,
		fmt.Sprintf("%d rows in sets (%s)", rows, ltoml.Duration(cost))), nil
}
