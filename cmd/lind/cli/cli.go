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

package cli

import (
	"github.com/spf13/cobra"
)

var (
	cliBrokerEndpoint string
	cliConfigFile     string
)

// NewCLICmd returns a new command line interface for communicating with lindb
func NewCLICmd() *cobra.Command {
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "LinDB command line interface",
	}

	cliCmd.PersistentFlags().StringVar(
		&cliBrokerEndpoint, "endpoint", "http://localhost:9000", "endpoint of any broker")
	cliCmd.PersistentFlags().StringVar(
		&cliConfigFile, "config", "./storage.toml", "path of a storage toml")

	cliCmd.AddCommand(
		createDatabaseCmd,
		listDatabaseCmd,
		getDatabaseCmd,
		deleteDatabaseCmd,
		addUserCmd,
		listUserCmd,
		getUserCmd,
		deleteUserCmd,
	)

	return cliCmd
}
