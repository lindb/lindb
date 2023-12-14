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

package sql

import (
	"errors"
	"sort"
	"strconv"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/models"
	optionpkg "github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

var errIntervalRetentionRequired = errors.New("both interval and retention are required")

// schemasStmtParser represents show schemas statement parser.
type schemasStmtParser struct {
	schema *stmt.Schema
	with   bool // whether to create a table using the WITH clause.
	err    error
}

// newSchemasStmtParse creates a show schemas statement parser.
func newSchemasStmtParse(schemaType stmt.SchemaType) *schemasStmtParser {
	return &schemasStmtParser{schema: &stmt.Schema{Type: schemaType}}
}

// visitName visits when production database config expression is entered.
func (s *schemasStmtParser) visitCfg(ctx *grammar.JsonContext) {
	s.schema.Value = ctx.GetText()
}

// visitWithCfg visits when production database config(with clause) expression is entered.
// the config format like this: https://github.com/lindb/lindb/issues/995#issuecomment-1851136998
func (s *schemasStmtParser) visitWithCfg(ctx *grammar.OptionClauseContext) {
	s.with = true

	var (
		database = &models.Database{
			Name:   ctx.DatabaseName().GetText(),
			Option: &optionpkg.DatabaseOption{},
		}
		pairs     = ctx.OptionPairs().AllOptionPair()
		intervals optionpkg.Intervals
	)

	s.fillDatabase(database, pairs)

	for _, option := range ctx.AllClosedOptionPairs() {
		pairs := option.OptionPairs().AllOptionPair()
		// interval and retention
		if len(pairs) != 2 {
			s.err = errIntervalRetentionRequired
			return
		}

		idx, idx2 := 0, 1
		key1, key2 := pairs[idx].OptionKey(), pairs[idx2].OptionKey()

		switch {
		case key1.T_INTERVAL() != nil && key2.T_RETENTION() != nil:
		case key2.T_INTERVAL() != nil && key1.T_RETENTION() != nil:
			idx, idx2 = idx2, idx
		default:
			return
		}
		var interval, retention timeutil.Interval
		if err := interval.ValueOf(s.parseOptionValue(pairs[idx].OptionValue())); err != nil {
			s.err = err
			return
		}
		if err := retention.ValueOf(s.parseOptionValue(pairs[idx2].OptionValue())); err != nil {
			s.err = err
			return
		}
		intervals = append(intervals, optionpkg.Interval{
			Interval:  interval,
			Retention: retention,
		})
	}

	sort.Sort(intervals)
	if err := intervals.IsValid(); err != nil {
		s.err = err
		return
	}

	database.Option.Intervals = intervals

	s.schema.Value = string(encoding.JSONMarshal(database))
}

// parseOptionValue parse option value
func (s *schemasStmtParser) parseOptionValue(ctx grammar.IOptionValueContext) string {
	text := ctx.GetText()
	if ctx.STRING() != nil {
		text = text[1 : len(text)-1]
	}
	return text
}

// parseNumber parse text to integer
func (s *schemasStmtParser) parseNumber(text string) int {
	n, err := strconv.Atoi(text)
	if err != nil {
		s.err = err
	}
	return n
}

// parseBool parse text to bool
func (s *schemasStmtParser) parseBool(text string) bool {
	n, err := strconv.ParseBool(text)
	if err != nil {
		s.err = err
	}
	return n
}

// fillDatabase initializes some properties of database
func (s *schemasStmtParser) fillDatabase(database *models.Database, pairs []grammar.IOptionPairContext) {
	for _, pair := range pairs {
		key, val := pair.OptionKey(), s.parseOptionValue(pair.OptionValue())
		switch {
		case key.T_STORAGE() != nil:
			database.Storage = val
		case key.T_NUM_OF_SHARD() != nil:
			database.NumOfShard = s.parseNumber(val)
		case key.T_REPLICA_FACTOR() != nil:
			database.ReplicaFactor = s.parseNumber(val)
		case key.T_AUTO_CREATE_NS() != nil:
			database.Option.AutoCreateNS = s.parseBool(val)
		case key.T_AHEAD() != nil:
			database.Option.Ahead = val
		case key.T_BEHEAD() != nil:
			database.Option.Behind = val
		}
	}
}

// visitDropDatabase visits when production database name expression is entered.
func (s *schemasStmtParser) visitDatabaseName(ctx *grammar.DatabaseNameContext) {
	if !s.with {
		s.schema.Value = strutil.GetStringValue(ctx.GetText())
	}
}

// build returns the state statement.
func (s *schemasStmtParser) build() (stmt.Statement, error) {
	return s.schema, s.err
}
