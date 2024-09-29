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
	"time"

	prompt "github.com/c-bata/go-prompt"
	"github.com/fatih/color"
	commonlogger "github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

// for testing
var (
	urlParse      = url.Parse
	newExecuteCli = client.NewExecuteCli
	runPromptFn   = runPrompt
	exit          = os.Exit
	newPrompt     = prompt.New
)

const (
	HTTPScheme = "http://"
)

type inputCtx struct {
	db string
}

var (
	endpoint string
	// tokens represents suggest token.
	tokens = []prompt.Suggest{
		{Text: "show"},
		{Text: "use"},
		{Text: "alive"},
		{Text: "master"},
		{Text: "storages"},
		{Text: "storage"},
		{Text: "broker"},
		{Text: "schemas"},
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
	spacesPattern = regexp.MustCompile(`\s+`)
	inputC        = &inputCtx{}
	query         = ""
	live          = true
	cli           client.ExecuteCli
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:9000", "Broker HTTP Endpoint")
}

// printErr prints error message.
func printErr(err error) {
	fmt.Println(color.RedString("ERROR:%s", err))
}

// executor executes command.
func executor(in string) {
	in = strings.TrimSpace(in)
	if strings.HasSuffix(in, ";") {
		query += in
		live = true

		defer func() {
			// reset query input buf
			query = ""
		}()

		query = strings.TrimSpace(strings.TrimSuffix(query, ";"))
		if query == "" {
			return
		}
		blocks := strings.Split(spacesPattern.ReplaceAllString(query, " "), " ")
		switch blocks[0] {
		case "exit":
			fmt.Println("Good Bye :)")
			exit(0)
			return
		case "use":
			if len(blocks) == 1 || strings.TrimSpace(blocks[0]) == "" {
				printErr(errors.New("database is required"))
				return
			}
			inputC.db = strings.TrimSpace(blocks[1])
			fmt.Println(color.GreenString("Database changed(current:%s)", inputC.db))
			return
		default:
			// query and print result in terminal
			executeAndPrint(models.ExecuteParam{SQL: query, Database: inputC.db})
		}
		return
	}

	query += strings.TrimSuffix(in, "\r\n") + " "
	live = false
}

func executeAndPrint(param models.ExecuteParam) {
	defer func() {
		if err0 := recover(); err0 != nil {
			printErr(fmt.Errorf("query error: %v", err0))
		}
	}()

	n := time.Now()
	rs, err := cli.Execute(param)
	cost := time.Since(n)
	if err != nil {
		printErr(err)
		return
	}
	if len(rs.Rows) == 0 {
		fmt.Println(color.GreenString("Query OK, 0 rows affected (%s)", ltoml.Duration(cost)))
		return
	}
	fmt.Printf("%s\n%s\n", rs.ToTable(),
		color.GreenString("%d rows in sets (%s)", len(rs.Rows), ltoml.Duration(cost)))
}

// completer returns prompt suggest.
func completer(bc string) []prompt.Suggest {
	if bc == "" {
		return nil
	}
	args := strings.Split(spacesPattern.ReplaceAllString(bc, " "), " ")
	cmdName := args[len(args)-1]
	return prompt.FilterHasPrefix(tokens, cmdName, true)
}

func main() {
	flag.Parse()

	if !strings.HasPrefix(endpoint, HTTPScheme) {
		endpoint = fmt.Sprintf("%s%s", HTTPScheme, endpoint)
	}
	endpointURL, err := urlParse(endpoint)
	if err != nil {
		printErr(err)
		return
	}
	commonlogger.IsCli = true
	_ = logger.InitLogger(commonlogger.Setting{Level: "error", Dir: "."}, "lin-cli.log")

	apiEndpoint := endpoint + constants.APIVersion1CliPath
	cli = newExecuteCli(apiEndpoint)

	// first retry connect and get master state
	rs, err := cli.Execute(models.ExecuteParam{Database: "information_schema", SQL: "select version from master"})
	if err != nil {
		printErr(err)
		return
	}
	if len(rs.Rows) == 0 {
		printErr(errors.New("no master found"))
		return
	}
	fmt.Println("Welcome to the LinDB.")
	fmt.Printf("Server version: %s\n", rs.Rows[0][0])
	endpointStr := fmt.Sprintf("lin@%s", endpointURL.Host)
	var spaces []string
	for i := 2; i < len(endpointStr); i++ {
		spaces = append(spaces, " ")
	}
	prefix := strings.Join(spaces, "")

	p := newPrompt(
		executor,
		func(document prompt.Document) []prompt.Suggest {
			return completer(document.TextBeforeCursor())
		},
		prompt.OptionLivePrefix(func() (string, bool) {
			if live {
				return endpointStr + "> ", true
			}
			return prefix + "- > ", true
		}),
		prompt.OptionTitle("LinDB Client"),

		prompt.OptionPrefixTextColor(prompt.Blue),
		prompt.OptionInputTextColor(prompt.White),

		prompt.OptionSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionDescriptionBGColor(prompt.White),
		prompt.OptionDescriptionTextColor(prompt.Black),

		prompt.OptionSelectedSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.Blue),
		prompt.OptionSelectedDescriptionTextColor(prompt.Black),
	)
	runPromptFn(p)
}

func runPrompt(p *prompt.Prompt) {
	p.Run()
}
