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

// Code generated from /Users/huangjie/go/src/github.com/lindb/lindb/sql/grammar/SQL.g4 by ANTLR 4.10.1. DO NOT EDIT.

package grammar // SQL
import (
	"fmt"
	"strconv"
  "sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}


type SQLParser struct {
	*antlr.BaseParser
}

var sqlParserStaticData struct {
  once                   sync.Once
  serializedATN          []int32
  literalNames           []string
  symbolicNames          []string
  ruleNames              []string
  predictionContextCache *antlr.PredictionContextCache
  atn                    *antlr.ATN
  decisionToDFA          []*antlr.DFA
}

func sqlParserInit() {
  staticData := &sqlParserStaticData
  staticData.literalNames = []string{
    "", "'true'", "'false'", "'null'", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 
    "", "", "", "", "", "", "", "", "", "", "'m'", "", "", "", "'M'", "", 
    "'.'", "':'", "'='", "'<>'", "'!='", "'>'", "'>='", "'<'", "'<='", "'=~'", 
    "'!~'", "','", "'{'", "'}'", "'['", "']'", "'('", "')'", "'+'", "'-'", 
    "'/'", "'*'", "'%'", "'_'",
  }
  staticData.symbolicNames = []string{
    "", "", "", "", "STRING", "WS", "T_CREATE", "T_UPDATE", "T_SET", "T_DROP", 
    "T_INTERVAL", "T_INTERVAL_NAME", "T_SHARD", "T_REPLICATION", "T_TTL", 
    "T_META_TTL", "T_PAST_TTL", "T_FUTURE_TTL", "T_KILL", "T_ON", "T_SHOW", 
    "T_USE", "T_STATE_REPO", "T_STATE_MACHINE", "T_MASTER", "T_METADATA", 
    "T_TYPES", "T_TYPE", "T_STORAGES", "T_STORAGE", "T_BROKER", "T_ALIVE", 
    "T_SCHEMAS", "T_DATASBAE", "T_DATASBAES", "T_NAMESPACE", "T_NAMESPACES", 
    "T_NODE", "T_METRICS", "T_METRIC", "T_FIELD", "T_FIELDS", "T_TAG", "T_INFO", 
    "T_KEYS", "T_KEY", "T_WITH", "T_VALUES", "T_VALUE", "T_FROM", "T_WHERE", 
    "T_LIMIT", "T_QUERIES", "T_QUERY", "T_EXPLAIN", "T_WITH_VALUE", "T_SELECT", 
    "T_AS", "T_AND", "T_OR", "T_FILL", "T_NULL", "T_PREVIOUS", "T_ORDER", 
    "T_ASC", "T_DESC", "T_LIKE", "T_NOT", "T_BETWEEN", "T_IS", "T_GROUP", 
    "T_HAVING", "T_BY", "T_FOR", "T_STATS", "T_TIME", "T_NOW", "T_IN", "T_LOG", 
    "T_PROFILE", "T_SUM", "T_MIN", "T_MAX", "T_COUNT", "T_LAST", "T_FIRST", 
    "T_AVG", "T_STDDEV", "T_QUANTILE", "T_RATE", "T_SECOND", "T_MINUTE", 
    "T_HOUR", "T_DAY", "T_WEEK", "T_MONTH", "T_YEAR", "T_DOT", "T_COLON", 
    "T_EQUAL", "T_NOTEQUAL", "T_NOTEQUAL2", "T_GREATER", "T_GREATEREQUAL", 
    "T_LESS", "T_LESSEQUAL", "T_REGEXP", "T_NEQREGEXP", "T_COMMA", "T_OPEN_B", 
    "T_CLOSE_B", "T_OPEN_SB", "T_CLOSE_SB", "T_OPEN_P", "T_CLOSE_P", "T_ADD", 
    "T_SUB", "T_DIV", "T_MUL", "T_MOD", "T_UNDERLINE", "L_ID", "L_INT", 
    "L_DEC",
  }
  staticData.ruleNames = []string{
    "statement", "useStmt", "showStmt", "showMasterStmt", "showStoragesStmt", 
    "showMetadataTypesStmt", "showBrokerMetaStmt", "showMasterMetaStmt", 
    "showStorageMetaStmt", "showAliveStmt", "showReplicationStmt", "showBrokerMetricStmt", 
    "showStorageMetricStmt", "createStorageStmt", "showSchemasStmt", "createDatabaseStmt", 
    "dropDatabaseStmt", "showDatabaseStmt", "showNameSpacesStmt", "showMetricsStmt", 
    "showFieldsStmt", "showTagKeysStmt", "showTagValuesStmt", "prefix", 
    "withTagKey", "namespace", "databaseName", "source", "queryStmt", "sourceAndSelect", 
    "selectExpr", "fields", "field", "alias", "storageFilter", "databaseFilter", 
    "typeFilter", "fromClause", "whereClause", "conditionExpr", "tagFilterExpr", 
    "tagValueList", "metricListFilter", "metricList", "timeRangeExpr", "timeExpr", 
    "nowExpr", "nowFunc", "groupByClause", "groupByKeys", "groupByKey", 
    "fillOption", "orderByClause", "sortField", "sortFields", "havingClause", 
    "boolExpr", "boolExprLogicalOp", "boolExprAtom", "binaryExpr", "binaryOperator", 
    "fieldExpr", "durationLit", "intervalItem", "exprFunc", "funcName", 
    "exprFuncParams", "funcParam", "exprAtom", "identFilter", "json", "obj", 
    "pair", "arr", "value", "intNumber", "decNumber", "limitClause", "metricName", 
    "tagKey", "tagValue", "ident", "nonReservedWords",
  }
  staticData.predictionContextCache = antlr.NewPredictionContextCache()
  staticData.serializedATN = []int32{
	4, 1, 123, 747, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 
	4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7, 
	10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15, 
	2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2, 
	21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26, 
	7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7, 
	31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36, 
	2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2, 
	42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47, 
	7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7, 
	52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57, 7, 57, 
	2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7, 62, 2, 
	63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67, 2, 68, 
	7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7, 71, 2, 72, 7, 72, 2, 73, 7, 
	73, 2, 74, 7, 74, 2, 75, 7, 75, 2, 76, 7, 76, 2, 77, 7, 77, 2, 78, 7, 78, 
	2, 79, 7, 79, 2, 80, 7, 80, 2, 81, 7, 81, 2, 82, 7, 82, 1, 0, 1, 0, 1, 
	0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 3, 0, 176, 8, 0, 1, 1, 1, 1, 1, 
	1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 
	2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 198, 8, 2, 1, 3, 1, 3, 1, 3, 1, 
	4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 
	6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 
	8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3, 8, 234, 8, 8, 1, 8, 1, 8, 1, 
	8, 3, 8, 239, 8, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 10, 
	1, 10, 3, 10, 250, 8, 10, 1, 10, 1, 10, 1, 10, 3, 10, 255, 8, 10, 1, 11, 
	1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 
	12, 3, 12, 269, 8, 12, 1, 12, 1, 12, 1, 12, 3, 12, 274, 8, 12, 1, 13, 1, 
	13, 1, 13, 1, 13, 1, 14, 1, 14, 1, 14, 1, 15, 1, 15, 1, 15, 1, 15, 1, 16, 
	1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1, 
	18, 1, 18, 3, 18, 300, 8, 18, 1, 18, 3, 18, 303, 8, 18, 1, 19, 1, 19, 1, 
	19, 1, 19, 3, 19, 309, 8, 19, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 315, 8, 
	19, 1, 19, 3, 19, 318, 8, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 
	1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1, 
	22, 1, 22, 3, 22, 338, 8, 22, 1, 22, 3, 22, 341, 8, 22, 1, 23, 1, 23, 1, 
	24, 1, 24, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1, 28, 3, 28, 354, 
	8, 28, 1, 28, 1, 28, 3, 28, 358, 8, 28, 1, 28, 3, 28, 361, 8, 28, 1, 28, 
	3, 28, 364, 8, 28, 1, 28, 3, 28, 367, 8, 28, 1, 28, 3, 28, 370, 8, 28, 
	1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 1, 29, 3, 29, 378, 8, 29, 1, 30, 1, 
	30, 1, 30, 1, 31, 1, 31, 1, 31, 5, 31, 386, 8, 31, 10, 31, 12, 31, 389, 
	9, 31, 1, 32, 1, 32, 3, 32, 393, 8, 32, 1, 33, 1, 33, 1, 33, 1, 34, 1, 
	34, 1, 34, 1, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36, 1, 36, 1, 36, 1, 36, 
	1, 37, 1, 37, 1, 37, 1, 37, 3, 37, 414, 8, 37, 1, 38, 1, 38, 1, 38, 1, 
	39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 3, 39, 427, 8, 39, 
	3, 39, 429, 8, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 
	40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 3, 40, 445, 8, 40, 1, 40, 
	1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 3, 40, 453, 8, 40, 1, 40, 1, 40, 1, 
	40, 1, 40, 3, 40, 459, 8, 40, 1, 40, 1, 40, 1, 40, 5, 40, 464, 8, 40, 10, 
	40, 12, 40, 467, 9, 40, 1, 41, 1, 41, 1, 41, 5, 41, 472, 8, 41, 10, 41, 
	12, 41, 475, 9, 41, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 
	43, 1, 43, 5, 43, 486, 8, 43, 10, 43, 12, 43, 489, 9, 43, 1, 44, 1, 44, 
	1, 44, 3, 44, 494, 8, 44, 1, 45, 1, 45, 1, 45, 1, 45, 3, 45, 500, 8, 45, 
	1, 46, 1, 46, 3, 46, 504, 8, 46, 1, 47, 1, 47, 1, 47, 3, 47, 509, 8, 47, 
	1, 47, 1, 47, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 3, 
	48, 521, 8, 48, 1, 48, 3, 48, 524, 8, 48, 1, 49, 1, 49, 1, 49, 5, 49, 529, 
	8, 49, 10, 49, 12, 49, 532, 9, 49, 1, 50, 1, 50, 1, 50, 1, 50, 1, 50, 1, 
	50, 3, 50, 540, 8, 50, 1, 51, 1, 51, 1, 52, 1, 52, 1, 52, 1, 52, 1, 53, 
	1, 53, 5, 53, 550, 8, 53, 10, 53, 12, 53, 553, 9, 53, 1, 54, 1, 54, 1, 
	54, 5, 54, 558, 8, 54, 10, 54, 12, 54, 561, 9, 54, 1, 55, 1, 55, 1, 55, 
	1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 572, 8, 56, 1, 56, 1, 
	56, 1, 56, 1, 56, 5, 56, 578, 8, 56, 10, 56, 12, 56, 581, 9, 56, 1, 57, 
	1, 57, 1, 58, 1, 58, 1, 59, 1, 59, 1, 59, 1, 59, 1, 60, 1, 60, 1, 60, 1, 
	60, 1, 60, 1, 60, 1, 60, 1, 60, 3, 60, 599, 8, 60, 1, 61, 1, 61, 1, 61, 
	1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 3, 61, 609, 8, 61, 1, 61, 1, 61, 1, 
	61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 5, 61, 
	623, 8, 61, 10, 61, 12, 61, 626, 9, 61, 1, 62, 1, 62, 1, 62, 1, 63, 1, 
	63, 1, 64, 1, 64, 1, 64, 3, 64, 636, 8, 64, 1, 64, 1, 64, 1, 65, 1, 65, 
	1, 66, 1, 66, 1, 66, 5, 66, 645, 8, 66, 10, 66, 12, 66, 648, 9, 66, 1, 
	67, 1, 67, 3, 67, 652, 8, 67, 1, 68, 1, 68, 3, 68, 656, 8, 68, 1, 68, 1, 
	68, 3, 68, 660, 8, 68, 1, 69, 1, 69, 1, 69, 1, 69, 1, 70, 1, 70, 1, 71, 
	1, 71, 1, 71, 1, 71, 5, 71, 672, 8, 71, 10, 71, 12, 71, 675, 9, 71, 1, 
	71, 1, 71, 1, 71, 1, 71, 3, 71, 681, 8, 71, 1, 72, 1, 72, 1, 72, 1, 72, 
	1, 73, 1, 73, 1, 73, 1, 73, 5, 73, 691, 8, 73, 10, 73, 12, 73, 694, 9, 
	73, 1, 73, 1, 73, 1, 73, 1, 73, 3, 73, 700, 8, 73, 1, 74, 1, 74, 1, 74, 
	1, 74, 1, 74, 1, 74, 1, 74, 1, 74, 3, 74, 710, 8, 74, 1, 75, 3, 75, 713, 
	8, 75, 1, 75, 1, 75, 1, 76, 3, 76, 718, 8, 76, 1, 76, 1, 76, 1, 77, 1, 
	77, 1, 77, 1, 78, 1, 78, 1, 79, 1, 79, 1, 80, 1, 80, 1, 81, 1, 81, 3, 81, 
	733, 8, 81, 1, 81, 1, 81, 1, 81, 3, 81, 738, 8, 81, 5, 81, 740, 8, 81, 
	10, 81, 12, 81, 743, 9, 81, 1, 82, 1, 82, 1, 82, 0, 3, 80, 112, 122, 83, 
	0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 
	38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 
	74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104, 106, 
	108, 110, 112, 114, 116, 118, 120, 122, 124, 126, 128, 130, 132, 134, 136, 
	138, 140, 142, 144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 164, 0, 
	10, 1, 0, 29, 30, 1, 0, 22, 23, 1, 0, 58, 59, 2, 0, 61, 62, 122, 123, 1, 
	0, 64, 65, 2, 0, 66, 66, 106, 106, 1, 0, 90, 96, 1, 0, 80, 89, 1, 0, 115, 
	116, 1, 0, 6, 96, 771, 0, 175, 1, 0, 0, 0, 2, 177, 1, 0, 0, 0, 4, 197, 
	1, 0, 0, 0, 6, 199, 1, 0, 0, 0, 8, 202, 1, 0, 0, 0, 10, 205, 1, 0, 0, 0, 
	12, 209, 1, 0, 0, 0, 14, 217, 1, 0, 0, 0, 16, 225, 1, 0, 0, 0, 18, 240, 
	1, 0, 0, 0, 20, 244, 1, 0, 0, 0, 22, 256, 1, 0, 0, 0, 24, 262, 1, 0, 0, 
	0, 26, 275, 1, 0, 0, 0, 28, 279, 1, 0, 0, 0, 30, 282, 1, 0, 0, 0, 32, 286, 
	1, 0, 0, 0, 34, 290, 1, 0, 0, 0, 36, 293, 1, 0, 0, 0, 38, 304, 1, 0, 0, 
	0, 40, 319, 1, 0, 0, 0, 42, 323, 1, 0, 0, 0, 44, 328, 1, 0, 0, 0, 46, 342, 
	1, 0, 0, 0, 48, 344, 1, 0, 0, 0, 50, 346, 1, 0, 0, 0, 52, 348, 1, 0, 0, 
	0, 54, 350, 1, 0, 0, 0, 56, 353, 1, 0, 0, 0, 58, 377, 1, 0, 0, 0, 60, 379, 
	1, 0, 0, 0, 62, 382, 1, 0, 0, 0, 64, 390, 1, 0, 0, 0, 66, 394, 1, 0, 0, 
	0, 68, 397, 1, 0, 0, 0, 70, 401, 1, 0, 0, 0, 72, 405, 1, 0, 0, 0, 74, 409, 
	1, 0, 0, 0, 76, 415, 1, 0, 0, 0, 78, 428, 1, 0, 0, 0, 80, 458, 1, 0, 0, 
	0, 82, 468, 1, 0, 0, 0, 84, 476, 1, 0, 0, 0, 86, 482, 1, 0, 0, 0, 88, 490, 
	1, 0, 0, 0, 90, 495, 1, 0, 0, 0, 92, 501, 1, 0, 0, 0, 94, 505, 1, 0, 0, 
	0, 96, 512, 1, 0, 0, 0, 98, 525, 1, 0, 0, 0, 100, 539, 1, 0, 0, 0, 102, 
	541, 1, 0, 0, 0, 104, 543, 1, 0, 0, 0, 106, 547, 1, 0, 0, 0, 108, 554, 
	1, 0, 0, 0, 110, 562, 1, 0, 0, 0, 112, 571, 1, 0, 0, 0, 114, 582, 1, 0, 
	0, 0, 116, 584, 1, 0, 0, 0, 118, 586, 1, 0, 0, 0, 120, 598, 1, 0, 0, 0, 
	122, 608, 1, 0, 0, 0, 124, 627, 1, 0, 0, 0, 126, 630, 1, 0, 0, 0, 128, 
	632, 1, 0, 0, 0, 130, 639, 1, 0, 0, 0, 132, 641, 1, 0, 0, 0, 134, 651, 
	1, 0, 0, 0, 136, 659, 1, 0, 0, 0, 138, 661, 1, 0, 0, 0, 140, 665, 1, 0, 
	0, 0, 142, 680, 1, 0, 0, 0, 144, 682, 1, 0, 0, 0, 146, 699, 1, 0, 0, 0, 
	148, 709, 1, 0, 0, 0, 150, 712, 1, 0, 0, 0, 152, 717, 1, 0, 0, 0, 154, 
	721, 1, 0, 0, 0, 156, 724, 1, 0, 0, 0, 158, 726, 1, 0, 0, 0, 160, 728, 
	1, 0, 0, 0, 162, 732, 1, 0, 0, 0, 164, 744, 1, 0, 0, 0, 166, 176, 3, 4, 
	2, 0, 167, 176, 3, 26, 13, 0, 168, 176, 3, 2, 1, 0, 169, 176, 3, 56, 28, 
	0, 170, 176, 3, 30, 15, 0, 171, 176, 3, 32, 16, 0, 172, 173, 3, 162, 81, 
	0, 173, 174, 5, 0, 0, 1, 174, 176, 1, 0, 0, 0, 175, 166, 1, 0, 0, 0, 175, 
	167, 1, 0, 0, 0, 175, 168, 1, 0, 0, 0, 175, 169, 1, 0, 0, 0, 175, 170, 
	1, 0, 0, 0, 175, 171, 1, 0, 0, 0, 175, 172, 1, 0, 0, 0, 176, 1, 1, 0, 0, 
	0, 177, 178, 5, 21, 0, 0, 178, 179, 3, 162, 81, 0, 179, 3, 1, 0, 0, 0, 
	180, 198, 3, 6, 3, 0, 181, 198, 3, 10, 5, 0, 182, 198, 3, 12, 6, 0, 183, 
	198, 3, 14, 7, 0, 184, 198, 3, 16, 8, 0, 185, 198, 3, 8, 4, 0, 186, 198, 
	3, 18, 9, 0, 187, 198, 3, 22, 11, 0, 188, 198, 3, 24, 12, 0, 189, 198, 
	3, 20, 10, 0, 190, 198, 3, 28, 14, 0, 191, 198, 3, 34, 17, 0, 192, 198, 
	3, 36, 18, 0, 193, 198, 3, 38, 19, 0, 194, 198, 3, 40, 20, 0, 195, 198, 
	3, 42, 21, 0, 196, 198, 3, 44, 22, 0, 197, 180, 1, 0, 0, 0, 197, 181, 1, 
	0, 0, 0, 197, 182, 1, 0, 0, 0, 197, 183, 1, 0, 0, 0, 197, 184, 1, 0, 0, 
	0, 197, 185, 1, 0, 0, 0, 197, 186, 1, 0, 0, 0, 197, 187, 1, 0, 0, 0, 197, 
	188, 1, 0, 0, 0, 197, 189, 1, 0, 0, 0, 197, 190, 1, 0, 0, 0, 197, 191, 
	1, 0, 0, 0, 197, 192, 1, 0, 0, 0, 197, 193, 1, 0, 0, 0, 197, 194, 1, 0, 
	0, 0, 197, 195, 1, 0, 0, 0, 197, 196, 1, 0, 0, 0, 198, 5, 1, 0, 0, 0, 199, 
	200, 5, 20, 0, 0, 200, 201, 5, 24, 0, 0, 201, 7, 1, 0, 0, 0, 202, 203, 
	5, 20, 0, 0, 203, 204, 5, 28, 0, 0, 204, 9, 1, 0, 0, 0, 205, 206, 5, 20, 
	0, 0, 206, 207, 5, 25, 0, 0, 207, 208, 5, 26, 0, 0, 208, 11, 1, 0, 0, 0, 
	209, 210, 5, 20, 0, 0, 210, 211, 5, 30, 0, 0, 211, 212, 5, 25, 0, 0, 212, 
	213, 5, 49, 0, 0, 213, 214, 3, 54, 27, 0, 214, 215, 5, 50, 0, 0, 215, 216, 
	3, 72, 36, 0, 216, 13, 1, 0, 0, 0, 217, 218, 5, 20, 0, 0, 218, 219, 5, 
	24, 0, 0, 219, 220, 5, 25, 0, 0, 220, 221, 5, 49, 0, 0, 221, 222, 3, 54, 
	27, 0, 222, 223, 5, 50, 0, 0, 223, 224, 3, 72, 36, 0, 224, 15, 1, 0, 0, 
	0, 225, 226, 5, 20, 0, 0, 226, 227, 5, 29, 0, 0, 227, 228, 5, 25, 0, 0, 
	228, 229, 5, 49, 0, 0, 229, 230, 3, 54, 27, 0, 230, 233, 5, 50, 0, 0, 231, 
	234, 3, 68, 34, 0, 232, 234, 3, 72, 36, 0, 233, 231, 1, 0, 0, 0, 233, 232, 
	1, 0, 0, 0, 234, 235, 1, 0, 0, 0, 235, 238, 5, 58, 0, 0, 236, 239, 3, 68, 
	34, 0, 237, 239, 3, 72, 36, 0, 238, 236, 1, 0, 0, 0, 238, 237, 1, 0, 0, 
	0, 239, 17, 1, 0, 0, 0, 240, 241, 5, 20, 0, 0, 241, 242, 7, 0, 0, 0, 242, 
	243, 5, 31, 0, 0, 243, 19, 1, 0, 0, 0, 244, 245, 5, 20, 0, 0, 245, 246, 
	5, 13, 0, 0, 246, 249, 5, 50, 0, 0, 247, 250, 3, 68, 34, 0, 248, 250, 3, 
	70, 35, 0, 249, 247, 1, 0, 0, 0, 249, 248, 1, 0, 0, 0, 250, 251, 1, 0, 
	0, 0, 251, 254, 5, 58, 0, 0, 252, 255, 3, 68, 34, 0, 253, 255, 3, 70, 35, 
	0, 254, 252, 1, 0, 0, 0, 254, 253, 1, 0, 0, 0, 255, 21, 1, 0, 0, 0, 256, 
	257, 5, 20, 0, 0, 257, 258, 5, 30, 0, 0, 258, 259, 5, 39, 0, 0, 259, 260, 
	5, 50, 0, 0, 260, 261, 3, 84, 42, 0, 261, 23, 1, 0, 0, 0, 262, 263, 5, 
	20, 0, 0, 263, 264, 5, 29, 0, 0, 264, 265, 5, 39, 0, 0, 265, 268, 5, 50, 
	0, 0, 266, 269, 3, 68, 34, 0, 267, 269, 3, 84, 42, 0, 268, 266, 1, 0, 0, 
	0, 268, 267, 1, 0, 0, 0, 269, 270, 1, 0, 0, 0, 270, 273, 5, 58, 0, 0, 271, 
	274, 3, 68, 34, 0, 272, 274, 3, 84, 42, 0, 273, 271, 1, 0, 0, 0, 273, 272, 
	1, 0, 0, 0, 274, 25, 1, 0, 0, 0, 275, 276, 5, 6, 0, 0, 276, 277, 5, 29, 
	0, 0, 277, 278, 3, 140, 70, 0, 278, 27, 1, 0, 0, 0, 279, 280, 5, 20, 0, 
	0, 280, 281, 5, 32, 0, 0, 281, 29, 1, 0, 0, 0, 282, 283, 5, 6, 0, 0, 283, 
	284, 5, 33, 0, 0, 284, 285, 3, 140, 70, 0, 285, 31, 1, 0, 0, 0, 286, 287, 
	5, 9, 0, 0, 287, 288, 5, 33, 0, 0, 288, 289, 3, 52, 26, 0, 289, 33, 1, 
	0, 0, 0, 290, 291, 5, 20, 0, 0, 291, 292, 5, 34, 0, 0, 292, 35, 1, 0, 0, 
	0, 293, 294, 5, 20, 0, 0, 294, 299, 5, 36, 0, 0, 295, 296, 5, 50, 0, 0, 
	296, 297, 5, 35, 0, 0, 297, 298, 5, 99, 0, 0, 298, 300, 3, 46, 23, 0, 299, 
	295, 1, 0, 0, 0, 299, 300, 1, 0, 0, 0, 300, 302, 1, 0, 0, 0, 301, 303, 
	3, 154, 77, 0, 302, 301, 1, 0, 0, 0, 302, 303, 1, 0, 0, 0, 303, 37, 1, 
	0, 0, 0, 304, 305, 5, 20, 0, 0, 305, 308, 5, 38, 0, 0, 306, 307, 5, 19, 
	0, 0, 307, 309, 3, 50, 25, 0, 308, 306, 1, 0, 0, 0, 308, 309, 1, 0, 0, 
	0, 309, 314, 1, 0, 0, 0, 310, 311, 5, 50, 0, 0, 311, 312, 5, 39, 0, 0, 
	312, 313, 5, 99, 0, 0, 313, 315, 3, 46, 23, 0, 314, 310, 1, 0, 0, 0, 314, 
	315, 1, 0, 0, 0, 315, 317, 1, 0, 0, 0, 316, 318, 3, 154, 77, 0, 317, 316, 
	1, 0, 0, 0, 317, 318, 1, 0, 0, 0, 318, 39, 1, 0, 0, 0, 319, 320, 5, 20, 
	0, 0, 320, 321, 5, 41, 0, 0, 321, 322, 3, 74, 37, 0, 322, 41, 1, 0, 0, 
	0, 323, 324, 5, 20, 0, 0, 324, 325, 5, 42, 0, 0, 325, 326, 5, 44, 0, 0, 
	326, 327, 3, 74, 37, 0, 327, 43, 1, 0, 0, 0, 328, 329, 5, 20, 0, 0, 329, 
	330, 5, 42, 0, 0, 330, 331, 5, 47, 0, 0, 331, 332, 3, 74, 37, 0, 332, 333, 
	5, 46, 0, 0, 333, 334, 5, 45, 0, 0, 334, 335, 5, 99, 0, 0, 335, 337, 3, 
	48, 24, 0, 336, 338, 3, 76, 38, 0, 337, 336, 1, 0, 0, 0, 337, 338, 1, 0, 
	0, 0, 338, 340, 1, 0, 0, 0, 339, 341, 3, 154, 77, 0, 340, 339, 1, 0, 0, 
	0, 340, 341, 1, 0, 0, 0, 341, 45, 1, 0, 0, 0, 342, 343, 3, 162, 81, 0, 
	343, 47, 1, 0, 0, 0, 344, 345, 3, 162, 81, 0, 345, 49, 1, 0, 0, 0, 346, 
	347, 3, 162, 81, 0, 347, 51, 1, 0, 0, 0, 348, 349, 3, 162, 81, 0, 349, 
	53, 1, 0, 0, 0, 350, 351, 7, 1, 0, 0, 351, 55, 1, 0, 0, 0, 352, 354, 5, 
	54, 0, 0, 353, 352, 1, 0, 0, 0, 353, 354, 1, 0, 0, 0, 354, 355, 1, 0, 0, 
	0, 355, 357, 3, 58, 29, 0, 356, 358, 3, 76, 38, 0, 357, 356, 1, 0, 0, 0, 
	357, 358, 1, 0, 0, 0, 358, 360, 1, 0, 0, 0, 359, 361, 3, 96, 48, 0, 360, 
	359, 1, 0, 0, 0, 360, 361, 1, 0, 0, 0, 361, 363, 1, 0, 0, 0, 362, 364, 
	3, 104, 52, 0, 363, 362, 1, 0, 0, 0, 363, 364, 1, 0, 0, 0, 364, 366, 1, 
	0, 0, 0, 365, 367, 3, 154, 77, 0, 366, 365, 1, 0, 0, 0, 366, 367, 1, 0, 
	0, 0, 367, 369, 1, 0, 0, 0, 368, 370, 5, 55, 0, 0, 369, 368, 1, 0, 0, 0, 
	369, 370, 1, 0, 0, 0, 370, 57, 1, 0, 0, 0, 371, 372, 3, 60, 30, 0, 372, 
	373, 3, 74, 37, 0, 373, 378, 1, 0, 0, 0, 374, 375, 3, 74, 37, 0, 375, 376, 
	3, 60, 30, 0, 376, 378, 1, 0, 0, 0, 377, 371, 1, 0, 0, 0, 377, 374, 1, 
	0, 0, 0, 378, 59, 1, 0, 0, 0, 379, 380, 5, 56, 0, 0, 380, 381, 3, 62, 31, 
	0, 381, 61, 1, 0, 0, 0, 382, 387, 3, 64, 32, 0, 383, 384, 5, 108, 0, 0, 
	384, 386, 3, 64, 32, 0, 385, 383, 1, 0, 0, 0, 386, 389, 1, 0, 0, 0, 387, 
	385, 1, 0, 0, 0, 387, 388, 1, 0, 0, 0, 388, 63, 1, 0, 0, 0, 389, 387, 1, 
	0, 0, 0, 390, 392, 3, 122, 61, 0, 391, 393, 3, 66, 33, 0, 392, 391, 1, 
	0, 0, 0, 392, 393, 1, 0, 0, 0, 393, 65, 1, 0, 0, 0, 394, 395, 5, 57, 0, 
	0, 395, 396, 3, 162, 81, 0, 396, 67, 1, 0, 0, 0, 397, 398, 5, 29, 0, 0, 
	398, 399, 5, 99, 0, 0, 399, 400, 3, 162, 81, 0, 400, 69, 1, 0, 0, 0, 401, 
	402, 5, 33, 0, 0, 402, 403, 5, 99, 0, 0, 403, 404, 3, 162, 81, 0, 404, 
	71, 1, 0, 0, 0, 405, 406, 5, 27, 0, 0, 406, 407, 5, 99, 0, 0, 407, 408, 
	3, 162, 81, 0, 408, 73, 1, 0, 0, 0, 409, 410, 5, 49, 0, 0, 410, 413, 3, 
	156, 78, 0, 411, 412, 5, 19, 0, 0, 412, 414, 3, 50, 25, 0, 413, 411, 1, 
	0, 0, 0, 413, 414, 1, 0, 0, 0, 414, 75, 1, 0, 0, 0, 415, 416, 5, 50, 0, 
	0, 416, 417, 3, 78, 39, 0, 417, 77, 1, 0, 0, 0, 418, 429, 3, 80, 40, 0, 
	419, 420, 3, 80, 40, 0, 420, 421, 5, 58, 0, 0, 421, 422, 3, 88, 44, 0, 
	422, 429, 1, 0, 0, 0, 423, 426, 3, 88, 44, 0, 424, 425, 5, 58, 0, 0, 425, 
	427, 3, 80, 40, 0, 426, 424, 1, 0, 0, 0, 426, 427, 1, 0, 0, 0, 427, 429, 
	1, 0, 0, 0, 428, 418, 1, 0, 0, 0, 428, 419, 1, 0, 0, 0, 428, 423, 1, 0, 
	0, 0, 429, 79, 1, 0, 0, 0, 430, 431, 6, 40, -1, 0, 431, 432, 5, 113, 0, 
	0, 432, 433, 3, 80, 40, 0, 433, 434, 5, 114, 0, 0, 434, 459, 1, 0, 0, 0, 
	435, 444, 3, 158, 79, 0, 436, 445, 5, 99, 0, 0, 437, 445, 5, 66, 0, 0, 
	438, 439, 5, 67, 0, 0, 439, 445, 5, 66, 0, 0, 440, 445, 5, 106, 0, 0, 441, 
	445, 5, 107, 0, 0, 442, 445, 5, 100, 0, 0, 443, 445, 5, 101, 0, 0, 444, 
	436, 1, 0, 0, 0, 444, 437, 1, 0, 0, 0, 444, 438, 1, 0, 0, 0, 444, 440, 
	1, 0, 0, 0, 444, 441, 1, 0, 0, 0, 444, 442, 1, 0, 0, 0, 444, 443, 1, 0, 
	0, 0, 445, 446, 1, 0, 0, 0, 446, 447, 3, 160, 80, 0, 447, 459, 1, 0, 0, 
	0, 448, 452, 3, 158, 79, 0, 449, 453, 5, 77, 0, 0, 450, 451, 5, 67, 0, 
	0, 451, 453, 5, 77, 0, 0, 452, 449, 1, 0, 0, 0, 452, 450, 1, 0, 0, 0, 453, 
	454, 1, 0, 0, 0, 454, 455, 5, 113, 0, 0, 455, 456, 3, 82, 41, 0, 456, 457, 
	5, 114, 0, 0, 457, 459, 1, 0, 0, 0, 458, 430, 1, 0, 0, 0, 458, 435, 1, 
	0, 0, 0, 458, 448, 1, 0, 0, 0, 459, 465, 1, 0, 0, 0, 460, 461, 10, 1, 0, 
	0, 461, 462, 7, 2, 0, 0, 462, 464, 3, 80, 40, 2, 463, 460, 1, 0, 0, 0, 
	464, 467, 1, 0, 0, 0, 465, 463, 1, 0, 0, 0, 465, 466, 1, 0, 0, 0, 466, 
	81, 1, 0, 0, 0, 467, 465, 1, 0, 0, 0, 468, 473, 3, 160, 80, 0, 469, 470, 
	5, 108, 0, 0, 470, 472, 3, 160, 80, 0, 471, 469, 1, 0, 0, 0, 472, 475, 
	1, 0, 0, 0, 473, 471, 1, 0, 0, 0, 473, 474, 1, 0, 0, 0, 474, 83, 1, 0, 
	0, 0, 475, 473, 1, 0, 0, 0, 476, 477, 5, 39, 0, 0, 477, 478, 5, 77, 0, 
	0, 478, 479, 5, 113, 0, 0, 479, 480, 3, 86, 43, 0, 480, 481, 5, 114, 0, 
	0, 481, 85, 1, 0, 0, 0, 482, 487, 3, 162, 81, 0, 483, 484, 5, 108, 0, 0, 
	484, 486, 3, 162, 81, 0, 485, 483, 1, 0, 0, 0, 486, 489, 1, 0, 0, 0, 487, 
	485, 1, 0, 0, 0, 487, 488, 1, 0, 0, 0, 488, 87, 1, 0, 0, 0, 489, 487, 1, 
	0, 0, 0, 490, 493, 3, 90, 45, 0, 491, 492, 5, 58, 0, 0, 492, 494, 3, 90, 
	45, 0, 493, 491, 1, 0, 0, 0, 493, 494, 1, 0, 0, 0, 494, 89, 1, 0, 0, 0, 
	495, 496, 5, 75, 0, 0, 496, 499, 3, 120, 60, 0, 497, 500, 3, 92, 46, 0, 
	498, 500, 3, 162, 81, 0, 499, 497, 1, 0, 0, 0, 499, 498, 1, 0, 0, 0, 500, 
	91, 1, 0, 0, 0, 501, 503, 3, 94, 47, 0, 502, 504, 3, 124, 62, 0, 503, 502, 
	1, 0, 0, 0, 503, 504, 1, 0, 0, 0, 504, 93, 1, 0, 0, 0, 505, 506, 5, 76, 
	0, 0, 506, 508, 5, 113, 0, 0, 507, 509, 3, 132, 66, 0, 508, 507, 1, 0, 
	0, 0, 508, 509, 1, 0, 0, 0, 509, 510, 1, 0, 0, 0, 510, 511, 5, 114, 0, 
	0, 511, 95, 1, 0, 0, 0, 512, 513, 5, 70, 0, 0, 513, 514, 5, 72, 0, 0, 514, 
	520, 3, 98, 49, 0, 515, 516, 5, 60, 0, 0, 516, 517, 5, 113, 0, 0, 517, 
	518, 3, 102, 51, 0, 518, 519, 5, 114, 0, 0, 519, 521, 1, 0, 0, 0, 520, 
	515, 1, 0, 0, 0, 520, 521, 1, 0, 0, 0, 521, 523, 1, 0, 0, 0, 522, 524, 
	3, 110, 55, 0, 523, 522, 1, 0, 0, 0, 523, 524, 1, 0, 0, 0, 524, 97, 1, 
	0, 0, 0, 525, 530, 3, 100, 50, 0, 526, 527, 5, 108, 0, 0, 527, 529, 3, 
	100, 50, 0, 528, 526, 1, 0, 0, 0, 529, 532, 1, 0, 0, 0, 530, 528, 1, 0, 
	0, 0, 530, 531, 1, 0, 0, 0, 531, 99, 1, 0, 0, 0, 532, 530, 1, 0, 0, 0, 
	533, 540, 3, 162, 81, 0, 534, 535, 5, 75, 0, 0, 535, 536, 5, 113, 0, 0, 
	536, 537, 3, 124, 62, 0, 537, 538, 5, 114, 0, 0, 538, 540, 1, 0, 0, 0, 
	539, 533, 1, 0, 0, 0, 539, 534, 1, 0, 0, 0, 540, 101, 1, 0, 0, 0, 541, 
	542, 7, 3, 0, 0, 542, 103, 1, 0, 0, 0, 543, 544, 5, 63, 0, 0, 544, 545, 
	5, 72, 0, 0, 545, 546, 3, 108, 54, 0, 546, 105, 1, 0, 0, 0, 547, 551, 3, 
	122, 61, 0, 548, 550, 7, 4, 0, 0, 549, 548, 1, 0, 0, 0, 550, 553, 1, 0, 
	0, 0, 551, 549, 1, 0, 0, 0, 551, 552, 1, 0, 0, 0, 552, 107, 1, 0, 0, 0, 
	553, 551, 1, 0, 0, 0, 554, 559, 3, 106, 53, 0, 555, 556, 5, 108, 0, 0, 
	556, 558, 3, 106, 53, 0, 557, 555, 1, 0, 0, 0, 558, 561, 1, 0, 0, 0, 559, 
	557, 1, 0, 0, 0, 559, 560, 1, 0, 0, 0, 560, 109, 1, 0, 0, 0, 561, 559, 
	1, 0, 0, 0, 562, 563, 5, 71, 0, 0, 563, 564, 3, 112, 56, 0, 564, 111, 1, 
	0, 0, 0, 565, 566, 6, 56, -1, 0, 566, 567, 5, 113, 0, 0, 567, 568, 3, 112, 
	56, 0, 568, 569, 5, 114, 0, 0, 569, 572, 1, 0, 0, 0, 570, 572, 3, 116, 
	58, 0, 571, 565, 1, 0, 0, 0, 571, 570, 1, 0, 0, 0, 572, 579, 1, 0, 0, 0, 
	573, 574, 10, 2, 0, 0, 574, 575, 3, 114, 57, 0, 575, 576, 3, 112, 56, 3, 
	576, 578, 1, 0, 0, 0, 577, 573, 1, 0, 0, 0, 578, 581, 1, 0, 0, 0, 579, 
	577, 1, 0, 0, 0, 579, 580, 1, 0, 0, 0, 580, 113, 1, 0, 0, 0, 581, 579, 
	1, 0, 0, 0, 582, 583, 7, 2, 0, 0, 583, 115, 1, 0, 0, 0, 584, 585, 3, 118, 
	59, 0, 585, 117, 1, 0, 0, 0, 586, 587, 3, 122, 61, 0, 587, 588, 3, 120, 
	60, 0, 588, 589, 3, 122, 61, 0, 589, 119, 1, 0, 0, 0, 590, 599, 5, 99, 
	0, 0, 591, 599, 5, 100, 0, 0, 592, 599, 5, 101, 0, 0, 593, 599, 5, 104, 
	0, 0, 594, 599, 5, 105, 0, 0, 595, 599, 5, 102, 0, 0, 596, 599, 5, 103, 
	0, 0, 597, 599, 7, 5, 0, 0, 598, 590, 1, 0, 0, 0, 598, 591, 1, 0, 0, 0, 
	598, 592, 1, 0, 0, 0, 598, 593, 1, 0, 0, 0, 598, 594, 1, 0, 0, 0, 598, 
	595, 1, 0, 0, 0, 598, 596, 1, 0, 0, 0, 598, 597, 1, 0, 0, 0, 599, 121, 
	1, 0, 0, 0, 600, 601, 6, 61, -1, 0, 601, 602, 5, 113, 0, 0, 602, 603, 3, 
	122, 61, 0, 603, 604, 5, 114, 0, 0, 604, 609, 1, 0, 0, 0, 605, 609, 3, 
	128, 64, 0, 606, 609, 3, 136, 68, 0, 607, 609, 3, 124, 62, 0, 608, 600, 
	1, 0, 0, 0, 608, 605, 1, 0, 0, 0, 608, 606, 1, 0, 0, 0, 608, 607, 1, 0, 
	0, 0, 609, 624, 1, 0, 0, 0, 610, 611, 10, 8, 0, 0, 611, 612, 5, 118, 0, 
	0, 612, 623, 3, 122, 61, 9, 613, 614, 10, 7, 0, 0, 614, 615, 5, 117, 0, 
	0, 615, 623, 3, 122, 61, 8, 616, 617, 10, 6, 0, 0, 617, 618, 5, 115, 0, 
	0, 618, 623, 3, 122, 61, 7, 619, 620, 10, 5, 0, 0, 620, 621, 5, 116, 0, 
	0, 621, 623, 3, 122, 61, 6, 622, 610, 1, 0, 0, 0, 622, 613, 1, 0, 0, 0, 
	622, 616, 1, 0, 0, 0, 622, 619, 1, 0, 0, 0, 623, 626, 1, 0, 0, 0, 624, 
	622, 1, 0, 0, 0, 624, 625, 1, 0, 0, 0, 625, 123, 1, 0, 0, 0, 626, 624, 
	1, 0, 0, 0, 627, 628, 3, 150, 75, 0, 628, 629, 3, 126, 63, 0, 629, 125, 
	1, 0, 0, 0, 630, 631, 7, 6, 0, 0, 631, 127, 1, 0, 0, 0, 632, 633, 3, 130, 
	65, 0, 633, 635, 5, 113, 0, 0, 634, 636, 3, 132, 66, 0, 635, 634, 1, 0, 
	0, 0, 635, 636, 1, 0, 0, 0, 636, 637, 1, 0, 0, 0, 637, 638, 5, 114, 0, 
	0, 638, 129, 1, 0, 0, 0, 639, 640, 7, 7, 0, 0, 640, 131, 1, 0, 0, 0, 641, 
	646, 3, 134, 67, 0, 642, 643, 5, 108, 0, 0, 643, 645, 3, 134, 67, 0, 644, 
	642, 1, 0, 0, 0, 645, 648, 1, 0, 0, 0, 646, 644, 1, 0, 0, 0, 646, 647, 
	1, 0, 0, 0, 647, 133, 1, 0, 0, 0, 648, 646, 1, 0, 0, 0, 649, 652, 3, 122, 
	61, 0, 650, 652, 3, 80, 40, 0, 651, 649, 1, 0, 0, 0, 651, 650, 1, 0, 0, 
	0, 652, 135, 1, 0, 0, 0, 653, 655, 3, 162, 81, 0, 654, 656, 3, 138, 69, 
	0, 655, 654, 1, 0, 0, 0, 655, 656, 1, 0, 0, 0, 656, 660, 1, 0, 0, 0, 657, 
	660, 3, 152, 76, 0, 658, 660, 3, 150, 75, 0, 659, 653, 1, 0, 0, 0, 659, 
	657, 1, 0, 0, 0, 659, 658, 1, 0, 0, 0, 660, 137, 1, 0, 0, 0, 661, 662, 
	5, 111, 0, 0, 662, 663, 3, 80, 40, 0, 663, 664, 5, 112, 0, 0, 664, 139, 
	1, 0, 0, 0, 665, 666, 3, 148, 74, 0, 666, 141, 1, 0, 0, 0, 667, 668, 5, 
	109, 0, 0, 668, 673, 3, 144, 72, 0, 669, 670, 5, 108, 0, 0, 670, 672, 3, 
	144, 72, 0, 671, 669, 1, 0, 0, 0, 672, 675, 1, 0, 0, 0, 673, 671, 1, 0, 
	0, 0, 673, 674, 1, 0, 0, 0, 674, 676, 1, 0, 0, 0, 675, 673, 1, 0, 0, 0, 
	676, 677, 5, 110, 0, 0, 677, 681, 1, 0, 0, 0, 678, 679, 5, 109, 0, 0, 679, 
	681, 5, 110, 0, 0, 680, 667, 1, 0, 0, 0, 680, 678, 1, 0, 0, 0, 681, 143, 
	1, 0, 0, 0, 682, 683, 5, 4, 0, 0, 683, 684, 5, 98, 0, 0, 684, 685, 3, 148, 
	74, 0, 685, 145, 1, 0, 0, 0, 686, 687, 5, 111, 0, 0, 687, 692, 3, 148, 
	74, 0, 688, 689, 5, 108, 0, 0, 689, 691, 3, 148, 74, 0, 690, 688, 1, 0, 
	0, 0, 691, 694, 1, 0, 0, 0, 692, 690, 1, 0, 0, 0, 692, 693, 1, 0, 0, 0, 
	693, 695, 1, 0, 0, 0, 694, 692, 1, 0, 0, 0, 695, 696, 5, 112, 0, 0, 696, 
	700, 1, 0, 0, 0, 697, 698, 5, 111, 0, 0, 698, 700, 5, 112, 0, 0, 699, 686, 
	1, 0, 0, 0, 699, 697, 1, 0, 0, 0, 700, 147, 1, 0, 0, 0, 701, 710, 5, 4, 
	0, 0, 702, 710, 3, 150, 75, 0, 703, 710, 3, 152, 76, 0, 704, 710, 3, 142, 
	71, 0, 705, 710, 3, 146, 73, 0, 706, 710, 5, 1, 0, 0, 707, 710, 5, 2, 0, 
	0, 708, 710, 5, 3, 0, 0, 709, 701, 1, 0, 0, 0, 709, 702, 1, 0, 0, 0, 709, 
	703, 1, 0, 0, 0, 709, 704, 1, 0, 0, 0, 709, 705, 1, 0, 0, 0, 709, 706, 
	1, 0, 0, 0, 709, 707, 1, 0, 0, 0, 709, 708, 1, 0, 0, 0, 710, 149, 1, 0, 
	0, 0, 711, 713, 7, 8, 0, 0, 712, 711, 1, 0, 0, 0, 712, 713, 1, 0, 0, 0, 
	713, 714, 1, 0, 0, 0, 714, 715, 5, 122, 0, 0, 715, 151, 1, 0, 0, 0, 716, 
	718, 7, 8, 0, 0, 717, 716, 1, 0, 0, 0, 717, 718, 1, 0, 0, 0, 718, 719, 
	1, 0, 0, 0, 719, 720, 5, 123, 0, 0, 720, 153, 1, 0, 0, 0, 721, 722, 5, 
	51, 0, 0, 722, 723, 5, 122, 0, 0, 723, 155, 1, 0, 0, 0, 724, 725, 3, 162, 
	81, 0, 725, 157, 1, 0, 0, 0, 726, 727, 3, 162, 81, 0, 727, 159, 1, 0, 0, 
	0, 728, 729, 3, 162, 81, 0, 729, 161, 1, 0, 0, 0, 730, 733, 5, 121, 0, 
	0, 731, 733, 3, 164, 82, 0, 732, 730, 1, 0, 0, 0, 732, 731, 1, 0, 0, 0, 
	733, 741, 1, 0, 0, 0, 734, 737, 5, 97, 0, 0, 735, 738, 5, 121, 0, 0, 736, 
	738, 3, 164, 82, 0, 737, 735, 1, 0, 0, 0, 737, 736, 1, 0, 0, 0, 738, 740, 
	1, 0, 0, 0, 739, 734, 1, 0, 0, 0, 740, 743, 1, 0, 0, 0, 741, 739, 1, 0, 
	0, 0, 741, 742, 1, 0, 0, 0, 742, 163, 1, 0, 0, 0, 743, 741, 1, 0, 0, 0, 
	744, 745, 7, 9, 0, 0, 745, 165, 1, 0, 0, 0, 64, 175, 197, 233, 238, 249, 
	254, 268, 273, 299, 302, 308, 314, 317, 337, 340, 353, 357, 360, 363, 366, 
	369, 377, 387, 392, 413, 426, 428, 444, 452, 458, 465, 473, 487, 493, 499, 
	503, 508, 520, 523, 530, 539, 551, 559, 571, 579, 598, 608, 622, 624, 635, 
	646, 651, 655, 659, 673, 680, 692, 699, 709, 712, 717, 732, 737, 741,
}
  deserializer := antlr.NewATNDeserializer(nil)
  staticData.atn = deserializer.Deserialize(staticData.serializedATN)
  atn := staticData.atn
  staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
  decisionToDFA := staticData.decisionToDFA
  for index, state := range atn.DecisionToState {
    decisionToDFA[index] = antlr.NewDFA(state, index)
  }
}

// SQLParserInit initializes any static state used to implement SQLParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewSQLParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SQLParserInit() {
  staticData := &sqlParserStaticData
  staticData.once.Do(sqlParserInit)
}

// NewSQLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewSQLParser(input antlr.TokenStream) *SQLParser {
	SQLParserInit()
	this := new(SQLParser)
	this.BaseParser = antlr.NewBaseParser(input)
  staticData := &sqlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "SQL.g4"

	return this
}


// SQLParser tokens.
const (
	SQLParserEOF = antlr.TokenEOF
	SQLParserT__0 = 1
	SQLParserT__1 = 2
	SQLParserT__2 = 3
	SQLParserSTRING = 4
	SQLParserWS = 5
	SQLParserT_CREATE = 6
	SQLParserT_UPDATE = 7
	SQLParserT_SET = 8
	SQLParserT_DROP = 9
	SQLParserT_INTERVAL = 10
	SQLParserT_INTERVAL_NAME = 11
	SQLParserT_SHARD = 12
	SQLParserT_REPLICATION = 13
	SQLParserT_TTL = 14
	SQLParserT_META_TTL = 15
	SQLParserT_PAST_TTL = 16
	SQLParserT_FUTURE_TTL = 17
	SQLParserT_KILL = 18
	SQLParserT_ON = 19
	SQLParserT_SHOW = 20
	SQLParserT_USE = 21
	SQLParserT_STATE_REPO = 22
	SQLParserT_STATE_MACHINE = 23
	SQLParserT_MASTER = 24
	SQLParserT_METADATA = 25
	SQLParserT_TYPES = 26
	SQLParserT_TYPE = 27
	SQLParserT_STORAGES = 28
	SQLParserT_STORAGE = 29
	SQLParserT_BROKER = 30
	SQLParserT_ALIVE = 31
	SQLParserT_SCHEMAS = 32
	SQLParserT_DATASBAE = 33
	SQLParserT_DATASBAES = 34
	SQLParserT_NAMESPACE = 35
	SQLParserT_NAMESPACES = 36
	SQLParserT_NODE = 37
	SQLParserT_METRICS = 38
	SQLParserT_METRIC = 39
	SQLParserT_FIELD = 40
	SQLParserT_FIELDS = 41
	SQLParserT_TAG = 42
	SQLParserT_INFO = 43
	SQLParserT_KEYS = 44
	SQLParserT_KEY = 45
	SQLParserT_WITH = 46
	SQLParserT_VALUES = 47
	SQLParserT_VALUE = 48
	SQLParserT_FROM = 49
	SQLParserT_WHERE = 50
	SQLParserT_LIMIT = 51
	SQLParserT_QUERIES = 52
	SQLParserT_QUERY = 53
	SQLParserT_EXPLAIN = 54
	SQLParserT_WITH_VALUE = 55
	SQLParserT_SELECT = 56
	SQLParserT_AS = 57
	SQLParserT_AND = 58
	SQLParserT_OR = 59
	SQLParserT_FILL = 60
	SQLParserT_NULL = 61
	SQLParserT_PREVIOUS = 62
	SQLParserT_ORDER = 63
	SQLParserT_ASC = 64
	SQLParserT_DESC = 65
	SQLParserT_LIKE = 66
	SQLParserT_NOT = 67
	SQLParserT_BETWEEN = 68
	SQLParserT_IS = 69
	SQLParserT_GROUP = 70
	SQLParserT_HAVING = 71
	SQLParserT_BY = 72
	SQLParserT_FOR = 73
	SQLParserT_STATS = 74
	SQLParserT_TIME = 75
	SQLParserT_NOW = 76
	SQLParserT_IN = 77
	SQLParserT_LOG = 78
	SQLParserT_PROFILE = 79
	SQLParserT_SUM = 80
	SQLParserT_MIN = 81
	SQLParserT_MAX = 82
	SQLParserT_COUNT = 83
	SQLParserT_LAST = 84
	SQLParserT_FIRST = 85
	SQLParserT_AVG = 86
	SQLParserT_STDDEV = 87
	SQLParserT_QUANTILE = 88
	SQLParserT_RATE = 89
	SQLParserT_SECOND = 90
	SQLParserT_MINUTE = 91
	SQLParserT_HOUR = 92
	SQLParserT_DAY = 93
	SQLParserT_WEEK = 94
	SQLParserT_MONTH = 95
	SQLParserT_YEAR = 96
	SQLParserT_DOT = 97
	SQLParserT_COLON = 98
	SQLParserT_EQUAL = 99
	SQLParserT_NOTEQUAL = 100
	SQLParserT_NOTEQUAL2 = 101
	SQLParserT_GREATER = 102
	SQLParserT_GREATEREQUAL = 103
	SQLParserT_LESS = 104
	SQLParserT_LESSEQUAL = 105
	SQLParserT_REGEXP = 106
	SQLParserT_NEQREGEXP = 107
	SQLParserT_COMMA = 108
	SQLParserT_OPEN_B = 109
	SQLParserT_CLOSE_B = 110
	SQLParserT_OPEN_SB = 111
	SQLParserT_CLOSE_SB = 112
	SQLParserT_OPEN_P = 113
	SQLParserT_CLOSE_P = 114
	SQLParserT_ADD = 115
	SQLParserT_SUB = 116
	SQLParserT_DIV = 117
	SQLParserT_MUL = 118
	SQLParserT_MOD = 119
	SQLParserT_UNDERLINE = 120
	SQLParserL_ID = 121
	SQLParserL_INT = 122
	SQLParserL_DEC = 123
)

