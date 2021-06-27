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
	"github.com/lindb/lindb/cmd/lind/cli"

	"github.com/spf13/cobra"
)

const linDBLogo = `
██╗     ██╗███╗   ██╗██████╗ ██████╗ 
██║     ██║████╗  ██║██╔══██╗██╔══██╗
██║     ██║██╔██╗ ██║██║  ██║██████╔╝
██║     ██║██║╚██╗██║██║  ██║██╔══██╗
███████╗██║██║ ╚████║██████╔╝██████╔╝
╚══════╝╚═╝╚═╝  ╚═══╝╚═════╝ ╚═════╝ 
`

const (
	linDBText = `
LinDB is a scalable, high performance, high availability, distributed time series database.
Complete documentation is available at https://lindb.io
`
)

// RootCmd command of cobra
var RootCmd = &cobra.Command{
	Use:   "lind",
	Short: "lind is the main command, used to control LinDB",
	Long:  linDBLogo + linDBText,
}

func init() {
	RootCmd.AddCommand(
		versionCmd,
		newStorageCmd(),
		newBrokerCmd(),
		newStandaloneCmd(),
		cli.NewCLICmd(),
	)
}
