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
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lindb/lindb/sql/grammar"
)

var keyWordsCmd = &cobra.Command{
	Use:   "keywords",
	Short: "Print the keywords for lin query language",
	Run: func(_ *cobra.Command, _ []string) {
		typ := reflect.TypeOf(&grammar.NonReservedWordsContext{})
		var keyWords []string
		for i := 0; i < typ.NumMethod(); i++ {
			methodName := typ.Method(i).Name
			if strings.HasPrefix(methodName, "T_") {
				keyWords = append(keyWords, strings.ReplaceAll(methodName, "T_", ""))
			}
		}
		sort.Strings(keyWords)
		count := 0
		for _, name := range keyWords {
			fmt.Printf("%-15s", name)
			count++
			if count%6 == 0 {
				fmt.Printf("\n")
			}
		}
		fmt.Println()
	},
}