// SQLParser rules.
const (
	SQLParserRULE_statement = 0
	SQLParserRULE_useStmt = 1
	SQLParserRULE_showStmt = 2
	SQLParserRULE_showMasterStmt = 3
	SQLParserRULE_showStoragesStmt = 4
	SQLParserRULE_showMetadataTypesStmt = 5
	SQLParserRULE_showBrokerMetaStmt = 6
	SQLParserRULE_showMasterMetaStmt = 7
	SQLParserRULE_showStorageMetaStmt = 8
	SQLParserRULE_showAliveStmt = 9
	SQLParserRULE_showReplicationStmt = 10
	SQLParserRULE_showBrokerMetricStmt = 11
	SQLParserRULE_showStorageMetricStmt = 12
	SQLParserRULE_createStorageStmt = 13
	SQLParserRULE_showSchemasStmt = 14
	SQLParserRULE_createDatabaseStmt = 15
	SQLParserRULE_dropDatabaseStmt = 16
	SQLParserRULE_showDatabaseStmt = 17
	SQLParserRULE_showNameSpacesStmt = 18
	SQLParserRULE_showMetricsStmt = 19
	SQLParserRULE_showFieldsStmt = 20
	SQLParserRULE_showTagKeysStmt = 21
	SQLParserRULE_showTagValuesStmt = 22
	SQLParserRULE_prefix = 23
	SQLParserRULE_withTagKey = 24
	SQLParserRULE_namespace = 25
	SQLParserRULE_databaseName = 26
	SQLParserRULE_source = 27
	SQLParserRULE_queryStmt = 28
	SQLParserRULE_sourceAndSelect = 29
	SQLParserRULE_selectExpr = 30
	SQLParserRULE_fields = 31
	SQLParserRULE_field = 32
	SQLParserRULE_alias = 33
	SQLParserRULE_storageFilter = 34
	SQLParserRULE_databaseFilter = 35
	SQLParserRULE_typeFilter = 36
	SQLParserRULE_fromClause = 37
	SQLParserRULE_whereClause = 38
	SQLParserRULE_conditionExpr = 39
	SQLParserRULE_tagFilterExpr = 40
	SQLParserRULE_tagValueList = 41
	SQLParserRULE_metricListFilter = 42
	SQLParserRULE_metricList = 43
	SQLParserRULE_timeRangeExpr = 44
	SQLParserRULE_timeExpr = 45
	SQLParserRULE_nowExpr = 46
	SQLParserRULE_nowFunc = 47
	SQLParserRULE_groupByClause = 48
	SQLParserRULE_groupByKeys = 49
	SQLParserRULE_groupByKey = 50
	SQLParserRULE_fillOption = 51
	SQLParserRULE_orderByClause = 52
	SQLParserRULE_sortField = 53
	SQLParserRULE_sortFields = 54
	SQLParserRULE_havingClause = 55
	SQLParserRULE_boolExpr = 56
	SQLParserRULE_boolExprLogicalOp = 57
	SQLParserRULE_boolExprAtom = 58
	SQLParserRULE_binaryExpr = 59
	SQLParserRULE_binaryOperator = 60
	SQLParserRULE_fieldExpr = 61
	SQLParserRULE_durationLit = 62
	SQLParserRULE_intervalItem = 63
	SQLParserRULE_exprFunc = 64
	SQLParserRULE_funcName = 65
	SQLParserRULE_exprFuncParams = 66
	SQLParserRULE_funcParam = 67
	SQLParserRULE_exprAtom = 68
	SQLParserRULE_identFilter = 69
	SQLParserRULE_json = 70
	SQLParserRULE_obj = 71
	SQLParserRULE_pair = 72
	SQLParserRULE_arr = 73
	SQLParserRULE_value = 74
	SQLParserRULE_intNumber = 75
	SQLParserRULE_decNumber = 76
	SQLParserRULE_limitClause = 77
	SQLParserRULE_metricName = 78
	SQLParserRULE_tagKey = 79
	SQLParserRULE_tagValue = 80
	SQLParserRULE_ident = 81
	SQLParserRULE_nonReservedWords = 82
)

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) ShowStmt() IShowStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowStmtContext)
}

func (s *StatementContext) CreateStorageStmt() ICreateStorageStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateStorageStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateStorageStmtContext)
}

func (s *StatementContext) UseStmt() IUseStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUseStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUseStmtContext)
}

func (s *StatementContext) QueryStmt() IQueryStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryStmtContext)
}

func (s *StatementContext) CreateDatabaseStmt() ICreateDatabaseStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateDatabaseStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateDatabaseStmtContext)
}

func (s *StatementContext) DropDatabaseStmt() IDropDatabaseStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDropDatabaseStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDropDatabaseStmtContext)
}

func (s *StatementContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *StatementContext) EOF() antlr.TerminalNode {
	return s.GetToken(SQLParserEOF, 0)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitStatement(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Statement() (localctx IStatementContext) {
	this := p
	_ = this

	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SQLParserRULE_statement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(175)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(166)
			p.ShowStmt()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(167)
			p.CreateStorageStmt()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(168)
			p.UseStmt()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(169)
			p.QueryStmt()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(170)
			p.CreateDatabaseStmt()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(171)
			p.DropDatabaseStmt()
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(172)
			p.Ident()
		}
		{
			p.SetState(173)
			p.Match(SQLParserEOF)
		}

	}


	return localctx
}


// IUseStmtContext is an interface to support dynamic dispatch.
type IUseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUseStmtContext differentiates from other interfaces.
	IsUseStmtContext()
}

type UseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUseStmtContext() *UseStmtContext {
	var p = new(UseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_useStmt
	return p
}

func (*UseStmtContext) IsUseStmtContext() {}

func NewUseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseStmtContext {
	var p = new(UseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_useStmt

	return p
}

func (s *UseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *UseStmtContext) T_USE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_USE, 0)
}

func (s *UseStmtContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *UseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *UseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterUseStmt(s)
	}
}

func (s *UseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitUseStmt(s)
	}
}

func (s *UseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitUseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) UseStmt() (localctx IUseStmtContext) {
	this := p
	_ = this

	localctx = NewUseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SQLParserRULE_useStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(177)
		p.Match(SQLParserT_USE)
	}
	{
		p.SetState(178)
		p.Ident()
	}



	return localctx
}


// IShowStmtContext is an interface to support dynamic dispatch.
type IShowStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStmtContext differentiates from other interfaces.
	IsShowStmtContext()
}

type ShowStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStmtContext() *ShowStmtContext {
	var p = new(ShowStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStmt
	return p
}

func (*ShowStmtContext) IsShowStmtContext() {}

func NewShowStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStmtContext {
	var p = new(ShowStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStmt

	return p
}

func (s *ShowStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStmtContext) ShowMasterStmt() IShowMasterStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowMasterStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowMasterStmtContext)
}

func (s *ShowStmtContext) ShowMetadataTypesStmt() IShowMetadataTypesStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowMetadataTypesStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowMetadataTypesStmtContext)
}

func (s *ShowStmtContext) ShowBrokerMetaStmt() IShowBrokerMetaStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowBrokerMetaStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowBrokerMetaStmtContext)
}

func (s *ShowStmtContext) ShowMasterMetaStmt() IShowMasterMetaStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowMasterMetaStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowMasterMetaStmtContext)
}

func (s *ShowStmtContext) ShowStorageMetaStmt() IShowStorageMetaStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowStorageMetaStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowStorageMetaStmtContext)
}

func (s *ShowStmtContext) ShowStoragesStmt() IShowStoragesStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowStoragesStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowStoragesStmtContext)
}

func (s *ShowStmtContext) ShowAliveStmt() IShowAliveStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowAliveStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowAliveStmtContext)
}

func (s *ShowStmtContext) ShowBrokerMetricStmt() IShowBrokerMetricStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowBrokerMetricStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowBrokerMetricStmtContext)
}

func (s *ShowStmtContext) ShowStorageMetricStmt() IShowStorageMetricStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowStorageMetricStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowStorageMetricStmtContext)
}

func (s *ShowStmtContext) ShowReplicationStmt() IShowReplicationStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowReplicationStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowReplicationStmtContext)
}

func (s *ShowStmtContext) ShowSchemasStmt() IShowSchemasStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowSchemasStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowSchemasStmtContext)
}

func (s *ShowStmtContext) ShowDatabaseStmt() IShowDatabaseStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowDatabaseStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowDatabaseStmtContext)
}

func (s *ShowStmtContext) ShowNameSpacesStmt() IShowNameSpacesStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowNameSpacesStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowNameSpacesStmtContext)
}

func (s *ShowStmtContext) ShowMetricsStmt() IShowMetricsStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowMetricsStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowMetricsStmtContext)
}

func (s *ShowStmtContext) ShowFieldsStmt() IShowFieldsStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowFieldsStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowFieldsStmtContext)
}

func (s *ShowStmtContext) ShowTagKeysStmt() IShowTagKeysStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowTagKeysStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowTagKeysStmtContext)
}

func (s *ShowStmtContext) ShowTagValuesStmt() IShowTagValuesStmtContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowTagValuesStmtContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowTagValuesStmtContext)
}

func (s *ShowStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStmt(s)
	}
}

func (s *ShowStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStmt(s)
	}
}

func (s *ShowStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowStmt() (localctx IShowStmtContext) {
	this := p
	_ = this

	localctx = NewShowStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_showStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(180)
			p.ShowMasterStmt()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(181)
			p.ShowMetadataTypesStmt()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(182)
			p.ShowBrokerMetaStmt()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(183)
			p.ShowMasterMetaStmt()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(184)
			p.ShowStorageMetaStmt()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(185)
			p.ShowStoragesStmt()
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(186)
			p.ShowAliveStmt()
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(187)
			p.ShowBrokerMetricStmt()
		}


	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(188)
			p.ShowStorageMetricStmt()
		}


	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(189)
			p.ShowReplicationStmt()
		}


	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(190)
			p.ShowSchemasStmt()
		}


	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(191)
			p.ShowDatabaseStmt()
		}


	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(192)
			p.ShowNameSpacesStmt()
		}


	case 14:
		p.EnterOuterAlt(localctx, 14)
		{
			p.SetState(193)
			p.ShowMetricsStmt()
		}


	case 15:
		p.EnterOuterAlt(localctx, 15)
		{
			p.SetState(194)
			p.ShowFieldsStmt()
		}


	case 16:
		p.EnterOuterAlt(localctx, 16)
		{
			p.SetState(195)
			p.ShowTagKeysStmt()
		}


	case 17:
		p.EnterOuterAlt(localctx, 17)
		{
			p.SetState(196)
			p.ShowTagValuesStmt()
		}

	}


	return localctx
}


// IShowMasterStmtContext is an interface to support dynamic dispatch.
type IShowMasterStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMasterStmtContext differentiates from other interfaces.
	IsShowMasterStmtContext()
}

type ShowMasterStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMasterStmtContext() *ShowMasterStmtContext {
	var p = new(ShowMasterStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMasterStmt
	return p
}

func (*ShowMasterStmtContext) IsShowMasterStmtContext() {}

func NewShowMasterStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMasterStmtContext {
	var p = new(ShowMasterStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMasterStmt

	return p
}

func (s *ShowMasterStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMasterStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMasterStmtContext) T_MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MASTER, 0)
}

func (s *ShowMasterStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMasterStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMasterStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMasterStmt(s)
	}
}

func (s *ShowMasterStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMasterStmt(s)
	}
}

func (s *ShowMasterStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowMasterStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowMasterStmt() (localctx IShowMasterStmtContext) {
	this := p
	_ = this

	localctx = NewShowMasterStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_showMasterStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(199)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(200)
		p.Match(SQLParserT_MASTER)
	}



	return localctx
}


// IShowStoragesStmtContext is an interface to support dynamic dispatch.
type IShowStoragesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStoragesStmtContext differentiates from other interfaces.
	IsShowStoragesStmtContext()
}

type ShowStoragesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStoragesStmtContext() *ShowStoragesStmtContext {
	var p = new(ShowStoragesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStoragesStmt
	return p
}

func (*ShowStoragesStmtContext) IsShowStoragesStmtContext() {}

func NewShowStoragesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStoragesStmtContext {
	var p = new(ShowStoragesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStoragesStmt

	return p
}

func (s *ShowStoragesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStoragesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStoragesStmtContext) T_STORAGES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGES, 0)
}

func (s *ShowStoragesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStoragesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStoragesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStoragesStmt(s)
	}
}

func (s *ShowStoragesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStoragesStmt(s)
	}
}

func (s *ShowStoragesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowStoragesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowStoragesStmt() (localctx IShowStoragesStmtContext) {
	this := p
	_ = this

	localctx = NewShowStoragesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_showStoragesStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(202)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(203)
		p.Match(SQLParserT_STORAGES)
	}



	return localctx
}


// IShowMetadataTypesStmtContext is an interface to support dynamic dispatch.
type IShowMetadataTypesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMetadataTypesStmtContext differentiates from other interfaces.
	IsShowMetadataTypesStmtContext()
}

type ShowMetadataTypesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMetadataTypesStmtContext() *ShowMetadataTypesStmtContext {
	var p = new(ShowMetadataTypesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMetadataTypesStmt
	return p
}

func (*ShowMetadataTypesStmtContext) IsShowMetadataTypesStmtContext() {}

func NewShowMetadataTypesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMetadataTypesStmtContext {
	var p = new(ShowMetadataTypesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMetadataTypesStmt

	return p
}

func (s *ShowMetadataTypesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMetadataTypesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMetadataTypesStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowMetadataTypesStmtContext) T_TYPES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPES, 0)
}

func (s *ShowMetadataTypesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetadataTypesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMetadataTypesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMetadataTypesStmt(s)
	}
}

func (s *ShowMetadataTypesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMetadataTypesStmt(s)
	}
}

func (s *ShowMetadataTypesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowMetadataTypesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowMetadataTypesStmt() (localctx IShowMetadataTypesStmtContext) {
	this := p
	_ = this

	localctx = NewShowMetadataTypesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_showMetadataTypesStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(205)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(206)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(207)
		p.Match(SQLParserT_TYPES)
	}



	return localctx
}


// IShowBrokerMetaStmtContext is an interface to support dynamic dispatch.
type IShowBrokerMetaStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowBrokerMetaStmtContext differentiates from other interfaces.
	IsShowBrokerMetaStmtContext()
}

type ShowBrokerMetaStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowBrokerMetaStmtContext() *ShowBrokerMetaStmtContext {
	var p = new(ShowBrokerMetaStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showBrokerMetaStmt
	return p
}

func (*ShowBrokerMetaStmtContext) IsShowBrokerMetaStmtContext() {}

func NewShowBrokerMetaStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowBrokerMetaStmtContext {
	var p = new(ShowBrokerMetaStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showBrokerMetaStmt

	return p
}

func (s *ShowBrokerMetaStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowBrokerMetaStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowBrokerMetaStmtContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *ShowBrokerMetaStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowBrokerMetaStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowBrokerMetaStmtContext) Source() ISourceContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *ShowBrokerMetaStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowBrokerMetaStmtContext) TypeFilter() ITypeFilterContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeFilterContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeFilterContext)
}

func (s *ShowBrokerMetaStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowBrokerMetaStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowBrokerMetaStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowBrokerMetaStmt(s)
	}
}

func (s *ShowBrokerMetaStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowBrokerMetaStmt(s)
	}
}

func (s *ShowBrokerMetaStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowBrokerMetaStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowBrokerMetaStmt() (localctx IShowBrokerMetaStmtContext) {
	this := p
	_ = this

	localctx = NewShowBrokerMetaStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_showBrokerMetaStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(209)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(210)
		p.Match(SQLParserT_BROKER)
	}
	{
		p.SetState(211)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(212)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(213)
		p.Source()
	}
	{
		p.SetState(214)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(215)
		p.TypeFilter()
	}



	return localctx
}


// IShowMasterMetaStmtContext is an interface to support dynamic dispatch.
type IShowMasterMetaStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMasterMetaStmtContext differentiates from other interfaces.
	IsShowMasterMetaStmtContext()
}

type ShowMasterMetaStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMasterMetaStmtContext() *ShowMasterMetaStmtContext {
	var p = new(ShowMasterMetaStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMasterMetaStmt
	return p
}

func (*ShowMasterMetaStmtContext) IsShowMasterMetaStmtContext() {}

func NewShowMasterMetaStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMasterMetaStmtContext {
	var p = new(ShowMasterMetaStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMasterMetaStmt

	return p
}

func (s *ShowMasterMetaStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMasterMetaStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMasterMetaStmtContext) T_MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MASTER, 0)
}

func (s *ShowMasterMetaStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowMasterMetaStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowMasterMetaStmtContext) Source() ISourceContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *ShowMasterMetaStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowMasterMetaStmtContext) TypeFilter() ITypeFilterContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeFilterContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeFilterContext)
}

func (s *ShowMasterMetaStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMasterMetaStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMasterMetaStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMasterMetaStmt(s)
	}
}

func (s *ShowMasterMetaStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMasterMetaStmt(s)
	}
}

func (s *ShowMasterMetaStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowMasterMetaStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowMasterMetaStmt() (localctx IShowMasterMetaStmtContext) {
	this := p
	_ = this

	localctx = NewShowMasterMetaStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_showMasterMetaStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(217)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(218)
		p.Match(SQLParserT_MASTER)
	}
	{
		p.SetState(219)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(220)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(221)
		p.Source()
	}
	{
		p.SetState(222)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(223)
		p.TypeFilter()
	}



	return localctx
}


// IShowStorageMetaStmtContext is an interface to support dynamic dispatch.
type IShowStorageMetaStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStorageMetaStmtContext differentiates from other interfaces.
	IsShowStorageMetaStmtContext()
}

type ShowStorageMetaStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStorageMetaStmtContext() *ShowStorageMetaStmtContext {
	var p = new(ShowStorageMetaStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStorageMetaStmt
	return p
}

func (*ShowStorageMetaStmtContext) IsShowStorageMetaStmtContext() {}

func NewShowStorageMetaStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStorageMetaStmtContext {
	var p = new(ShowStorageMetaStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStorageMetaStmt

	return p
}

func (s *ShowStorageMetaStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStorageMetaStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStorageMetaStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *ShowStorageMetaStmtContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *ShowStorageMetaStmtContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *ShowStorageMetaStmtContext) Source() ISourceContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceContext)
}

func (s *ShowStorageMetaStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowStorageMetaStmtContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ShowStorageMetaStmtContext) AllStorageFilter() []IStorageFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStorageFilterContext); ok {
			len++
		}
	}

	tst := make([]IStorageFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStorageFilterContext); ok {
			tst[i] = t.(IStorageFilterContext)
			i++
		}
	}

	return tst
}

func (s *ShowStorageMetaStmtContext) StorageFilter(i int) IStorageFilterContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStorageFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStorageFilterContext)
}

func (s *ShowStorageMetaStmtContext) AllTypeFilter() []ITypeFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeFilterContext); ok {
			len++
		}
	}

	tst := make([]ITypeFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeFilterContext); ok {
			tst[i] = t.(ITypeFilterContext)
			i++
		}
	}

	return tst
}

func (s *ShowStorageMetaStmtContext) TypeFilter(i int) ITypeFilterContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeFilterContext)
}

func (s *ShowStorageMetaStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStorageMetaStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStorageMetaStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStorageMetaStmt(s)
	}
}

func (s *ShowStorageMetaStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStorageMetaStmt(s)
	}
}

func (s *ShowStorageMetaStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowStorageMetaStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowStorageMetaStmt() (localctx IShowStorageMetaStmtContext) {
	this := p
	_ = this

	localctx = NewShowStorageMetaStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_showStorageMetaStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(225)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(226)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(227)
		p.Match(SQLParserT_METADATA)
	}
	{
		p.SetState(228)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(229)
		p.Source()
	}
	{
		p.SetState(230)
		p.Match(SQLParserT_WHERE)
	}
	p.SetState(233)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(231)
			p.StorageFilter()
		}


	case SQLParserT_TYPE:
		{
			p.SetState(232)
			p.TypeFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(235)
		p.Match(SQLParserT_AND)
	}
	p.SetState(238)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(236)
			p.StorageFilter()
		}


	case SQLParserT_TYPE:
		{
			p.SetState(237)
			p.TypeFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}



	return localctx
}


// IShowAliveStmtContext is an interface to support dynamic dispatch.
type IShowAliveStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowAliveStmtContext differentiates from other interfaces.
	IsShowAliveStmtContext()
}

type ShowAliveStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowAliveStmtContext() *ShowAliveStmtContext {
	var p = new(ShowAliveStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showAliveStmt
	return p
}

func (*ShowAliveStmtContext) IsShowAliveStmtContext() {}

func NewShowAliveStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowAliveStmtContext {
	var p = new(ShowAliveStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showAliveStmt

	return p
}

func (s *ShowAliveStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowAliveStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowAliveStmtContext) T_ALIVE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ALIVE, 0)
}

func (s *ShowAliveStmtContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *ShowAliveStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *ShowAliveStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowAliveStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowAliveStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowAliveStmt(s)
	}
}

func (s *ShowAliveStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowAliveStmt(s)
	}
}

func (s *ShowAliveStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowAliveStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowAliveStmt() (localctx IShowAliveStmtContext) {
	this := p
	_ = this

	localctx = NewShowAliveStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_showAliveStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(240)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(241)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_STORAGE || _la == SQLParserT_BROKER) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(242)
		p.Match(SQLParserT_ALIVE)
	}



	return localctx
}


// IShowReplicationStmtContext is an interface to support dynamic dispatch.
type IShowReplicationStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowReplicationStmtContext differentiates from other interfaces.
	IsShowReplicationStmtContext()
}

type ShowReplicationStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowReplicationStmtContext() *ShowReplicationStmtContext {
	var p = new(ShowReplicationStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showReplicationStmt
	return p
}

func (*ShowReplicationStmtContext) IsShowReplicationStmtContext() {}

func NewShowReplicationStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowReplicationStmtContext {
	var p = new(ShowReplicationStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showReplicationStmt

	return p
}

func (s *ShowReplicationStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowReplicationStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowReplicationStmtContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *ShowReplicationStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowReplicationStmtContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ShowReplicationStmtContext) AllStorageFilter() []IStorageFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStorageFilterContext); ok {
			len++
		}
	}

	tst := make([]IStorageFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStorageFilterContext); ok {
			tst[i] = t.(IStorageFilterContext)
			i++
		}
	}

	return tst
}

func (s *ShowReplicationStmtContext) StorageFilter(i int) IStorageFilterContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStorageFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStorageFilterContext)
}

