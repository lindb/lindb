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
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/fatih/color"

	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var (
	endpoint string
	tokens   = []prompt.Suggest{
		{Text: "show"},
		{Text: "use"},
		{Text: "master"},
		{Text: "database"},
		{Text: "databases"},
		{Text: "group by"},
		{Text: "select"},
		{Text: "from"},
		{Text: "where"},
		{Text: "namespaces"},
		{Text: "namespace"},
		{Text: "metrics"},
		{Text: "metric"},
		{Text: "fields"},
		{Text: "tag"},
		{Text: "with"},
		{Text: "keys"},
		{Text: "key"},
		{Text: "values"},
		{Text: "and"},
	}
	spaces = regexp.MustCompile(`\s+`)
)

const (
	HTTPScheme = "http://"
)

type inputCtx struct {
	db string
}

func init() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:9000", "Broker HTTP Endpoint")
}

func printErr(err error) {
	fmt.Println(color.RedString("ERROR:%s", err))
}

func main() {
	flag.Parse()
	var history []string

	if !strings.HasPrefix(endpoint, HTTPScheme) {
		endpoint = fmt.Sprintf("%s%s", HTTPScheme, endpoint)
	}
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		printErr(err)
		return
	}

	apiEndpoint := fmt.Sprintf("%s/api", endpoint)
	cli := client.NewExecuteCli(apiEndpoint)

	// first retry connect and get master state
	master := &models.Master{}
	err = cli.Execute(models.ExecuteParam{SQL: "show master"}, &master)
	if err != nil || master.Node == nil {
		printErr(err)
		return
	}
	fmt.Println("Welcome to the LinDB.")
	fmt.Printf("Server version: %s\n", master.Node.Version)
	inputC := &inputCtx{}

	p := prompt.New(
		func(in string) {
			in = strings.TrimSpace(in)
			history = append(history, in)
			blocks := strings.Split(spaces.ReplaceAllString(in, " "), " ")
			switch blocks[0] {
			case "exit":
				fmt.Println("Good Bye :)")
				os.Exit(0)
			default:
				stmt, err := sql.Parse(in)
				if err != nil {
					printErr(err)
					return
				}
				var result interface{}
				switch s := stmt.(type) {
				case *stmtpkg.Use:
					inputC.db = s.Name
					fmt.Printf("Database changed(current:%s)\n", inputC.db)
					return
				case *stmtpkg.State:
					// execute state query
					if s.Type == stmtpkg.Master {
						result = &models.Master{}
					}
				case *stmtpkg.Metadata:
					result = &models.Metadata{}
				case *stmtpkg.Query:
					result = &models.ResultSet{}
					if strings.TrimSpace(inputC.db) == "" {
						printErr(errors.New("please select database(use ...)"))
						return
					}
				}
				rs, err := cli.ExecuteAsResult(models.ExecuteParam{SQL: in, Database: inputC.db}, result)
				if err != nil {
					printErr(err)
					return
				}
				// print result in terminal
				fmt.Println(rs)
			}
		},
		func(d prompt.Document) []prompt.Suggest {
			bc := d.TextBeforeCursor()
			if bc == "" {
				return nil
			}
			args := strings.Split(spaces.ReplaceAllString(bc, " "), " ")
			cmdName := args[len(args)-1]
			return prompt.FilterHasPrefix(tokens, cmdName, true)
		},
		prompt.OptionLivePrefix(func() (prefix string, useLivePrefix bool) {
			return fmt.Sprintf("lin@%s> ", endpointUrl.Host), true
		}),
		prompt.OptionHistory(history),

		prompt.OptionPrefixTextColor(prompt.DarkGreen),
		prompt.OptionInputTextColor(prompt.DarkBlue),

		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.White),
		prompt.OptionDescriptionTextColor(prompt.Black),

		prompt.OptionSelectedSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.Blue),
		prompt.OptionSelectedDescriptionTextColor(prompt.Black),
	)

	p.Run()
}
