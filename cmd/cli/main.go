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
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/go-prompt"
	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/fatih/color"
	commonlogger "github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/ltoml"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/sql/grammar"
)

// for testing
var (
	urlParse      = url.Parse
	newExecuteCli = client.NewExecuteCli
	runPromptFn   = runPrompt
	newPrompt     = prompt.New
)

const (
	HTTPScheme = "http://"
)

type inputCtx struct {
	db string
}

var (
	endpoint      string
	spacesPattern = regexp.MustCompile(`\s+`)
	inputC        = &inputCtx{}
	query         = ""
	live          = true
	cli           client.ExecuteCli
	suggestItems  = suggestTokens()
)

// suggestTokens returns prompt suggest tokens.
func suggestTokens() (tokens []prompt.Suggest) {
	typ := reflect.TypeOf(&grammar.NonReservedContext{})
	var keyWords []string
	for i := 0; i < typ.NumMethod(); i++ {
		methodName := typ.Method(i).Name
		if strings.ToUpper(methodName) == methodName {
			keyWords = append(keyWords, methodName)
		}
	}
	keyWords = append(keyWords, "EXIT")
	keyWords = append(keyWords, "information_schema")
	sort.Strings(keyWords)
	for _, word := range keyWords {
		tokens = append(tokens, prompt.Suggest{Text: word})
	}
	return
}

func init() {
	flag.StringVar(&endpoint, "endpoint", "http://localhost:9000", "Broker HTTP Endpoint")
}

// printErr prints error message.
func printErr(err error) {
	fmt.Println(color.RedString("ERROR:%s", err))
}

func exit() {
	fmt.Println("Good Bye :)")
	os.Exit(0)
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
		switch strings.ToLower(blocks[0]) {
		case "exit":
			exit()
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
	version := rs.Rows[0][0]
	fmt.Println("Welcome to the LinDB. Commands end with ; .")
	fmt.Printf("Server version: %s\n", version)
	endpointStr := fmt.Sprintf("lin@%s", endpointURL.Host)
	var spaces []string
	for i := 2; i < len(endpointStr); i++ {
		spaces = append(spaces, " ")
	}
	prefix := strings.Join(spaces, "")

	p := newPrompt(
		executor,
		prompt.WithCompleter(func(doc prompt.Document) (suggestions []prompt.Suggest, startChar istrings.RuneNumber, endChar istrings.RuneNumber) {
			endIndex := doc.CurrentRuneIndex()
			w := doc.GetWordBeforeCursor()
			startIndex := endIndex - istrings.RuneCountInString(w)
			if w == "" {
				return nil, startIndex, endIndex
			}
			return prompt.FilterHasPrefix(suggestItems, w, true), startIndex, endIndex
		}),
		prompt.WithPrefixCallback(func() string {
			if live {
				return endpointStr + "> "
			}
			return prefix + "- > "
		}),
		prompt.WithTitle(fmt.Sprintf("LinDB %s Command Line Client", version)),

		prompt.WithPrefixTextColor(prompt.Blue),
		prompt.WithInputTextColor(prompt.White),

		prompt.WithSuggestionBGColor(prompt.LightGray),
		prompt.WithSuggestionTextColor(prompt.Black),
		prompt.WithDescriptionBGColor(prompt.White),
		prompt.WithDescriptionTextColor(prompt.Black),

		prompt.WithSelectedSuggestionBGColor(prompt.DarkBlue),
		prompt.WithSelectedSuggestionTextColor(prompt.Black),
		prompt.WithSelectedDescriptionBGColor(prompt.Blue),
		prompt.WithSelectedDescriptionTextColor(prompt.Black),

		// key bind
		prompt.WithKeyBind(prompt.KeyBind{Key: prompt.ControlC, Fn: func(buf *prompt.Prompt) bool {
			exit()
			return false
		}}),

		// highlight
		prompt.WithLexer(prompt.NewEagerLexer(func(line string) []prompt.Token {
			if len(line) == 0 {
				return nil
			}

			var elements []prompt.Token
			var currentByte istrings.ByteNumber
			var firstByte istrings.ByteNumber
			var firstCharSeen bool
			var lastChar rune

			isKeyWord := func(key string) bool {
				_, ok := lo.Find(suggestItems, func(item prompt.Suggest) bool {
					return strings.ToUpper(key) == strings.ToUpper(item.Text)
				})
				return ok
			}
			var color prompt.Color
			for i, char := range line {
				currentByte = istrings.ByteNumber(i)
				lastChar = char
				if unicode.IsSpace(char) {
					if !firstCharSeen {
						continue
					}

					if isKeyWord(line[firstByte:currentByte]) {
						color = prompt.Purple
						element := prompt.NewSimpleToken(
							firstByte,
							currentByte-1,
							prompt.SimpleTokenWithColor(color),
						)
						elements = append(elements, element)
					}
					firstCharSeen = false
					continue
				}
				if !firstCharSeen {
					firstByte = istrings.ByteNumber(i)
					firstCharSeen = true
				}
			}
			if !unicode.IsSpace(lastChar) {
				start := currentByte + istrings.ByteNumber(utf8.RuneLen(lastChar)) - 1
				if isKeyWord(line[firstByte:start]) {
					color = prompt.Purple
					element := prompt.NewSimpleToken(
						firstByte,
						start,
						prompt.SimpleTokenWithColor(color),
					)
					elements = append(elements, element)
				}
			}

			return elements
		})),
	)
	runPromptFn(p)
}

func runPrompt(p *prompt.Prompt) {
	p.Run()
}