func (s *ShowReplicationStmtContext) AllDatabaseFilter() []IDatabaseFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDatabaseFilterContext); ok {
			len++
		}
	}

	tst := make([]IDatabaseFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDatabaseFilterContext); ok {
			tst[i] = t.(IDatabaseFilterContext)
			i++
		}
	}

	return tst
}

func (s *ShowReplicationStmtContext) DatabaseFilter(i int) IDatabaseFilterContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDatabaseFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDatabaseFilterContext)
}

func (s *ShowReplicationStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowReplicationStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowReplicationStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowReplicationStmt(s)
	}
}

func (s *ShowReplicationStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowReplicationStmt(s)
	}
}

func (s *ShowReplicationStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowReplicationStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowReplicationStmt() (localctx IShowReplicationStmtContext) {
	this := p
	_ = this

	localctx = NewShowReplicationStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_showReplicationStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(244)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(245)
		p.Match(SQLParserT_REPLICATION)
	}
	{
		p.SetState(246)
		p.Match(SQLParserT_WHERE)
	}
	p.SetState(249)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(247)
			p.StorageFilter()
		}


	case SQLParserT_DATASBAE:
		{
			p.SetState(248)
			p.DatabaseFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(251)
		p.Match(SQLParserT_AND)
	}
	p.SetState(254)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(252)
			p.StorageFilter()
		}


	case SQLParserT_DATASBAE:
		{
			p.SetState(253)
			p.DatabaseFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}



	return localctx
}


// IShowBrokerMetricStmtContext is an interface to support dynamic dispatch.
type IShowBrokerMetricStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowBrokerMetricStmtContext differentiates from other interfaces.
	IsShowBrokerMetricStmtContext()
}

type ShowBrokerMetricStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowBrokerMetricStmtContext() *ShowBrokerMetricStmtContext {
	var p = new(ShowBrokerMetricStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showBrokerMetricStmt
	return p
}

func (*ShowBrokerMetricStmtContext) IsShowBrokerMetricStmtContext() {}

func NewShowBrokerMetricStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowBrokerMetricStmtContext {
	var p = new(ShowBrokerMetricStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showBrokerMetricStmt

	return p
}

func (s *ShowBrokerMetricStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowBrokerMetricStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowBrokerMetricStmtContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *ShowBrokerMetricStmtContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *ShowBrokerMetricStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowBrokerMetricStmtContext) MetricListFilter() IMetricListFilterContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMetricListFilterContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMetricListFilterContext)
}

func (s *ShowBrokerMetricStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowBrokerMetricStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowBrokerMetricStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowBrokerMetricStmt(s)
	}
}

func (s *ShowBrokerMetricStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowBrokerMetricStmt(s)
	}
}

func (s *ShowBrokerMetricStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowBrokerMetricStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowBrokerMetricStmt() (localctx IShowBrokerMetricStmtContext) {
	this := p
	_ = this

	localctx = NewShowBrokerMetricStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_showBrokerMetricStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(256)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(257)
		p.Match(SQLParserT_BROKER)
	}
	{
		p.SetState(258)
		p.Match(SQLParserT_METRIC)
	}
	{
		p.SetState(259)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(260)
		p.MetricListFilter()
	}



	return localctx
}


// IShowStorageMetricStmtContext is an interface to support dynamic dispatch.
type IShowStorageMetricStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowStorageMetricStmtContext differentiates from other interfaces.
	IsShowStorageMetricStmtContext()
}

type ShowStorageMetricStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStorageMetricStmtContext() *ShowStorageMetricStmtContext {
	var p = new(ShowStorageMetricStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showStorageMetricStmt
	return p
}

func (*ShowStorageMetricStmtContext) IsShowStorageMetricStmtContext() {}

func NewShowStorageMetricStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStorageMetricStmtContext {
	var p = new(ShowStorageMetricStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStorageMetricStmt

	return p
}

func (s *ShowStorageMetricStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStorageMetricStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowStorageMetricStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *ShowStorageMetricStmtContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *ShowStorageMetricStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowStorageMetricStmtContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ShowStorageMetricStmtContext) AllStorageFilter() []IStorageFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStorageFilterContext); ok {
			len++
		}
	}

	tst := make([]IStorageFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStorageFilterContext); ok {
			tst[i] = t.(IStorageFilterContext)
			i++
		}
	}

	return tst
}

func (s *ShowStorageMetricStmtContext) StorageFilter(i int) IStorageFilterContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStorageFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStorageFilterContext)
}

func (s *ShowStorageMetricStmtContext) AllMetricListFilter() []IMetricListFilterContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMetricListFilterContext); ok {
			len++
		}
	}

	tst := make([]IMetricListFilterContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMetricListFilterContext); ok {
			tst[i] = t.(IMetricListFilterContext)
			i++
		}
	}

	return tst
}

func (s *ShowStorageMetricStmtContext) MetricListFilter(i int) IMetricListFilterContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMetricListFilterContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMetricListFilterContext)
}

func (s *ShowStorageMetricStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStorageMetricStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowStorageMetricStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowStorageMetricStmt(s)
	}
}

func (s *ShowStorageMetricStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowStorageMetricStmt(s)
	}
}

func (s *ShowStorageMetricStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowStorageMetricStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowStorageMetricStmt() (localctx IShowStorageMetricStmtContext) {
	this := p
	_ = this

	localctx = NewShowStorageMetricStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_showStorageMetricStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(262)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(263)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(264)
		p.Match(SQLParserT_METRIC)
	}
	{
		p.SetState(265)
		p.Match(SQLParserT_WHERE)
	}
	p.SetState(268)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(266)
			p.StorageFilter()
		}


	case SQLParserT_METRIC:
		{
			p.SetState(267)
			p.MetricListFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	{
		p.SetState(270)
		p.Match(SQLParserT_AND)
	}
	p.SetState(273)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_STORAGE:
		{
			p.SetState(271)
			p.StorageFilter()
		}


	case SQLParserT_METRIC:
		{
			p.SetState(272)
			p.MetricListFilter()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}



	return localctx
}


// ICreateStorageStmtContext is an interface to support dynamic dispatch.
type ICreateStorageStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCreateStorageStmtContext differentiates from other interfaces.
	IsCreateStorageStmtContext()
}

type CreateStorageStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateStorageStmtContext() *CreateStorageStmtContext {
	var p = new(CreateStorageStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_createStorageStmt
	return p
}

func (*CreateStorageStmtContext) IsCreateStorageStmtContext() {}

func NewCreateStorageStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateStorageStmtContext {
	var p = new(CreateStorageStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createStorageStmt

	return p
}

func (s *CreateStorageStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateStorageStmtContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *CreateStorageStmtContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *CreateStorageStmtContext) Json() IJsonContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJsonContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJsonContext)
}

func (s *CreateStorageStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateStorageStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CreateStorageStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterCreateStorageStmt(s)
	}
}

func (s *CreateStorageStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitCreateStorageStmt(s)
	}
}

func (s *CreateStorageStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitCreateStorageStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) CreateStorageStmt() (localctx ICreateStorageStmtContext) {
	this := p
	_ = this

	localctx = NewCreateStorageStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SQLParserRULE_createStorageStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(275)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(276)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(277)
		p.Json()
	}



	return localctx
}


// IShowSchemasStmtContext is an interface to support dynamic dispatch.
type IShowSchemasStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowSchemasStmtContext differentiates from other interfaces.
	IsShowSchemasStmtContext()
}

type ShowSchemasStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowSchemasStmtContext() *ShowSchemasStmtContext {
	var p = new(ShowSchemasStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showSchemasStmt
	return p
}

func (*ShowSchemasStmtContext) IsShowSchemasStmtContext() {}

func NewShowSchemasStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowSchemasStmtContext {
	var p = new(ShowSchemasStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showSchemasStmt

	return p
}

func (s *ShowSchemasStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowSchemasStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowSchemasStmtContext) T_SCHEMAS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SCHEMAS, 0)
}

func (s *ShowSchemasStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowSchemasStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowSchemasStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowSchemasStmt(s)
	}
}

func (s *ShowSchemasStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowSchemasStmt(s)
	}
}

func (s *ShowSchemasStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowSchemasStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowSchemasStmt() (localctx IShowSchemasStmtContext) {
	this := p
	_ = this

	localctx = NewShowSchemasStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SQLParserRULE_showSchemasStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(279)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(280)
		p.Match(SQLParserT_SCHEMAS)
	}



	return localctx
}


// ICreateDatabaseStmtContext is an interface to support dynamic dispatch.
type ICreateDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCreateDatabaseStmtContext differentiates from other interfaces.
	IsCreateDatabaseStmtContext()
}

type CreateDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateDatabaseStmtContext() *CreateDatabaseStmtContext {
	var p = new(CreateDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_createDatabaseStmt
	return p
}

func (*CreateDatabaseStmtContext) IsCreateDatabaseStmtContext() {}

func NewCreateDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateDatabaseStmtContext {
	var p = new(CreateDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createDatabaseStmt

	return p
}

func (s *CreateDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateDatabaseStmtContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *CreateDatabaseStmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *CreateDatabaseStmtContext) Json() IJsonContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJsonContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJsonContext)
}

func (s *CreateDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *CreateDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterCreateDatabaseStmt(s)
	}
}

func (s *CreateDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitCreateDatabaseStmt(s)
	}
}

func (s *CreateDatabaseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitCreateDatabaseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) CreateDatabaseStmt() (localctx ICreateDatabaseStmtContext) {
	this := p
	_ = this

	localctx = NewCreateDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SQLParserRULE_createDatabaseStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(282)
		p.Match(SQLParserT_CREATE)
	}
	{
		p.SetState(283)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(284)
		p.Json()
	}



	return localctx
}


// IDropDatabaseStmtContext is an interface to support dynamic dispatch.
type IDropDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDropDatabaseStmtContext differentiates from other interfaces.
	IsDropDatabaseStmtContext()
}

type DropDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDropDatabaseStmtContext() *DropDatabaseStmtContext {
	var p = new(DropDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_dropDatabaseStmt
	return p
}

func (*DropDatabaseStmtContext) IsDropDatabaseStmtContext() {}

func NewDropDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DropDatabaseStmtContext {
	var p = new(DropDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dropDatabaseStmt

	return p
}

func (s *DropDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *DropDatabaseStmtContext) T_DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DROP, 0)
}

func (s *DropDatabaseStmtContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *DropDatabaseStmtContext) DatabaseName() IDatabaseNameContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDatabaseNameContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDatabaseNameContext)
}

func (s *DropDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DropDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DropDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDropDatabaseStmt(s)
	}
}

func (s *DropDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDropDatabaseStmt(s)
	}
}

func (s *DropDatabaseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDropDatabaseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) DropDatabaseStmt() (localctx IDropDatabaseStmtContext) {
	this := p
	_ = this

	localctx = NewDropDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_dropDatabaseStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(286)
		p.Match(SQLParserT_DROP)
	}
	{
		p.SetState(287)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(288)
		p.DatabaseName()
	}



	return localctx
}


// IShowDatabaseStmtContext is an interface to support dynamic dispatch.
type IShowDatabaseStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowDatabaseStmtContext differentiates from other interfaces.
	IsShowDatabaseStmtContext()
}

type ShowDatabaseStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowDatabaseStmtContext() *ShowDatabaseStmtContext {
	var p = new(ShowDatabaseStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showDatabaseStmt
	return p
}

func (*ShowDatabaseStmtContext) IsShowDatabaseStmtContext() {}

func NewShowDatabaseStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowDatabaseStmtContext {
	var p = new(ShowDatabaseStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showDatabaseStmt

	return p
}

func (s *ShowDatabaseStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowDatabaseStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowDatabaseStmtContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *ShowDatabaseStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowDatabaseStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowDatabaseStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowDatabaseStmt(s)
	}
}

func (s *ShowDatabaseStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowDatabaseStmt(s)
	}
}

func (s *ShowDatabaseStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowDatabaseStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowDatabaseStmt() (localctx IShowDatabaseStmtContext) {
	this := p
	_ = this

	localctx = NewShowDatabaseStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SQLParserRULE_showDatabaseStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(290)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(291)
		p.Match(SQLParserT_DATASBAES)
	}



	return localctx
}


// IShowNameSpacesStmtContext is an interface to support dynamic dispatch.
type IShowNameSpacesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowNameSpacesStmtContext differentiates from other interfaces.
	IsShowNameSpacesStmtContext()
}

type ShowNameSpacesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowNameSpacesStmtContext() *ShowNameSpacesStmtContext {
	var p = new(ShowNameSpacesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showNameSpacesStmt
	return p
}

func (*ShowNameSpacesStmtContext) IsShowNameSpacesStmtContext() {}

func NewShowNameSpacesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowNameSpacesStmtContext {
	var p = new(ShowNameSpacesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showNameSpacesStmt

	return p
}

func (s *ShowNameSpacesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowNameSpacesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowNameSpacesStmtContext) T_NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACES, 0)
}

func (s *ShowNameSpacesStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowNameSpacesStmtContext) T_NAMESPACE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACE, 0)
}

func (s *ShowNameSpacesStmtContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *ShowNameSpacesStmtContext) Prefix() IPrefixContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrefixContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrefixContext)
}

func (s *ShowNameSpacesStmtContext) LimitClause() ILimitClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowNameSpacesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowNameSpacesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowNameSpacesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowNameSpacesStmt(s)
	}
}

func (s *ShowNameSpacesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowNameSpacesStmt(s)
	}
}

func (s *ShowNameSpacesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowNameSpacesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowNameSpacesStmt() (localctx IShowNameSpacesStmtContext) {
	this := p
	_ = this

	localctx = NewShowNameSpacesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SQLParserRULE_showNameSpacesStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(293)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(294)
		p.Match(SQLParserT_NAMESPACES)
	}
	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(295)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(296)
			p.Match(SQLParserT_NAMESPACE)
		}
		{
			p.SetState(297)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(298)
			p.Prefix()
		}

	}
	p.SetState(302)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(301)
			p.LimitClause()
		}

	}



	return localctx
}


// IShowMetricsStmtContext is an interface to support dynamic dispatch.
type IShowMetricsStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowMetricsStmtContext differentiates from other interfaces.
	IsShowMetricsStmtContext()
}

type ShowMetricsStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMetricsStmtContext() *ShowMetricsStmtContext {
	var p = new(ShowMetricsStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showMetricsStmt
	return p
}

func (*ShowMetricsStmtContext) IsShowMetricsStmtContext() {}

func NewShowMetricsStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMetricsStmtContext {
	var p = new(ShowMetricsStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMetricsStmt

	return p
}

func (s *ShowMetricsStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMetricsStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowMetricsStmtContext) T_METRICS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRICS, 0)
}

func (s *ShowMetricsStmtContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *ShowMetricsStmtContext) Namespace() INamespaceContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *ShowMetricsStmtContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *ShowMetricsStmtContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *ShowMetricsStmtContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *ShowMetricsStmtContext) Prefix() IPrefixContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrefixContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrefixContext)
}

func (s *ShowMetricsStmtContext) LimitClause() ILimitClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowMetricsStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetricsStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowMetricsStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowMetricsStmt(s)
	}
}

func (s *ShowMetricsStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowMetricsStmt(s)
	}
}

func (s *ShowMetricsStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowMetricsStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowMetricsStmt() (localctx IShowMetricsStmtContext) {
	this := p
	_ = this

	localctx = NewShowMetricsStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SQLParserRULE_showMetricsStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(304)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(305)
		p.Match(SQLParserT_METRICS)
	}
	p.SetState(308)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(306)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(307)
			p.Namespace()
		}

	}
	p.SetState(314)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(310)
			p.Match(SQLParserT_WHERE)
		}
		{
			p.SetState(311)
			p.Match(SQLParserT_METRIC)
		}
		{
			p.SetState(312)
			p.Match(SQLParserT_EQUAL)
		}
		{
			p.SetState(313)
			p.Prefix()
		}

	}
	p.SetState(317)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(316)
			p.LimitClause()
		}

	}



	return localctx
}


// IShowFieldsStmtContext is an interface to support dynamic dispatch.
type IShowFieldsStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowFieldsStmtContext differentiates from other interfaces.
	IsShowFieldsStmtContext()
}

type ShowFieldsStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowFieldsStmtContext() *ShowFieldsStmtContext {
	var p = new(ShowFieldsStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showFieldsStmt
	return p
}

func (*ShowFieldsStmtContext) IsShowFieldsStmtContext() {}

func NewShowFieldsStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowFieldsStmtContext {
	var p = new(ShowFieldsStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showFieldsStmt

	return p
}

func (s *ShowFieldsStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowFieldsStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowFieldsStmtContext) T_FIELDS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELDS, 0)
}

func (s *ShowFieldsStmtContext) FromClause() IFromClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *ShowFieldsStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowFieldsStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowFieldsStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowFieldsStmt(s)
	}
}

func (s *ShowFieldsStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowFieldsStmt(s)
	}
}

func (s *ShowFieldsStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowFieldsStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowFieldsStmt() (localctx IShowFieldsStmtContext) {
	this := p
	_ = this

	localctx = NewShowFieldsStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SQLParserRULE_showFieldsStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(319)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(320)
		p.Match(SQLParserT_FIELDS)
	}
	{
		p.SetState(321)
		p.FromClause()
	}



	return localctx
}


// IShowTagKeysStmtContext is an interface to support dynamic dispatch.
type IShowTagKeysStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowTagKeysStmtContext differentiates from other interfaces.
	IsShowTagKeysStmtContext()
}

type ShowTagKeysStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowTagKeysStmtContext() *ShowTagKeysStmtContext {
	var p = new(ShowTagKeysStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showTagKeysStmt
	return p
}

func (*ShowTagKeysStmtContext) IsShowTagKeysStmtContext() {}

func NewShowTagKeysStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowTagKeysStmtContext {
	var p = new(ShowTagKeysStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showTagKeysStmt

	return p
}

func (s *ShowTagKeysStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowTagKeysStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowTagKeysStmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *ShowTagKeysStmtContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *ShowTagKeysStmtContext) FromClause() IFromClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *ShowTagKeysStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagKeysStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowTagKeysStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowTagKeysStmt(s)
	}
}

func (s *ShowTagKeysStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowTagKeysStmt(s)
	}
}

func (s *ShowTagKeysStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowTagKeysStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowTagKeysStmt() (localctx IShowTagKeysStmtContext) {
	this := p
	_ = this

	localctx = NewShowTagKeysStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SQLParserRULE_showTagKeysStmt)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(323)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(324)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(325)
		p.Match(SQLParserT_KEYS)
	}
	{
		p.SetState(326)
		p.FromClause()
	}



	return localctx
}


// IShowTagValuesStmtContext is an interface to support dynamic dispatch.
type IShowTagValuesStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowTagValuesStmtContext differentiates from other interfaces.
	IsShowTagValuesStmtContext()
}

type ShowTagValuesStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowTagValuesStmtContext() *ShowTagValuesStmtContext {
	var p = new(ShowTagValuesStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_showTagValuesStmt
	return p
}

func (*ShowTagValuesStmtContext) IsShowTagValuesStmtContext() {}

func NewShowTagValuesStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowTagValuesStmtContext {
	var p = new(ShowTagValuesStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showTagValuesStmt

	return p
}

func (s *ShowTagValuesStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowTagValuesStmtContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *ShowTagValuesStmtContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *ShowTagValuesStmtContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *ShowTagValuesStmtContext) FromClause() IFromClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *ShowTagValuesStmtContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *ShowTagValuesStmtContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *ShowTagValuesStmtContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *ShowTagValuesStmtContext) WithTagKey() IWithTagKeyContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWithTagKeyContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWithTagKeyContext)
}

func (s *ShowTagValuesStmtContext) WhereClause() IWhereClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereClauseContext)
}

func (s *ShowTagValuesStmtContext) LimitClause() ILimitClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *ShowTagValuesStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagValuesStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ShowTagValuesStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterShowTagValuesStmt(s)
	}
}

func (s *ShowTagValuesStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitShowTagValuesStmt(s)
	}
}

func (s *ShowTagValuesStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitShowTagValuesStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ShowTagValuesStmt() (localctx IShowTagValuesStmtContext) {
	this := p
	_ = this

	localctx = NewShowTagValuesStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SQLParserRULE_showTagValuesStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(328)
		p.Match(SQLParserT_SHOW)
	}
	{
		p.SetState(329)
		p.Match(SQLParserT_TAG)
	}
	{
		p.SetState(330)
		p.Match(SQLParserT_VALUES)
	}
	{
		p.SetState(331)
		p.FromClause()
	}
	{
		p.SetState(332)
		p.Match(SQLParserT_WITH)
	}
	{
		p.SetState(333)
		p.Match(SQLParserT_KEY)
	}
	{
		p.SetState(334)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(335)
		p.WithTagKey()
	}
	p.SetState(337)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(336)
			p.WhereClause()
		}

	}
	p.SetState(340)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(339)
			p.LimitClause()
		}

	}



	return localctx
}


// IPrefixContext is an interface to support dynamic dispatch.
type IPrefixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPrefixContext differentiates from other interfaces.
	IsPrefixContext()
}

type PrefixContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrefixContext() *PrefixContext {
	var p = new(PrefixContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_prefix
	return p
}

func (*PrefixContext) IsPrefixContext() {}

func NewPrefixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrefixContext {
	var p = new(PrefixContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_prefix

	return p
}

func (s *PrefixContext) GetParser() antlr.Parser { return s.parser }

func (s *PrefixContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *PrefixContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrefixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PrefixContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterPrefix(s)
	}
}

func (s *PrefixContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitPrefix(s)
	}
}

func (s *PrefixContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitPrefix(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Prefix() (localctx IPrefixContext) {
	this := p
	_ = this

	localctx = NewPrefixContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, SQLParserRULE_prefix)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(342)
		p.Ident()
	}



	return localctx
}


// IWithTagKeyContext is an interface to support dynamic dispatch.
type IWithTagKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWithTagKeyContext differentiates from other interfaces.
	IsWithTagKeyContext()
}

type WithTagKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithTagKeyContext() *WithTagKeyContext {
	var p = new(WithTagKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_withTagKey
	return p
}

func (*WithTagKeyContext) IsWithTagKeyContext() {}

func NewWithTagKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithTagKeyContext {
	var p = new(WithTagKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_withTagKey

	return p
}

func (s *WithTagKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *WithTagKeyContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *WithTagKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithTagKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WithTagKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWithTagKey(s)
	}
}

func (s *WithTagKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWithTagKey(s)
	}
}

func (s *WithTagKeyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWithTagKey(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) WithTagKey() (localctx IWithTagKeyContext) {
	this := p
	_ = this

	localctx = NewWithTagKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, SQLParserRULE_withTagKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(344)
		p.Ident()
	}



	return localctx
}


// INamespaceContext is an interface to support dynamic dispatch.
type INamespaceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNamespaceContext differentiates from other interfaces.
	IsNamespaceContext()
}

type NamespaceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamespaceContext() *NamespaceContext {
	var p = new(NamespaceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_namespace
	return p
}

func (*NamespaceContext) IsNamespaceContext() {}

func NewNamespaceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceContext {
	var p = new(NamespaceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_namespace

	return p
}

func (s *NamespaceContext) GetParser() antlr.Parser { return s.parser }

func (s *NamespaceContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *NamespaceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamespaceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NamespaceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNamespace(s)
	}
}

func (s *NamespaceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNamespace(s)
	}
}

func (s *NamespaceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNamespace(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Namespace() (localctx INamespaceContext) {
	this := p
	_ = this

	localctx = NewNamespaceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SQLParserRULE_namespace)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(346)
		p.Ident()
	}



	return localctx
}


// IDatabaseNameContext is an interface to support dynamic dispatch.
type IDatabaseNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDatabaseNameContext differentiates from other interfaces.
	IsDatabaseNameContext()
}

type DatabaseNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDatabaseNameContext() *DatabaseNameContext {
	var p = new(DatabaseNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_databaseName
	return p
}

func (*DatabaseNameContext) IsDatabaseNameContext() {}

func NewDatabaseNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DatabaseNameContext {
	var p = new(DatabaseNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_databaseName

	return p
}

func (s *DatabaseNameContext) GetParser() antlr.Parser { return s.parser }

func (s *DatabaseNameContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *DatabaseNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DatabaseNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DatabaseNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDatabaseName(s)
	}
}

func (s *DatabaseNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDatabaseName(s)
	}
}

func (s *DatabaseNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDatabaseName(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) DatabaseName() (localctx IDatabaseNameContext) {
	this := p
	_ = this

	localctx = NewDatabaseNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, SQLParserRULE_databaseName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(348)
		p.Ident()
	}



	return localctx
}


// ISourceContext is an interface to support dynamic dispatch.
type ISourceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSourceContext differentiates from other interfaces.
	IsSourceContext()
}

type SourceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySourceContext() *SourceContext {
	var p = new(SourceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_source
	return p
}

func (*SourceContext) IsSourceContext() {}

func NewSourceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceContext {
	var p = new(SourceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_source

	return p
}

func (s *SourceContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceContext) T_STATE_MACHINE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_MACHINE, 0)
}

func (s *SourceContext) T_STATE_REPO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_REPO, 0)
}

func (s *SourceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SourceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSource(s)
	}
}

func (s *SourceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSource(s)
	}
}

