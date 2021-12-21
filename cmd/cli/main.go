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
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/ltoml"

	"github.com/c-bata/go-prompt"
	etcdcliv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var (
	cfgPath        string
	brokerEndpoint string
	cfg            coordinatorCfg
	etcdClient     *etcdcliv3.Client
)

type coordinatorCfg struct {
	Coordinator config.RepoState `toml:"coordinator"`
}

func init() {
	flag.StringVar(&brokerEndpoint, "endpoint", "http://localhost:9000", "Broker HTTP Endpoint")
	flag.StringVar(&cfgPath, "config", "broker.toml", "Use either toml config of broker or storage")
}

const (
	createDatabaseSuggest = "create database"
	showDatabaseSuggest   = "show databases"
	etcdSuggest           = "etcd"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "create database", Description: "Creates the database for LinDB"},
		{Text: "show databases", Description: "Lists all databases of LinDB"},
		{Text: "etcd", Description: "Show meta data stored in ETCD"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func parseConfig() error {
	if cfgPath == "" {
		return fmt.Errorf("config path is empty")
	}
	if err := ltoml.LoadConfig(cfgPath, "", &cfg); err != nil {
		return fmt.Errorf("decode config file error: %s", err)
	}
	etcdCfg := etcdcliv3.Config{
		Endpoints:            cfg.Coordinator.Endpoints,
		DialTimeout:          time.Second * 5,
		DialKeepAliveTime:    time.Second * 5,
		DialKeepAliveTimeout: time.Second * 5,
		DialOptions:          []grpc.DialOption{grpc.WithBlock()},
		Username:             cfg.Coordinator.Username,
		Password:             cfg.Coordinator.Password,
	}
	var err error
	etcdClient, err = etcdcliv3.New(etcdCfg)
	if err != nil {
		return fmt.Errorf("create etcd client error: %s", err)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := parseConfig(); err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Please select a Command ")
	chosen := prompt.Input("> ", completer)
	switch chosen {
	case createDatabaseSuggest:
		createDatabase()
	case showDatabaseSuggest:
	case etcdSuggest:
		etcdCtx = &etcdContext{path: []string{cfg.Coordinator.Namespace}}
		p := prompt.New(
			etcdExecutor,
			etcdCompleter,
			prompt.OptionPrefix("[ETCD] >> "),
			prompt.OptionLivePrefix(etcdLivePrefix),
			prompt.OptionTitle("etcd-prompt"),
		)
		p.Run()
	}
}
