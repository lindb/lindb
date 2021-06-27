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

import "github.com/spf13/cobra"

var addUserCmd = &cobra.Command{
	Use:   "user-add",
	Short: "Adds a new admin user",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

var listUserCmd = &cobra.Command{
	Use:   "user-list",
	Short: "Lists all users",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

var getUserCmd = &cobra.Command{
	Use:   "user-get",
	Short: "Gets detailed information of a user",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user-delete",
	Short: "Deletes a user",
	Run: func(cmd *cobra.Command, args []string) {
		// todo: @codingcrush
	},
}