func (s *SourceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSource(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Source() (localctx ISourceContext) {
	this := p
	_ = this

	localctx = NewSourceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, SQLParserRULE_source)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(350)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_STATE_REPO || _la == SQLParserT_STATE_MACHINE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


// IQueryStmtContext is an interface to support dynamic dispatch.
type IQueryStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQueryStmtContext differentiates from other interfaces.
	IsQueryStmtContext()
}

type QueryStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryStmtContext() *QueryStmtContext {
	var p = new(QueryStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_queryStmt
	return p
}

func (*QueryStmtContext) IsQueryStmtContext() {}

func NewQueryStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryStmtContext {
	var p = new(QueryStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_queryStmt

	return p
}

func (s *QueryStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryStmtContext) SourceAndSelect() ISourceAndSelectContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISourceAndSelectContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISourceAndSelectContext)
}

func (s *QueryStmtContext) T_EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EXPLAIN, 0)
}

func (s *QueryStmtContext) WhereClause() IWhereClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereClauseContext)
}

func (s *QueryStmtContext) GroupByClause() IGroupByClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByClauseContext)
}

func (s *QueryStmtContext) OrderByClause() IOrderByClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrderByClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrderByClauseContext)
}

func (s *QueryStmtContext) LimitClause() ILimitClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitClauseContext)
}

func (s *QueryStmtContext) T_WITH_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH_VALUE, 0)
}

func (s *QueryStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *QueryStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterQueryStmt(s)
	}
}

func (s *QueryStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitQueryStmt(s)
	}
}

func (s *QueryStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitQueryStmt(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) QueryStmt() (localctx IQueryStmtContext) {
	this := p
	_ = this

	localctx = NewQueryStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, SQLParserRULE_queryStmt)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(353)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_EXPLAIN {
		{
			p.SetState(352)
			p.Match(SQLParserT_EXPLAIN)
		}

	}
	{
		p.SetState(355)
		p.SourceAndSelect()
	}
	p.SetState(357)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WHERE {
		{
			p.SetState(356)
			p.WhereClause()
		}

	}
	p.SetState(360)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_GROUP {
		{
			p.SetState(359)
			p.GroupByClause()
		}

	}
	p.SetState(363)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ORDER {
		{
			p.SetState(362)
			p.OrderByClause()
		}

	}
	p.SetState(366)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_LIMIT {
		{
			p.SetState(365)
			p.LimitClause()
		}

	}
	p.SetState(369)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_WITH_VALUE {
		{
			p.SetState(368)
			p.Match(SQLParserT_WITH_VALUE)
		}

	}



	return localctx
}


// ISourceAndSelectContext is an interface to support dynamic dispatch.
type ISourceAndSelectContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSourceAndSelectContext differentiates from other interfaces.
	IsSourceAndSelectContext()
}

type SourceAndSelectContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySourceAndSelectContext() *SourceAndSelectContext {
	var p = new(SourceAndSelectContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sourceAndSelect
	return p
}

func (*SourceAndSelectContext) IsSourceAndSelectContext() {}

func NewSourceAndSelectContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SourceAndSelectContext {
	var p = new(SourceAndSelectContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sourceAndSelect

	return p
}

func (s *SourceAndSelectContext) GetParser() antlr.Parser { return s.parser }

func (s *SourceAndSelectContext) SelectExpr() ISelectExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectExprContext)
}

func (s *SourceAndSelectContext) FromClause() IFromClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *SourceAndSelectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SourceAndSelectContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SourceAndSelectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSourceAndSelect(s)
	}
}

func (s *SourceAndSelectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSourceAndSelect(s)
	}
}

func (s *SourceAndSelectContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSourceAndSelect(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) SourceAndSelect() (localctx ISourceAndSelectContext) {
	this := p
	_ = this

	localctx = NewSourceAndSelectContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, SQLParserRULE_sourceAndSelect)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(377)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_SELECT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(371)
			p.SelectExpr()
		}
		{
			p.SetState(372)
			p.FromClause()
		}


	case SQLParserT_FROM:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(374)
			p.FromClause()
		}
		{
			p.SetState(375)
			p.SelectExpr()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}


	return localctx
}


// ISelectExprContext is an interface to support dynamic dispatch.
type ISelectExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSelectExprContext differentiates from other interfaces.
	IsSelectExprContext()
}

type SelectExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectExprContext() *SelectExprContext {
	var p = new(SelectExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_selectExpr
	return p
}

func (*SelectExprContext) IsSelectExprContext() {}

func NewSelectExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectExprContext {
	var p = new(SelectExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_selectExpr

	return p
}

func (s *SelectExprContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectExprContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *SelectExprContext) Fields() IFieldsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldsContext)
}

func (s *SelectExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SelectExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSelectExpr(s)
	}
}

func (s *SelectExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSelectExpr(s)
	}
}

func (s *SelectExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSelectExpr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) SelectExpr() (localctx ISelectExprContext) {
	this := p
	_ = this

	localctx = NewSelectExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, SQLParserRULE_selectExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(379)
		p.Match(SQLParserT_SELECT)
	}
	{
		p.SetState(380)
		p.Fields()
	}



	return localctx
}


// IFieldsContext is an interface to support dynamic dispatch.
type IFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldsContext differentiates from other interfaces.
	IsFieldsContext()
}

type FieldsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldsContext() *FieldsContext {
	var p = new(FieldsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fields
	return p
}

func (*FieldsContext) IsFieldsContext() {}

func NewFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldsContext {
	var p = new(FieldsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fields

	return p
}

func (s *FieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldsContext) AllField() []IFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldContext); ok {
			len++
		}
	}

	tst := make([]IFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldContext); ok {
			tst[i] = t.(IFieldContext)
			i++
		}
	}

	return tst
}

func (s *FieldsContext) Field(i int) IFieldContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldContext)
}

func (s *FieldsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *FieldsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *FieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFields(s)
	}
}

func (s *FieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFields(s)
	}
}

func (s *FieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFields(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Fields() (localctx IFieldsContext) {
	this := p
	_ = this

	localctx = NewFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SQLParserRULE_fields)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(382)
		p.Field()
	}
	p.SetState(387)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(383)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(384)
			p.Field()
		}


		p.SetState(389)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IFieldContext is an interface to support dynamic dispatch.
type IFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldContext differentiates from other interfaces.
	IsFieldContext()
}

type FieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldContext() *FieldContext {
	var p = new(FieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_field
	return p
}

func (*FieldContext) IsFieldContext() {}

func NewFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldContext {
	var p = new(FieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_field

	return p
}

func (s *FieldContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldContext) FieldExpr() IFieldExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *FieldContext) Alias() IAliasContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *FieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterField(s)
	}
}

func (s *FieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitField(s)
	}
}

func (s *FieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitField(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Field() (localctx IFieldContext) {
	this := p
	_ = this

	localctx = NewFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SQLParserRULE_field)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(390)
		p.fieldExpr(0)
	}
	p.SetState(392)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_AS {
		{
			p.SetState(391)
			p.Alias()
		}

	}



	return localctx
}


// IAliasContext is an interface to support dynamic dispatch.
type IAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAliasContext differentiates from other interfaces.
	IsAliasContext()
}

type AliasContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAliasContext() *AliasContext {
	var p = new(AliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_alias
	return p
}

func (*AliasContext) IsAliasContext() {}

func NewAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AliasContext {
	var p = new(AliasContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_alias

	return p
}

func (s *AliasContext) GetParser() antlr.Parser { return s.parser }

func (s *AliasContext) T_AS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AS, 0)
}

func (s *AliasContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *AliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *AliasContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterAlias(s)
	}
}

func (s *AliasContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitAlias(s)
	}
}

func (s *AliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitAlias(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Alias() (localctx IAliasContext) {
	this := p
	_ = this

	localctx = NewAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, SQLParserRULE_alias)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(394)
		p.Match(SQLParserT_AS)
	}
	{
		p.SetState(395)
		p.Ident()
	}



	return localctx
}


// IStorageFilterContext is an interface to support dynamic dispatch.
type IStorageFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStorageFilterContext differentiates from other interfaces.
	IsStorageFilterContext()
}

type StorageFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStorageFilterContext() *StorageFilterContext {
	var p = new(StorageFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_storageFilter
	return p
}

func (*StorageFilterContext) IsStorageFilterContext() {}

func NewStorageFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StorageFilterContext {
	var p = new(StorageFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_storageFilter

	return p
}

func (s *StorageFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *StorageFilterContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *StorageFilterContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *StorageFilterContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *StorageFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StorageFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *StorageFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterStorageFilter(s)
	}
}

func (s *StorageFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitStorageFilter(s)
	}
}

func (s *StorageFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitStorageFilter(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) StorageFilter() (localctx IStorageFilterContext) {
	this := p
	_ = this

	localctx = NewStorageFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, SQLParserRULE_storageFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(397)
		p.Match(SQLParserT_STORAGE)
	}
	{
		p.SetState(398)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(399)
		p.Ident()
	}



	return localctx
}


// IDatabaseFilterContext is an interface to support dynamic dispatch.
type IDatabaseFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDatabaseFilterContext differentiates from other interfaces.
	IsDatabaseFilterContext()
}

type DatabaseFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDatabaseFilterContext() *DatabaseFilterContext {
	var p = new(DatabaseFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_databaseFilter
	return p
}

func (*DatabaseFilterContext) IsDatabaseFilterContext() {}

func NewDatabaseFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DatabaseFilterContext {
	var p = new(DatabaseFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_databaseFilter

	return p
}

func (s *DatabaseFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *DatabaseFilterContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *DatabaseFilterContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *DatabaseFilterContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *DatabaseFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DatabaseFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DatabaseFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDatabaseFilter(s)
	}
}

func (s *DatabaseFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDatabaseFilter(s)
	}
}

func (s *DatabaseFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDatabaseFilter(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) DatabaseFilter() (localctx IDatabaseFilterContext) {
	this := p
	_ = this

	localctx = NewDatabaseFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SQLParserRULE_databaseFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(401)
		p.Match(SQLParserT_DATASBAE)
	}
	{
		p.SetState(402)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(403)
		p.Ident()
	}



	return localctx
}


// ITypeFilterContext is an interface to support dynamic dispatch.
type ITypeFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTypeFilterContext differentiates from other interfaces.
	IsTypeFilterContext()
}

type TypeFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeFilterContext() *TypeFilterContext {
	var p = new(TypeFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_typeFilter
	return p
}

func (*TypeFilterContext) IsTypeFilterContext() {}

func NewTypeFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeFilterContext {
	var p = new(TypeFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_typeFilter

	return p
}

func (s *TypeFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeFilterContext) T_TYPE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPE, 0)
}

func (s *TypeFilterContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *TypeFilterContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TypeFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TypeFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTypeFilter(s)
	}
}

func (s *TypeFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTypeFilter(s)
	}
}

func (s *TypeFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTypeFilter(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) TypeFilter() (localctx ITypeFilterContext) {
	this := p
	_ = this

	localctx = NewTypeFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, SQLParserRULE_typeFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(405)
		p.Match(SQLParserT_TYPE)
	}
	{
		p.SetState(406)
		p.Match(SQLParserT_EQUAL)
	}
	{
		p.SetState(407)
		p.Ident()
	}



	return localctx
}


// IFromClauseContext is an interface to support dynamic dispatch.
type IFromClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFromClauseContext differentiates from other interfaces.
	IsFromClauseContext()
}

type FromClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromClauseContext() *FromClauseContext {
	var p = new(FromClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fromClause
	return p
}

func (*FromClauseContext) IsFromClauseContext() {}

func NewFromClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromClauseContext {
	var p = new(FromClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fromClause

	return p
}

func (s *FromClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FromClauseContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *FromClauseContext) MetricName() IMetricNameContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMetricNameContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMetricNameContext)
}

func (s *FromClauseContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *FromClauseContext) Namespace() INamespaceContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamespaceContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamespaceContext)
}

func (s *FromClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FromClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFromClause(s)
	}
}

func (s *FromClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFromClause(s)
	}
}

func (s *FromClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFromClause(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) FromClause() (localctx IFromClauseContext) {
	this := p
	_ = this

	localctx = NewFromClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, SQLParserRULE_fromClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(409)
		p.Match(SQLParserT_FROM)
	}
	{
		p.SetState(410)
		p.MetricName()
	}
	p.SetState(413)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ON {
		{
			p.SetState(411)
			p.Match(SQLParserT_ON)
		}
		{
			p.SetState(412)
			p.Namespace()
		}

	}



	return localctx
}


// IWhereClauseContext is an interface to support dynamic dispatch.
type IWhereClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhereClauseContext differentiates from other interfaces.
	IsWhereClauseContext()
}

type WhereClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereClauseContext() *WhereClauseContext {
	var p = new(WhereClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_whereClause
	return p
}

func (*WhereClauseContext) IsWhereClauseContext() {}

func NewWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereClauseContext {
	var p = new(WhereClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_whereClause

	return p
}

func (s *WhereClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereClauseContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *WhereClauseContext) ConditionExpr() IConditionExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionExprContext)
}

func (s *WhereClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *WhereClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterWhereClause(s)
	}
}

func (s *WhereClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitWhereClause(s)
	}
}

func (s *WhereClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitWhereClause(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) WhereClause() (localctx IWhereClauseContext) {
	this := p
	_ = this

	localctx = NewWhereClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, SQLParserRULE_whereClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(415)
		p.Match(SQLParserT_WHERE)
	}
	{
		p.SetState(416)
		p.ConditionExpr()
	}



	return localctx
}


// IConditionExprContext is an interface to support dynamic dispatch.
type IConditionExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsConditionExprContext differentiates from other interfaces.
	IsConditionExprContext()
}

type ConditionExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConditionExprContext() *ConditionExprContext {
	var p = new(ConditionExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_conditionExpr
	return p
}

func (*ConditionExprContext) IsConditionExprContext() {}

func NewConditionExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionExprContext {
	var p = new(ConditionExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_conditionExpr

	return p
}

func (s *ConditionExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionExprContext) TagFilterExpr() ITagFilterExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagFilterExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
}

func (s *ConditionExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *ConditionExprContext) TimeRangeExpr() ITimeRangeExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeRangeExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeRangeExprContext)
}

func (s *ConditionExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ConditionExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterConditionExpr(s)
	}
}

func (s *ConditionExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitConditionExpr(s)
	}
}

func (s *ConditionExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitConditionExpr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ConditionExpr() (localctx IConditionExprContext) {
	this := p
	_ = this

	localctx = NewConditionExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, SQLParserRULE_conditionExpr)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(428)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(418)
			p.tagFilterExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(419)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(420)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(421)
			p.TimeRangeExpr()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(423)
			p.TimeRangeExpr()
		}
		p.SetState(426)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		if _la == SQLParserT_AND {
			{
				p.SetState(424)
				p.Match(SQLParserT_AND)
			}
			{
				p.SetState(425)
				p.tagFilterExpr(0)
			}

		}

	}


	return localctx
}


// ITagFilterExprContext is an interface to support dynamic dispatch.
type ITagFilterExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagFilterExprContext differentiates from other interfaces.
	IsTagFilterExprContext()
}

type TagFilterExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagFilterExprContext() *TagFilterExprContext {
	var p = new(TagFilterExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagFilterExpr
	return p
}

func (*TagFilterExprContext) IsTagFilterExprContext() {}

func NewTagFilterExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagFilterExprContext {
	var p = new(TagFilterExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagFilterExpr

	return p
}

func (s *TagFilterExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TagFilterExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *TagFilterExprContext) AllTagFilterExpr() []ITagFilterExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITagFilterExprContext); ok {
			len++
		}
	}

	tst := make([]ITagFilterExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITagFilterExprContext); ok {
			tst[i] = t.(ITagFilterExprContext)
			i++
		}
	}

	return tst
}

func (s *TagFilterExprContext) TagFilterExpr(i int) ITagFilterExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagFilterExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
}

func (s *TagFilterExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *TagFilterExprContext) TagKey() ITagKeyContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagKeyContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagKeyContext)
}

func (s *TagFilterExprContext) TagValue() ITagValueContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagValueContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagValueContext)
}

func (s *TagFilterExprContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *TagFilterExprContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *TagFilterExprContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *TagFilterExprContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *TagFilterExprContext) T_NEQREGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NEQREGEXP, 0)
}

func (s *TagFilterExprContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *TagFilterExprContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *TagFilterExprContext) TagValueList() ITagValueListContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagValueListContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagValueListContext)
}

func (s *TagFilterExprContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *TagFilterExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TagFilterExprContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *TagFilterExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagFilterExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TagFilterExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagFilterExpr(s)
	}
}

func (s *TagFilterExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagFilterExpr(s)
	}
}

func (s *TagFilterExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagFilterExpr(s)

	default:
		return t.VisitChildren(s)
	}
}





func (p *SQLParser) TagFilterExpr() (localctx ITagFilterExprContext) {
	return p.tagFilterExpr(0)
}

func (p *SQLParser) tagFilterExpr(_p int) (localctx ITagFilterExprContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewTagFilterExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ITagFilterExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 80
	p.EnterRecursionRule(localctx, 80, SQLParserRULE_tagFilterExpr, _p)
	var _la int


	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(458)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(431)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(432)
			p.tagFilterExpr(0)
		}
		{
			p.SetState(433)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(435)
			p.TagKey()
		}
		p.SetState(444)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_EQUAL:
			{
				p.SetState(436)
				p.Match(SQLParserT_EQUAL)
			}


		case SQLParserT_LIKE:
			{
				p.SetState(437)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_NOT:
			{
				p.SetState(438)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(439)
				p.Match(SQLParserT_LIKE)
			}


		case SQLParserT_REGEXP:
			{
				p.SetState(440)
				p.Match(SQLParserT_REGEXP)
			}


		case SQLParserT_NEQREGEXP:
			{
				p.SetState(441)
				p.Match(SQLParserT_NEQREGEXP)
			}


		case SQLParserT_NOTEQUAL:
			{
				p.SetState(442)
				p.Match(SQLParserT_NOTEQUAL)
			}


		case SQLParserT_NOTEQUAL2:
			{
				p.SetState(443)
				p.Match(SQLParserT_NOTEQUAL2)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(446)
			p.TagValue()
		}


	case 3:
		{
			p.SetState(448)
			p.TagKey()
		}
		p.SetState(452)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SQLParserT_IN:
			{
				p.SetState(449)
				p.Match(SQLParserT_IN)
			}


		case SQLParserT_NOT:
			{
				p.SetState(450)
				p.Match(SQLParserT_NOT)
			}
			{
				p.SetState(451)
				p.Match(SQLParserT_IN)
			}



		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(454)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(455)
			p.TagValueList()
		}
		{
			p.SetState(456)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(465)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTagFilterExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_tagFilterExpr)
			p.SetState(460)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(461)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(462)
				p.tagFilterExpr(2)
			}


		}
		p.SetState(467)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext())
	}



	return localctx
}


// ITagValueListContext is an interface to support dynamic dispatch.
type ITagValueListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagValueListContext differentiates from other interfaces.
	IsTagValueListContext()
}

type TagValueListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagValueListContext() *TagValueListContext {
	var p = new(TagValueListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagValueList
	return p
}

func (*TagValueListContext) IsTagValueListContext() {}

func NewTagValueListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagValueListContext {
	var p = new(TagValueListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagValueList

	return p
}

func (s *TagValueListContext) GetParser() antlr.Parser { return s.parser }

func (s *TagValueListContext) AllTagValue() []ITagValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITagValueContext); ok {
			len++
		}
	}

	tst := make([]ITagValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITagValueContext); ok {
			tst[i] = t.(ITagValueContext)
			i++
		}
	}

	return tst
}

func (s *TagValueListContext) TagValue(i int) ITagValueContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagValueContext)
}

func (s *TagValueListContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *TagValueListContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *TagValueListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagValueListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TagValueListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagValueList(s)
	}
}

func (s *TagValueListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagValueList(s)
	}
}

func (s *TagValueListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagValueList(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) TagValueList() (localctx ITagValueListContext) {
	this := p
	_ = this

	localctx = NewTagValueListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, SQLParserRULE_tagValueList)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(468)
		p.TagValue()
	}
	p.SetState(473)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(469)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(470)
			p.TagValue()
		}


		p.SetState(475)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IMetricListFilterContext is an interface to support dynamic dispatch.
type IMetricListFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetricListFilterContext differentiates from other interfaces.
	IsMetricListFilterContext()
}

type MetricListFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricListFilterContext() *MetricListFilterContext {
	var p = new(MetricListFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metricListFilter
	return p
}

func (*MetricListFilterContext) IsMetricListFilterContext() {}

func NewMetricListFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricListFilterContext {
	var p = new(MetricListFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metricListFilter

	return p
}

func (s *MetricListFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricListFilterContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *MetricListFilterContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *MetricListFilterContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *MetricListFilterContext) MetricList() IMetricListContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMetricListContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMetricListContext)
}

func (s *MetricListFilterContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *MetricListFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricListFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MetricListFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetricListFilter(s)
	}
}

func (s *MetricListFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetricListFilter(s)
	}
}

func (s *MetricListFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetricListFilter(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) MetricListFilter() (localctx IMetricListFilterContext) {
	this := p
	_ = this

	localctx = NewMetricListFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, SQLParserRULE_metricListFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(476)
		p.Match(SQLParserT_METRIC)
	}
	{
		p.SetState(477)
		p.Match(SQLParserT_IN)
	}

	{
		p.SetState(478)
		p.Match(SQLParserT_OPEN_P)
	}
	{
		p.SetState(479)
		p.MetricList()
	}
	{
		p.SetState(480)
		p.Match(SQLParserT_CLOSE_P)
	}




	return localctx
}


// IMetricListContext is an interface to support dynamic dispatch.
type IMetricListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetricListContext differentiates from other interfaces.
	IsMetricListContext()
}

type MetricListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricListContext() *MetricListContext {
	var p = new(MetricListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metricList
	return p
}

func (*MetricListContext) IsMetricListContext() {}

func NewMetricListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricListContext {
	var p = new(MetricListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metricList

	return p
}

func (s *MetricListContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricListContext) AllIdent() []IIdentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentContext); ok {
			len++
		}
	}

	tst := make([]IIdentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentContext); ok {
			tst[i] = t.(IIdentContext)
			i++
		}
	}

	return tst
}

func (s *MetricListContext) Ident(i int) IIdentContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *MetricListContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *MetricListContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *MetricListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MetricListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetricList(s)
	}
}

func (s *MetricListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetricList(s)
	}
}

func (s *MetricListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetricList(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) MetricList() (localctx IMetricListContext) {
	this := p
	_ = this

	localctx = NewMetricListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SQLParserRULE_metricList)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(482)
		p.Ident()
	}
	p.SetState(487)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(483)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(484)
			p.Ident()
		}


		p.SetState(489)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// ITimeRangeExprContext is an interface to support dynamic dispatch.
type ITimeRangeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeRangeExprContext differentiates from other interfaces.
	IsTimeRangeExprContext()
}

type TimeRangeExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeRangeExprContext() *TimeRangeExprContext {
	var p = new(TimeRangeExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_timeRangeExpr
	return p
}

func (*TimeRangeExprContext) IsTimeRangeExprContext() {}

func NewTimeRangeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeRangeExprContext {
	var p = new(TimeRangeExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_timeRangeExpr

	return p
}

func (s *TimeRangeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeRangeExprContext) AllTimeExpr() []ITimeExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITimeExprContext); ok {
			len++
		}
	}

	tst := make([]ITimeExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITimeExprContext); ok {
			tst[i] = t.(ITimeExprContext)
			i++
		}
	}

	return tst
}

func (s *TimeRangeExprContext) TimeExpr(i int) ITimeExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITimeExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITimeExprContext)
}

func (s *TimeRangeExprContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *TimeRangeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeRangeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TimeRangeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTimeRangeExpr(s)
	}
}

func (s *TimeRangeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTimeRangeExpr(s)
	}
}

func (s *TimeRangeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTimeRangeExpr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) TimeRangeExpr() (localctx ITimeRangeExprContext) {
	this := p
	_ = this

	localctx = NewTimeRangeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, SQLParserRULE_timeRangeExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(490)
		p.TimeExpr()
	}
	p.SetState(493)
	p.GetErrorHandler().Sync(p)


	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(491)
			p.Match(SQLParserT_AND)
		}
		{
			p.SetState(492)
			p.TimeExpr()
		}


	}



	return localctx
}


// ITimeExprContext is an interface to support dynamic dispatch.
type ITimeExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTimeExprContext differentiates from other interfaces.
	IsTimeExprContext()
}

type TimeExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTimeExprContext() *TimeExprContext {
	var p = new(TimeExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_timeExpr
	return p
}

func (*TimeExprContext) IsTimeExprContext() {}

func NewTimeExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TimeExprContext {
	var p = new(TimeExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_timeExpr

	return p
}

func (s *TimeExprContext) GetParser() antlr.Parser { return s.parser }

func (s *TimeExprContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *TimeExprContext) BinaryOperator() IBinaryOperatorContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBinaryOperatorContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBinaryOperatorContext)
}

func (s *TimeExprContext) NowExpr() INowExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INowExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INowExprContext)
}

func (s *TimeExprContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TimeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TimeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTimeExpr(s)
	}
}

func (s *TimeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTimeExpr(s)
	}
}

func (s *TimeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTimeExpr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) TimeExpr() (localctx ITimeExprContext) {
	this := p
	_ = this

	localctx = NewTimeExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, SQLParserRULE_timeExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(495)
		p.Match(SQLParserT_TIME)
	}
	{
		p.SetState(496)
		p.BinaryOperator()
	}
	p.SetState(499)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(497)
			p.NowExpr()
		}


	case 2:
		{
			p.SetState(498)
			p.Ident()
		}

	}



	return localctx
}


