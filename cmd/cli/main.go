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

	prompt "github.com/c-bata/go-prompt"
	"github.com/fatih/color"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/sql"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
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

		query = strings.TrimSuffix(query, ";")
		if query == "" {
			return
		}
		blocks := strings.Split(spacesPattern.ReplaceAllString(query, " "), " ")
		switch blocks[0] {
		case "exit":
			fmt.Println("Good Bye :)")
			exit(0)
			return
		default:
			stmt, err := sql.Parse(query)
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
			case *stmtpkg.Storage:
				if s.Type == stmtpkg.StorageOpShow {
					result = &models.Storages{}
				}
			case *stmtpkg.State:
				// execute state query
				switch s.Type {
				case stmtpkg.Master:
					result = &models.Master{}
				case stmtpkg.BrokerAlive:
					result = &models.StatelessNodes{}
				}
			case *stmtpkg.Schema:
				switch s.Type {
				case stmtpkg.DatabaseNameSchemaType:
					result = &models.DatabaseNames{}
				case stmtpkg.DatabaseSchemaType:
					result = &models.Databases{}
				}
			case *stmtpkg.MetricMetadata:
				if strings.TrimSpace(inputC.db) == "" {
					printErr(errors.New("please select database(use ...)"))
					return
				}
				result = &models.Metadata{}
			case *stmtpkg.Query:
				result = &models.ResultSet{}
				if strings.TrimSpace(inputC.db) == "" {
					printErr(errors.New("please select database(use ...)"))
					return
				}
			}
			rs, err := cli.ExecuteAsResult(models.ExecuteParam{SQL: query, Database: inputC.db}, result)
			if err != nil {
				printErr(err)
				return
			}
			// print result in terminal
			fmt.Println(rs)
		}
		return
	}

	query += strings.TrimSuffix(in, "\r\n") + " "
	live = false
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
	logger.IsCli = true
	_ = logger.InitLogger(config.Logging{Level: "error", Dir: "."}, "lin-cli.log")

	apiEndpoint := endpoint + constants.APIVersion1CliPath
	cli = newExecuteCli(apiEndpoint)

	// first retry connect and get master state
	master := &models.Master{}
	err = cli.Execute(models.ExecuteParam{SQL: "show master"}, master)
	if err != nil || master.Node == nil {
		printErr(err)
		return
	}
	fmt.Println("Welcome to the LinDB.")
	fmt.Printf("Server version: %s\n", master.Node.Version)
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
	runPromptFn(p)
}

func runPrompt(p *prompt.Prompt) {
	p.Run()
}
