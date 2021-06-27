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

package lind

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// These variables are populated via the Go linker.
var (
	// release version, ldflags
	version = ""
	// binary build-time, ldflags
	buildTime = "unknown"
	// debug mode
	debug = false
	// cfg path
	cfg = ""
)

const defaultVersion = "0.0.0"

func getVersion() string {
	if version == "" {
		return defaultVersion
	}
	return version
}

func printVersion() {
	fmt.Printf("LinDB: %v, BuildDate: %v\n", getVersion(), buildTime)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
		fmt.Printf("GOOS=%q\n", runtime.GOOS)
		fmt.Printf("GOARCH=%q\n", runtime.GOARCH)
		fmt.Printf("GOVERSION=%q\n", runtime.Version())
	},
}