// INowExprContext is an interface to support dynamic dispatch.
type INowExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNowExprContext differentiates from other interfaces.
	IsNowExprContext()
}

type NowExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNowExprContext() *NowExprContext {
	var p = new(NowExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_nowExpr
	return p
}

func (*NowExprContext) IsNowExprContext() {}

func NewNowExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NowExprContext {
	var p = new(NowExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nowExpr

	return p
}

func (s *NowExprContext) GetParser() antlr.Parser { return s.parser }

func (s *NowExprContext) NowFunc() INowFuncContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INowFuncContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INowFuncContext)
}

func (s *NowExprContext) DurationLit() IDurationLitContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDurationLitContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *NowExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NowExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NowExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNowExpr(s)
	}
}

func (s *NowExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNowExpr(s)
	}
}

func (s *NowExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNowExpr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) NowExpr() (localctx INowExprContext) {
	this := p
	_ = this

	localctx = NewNowExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, SQLParserRULE_nowExpr)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(501)
		p.NowFunc()
	}
	p.SetState(503)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if ((((_la - 115)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 115))) & ((1 << (SQLParserT_ADD - 115)) | (1 << (SQLParserT_SUB - 115)) | (1 << (SQLParserL_INT - 115)))) != 0) {
		{
			p.SetState(502)
			p.DurationLit()
		}

	}



	return localctx
}


// INowFuncContext is an interface to support dynamic dispatch.
type INowFuncContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNowFuncContext differentiates from other interfaces.
	IsNowFuncContext()
}

type NowFuncContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNowFuncContext() *NowFuncContext {
	var p = new(NowFuncContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_nowFunc
	return p
}

func (*NowFuncContext) IsNowFuncContext() {}

func NewNowFuncContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NowFuncContext {
	var p = new(NowFuncContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nowFunc

	return p
}

func (s *NowFuncContext) GetParser() antlr.Parser { return s.parser }

func (s *NowFuncContext) T_NOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOW, 0)
}

func (s *NowFuncContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *NowFuncContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *NowFuncContext) ExprFuncParams() IExprFuncParamsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprFuncParamsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprFuncParamsContext)
}

func (s *NowFuncContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NowFuncContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NowFuncContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNowFunc(s)
	}
}

func (s *NowFuncContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNowFunc(s)
	}
}

func (s *NowFuncContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNowFunc(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) NowFunc() (localctx INowFuncContext) {
	this := p
	_ = this

	localctx = NewNowFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, SQLParserRULE_nowFunc)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(505)
		p.Match(SQLParserT_NOW)
	}
	{
		p.SetState(506)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(508)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_STATE_REPO) | (1 << SQLParserT_STATE_MACHINE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_METADATA) | (1 << SQLParserT_TYPES) | (1 << SQLParserT_TYPE) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_SCHEMAS - 32)) | (1 << (SQLParserT_DATASBAE - 32)) | (1 << (SQLParserT_DATASBAES - 32)) | (1 << (SQLParserT_NAMESPACE - 32)) | (1 << (SQLParserT_NAMESPACES - 32)) | (1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_ASC - 64)) | (1 << (SQLParserT_DESC - 64)) | (1 << (SQLParserT_LIKE - 64)) | (1 << (SQLParserT_NOT - 64)) | (1 << (SQLParserT_BETWEEN - 64)) | (1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_LAST - 64)) | (1 << (SQLParserT_FIRST - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_RATE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)))) != 0) || ((((_la - 96)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 96))) & ((1 << (SQLParserT_YEAR - 96)) | (1 << (SQLParserT_OPEN_P - 96)) | (1 << (SQLParserT_ADD - 96)) | (1 << (SQLParserT_SUB - 96)) | (1 << (SQLParserL_ID - 96)) | (1 << (SQLParserL_INT - 96)) | (1 << (SQLParserL_DEC - 96)))) != 0) {
		{
			p.SetState(507)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(510)
		p.Match(SQLParserT_CLOSE_P)
	}



	return localctx
}


// IGroupByClauseContext is an interface to support dynamic dispatch.
type IGroupByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByClauseContext differentiates from other interfaces.
	IsGroupByClauseContext()
}

type GroupByClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByClauseContext() *GroupByClauseContext {
	var p = new(GroupByClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_groupByClause
	return p
}

func (*GroupByClauseContext) IsGroupByClauseContext() {}

func NewGroupByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByClauseContext {
	var p = new(GroupByClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupByClause

	return p
}

func (s *GroupByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByClauseContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *GroupByClauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *GroupByClauseContext) GroupByKeys() IGroupByKeysContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByKeysContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByKeysContext)
}

func (s *GroupByClauseContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *GroupByClauseContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *GroupByClauseContext) FillOption() IFillOptionContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFillOptionContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFillOptionContext)
}

func (s *GroupByClauseContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *GroupByClauseContext) HavingClause() IHavingClauseContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHavingClauseContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHavingClauseContext)
}

func (s *GroupByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GroupByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroupByClause(s)
	}
}

func (s *GroupByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroupByClause(s)
	}
}

func (s *GroupByClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitGroupByClause(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) GroupByClause() (localctx IGroupByClauseContext) {
	this := p
	_ = this

	localctx = NewGroupByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, SQLParserRULE_groupByClause)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(512)
		p.Match(SQLParserT_GROUP)
	}
	{
		p.SetState(513)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(514)
		p.GroupByKeys()
	}
	p.SetState(520)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_FILL {
		{
			p.SetState(515)
			p.Match(SQLParserT_FILL)
		}
		{
			p.SetState(516)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(517)
			p.FillOption()
		}
		{
			p.SetState(518)
			p.Match(SQLParserT_CLOSE_P)
		}

	}
	p.SetState(523)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_HAVING {
		{
			p.SetState(522)
			p.HavingClause()
		}

	}



	return localctx
}


// IGroupByKeysContext is an interface to support dynamic dispatch.
type IGroupByKeysContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByKeysContext differentiates from other interfaces.
	IsGroupByKeysContext()
}

type GroupByKeysContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByKeysContext() *GroupByKeysContext {
	var p = new(GroupByKeysContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_groupByKeys
	return p
}

func (*GroupByKeysContext) IsGroupByKeysContext() {}

func NewGroupByKeysContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByKeysContext {
	var p = new(GroupByKeysContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupByKeys

	return p
}

func (s *GroupByKeysContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByKeysContext) AllGroupByKey() []IGroupByKeyContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IGroupByKeyContext); ok {
			len++
		}
	}

	tst := make([]IGroupByKeyContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IGroupByKeyContext); ok {
			tst[i] = t.(IGroupByKeyContext)
			i++
		}
	}

	return tst
}

func (s *GroupByKeysContext) GroupByKey(i int) IGroupByKeyContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByKeyContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByKeyContext)
}

func (s *GroupByKeysContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *GroupByKeysContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *GroupByKeysContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByKeysContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GroupByKeysContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroupByKeys(s)
	}
}

func (s *GroupByKeysContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroupByKeys(s)
	}
}

func (s *GroupByKeysContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitGroupByKeys(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) GroupByKeys() (localctx IGroupByKeysContext) {
	this := p
	_ = this

	localctx = NewGroupByKeysContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, SQLParserRULE_groupByKeys)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(525)
		p.GroupByKey()
	}
	p.SetState(530)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(526)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(527)
			p.GroupByKey()
		}


		p.SetState(532)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IGroupByKeyContext is an interface to support dynamic dispatch.
type IGroupByKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByKeyContext differentiates from other interfaces.
	IsGroupByKeyContext()
}

type GroupByKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByKeyContext() *GroupByKeyContext {
	var p = new(GroupByKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_groupByKey
	return p
}

func (*GroupByKeyContext) IsGroupByKeyContext() {}

func NewGroupByKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByKeyContext {
	var p = new(GroupByKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupByKey

	return p
}

func (s *GroupByKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByKeyContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *GroupByKeyContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *GroupByKeyContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *GroupByKeyContext) DurationLit() IDurationLitContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDurationLitContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *GroupByKeyContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *GroupByKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *GroupByKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterGroupByKey(s)
	}
}

func (s *GroupByKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitGroupByKey(s)
	}
}

func (s *GroupByKeyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitGroupByKey(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) GroupByKey() (localctx IGroupByKeyContext) {
	this := p
	_ = this

	localctx = NewGroupByKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, SQLParserRULE_groupByKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(539)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(533)
			p.Ident()
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(534)
			p.Match(SQLParserT_TIME)
		}
		{
			p.SetState(535)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(536)
			p.DurationLit()
		}
		{
			p.SetState(537)
			p.Match(SQLParserT_CLOSE_P)
		}

	}


	return localctx
}


// IFillOptionContext is an interface to support dynamic dispatch.
type IFillOptionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFillOptionContext differentiates from other interfaces.
	IsFillOptionContext()
}

type FillOptionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillOptionContext() *FillOptionContext {
	var p = new(FillOptionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fillOption
	return p
}

func (*FillOptionContext) IsFillOptionContext() {}

func NewFillOptionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillOptionContext {
	var p = new(FillOptionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fillOption

	return p
}

func (s *FillOptionContext) GetParser() antlr.Parser { return s.parser }

func (s *FillOptionContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *FillOptionContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
}

func (s *FillOptionContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *FillOptionContext) L_DEC() antlr.TerminalNode {
	return s.GetToken(SQLParserL_DEC, 0)
}

func (s *FillOptionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillOptionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FillOptionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFillOption(s)
	}
}

func (s *FillOptionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFillOption(s)
	}
}

func (s *FillOptionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFillOption(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) FillOption() (localctx IFillOptionContext) {
	this := p
	_ = this

	localctx = NewFillOptionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, SQLParserRULE_fillOption)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(541)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_NULL || _la == SQLParserT_PREVIOUS || _la == SQLParserL_INT || _la == SQLParserL_DEC) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


// IOrderByClauseContext is an interface to support dynamic dispatch.
type IOrderByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOrderByClauseContext differentiates from other interfaces.
	IsOrderByClauseContext()
}

type OrderByClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrderByClauseContext() *OrderByClauseContext {
	var p = new(OrderByClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_orderByClause
	return p
}

func (*OrderByClauseContext) IsOrderByClauseContext() {}

func NewOrderByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByClauseContext {
	var p = new(OrderByClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_orderByClause

	return p
}

func (s *OrderByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByClauseContext) T_ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ORDER, 0)
}

func (s *OrderByClauseContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *OrderByClauseContext) SortFields() ISortFieldsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortFieldsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortFieldsContext)
}

func (s *OrderByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *OrderByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterOrderByClause(s)
	}
}

func (s *OrderByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitOrderByClause(s)
	}
}

func (s *OrderByClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitOrderByClause(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) OrderByClause() (localctx IOrderByClauseContext) {
	this := p
	_ = this

	localctx = NewOrderByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, SQLParserRULE_orderByClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(543)
		p.Match(SQLParserT_ORDER)
	}
	{
		p.SetState(544)
		p.Match(SQLParserT_BY)
	}
	{
		p.SetState(545)
		p.SortFields()
	}



	return localctx
}


// ISortFieldContext is an interface to support dynamic dispatch.
type ISortFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortFieldContext differentiates from other interfaces.
	IsSortFieldContext()
}

type SortFieldContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortFieldContext() *SortFieldContext {
	var p = new(SortFieldContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sortField
	return p
}

func (*SortFieldContext) IsSortFieldContext() {}

func NewSortFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortFieldContext {
	var p = new(SortFieldContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sortField

	return p
}

func (s *SortFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *SortFieldContext) FieldExpr() IFieldExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *SortFieldContext) AllT_ASC() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_ASC)
}

func (s *SortFieldContext) T_ASC(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_ASC, i)
}

func (s *SortFieldContext) AllT_DESC() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_DESC)
}

func (s *SortFieldContext) T_DESC(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_DESC, i)
}

func (s *SortFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SortFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSortField(s)
	}
}

func (s *SortFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSortField(s)
	}
}

func (s *SortFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSortField(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) SortField() (localctx ISortFieldContext) {
	this := p
	_ = this

	localctx = NewSortFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, SQLParserRULE_sortField)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(547)
		p.fieldExpr(0)
	}
	p.SetState(551)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_ASC || _la == SQLParserT_DESC {
		{
			p.SetState(548)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ASC || _la == SQLParserT_DESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}


		p.SetState(553)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// ISortFieldsContext is an interface to support dynamic dispatch.
type ISortFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSortFieldsContext differentiates from other interfaces.
	IsSortFieldsContext()
}

type SortFieldsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySortFieldsContext() *SortFieldsContext {
	var p = new(SortFieldsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_sortFields
	return p
}

func (*SortFieldsContext) IsSortFieldsContext() {}

func NewSortFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortFieldsContext {
	var p = new(SortFieldsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sortFields

	return p
}

func (s *SortFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *SortFieldsContext) AllSortField() []ISortFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISortFieldContext); ok {
			len++
		}
	}

	tst := make([]ISortFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISortFieldContext); ok {
			tst[i] = t.(ISortFieldContext)
			i++
		}
	}

	return tst
}

func (s *SortFieldsContext) SortField(i int) ISortFieldContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortFieldContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortFieldContext)
}

func (s *SortFieldsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *SortFieldsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *SortFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *SortFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterSortFields(s)
	}
}

func (s *SortFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitSortFields(s)
	}
}

func (s *SortFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitSortFields(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) SortFields() (localctx ISortFieldsContext) {
	this := p
	_ = this

	localctx = NewSortFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, SQLParserRULE_sortFields)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(554)
		p.SortField()
	}
	p.SetState(559)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(555)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(556)
			p.SortField()
		}


		p.SetState(561)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IHavingClauseContext is an interface to support dynamic dispatch.
type IHavingClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHavingClauseContext differentiates from other interfaces.
	IsHavingClauseContext()
}

type HavingClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHavingClauseContext() *HavingClauseContext {
	var p = new(HavingClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_havingClause
	return p
}

func (*HavingClauseContext) IsHavingClauseContext() {}

func NewHavingClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HavingClauseContext {
	var p = new(HavingClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_havingClause

	return p
}

func (s *HavingClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *HavingClauseContext) T_HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HAVING, 0)
}

func (s *HavingClauseContext) BoolExpr() IBoolExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoolExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBoolExprContext)
}

func (s *HavingClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HavingClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *HavingClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterHavingClause(s)
	}
}

func (s *HavingClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitHavingClause(s)
	}
}

func (s *HavingClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitHavingClause(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) HavingClause() (localctx IHavingClauseContext) {
	this := p
	_ = this

	localctx = NewHavingClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, SQLParserRULE_havingClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(562)
		p.Match(SQLParserT_HAVING)
	}
	{
		p.SetState(563)
		p.boolExpr(0)
	}



	return localctx
}


// IBoolExprContext is an interface to support dynamic dispatch.
type IBoolExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprContext differentiates from other interfaces.
	IsBoolExprContext()
}

type BoolExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprContext() *BoolExprContext {
	var p = new(BoolExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExpr
	return p
}

func (*BoolExprContext) IsBoolExprContext() {}

func NewBoolExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprContext {
	var p = new(BoolExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExpr

	return p
}

func (s *BoolExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *BoolExprContext) AllBoolExpr() []IBoolExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBoolExprContext); ok {
			len++
		}
	}

	tst := make([]IBoolExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBoolExprContext); ok {
			tst[i] = t.(IBoolExprContext)
			i++
		}
	}

	return tst
}

func (s *BoolExprContext) BoolExpr(i int) IBoolExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoolExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBoolExprContext)
}

func (s *BoolExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *BoolExprContext) BoolExprAtom() IBoolExprAtomContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoolExprAtomContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBoolExprAtomContext)
}

func (s *BoolExprContext) BoolExprLogicalOp() IBoolExprLogicalOpContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoolExprLogicalOpContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBoolExprLogicalOpContext)
}

func (s *BoolExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BoolExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExpr(s)
	}
}

func (s *BoolExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExpr(s)
	}
}

func (s *BoolExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExpr(s)

	default:
		return t.VisitChildren(s)
	}
}





func (p *SQLParser) BoolExpr() (localctx IBoolExprContext) {
	return p.boolExpr(0)
}

func (p *SQLParser) boolExpr(_p int) (localctx IBoolExprContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewBoolExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IBoolExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 112
	p.EnterRecursionRule(localctx, 112, SQLParserRULE_boolExpr, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(571)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(566)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(567)
			p.boolExpr(0)
		}
		{
			p.SetState(568)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(570)
			p.BoolExprAtom()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(579)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewBoolExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_boolExpr)
			p.SetState(573)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(574)
				p.BoolExprLogicalOp()
			}
			{
				p.SetState(575)
				p.boolExpr(3)
			}


		}
		p.SetState(581)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext())
	}



	return localctx
}


// IBoolExprLogicalOpContext is an interface to support dynamic dispatch.
type IBoolExprLogicalOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprLogicalOpContext differentiates from other interfaces.
	IsBoolExprLogicalOpContext()
}

type BoolExprLogicalOpContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprLogicalOpContext() *BoolExprLogicalOpContext {
	var p = new(BoolExprLogicalOpContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExprLogicalOp
	return p
}

func (*BoolExprLogicalOpContext) IsBoolExprLogicalOpContext() {}

func NewBoolExprLogicalOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprLogicalOpContext {
	var p = new(BoolExprLogicalOpContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExprLogicalOp

	return p
}

func (s *BoolExprLogicalOpContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprLogicalOpContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *BoolExprLogicalOpContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *BoolExprLogicalOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprLogicalOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BoolExprLogicalOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExprLogicalOp(s)
	}
}

func (s *BoolExprLogicalOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExprLogicalOp(s)
	}
}

func (s *BoolExprLogicalOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExprLogicalOp(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) BoolExprLogicalOp() (localctx IBoolExprLogicalOpContext) {
	this := p
	_ = this

	localctx = NewBoolExprLogicalOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SQLParserRULE_boolExprLogicalOp)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(582)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserT_AND || _la == SQLParserT_OR) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


// IBoolExprAtomContext is an interface to support dynamic dispatch.
type IBoolExprAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBoolExprAtomContext differentiates from other interfaces.
	IsBoolExprAtomContext()
}

type BoolExprAtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolExprAtomContext() *BoolExprAtomContext {
	var p = new(BoolExprAtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_boolExprAtom
	return p
}

func (*BoolExprAtomContext) IsBoolExprAtomContext() {}

func NewBoolExprAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolExprAtomContext {
	var p = new(BoolExprAtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_boolExprAtom

	return p
}

func (s *BoolExprAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolExprAtomContext) BinaryExpr() IBinaryExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBinaryExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBinaryExprContext)
}

func (s *BoolExprAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolExprAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BoolExprAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBoolExprAtom(s)
	}
}

func (s *BoolExprAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBoolExprAtom(s)
	}
}

func (s *BoolExprAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBoolExprAtom(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) BoolExprAtom() (localctx IBoolExprAtomContext) {
	this := p
	_ = this

	localctx = NewBoolExprAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, SQLParserRULE_boolExprAtom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(584)
		p.BinaryExpr()
	}



	return localctx
}


// IBinaryExprContext is an interface to support dynamic dispatch.
type IBinaryExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinaryExprContext differentiates from other interfaces.
	IsBinaryExprContext()
}

type BinaryExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryExprContext() *BinaryExprContext {
	var p = new(BinaryExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_binaryExpr
	return p
}

func (*BinaryExprContext) IsBinaryExprContext() {}

func NewBinaryExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryExprContext {
	var p = new(BinaryExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_binaryExpr

	return p
}

func (s *BinaryExprContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryExprContext) AllFieldExpr() []IFieldExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldExprContext); ok {
			len++
		}
	}

	tst := make([]IFieldExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldExprContext); ok {
			tst[i] = t.(IFieldExprContext)
			i++
		}
	}

	return tst
}

func (s *BinaryExprContext) FieldExpr(i int) IFieldExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *BinaryExprContext) BinaryOperator() IBinaryOperatorContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBinaryOperatorContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBinaryOperatorContext)
}

func (s *BinaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BinaryExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBinaryExpr(s)
	}
}

func (s *BinaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBinaryExpr(s)
	}
}

func (s *BinaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBinaryExpr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) BinaryExpr() (localctx IBinaryExprContext) {
	this := p
	_ = this

	localctx = NewBinaryExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 118, SQLParserRULE_binaryExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(586)
		p.fieldExpr(0)
	}
	{
		p.SetState(587)
		p.BinaryOperator()
	}
	{
		p.SetState(588)
		p.fieldExpr(0)
	}



	return localctx
}


// IBinaryOperatorContext is an interface to support dynamic dispatch.
type IBinaryOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBinaryOperatorContext differentiates from other interfaces.
	IsBinaryOperatorContext()
}

type BinaryOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBinaryOperatorContext() *BinaryOperatorContext {
	var p = new(BinaryOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_binaryOperator
	return p
}

func (*BinaryOperatorContext) IsBinaryOperatorContext() {}

func NewBinaryOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BinaryOperatorContext {
	var p = new(BinaryOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_binaryOperator

	return p
}

func (s *BinaryOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *BinaryOperatorContext) T_EQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EQUAL, 0)
}

func (s *BinaryOperatorContext) T_NOTEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL, 0)
}

func (s *BinaryOperatorContext) T_NOTEQUAL2() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOTEQUAL2, 0)
}

func (s *BinaryOperatorContext) T_LESS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESS, 0)
}

func (s *BinaryOperatorContext) T_LESSEQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LESSEQUAL, 0)
}

func (s *BinaryOperatorContext) T_GREATER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATER, 0)
}

func (s *BinaryOperatorContext) T_GREATEREQUAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GREATEREQUAL, 0)
}

func (s *BinaryOperatorContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *BinaryOperatorContext) T_REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REGEXP, 0)
}

func (s *BinaryOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *BinaryOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterBinaryOperator(s)
	}
}

func (s *BinaryOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitBinaryOperator(s)
	}
}

func (s *BinaryOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitBinaryOperator(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) BinaryOperator() (localctx IBinaryOperatorContext) {
	this := p
	_ = this

	localctx = NewBinaryOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 120, SQLParserRULE_binaryOperator)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(598)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserT_EQUAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(590)
			p.Match(SQLParserT_EQUAL)
		}


	case SQLParserT_NOTEQUAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(591)
			p.Match(SQLParserT_NOTEQUAL)
		}


	case SQLParserT_NOTEQUAL2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(592)
			p.Match(SQLParserT_NOTEQUAL2)
		}


	case SQLParserT_LESS:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(593)
			p.Match(SQLParserT_LESS)
		}


	case SQLParserT_LESSEQUAL:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(594)
			p.Match(SQLParserT_LESSEQUAL)
		}


	case SQLParserT_GREATER:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(595)
			p.Match(SQLParserT_GREATER)
		}


	case SQLParserT_GREATEREQUAL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(596)
			p.Match(SQLParserT_GREATEREQUAL)
		}


	case SQLParserT_LIKE, SQLParserT_REGEXP:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(597)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_LIKE || _la == SQLParserT_REGEXP) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}


	return localctx
}


// IFieldExprContext is an interface to support dynamic dispatch.
type IFieldExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFieldExprContext differentiates from other interfaces.
	IsFieldExprContext()
}

type FieldExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFieldExprContext() *FieldExprContext {
	var p = new(FieldExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_fieldExpr
	return p
}

func (*FieldExprContext) IsFieldExprContext() {}

func NewFieldExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldExprContext {
	var p = new(FieldExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_fieldExpr

	return p
}

func (s *FieldExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldExprContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *FieldExprContext) AllFieldExpr() []IFieldExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFieldExprContext); ok {
			len++
		}
	}

	tst := make([]IFieldExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFieldExprContext); ok {
			tst[i] = t.(IFieldExprContext)
			i++
		}
	}

	return tst
}

func (s *FieldExprContext) FieldExpr(i int) IFieldExprContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *FieldExprContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *FieldExprContext) ExprFunc() IExprFuncContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprFuncContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprFuncContext)
}

func (s *FieldExprContext) ExprAtom() IExprAtomContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprAtomContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprAtomContext)
}

func (s *FieldExprContext) DurationLit() IDurationLitContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDurationLitContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDurationLitContext)
}

func (s *FieldExprContext) T_MUL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MUL, 0)
}

func (s *FieldExprContext) T_DIV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DIV, 0)
}

func (s *FieldExprContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *FieldExprContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *FieldExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FieldExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFieldExpr(s)
	}
}

func (s *FieldExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFieldExpr(s)
	}
}

func (s *FieldExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFieldExpr(s)

	default:
		return t.VisitChildren(s)
	}
}





func (p *SQLParser) FieldExpr() (localctx IFieldExprContext) {
	return p.fieldExpr(0)
}

