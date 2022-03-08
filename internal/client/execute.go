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
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/app/broker/api/exec"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/ltoml"
)

// ExecuteCli represents lin query language execute client.
type ExecuteCli struct {
	Base
}

// NewExecuteCli creates a lin query language execute client instance.
func NewExecuteCli(endpoint string) *ExecuteCli {
	cli := resty.New()
	cli.SetBaseURL(endpoint)
	return &ExecuteCli{
		Base{
			cli: cli,
		}}
}

// Execute executes lin query language, then returns execute result.
func (cli *ExecuteCli) Execute(param models.ExecuteParam, rs interface{}) error {
	// send request
	resp, err := cli.cli.R().
		SetBody(param).
		SetHeader("Accept", "application/json").
		Put(exec.ExecutePath)
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
	return fmt.Errorf(resp.Status())
}

// ExecuteAsResult executes lin query language, then returns terminal result.
func (cli *ExecuteCli) ExecuteAsResult(param models.ExecuteParam, rs interface{}) (string, error) {
	n := time.Now()
	err := cli.Execute(param, rs)
	cost := time.Since(n)
	if err != nil {
		return "", err
	}
	rows := getLen(rs)
	if rows == 0 {
		return fmt.Sprintf("Query OK, 0 rows affected (%s)", ltoml.Duration(cost)), nil
	}
	formatter, ok := rs.(models.TableFormatter)
	result := ""
	if ok {
		result = formatter.ToTable() + "\n"
	}
	return fmt.Sprintf("%s%s", result,
		fmt.Sprintf("%d rows in sets (%s)", rows, ltoml.Duration(cost))), nil
}

// getLen returns the length of v.
func getLen(v interface{}) int {
	objValue := reflect.ValueOf(v)
	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len()
	case reflect.Struct:
		return 1
	case reflect.Ptr:
		deref := objValue.Elem().Interface()
		return getLen(deref)
	}
	return 0
}
