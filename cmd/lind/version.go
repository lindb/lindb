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
	"runtime"

	"github.com/spf13/cobra"

	"github.com/lindb/lindb/config"
)

var (
	// pprof mode
	pprof = false
	// cfg path
	cfg = ""
	// enable swagger api doc
	doc = false
	// storage myid
	myID = 1
	// if enable embed etcd
	embedEtcd = true
)

func printVersion() {
	fmt.Printf("LinDB: %v, BuildDate: %v\n", config.Version, config.BuildTime)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(_ *cobra.Command, _ []string) {
		printVersion()
		fmt.Printf("GOOS=%q\n", runtime.GOOS)
		fmt.Printf("GOARCH=%q\n", runtime.GOARCH)
		fmt.Printf("GOVERSION=%q\n", runtime.Version())
	},
}
