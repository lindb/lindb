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
	"encoding/json"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./state_machine.go -destination=./state_machine_mock.go -package=client

// StateMachineCli represents state machine explore client.
type StateMachineCli interface {
	// FetchStateByNode fetches the state from state machine by target node.
	FetchStateByNode(params map[string]string, node models.Node) (interface{}, error)
	// FetchStateByNodes fetches the state from state machine by target nodes.
	FetchStateByNodes(params map[string]string, nodes []models.Node) interface{}
}

// stateMachineCli implements StateMachineCli interface.
type stateMachineCli struct {
	logger *logger.Logger
}

// NewStateMachineCli creates a state machine explore client instance.
func NewStateMachineCli() StateMachineCli {
	return &stateMachineCli{
		logger: logger.GetLogger("client", "StateMachine"),
	}
}

// FetchStateByNode fetches the state from state machine by target node.
func (cli *stateMachineCli) FetchStateByNode(params map[string]string, node models.Node) (interface{}, error) {
	address := node.HTTPAddress()
	var r json.RawMessage
	_, err := resty.New().R().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetResult(&r).
		Get(address + "/api/state/machine/explore")
	if err != nil {
		return nil, err
	}
	return r, nil
}

// FetchStateByNodes fetches the state from state machine by target nodes.
func (cli *stateMachineCli) FetchStateByNodes(params map[string]string, nodes []models.Node) interface{} {
	var wait sync.WaitGroup
	wait.Add(len(nodes))
	result := make(map[string]interface{})
	for idx := range nodes {
		i := idx
		go func() {
			defer wait.Done()
			node := nodes[i]
			r, err := cli.FetchStateByNode(params, node)
			if err != nil {
				cli.logger.Error("get state from alive node",
					logger.String("url", node.HTTPAddress()),
					logger.Error(err))
				return
			}
			result[node.Indicator()] = r
		}()
	}
	wait.Wait()
	return result
}