func (p *SQLParser) fieldExpr(_p int) (localctx IFieldExprContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewFieldExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IFieldExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 122
	p.EnterRecursionRule(localctx, 122, SQLParserRULE_fieldExpr, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(608)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 46, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(601)
			p.Match(SQLParserT_OPEN_P)
		}
		{
			p.SetState(602)
			p.fieldExpr(0)
		}
		{
			p.SetState(603)
			p.Match(SQLParserT_CLOSE_P)
		}


	case 2:
		{
			p.SetState(605)
			p.ExprFunc()
		}


	case 3:
		{
			p.SetState(606)
			p.ExprAtom()
		}


	case 4:
		{
			p.SetState(607)
			p.DurationLit()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(624)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(622)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 47, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(610)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
				}
				{
					p.SetState(611)
					p.Match(SQLParserT_MUL)
				}
				{
					p.SetState(612)
					p.fieldExpr(9)
				}


			case 2:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(613)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
				}
				{
					p.SetState(614)
					p.Match(SQLParserT_DIV)
				}
				{
					p.SetState(615)
					p.fieldExpr(8)
				}


			case 3:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(616)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
				}
				{
					p.SetState(617)
					p.Match(SQLParserT_ADD)
				}
				{
					p.SetState(618)
					p.fieldExpr(7)
				}


			case 4:
				localctx = NewFieldExprContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_fieldExpr)
				p.SetState(619)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
				}
				{
					p.SetState(620)
					p.Match(SQLParserT_SUB)
				}
				{
					p.SetState(621)
					p.fieldExpr(6)
				}

			}

		}
		p.SetState(626)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())
	}



	return localctx
}


// IDurationLitContext is an interface to support dynamic dispatch.
type IDurationLitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDurationLitContext differentiates from other interfaces.
	IsDurationLitContext()
}

type DurationLitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDurationLitContext() *DurationLitContext {
	var p = new(DurationLitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_durationLit
	return p
}

func (*DurationLitContext) IsDurationLitContext() {}

func NewDurationLitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DurationLitContext {
	var p = new(DurationLitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_durationLit

	return p
}

func (s *DurationLitContext) GetParser() antlr.Parser { return s.parser }

func (s *DurationLitContext) IntNumber() IIntNumberContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntNumberContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *DurationLitContext) IntervalItem() IIntervalItemContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntervalItemContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntervalItemContext)
}

func (s *DurationLitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DurationLitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DurationLitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDurationLit(s)
	}
}

func (s *DurationLitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDurationLit(s)
	}
}

func (s *DurationLitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDurationLit(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) DurationLit() (localctx IDurationLitContext) {
	this := p
	_ = this

	localctx = NewDurationLitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 124, SQLParserRULE_durationLit)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(627)
		p.IntNumber()
	}
	{
		p.SetState(628)
		p.IntervalItem()
	}



	return localctx
}


// IIntervalItemContext is an interface to support dynamic dispatch.
type IIntervalItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntervalItemContext differentiates from other interfaces.
	IsIntervalItemContext()
}

type IntervalItemContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalItemContext() *IntervalItemContext {
	var p = new(IntervalItemContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intervalItem
	return p
}

func (*IntervalItemContext) IsIntervalItemContext() {}

func NewIntervalItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalItemContext {
	var p = new(IntervalItemContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalItem

	return p
}

func (s *IntervalItemContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalItemContext) T_SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SECOND, 0)
}

func (s *IntervalItemContext) T_MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MINUTE, 0)
}

func (s *IntervalItemContext) T_HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HOUR, 0)
}

func (s *IntervalItemContext) T_DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DAY, 0)
}

func (s *IntervalItemContext) T_WEEK() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WEEK, 0)
}

func (s *IntervalItemContext) T_MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MONTH, 0)
}

func (s *IntervalItemContext) T_YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_YEAR, 0)
}

func (s *IntervalItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IntervalItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntervalItem(s)
	}
}

func (s *IntervalItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntervalItem(s)
	}
}

func (s *IntervalItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntervalItem(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) IntervalItem() (localctx IIntervalItemContext) {
	this := p
	_ = this

	localctx = NewIntervalItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 126, SQLParserRULE_intervalItem)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(630)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 90)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 90))) & ((1 << (SQLParserT_SECOND - 90)) | (1 << (SQLParserT_MINUTE - 90)) | (1 << (SQLParserT_HOUR - 90)) | (1 << (SQLParserT_DAY - 90)) | (1 << (SQLParserT_WEEK - 90)) | (1 << (SQLParserT_MONTH - 90)) | (1 << (SQLParserT_YEAR - 90)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


// IExprFuncContext is an interface to support dynamic dispatch.
type IExprFuncContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprFuncContext differentiates from other interfaces.
	IsExprFuncContext()
}

type ExprFuncContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprFuncContext() *ExprFuncContext {
	var p = new(ExprFuncContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_exprFunc
	return p
}

func (*ExprFuncContext) IsExprFuncContext() {}

func NewExprFuncContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprFuncContext {
	var p = new(ExprFuncContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_exprFunc

	return p
}

func (s *ExprFuncContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprFuncContext) FuncName() IFuncNameContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFuncNameContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFuncNameContext)
}

func (s *ExprFuncContext) T_OPEN_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_P, 0)
}

func (s *ExprFuncContext) T_CLOSE_P() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_P, 0)
}

func (s *ExprFuncContext) ExprFuncParams() IExprFuncParamsContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprFuncParamsContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprFuncParamsContext)
}

func (s *ExprFuncContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprFuncContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ExprFuncContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExprFunc(s)
	}
}

func (s *ExprFuncContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExprFunc(s)
	}
}

func (s *ExprFuncContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExprFunc(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ExprFunc() (localctx IExprFuncContext) {
	this := p
	_ = this

	localctx = NewExprFuncContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 128, SQLParserRULE_exprFunc)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(632)
		p.FuncName()
	}
	{
		p.SetState(633)
		p.Match(SQLParserT_OPEN_P)
	}
	p.SetState(635)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if (((_la) & -(0x1f+1)) == 0 && ((1 << uint(_la)) & ((1 << SQLParserT_CREATE) | (1 << SQLParserT_UPDATE) | (1 << SQLParserT_SET) | (1 << SQLParserT_DROP) | (1 << SQLParserT_INTERVAL) | (1 << SQLParserT_INTERVAL_NAME) | (1 << SQLParserT_SHARD) | (1 << SQLParserT_REPLICATION) | (1 << SQLParserT_TTL) | (1 << SQLParserT_META_TTL) | (1 << SQLParserT_PAST_TTL) | (1 << SQLParserT_FUTURE_TTL) | (1 << SQLParserT_KILL) | (1 << SQLParserT_ON) | (1 << SQLParserT_SHOW) | (1 << SQLParserT_USE) | (1 << SQLParserT_STATE_REPO) | (1 << SQLParserT_STATE_MACHINE) | (1 << SQLParserT_MASTER) | (1 << SQLParserT_METADATA) | (1 << SQLParserT_TYPES) | (1 << SQLParserT_TYPE) | (1 << SQLParserT_STORAGES) | (1 << SQLParserT_STORAGE) | (1 << SQLParserT_BROKER) | (1 << SQLParserT_ALIVE))) != 0) || ((((_la - 32)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 32))) & ((1 << (SQLParserT_SCHEMAS - 32)) | (1 << (SQLParserT_DATASBAE - 32)) | (1 << (SQLParserT_DATASBAES - 32)) | (1 << (SQLParserT_NAMESPACE - 32)) | (1 << (SQLParserT_NAMESPACES - 32)) | (1 << (SQLParserT_NODE - 32)) | (1 << (SQLParserT_METRICS - 32)) | (1 << (SQLParserT_METRIC - 32)) | (1 << (SQLParserT_FIELD - 32)) | (1 << (SQLParserT_FIELDS - 32)) | (1 << (SQLParserT_TAG - 32)) | (1 << (SQLParserT_INFO - 32)) | (1 << (SQLParserT_KEYS - 32)) | (1 << (SQLParserT_KEY - 32)) | (1 << (SQLParserT_WITH - 32)) | (1 << (SQLParserT_VALUES - 32)) | (1 << (SQLParserT_VALUE - 32)) | (1 << (SQLParserT_FROM - 32)) | (1 << (SQLParserT_WHERE - 32)) | (1 << (SQLParserT_LIMIT - 32)) | (1 << (SQLParserT_QUERIES - 32)) | (1 << (SQLParserT_QUERY - 32)) | (1 << (SQLParserT_EXPLAIN - 32)) | (1 << (SQLParserT_WITH_VALUE - 32)) | (1 << (SQLParserT_SELECT - 32)) | (1 << (SQLParserT_AS - 32)) | (1 << (SQLParserT_AND - 32)) | (1 << (SQLParserT_OR - 32)) | (1 << (SQLParserT_FILL - 32)) | (1 << (SQLParserT_NULL - 32)) | (1 << (SQLParserT_PREVIOUS - 32)) | (1 << (SQLParserT_ORDER - 32)))) != 0) || ((((_la - 64)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 64))) & ((1 << (SQLParserT_ASC - 64)) | (1 << (SQLParserT_DESC - 64)) | (1 << (SQLParserT_LIKE - 64)) | (1 << (SQLParserT_NOT - 64)) | (1 << (SQLParserT_BETWEEN - 64)) | (1 << (SQLParserT_IS - 64)) | (1 << (SQLParserT_GROUP - 64)) | (1 << (SQLParserT_HAVING - 64)) | (1 << (SQLParserT_BY - 64)) | (1 << (SQLParserT_FOR - 64)) | (1 << (SQLParserT_STATS - 64)) | (1 << (SQLParserT_TIME - 64)) | (1 << (SQLParserT_NOW - 64)) | (1 << (SQLParserT_IN - 64)) | (1 << (SQLParserT_LOG - 64)) | (1 << (SQLParserT_PROFILE - 64)) | (1 << (SQLParserT_SUM - 64)) | (1 << (SQLParserT_MIN - 64)) | (1 << (SQLParserT_MAX - 64)) | (1 << (SQLParserT_COUNT - 64)) | (1 << (SQLParserT_LAST - 64)) | (1 << (SQLParserT_FIRST - 64)) | (1 << (SQLParserT_AVG - 64)) | (1 << (SQLParserT_STDDEV - 64)) | (1 << (SQLParserT_QUANTILE - 64)) | (1 << (SQLParserT_RATE - 64)) | (1 << (SQLParserT_SECOND - 64)) | (1 << (SQLParserT_MINUTE - 64)) | (1 << (SQLParserT_HOUR - 64)) | (1 << (SQLParserT_DAY - 64)) | (1 << (SQLParserT_WEEK - 64)) | (1 << (SQLParserT_MONTH - 64)))) != 0) || ((((_la - 96)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 96))) & ((1 << (SQLParserT_YEAR - 96)) | (1 << (SQLParserT_OPEN_P - 96)) | (1 << (SQLParserT_ADD - 96)) | (1 << (SQLParserT_SUB - 96)) | (1 << (SQLParserL_ID - 96)) | (1 << (SQLParserL_INT - 96)) | (1 << (SQLParserL_DEC - 96)))) != 0) {
		{
			p.SetState(634)
			p.ExprFuncParams()
		}

	}
	{
		p.SetState(637)
		p.Match(SQLParserT_CLOSE_P)
	}



	return localctx
}


// IFuncNameContext is an interface to support dynamic dispatch.
type IFuncNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncNameContext differentiates from other interfaces.
	IsFuncNameContext()
}

type FuncNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncNameContext() *FuncNameContext {
	var p = new(FuncNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_funcName
	return p
}

func (*FuncNameContext) IsFuncNameContext() {}

func NewFuncNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncNameContext {
	var p = new(FuncNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_funcName

	return p
}

func (s *FuncNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncNameContext) T_SUM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUM, 0)
}

func (s *FuncNameContext) T_MIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MIN, 0)
}

func (s *FuncNameContext) T_MAX() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MAX, 0)
}

func (s *FuncNameContext) T_AVG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AVG, 0)
}

func (s *FuncNameContext) T_COUNT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COUNT, 0)
}

func (s *FuncNameContext) T_LAST() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LAST, 0)
}

func (s *FuncNameContext) T_FIRST() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIRST, 0)
}

func (s *FuncNameContext) T_STDDEV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STDDEV, 0)
}

func (s *FuncNameContext) T_QUANTILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUANTILE, 0)
}

func (s *FuncNameContext) T_RATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_RATE, 0)
}

func (s *FuncNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FuncNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFuncName(s)
	}
}

func (s *FuncNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFuncName(s)
	}
}

func (s *FuncNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFuncName(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) FuncName() (localctx IFuncNameContext) {
	this := p
	_ = this

	localctx = NewFuncNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 130, SQLParserRULE_funcName)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(639)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 80)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 80))) & ((1 << (SQLParserT_SUM - 80)) | (1 << (SQLParserT_MIN - 80)) | (1 << (SQLParserT_MAX - 80)) | (1 << (SQLParserT_COUNT - 80)) | (1 << (SQLParserT_LAST - 80)) | (1 << (SQLParserT_FIRST - 80)) | (1 << (SQLParserT_AVG - 80)) | (1 << (SQLParserT_STDDEV - 80)) | (1 << (SQLParserT_QUANTILE - 80)) | (1 << (SQLParserT_RATE - 80)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


// IExprFuncParamsContext is an interface to support dynamic dispatch.
type IExprFuncParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprFuncParamsContext differentiates from other interfaces.
	IsExprFuncParamsContext()
}

type ExprFuncParamsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprFuncParamsContext() *ExprFuncParamsContext {
	var p = new(ExprFuncParamsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_exprFuncParams
	return p
}

func (*ExprFuncParamsContext) IsExprFuncParamsContext() {}

func NewExprFuncParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprFuncParamsContext {
	var p = new(ExprFuncParamsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_exprFuncParams

	return p
}

func (s *ExprFuncParamsContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprFuncParamsContext) AllFuncParam() []IFuncParamContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFuncParamContext); ok {
			len++
		}
	}

	tst := make([]IFuncParamContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFuncParamContext); ok {
			tst[i] = t.(IFuncParamContext)
			i++
		}
	}

	return tst
}

func (s *ExprFuncParamsContext) FuncParam(i int) IFuncParamContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFuncParamContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFuncParamContext)
}

func (s *ExprFuncParamsContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *ExprFuncParamsContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *ExprFuncParamsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprFuncParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ExprFuncParamsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExprFuncParams(s)
	}
}

func (s *ExprFuncParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExprFuncParams(s)
	}
}

func (s *ExprFuncParamsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExprFuncParams(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ExprFuncParams() (localctx IExprFuncParamsContext) {
	this := p
	_ = this

	localctx = NewExprFuncParamsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 132, SQLParserRULE_exprFuncParams)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(641)
		p.FuncParam()
	}
	p.SetState(646)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	for _la == SQLParserT_COMMA {
		{
			p.SetState(642)
			p.Match(SQLParserT_COMMA)
		}
		{
			p.SetState(643)
			p.FuncParam()
		}


		p.SetState(648)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}



	return localctx
}


// IFuncParamContext is an interface to support dynamic dispatch.
type IFuncParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFuncParamContext differentiates from other interfaces.
	IsFuncParamContext()
}

type FuncParamContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFuncParamContext() *FuncParamContext {
	var p = new(FuncParamContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_funcParam
	return p
}

func (*FuncParamContext) IsFuncParamContext() {}

func NewFuncParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncParamContext {
	var p = new(FuncParamContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_funcParam

	return p
}

func (s *FuncParamContext) GetParser() antlr.Parser { return s.parser }

func (s *FuncParamContext) FieldExpr() IFieldExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldExprContext)
}

func (s *FuncParamContext) TagFilterExpr() ITagFilterExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagFilterExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
}

func (s *FuncParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *FuncParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterFuncParam(s)
	}
}

func (s *FuncParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitFuncParam(s)
	}
}

func (s *FuncParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitFuncParam(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) FuncParam() (localctx IFuncParamContext) {
	this := p
	_ = this

	localctx = NewFuncParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 134, SQLParserRULE_funcParam)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(651)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(649)
			p.fieldExpr(0)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(650)
			p.tagFilterExpr(0)
		}

	}


	return localctx
}


// IExprAtomContext is an interface to support dynamic dispatch.
type IExprAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprAtomContext differentiates from other interfaces.
	IsExprAtomContext()
}

type ExprAtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprAtomContext() *ExprAtomContext {
	var p = new(ExprAtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_exprAtom
	return p
}

func (*ExprAtomContext) IsExprAtomContext() {}

func NewExprAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprAtomContext {
	var p = new(ExprAtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_exprAtom

	return p
}

func (s *ExprAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprAtomContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ExprAtomContext) IdentFilter() IIdentFilterContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentFilterContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentFilterContext)
}

func (s *ExprAtomContext) DecNumber() IDecNumberContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDecNumberContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDecNumberContext)
}

func (s *ExprAtomContext) IntNumber() IIntNumberContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntNumberContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *ExprAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ExprAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterExprAtom(s)
	}
}

func (s *ExprAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitExprAtom(s)
	}
}

func (s *ExprAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitExprAtom(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) ExprAtom() (localctx IExprAtomContext) {
	this := p
	_ = this

	localctx = NewExprAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 136, SQLParserRULE_exprAtom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(659)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 53, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(653)
			p.Ident()
		}
		p.SetState(655)
		p.GetErrorHandler().Sync(p)


		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 52, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(654)
				p.IdentFilter()
			}


		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(657)
			p.DecNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(658)
			p.IntNumber()
		}

	}


	return localctx
}


// IIdentFilterContext is an interface to support dynamic dispatch.
type IIdentFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentFilterContext differentiates from other interfaces.
	IsIdentFilterContext()
}

type IdentFilterContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentFilterContext() *IdentFilterContext {
	var p = new(IdentFilterContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_identFilter
	return p
}

func (*IdentFilterContext) IsIdentFilterContext() {}

func NewIdentFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentFilterContext {
	var p = new(IdentFilterContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_identFilter

	return p
}

func (s *IdentFilterContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentFilterContext) T_OPEN_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_SB, 0)
}

func (s *IdentFilterContext) TagFilterExpr() ITagFilterExprContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagFilterExprContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagFilterExprContext)
}

func (s *IdentFilterContext) T_CLOSE_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_SB, 0)
}

func (s *IdentFilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentFilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IdentFilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIdentFilter(s)
	}
}

func (s *IdentFilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIdentFilter(s)
	}
}

func (s *IdentFilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIdentFilter(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) IdentFilter() (localctx IIdentFilterContext) {
	this := p
	_ = this

	localctx = NewIdentFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 138, SQLParserRULE_identFilter)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(661)
		p.Match(SQLParserT_OPEN_SB)
	}
	{
		p.SetState(662)
		p.tagFilterExpr(0)
	}
	{
		p.SetState(663)
		p.Match(SQLParserT_CLOSE_SB)
	}



	return localctx
}


// IJsonContext is an interface to support dynamic dispatch.
type IJsonContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsJsonContext differentiates from other interfaces.
	IsJsonContext()
}

type JsonContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJsonContext() *JsonContext {
	var p = new(JsonContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_json
	return p
}

func (*JsonContext) IsJsonContext() {}

func NewJsonContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JsonContext {
	var p = new(JsonContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_json

	return p
}

func (s *JsonContext) GetParser() antlr.Parser { return s.parser }

func (s *JsonContext) Value() IValueContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *JsonContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JsonContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *JsonContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterJson(s)
	}
}

func (s *JsonContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitJson(s)
	}
}

func (s *JsonContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitJson(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Json() (localctx IJsonContext) {
	this := p
	_ = this

	localctx = NewJsonContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 140, SQLParserRULE_json)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(665)
		p.Value()
	}



	return localctx
}


// IObjContext is an interface to support dynamic dispatch.
type IObjContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsObjContext differentiates from other interfaces.
	IsObjContext()
}

type ObjContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyObjContext() *ObjContext {
	var p = new(ObjContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_obj
	return p
}

func (*ObjContext) IsObjContext() {}

func NewObjContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjContext {
	var p = new(ObjContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_obj

	return p
}

func (s *ObjContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjContext) T_OPEN_B() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_B, 0)
}

func (s *ObjContext) AllPair() []IPairContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPairContext); ok {
			len++
		}
	}

	tst := make([]IPairContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPairContext); ok {
			tst[i] = t.(IPairContext)
			i++
		}
	}

	return tst
}

func (s *ObjContext) Pair(i int) IPairContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPairContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPairContext)
}

func (s *ObjContext) T_CLOSE_B() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_B, 0)
}

func (s *ObjContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *ObjContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *ObjContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ObjContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterObj(s)
	}
}

func (s *ObjContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitObj(s)
	}
}

func (s *ObjContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitObj(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Obj() (localctx IObjContext) {
	this := p
	_ = this

	localctx = NewObjContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 142, SQLParserRULE_obj)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(680)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 55, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(667)
			p.Match(SQLParserT_OPEN_B)
		}
		{
			p.SetState(668)
			p.Pair()
		}
		p.SetState(673)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == SQLParserT_COMMA {
			{
				p.SetState(669)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(670)
				p.Pair()
			}


			p.SetState(675)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(676)
			p.Match(SQLParserT_CLOSE_B)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(678)
			p.Match(SQLParserT_OPEN_B)
		}
		{
			p.SetState(679)
			p.Match(SQLParserT_CLOSE_B)
		}

	}


	return localctx
}


// IPairContext is an interface to support dynamic dispatch.
type IPairContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPairContext differentiates from other interfaces.
	IsPairContext()
}

type PairContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPairContext() *PairContext {
	var p = new(PairContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_pair
	return p
}

func (*PairContext) IsPairContext() {}

func NewPairContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PairContext {
	var p = new(PairContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_pair

	return p
}

func (s *PairContext) GetParser() antlr.Parser { return s.parser }

func (s *PairContext) STRING() antlr.TerminalNode {
	return s.GetToken(SQLParserSTRING, 0)
}

func (s *PairContext) T_COLON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COLON, 0)
}

func (s *PairContext) Value() IValueContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *PairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *PairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterPair(s)
	}
}

func (s *PairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitPair(s)
	}
}

func (s *PairContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitPair(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Pair() (localctx IPairContext) {
	this := p
	_ = this

	localctx = NewPairContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 144, SQLParserRULE_pair)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(682)
		p.Match(SQLParserSTRING)
	}
	{
		p.SetState(683)
		p.Match(SQLParserT_COLON)
	}
	{
		p.SetState(684)
		p.Value()
	}



	return localctx
}


// IArrContext is an interface to support dynamic dispatch.
type IArrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArrContext differentiates from other interfaces.
	IsArrContext()
}

type ArrContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArrContext() *ArrContext {
	var p = new(ArrContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_arr
	return p
}

func (*ArrContext) IsArrContext() {}

func NewArrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArrContext {
	var p = new(ArrContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_arr

	return p
}

func (s *ArrContext) GetParser() antlr.Parser { return s.parser }

func (s *ArrContext) T_OPEN_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OPEN_SB, 0)
}

func (s *ArrContext) AllValue() []IValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueContext); ok {
			len++
		}
	}

	tst := make([]IValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueContext); ok {
			tst[i] = t.(IValueContext)
			i++
		}
	}

	return tst
}

func (s *ArrContext) Value(i int) IValueContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *ArrContext) T_CLOSE_SB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CLOSE_SB, 0)
}

func (s *ArrContext) AllT_COMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_COMMA)
}

func (s *ArrContext) T_COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_COMMA, i)
}

func (s *ArrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ArrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterArr(s)
	}
}

func (s *ArrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitArr(s)
	}
}

func (s *ArrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitArr(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Arr() (localctx IArrContext) {
	this := p
	_ = this

	localctx = NewArrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 146, SQLParserRULE_arr)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(699)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 57, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(686)
			p.Match(SQLParserT_OPEN_SB)
		}
		{
			p.SetState(687)
			p.Value()
		}
		p.SetState(692)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)


		for _la == SQLParserT_COMMA {
			{
				p.SetState(688)
				p.Match(SQLParserT_COMMA)
			}
			{
				p.SetState(689)
				p.Value()
			}


			p.SetState(694)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(695)
			p.Match(SQLParserT_CLOSE_SB)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(697)
			p.Match(SQLParserT_OPEN_SB)
		}
		{
			p.SetState(698)
			p.Match(SQLParserT_CLOSE_SB)
		}

	}


	return localctx
}


// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_value
	return p
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) STRING() antlr.TerminalNode {
	return s.GetToken(SQLParserSTRING, 0)
}

func (s *ValueContext) IntNumber() IIntNumberContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntNumberContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntNumberContext)
}

func (s *ValueContext) DecNumber() IDecNumberContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDecNumberContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDecNumberContext)
}

func (s *ValueContext) Obj() IObjContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjContext)
}

func (s *ValueContext) Arr() IArrContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArrContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArrContext)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitValue(s)
	}
}

func (s *ValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitValue(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Value() (localctx IValueContext) {
	this := p
	_ = this

	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 148, SQLParserRULE_value)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(709)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 58, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(701)
			p.Match(SQLParserSTRING)
		}


	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(702)
			p.IntNumber()
		}


	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(703)
			p.DecNumber()
		}


	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(704)
			p.Obj()
		}


	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(705)
			p.Arr()
		}


	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(706)
			p.Match(SQLParserT__0)
		}


	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(707)
			p.Match(SQLParserT__1)
		}


	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(708)
			p.Match(SQLParserT__2)
		}

	}


	return localctx
}


// IIntNumberContext is an interface to support dynamic dispatch.
type IIntNumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntNumberContext differentiates from other interfaces.
	IsIntNumberContext()
}

type IntNumberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntNumberContext() *IntNumberContext {
	var p = new(IntNumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_intNumber
	return p
}

func (*IntNumberContext) IsIntNumberContext() {}

func NewIntNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntNumberContext {
	var p = new(IntNumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intNumber

	return p
}

func (s *IntNumberContext) GetParser() antlr.Parser { return s.parser }

func (s *IntNumberContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *IntNumberContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *IntNumberContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *IntNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntNumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IntNumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIntNumber(s)
	}
}

func (s *IntNumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIntNumber(s)
	}
}

