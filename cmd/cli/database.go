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
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/c-bata/go-prompt"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/bootstrap"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/option"
)

func emptyCompleter(document prompt.Document) []prompt.Suggest { return nil }

func createDatabase() {
	database := prompt.Input("database name?: ", emptyCompleter)
	if len(database) == 0 {
		_, _ = fmt.Fprint(os.Stderr, errors.New("database name is empty"))
		os.Exit(1)
	}
	numOfShardsString := prompt.Input("number of shards?: ", emptyCompleter)
	numOfShards, err := strconv.ParseInt(numOfShardsString, 10, 64)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	replicaFactorString := prompt.Input("replica factor?: ", emptyCompleter)
	replicaFactor, err := strconv.ParseInt(replicaFactorString, 10, 64)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	intervalString := prompt.Input("write interval?: ", emptyCompleter)
	var d ltoml.Duration
	if err := d.UnmarshalText([]byte(intervalString)); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	initializer := bootstrap.NewClusterInitializer(brokerEndpoint)
	if err := initializer.InitStorageCluster(config.StorageCluster{
		Name:   database,
		Config: cfg.Coordinator},
	); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	if err := initializer.InitInternalDatabase(models.Database{
		Name:          database,
		Storage:       database,
		NumOfShard:    int(numOfShards),
		ReplicaFactor: int(replicaFactor),
		Option: option.DatabaseOption{
			Interval: intervalString,
		},
	}); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