func (s *IntNumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIntNumber(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) IntNumber() (localctx IIntNumberContext) {
	this := p
	_ = this

	localctx = NewIntNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 150, SQLParserRULE_intNumber)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(712)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(711)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ADD || _la == SQLParserT_SUB) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(714)
		p.Match(SQLParserL_INT)
	}



	return localctx
}


// IDecNumberContext is an interface to support dynamic dispatch.
type IDecNumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDecNumberContext differentiates from other interfaces.
	IsDecNumberContext()
}

type DecNumberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDecNumberContext() *DecNumberContext {
	var p = new(DecNumberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_decNumber
	return p
}

func (*DecNumberContext) IsDecNumberContext() {}

func NewDecNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DecNumberContext {
	var p = new(DecNumberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_decNumber

	return p
}

func (s *DecNumberContext) GetParser() antlr.Parser { return s.parser }

func (s *DecNumberContext) L_DEC() antlr.TerminalNode {
	return s.GetToken(SQLParserL_DEC, 0)
}

func (s *DecNumberContext) T_SUB() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUB, 0)
}

func (s *DecNumberContext) T_ADD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ADD, 0)
}

func (s *DecNumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecNumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *DecNumberContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterDecNumber(s)
	}
}

func (s *DecNumberContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitDecNumber(s)
	}
}

func (s *DecNumberContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitDecNumber(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) DecNumber() (localctx IDecNumberContext) {
	this := p
	_ = this

	localctx = NewDecNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 152, SQLParserRULE_decNumber)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(717)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)


	if _la == SQLParserT_ADD || _la == SQLParserT_SUB {
		{
			p.SetState(716)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserT_ADD || _la == SQLParserT_SUB) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(719)
		p.Match(SQLParserL_DEC)
	}



	return localctx
}


// ILimitClauseContext is an interface to support dynamic dispatch.
type ILimitClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLimitClauseContext differentiates from other interfaces.
	IsLimitClauseContext()
}

type LimitClauseContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitClauseContext() *LimitClauseContext {
	var p = new(LimitClauseContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_limitClause
	return p
}

func (*LimitClauseContext) IsLimitClauseContext() {}

func NewLimitClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitClauseContext {
	var p = new(LimitClauseContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_limitClause

	return p
}

func (s *LimitClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitClauseContext) T_LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIMIT, 0)
}

func (s *LimitClauseContext) L_INT() antlr.TerminalNode {
	return s.GetToken(SQLParserL_INT, 0)
}

func (s *LimitClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *LimitClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterLimitClause(s)
	}
}

func (s *LimitClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitLimitClause(s)
	}
}

func (s *LimitClauseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitLimitClause(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) LimitClause() (localctx ILimitClauseContext) {
	this := p
	_ = this

	localctx = NewLimitClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 154, SQLParserRULE_limitClause)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(721)
		p.Match(SQLParserT_LIMIT)
	}
	{
		p.SetState(722)
		p.Match(SQLParserL_INT)
	}



	return localctx
}


// IMetricNameContext is an interface to support dynamic dispatch.
type IMetricNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMetricNameContext differentiates from other interfaces.
	IsMetricNameContext()
}

type MetricNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMetricNameContext() *MetricNameContext {
	var p = new(MetricNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_metricName
	return p
}

func (*MetricNameContext) IsMetricNameContext() {}

func NewMetricNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MetricNameContext {
	var p = new(MetricNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_metricName

	return p
}

func (s *MetricNameContext) GetParser() antlr.Parser { return s.parser }

func (s *MetricNameContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *MetricNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MetricNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *MetricNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterMetricName(s)
	}
}

func (s *MetricNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitMetricName(s)
	}
}

func (s *MetricNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitMetricName(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) MetricName() (localctx IMetricNameContext) {
	this := p
	_ = this

	localctx = NewMetricNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 156, SQLParserRULE_metricName)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(724)
		p.Ident()
	}



	return localctx
}


// ITagKeyContext is an interface to support dynamic dispatch.
type ITagKeyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagKeyContext differentiates from other interfaces.
	IsTagKeyContext()
}

type TagKeyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagKeyContext() *TagKeyContext {
	var p = new(TagKeyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagKey
	return p
}

func (*TagKeyContext) IsTagKeyContext() {}

func NewTagKeyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagKeyContext {
	var p = new(TagKeyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagKey

	return p
}

func (s *TagKeyContext) GetParser() antlr.Parser { return s.parser }

func (s *TagKeyContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TagKeyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagKeyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TagKeyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagKey(s)
	}
}

func (s *TagKeyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagKey(s)
	}
}

func (s *TagKeyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagKey(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) TagKey() (localctx ITagKeyContext) {
	this := p
	_ = this

	localctx = NewTagKeyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 158, SQLParserRULE_tagKey)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(726)
		p.Ident()
	}



	return localctx
}


// ITagValueContext is an interface to support dynamic dispatch.
type ITagValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagValueContext differentiates from other interfaces.
	IsTagValueContext()
}

type TagValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagValueContext() *TagValueContext {
	var p = new(TagValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_tagValue
	return p
}

func (*TagValueContext) IsTagValueContext() {}

func NewTagValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagValueContext {
	var p = new(TagValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_tagValue

	return p
}

func (s *TagValueContext) GetParser() antlr.Parser { return s.parser }

func (s *TagValueContext) Ident() IIdentContext {
	var t antlr.RuleContext;
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext);
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *TagValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *TagValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterTagValue(s)
	}
}

func (s *TagValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitTagValue(s)
	}
}

func (s *TagValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitTagValue(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) TagValue() (localctx ITagValueContext) {
	this := p
	_ = this

	localctx = NewTagValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 160, SQLParserRULE_tagValue)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(728)
		p.Ident()
	}



	return localctx
}


// IIdentContext is an interface to support dynamic dispatch.
type IIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentContext differentiates from other interfaces.
	IsIdentContext()
}

type IdentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentContext() *IdentContext {
	var p = new(IdentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_ident
	return p
}

func (*IdentContext) IsIdentContext() {}

func NewIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentContext {
	var p = new(IdentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_ident

	return p
}

func (s *IdentContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentContext) AllL_ID() []antlr.TerminalNode {
	return s.GetTokens(SQLParserL_ID)
}

func (s *IdentContext) L_ID(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserL_ID, i)
}

func (s *IdentContext) AllNonReservedWords() []INonReservedWordsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INonReservedWordsContext); ok {
			len++
		}
	}

	tst := make([]INonReservedWordsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INonReservedWordsContext); ok {
			tst[i] = t.(INonReservedWordsContext)
			i++
		}
	}

	return tst
}

func (s *IdentContext) NonReservedWords(i int) INonReservedWordsContext {
	var t antlr.RuleContext;
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INonReservedWordsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext);
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INonReservedWordsContext)
}

func (s *IdentContext) AllT_DOT() []antlr.TerminalNode {
	return s.GetTokens(SQLParserT_DOT)
}

func (s *IdentContext) T_DOT(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserT_DOT, i)
}

func (s *IdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *IdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterIdent(s)
	}
}

func (s *IdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitIdent(s)
	}
}

func (s *IdentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitIdent(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) Ident() (localctx IIdentContext) {
	this := p
	_ = this

	localctx = NewIdentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 162, SQLParserRULE_ident)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(732)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SQLParserL_ID:
		{
			p.SetState(730)
			p.Match(SQLParserL_ID)
		}


	case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_USE, SQLParserT_STATE_REPO, SQLParserT_STATE_MACHINE, SQLParserT_MASTER, SQLParserT_METADATA, SQLParserT_TYPES, SQLParserT_TYPE, SQLParserT_STORAGES, SQLParserT_STORAGE, SQLParserT_BROKER, SQLParserT_ALIVE, SQLParserT_SCHEMAS, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_LAST, SQLParserT_FIRST, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_RATE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
		{
			p.SetState(731)
			p.NonReservedWords()
		}



	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(741)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 63, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(734)
				p.Match(SQLParserT_DOT)
			}
			p.SetState(737)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case SQLParserL_ID:
				{
					p.SetState(735)
					p.Match(SQLParserL_ID)
				}


			case SQLParserT_CREATE, SQLParserT_UPDATE, SQLParserT_SET, SQLParserT_DROP, SQLParserT_INTERVAL, SQLParserT_INTERVAL_NAME, SQLParserT_SHARD, SQLParserT_REPLICATION, SQLParserT_TTL, SQLParserT_META_TTL, SQLParserT_PAST_TTL, SQLParserT_FUTURE_TTL, SQLParserT_KILL, SQLParserT_ON, SQLParserT_SHOW, SQLParserT_USE, SQLParserT_STATE_REPO, SQLParserT_STATE_MACHINE, SQLParserT_MASTER, SQLParserT_METADATA, SQLParserT_TYPES, SQLParserT_TYPE, SQLParserT_STORAGES, SQLParserT_STORAGE, SQLParserT_BROKER, SQLParserT_ALIVE, SQLParserT_SCHEMAS, SQLParserT_DATASBAE, SQLParserT_DATASBAES, SQLParserT_NAMESPACE, SQLParserT_NAMESPACES, SQLParserT_NODE, SQLParserT_METRICS, SQLParserT_METRIC, SQLParserT_FIELD, SQLParserT_FIELDS, SQLParserT_TAG, SQLParserT_INFO, SQLParserT_KEYS, SQLParserT_KEY, SQLParserT_WITH, SQLParserT_VALUES, SQLParserT_VALUE, SQLParserT_FROM, SQLParserT_WHERE, SQLParserT_LIMIT, SQLParserT_QUERIES, SQLParserT_QUERY, SQLParserT_EXPLAIN, SQLParserT_WITH_VALUE, SQLParserT_SELECT, SQLParserT_AS, SQLParserT_AND, SQLParserT_OR, SQLParserT_FILL, SQLParserT_NULL, SQLParserT_PREVIOUS, SQLParserT_ORDER, SQLParserT_ASC, SQLParserT_DESC, SQLParserT_LIKE, SQLParserT_NOT, SQLParserT_BETWEEN, SQLParserT_IS, SQLParserT_GROUP, SQLParserT_HAVING, SQLParserT_BY, SQLParserT_FOR, SQLParserT_STATS, SQLParserT_TIME, SQLParserT_NOW, SQLParserT_IN, SQLParserT_LOG, SQLParserT_PROFILE, SQLParserT_SUM, SQLParserT_MIN, SQLParserT_MAX, SQLParserT_COUNT, SQLParserT_LAST, SQLParserT_FIRST, SQLParserT_AVG, SQLParserT_STDDEV, SQLParserT_QUANTILE, SQLParserT_RATE, SQLParserT_SECOND, SQLParserT_MINUTE, SQLParserT_HOUR, SQLParserT_DAY, SQLParserT_WEEK, SQLParserT_MONTH, SQLParserT_YEAR:
				{
					p.SetState(736)
					p.NonReservedWords()
				}



			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}


		}
		p.SetState(743)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 63, p.GetParserRuleContext())
	}



	return localctx
}


// INonReservedWordsContext is an interface to support dynamic dispatch.
type INonReservedWordsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNonReservedWordsContext differentiates from other interfaces.
	IsNonReservedWordsContext()
}

type NonReservedWordsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNonReservedWordsContext() *NonReservedWordsContext {
	var p = new(NonReservedWordsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SQLParserRULE_nonReservedWords
	return p
}

func (*NonReservedWordsContext) IsNonReservedWordsContext() {}

func NewNonReservedWordsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NonReservedWordsContext {
	var p = new(NonReservedWordsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nonReservedWords

	return p
}

func (s *NonReservedWordsContext) GetParser() antlr.Parser { return s.parser }

func (s *NonReservedWordsContext) T_CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_CREATE, 0)
}

func (s *NonReservedWordsContext) T_UPDATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_UPDATE, 0)
}

func (s *NonReservedWordsContext) T_SET() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SET, 0)
}

func (s *NonReservedWordsContext) T_DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DROP, 0)
}

func (s *NonReservedWordsContext) T_INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL, 0)
}

func (s *NonReservedWordsContext) T_INTERVAL_NAME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INTERVAL_NAME, 0)
}

func (s *NonReservedWordsContext) T_SHARD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHARD, 0)
}

func (s *NonReservedWordsContext) T_REPLICATION() antlr.TerminalNode {
	return s.GetToken(SQLParserT_REPLICATION, 0)
}

func (s *NonReservedWordsContext) T_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TTL, 0)
}

func (s *NonReservedWordsContext) T_META_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_META_TTL, 0)
}

func (s *NonReservedWordsContext) T_PAST_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PAST_TTL, 0)
}

func (s *NonReservedWordsContext) T_FUTURE_TTL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FUTURE_TTL, 0)
}

func (s *NonReservedWordsContext) T_KILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KILL, 0)
}

func (s *NonReservedWordsContext) T_ON() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ON, 0)
}

func (s *NonReservedWordsContext) T_SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SHOW, 0)
}

func (s *NonReservedWordsContext) T_DATASBAE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAE, 0)
}

func (s *NonReservedWordsContext) T_DATASBAES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DATASBAES, 0)
}

func (s *NonReservedWordsContext) T_NAMESPACE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACE, 0)
}

func (s *NonReservedWordsContext) T_NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NAMESPACES, 0)
}

func (s *NonReservedWordsContext) T_NODE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NODE, 0)
}

func (s *NonReservedWordsContext) T_METRICS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRICS, 0)
}

func (s *NonReservedWordsContext) T_METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METRIC, 0)
}

func (s *NonReservedWordsContext) T_FIELD() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELD, 0)
}

func (s *NonReservedWordsContext) T_FIELDS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIELDS, 0)
}

func (s *NonReservedWordsContext) T_TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TAG, 0)
}

func (s *NonReservedWordsContext) T_INFO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_INFO, 0)
}

func (s *NonReservedWordsContext) T_KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEYS, 0)
}

func (s *NonReservedWordsContext) T_KEY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_KEY, 0)
}

func (s *NonReservedWordsContext) T_WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH, 0)
}

func (s *NonReservedWordsContext) T_VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUES, 0)
}

func (s *NonReservedWordsContext) T_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_VALUE, 0)
}

func (s *NonReservedWordsContext) T_FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FROM, 0)
}

func (s *NonReservedWordsContext) T_WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WHERE, 0)
}

func (s *NonReservedWordsContext) T_LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIMIT, 0)
}

func (s *NonReservedWordsContext) T_QUERIES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERIES, 0)
}

func (s *NonReservedWordsContext) T_QUERY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUERY, 0)
}

func (s *NonReservedWordsContext) T_EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_EXPLAIN, 0)
}

func (s *NonReservedWordsContext) T_WITH_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WITH_VALUE, 0)
}

func (s *NonReservedWordsContext) T_SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SELECT, 0)
}

func (s *NonReservedWordsContext) T_AS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AS, 0)
}

func (s *NonReservedWordsContext) T_AND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AND, 0)
}

func (s *NonReservedWordsContext) T_OR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_OR, 0)
}

func (s *NonReservedWordsContext) T_FILL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FILL, 0)
}

func (s *NonReservedWordsContext) T_NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NULL, 0)
}

func (s *NonReservedWordsContext) T_PREVIOUS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PREVIOUS, 0)
}

func (s *NonReservedWordsContext) T_ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ORDER, 0)
}

func (s *NonReservedWordsContext) T_ASC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ASC, 0)
}

func (s *NonReservedWordsContext) T_DESC() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DESC, 0)
}

func (s *NonReservedWordsContext) T_LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LIKE, 0)
}

func (s *NonReservedWordsContext) T_NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOT, 0)
}

func (s *NonReservedWordsContext) T_BETWEEN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BETWEEN, 0)
}

func (s *NonReservedWordsContext) T_IS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IS, 0)
}

func (s *NonReservedWordsContext) T_GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserT_GROUP, 0)
}

func (s *NonReservedWordsContext) T_HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HAVING, 0)
}

func (s *NonReservedWordsContext) T_BY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BY, 0)
}

func (s *NonReservedWordsContext) T_FOR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FOR, 0)
}

func (s *NonReservedWordsContext) T_STATS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATS, 0)
}

func (s *NonReservedWordsContext) T_TIME() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TIME, 0)
}

func (s *NonReservedWordsContext) T_NOW() antlr.TerminalNode {
	return s.GetToken(SQLParserT_NOW, 0)
}

func (s *NonReservedWordsContext) T_IN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_IN, 0)
}

func (s *NonReservedWordsContext) T_LOG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LOG, 0)
}

func (s *NonReservedWordsContext) T_PROFILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_PROFILE, 0)
}

func (s *NonReservedWordsContext) T_SUM() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SUM, 0)
}

func (s *NonReservedWordsContext) T_MIN() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MIN, 0)
}

func (s *NonReservedWordsContext) T_MAX() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MAX, 0)
}

func (s *NonReservedWordsContext) T_COUNT() antlr.TerminalNode {
	return s.GetToken(SQLParserT_COUNT, 0)
}

func (s *NonReservedWordsContext) T_LAST() antlr.TerminalNode {
	return s.GetToken(SQLParserT_LAST, 0)
}

func (s *NonReservedWordsContext) T_FIRST() antlr.TerminalNode {
	return s.GetToken(SQLParserT_FIRST, 0)
}

func (s *NonReservedWordsContext) T_AVG() antlr.TerminalNode {
	return s.GetToken(SQLParserT_AVG, 0)
}

func (s *NonReservedWordsContext) T_STDDEV() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STDDEV, 0)
}

func (s *NonReservedWordsContext) T_QUANTILE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_QUANTILE, 0)
}

func (s *NonReservedWordsContext) T_RATE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_RATE, 0)
}

func (s *NonReservedWordsContext) T_SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SECOND, 0)
}

func (s *NonReservedWordsContext) T_MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MINUTE, 0)
}

func (s *NonReservedWordsContext) T_HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_HOUR, 0)
}

func (s *NonReservedWordsContext) T_DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserT_DAY, 0)
}

func (s *NonReservedWordsContext) T_WEEK() antlr.TerminalNode {
	return s.GetToken(SQLParserT_WEEK, 0)
}

func (s *NonReservedWordsContext) T_MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MONTH, 0)
}

func (s *NonReservedWordsContext) T_YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserT_YEAR, 0)
}

func (s *NonReservedWordsContext) T_USE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_USE, 0)
}

func (s *NonReservedWordsContext) T_MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_MASTER, 0)
}

func (s *NonReservedWordsContext) T_METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserT_METADATA, 0)
}

func (s *NonReservedWordsContext) T_TYPE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPE, 0)
}

func (s *NonReservedWordsContext) T_TYPES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_TYPES, 0)
}

func (s *NonReservedWordsContext) T_STORAGES() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGES, 0)
}

func (s *NonReservedWordsContext) T_STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STORAGE, 0)
}

func (s *NonReservedWordsContext) T_ALIVE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_ALIVE, 0)
}

func (s *NonReservedWordsContext) T_BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserT_BROKER, 0)
}

func (s *NonReservedWordsContext) T_SCHEMAS() antlr.TerminalNode {
	return s.GetToken(SQLParserT_SCHEMAS, 0)
}

func (s *NonReservedWordsContext) T_STATE_REPO() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_REPO, 0)
}

func (s *NonReservedWordsContext) T_STATE_MACHINE() antlr.TerminalNode {
	return s.GetToken(SQLParserT_STATE_MACHINE, 0)
}

func (s *NonReservedWordsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NonReservedWordsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}


func (s *NonReservedWordsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.EnterNonReservedWords(s)
	}
}

func (s *NonReservedWordsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLListener); ok {
		listenerT.ExitNonReservedWords(s)
	}
}

func (s *NonReservedWordsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLVisitor:
		return t.VisitNonReservedWords(s)

	default:
		return t.VisitChildren(s)
	}
}




func (p *SQLParser) NonReservedWords() (localctx INonReservedWordsContext) {
	this := p
	_ = this

	localctx = NewNonReservedWordsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 164, SQLParserRULE_nonReservedWords)
	var _la int


	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(744)
		_la = p.GetTokenStream().LA(1)

		if !(((((_la - 6)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 6))) & ((1 << (SQLParserT_CREATE - 6)) | (1 << (SQLParserT_UPDATE - 6)) | (1 << (SQLParserT_SET - 6)) | (1 << (SQLParserT_DROP - 6)) | (1 << (SQLParserT_INTERVAL - 6)) | (1 << (SQLParserT_INTERVAL_NAME - 6)) | (1 << (SQLParserT_SHARD - 6)) | (1 << (SQLParserT_REPLICATION - 6)) | (1 << (SQLParserT_TTL - 6)) | (1 << (SQLParserT_META_TTL - 6)) | (1 << (SQLParserT_PAST_TTL - 6)) | (1 << (SQLParserT_FUTURE_TTL - 6)) | (1 << (SQLParserT_KILL - 6)) | (1 << (SQLParserT_ON - 6)) | (1 << (SQLParserT_SHOW - 6)) | (1 << (SQLParserT_USE - 6)) | (1 << (SQLParserT_STATE_REPO - 6)) | (1 << (SQLParserT_STATE_MACHINE - 6)) | (1 << (SQLParserT_MASTER - 6)) | (1 << (SQLParserT_METADATA - 6)) | (1 << (SQLParserT_TYPES - 6)) | (1 << (SQLParserT_TYPE - 6)) | (1 << (SQLParserT_STORAGES - 6)) | (1 << (SQLParserT_STORAGE - 6)) | (1 << (SQLParserT_BROKER - 6)) | (1 << (SQLParserT_ALIVE - 6)) | (1 << (SQLParserT_SCHEMAS - 6)) | (1 << (SQLParserT_DATASBAE - 6)) | (1 << (SQLParserT_DATASBAES - 6)) | (1 << (SQLParserT_NAMESPACE - 6)) | (1 << (SQLParserT_NAMESPACES - 6)) | (1 << (SQLParserT_NODE - 6)))) != 0) || ((((_la - 38)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 38))) & ((1 << (SQLParserT_METRICS - 38)) | (1 << (SQLParserT_METRIC - 38)) | (1 << (SQLParserT_FIELD - 38)) | (1 << (SQLParserT_FIELDS - 38)) | (1 << (SQLParserT_TAG - 38)) | (1 << (SQLParserT_INFO - 38)) | (1 << (SQLParserT_KEYS - 38)) | (1 << (SQLParserT_KEY - 38)) | (1 << (SQLParserT_WITH - 38)) | (1 << (SQLParserT_VALUES - 38)) | (1 << (SQLParserT_VALUE - 38)) | (1 << (SQLParserT_FROM - 38)) | (1 << (SQLParserT_WHERE - 38)) | (1 << (SQLParserT_LIMIT - 38)) | (1 << (SQLParserT_QUERIES - 38)) | (1 << (SQLParserT_QUERY - 38)) | (1 << (SQLParserT_EXPLAIN - 38)) | (1 << (SQLParserT_WITH_VALUE - 38)) | (1 << (SQLParserT_SELECT - 38)) | (1 << (SQLParserT_AS - 38)) | (1 << (SQLParserT_AND - 38)) | (1 << (SQLParserT_OR - 38)) | (1 << (SQLParserT_FILL - 38)) | (1 << (SQLParserT_NULL - 38)) | (1 << (SQLParserT_PREVIOUS - 38)) | (1 << (SQLParserT_ORDER - 38)) | (1 << (SQLParserT_ASC - 38)) | (1 << (SQLParserT_DESC - 38)) | (1 << (SQLParserT_LIKE - 38)) | (1 << (SQLParserT_NOT - 38)) | (1 << (SQLParserT_BETWEEN - 38)) | (1 << (SQLParserT_IS - 38)))) != 0) || ((((_la - 70)) & -(0x1f+1)) == 0 && ((1 << uint((_la - 70))) & ((1 << (SQLParserT_GROUP - 70)) | (1 << (SQLParserT_HAVING - 70)) | (1 << (SQLParserT_BY - 70)) | (1 << (SQLParserT_FOR - 70)) | (1 << (SQLParserT_STATS - 70)) | (1 << (SQLParserT_TIME - 70)) | (1 << (SQLParserT_NOW - 70)) | (1 << (SQLParserT_IN - 70)) | (1 << (SQLParserT_LOG - 70)) | (1 << (SQLParserT_PROFILE - 70)) | (1 << (SQLParserT_SUM - 70)) | (1 << (SQLParserT_MIN - 70)) | (1 << (SQLParserT_MAX - 70)) | (1 << (SQLParserT_COUNT - 70)) | (1 << (SQLParserT_LAST - 70)) | (1 << (SQLParserT_FIRST - 70)) | (1 << (SQLParserT_AVG - 70)) | (1 << (SQLParserT_STDDEV - 70)) | (1 << (SQLParserT_QUANTILE - 70)) | (1 << (SQLParserT_RATE - 70)) | (1 << (SQLParserT_SECOND - 70)) | (1 << (SQLParserT_MINUTE - 70)) | (1 << (SQLParserT_HOUR - 70)) | (1 << (SQLParserT_DAY - 70)) | (1 << (SQLParserT_WEEK - 70)) | (1 << (SQLParserT_MONTH - 70)) | (1 << (SQLParserT_YEAR - 70)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}



	return localctx
}


func (p *SQLParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 40:
			var t *TagFilterExprContext = nil
			if localctx != nil { t = localctx.(*TagFilterExprContext) }
			return p.TagFilterExpr_Sempred(t, predIndex)

	case 56:
			var t *BoolExprContext = nil
			if localctx != nil { t = localctx.(*BoolExprContext) }
			return p.BoolExpr_Sempred(t, predIndex)

	case 61:
			var t *FieldExprContext = nil
			if localctx != nil { t = localctx.(*FieldExprContext) }
			return p.FieldExpr_Sempred(t, predIndex)


	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *SQLParser) TagFilterExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
			return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) BoolExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 1:
			return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) FieldExpr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 2:
			return p.Precpred(p.GetParserRuleContext(), 8)

	case 3:
			return p.Precpred(p.GetParserRuleContext(), 7)

	case 4:
			return p.Precpred(p.GetParserRuleContext(), 6)

	case 5:
			return p.Precpred(p.GetParserRuleContext(), 5)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

