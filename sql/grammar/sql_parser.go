// Code generated from ./sql/grammar/SQLParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar // SQLParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type SQLParser struct {
	*antlr.BaseParser
}

var SQLParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sqlparserParserInit() {
	staticData := &SQLParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "'ALL'", "'ALIVE'", "'AND'", "'ANALYZE'", "'AS'", "'ASC'",
		"'BROKER'", "'BROKERS'", "'BY'", "'COMPACT'", "'CREATE'", "'CROSS'",
		"'DATABASE'", "'DATABASES'", "'DEFAULT'", "'DESC'", "'DISTRIBUTED'",
		"'DROP'", "'ENGINE'", "'ESCAPE'", "'EXPLAIN'", "'EXISTS'", "'FALSE'",
		"'FIELDS'", "'FILTER'", "'FLUSH'", "'FROM'", "'GROUP'", "'HAVING'",
		"'IF'", "'IN'", "'JOIN'", "'KEYS'", "'LEFT'", "'LIKE'", "'LIMIT'", "'LOG'",
		"'LOGICAL'", "'MASTER'", "'METRICS'", "'METRIC'", "'METADATA'", "'METADATAS'",
		"'NAMESPACE'", "'NAMESPACES'", "'NOT'", "'ON'", "'OR'", "'ORDER'", "'PLAN'",
		"'REQUESTS'", "'REPLICATIONS'", "'RIGHT'", "'ROLLUP'", "'SELECT'", "'SHOW'",
		"'STATE'", "'STORAGE'", "'TAG'", "'TRACE'", "'TRUE'", "'TYPE'", "'TYPES'",
		"'VALUES'", "'WHERE'", "'WITH'", "'WITHIN'", "'USING'", "'USE'", "'='",
		"", "'<'", "'<='", "'>'", "'>='", "'+'", "'-'", "'*'", "'/'", "'%'",
		"'=~'", "'!~'", "'!'", "'.'", "'('", "')'", "','",
	}
	staticData.SymbolicNames = []string{
		"", "SIMPLE_COMMENT", "BRACKETED_COMMENT", "WS", "ALL", "ALIVE", "AND",
		"ANALYZE", "AS", "ASC", "BROKER", "BROKERS", "BY", "COMPACT", "CREATE",
		"CROSS", "DATABASE", "DATABASES", "DEFAULT", "DESC", "DISTRIBUTED",
		"DROP", "ENGINE", "ESCAPE", "EXPLAIN", "EXISTS", "FALSE", "FIELDS",
		"FILTER", "FLUSH", "FROM", "GROUP", "HAVING", "IF", "IN", "JOIN", "KEYS",
		"LEFT", "LIKE", "LIMIT", "LOG", "LOGICAL", "MASTER", "METRICS", "METRIC",
		"METADATA", "METADATAS", "NAMESPACE", "NAMESPACES", "NOT", "ON", "OR",
		"ORDER", "PLAN", "REQUESTS", "REPLICATIONS", "RIGHT", "ROLLUP", "SELECT",
		"SHOW", "STATE", "STORAGE", "TAG", "TRACE", "TRUE", "TYPE", "TYPES",
		"VALUES", "WHERE", "WITH", "WITHIN", "USING", "USE", "EQ", "NEQ", "LT",
		"LTE", "GT", "GTE", "PLUS", "MINUS", "ASTERISK", "SLASH", "PERCENT",
		"REGEXP", "NEQREGEXP", "EXCLAMATION_SYMBOL", "DOT", "LR_BRACKET", "RR_BRACKET",
		"COMMA", "STRING", "INTEGER_VALUE", "DECIMAL_VALUE", "DOUBLE_VALUE",
		"IDENTIFIER", "DIGIT_IDENTIFIER", "QUOTED_IDENTIFIER", "BACKQUOTED_IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"statement", "ddlStatement", "dmlStatement", "adminStatement", "utilityStatement",
		"explainOption", "createDatabase", "createDatabaseOptions", "rollupOptions",
		"dropDatabase", "createBroker", "flushDatabase", "compactDatabase",
		"showStatement", "useStatement", "showMetricMetadata", "query", "with",
		"namedQuery", "queryNoWith", "queryTerm", "queryPrimary", "querySpecification",
		"selectItem", "relation", "joinType", "joinCriteria", "aliasedRelation",
		"relationPrimary", "groupBy", "groupingElement", "groupingSet", "having",
		"orderBy", "sortItem", "limitRowCount", "expression", "booleanExpression",
		"valueExpression", "primaryExpression", "predicate", "comparisonOperator",
		"filter", "qualifiedName", "properties", "propertyAssignments", "property",
		"propertyValue", "booleanValue", "string", "identifier", "number", "nonReserved",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 98, 600, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 3, 0, 113, 8, 0, 1, 1, 1, 1, 1,
		1, 3, 1, 118, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 126, 8, 2,
		10, 2, 12, 2, 129, 9, 2, 1, 2, 1, 2, 3, 2, 133, 8, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 3, 2, 139, 8, 2, 1, 3, 1, 3, 1, 3, 3, 3, 144, 8, 3, 1, 4, 1, 4, 1,
		5, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 5, 6, 157, 8, 6, 10,
		6, 12, 6, 160, 9, 6, 3, 6, 162, 8, 6, 1, 6, 1, 6, 3, 6, 166, 8, 6, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 6, 5, 6, 173, 8, 6, 10, 6, 12, 6, 176, 9, 6, 1, 6,
		1, 6, 3, 6, 180, 8, 6, 1, 7, 1, 7, 3, 7, 184, 8, 7, 1, 7, 1, 7, 1, 8, 1,
		8, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 194, 8, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1,
		10, 1, 10, 1, 10, 3, 10, 203, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 12,
		1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1,
		13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13,
		1, 13, 1, 13, 1, 13, 3, 13, 234, 8, 13, 1, 14, 1, 14, 1, 14, 1, 15, 1,
		15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15,
		3, 15, 251, 8, 15, 1, 16, 3, 16, 254, 8, 16, 1, 16, 1, 16, 1, 17, 1, 17,
		1, 17, 1, 17, 5, 17, 262, 8, 17, 10, 17, 12, 17, 265, 9, 17, 1, 18, 1,
		18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 3, 19, 277,
		8, 19, 1, 19, 1, 19, 3, 19, 281, 8, 19, 1, 20, 1, 20, 1, 21, 1, 21, 1,
		21, 1, 21, 1, 21, 3, 21, 290, 8, 21, 1, 22, 1, 22, 1, 22, 1, 22, 5, 22,
		296, 8, 22, 10, 22, 12, 22, 299, 9, 22, 1, 22, 1, 22, 1, 22, 1, 22, 5,
		22, 305, 8, 22, 10, 22, 12, 22, 308, 9, 22, 3, 22, 310, 8, 22, 1, 22, 1,
		22, 3, 22, 314, 8, 22, 1, 22, 1, 22, 1, 22, 3, 22, 319, 8, 22, 1, 22, 1,
		22, 3, 22, 323, 8, 22, 1, 23, 1, 23, 3, 23, 327, 8, 23, 1, 23, 3, 23, 330,
		8, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 337, 8, 23, 1, 24, 1,
		24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24,
		3, 24, 351, 8, 24, 5, 24, 353, 8, 24, 10, 24, 12, 24, 356, 9, 24, 1, 25,
		1, 25, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 5, 26, 367, 8,
		26, 10, 26, 12, 26, 370, 9, 26, 1, 26, 1, 26, 3, 26, 374, 8, 26, 1, 27,
		1, 27, 3, 27, 378, 8, 27, 1, 27, 3, 27, 381, 8, 27, 1, 28, 1, 28, 1, 28,
		1, 28, 1, 28, 3, 28, 388, 8, 28, 1, 29, 1, 29, 1, 29, 5, 29, 393, 8, 29,
		10, 29, 12, 29, 396, 9, 29, 1, 30, 1, 30, 3, 30, 400, 8, 30, 1, 31, 1,
		31, 1, 31, 1, 31, 5, 31, 406, 8, 31, 10, 31, 12, 31, 409, 9, 31, 3, 31,
		411, 8, 31, 1, 31, 1, 31, 3, 31, 415, 8, 31, 1, 32, 1, 32, 1, 33, 1, 33,
		1, 33, 5, 33, 422, 8, 33, 10, 33, 12, 33, 425, 9, 33, 1, 34, 1, 34, 3,
		34, 429, 8, 34, 1, 35, 1, 35, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37,
		3, 37, 439, 8, 37, 1, 37, 1, 37, 1, 37, 1, 37, 1, 37, 1, 37, 5, 37, 447,
		8, 37, 10, 37, 12, 37, 450, 9, 37, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1,
		38, 1, 38, 1, 38, 1, 38, 5, 38, 461, 8, 38, 10, 38, 12, 38, 464, 9, 38,
		1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 5, 39, 475,
		8, 39, 10, 39, 12, 39, 478, 9, 39, 3, 39, 480, 8, 39, 1, 39, 1, 39, 1,
		39, 1, 39, 1, 39, 1, 39, 1, 39, 3, 39, 489, 8, 39, 1, 39, 1, 39, 1, 39,
		5, 39, 494, 8, 39, 10, 39, 12, 39, 497, 9, 39, 1, 40, 1, 40, 1, 40, 1,
		40, 1, 40, 1, 40, 3, 40, 505, 8, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40,
		5, 40, 512, 8, 40, 10, 40, 12, 40, 515, 9, 40, 1, 40, 1, 40, 1, 40, 1,
		40, 3, 40, 521, 8, 40, 1, 40, 1, 40, 1, 40, 1, 40, 3, 40, 527, 8, 40, 1,
		40, 1, 40, 1, 40, 3, 40, 532, 8, 40, 1, 40, 3, 40, 535, 8, 40, 1, 41, 1,
		41, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 5, 43,
		548, 8, 43, 10, 43, 12, 43, 551, 9, 43, 1, 44, 1, 44, 1, 44, 1, 44, 1,
		45, 1, 45, 1, 45, 5, 45, 560, 8, 45, 10, 45, 12, 45, 563, 9, 45, 1, 46,
		1, 46, 1, 46, 1, 46, 1, 47, 1, 47, 3, 47, 571, 8, 47, 1, 48, 1, 48, 1,
		49, 1, 49, 1, 50, 1, 50, 1, 50, 1, 50, 1, 50, 3, 50, 582, 8, 50, 1, 51,
		3, 51, 585, 8, 51, 1, 51, 1, 51, 3, 51, 589, 8, 51, 1, 51, 1, 51, 3, 51,
		593, 8, 51, 1, 51, 3, 51, 596, 8, 51, 1, 52, 1, 52, 1, 52, 0, 4, 48, 74,
		76, 78, 53, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
		32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66,
		68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102,
		104, 0, 11, 2, 0, 20, 20, 41, 41, 3, 0, 40, 40, 44, 44, 63, 63, 2, 0, 37,
		37, 56, 56, 2, 0, 9, 9, 19, 19, 2, 0, 49, 49, 86, 86, 1, 0, 81, 83, 1,
		0, 79, 80, 1, 0, 84, 85, 1, 0, 73, 78, 2, 0, 26, 26, 64, 64, 2, 0, 4, 6,
		8, 72, 641, 0, 112, 1, 0, 0, 0, 2, 117, 1, 0, 0, 0, 4, 138, 1, 0, 0, 0,
		6, 143, 1, 0, 0, 0, 8, 145, 1, 0, 0, 0, 10, 147, 1, 0, 0, 0, 12, 150, 1,
		0, 0, 0, 14, 181, 1, 0, 0, 0, 16, 187, 1, 0, 0, 0, 18, 189, 1, 0, 0, 0,
		20, 197, 1, 0, 0, 0, 22, 204, 1, 0, 0, 0, 24, 208, 1, 0, 0, 0, 26, 233,
		1, 0, 0, 0, 28, 235, 1, 0, 0, 0, 30, 250, 1, 0, 0, 0, 32, 253, 1, 0, 0,
		0, 34, 257, 1, 0, 0, 0, 36, 266, 1, 0, 0, 0, 38, 272, 1, 0, 0, 0, 40, 282,
		1, 0, 0, 0, 42, 289, 1, 0, 0, 0, 44, 291, 1, 0, 0, 0, 46, 336, 1, 0, 0,
		0, 48, 338, 1, 0, 0, 0, 50, 357, 1, 0, 0, 0, 52, 373, 1, 0, 0, 0, 54, 375,
		1, 0, 0, 0, 56, 387, 1, 0, 0, 0, 58, 389, 1, 0, 0, 0, 60, 399, 1, 0, 0,
		0, 62, 414, 1, 0, 0, 0, 64, 416, 1, 0, 0, 0, 66, 418, 1, 0, 0, 0, 68, 426,
		1, 0, 0, 0, 70, 430, 1, 0, 0, 0, 72, 432, 1, 0, 0, 0, 74, 438, 1, 0, 0,
		0, 76, 451, 1, 0, 0, 0, 78, 488, 1, 0, 0, 0, 80, 534, 1, 0, 0, 0, 82, 536,
		1, 0, 0, 0, 84, 538, 1, 0, 0, 0, 86, 544, 1, 0, 0, 0, 88, 552, 1, 0, 0,
		0, 90, 556, 1, 0, 0, 0, 92, 564, 1, 0, 0, 0, 94, 570, 1, 0, 0, 0, 96, 572,
		1, 0, 0, 0, 98, 574, 1, 0, 0, 0, 100, 581, 1, 0, 0, 0, 102, 595, 1, 0,
		0, 0, 104, 597, 1, 0, 0, 0, 106, 113, 3, 2, 1, 0, 107, 113, 3, 4, 2, 0,
		108, 113, 3, 6, 3, 0, 109, 110, 3, 8, 4, 0, 110, 111, 5, 0, 0, 1, 111,
		113, 1, 0, 0, 0, 112, 106, 1, 0, 0, 0, 112, 107, 1, 0, 0, 0, 112, 108,
		1, 0, 0, 0, 112, 109, 1, 0, 0, 0, 113, 1, 1, 0, 0, 0, 114, 118, 3, 12,
		6, 0, 115, 118, 3, 18, 9, 0, 116, 118, 3, 20, 10, 0, 117, 114, 1, 0, 0,
		0, 117, 115, 1, 0, 0, 0, 117, 116, 1, 0, 0, 0, 118, 3, 1, 0, 0, 0, 119,
		139, 3, 32, 16, 0, 120, 132, 5, 24, 0, 0, 121, 122, 5, 88, 0, 0, 122, 127,
		3, 10, 5, 0, 123, 124, 5, 90, 0, 0, 124, 126, 3, 10, 5, 0, 125, 123, 1,
		0, 0, 0, 126, 129, 1, 0, 0, 0, 127, 125, 1, 0, 0, 0, 127, 128, 1, 0, 0,
		0, 128, 130, 1, 0, 0, 0, 129, 127, 1, 0, 0, 0, 130, 131, 5, 89, 0, 0, 131,
		133, 1, 0, 0, 0, 132, 121, 1, 0, 0, 0, 132, 133, 1, 0, 0, 0, 133, 134,
		1, 0, 0, 0, 134, 139, 3, 4, 2, 0, 135, 136, 5, 24, 0, 0, 136, 137, 5, 7,
		0, 0, 137, 139, 3, 4, 2, 0, 138, 119, 1, 0, 0, 0, 138, 120, 1, 0, 0, 0,
		138, 135, 1, 0, 0, 0, 139, 5, 1, 0, 0, 0, 140, 144, 3, 22, 11, 0, 141,
		144, 3, 24, 12, 0, 142, 144, 3, 26, 13, 0, 143, 140, 1, 0, 0, 0, 143, 141,
		1, 0, 0, 0, 143, 142, 1, 0, 0, 0, 144, 7, 1, 0, 0, 0, 145, 146, 3, 28,
		14, 0, 146, 9, 1, 0, 0, 0, 147, 148, 5, 65, 0, 0, 148, 149, 7, 0, 0, 0,
		149, 11, 1, 0, 0, 0, 150, 151, 5, 14, 0, 0, 151, 152, 5, 16, 0, 0, 152,
		161, 3, 86, 43, 0, 153, 158, 3, 14, 7, 0, 154, 155, 5, 90, 0, 0, 155, 157,
		3, 14, 7, 0, 156, 154, 1, 0, 0, 0, 157, 160, 1, 0, 0, 0, 158, 156, 1, 0,
		0, 0, 158, 159, 1, 0, 0, 0, 159, 162, 1, 0, 0, 0, 160, 158, 1, 0, 0, 0,
		161, 153, 1, 0, 0, 0, 161, 162, 1, 0, 0, 0, 162, 165, 1, 0, 0, 0, 163,
		164, 5, 69, 0, 0, 164, 166, 3, 88, 44, 0, 165, 163, 1, 0, 0, 0, 165, 166,
		1, 0, 0, 0, 166, 179, 1, 0, 0, 0, 167, 168, 5, 57, 0, 0, 168, 169, 5, 88,
		0, 0, 169, 174, 3, 16, 8, 0, 170, 171, 5, 90, 0, 0, 171, 173, 3, 16, 8,
		0, 172, 170, 1, 0, 0, 0, 173, 176, 1, 0, 0, 0, 174, 172, 1, 0, 0, 0, 174,
		175, 1, 0, 0, 0, 175, 177, 1, 0, 0, 0, 176, 174, 1, 0, 0, 0, 177, 178,
		5, 89, 0, 0, 178, 180, 1, 0, 0, 0, 179, 167, 1, 0, 0, 0, 179, 180, 1, 0,
		0, 0, 180, 13, 1, 0, 0, 0, 181, 183, 5, 22, 0, 0, 182, 184, 5, 73, 0, 0,
		183, 182, 1, 0, 0, 0, 183, 184, 1, 0, 0, 0, 184, 185, 1, 0, 0, 0, 185,
		186, 7, 1, 0, 0, 186, 15, 1, 0, 0, 0, 187, 188, 3, 88, 44, 0, 188, 17,
		1, 0, 0, 0, 189, 190, 5, 21, 0, 0, 190, 193, 5, 16, 0, 0, 191, 192, 5,
		33, 0, 0, 192, 194, 5, 25, 0, 0, 193, 191, 1, 0, 0, 0, 193, 194, 1, 0,
		0, 0, 194, 195, 1, 0, 0, 0, 195, 196, 3, 86, 43, 0, 196, 19, 1, 0, 0, 0,
		197, 198, 5, 14, 0, 0, 198, 199, 5, 10, 0, 0, 199, 202, 3, 86, 43, 0, 200,
		201, 5, 69, 0, 0, 201, 203, 3, 88, 44, 0, 202, 200, 1, 0, 0, 0, 202, 203,
		1, 0, 0, 0, 203, 21, 1, 0, 0, 0, 204, 205, 5, 29, 0, 0, 205, 206, 5, 16,
		0, 0, 206, 207, 3, 86, 43, 0, 207, 23, 1, 0, 0, 0, 208, 209, 5, 13, 0,
		0, 209, 210, 5, 16, 0, 0, 210, 211, 3, 86, 43, 0, 211, 25, 1, 0, 0, 0,
		212, 213, 5, 59, 0, 0, 213, 234, 5, 42, 0, 0, 214, 215, 5, 59, 0, 0, 215,
		234, 5, 11, 0, 0, 216, 217, 5, 59, 0, 0, 217, 234, 5, 54, 0, 0, 218, 219,
		5, 59, 0, 0, 219, 234, 5, 39, 0, 0, 220, 221, 5, 59, 0, 0, 221, 222, 5,
		45, 0, 0, 222, 234, 5, 66, 0, 0, 223, 224, 5, 59, 0, 0, 224, 234, 5, 46,
		0, 0, 225, 226, 5, 59, 0, 0, 226, 234, 5, 5, 0, 0, 227, 228, 5, 59, 0,
		0, 228, 234, 5, 55, 0, 0, 229, 230, 5, 59, 0, 0, 230, 234, 5, 60, 0, 0,
		231, 232, 5, 59, 0, 0, 232, 234, 5, 17, 0, 0, 233, 212, 1, 0, 0, 0, 233,
		214, 1, 0, 0, 0, 233, 216, 1, 0, 0, 0, 233, 218, 1, 0, 0, 0, 233, 220,
		1, 0, 0, 0, 233, 223, 1, 0, 0, 0, 233, 225, 1, 0, 0, 0, 233, 227, 1, 0,
		0, 0, 233, 229, 1, 0, 0, 0, 233, 231, 1, 0, 0, 0, 234, 27, 1, 0, 0, 0,
		235, 236, 5, 72, 0, 0, 236, 237, 3, 100, 50, 0, 237, 29, 1, 0, 0, 0, 238,
		239, 5, 59, 0, 0, 239, 251, 5, 48, 0, 0, 240, 241, 5, 59, 0, 0, 241, 251,
		5, 43, 0, 0, 242, 243, 5, 59, 0, 0, 243, 251, 5, 27, 0, 0, 244, 245, 5,
		59, 0, 0, 245, 246, 5, 62, 0, 0, 246, 251, 5, 36, 0, 0, 247, 248, 5, 59,
		0, 0, 248, 249, 5, 62, 0, 0, 249, 251, 5, 67, 0, 0, 250, 238, 1, 0, 0,
		0, 250, 240, 1, 0, 0, 0, 250, 242, 1, 0, 0, 0, 250, 244, 1, 0, 0, 0, 250,
		247, 1, 0, 0, 0, 251, 31, 1, 0, 0, 0, 252, 254, 3, 34, 17, 0, 253, 252,
		1, 0, 0, 0, 253, 254, 1, 0, 0, 0, 254, 255, 1, 0, 0, 0, 255, 256, 3, 38,
		19, 0, 256, 33, 1, 0, 0, 0, 257, 258, 5, 69, 0, 0, 258, 263, 3, 36, 18,
		0, 259, 260, 5, 90, 0, 0, 260, 262, 3, 36, 18, 0, 261, 259, 1, 0, 0, 0,
		262, 265, 1, 0, 0, 0, 263, 261, 1, 0, 0, 0, 263, 264, 1, 0, 0, 0, 264,
		35, 1, 0, 0, 0, 265, 263, 1, 0, 0, 0, 266, 267, 3, 100, 50, 0, 267, 268,
		5, 8, 0, 0, 268, 269, 5, 88, 0, 0, 269, 270, 3, 32, 16, 0, 270, 271, 5,
		89, 0, 0, 271, 37, 1, 0, 0, 0, 272, 276, 3, 40, 20, 0, 273, 274, 5, 52,
		0, 0, 274, 275, 5, 12, 0, 0, 275, 277, 3, 66, 33, 0, 276, 273, 1, 0, 0,
		0, 276, 277, 1, 0, 0, 0, 277, 280, 1, 0, 0, 0, 278, 279, 5, 39, 0, 0, 279,
		281, 3, 70, 35, 0, 280, 278, 1, 0, 0, 0, 280, 281, 1, 0, 0, 0, 281, 39,
		1, 0, 0, 0, 282, 283, 3, 42, 21, 0, 283, 41, 1, 0, 0, 0, 284, 290, 3, 44,
		22, 0, 285, 286, 5, 88, 0, 0, 286, 287, 3, 38, 19, 0, 287, 288, 5, 89,
		0, 0, 288, 290, 1, 0, 0, 0, 289, 284, 1, 0, 0, 0, 289, 285, 1, 0, 0, 0,
		290, 43, 1, 0, 0, 0, 291, 292, 5, 58, 0, 0, 292, 297, 3, 46, 23, 0, 293,
		294, 5, 90, 0, 0, 294, 296, 3, 46, 23, 0, 295, 293, 1, 0, 0, 0, 296, 299,
		1, 0, 0, 0, 297, 295, 1, 0, 0, 0, 297, 298, 1, 0, 0, 0, 298, 309, 1, 0,
		0, 0, 299, 297, 1, 0, 0, 0, 300, 301, 5, 30, 0, 0, 301, 306, 3, 48, 24,
		0, 302, 303, 5, 90, 0, 0, 303, 305, 3, 48, 24, 0, 304, 302, 1, 0, 0, 0,
		305, 308, 1, 0, 0, 0, 306, 304, 1, 0, 0, 0, 306, 307, 1, 0, 0, 0, 307,
		310, 1, 0, 0, 0, 308, 306, 1, 0, 0, 0, 309, 300, 1, 0, 0, 0, 309, 310,
		1, 0, 0, 0, 310, 313, 1, 0, 0, 0, 311, 312, 5, 68, 0, 0, 312, 314, 3, 74,
		37, 0, 313, 311, 1, 0, 0, 0, 313, 314, 1, 0, 0, 0, 314, 318, 1, 0, 0, 0,
		315, 316, 5, 31, 0, 0, 316, 317, 5, 12, 0, 0, 317, 319, 3, 58, 29, 0, 318,
		315, 1, 0, 0, 0, 318, 319, 1, 0, 0, 0, 319, 322, 1, 0, 0, 0, 320, 321,
		5, 32, 0, 0, 321, 323, 3, 64, 32, 0, 322, 320, 1, 0, 0, 0, 322, 323, 1,
		0, 0, 0, 323, 45, 1, 0, 0, 0, 324, 329, 3, 72, 36, 0, 325, 327, 5, 8, 0,
		0, 326, 325, 1, 0, 0, 0, 326, 327, 1, 0, 0, 0, 327, 328, 1, 0, 0, 0, 328,
		330, 3, 100, 50, 0, 329, 326, 1, 0, 0, 0, 329, 330, 1, 0, 0, 0, 330, 337,
		1, 0, 0, 0, 331, 332, 3, 78, 39, 0, 332, 333, 5, 87, 0, 0, 333, 334, 5,
		81, 0, 0, 334, 337, 1, 0, 0, 0, 335, 337, 5, 81, 0, 0, 336, 324, 1, 0,
		0, 0, 336, 331, 1, 0, 0, 0, 336, 335, 1, 0, 0, 0, 337, 47, 1, 0, 0, 0,
		338, 339, 6, 24, -1, 0, 339, 340, 3, 54, 27, 0, 340, 354, 1, 0, 0, 0, 341,
		350, 10, 2, 0, 0, 342, 343, 5, 15, 0, 0, 343, 344, 5, 35, 0, 0, 344, 351,
		3, 48, 24, 0, 345, 346, 3, 50, 25, 0, 346, 347, 5, 35, 0, 0, 347, 348,
		3, 48, 24, 0, 348, 349, 3, 52, 26, 0, 349, 351, 1, 0, 0, 0, 350, 342, 1,
		0, 0, 0, 350, 345, 1, 0, 0, 0, 351, 353, 1, 0, 0, 0, 352, 341, 1, 0, 0,
		0, 353, 356, 1, 0, 0, 0, 354, 352, 1, 0, 0, 0, 354, 355, 1, 0, 0, 0, 355,
		49, 1, 0, 0, 0, 356, 354, 1, 0, 0, 0, 357, 358, 7, 2, 0, 0, 358, 51, 1,
		0, 0, 0, 359, 360, 5, 50, 0, 0, 360, 374, 3, 74, 37, 0, 361, 362, 5, 71,
		0, 0, 362, 363, 5, 88, 0, 0, 363, 368, 3, 100, 50, 0, 364, 365, 5, 90,
		0, 0, 365, 367, 3, 100, 50, 0, 366, 364, 1, 0, 0, 0, 367, 370, 1, 0, 0,
		0, 368, 366, 1, 0, 0, 0, 368, 369, 1, 0, 0, 0, 369, 371, 1, 0, 0, 0, 370,
		368, 1, 0, 0, 0, 371, 372, 5, 89, 0, 0, 372, 374, 1, 0, 0, 0, 373, 359,
		1, 0, 0, 0, 373, 361, 1, 0, 0, 0, 374, 53, 1, 0, 0, 0, 375, 380, 3, 56,
		28, 0, 376, 378, 5, 8, 0, 0, 377, 376, 1, 0, 0, 0, 377, 378, 1, 0, 0, 0,
		378, 379, 1, 0, 0, 0, 379, 381, 3, 100, 50, 0, 380, 377, 1, 0, 0, 0, 380,
		381, 1, 0, 0, 0, 381, 55, 1, 0, 0, 0, 382, 388, 3, 86, 43, 0, 383, 384,
		5, 88, 0, 0, 384, 385, 3, 32, 16, 0, 385, 386, 5, 89, 0, 0, 386, 388, 1,
		0, 0, 0, 387, 382, 1, 0, 0, 0, 387, 383, 1, 0, 0, 0, 388, 57, 1, 0, 0,
		0, 389, 394, 3, 60, 30, 0, 390, 391, 5, 90, 0, 0, 391, 393, 3, 60, 30,
		0, 392, 390, 1, 0, 0, 0, 393, 396, 1, 0, 0, 0, 394, 392, 1, 0, 0, 0, 394,
		395, 1, 0, 0, 0, 395, 59, 1, 0, 0, 0, 396, 394, 1, 0, 0, 0, 397, 400, 3,
		62, 31, 0, 398, 400, 5, 4, 0, 0, 399, 397, 1, 0, 0, 0, 399, 398, 1, 0,
		0, 0, 400, 61, 1, 0, 0, 0, 401, 410, 5, 88, 0, 0, 402, 407, 3, 72, 36,
		0, 403, 404, 5, 90, 0, 0, 404, 406, 3, 72, 36, 0, 405, 403, 1, 0, 0, 0,
		406, 409, 1, 0, 0, 0, 407, 405, 1, 0, 0, 0, 407, 408, 1, 0, 0, 0, 408,
		411, 1, 0, 0, 0, 409, 407, 1, 0, 0, 0, 410, 402, 1, 0, 0, 0, 410, 411,
		1, 0, 0, 0, 411, 412, 1, 0, 0, 0, 412, 415, 5, 89, 0, 0, 413, 415, 3, 72,
		36, 0, 414, 401, 1, 0, 0, 0, 414, 413, 1, 0, 0, 0, 415, 63, 1, 0, 0, 0,
		416, 417, 3, 74, 37, 0, 417, 65, 1, 0, 0, 0, 418, 423, 3, 68, 34, 0, 419,
		420, 5, 90, 0, 0, 420, 422, 3, 68, 34, 0, 421, 419, 1, 0, 0, 0, 422, 425,
		1, 0, 0, 0, 423, 421, 1, 0, 0, 0, 423, 424, 1, 0, 0, 0, 424, 67, 1, 0,
		0, 0, 425, 423, 1, 0, 0, 0, 426, 428, 3, 72, 36, 0, 427, 429, 7, 3, 0,
		0, 428, 427, 1, 0, 0, 0, 428, 429, 1, 0, 0, 0, 429, 69, 1, 0, 0, 0, 430,
		431, 5, 92, 0, 0, 431, 71, 1, 0, 0, 0, 432, 433, 3, 74, 37, 0, 433, 73,
		1, 0, 0, 0, 434, 435, 6, 37, -1, 0, 435, 436, 7, 4, 0, 0, 436, 439, 3,
		74, 37, 4, 437, 439, 3, 80, 40, 0, 438, 434, 1, 0, 0, 0, 438, 437, 1, 0,
		0, 0, 439, 448, 1, 0, 0, 0, 440, 441, 10, 3, 0, 0, 441, 442, 5, 6, 0, 0,
		442, 447, 3, 74, 37, 4, 443, 444, 10, 2, 0, 0, 444, 445, 5, 51, 0, 0, 445,
		447, 3, 74, 37, 3, 446, 440, 1, 0, 0, 0, 446, 443, 1, 0, 0, 0, 447, 450,
		1, 0, 0, 0, 448, 446, 1, 0, 0, 0, 448, 449, 1, 0, 0, 0, 449, 75, 1, 0,
		0, 0, 450, 448, 1, 0, 0, 0, 451, 452, 6, 38, -1, 0, 452, 453, 3, 78, 39,
		0, 453, 462, 1, 0, 0, 0, 454, 455, 10, 2, 0, 0, 455, 456, 7, 5, 0, 0, 456,
		461, 3, 76, 38, 3, 457, 458, 10, 1, 0, 0, 458, 459, 7, 6, 0, 0, 459, 461,
		3, 76, 38, 2, 460, 454, 1, 0, 0, 0, 460, 457, 1, 0, 0, 0, 461, 464, 1,
		0, 0, 0, 462, 460, 1, 0, 0, 0, 462, 463, 1, 0, 0, 0, 463, 77, 1, 0, 0,
		0, 464, 462, 1, 0, 0, 0, 465, 466, 6, 39, -1, 0, 466, 489, 3, 102, 51,
		0, 467, 489, 3, 96, 48, 0, 468, 489, 3, 98, 49, 0, 469, 470, 3, 86, 43,
		0, 470, 479, 5, 88, 0, 0, 471, 476, 3, 72, 36, 0, 472, 473, 5, 90, 0, 0,
		473, 475, 3, 72, 36, 0, 474, 472, 1, 0, 0, 0, 475, 478, 1, 0, 0, 0, 476,
		474, 1, 0, 0, 0, 476, 477, 1, 0, 0, 0, 477, 480, 1, 0, 0, 0, 478, 476,
		1, 0, 0, 0, 479, 471, 1, 0, 0, 0, 479, 480, 1, 0, 0, 0, 480, 481, 1, 0,
		0, 0, 481, 482, 5, 89, 0, 0, 482, 489, 1, 0, 0, 0, 483, 489, 3, 100, 50,
		0, 484, 485, 5, 88, 0, 0, 485, 486, 3, 72, 36, 0, 486, 487, 5, 89, 0, 0,
		487, 489, 1, 0, 0, 0, 488, 465, 1, 0, 0, 0, 488, 467, 1, 0, 0, 0, 488,
		468, 1, 0, 0, 0, 488, 469, 1, 0, 0, 0, 488, 483, 1, 0, 0, 0, 488, 484,
		1, 0, 0, 0, 489, 495, 1, 0, 0, 0, 490, 491, 10, 2, 0, 0, 491, 492, 5, 87,
		0, 0, 492, 494, 3, 100, 50, 0, 493, 490, 1, 0, 0, 0, 494, 497, 1, 0, 0,
		0, 495, 493, 1, 0, 0, 0, 495, 496, 1, 0, 0, 0, 496, 79, 1, 0, 0, 0, 497,
		495, 1, 0, 0, 0, 498, 499, 3, 76, 38, 0, 499, 500, 3, 82, 41, 0, 500, 501,
		3, 76, 38, 0, 501, 535, 1, 0, 0, 0, 502, 504, 3, 76, 38, 0, 503, 505, 5,
		49, 0, 0, 504, 503, 1, 0, 0, 0, 504, 505, 1, 0, 0, 0, 505, 506, 1, 0, 0,
		0, 506, 507, 5, 34, 0, 0, 507, 508, 5, 88, 0, 0, 508, 513, 3, 72, 36, 0,
		509, 510, 5, 90, 0, 0, 510, 512, 3, 72, 36, 0, 511, 509, 1, 0, 0, 0, 512,
		515, 1, 0, 0, 0, 513, 511, 1, 0, 0, 0, 513, 514, 1, 0, 0, 0, 514, 516,
		1, 0, 0, 0, 515, 513, 1, 0, 0, 0, 516, 517, 5, 89, 0, 0, 517, 535, 1, 0,
		0, 0, 518, 520, 3, 76, 38, 0, 519, 521, 5, 49, 0, 0, 520, 519, 1, 0, 0,
		0, 520, 521, 1, 0, 0, 0, 521, 522, 1, 0, 0, 0, 522, 523, 5, 38, 0, 0, 523,
		526, 3, 76, 38, 0, 524, 525, 5, 23, 0, 0, 525, 527, 3, 76, 38, 0, 526,
		524, 1, 0, 0, 0, 526, 527, 1, 0, 0, 0, 527, 535, 1, 0, 0, 0, 528, 529,
		3, 76, 38, 0, 529, 531, 7, 7, 0, 0, 530, 532, 3, 76, 38, 0, 531, 530, 1,
		0, 0, 0, 531, 532, 1, 0, 0, 0, 532, 535, 1, 0, 0, 0, 533, 535, 3, 76, 38,
		0, 534, 498, 1, 0, 0, 0, 534, 502, 1, 0, 0, 0, 534, 518, 1, 0, 0, 0, 534,
		528, 1, 0, 0, 0, 534, 533, 1, 0, 0, 0, 535, 81, 1, 0, 0, 0, 536, 537, 7,
		8, 0, 0, 537, 83, 1, 0, 0, 0, 538, 539, 5, 28, 0, 0, 539, 540, 5, 88, 0,
		0, 540, 541, 5, 68, 0, 0, 541, 542, 3, 74, 37, 0, 542, 543, 5, 89, 0, 0,
		543, 85, 1, 0, 0, 0, 544, 549, 3, 100, 50, 0, 545, 546, 5, 87, 0, 0, 546,
		548, 3, 100, 50, 0, 547, 545, 1, 0, 0, 0, 548, 551, 1, 0, 0, 0, 549, 547,
		1, 0, 0, 0, 549, 550, 1, 0, 0, 0, 550, 87, 1, 0, 0, 0, 551, 549, 1, 0,
		0, 0, 552, 553, 5, 88, 0, 0, 553, 554, 3, 90, 45, 0, 554, 555, 5, 89, 0,
		0, 555, 89, 1, 0, 0, 0, 556, 561, 3, 92, 46, 0, 557, 558, 5, 90, 0, 0,
		558, 560, 3, 92, 46, 0, 559, 557, 1, 0, 0, 0, 560, 563, 1, 0, 0, 0, 561,
		559, 1, 0, 0, 0, 561, 562, 1, 0, 0, 0, 562, 91, 1, 0, 0, 0, 563, 561, 1,
		0, 0, 0, 564, 565, 3, 100, 50, 0, 565, 566, 5, 73, 0, 0, 566, 567, 3, 94,
		47, 0, 567, 93, 1, 0, 0, 0, 568, 571, 5, 18, 0, 0, 569, 571, 3, 72, 36,
		0, 570, 568, 1, 0, 0, 0, 570, 569, 1, 0, 0, 0, 571, 95, 1, 0, 0, 0, 572,
		573, 7, 9, 0, 0, 573, 97, 1, 0, 0, 0, 574, 575, 5, 91, 0, 0, 575, 99, 1,
		0, 0, 0, 576, 582, 5, 95, 0, 0, 577, 582, 5, 97, 0, 0, 578, 582, 3, 104,
		52, 0, 579, 582, 5, 98, 0, 0, 580, 582, 5, 96, 0, 0, 581, 576, 1, 0, 0,
		0, 581, 577, 1, 0, 0, 0, 581, 578, 1, 0, 0, 0, 581, 579, 1, 0, 0, 0, 581,
		580, 1, 0, 0, 0, 582, 101, 1, 0, 0, 0, 583, 585, 5, 80, 0, 0, 584, 583,
		1, 0, 0, 0, 584, 585, 1, 0, 0, 0, 585, 586, 1, 0, 0, 0, 586, 596, 5, 93,
		0, 0, 587, 589, 5, 80, 0, 0, 588, 587, 1, 0, 0, 0, 588, 589, 1, 0, 0, 0,
		589, 590, 1, 0, 0, 0, 590, 596, 5, 94, 0, 0, 591, 593, 5, 80, 0, 0, 592,
		591, 1, 0, 0, 0, 592, 593, 1, 0, 0, 0, 593, 594, 1, 0, 0, 0, 594, 596,
		5, 92, 0, 0, 595, 584, 1, 0, 0, 0, 595, 588, 1, 0, 0, 0, 595, 592, 1, 0,
		0, 0, 596, 103, 1, 0, 0, 0, 597, 598, 7, 10, 0, 0, 598, 105, 1, 0, 0, 0,
		67, 112, 117, 127, 132, 138, 143, 158, 161, 165, 174, 179, 183, 193, 202,
		233, 250, 253, 263, 276, 280, 289, 297, 306, 309, 313, 318, 322, 326, 329,
		336, 350, 354, 368, 373, 377, 380, 387, 394, 399, 407, 410, 414, 423, 428,
		438, 446, 448, 460, 462, 476, 479, 488, 495, 504, 513, 520, 526, 531, 534,
		549, 561, 570, 581, 584, 588, 592, 595,
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
	staticData := &SQLParserParserStaticData
	staticData.once.Do(sqlparserParserInit)
}

// NewSQLParser produces a new parser instance for the optional input antlr.TokenStream.
func NewSQLParser(input antlr.TokenStream) *SQLParser {
	SQLParserInit()
	this := new(SQLParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &SQLParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "SQLParser.g4"

	return this
}

// SQLParser tokens.
const (
	SQLParserEOF                   = antlr.TokenEOF
	SQLParserSIMPLE_COMMENT        = 1
	SQLParserBRACKETED_COMMENT     = 2
	SQLParserWS                    = 3
	SQLParserALL                   = 4
	SQLParserALIVE                 = 5
	SQLParserAND                   = 6
	SQLParserANALYZE               = 7
	SQLParserAS                    = 8
	SQLParserASC                   = 9
	SQLParserBROKER                = 10
	SQLParserBROKERS               = 11
	SQLParserBY                    = 12
	SQLParserCOMPACT               = 13
	SQLParserCREATE                = 14
	SQLParserCROSS                 = 15
	SQLParserDATABASE              = 16
	SQLParserDATABASES             = 17
	SQLParserDEFAULT               = 18
	SQLParserDESC                  = 19
	SQLParserDISTRIBUTED           = 20
	SQLParserDROP                  = 21
	SQLParserENGINE                = 22
	SQLParserESCAPE                = 23
	SQLParserEXPLAIN               = 24
	SQLParserEXISTS                = 25
	SQLParserFALSE                 = 26
	SQLParserFIELDS                = 27
	SQLParserFILTER                = 28
	SQLParserFLUSH                 = 29
	SQLParserFROM                  = 30
	SQLParserGROUP                 = 31
	SQLParserHAVING                = 32
	SQLParserIF                    = 33
	SQLParserIN                    = 34
	SQLParserJOIN                  = 35
	SQLParserKEYS                  = 36
	SQLParserLEFT                  = 37
	SQLParserLIKE                  = 38
	SQLParserLIMIT                 = 39
	SQLParserLOG                   = 40
	SQLParserLOGICAL               = 41
	SQLParserMASTER                = 42
	SQLParserMETRICS               = 43
	SQLParserMETRIC                = 44
	SQLParserMETADATA              = 45
	SQLParserMETADATAS             = 46
	SQLParserNAMESPACE             = 47
	SQLParserNAMESPACES            = 48
	SQLParserNOT                   = 49
	SQLParserON                    = 50
	SQLParserOR                    = 51
	SQLParserORDER                 = 52
	SQLParserPLAN                  = 53
	SQLParserREQUESTS              = 54
	SQLParserREPLICATIONS          = 55
	SQLParserRIGHT                 = 56
	SQLParserROLLUP                = 57
	SQLParserSELECT                = 58
	SQLParserSHOW                  = 59
	SQLParserSTATE                 = 60
	SQLParserSTORAGE               = 61
	SQLParserTAG                   = 62
	SQLParserTRACE                 = 63
	SQLParserTRUE                  = 64
	SQLParserTYPE                  = 65
	SQLParserTYPES                 = 66
	SQLParserVALUES                = 67
	SQLParserWHERE                 = 68
	SQLParserWITH                  = 69
	SQLParserWITHIN                = 70
	SQLParserUSING                 = 71
	SQLParserUSE                   = 72
	SQLParserEQ                    = 73
	SQLParserNEQ                   = 74
	SQLParserLT                    = 75
	SQLParserLTE                   = 76
	SQLParserGT                    = 77
	SQLParserGTE                   = 78
	SQLParserPLUS                  = 79
	SQLParserMINUS                 = 80
	SQLParserASTERISK              = 81
	SQLParserSLASH                 = 82
	SQLParserPERCENT               = 83
	SQLParserREGEXP                = 84
	SQLParserNEQREGEXP             = 85
	SQLParserEXCLAMATION_SYMBOL    = 86
	SQLParserDOT                   = 87
	SQLParserLR_BRACKET            = 88
	SQLParserRR_BRACKET            = 89
	SQLParserCOMMA                 = 90
	SQLParserSTRING                = 91
	SQLParserINTEGER_VALUE         = 92
	SQLParserDECIMAL_VALUE         = 93
	SQLParserDOUBLE_VALUE          = 94
	SQLParserIDENTIFIER            = 95
	SQLParserDIGIT_IDENTIFIER      = 96
	SQLParserQUOTED_IDENTIFIER     = 97
	SQLParserBACKQUOTED_IDENTIFIER = 98
)

// SQLParser rules.
const (
	SQLParserRULE_statement             = 0
	SQLParserRULE_ddlStatement          = 1
	SQLParserRULE_dmlStatement          = 2
	SQLParserRULE_adminStatement        = 3
	SQLParserRULE_utilityStatement      = 4
	SQLParserRULE_explainOption         = 5
	SQLParserRULE_createDatabase        = 6
	SQLParserRULE_createDatabaseOptions = 7
	SQLParserRULE_rollupOptions         = 8
	SQLParserRULE_dropDatabase          = 9
	SQLParserRULE_createBroker          = 10
	SQLParserRULE_flushDatabase         = 11
	SQLParserRULE_compactDatabase       = 12
	SQLParserRULE_showStatement         = 13
	SQLParserRULE_useStatement          = 14
	SQLParserRULE_showMetricMetadata    = 15
	SQLParserRULE_query                 = 16
	SQLParserRULE_with                  = 17
	SQLParserRULE_namedQuery            = 18
	SQLParserRULE_queryNoWith           = 19
	SQLParserRULE_queryTerm             = 20
	SQLParserRULE_queryPrimary          = 21
	SQLParserRULE_querySpecification    = 22
	SQLParserRULE_selectItem            = 23
	SQLParserRULE_relation              = 24
	SQLParserRULE_joinType              = 25
	SQLParserRULE_joinCriteria          = 26
	SQLParserRULE_aliasedRelation       = 27
	SQLParserRULE_relationPrimary       = 28
	SQLParserRULE_groupBy               = 29
	SQLParserRULE_groupingElement       = 30
	SQLParserRULE_groupingSet           = 31
	SQLParserRULE_having                = 32
	SQLParserRULE_orderBy               = 33
	SQLParserRULE_sortItem              = 34
	SQLParserRULE_limitRowCount         = 35
	SQLParserRULE_expression            = 36
	SQLParserRULE_booleanExpression     = 37
	SQLParserRULE_valueExpression       = 38
	SQLParserRULE_primaryExpression     = 39
	SQLParserRULE_predicate             = 40
	SQLParserRULE_comparisonOperator    = 41
	SQLParserRULE_filter                = 42
	SQLParserRULE_qualifiedName         = 43
	SQLParserRULE_properties            = 44
	SQLParserRULE_propertyAssignments   = 45
	SQLParserRULE_property              = 46
	SQLParserRULE_propertyValue         = 47
	SQLParserRULE_booleanValue          = 48
	SQLParserRULE_string                = 49
	SQLParserRULE_identifier            = 50
	SQLParserRULE_number                = 51
	SQLParserRULE_nonReserved           = 52
)

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DdlStatement() IDdlStatementContext
	DmlStatement() IDmlStatementContext
	AdminStatement() IAdminStatementContext
	UtilityStatement() IUtilityStatementContext
	EOF() antlr.TerminalNode

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_statement
	return p
}

func InitEmptyStatementContext(p *StatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_statement
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) DdlStatement() IDdlStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDdlStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDdlStatementContext)
}

func (s *StatementContext) DmlStatement() IDmlStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDmlStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDmlStatementContext)
}

func (s *StatementContext) AdminStatement() IAdminStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAdminStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAdminStatementContext)
}

func (s *StatementContext) UtilityStatement() IUtilityStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUtilityStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUtilityStatementContext)
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
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SQLParserRULE_statement)
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserCREATE, SQLParserDROP:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(106)
			p.DdlStatement()
		}

	case SQLParserEXPLAIN, SQLParserSELECT, SQLParserWITH, SQLParserLR_BRACKET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(107)
			p.DmlStatement()
		}

	case SQLParserCOMPACT, SQLParserFLUSH, SQLParserSHOW:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(108)
			p.AdminStatement()
		}

	case SQLParserUSE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(109)
			p.UtilityStatement()
		}
		{
			p.SetState(110)
			p.Match(SQLParserEOF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDdlStatementContext is an interface to support dynamic dispatch.
type IDdlStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CreateDatabase() ICreateDatabaseContext
	DropDatabase() IDropDatabaseContext
	CreateBroker() ICreateBrokerContext

	// IsDdlStatementContext differentiates from other interfaces.
	IsDdlStatementContext()
}

type DdlStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDdlStatementContext() *DdlStatementContext {
	var p = new(DdlStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_ddlStatement
	return p
}

func InitEmptyDdlStatementContext(p *DdlStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_ddlStatement
}

func (*DdlStatementContext) IsDdlStatementContext() {}

func NewDdlStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DdlStatementContext {
	var p = new(DdlStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_ddlStatement

	return p
}

func (s *DdlStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DdlStatementContext) CreateDatabase() ICreateDatabaseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateDatabaseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateDatabaseContext)
}

func (s *DdlStatementContext) DropDatabase() IDropDatabaseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDropDatabaseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDropDatabaseContext)
}

func (s *DdlStatementContext) CreateBroker() ICreateBrokerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateBrokerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateBrokerContext)
}

func (s *DdlStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DdlStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DdlStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDdlStatement(s)
	}
}

func (s *DdlStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDdlStatement(s)
	}
}

func (s *DdlStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDdlStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DdlStatement() (localctx IDdlStatementContext) {
	localctx = NewDdlStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SQLParserRULE_ddlStatement)
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(114)
			p.CreateDatabase()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(115)
			p.DropDatabase()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(116)
			p.CreateBroker()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDmlStatementContext is an interface to support dynamic dispatch.
type IDmlStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsDmlStatementContext differentiates from other interfaces.
	IsDmlStatementContext()
}

type DmlStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDmlStatementContext() *DmlStatementContext {
	var p = new(DmlStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_dmlStatement
	return p
}

func InitEmptyDmlStatementContext(p *DmlStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_dmlStatement
}

func (*DmlStatementContext) IsDmlStatementContext() {}

func NewDmlStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DmlStatementContext {
	var p = new(DmlStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dmlStatement

	return p
}

func (s *DmlStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DmlStatementContext) CopyAll(ctx *DmlStatementContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *DmlStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DmlStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ExplainContext struct {
	DmlStatementContext
}

func NewExplainContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExplainContext {
	var p = new(ExplainContext)

	InitEmptyDmlStatementContext(&p.DmlStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*DmlStatementContext))

	return p
}

func (s *ExplainContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExplainContext) EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserEXPLAIN, 0)
}

func (s *ExplainContext) DmlStatement() IDmlStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDmlStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDmlStatementContext)
}

func (s *ExplainContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *ExplainContext) AllExplainOption() []IExplainOptionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExplainOptionContext); ok {
			len++
		}
	}

	tst := make([]IExplainOptionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExplainOptionContext); ok {
			tst[i] = t.(IExplainOptionContext)
			i++
		}
	}

	return tst
}

func (s *ExplainContext) ExplainOption(i int) IExplainOptionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExplainOptionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExplainOptionContext)
}

func (s *ExplainContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *ExplainContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *ExplainContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *ExplainContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterExplain(s)
	}
}

func (s *ExplainContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitExplain(s)
	}
}

func (s *ExplainContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitExplain(s)

	default:
		return t.VisitChildren(s)
	}
}

type StatementDefaultContext struct {
	DmlStatementContext
}

func NewStatementDefaultContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StatementDefaultContext {
	var p = new(StatementDefaultContext)

	InitEmptyDmlStatementContext(&p.DmlStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*DmlStatementContext))

	return p
}

func (s *StatementDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementDefaultContext) Query() IQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryContext)
}

func (s *StatementDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterStatementDefault(s)
	}
}

func (s *StatementDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitStatementDefault(s)
	}
}

func (s *StatementDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitStatementDefault(s)

	default:
		return t.VisitChildren(s)
	}
}

type ExplainAnalyzeContext struct {
	DmlStatementContext
}

func NewExplainAnalyzeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExplainAnalyzeContext {
	var p = new(ExplainAnalyzeContext)

	InitEmptyDmlStatementContext(&p.DmlStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*DmlStatementContext))

	return p
}

func (s *ExplainAnalyzeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExplainAnalyzeContext) EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserEXPLAIN, 0)
}

func (s *ExplainAnalyzeContext) ANALYZE() antlr.TerminalNode {
	return s.GetToken(SQLParserANALYZE, 0)
}

func (s *ExplainAnalyzeContext) DmlStatement() IDmlStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDmlStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDmlStatementContext)
}

func (s *ExplainAnalyzeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterExplainAnalyze(s)
	}
}

func (s *ExplainAnalyzeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitExplainAnalyze(s)
	}
}

func (s *ExplainAnalyzeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitExplainAnalyze(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DmlStatement() (localctx IDmlStatementContext) {
	localctx = NewDmlStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_dmlStatement)
	var _la int

	p.SetState(138)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		localctx = NewStatementDefaultContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(119)
			p.Query()
		}

	case 2:
		localctx = NewExplainContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(120)
			p.Match(SQLParserEXPLAIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(132)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(121)
				p.Match(SQLParserLR_BRACKET)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(122)
				p.ExplainOption()
			}
			p.SetState(127)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == SQLParserCOMMA {
				{
					p.SetState(123)
					p.Match(SQLParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(124)
					p.ExplainOption()
				}

				p.SetState(129)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(130)
				p.Match(SQLParserRR_BRACKET)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}
		{
			p.SetState(134)
			p.DmlStatement()
		}

	case 3:
		localctx = NewExplainAnalyzeContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(135)
			p.Match(SQLParserEXPLAIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(136)
			p.Match(SQLParserANALYZE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(137)
			p.DmlStatement()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAdminStatementContext is an interface to support dynamic dispatch.
type IAdminStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FlushDatabase() IFlushDatabaseContext
	CompactDatabase() ICompactDatabaseContext
	ShowStatement() IShowStatementContext

	// IsAdminStatementContext differentiates from other interfaces.
	IsAdminStatementContext()
}

type AdminStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAdminStatementContext() *AdminStatementContext {
	var p = new(AdminStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_adminStatement
	return p
}

func InitEmptyAdminStatementContext(p *AdminStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_adminStatement
}

func (*AdminStatementContext) IsAdminStatementContext() {}

func NewAdminStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AdminStatementContext {
	var p = new(AdminStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_adminStatement

	return p
}

func (s *AdminStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *AdminStatementContext) FlushDatabase() IFlushDatabaseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFlushDatabaseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFlushDatabaseContext)
}

func (s *AdminStatementContext) CompactDatabase() ICompactDatabaseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompactDatabaseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompactDatabaseContext)
}

func (s *AdminStatementContext) ShowStatement() IShowStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowStatementContext)
}

func (s *AdminStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdminStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AdminStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterAdminStatement(s)
	}
}

func (s *AdminStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitAdminStatement(s)
	}
}

func (s *AdminStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitAdminStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) AdminStatement() (localctx IAdminStatementContext) {
	localctx = NewAdminStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SQLParserRULE_adminStatement)
	p.SetState(143)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserFLUSH:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(140)
			p.FlushDatabase()
		}

	case SQLParserCOMPACT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(141)
			p.CompactDatabase()
		}

	case SQLParserSHOW:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(142)
			p.ShowStatement()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUtilityStatementContext is an interface to support dynamic dispatch.
type IUtilityStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UseStatement() IUseStatementContext

	// IsUtilityStatementContext differentiates from other interfaces.
	IsUtilityStatementContext()
}

type UtilityStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUtilityStatementContext() *UtilityStatementContext {
	var p = new(UtilityStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_utilityStatement
	return p
}

func InitEmptyUtilityStatementContext(p *UtilityStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_utilityStatement
}

func (*UtilityStatementContext) IsUtilityStatementContext() {}

func NewUtilityStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UtilityStatementContext {
	var p = new(UtilityStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_utilityStatement

	return p
}

func (s *UtilityStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *UtilityStatementContext) UseStatement() IUseStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUseStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUseStatementContext)
}

func (s *UtilityStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UtilityStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UtilityStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterUtilityStatement(s)
	}
}

func (s *UtilityStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitUtilityStatement(s)
	}
}

func (s *UtilityStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitUtilityStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) UtilityStatement() (localctx IUtilityStatementContext) {
	localctx = NewUtilityStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SQLParserRULE_utilityStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(145)
		p.UseStatement()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExplainOptionContext is an interface to support dynamic dispatch.
type IExplainOptionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsExplainOptionContext differentiates from other interfaces.
	IsExplainOptionContext()
}

type ExplainOptionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExplainOptionContext() *ExplainOptionContext {
	var p = new(ExplainOptionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_explainOption
	return p
}

func InitEmptyExplainOptionContext(p *ExplainOptionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_explainOption
}

func (*ExplainOptionContext) IsExplainOptionContext() {}

func NewExplainOptionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExplainOptionContext {
	var p = new(ExplainOptionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_explainOption

	return p
}

func (s *ExplainOptionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExplainOptionContext) CopyAll(ctx *ExplainOptionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ExplainOptionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExplainOptionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ExplainTypeContext struct {
	ExplainOptionContext
	value antlr.Token
}

func NewExplainTypeContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExplainTypeContext {
	var p = new(ExplainTypeContext)

	InitEmptyExplainOptionContext(&p.ExplainOptionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExplainOptionContext))

	return p
}

func (s *ExplainTypeContext) GetValue() antlr.Token { return s.value }

func (s *ExplainTypeContext) SetValue(v antlr.Token) { s.value = v }

func (s *ExplainTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExplainTypeContext) TYPE() antlr.TerminalNode {
	return s.GetToken(SQLParserTYPE, 0)
}

func (s *ExplainTypeContext) LOGICAL() antlr.TerminalNode {
	return s.GetToken(SQLParserLOGICAL, 0)
}

func (s *ExplainTypeContext) DISTRIBUTED() antlr.TerminalNode {
	return s.GetToken(SQLParserDISTRIBUTED, 0)
}

func (s *ExplainTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterExplainType(s)
	}
}

func (s *ExplainTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitExplainType(s)
	}
}

func (s *ExplainTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitExplainType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ExplainOption() (localctx IExplainOptionContext) {
	localctx = NewExplainOptionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SQLParserRULE_explainOption)
	var _la int

	localctx = NewExplainTypeContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(147)
		p.Match(SQLParserTYPE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(148)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*ExplainTypeContext).value = _lt

		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserDISTRIBUTED || _la == SQLParserLOGICAL) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*ExplainTypeContext).value = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateDatabaseContext is an interface to support dynamic dispatch.
type ICreateDatabaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name rule contexts.
	GetName() IQualifiedNameContext

	// SetName sets the name rule contexts.
	SetName(IQualifiedNameContext)

	// Getter signatures
	CREATE() antlr.TerminalNode
	DATABASE() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	AllCreateDatabaseOptions() []ICreateDatabaseOptionsContext
	CreateDatabaseOptions(i int) ICreateDatabaseOptionsContext
	WITH() antlr.TerminalNode
	Properties() IPropertiesContext
	ROLLUP() antlr.TerminalNode
	LR_BRACKET() antlr.TerminalNode
	AllRollupOptions() []IRollupOptionsContext
	RollupOptions(i int) IRollupOptionsContext
	RR_BRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCreateDatabaseContext differentiates from other interfaces.
	IsCreateDatabaseContext()
}

type CreateDatabaseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	name   IQualifiedNameContext
}

func NewEmptyCreateDatabaseContext() *CreateDatabaseContext {
	var p = new(CreateDatabaseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createDatabase
	return p
}

func InitEmptyCreateDatabaseContext(p *CreateDatabaseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createDatabase
}

func (*CreateDatabaseContext) IsCreateDatabaseContext() {}

func NewCreateDatabaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateDatabaseContext {
	var p = new(CreateDatabaseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createDatabase

	return p
}

func (s *CreateDatabaseContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateDatabaseContext) GetName() IQualifiedNameContext { return s.name }

func (s *CreateDatabaseContext) SetName(v IQualifiedNameContext) { s.name = v }

func (s *CreateDatabaseContext) CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserCREATE, 0)
}

func (s *CreateDatabaseContext) DATABASE() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASE, 0)
}

func (s *CreateDatabaseContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *CreateDatabaseContext) AllCreateDatabaseOptions() []ICreateDatabaseOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICreateDatabaseOptionsContext); ok {
			len++
		}
	}

	tst := make([]ICreateDatabaseOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICreateDatabaseOptionsContext); ok {
			tst[i] = t.(ICreateDatabaseOptionsContext)
			i++
		}
	}

	return tst
}

func (s *CreateDatabaseContext) CreateDatabaseOptions(i int) ICreateDatabaseOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateDatabaseOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateDatabaseOptionsContext)
}

func (s *CreateDatabaseContext) WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserWITH, 0)
}

func (s *CreateDatabaseContext) Properties() IPropertiesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertiesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertiesContext)
}

func (s *CreateDatabaseContext) ROLLUP() antlr.TerminalNode {
	return s.GetToken(SQLParserROLLUP, 0)
}

func (s *CreateDatabaseContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *CreateDatabaseContext) AllRollupOptions() []IRollupOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRollupOptionsContext); ok {
			len++
		}
	}

	tst := make([]IRollupOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRollupOptionsContext); ok {
			tst[i] = t.(IRollupOptionsContext)
			i++
		}
	}

	return tst
}

func (s *CreateDatabaseContext) RollupOptions(i int) IRollupOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRollupOptionsContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRollupOptionsContext)
}

func (s *CreateDatabaseContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *CreateDatabaseContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *CreateDatabaseContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *CreateDatabaseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateDatabaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateDatabaseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterCreateDatabase(s)
	}
}

func (s *CreateDatabaseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitCreateDatabase(s)
	}
}

func (s *CreateDatabaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitCreateDatabase(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) CreateDatabase() (localctx ICreateDatabaseContext) {
	localctx = NewCreateDatabaseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SQLParserRULE_createDatabase)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(150)
		p.Match(SQLParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(151)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(152)

		var _x = p.QualifiedName()

		localctx.(*CreateDatabaseContext).name = _x
	}
	p.SetState(161)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserENGINE {
		{
			p.SetState(153)
			p.CreateDatabaseOptions()
		}
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(154)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(155)
				p.CreateDatabaseOptions()
			}

			p.SetState(160)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWITH {
		{
			p.SetState(163)
			p.Match(SQLParserWITH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(164)
			p.Properties()
		}

	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserROLLUP {
		{
			p.SetState(167)
			p.Match(SQLParserROLLUP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(168)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(169)
			p.RollupOptions()
		}
		p.SetState(174)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(170)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(171)
				p.RollupOptions()
			}

			p.SetState(176)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(177)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateDatabaseOptionsContext is an interface to support dynamic dispatch.
type ICreateDatabaseOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsCreateDatabaseOptionsContext differentiates from other interfaces.
	IsCreateDatabaseOptionsContext()
}

type CreateDatabaseOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateDatabaseOptionsContext() *CreateDatabaseOptionsContext {
	var p = new(CreateDatabaseOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createDatabaseOptions
	return p
}

func InitEmptyCreateDatabaseOptionsContext(p *CreateDatabaseOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createDatabaseOptions
}

func (*CreateDatabaseOptionsContext) IsCreateDatabaseOptionsContext() {}

func NewCreateDatabaseOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateDatabaseOptionsContext {
	var p = new(CreateDatabaseOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createDatabaseOptions

	return p
}

func (s *CreateDatabaseOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateDatabaseOptionsContext) CopyAll(ctx *CreateDatabaseOptionsContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *CreateDatabaseOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateDatabaseOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type EngineOptionContext struct {
	CreateDatabaseOptionsContext
	value antlr.Token
}

func NewEngineOptionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EngineOptionContext {
	var p = new(EngineOptionContext)

	InitEmptyCreateDatabaseOptionsContext(&p.CreateDatabaseOptionsContext)
	p.parser = parser
	p.CopyAll(ctx.(*CreateDatabaseOptionsContext))

	return p
}

func (s *EngineOptionContext) GetValue() antlr.Token { return s.value }

func (s *EngineOptionContext) SetValue(v antlr.Token) { s.value = v }

func (s *EngineOptionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EngineOptionContext) ENGINE() antlr.TerminalNode {
	return s.GetToken(SQLParserENGINE, 0)
}

func (s *EngineOptionContext) METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserMETRIC, 0)
}

func (s *EngineOptionContext) LOG() antlr.TerminalNode {
	return s.GetToken(SQLParserLOG, 0)
}

func (s *EngineOptionContext) TRACE() antlr.TerminalNode {
	return s.GetToken(SQLParserTRACE, 0)
}

func (s *EngineOptionContext) EQ() antlr.TerminalNode {
	return s.GetToken(SQLParserEQ, 0)
}

func (s *EngineOptionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterEngineOption(s)
	}
}

func (s *EngineOptionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitEngineOption(s)
	}
}

func (s *EngineOptionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitEngineOption(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) CreateDatabaseOptions() (localctx ICreateDatabaseOptionsContext) {
	localctx = NewCreateDatabaseOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SQLParserRULE_createDatabaseOptions)
	var _la int

	localctx = NewEngineOptionContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(181)
		p.Match(SQLParserENGINE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserEQ {
		{
			p.SetState(182)
			p.Match(SQLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(185)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*EngineOptionContext).value = _lt

		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-9223353345157103616) != 0) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*EngineOptionContext).value = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRollupOptionsContext is an interface to support dynamic dispatch.
type IRollupOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Properties() IPropertiesContext

	// IsRollupOptionsContext differentiates from other interfaces.
	IsRollupOptionsContext()
}

type RollupOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRollupOptionsContext() *RollupOptionsContext {
	var p = new(RollupOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_rollupOptions
	return p
}

func InitEmptyRollupOptionsContext(p *RollupOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_rollupOptions
}

func (*RollupOptionsContext) IsRollupOptionsContext() {}

func NewRollupOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RollupOptionsContext {
	var p = new(RollupOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_rollupOptions

	return p
}

func (s *RollupOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *RollupOptionsContext) Properties() IPropertiesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertiesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertiesContext)
}

func (s *RollupOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RollupOptionsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterRollupOptions(s)
	}
}

func (s *RollupOptionsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitRollupOptions(s)
	}
}

func (s *RollupOptionsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitRollupOptions(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) RollupOptions() (localctx IRollupOptionsContext) {
	localctx = NewRollupOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SQLParserRULE_rollupOptions)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(187)
		p.Properties()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDropDatabaseContext is an interface to support dynamic dispatch.
type IDropDatabaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDatabase returns the database rule contexts.
	GetDatabase() IQualifiedNameContext

	// SetDatabase sets the database rule contexts.
	SetDatabase(IQualifiedNameContext)

	// Getter signatures
	DROP() antlr.TerminalNode
	DATABASE() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	IF() antlr.TerminalNode
	EXISTS() antlr.TerminalNode

	// IsDropDatabaseContext differentiates from other interfaces.
	IsDropDatabaseContext()
}

type DropDatabaseContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	database IQualifiedNameContext
}

func NewEmptyDropDatabaseContext() *DropDatabaseContext {
	var p = new(DropDatabaseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_dropDatabase
	return p
}

func InitEmptyDropDatabaseContext(p *DropDatabaseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_dropDatabase
}

func (*DropDatabaseContext) IsDropDatabaseContext() {}

func NewDropDatabaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DropDatabaseContext {
	var p = new(DropDatabaseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_dropDatabase

	return p
}

func (s *DropDatabaseContext) GetParser() antlr.Parser { return s.parser }

func (s *DropDatabaseContext) GetDatabase() IQualifiedNameContext { return s.database }

func (s *DropDatabaseContext) SetDatabase(v IQualifiedNameContext) { s.database = v }

func (s *DropDatabaseContext) DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserDROP, 0)
}

func (s *DropDatabaseContext) DATABASE() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASE, 0)
}

func (s *DropDatabaseContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *DropDatabaseContext) IF() antlr.TerminalNode {
	return s.GetToken(SQLParserIF, 0)
}

func (s *DropDatabaseContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(SQLParserEXISTS, 0)
}

func (s *DropDatabaseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DropDatabaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DropDatabaseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDropDatabase(s)
	}
}

func (s *DropDatabaseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDropDatabase(s)
	}
}

func (s *DropDatabaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDropDatabase(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DropDatabase() (localctx IDropDatabaseContext) {
	localctx = NewDropDatabaseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_dropDatabase)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(189)
		p.Match(SQLParserDROP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(190)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(193)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 12, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(191)
			p.Match(SQLParserIF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(192)
			p.Match(SQLParserEXISTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}
	{
		p.SetState(195)

		var _x = p.QualifiedName()

		localctx.(*DropDatabaseContext).database = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateBrokerContext is an interface to support dynamic dispatch.
type ICreateBrokerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name rule contexts.
	GetName() IQualifiedNameContext

	// SetName sets the name rule contexts.
	SetName(IQualifiedNameContext)

	// Getter signatures
	CREATE() antlr.TerminalNode
	BROKER() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	WITH() antlr.TerminalNode
	Properties() IPropertiesContext

	// IsCreateBrokerContext differentiates from other interfaces.
	IsCreateBrokerContext()
}

type CreateBrokerContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	name   IQualifiedNameContext
}

func NewEmptyCreateBrokerContext() *CreateBrokerContext {
	var p = new(CreateBrokerContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createBroker
	return p
}

func InitEmptyCreateBrokerContext(p *CreateBrokerContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_createBroker
}

func (*CreateBrokerContext) IsCreateBrokerContext() {}

func NewCreateBrokerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateBrokerContext {
	var p = new(CreateBrokerContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_createBroker

	return p
}

func (s *CreateBrokerContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateBrokerContext) GetName() IQualifiedNameContext { return s.name }

func (s *CreateBrokerContext) SetName(v IQualifiedNameContext) { s.name = v }

func (s *CreateBrokerContext) CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserCREATE, 0)
}

func (s *CreateBrokerContext) BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserBROKER, 0)
}

func (s *CreateBrokerContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *CreateBrokerContext) WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserWITH, 0)
}

func (s *CreateBrokerContext) Properties() IPropertiesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertiesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertiesContext)
}

func (s *CreateBrokerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateBrokerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateBrokerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterCreateBroker(s)
	}
}

func (s *CreateBrokerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitCreateBroker(s)
	}
}

func (s *CreateBrokerContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitCreateBroker(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) CreateBroker() (localctx ICreateBrokerContext) {
	localctx = NewCreateBrokerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SQLParserRULE_createBroker)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(197)
		p.Match(SQLParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(198)
		p.Match(SQLParserBROKER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(199)

		var _x = p.QualifiedName()

		localctx.(*CreateBrokerContext).name = _x
	}
	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWITH {
		{
			p.SetState(200)
			p.Match(SQLParserWITH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(201)
			p.Properties()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFlushDatabaseContext is an interface to support dynamic dispatch.
type IFlushDatabaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDatabase returns the database rule contexts.
	GetDatabase() IQualifiedNameContext

	// SetDatabase sets the database rule contexts.
	SetDatabase(IQualifiedNameContext)

	// Getter signatures
	FLUSH() antlr.TerminalNode
	DATABASE() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext

	// IsFlushDatabaseContext differentiates from other interfaces.
	IsFlushDatabaseContext()
}

type FlushDatabaseContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	database IQualifiedNameContext
}

func NewEmptyFlushDatabaseContext() *FlushDatabaseContext {
	var p = new(FlushDatabaseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_flushDatabase
	return p
}

func InitEmptyFlushDatabaseContext(p *FlushDatabaseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_flushDatabase
}

func (*FlushDatabaseContext) IsFlushDatabaseContext() {}

func NewFlushDatabaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FlushDatabaseContext {
	var p = new(FlushDatabaseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_flushDatabase

	return p
}

func (s *FlushDatabaseContext) GetParser() antlr.Parser { return s.parser }

func (s *FlushDatabaseContext) GetDatabase() IQualifiedNameContext { return s.database }

func (s *FlushDatabaseContext) SetDatabase(v IQualifiedNameContext) { s.database = v }

func (s *FlushDatabaseContext) FLUSH() antlr.TerminalNode {
	return s.GetToken(SQLParserFLUSH, 0)
}

func (s *FlushDatabaseContext) DATABASE() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASE, 0)
}

func (s *FlushDatabaseContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *FlushDatabaseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FlushDatabaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FlushDatabaseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterFlushDatabase(s)
	}
}

func (s *FlushDatabaseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitFlushDatabase(s)
	}
}

func (s *FlushDatabaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitFlushDatabase(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) FlushDatabase() (localctx IFlushDatabaseContext) {
	localctx = NewFlushDatabaseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SQLParserRULE_flushDatabase)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(204)
		p.Match(SQLParserFLUSH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(205)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(206)

		var _x = p.QualifiedName()

		localctx.(*FlushDatabaseContext).database = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompactDatabaseContext is an interface to support dynamic dispatch.
type ICompactDatabaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDatabase returns the database rule contexts.
	GetDatabase() IQualifiedNameContext

	// SetDatabase sets the database rule contexts.
	SetDatabase(IQualifiedNameContext)

	// Getter signatures
	COMPACT() antlr.TerminalNode
	DATABASE() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext

	// IsCompactDatabaseContext differentiates from other interfaces.
	IsCompactDatabaseContext()
}

type CompactDatabaseContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	database IQualifiedNameContext
}

func NewEmptyCompactDatabaseContext() *CompactDatabaseContext {
	var p = new(CompactDatabaseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_compactDatabase
	return p
}

func InitEmptyCompactDatabaseContext(p *CompactDatabaseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_compactDatabase
}

func (*CompactDatabaseContext) IsCompactDatabaseContext() {}

func NewCompactDatabaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompactDatabaseContext {
	var p = new(CompactDatabaseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_compactDatabase

	return p
}

func (s *CompactDatabaseContext) GetParser() antlr.Parser { return s.parser }

func (s *CompactDatabaseContext) GetDatabase() IQualifiedNameContext { return s.database }

func (s *CompactDatabaseContext) SetDatabase(v IQualifiedNameContext) { s.database = v }

func (s *CompactDatabaseContext) COMPACT() antlr.TerminalNode {
	return s.GetToken(SQLParserCOMPACT, 0)
}

func (s *CompactDatabaseContext) DATABASE() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASE, 0)
}

func (s *CompactDatabaseContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *CompactDatabaseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompactDatabaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompactDatabaseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterCompactDatabase(s)
	}
}

func (s *CompactDatabaseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitCompactDatabase(s)
	}
}

func (s *CompactDatabaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitCompactDatabase(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) CompactDatabase() (localctx ICompactDatabaseContext) {
	localctx = NewCompactDatabaseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, SQLParserRULE_compactDatabase)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(208)
		p.Match(SQLParserCOMPACT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(209)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(210)

		var _x = p.QualifiedName()

		localctx.(*CompactDatabaseContext).database = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShowStatementContext is an interface to support dynamic dispatch.
type IShowStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsShowStatementContext differentiates from other interfaces.
	IsShowStatementContext()
}

type ShowStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStatementContext() *ShowStatementContext {
	var p = new(ShowStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_showStatement
	return p
}

func InitEmptyShowStatementContext(p *ShowStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_showStatement
}

func (*ShowStatementContext) IsShowStatementContext() {}

func NewShowStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStatementContext {
	var p = new(ShowStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showStatement

	return p
}

func (s *ShowStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStatementContext) CopyAll(ctx *ShowStatementContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ShowStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ShowStateContext struct {
	ShowStatementContext
}

func NewShowStateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowStateContext {
	var p = new(ShowStateContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowStateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStateContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowStateContext) STATE() antlr.TerminalNode {
	return s.GetToken(SQLParserSTATE, 0)
}

func (s *ShowStateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowState(s)
	}
}

func (s *ShowStateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowState(s)
	}
}

func (s *ShowStateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowState(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowBrokersContext struct {
	ShowStatementContext
}

func NewShowBrokersContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowBrokersContext {
	var p = new(ShowBrokersContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowBrokersContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowBrokersContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowBrokersContext) BROKERS() antlr.TerminalNode {
	return s.GetToken(SQLParserBROKERS, 0)
}

func (s *ShowBrokersContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowBrokers(s)
	}
}

func (s *ShowBrokersContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowBrokers(s)
	}
}

func (s *ShowBrokersContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowBrokers(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowLimitContext struct {
	ShowStatementContext
}

func NewShowLimitContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowLimitContext {
	var p = new(ShowLimitContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowLimitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowLimitContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowLimitContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserLIMIT, 0)
}

func (s *ShowLimitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowLimit(s)
	}
}

func (s *ShowLimitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowLimit(s)
	}
}

func (s *ShowLimitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowLimit(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowMetadataTypesContext struct {
	ShowStatementContext
}

func NewShowMetadataTypesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowMetadataTypesContext {
	var p = new(ShowMetadataTypesContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowMetadataTypesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetadataTypesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowMetadataTypesContext) METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserMETADATA, 0)
}

func (s *ShowMetadataTypesContext) TYPES() antlr.TerminalNode {
	return s.GetToken(SQLParserTYPES, 0)
}

func (s *ShowMetadataTypesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowMetadataTypes(s)
	}
}

func (s *ShowMetadataTypesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowMetadataTypes(s)
	}
}

func (s *ShowMetadataTypesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowMetadataTypes(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowAliveContext struct {
	ShowStatementContext
}

func NewShowAliveContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowAliveContext {
	var p = new(ShowAliveContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowAliveContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowAliveContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowAliveContext) ALIVE() antlr.TerminalNode {
	return s.GetToken(SQLParserALIVE, 0)
}

func (s *ShowAliveContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowAlive(s)
	}
}

func (s *ShowAliveContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowAlive(s)
	}
}

func (s *ShowAliveContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowAlive(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowDatabasesContext struct {
	ShowStatementContext
}

func NewShowDatabasesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowDatabasesContext {
	var p = new(ShowDatabasesContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowDatabasesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowDatabasesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowDatabasesContext) DATABASES() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASES, 0)
}

func (s *ShowDatabasesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowDatabases(s)
	}
}

func (s *ShowDatabasesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowDatabases(s)
	}
}

func (s *ShowDatabasesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowDatabases(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowRequestsContext struct {
	ShowStatementContext
}

func NewShowRequestsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowRequestsContext {
	var p = new(ShowRequestsContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowRequestsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowRequestsContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowRequestsContext) REQUESTS() antlr.TerminalNode {
	return s.GetToken(SQLParserREQUESTS, 0)
}

func (s *ShowRequestsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowRequests(s)
	}
}

func (s *ShowRequestsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowRequests(s)
	}
}

func (s *ShowRequestsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowRequests(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowReplicationsContext struct {
	ShowStatementContext
}

func NewShowReplicationsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowReplicationsContext {
	var p = new(ShowReplicationsContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowReplicationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowReplicationsContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowReplicationsContext) REPLICATIONS() antlr.TerminalNode {
	return s.GetToken(SQLParserREPLICATIONS, 0)
}

func (s *ShowReplicationsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowReplications(s)
	}
}

func (s *ShowReplicationsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowReplications(s)
	}
}

func (s *ShowReplicationsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowReplications(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowMasterContext struct {
	ShowStatementContext
}

func NewShowMasterContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowMasterContext {
	var p = new(ShowMasterContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowMasterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMasterContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowMasterContext) MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserMASTER, 0)
}

func (s *ShowMasterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowMaster(s)
	}
}

func (s *ShowMasterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowMaster(s)
	}
}

func (s *ShowMasterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowMaster(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowMetadatasContext struct {
	ShowStatementContext
}

func NewShowMetadatasContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowMetadatasContext {
	var p = new(ShowMetadatasContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowMetadatasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetadatasContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowMetadatasContext) METADATAS() antlr.TerminalNode {
	return s.GetToken(SQLParserMETADATAS, 0)
}

func (s *ShowMetadatasContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowMetadatas(s)
	}
}

func (s *ShowMetadatasContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowMetadatas(s)
	}
}

func (s *ShowMetadatasContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowMetadatas(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowStatement() (localctx IShowStatementContext) {
	localctx = NewShowStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SQLParserRULE_showStatement)
	p.SetState(233)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 14, p.GetParserRuleContext()) {
	case 1:
		localctx = NewShowMasterContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(212)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(213)
			p.Match(SQLParserMASTER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewShowBrokersContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(214)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(215)
			p.Match(SQLParserBROKERS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewShowRequestsContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(216)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(217)
			p.Match(SQLParserREQUESTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewShowLimitContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(218)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(219)
			p.Match(SQLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewShowMetadataTypesContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(220)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(221)
			p.Match(SQLParserMETADATA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(222)
			p.Match(SQLParserTYPES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		localctx = NewShowMetadatasContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(223)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(224)
			p.Match(SQLParserMETADATAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		localctx = NewShowAliveContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(225)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(226)
			p.Match(SQLParserALIVE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		localctx = NewShowReplicationsContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(227)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(228)
			p.Match(SQLParserREPLICATIONS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 9:
		localctx = NewShowStateContext(p, localctx)
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(229)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(230)
			p.Match(SQLParserSTATE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 10:
		localctx = NewShowDatabasesContext(p, localctx)
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(231)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(232)
			p.Match(SQLParserDATABASES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUseStatementContext is an interface to support dynamic dispatch.
type IUseStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetDatabase returns the database rule contexts.
	GetDatabase() IIdentifierContext

	// SetDatabase sets the database rule contexts.
	SetDatabase(IIdentifierContext)

	// Getter signatures
	USE() antlr.TerminalNode
	Identifier() IIdentifierContext

	// IsUseStatementContext differentiates from other interfaces.
	IsUseStatementContext()
}

type UseStatementContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	database IIdentifierContext
}

func NewEmptyUseStatementContext() *UseStatementContext {
	var p = new(UseStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_useStatement
	return p
}

func InitEmptyUseStatementContext(p *UseStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_useStatement
}

func (*UseStatementContext) IsUseStatementContext() {}

func NewUseStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UseStatementContext {
	var p = new(UseStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_useStatement

	return p
}

func (s *UseStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *UseStatementContext) GetDatabase() IIdentifierContext { return s.database }

func (s *UseStatementContext) SetDatabase(v IIdentifierContext) { s.database = v }

func (s *UseStatementContext) USE() antlr.TerminalNode {
	return s.GetToken(SQLParserUSE, 0)
}

func (s *UseStatementContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *UseStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UseStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UseStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterUseStatement(s)
	}
}

func (s *UseStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitUseStatement(s)
	}
}

func (s *UseStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitUseStatement(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) UseStatement() (localctx IUseStatementContext) {
	localctx = NewUseStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SQLParserRULE_useStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(235)
		p.Match(SQLParserUSE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(236)

		var _x = p.Identifier()

		localctx.(*UseStatementContext).database = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShowMetricMetadataContext is an interface to support dynamic dispatch.
type IShowMetricMetadataContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsShowMetricMetadataContext differentiates from other interfaces.
	IsShowMetricMetadataContext()
}

type ShowMetricMetadataContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowMetricMetadataContext() *ShowMetricMetadataContext {
	var p = new(ShowMetricMetadataContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_showMetricMetadata
	return p
}

func InitEmptyShowMetricMetadataContext(p *ShowMetricMetadataContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_showMetricMetadata
}

func (*ShowMetricMetadataContext) IsShowMetricMetadataContext() {}

func NewShowMetricMetadataContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowMetricMetadataContext {
	var p = new(ShowMetricMetadataContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_showMetricMetadata

	return p
}

func (s *ShowMetricMetadataContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowMetricMetadataContext) CopyAll(ctx *ShowMetricMetadataContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ShowMetricMetadataContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetricMetadataContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ShowTagValuesContext struct {
	ShowMetricMetadataContext
}

func NewShowTagValuesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowTagValuesContext {
	var p = new(ShowTagValuesContext)

	InitEmptyShowMetricMetadataContext(&p.ShowMetricMetadataContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowMetricMetadataContext))

	return p
}

func (s *ShowTagValuesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagValuesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowTagValuesContext) TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserTAG, 0)
}

func (s *ShowTagValuesContext) VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserVALUES, 0)
}

func (s *ShowTagValuesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowTagValues(s)
	}
}

func (s *ShowTagValuesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowTagValues(s)
	}
}

func (s *ShowTagValuesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowTagValues(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowFieldsContext struct {
	ShowMetricMetadataContext
}

func NewShowFieldsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowFieldsContext {
	var p = new(ShowFieldsContext)

	InitEmptyShowMetricMetadataContext(&p.ShowMetricMetadataContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowMetricMetadataContext))

	return p
}

func (s *ShowFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowFieldsContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowFieldsContext) FIELDS() antlr.TerminalNode {
	return s.GetToken(SQLParserFIELDS, 0)
}

func (s *ShowFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowFields(s)
	}
}

func (s *ShowFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowFields(s)
	}
}

func (s *ShowFieldsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowFields(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowTagKeysContext struct {
	ShowMetricMetadataContext
}

func NewShowTagKeysContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowTagKeysContext {
	var p = new(ShowTagKeysContext)

	InitEmptyShowMetricMetadataContext(&p.ShowMetricMetadataContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowMetricMetadataContext))

	return p
}

func (s *ShowTagKeysContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTagKeysContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowTagKeysContext) TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserTAG, 0)
}

func (s *ShowTagKeysContext) KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserKEYS, 0)
}

func (s *ShowTagKeysContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowTagKeys(s)
	}
}

func (s *ShowTagKeysContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowTagKeys(s)
	}
}

func (s *ShowTagKeysContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowTagKeys(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowMetricsContext struct {
	ShowMetricMetadataContext
}

func NewShowMetricsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowMetricsContext {
	var p = new(ShowMetricsContext)

	InitEmptyShowMetricMetadataContext(&p.ShowMetricMetadataContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowMetricMetadataContext))

	return p
}

func (s *ShowMetricsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMetricsContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowMetricsContext) METRICS() antlr.TerminalNode {
	return s.GetToken(SQLParserMETRICS, 0)
}

func (s *ShowMetricsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowMetrics(s)
	}
}

func (s *ShowMetricsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowMetrics(s)
	}
}

func (s *ShowMetricsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowMetrics(s)

	default:
		return t.VisitChildren(s)
	}
}

type ShowNamespacesContext struct {
	ShowMetricMetadataContext
}

func NewShowNamespacesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowNamespacesContext {
	var p = new(ShowNamespacesContext)

	InitEmptyShowMetricMetadataContext(&p.ShowMetricMetadataContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowMetricMetadataContext))

	return p
}

func (s *ShowNamespacesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowNamespacesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowNamespacesContext) NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserNAMESPACES, 0)
}

func (s *ShowNamespacesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowNamespaces(s)
	}
}

func (s *ShowNamespacesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowNamespaces(s)
	}
}

func (s *ShowNamespacesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowNamespaces(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowMetricMetadata() (localctx IShowMetricMetadataContext) {
	localctx = NewShowMetricMetadataContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, SQLParserRULE_showMetricMetadata)
	p.SetState(250)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		localctx = NewShowNamespacesContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(238)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(239)
			p.Match(SQLParserNAMESPACES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewShowMetricsContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(240)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(241)
			p.Match(SQLParserMETRICS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewShowFieldsContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(242)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(243)
			p.Match(SQLParserFIELDS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewShowTagKeysContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(244)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(245)
			p.Match(SQLParserTAG)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(246)
			p.Match(SQLParserKEYS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewShowTagValuesContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(247)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(248)
			p.Match(SQLParserTAG)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(249)
			p.Match(SQLParserVALUES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QueryNoWith() IQueryNoWithContext
	With() IWithContext

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) QueryNoWith() IQueryNoWithContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryNoWithContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryNoWithContext)
}

func (s *QueryContext) With() IWithContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWithContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWithContext)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (s *QueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_query)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWITH {
		{
			p.SetState(252)
			p.With()
		}

	}
	{
		p.SetState(255)
		p.QueryNoWith()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWithContext is an interface to support dynamic dispatch.
type IWithContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WITH() antlr.TerminalNode
	AllNamedQuery() []INamedQueryContext
	NamedQuery(i int) INamedQueryContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsWithContext differentiates from other interfaces.
	IsWithContext()
}

type WithContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWithContext() *WithContext {
	var p = new(WithContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_with
	return p
}

func InitEmptyWithContext(p *WithContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_with
}

func (*WithContext) IsWithContext() {}

func NewWithContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WithContext {
	var p = new(WithContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_with

	return p
}

func (s *WithContext) GetParser() antlr.Parser { return s.parser }

func (s *WithContext) WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserWITH, 0)
}

func (s *WithContext) AllNamedQuery() []INamedQueryContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INamedQueryContext); ok {
			len++
		}
	}

	tst := make([]INamedQueryContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INamedQueryContext); ok {
			tst[i] = t.(INamedQueryContext)
			i++
		}
	}

	return tst
}

func (s *WithContext) NamedQuery(i int) INamedQueryContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamedQueryContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamedQueryContext)
}

func (s *WithContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *WithContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *WithContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WithContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterWith(s)
	}
}

func (s *WithContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitWith(s)
	}
}

func (s *WithContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitWith(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) With() (localctx IWithContext) {
	localctx = NewWithContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, SQLParserRULE_with)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(257)
		p.Match(SQLParserWITH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(258)
		p.NamedQuery()
	}
	p.SetState(263)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(259)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(260)
			p.NamedQuery()
		}

		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamedQueryContext is an interface to support dynamic dispatch.
type INamedQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name rule contexts.
	GetName() IIdentifierContext

	// SetName sets the name rule contexts.
	SetName(IIdentifierContext)

	// Getter signatures
	AS() antlr.TerminalNode
	LR_BRACKET() antlr.TerminalNode
	Query() IQueryContext
	RR_BRACKET() antlr.TerminalNode
	Identifier() IIdentifierContext

	// IsNamedQueryContext differentiates from other interfaces.
	IsNamedQueryContext()
}

type NamedQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	name   IIdentifierContext
}

func NewEmptyNamedQueryContext() *NamedQueryContext {
	var p = new(NamedQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_namedQuery
	return p
}

func InitEmptyNamedQueryContext(p *NamedQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_namedQuery
}

func (*NamedQueryContext) IsNamedQueryContext() {}

func NewNamedQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamedQueryContext {
	var p = new(NamedQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_namedQuery

	return p
}

func (s *NamedQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *NamedQueryContext) GetName() IIdentifierContext { return s.name }

func (s *NamedQueryContext) SetName(v IIdentifierContext) { s.name = v }

func (s *NamedQueryContext) AS() antlr.TerminalNode {
	return s.GetToken(SQLParserAS, 0)
}

func (s *NamedQueryContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *NamedQueryContext) Query() IQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryContext)
}

func (s *NamedQueryContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *NamedQueryContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *NamedQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamedQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamedQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterNamedQuery(s)
	}
}

func (s *NamedQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitNamedQuery(s)
	}
}

func (s *NamedQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitNamedQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) NamedQuery() (localctx INamedQueryContext) {
	localctx = NewNamedQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, SQLParserRULE_namedQuery)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(266)

		var _x = p.Identifier()

		localctx.(*NamedQueryContext).name = _x
	}
	{
		p.SetState(267)
		p.Match(SQLParserAS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(268)
		p.Match(SQLParserLR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(269)
		p.Query()
	}
	{
		p.SetState(270)
		p.Match(SQLParserRR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQueryNoWithContext is an interface to support dynamic dispatch.
type IQueryNoWithContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLimit returns the limit rule contexts.
	GetLimit() ILimitRowCountContext

	// SetLimit sets the limit rule contexts.
	SetLimit(ILimitRowCountContext)

	// Getter signatures
	QueryTerm() IQueryTermContext
	ORDER() antlr.TerminalNode
	BY() antlr.TerminalNode
	OrderBy() IOrderByContext
	LIMIT() antlr.TerminalNode
	LimitRowCount() ILimitRowCountContext

	// IsQueryNoWithContext differentiates from other interfaces.
	IsQueryNoWithContext()
}

type QueryNoWithContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	limit  ILimitRowCountContext
}

func NewEmptyQueryNoWithContext() *QueryNoWithContext {
	var p = new(QueryNoWithContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_queryNoWith
	return p
}

func InitEmptyQueryNoWithContext(p *QueryNoWithContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_queryNoWith
}

func (*QueryNoWithContext) IsQueryNoWithContext() {}

func NewQueryNoWithContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryNoWithContext {
	var p = new(QueryNoWithContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_queryNoWith

	return p
}

func (s *QueryNoWithContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryNoWithContext) GetLimit() ILimitRowCountContext { return s.limit }

func (s *QueryNoWithContext) SetLimit(v ILimitRowCountContext) { s.limit = v }

func (s *QueryNoWithContext) QueryTerm() IQueryTermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryTermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryTermContext)
}

func (s *QueryNoWithContext) ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserORDER, 0)
}

func (s *QueryNoWithContext) BY() antlr.TerminalNode {
	return s.GetToken(SQLParserBY, 0)
}

func (s *QueryNoWithContext) OrderBy() IOrderByContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrderByContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrderByContext)
}

func (s *QueryNoWithContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserLIMIT, 0)
}

func (s *QueryNoWithContext) LimitRowCount() ILimitRowCountContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitRowCountContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitRowCountContext)
}

func (s *QueryNoWithContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryNoWithContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryNoWithContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQueryNoWith(s)
	}
}

func (s *QueryNoWithContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQueryNoWith(s)
	}
}

func (s *QueryNoWithContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQueryNoWith(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QueryNoWith() (localctx IQueryNoWithContext) {
	localctx = NewQueryNoWithContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, SQLParserRULE_queryNoWith)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(272)
		p.QueryTerm()
	}
	p.SetState(276)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserORDER {
		{
			p.SetState(273)
			p.Match(SQLParserORDER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(274)
			p.Match(SQLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(275)
			p.OrderBy()
		}

	}
	p.SetState(280)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserLIMIT {
		{
			p.SetState(278)
			p.Match(SQLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(279)

			var _x = p.LimitRowCount()

			localctx.(*QueryNoWithContext).limit = _x
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQueryTermContext is an interface to support dynamic dispatch.
type IQueryTermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsQueryTermContext differentiates from other interfaces.
	IsQueryTermContext()
}

type QueryTermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryTermContext() *QueryTermContext {
	var p = new(QueryTermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_queryTerm
	return p
}

func InitEmptyQueryTermContext(p *QueryTermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_queryTerm
}

func (*QueryTermContext) IsQueryTermContext() {}

func NewQueryTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryTermContext {
	var p = new(QueryTermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_queryTerm

	return p
}

func (s *QueryTermContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryTermContext) CopyAll(ctx *QueryTermContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *QueryTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryTermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type QueryTermDefaultContext struct {
	QueryTermContext
}

func NewQueryTermDefaultContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *QueryTermDefaultContext {
	var p = new(QueryTermDefaultContext)

	InitEmptyQueryTermContext(&p.QueryTermContext)
	p.parser = parser
	p.CopyAll(ctx.(*QueryTermContext))

	return p
}

func (s *QueryTermDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryTermDefaultContext) QueryPrimary() IQueryPrimaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryPrimaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryPrimaryContext)
}

func (s *QueryTermDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQueryTermDefault(s)
	}
}

func (s *QueryTermDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQueryTermDefault(s)
	}
}

func (s *QueryTermDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQueryTermDefault(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QueryTerm() (localctx IQueryTermContext) {
	localctx = NewQueryTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, SQLParserRULE_queryTerm)
	localctx = NewQueryTermDefaultContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(282)
		p.QueryPrimary()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQueryPrimaryContext is an interface to support dynamic dispatch.
type IQueryPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsQueryPrimaryContext differentiates from other interfaces.
	IsQueryPrimaryContext()
}

type QueryPrimaryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryPrimaryContext() *QueryPrimaryContext {
	var p = new(QueryPrimaryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_queryPrimary
	return p
}

func InitEmptyQueryPrimaryContext(p *QueryPrimaryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_queryPrimary
}

func (*QueryPrimaryContext) IsQueryPrimaryContext() {}

func NewQueryPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryPrimaryContext {
	var p = new(QueryPrimaryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_queryPrimary

	return p
}

func (s *QueryPrimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryPrimaryContext) CopyAll(ctx *QueryPrimaryContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *QueryPrimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryPrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SubqueryContext struct {
	QueryPrimaryContext
}

func NewSubqueryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubqueryContext {
	var p = new(SubqueryContext)

	InitEmptyQueryPrimaryContext(&p.QueryPrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*QueryPrimaryContext))

	return p
}

func (s *SubqueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubqueryContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *SubqueryContext) QueryNoWith() IQueryNoWithContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryNoWithContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryNoWithContext)
}

func (s *SubqueryContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *SubqueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSubquery(s)
	}
}

func (s *SubqueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSubquery(s)
	}
}

func (s *SubqueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSubquery(s)

	default:
		return t.VisitChildren(s)
	}
}

type QueryPrimaryDefaultContext struct {
	QueryPrimaryContext
}

func NewQueryPrimaryDefaultContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *QueryPrimaryDefaultContext {
	var p = new(QueryPrimaryDefaultContext)

	InitEmptyQueryPrimaryContext(&p.QueryPrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*QueryPrimaryContext))

	return p
}

func (s *QueryPrimaryDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryPrimaryDefaultContext) QuerySpecification() IQuerySpecificationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQuerySpecificationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQuerySpecificationContext)
}

func (s *QueryPrimaryDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQueryPrimaryDefault(s)
	}
}

func (s *QueryPrimaryDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQueryPrimaryDefault(s)
	}
}

func (s *QueryPrimaryDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQueryPrimaryDefault(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QueryPrimary() (localctx IQueryPrimaryContext) {
	localctx = NewQueryPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, SQLParserRULE_queryPrimary)
	p.SetState(289)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserSELECT:
		localctx = NewQueryPrimaryDefaultContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(284)
			p.QuerySpecification()
		}

	case SQLParserLR_BRACKET:
		localctx = NewSubqueryContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(285)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(286)
			p.QueryNoWith()
		}
		{
			p.SetState(287)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQuerySpecificationContext is an interface to support dynamic dispatch.
type IQuerySpecificationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetWhere returns the where rule contexts.
	GetWhere() IBooleanExpressionContext

	// SetWhere sets the where rule contexts.
	SetWhere(IBooleanExpressionContext)

	// Getter signatures
	SELECT() antlr.TerminalNode
	AllSelectItem() []ISelectItemContext
	SelectItem(i int) ISelectItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	FROM() antlr.TerminalNode
	AllRelation() []IRelationContext
	Relation(i int) IRelationContext
	WHERE() antlr.TerminalNode
	GROUP() antlr.TerminalNode
	BY() antlr.TerminalNode
	GroupBy() IGroupByContext
	HAVING() antlr.TerminalNode
	Having() IHavingContext
	BooleanExpression() IBooleanExpressionContext

	// IsQuerySpecificationContext differentiates from other interfaces.
	IsQuerySpecificationContext()
}

type QuerySpecificationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	where  IBooleanExpressionContext
}

func NewEmptyQuerySpecificationContext() *QuerySpecificationContext {
	var p = new(QuerySpecificationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_querySpecification
	return p
}

func InitEmptyQuerySpecificationContext(p *QuerySpecificationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_querySpecification
}

func (*QuerySpecificationContext) IsQuerySpecificationContext() {}

func NewQuerySpecificationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QuerySpecificationContext {
	var p = new(QuerySpecificationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_querySpecification

	return p
}

func (s *QuerySpecificationContext) GetParser() antlr.Parser { return s.parser }

func (s *QuerySpecificationContext) GetWhere() IBooleanExpressionContext { return s.where }

func (s *QuerySpecificationContext) SetWhere(v IBooleanExpressionContext) { s.where = v }

func (s *QuerySpecificationContext) SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserSELECT, 0)
}

func (s *QuerySpecificationContext) AllSelectItem() []ISelectItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISelectItemContext); ok {
			len++
		}
	}

	tst := make([]ISelectItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISelectItemContext); ok {
			tst[i] = t.(ISelectItemContext)
			i++
		}
	}

	return tst
}

func (s *QuerySpecificationContext) SelectItem(i int) ISelectItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectItemContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectItemContext)
}

func (s *QuerySpecificationContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *QuerySpecificationContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *QuerySpecificationContext) FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserFROM, 0)
}

func (s *QuerySpecificationContext) AllRelation() []IRelationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelationContext); ok {
			len++
		}
	}

	tst := make([]IRelationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelationContext); ok {
			tst[i] = t.(IRelationContext)
			i++
		}
	}

	return tst
}

func (s *QuerySpecificationContext) Relation(i int) IRelationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *QuerySpecificationContext) WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserWHERE, 0)
}

func (s *QuerySpecificationContext) GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserGROUP, 0)
}

func (s *QuerySpecificationContext) BY() antlr.TerminalNode {
	return s.GetToken(SQLParserBY, 0)
}

func (s *QuerySpecificationContext) GroupBy() IGroupByContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByContext)
}

func (s *QuerySpecificationContext) HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserHAVING, 0)
}

func (s *QuerySpecificationContext) Having() IHavingContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHavingContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHavingContext)
}

func (s *QuerySpecificationContext) BooleanExpression() IBooleanExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *QuerySpecificationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuerySpecificationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QuerySpecificationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQuerySpecification(s)
	}
}

func (s *QuerySpecificationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQuerySpecification(s)
	}
}

func (s *QuerySpecificationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQuerySpecification(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QuerySpecification() (localctx IQuerySpecificationContext) {
	localctx = NewQuerySpecificationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, SQLParserRULE_querySpecification)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(291)
		p.Match(SQLParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(292)
		p.SelectItem()
	}
	p.SetState(297)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(293)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(294)
			p.SelectItem()
		}

		p.SetState(299)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(309)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserFROM {
		{
			p.SetState(300)
			p.Match(SQLParserFROM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(301)
			p.relation(0)
		}
		p.SetState(306)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(302)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(303)
				p.relation(0)
			}

			p.SetState(308)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	p.SetState(313)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWHERE {
		{
			p.SetState(311)
			p.Match(SQLParserWHERE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(312)

			var _x = p.booleanExpression(0)

			localctx.(*QuerySpecificationContext).where = _x
		}

	}
	p.SetState(318)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserGROUP {
		{
			p.SetState(315)
			p.Match(SQLParserGROUP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(316)
			p.Match(SQLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(317)
			p.GroupBy()
		}

	}
	p.SetState(322)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserHAVING {
		{
			p.SetState(320)
			p.Match(SQLParserHAVING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(321)
			p.Having()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectItemContext is an interface to support dynamic dispatch.
type ISelectItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSelectItemContext differentiates from other interfaces.
	IsSelectItemContext()
}

type SelectItemContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectItemContext() *SelectItemContext {
	var p = new(SelectItemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_selectItem
	return p
}

func InitEmptySelectItemContext(p *SelectItemContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_selectItem
}

func (*SelectItemContext) IsSelectItemContext() {}

func NewSelectItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectItemContext {
	var p = new(SelectItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_selectItem

	return p
}

func (s *SelectItemContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectItemContext) CopyAll(ctx *SelectItemContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *SelectItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectAllContext struct {
	SelectItemContext
}

func NewSelectAllContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectAllContext {
	var p = new(SelectAllContext)

	InitEmptySelectItemContext(&p.SelectItemContext)
	p.parser = parser
	p.CopyAll(ctx.(*SelectItemContext))

	return p
}

func (s *SelectAllContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectAllContext) PrimaryExpression() IPrimaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionContext)
}

func (s *SelectAllContext) DOT() antlr.TerminalNode {
	return s.GetToken(SQLParserDOT, 0)
}

func (s *SelectAllContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(SQLParserASTERISK, 0)
}

func (s *SelectAllContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSelectAll(s)
	}
}

func (s *SelectAllContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSelectAll(s)
	}
}

func (s *SelectAllContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSelectAll(s)

	default:
		return t.VisitChildren(s)
	}
}

type SelectSingleContext struct {
	SelectItemContext
}

func NewSelectSingleContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectSingleContext {
	var p = new(SelectSingleContext)

	InitEmptySelectItemContext(&p.SelectItemContext)
	p.parser = parser
	p.CopyAll(ctx.(*SelectItemContext))

	return p
}

func (s *SelectSingleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectSingleContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SelectSingleContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *SelectSingleContext) AS() antlr.TerminalNode {
	return s.GetToken(SQLParserAS, 0)
}

func (s *SelectSingleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSelectSingle(s)
	}
}

func (s *SelectSingleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSelectSingle(s)
	}
}

func (s *SelectSingleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSelectSingle(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SelectItem() (localctx ISelectItemContext) {
	localctx = NewSelectItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, SQLParserRULE_selectItem)
	p.SetState(336)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSelectSingleContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(324)
			p.Expression()
		}
		p.SetState(329)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) == 1 {
			p.SetState(326)
			p.GetErrorHandler().Sync(p)

			if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 27, p.GetParserRuleContext()) == 1 {
				{
					p.SetState(325)
					p.Match(SQLParserAS)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			} else if p.HasError() { // JIM
				goto errorExit
			}
			{
				p.SetState(328)
				p.Identifier()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 2:
		localctx = NewSelectAllContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(331)
			p.primaryExpression(0)
		}
		{
			p.SetState(332)
			p.Match(SQLParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(333)
			p.Match(SQLParserASTERISK)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewSelectAllContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(335)
			p.Match(SQLParserASTERISK)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsRelationContext differentiates from other interfaces.
	IsRelationContext()
}

type RelationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationContext() *RelationContext {
	var p = new(RelationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) CopyAll(ctx *RelationContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type RelationDefaultContext struct {
	RelationContext
}

func NewRelationDefaultContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RelationDefaultContext {
	var p = new(RelationDefaultContext)

	InitEmptyRelationContext(&p.RelationContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelationContext))

	return p
}

func (s *RelationDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationDefaultContext) AliasedRelation() IAliasedRelationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasedRelationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasedRelationContext)
}

func (s *RelationDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterRelationDefault(s)
	}
}

func (s *RelationDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitRelationDefault(s)
	}
}

func (s *RelationDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitRelationDefault(s)

	default:
		return t.VisitChildren(s)
	}
}

type JoinRelationContext struct {
	RelationContext
	left          IRelationContext
	right         IRelationContext
	rightRelation IRelationContext
}

func NewJoinRelationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *JoinRelationContext {
	var p = new(JoinRelationContext)

	InitEmptyRelationContext(&p.RelationContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelationContext))

	return p
}

func (s *JoinRelationContext) GetLeft() IRelationContext { return s.left }

func (s *JoinRelationContext) GetRight() IRelationContext { return s.right }

func (s *JoinRelationContext) GetRightRelation() IRelationContext { return s.rightRelation }

func (s *JoinRelationContext) SetLeft(v IRelationContext) { s.left = v }

func (s *JoinRelationContext) SetRight(v IRelationContext) { s.right = v }

func (s *JoinRelationContext) SetRightRelation(v IRelationContext) { s.rightRelation = v }

func (s *JoinRelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinRelationContext) AllRelation() []IRelationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelationContext); ok {
			len++
		}
	}

	tst := make([]IRelationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelationContext); ok {
			tst[i] = t.(IRelationContext)
			i++
		}
	}

	return tst
}

func (s *JoinRelationContext) Relation(i int) IRelationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *JoinRelationContext) CROSS() antlr.TerminalNode {
	return s.GetToken(SQLParserCROSS, 0)
}

func (s *JoinRelationContext) JOIN() antlr.TerminalNode {
	return s.GetToken(SQLParserJOIN, 0)
}

func (s *JoinRelationContext) JoinType() IJoinTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinTypeContext)
}

func (s *JoinRelationContext) JoinCriteria() IJoinCriteriaContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinCriteriaContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinCriteriaContext)
}

func (s *JoinRelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterJoinRelation(s)
	}
}

func (s *JoinRelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitJoinRelation(s)
	}
}

func (s *JoinRelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitJoinRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Relation() (localctx IRelationContext) {
	return p.relation(0)
}

func (p *SQLParser) relation(_p int) (localctx IRelationContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewRelationContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IRelationContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 48
	p.EnterRecursionRule(localctx, 48, SQLParserRULE_relation, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewRelationDefaultContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(339)
		p.AliasedRelation()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(354)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewJoinRelationContext(p, NewRelationContext(p, _parentctx, _parentState))
			localctx.(*JoinRelationContext).left = _prevctx

			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_relation)
			p.SetState(341)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				goto errorExit
			}
			p.SetState(350)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case SQLParserCROSS:
				{
					p.SetState(342)
					p.Match(SQLParserCROSS)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(343)
					p.Match(SQLParserJOIN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(344)

					var _x = p.relation(0)

					localctx.(*JoinRelationContext).right = _x
				}

			case SQLParserLEFT, SQLParserRIGHT:
				{
					p.SetState(345)
					p.JoinType()
				}
				{
					p.SetState(346)
					p.Match(SQLParserJOIN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(347)

					var _x = p.relation(0)

					localctx.(*JoinRelationContext).rightRelation = _x
				}
				{
					p.SetState(348)
					p.JoinCriteria()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

		}
		p.SetState(356)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 31, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IJoinTypeContext is an interface to support dynamic dispatch.
type IJoinTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LEFT() antlr.TerminalNode
	RIGHT() antlr.TerminalNode

	// IsJoinTypeContext differentiates from other interfaces.
	IsJoinTypeContext()
}

type JoinTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoinTypeContext() *JoinTypeContext {
	var p = new(JoinTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_joinType
	return p
}

func InitEmptyJoinTypeContext(p *JoinTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_joinType
}

func (*JoinTypeContext) IsJoinTypeContext() {}

func NewJoinTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinTypeContext {
	var p = new(JoinTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_joinType

	return p
}

func (s *JoinTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinTypeContext) LEFT() antlr.TerminalNode {
	return s.GetToken(SQLParserLEFT, 0)
}

func (s *JoinTypeContext) RIGHT() antlr.TerminalNode {
	return s.GetToken(SQLParserRIGHT, 0)
}

func (s *JoinTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *JoinTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterJoinType(s)
	}
}

func (s *JoinTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitJoinType(s)
	}
}

func (s *JoinTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitJoinType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) JoinType() (localctx IJoinTypeContext) {
	localctx = NewJoinTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, SQLParserRULE_joinType)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(357)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserLEFT || _la == SQLParserRIGHT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IJoinCriteriaContext is an interface to support dynamic dispatch.
type IJoinCriteriaContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ON() antlr.TerminalNode
	BooleanExpression() IBooleanExpressionContext
	USING() antlr.TerminalNode
	LR_BRACKET() antlr.TerminalNode
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext
	RR_BRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsJoinCriteriaContext differentiates from other interfaces.
	IsJoinCriteriaContext()
}

type JoinCriteriaContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoinCriteriaContext() *JoinCriteriaContext {
	var p = new(JoinCriteriaContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_joinCriteria
	return p
}

func InitEmptyJoinCriteriaContext(p *JoinCriteriaContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_joinCriteria
}

func (*JoinCriteriaContext) IsJoinCriteriaContext() {}

func NewJoinCriteriaContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinCriteriaContext {
	var p = new(JoinCriteriaContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_joinCriteria

	return p
}

func (s *JoinCriteriaContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinCriteriaContext) ON() antlr.TerminalNode {
	return s.GetToken(SQLParserON, 0)
}

func (s *JoinCriteriaContext) BooleanExpression() IBooleanExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *JoinCriteriaContext) USING() antlr.TerminalNode {
	return s.GetToken(SQLParserUSING, 0)
}

func (s *JoinCriteriaContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *JoinCriteriaContext) AllIdentifier() []IIdentifierContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentifierContext); ok {
			len++
		}
	}

	tst := make([]IIdentifierContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentifierContext); ok {
			tst[i] = t.(IIdentifierContext)
			i++
		}
	}

	return tst
}

func (s *JoinCriteriaContext) Identifier(i int) IIdentifierContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *JoinCriteriaContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *JoinCriteriaContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *JoinCriteriaContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *JoinCriteriaContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinCriteriaContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *JoinCriteriaContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterJoinCriteria(s)
	}
}

func (s *JoinCriteriaContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitJoinCriteria(s)
	}
}

func (s *JoinCriteriaContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitJoinCriteria(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) JoinCriteria() (localctx IJoinCriteriaContext) {
	localctx = NewJoinCriteriaContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, SQLParserRULE_joinCriteria)
	var _la int

	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserON:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(359)
			p.Match(SQLParserON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(360)
			p.booleanExpression(0)
		}

	case SQLParserUSING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(361)
			p.Match(SQLParserUSING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(362)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(363)
			p.Identifier()
		}
		p.SetState(368)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(364)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(365)
				p.Identifier()
			}

			p.SetState(370)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(371)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAliasedRelationContext is an interface to support dynamic dispatch.
type IAliasedRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RelationPrimary() IRelationPrimaryContext
	Identifier() IIdentifierContext
	AS() antlr.TerminalNode

	// IsAliasedRelationContext differentiates from other interfaces.
	IsAliasedRelationContext()
}

type AliasedRelationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAliasedRelationContext() *AliasedRelationContext {
	var p = new(AliasedRelationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_aliasedRelation
	return p
}

func InitEmptyAliasedRelationContext(p *AliasedRelationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_aliasedRelation
}

func (*AliasedRelationContext) IsAliasedRelationContext() {}

func NewAliasedRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AliasedRelationContext {
	var p = new(AliasedRelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_aliasedRelation

	return p
}

func (s *AliasedRelationContext) GetParser() antlr.Parser { return s.parser }

func (s *AliasedRelationContext) RelationPrimary() IRelationPrimaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationPrimaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationPrimaryContext)
}

func (s *AliasedRelationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *AliasedRelationContext) AS() antlr.TerminalNode {
	return s.GetToken(SQLParserAS, 0)
}

func (s *AliasedRelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AliasedRelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AliasedRelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterAliasedRelation(s)
	}
}

func (s *AliasedRelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitAliasedRelation(s)
	}
}

func (s *AliasedRelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitAliasedRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) AliasedRelation() (localctx IAliasedRelationContext) {
	localctx = NewAliasedRelationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, SQLParserRULE_aliasedRelation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(375)
		p.RelationPrimary()
	}
	p.SetState(380)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext()) == 1 {
		p.SetState(377)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(376)
				p.Match(SQLParserAS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}
		{
			p.SetState(379)
			p.Identifier()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRelationPrimaryContext is an interface to support dynamic dispatch.
type IRelationPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsRelationPrimaryContext differentiates from other interfaces.
	IsRelationPrimaryContext()
}

type RelationPrimaryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRelationPrimaryContext() *RelationPrimaryContext {
	var p = new(RelationPrimaryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_relationPrimary
	return p
}

func InitEmptyRelationPrimaryContext(p *RelationPrimaryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_relationPrimary
}

func (*RelationPrimaryContext) IsRelationPrimaryContext() {}

func NewRelationPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationPrimaryContext {
	var p = new(RelationPrimaryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_relationPrimary

	return p
}

func (s *RelationPrimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationPrimaryContext) CopyAll(ctx *RelationPrimaryContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *RelationPrimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationPrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SubQueryRelationContext struct {
	RelationPrimaryContext
}

func NewSubQueryRelationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubQueryRelationContext {
	var p = new(SubQueryRelationContext)

	InitEmptyRelationPrimaryContext(&p.RelationPrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelationPrimaryContext))

	return p
}

func (s *SubQueryRelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubQueryRelationContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *SubQueryRelationContext) Query() IQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQueryContext)
}

func (s *SubQueryRelationContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *SubQueryRelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSubQueryRelation(s)
	}
}

func (s *SubQueryRelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSubQueryRelation(s)
	}
}

func (s *SubQueryRelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSubQueryRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

type TableNameContext struct {
	RelationPrimaryContext
}

func NewTableNameContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TableNameContext {
	var p = new(TableNameContext)

	InitEmptyRelationPrimaryContext(&p.RelationPrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*RelationPrimaryContext))

	return p
}

func (s *TableNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableNameContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *TableNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterTableName(s)
	}
}

func (s *TableNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitTableName(s)
	}
}

func (s *TableNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitTableName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) RelationPrimary() (localctx IRelationPrimaryContext) {
	localctx = NewRelationPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, SQLParserRULE_relationPrimary)
	p.SetState(387)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserALL, SQLParserALIVE, SQLParserAND, SQLParserAS, SQLParserASC, SQLParserBROKER, SQLParserBROKERS, SQLParserBY, SQLParserCOMPACT, SQLParserCREATE, SQLParserCROSS, SQLParserDATABASE, SQLParserDATABASES, SQLParserDEFAULT, SQLParserDESC, SQLParserDISTRIBUTED, SQLParserDROP, SQLParserENGINE, SQLParserESCAPE, SQLParserEXPLAIN, SQLParserEXISTS, SQLParserFALSE, SQLParserFIELDS, SQLParserFILTER, SQLParserFLUSH, SQLParserFROM, SQLParserGROUP, SQLParserHAVING, SQLParserIF, SQLParserIN, SQLParserJOIN, SQLParserKEYS, SQLParserLEFT, SQLParserLIKE, SQLParserLIMIT, SQLParserLOG, SQLParserLOGICAL, SQLParserMASTER, SQLParserMETRICS, SQLParserMETRIC, SQLParserMETADATA, SQLParserMETADATAS, SQLParserNAMESPACE, SQLParserNAMESPACES, SQLParserNOT, SQLParserON, SQLParserOR, SQLParserORDER, SQLParserPLAN, SQLParserREQUESTS, SQLParserREPLICATIONS, SQLParserRIGHT, SQLParserROLLUP, SQLParserSELECT, SQLParserSHOW, SQLParserSTATE, SQLParserSTORAGE, SQLParserTAG, SQLParserTRACE, SQLParserTRUE, SQLParserTYPE, SQLParserTYPES, SQLParserVALUES, SQLParserWHERE, SQLParserWITH, SQLParserWITHIN, SQLParserUSING, SQLParserUSE, SQLParserIDENTIFIER, SQLParserDIGIT_IDENTIFIER, SQLParserQUOTED_IDENTIFIER, SQLParserBACKQUOTED_IDENTIFIER:
		localctx = NewTableNameContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(382)
			p.QualifiedName()
		}

	case SQLParserLR_BRACKET:
		localctx = NewSubQueryRelationContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(383)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(384)
			p.Query()
		}
		{
			p.SetState(385)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupByContext is an interface to support dynamic dispatch.
type IGroupByContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllGroupingElement() []IGroupingElementContext
	GroupingElement(i int) IGroupingElementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGroupByContext differentiates from other interfaces.
	IsGroupByContext()
}

type GroupByContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByContext() *GroupByContext {
	var p = new(GroupByContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_groupBy
	return p
}

func InitEmptyGroupByContext(p *GroupByContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_groupBy
}

func (*GroupByContext) IsGroupByContext() {}

func NewGroupByContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByContext {
	var p = new(GroupByContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupBy

	return p
}

func (s *GroupByContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByContext) AllGroupingElement() []IGroupingElementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IGroupingElementContext); ok {
			len++
		}
	}

	tst := make([]IGroupingElementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IGroupingElementContext); ok {
			tst[i] = t.(IGroupingElementContext)
			i++
		}
	}

	return tst
}

func (s *GroupByContext) GroupingElement(i int) IGroupingElementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupingElementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupingElementContext)
}

func (s *GroupByContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *GroupByContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *GroupByContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupByContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterGroupBy(s)
	}
}

func (s *GroupByContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitGroupBy(s)
	}
}

func (s *GroupByContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitGroupBy(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) GroupBy() (localctx IGroupByContext) {
	localctx = NewGroupByContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, SQLParserRULE_groupBy)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(389)
		p.GroupingElement()
	}
	p.SetState(394)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(390)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(391)
			p.GroupingElement()
		}

		p.SetState(396)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupingElementContext is an interface to support dynamic dispatch.
type IGroupingElementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsGroupingElementContext differentiates from other interfaces.
	IsGroupingElementContext()
}

type GroupingElementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupingElementContext() *GroupingElementContext {
	var p = new(GroupingElementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_groupingElement
	return p
}

func InitEmptyGroupingElementContext(p *GroupingElementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_groupingElement
}

func (*GroupingElementContext) IsGroupingElementContext() {}

func NewGroupingElementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupingElementContext {
	var p = new(GroupingElementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupingElement

	return p
}

func (s *GroupingElementContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupingElementContext) CopyAll(ctx *GroupingElementContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *GroupingElementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupingElementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SingleGroupingSetContext struct {
	GroupingElementContext
}

func NewSingleGroupingSetContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SingleGroupingSetContext {
	var p = new(SingleGroupingSetContext)

	InitEmptyGroupingElementContext(&p.GroupingElementContext)
	p.parser = parser
	p.CopyAll(ctx.(*GroupingElementContext))

	return p
}

func (s *SingleGroupingSetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleGroupingSetContext) GroupingSet() IGroupingSetContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupingSetContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupingSetContext)
}

func (s *SingleGroupingSetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSingleGroupingSet(s)
	}
}

func (s *SingleGroupingSetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSingleGroupingSet(s)
	}
}

func (s *SingleGroupingSetContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSingleGroupingSet(s)

	default:
		return t.VisitChildren(s)
	}
}

type GroupByAllColumnsContext struct {
	GroupingElementContext
}

func NewGroupByAllColumnsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GroupByAllColumnsContext {
	var p = new(GroupByAllColumnsContext)

	InitEmptyGroupingElementContext(&p.GroupingElementContext)
	p.parser = parser
	p.CopyAll(ctx.(*GroupingElementContext))

	return p
}

func (s *GroupByAllColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByAllColumnsContext) ALL() antlr.TerminalNode {
	return s.GetToken(SQLParserALL, 0)
}

func (s *GroupByAllColumnsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterGroupByAllColumns(s)
	}
}

func (s *GroupByAllColumnsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitGroupByAllColumns(s)
	}
}

func (s *GroupByAllColumnsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitGroupByAllColumns(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) GroupingElement() (localctx IGroupingElementContext) {
	localctx = NewGroupingElementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, SQLParserRULE_groupingElement)
	p.SetState(399)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSingleGroupingSetContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(397)
			p.GroupingSet()
		}

	case 2:
		localctx = NewGroupByAllColumnsContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(398)
			p.Match(SQLParserALL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupingSetContext is an interface to support dynamic dispatch.
type IGroupingSetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LR_BRACKET() antlr.TerminalNode
	RR_BRACKET() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGroupingSetContext differentiates from other interfaces.
	IsGroupingSetContext()
}

type GroupingSetContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupingSetContext() *GroupingSetContext {
	var p = new(GroupingSetContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_groupingSet
	return p
}

func InitEmptyGroupingSetContext(p *GroupingSetContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_groupingSet
}

func (*GroupingSetContext) IsGroupingSetContext() {}

func NewGroupingSetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupingSetContext {
	var p = new(GroupingSetContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_groupingSet

	return p
}

func (s *GroupingSetContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupingSetContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *GroupingSetContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *GroupingSetContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *GroupingSetContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *GroupingSetContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *GroupingSetContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *GroupingSetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupingSetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupingSetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterGroupingSet(s)
	}
}

func (s *GroupingSetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitGroupingSet(s)
	}
}

func (s *GroupingSetContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitGroupingSet(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) GroupingSet() (localctx IGroupingSetContext) {
	localctx = NewGroupingSetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, SQLParserRULE_groupingSet)
	var _la int

	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(401)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(410)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-144) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&34246558207) != 0) {
			{
				p.SetState(402)
				p.Expression()
			}
			p.SetState(407)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == SQLParserCOMMA {
				{
					p.SetState(403)
					p.Match(SQLParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(404)
					p.Expression()
				}

				p.SetState(409)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(412)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(413)
			p.Expression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHavingContext is an interface to support dynamic dispatch.
type IHavingContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BooleanExpression() IBooleanExpressionContext

	// IsHavingContext differentiates from other interfaces.
	IsHavingContext()
}

type HavingContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHavingContext() *HavingContext {
	var p = new(HavingContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_having
	return p
}

func InitEmptyHavingContext(p *HavingContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_having
}

func (*HavingContext) IsHavingContext() {}

func NewHavingContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HavingContext {
	var p = new(HavingContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_having

	return p
}

func (s *HavingContext) GetParser() antlr.Parser { return s.parser }

func (s *HavingContext) BooleanExpression() IBooleanExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *HavingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HavingContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HavingContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterHaving(s)
	}
}

func (s *HavingContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitHaving(s)
	}
}

func (s *HavingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitHaving(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Having() (localctx IHavingContext) {
	localctx = NewHavingContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, SQLParserRULE_having)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(416)
		p.booleanExpression(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOrderByContext is an interface to support dynamic dispatch.
type IOrderByContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSortItem() []ISortItemContext
	SortItem(i int) ISortItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsOrderByContext differentiates from other interfaces.
	IsOrderByContext()
}

type OrderByContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrderByContext() *OrderByContext {
	var p = new(OrderByContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_orderBy
	return p
}

func InitEmptyOrderByContext(p *OrderByContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_orderBy
}

func (*OrderByContext) IsOrderByContext() {}

func NewOrderByContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByContext {
	var p = new(OrderByContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_orderBy

	return p
}

func (s *OrderByContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByContext) AllSortItem() []ISortItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISortItemContext); ok {
			len++
		}
	}

	tst := make([]ISortItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISortItemContext); ok {
			tst[i] = t.(ISortItemContext)
			i++
		}
	}

	return tst
}

func (s *OrderByContext) SortItem(i int) ISortItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISortItemContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISortItemContext)
}

func (s *OrderByContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *OrderByContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *OrderByContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrderByContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterOrderBy(s)
	}
}

func (s *OrderByContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitOrderBy(s)
	}
}

func (s *OrderByContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitOrderBy(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) OrderBy() (localctx IOrderByContext) {
	localctx = NewOrderByContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, SQLParserRULE_orderBy)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(418)
		p.SortItem()
	}
	p.SetState(423)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(419)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(420)
			p.SortItem()
		}

		p.SetState(425)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISortItemContext is an interface to support dynamic dispatch.
type ISortItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOrdering returns the ordering token.
	GetOrdering() antlr.Token

	// SetOrdering sets the ordering token.
	SetOrdering(antlr.Token)

	// Getter signatures
	Expression() IExpressionContext
	ASC() antlr.TerminalNode
	DESC() antlr.TerminalNode

	// IsSortItemContext differentiates from other interfaces.
	IsSortItemContext()
}

type SortItemContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	ordering antlr.Token
}

func NewEmptySortItemContext() *SortItemContext {
	var p = new(SortItemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_sortItem
	return p
}

func InitEmptySortItemContext(p *SortItemContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_sortItem
}

func (*SortItemContext) IsSortItemContext() {}

func NewSortItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SortItemContext {
	var p = new(SortItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_sortItem

	return p
}

func (s *SortItemContext) GetParser() antlr.Parser { return s.parser }

func (s *SortItemContext) GetOrdering() antlr.Token { return s.ordering }

func (s *SortItemContext) SetOrdering(v antlr.Token) { s.ordering = v }

func (s *SortItemContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SortItemContext) ASC() antlr.TerminalNode {
	return s.GetToken(SQLParserASC, 0)
}

func (s *SortItemContext) DESC() antlr.TerminalNode {
	return s.GetToken(SQLParserDESC, 0)
}

func (s *SortItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SortItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SortItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterSortItem(s)
	}
}

func (s *SortItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitSortItem(s)
	}
}

func (s *SortItemContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitSortItem(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) SortItem() (localctx ISortItemContext) {
	localctx = NewSortItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, SQLParserRULE_sortItem)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(426)
		p.Expression()
	}
	p.SetState(428)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserASC || _la == SQLParserDESC {
		{
			p.SetState(427)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*SortItemContext).ordering = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserASC || _la == SQLParserDESC) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*SortItemContext).ordering = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILimitRowCountContext is an interface to support dynamic dispatch.
type ILimitRowCountContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTEGER_VALUE() antlr.TerminalNode

	// IsLimitRowCountContext differentiates from other interfaces.
	IsLimitRowCountContext()
}

type LimitRowCountContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitRowCountContext() *LimitRowCountContext {
	var p = new(LimitRowCountContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_limitRowCount
	return p
}

func InitEmptyLimitRowCountContext(p *LimitRowCountContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_limitRowCount
}

func (*LimitRowCountContext) IsLimitRowCountContext() {}

func NewLimitRowCountContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitRowCountContext {
	var p = new(LimitRowCountContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_limitRowCount

	return p
}

func (s *LimitRowCountContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitRowCountContext) INTEGER_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserINTEGER_VALUE, 0)
}

func (s *LimitRowCountContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitRowCountContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitRowCountContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterLimitRowCount(s)
	}
}

func (s *LimitRowCountContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitLimitRowCount(s)
	}
}

func (s *LimitRowCountContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitLimitRowCount(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) LimitRowCount() (localctx ILimitRowCountContext) {
	localctx = NewLimitRowCountContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, SQLParserRULE_limitRowCount)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(430)
		p.Match(SQLParserINTEGER_VALUE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BooleanExpression() IBooleanExpressionContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) BooleanExpression() IBooleanExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (s *ExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, SQLParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(432)
		p.booleanExpression(0)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBooleanExpressionContext is an interface to support dynamic dispatch.
type IBooleanExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsBooleanExpressionContext differentiates from other interfaces.
	IsBooleanExpressionContext()
}

type BooleanExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanExpressionContext() *BooleanExpressionContext {
	var p = new(BooleanExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_booleanExpression
	return p
}

func InitEmptyBooleanExpressionContext(p *BooleanExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_booleanExpression
}

func (*BooleanExpressionContext) IsBooleanExpressionContext() {}

func NewBooleanExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanExpressionContext {
	var p = new(BooleanExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_booleanExpression

	return p
}

func (s *BooleanExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanExpressionContext) CopyAll(ctx *BooleanExpressionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *BooleanExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LogicalNotContext struct {
	BooleanExpressionContext
	notOperator antlr.Token
}

func NewLogicalNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalNotContext {
	var p = new(LogicalNotContext)

	InitEmptyBooleanExpressionContext(&p.BooleanExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*BooleanExpressionContext))

	return p
}

func (s *LogicalNotContext) GetNotOperator() antlr.Token { return s.notOperator }

func (s *LogicalNotContext) SetNotOperator(v antlr.Token) { s.notOperator = v }

func (s *LogicalNotContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalNotContext) BooleanExpression() IBooleanExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *LogicalNotContext) NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserNOT, 0)
}

func (s *LogicalNotContext) EXCLAMATION_SYMBOL() antlr.TerminalNode {
	return s.GetToken(SQLParserEXCLAMATION_SYMBOL, 0)
}

func (s *LogicalNotContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterLogicalNot(s)
	}
}

func (s *LogicalNotContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitLogicalNot(s)
	}
}

func (s *LogicalNotContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitLogicalNot(s)

	default:
		return t.VisitChildren(s)
	}
}

type PredicatedExpressionContext struct {
	BooleanExpressionContext
}

func NewPredicatedExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PredicatedExpressionContext {
	var p = new(PredicatedExpressionContext)

	InitEmptyBooleanExpressionContext(&p.BooleanExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*BooleanExpressionContext))

	return p
}

func (s *PredicatedExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicatedExpressionContext) Predicate() IPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPredicateContext)
}

func (s *PredicatedExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterPredicatedExpression(s)
	}
}

func (s *PredicatedExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitPredicatedExpression(s)
	}
}

func (s *PredicatedExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitPredicatedExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type OrContext struct {
	BooleanExpressionContext
}

func NewOrContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OrContext {
	var p = new(OrContext)

	InitEmptyBooleanExpressionContext(&p.BooleanExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*BooleanExpressionContext))

	return p
}

func (s *OrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrContext) AllBooleanExpression() []IBooleanExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			len++
		}
	}

	tst := make([]IBooleanExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBooleanExpressionContext); ok {
			tst[i] = t.(IBooleanExpressionContext)
			i++
		}
	}

	return tst
}

func (s *OrContext) BooleanExpression(i int) IBooleanExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *OrContext) OR() antlr.TerminalNode {
	return s.GetToken(SQLParserOR, 0)
}

func (s *OrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterOr(s)
	}
}

func (s *OrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitOr(s)
	}
}

func (s *OrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitOr(s)

	default:
		return t.VisitChildren(s)
	}
}

type AndContext struct {
	BooleanExpressionContext
}

func NewAndContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AndContext {
	var p = new(AndContext)

	InitEmptyBooleanExpressionContext(&p.BooleanExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*BooleanExpressionContext))

	return p
}

func (s *AndContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndContext) AllBooleanExpression() []IBooleanExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			len++
		}
	}

	tst := make([]IBooleanExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBooleanExpressionContext); ok {
			tst[i] = t.(IBooleanExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AndContext) BooleanExpression(i int) IBooleanExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *AndContext) AND() antlr.TerminalNode {
	return s.GetToken(SQLParserAND, 0)
}

func (s *AndContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterAnd(s)
	}
}

func (s *AndContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitAnd(s)
	}
}

func (s *AndContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BooleanExpression() (localctx IBooleanExpressionContext) {
	return p.booleanExpression(0)
}

func (p *SQLParser) booleanExpression(_p int) (localctx IBooleanExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewBooleanExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IBooleanExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 74
	p.EnterRecursionRule(localctx, 74, SQLParserRULE_booleanExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(438)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		localctx = NewLogicalNotContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(435)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*LogicalNotContext).notOperator = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserNOT || _la == SQLParserEXCLAMATION_SYMBOL) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*LogicalNotContext).notOperator = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(436)
			p.booleanExpression(4)
		}

	case 2:
		localctx = NewPredicatedExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(437)
			p.Predicate()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(448)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(446)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 45, p.GetParserRuleContext()) {
			case 1:
				localctx = NewAndContext(p, NewBooleanExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_booleanExpression)
				p.SetState(440)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(441)
					p.Match(SQLParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(442)
					p.booleanExpression(4)
				}

			case 2:
				localctx = NewOrContext(p, NewBooleanExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_booleanExpression)
				p.SetState(443)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(444)
					p.Match(SQLParserOR)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(445)
					p.booleanExpression(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(450)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueExpressionContext is an interface to support dynamic dispatch.
type IValueExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsValueExpressionContext differentiates from other interfaces.
	IsValueExpressionContext()
}

type ValueExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueExpressionContext() *ValueExpressionContext {
	var p = new(ValueExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_valueExpression
	return p
}

func InitEmptyValueExpressionContext(p *ValueExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_valueExpression
}

func (*ValueExpressionContext) IsValueExpressionContext() {}

func NewValueExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueExpressionContext {
	var p = new(ValueExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_valueExpression

	return p
}

func (s *ValueExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueExpressionContext) CopyAll(ctx *ValueExpressionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ValueExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ValueExpressionDefaultContext struct {
	ValueExpressionContext
}

func NewValueExpressionDefaultContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ValueExpressionDefaultContext {
	var p = new(ValueExpressionDefaultContext)

	InitEmptyValueExpressionContext(&p.ValueExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ValueExpressionContext))

	return p
}

func (s *ValueExpressionDefaultContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueExpressionDefaultContext) PrimaryExpression() IPrimaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionContext)
}

func (s *ValueExpressionDefaultContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterValueExpressionDefault(s)
	}
}

func (s *ValueExpressionDefaultContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitValueExpressionDefault(s)
	}
}

func (s *ValueExpressionDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitValueExpressionDefault(s)

	default:
		return t.VisitChildren(s)
	}
}

type ArithmeticBinaryContext struct {
	ValueExpressionContext
	left     IValueExpressionContext
	operator antlr.Token
	right    IValueExpressionContext
}

func NewArithmeticBinaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ArithmeticBinaryContext {
	var p = new(ArithmeticBinaryContext)

	InitEmptyValueExpressionContext(&p.ValueExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ValueExpressionContext))

	return p
}

func (s *ArithmeticBinaryContext) GetOperator() antlr.Token { return s.operator }

func (s *ArithmeticBinaryContext) SetOperator(v antlr.Token) { s.operator = v }

func (s *ArithmeticBinaryContext) GetLeft() IValueExpressionContext { return s.left }

func (s *ArithmeticBinaryContext) GetRight() IValueExpressionContext { return s.right }

func (s *ArithmeticBinaryContext) SetLeft(v IValueExpressionContext) { s.left = v }

func (s *ArithmeticBinaryContext) SetRight(v IValueExpressionContext) { s.right = v }

func (s *ArithmeticBinaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArithmeticBinaryContext) AllValueExpression() []IValueExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueExpressionContext); ok {
			len++
		}
	}

	tst := make([]IValueExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueExpressionContext); ok {
			tst[i] = t.(IValueExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ArithmeticBinaryContext) ValueExpression(i int) IValueExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueExpressionContext)
}

func (s *ArithmeticBinaryContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(SQLParserASTERISK, 0)
}

func (s *ArithmeticBinaryContext) SLASH() antlr.TerminalNode {
	return s.GetToken(SQLParserSLASH, 0)
}

func (s *ArithmeticBinaryContext) PERCENT() antlr.TerminalNode {
	return s.GetToken(SQLParserPERCENT, 0)
}

func (s *ArithmeticBinaryContext) PLUS() antlr.TerminalNode {
	return s.GetToken(SQLParserPLUS, 0)
}

func (s *ArithmeticBinaryContext) MINUS() antlr.TerminalNode {
	return s.GetToken(SQLParserMINUS, 0)
}

func (s *ArithmeticBinaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterArithmeticBinary(s)
	}
}

func (s *ArithmeticBinaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitArithmeticBinary(s)
	}
}

func (s *ArithmeticBinaryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitArithmeticBinary(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ValueExpression() (localctx IValueExpressionContext) {
	return p.valueExpression(0)
}

func (p *SQLParser) valueExpression(_p int) (localctx IValueExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewValueExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IValueExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 76
	p.EnterRecursionRule(localctx, 76, SQLParserRULE_valueExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewValueExpressionDefaultContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(452)
		p.primaryExpression(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(462)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(460)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 47, p.GetParserRuleContext()) {
			case 1:
				localctx = NewArithmeticBinaryContext(p, NewValueExpressionContext(p, _parentctx, _parentState))
				localctx.(*ArithmeticBinaryContext).left = _prevctx

				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_valueExpression)
				p.SetState(454)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(455)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ArithmeticBinaryContext).operator = _lt

					_la = p.GetTokenStream().LA(1)

					if !((int64((_la-81)) & ^0x3f) == 0 && ((int64(1)<<(_la-81))&7) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ArithmeticBinaryContext).operator = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(456)

					var _x = p.valueExpression(3)

					localctx.(*ArithmeticBinaryContext).right = _x
				}

			case 2:
				localctx = NewArithmeticBinaryContext(p, NewValueExpressionContext(p, _parentctx, _parentState))
				localctx.(*ArithmeticBinaryContext).left = _prevctx

				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_valueExpression)
				p.SetState(457)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(458)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ArithmeticBinaryContext).operator = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == SQLParserPLUS || _la == SQLParserMINUS) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ArithmeticBinaryContext).operator = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(459)

					var _x = p.valueExpression(2)

					localctx.(*ArithmeticBinaryContext).right = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(464)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimaryExpressionContext is an interface to support dynamic dispatch.
type IPrimaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPrimaryExpressionContext differentiates from other interfaces.
	IsPrimaryExpressionContext()
}

type PrimaryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExpressionContext() *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_primaryExpression
	return p
}

func InitEmptyPrimaryExpressionContext(p *PrimaryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_primaryExpression
}

func (*PrimaryExpressionContext) IsPrimaryExpressionContext() {}

func NewPrimaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_primaryExpression

	return p
}

func (s *PrimaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExpressionContext) CopyAll(ctx *PrimaryExpressionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PrimaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DereferenceContext struct {
	PrimaryExpressionContext
	base      IPrimaryExpressionContext
	fieldName IIdentifierContext
}

func NewDereferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DereferenceContext {
	var p = new(DereferenceContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *DereferenceContext) GetBase() IPrimaryExpressionContext { return s.base }

func (s *DereferenceContext) GetFieldName() IIdentifierContext { return s.fieldName }

func (s *DereferenceContext) SetBase(v IPrimaryExpressionContext) { s.base = v }

func (s *DereferenceContext) SetFieldName(v IIdentifierContext) { s.fieldName = v }

func (s *DereferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DereferenceContext) DOT() antlr.TerminalNode {
	return s.GetToken(SQLParserDOT, 0)
}

func (s *DereferenceContext) PrimaryExpression() IPrimaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionContext)
}

func (s *DereferenceContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *DereferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDereference(s)
	}
}

func (s *DereferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDereference(s)
	}
}

func (s *DereferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDereference(s)

	default:
		return t.VisitChildren(s)
	}
}

type ColumnReferenceContext struct {
	PrimaryExpressionContext
}

func NewColumnReferenceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ColumnReferenceContext {
	var p = new(ColumnReferenceContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *ColumnReferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnReferenceContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *ColumnReferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnReference(s)
	}
}

func (s *ColumnReferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnReference(s)
	}
}

func (s *ColumnReferenceContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnReference(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringLiteralContext struct {
	PrimaryExpressionContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) String_() IStringContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringContext)
}

func (s *StringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterStringLiteral(s)
	}
}

func (s *StringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitStringLiteral(s)
	}
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionCallContext struct {
	PrimaryExpressionContext
}

func NewFunctionCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionCallContext {
	var p = new(FunctionCallContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *FunctionCallContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *FunctionCallContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *FunctionCallContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FunctionCallContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *FunctionCallContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (s *FunctionCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitFunctionCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenExpressionContext struct {
	PrimaryExpressionContext
}

func NewParenExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenExpressionContext {
	var p = new(ParenExpressionContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *ParenExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenExpressionContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *ParenExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenExpressionContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *ParenExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterParenExpression(s)
	}
}

func (s *ParenExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitParenExpression(s)
	}
}

func (s *ParenExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitParenExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type NumericLiteralContext struct {
	PrimaryExpressionContext
}

func NewNumericLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NumericLiteralContext {
	var p = new(NumericLiteralContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *NumericLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumericLiteralContext) Number() INumberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumberContext)
}

func (s *NumericLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterNumericLiteral(s)
	}
}

func (s *NumericLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitNumericLiteral(s)
	}
}

func (s *NumericLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitNumericLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type BooleanLiteralContext struct {
	PrimaryExpressionContext
}

func NewBooleanLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BooleanLiteralContext {
	var p = new(BooleanLiteralContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *BooleanLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanLiteralContext) BooleanValue() IBooleanValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanValueContext)
}

func (s *BooleanLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBooleanLiteral(s)
	}
}

func (s *BooleanLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBooleanLiteral(s)
	}
}

func (s *BooleanLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBooleanLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) PrimaryExpression() (localctx IPrimaryExpressionContext) {
	return p.primaryExpression(0)
}

func (p *SQLParser) primaryExpression(_p int) (localctx IPrimaryExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewPrimaryExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPrimaryExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 78
	p.EnterRecursionRule(localctx, 78, SQLParserRULE_primaryExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(488)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
	case 1:
		localctx = NewNumericLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(466)
			p.Number()
		}

	case 2:
		localctx = NewBooleanLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(467)
			p.BooleanValue()
		}

	case 3:
		localctx = NewStringLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(468)
			p.String_()
		}

	case 4:
		localctx = NewFunctionCallContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(469)
			p.QualifiedName()
		}
		{
			p.SetState(470)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(479)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-144) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&34246558207) != 0) {
			{
				p.SetState(471)
				p.Expression()
			}
			p.SetState(476)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == SQLParserCOMMA {
				{
					p.SetState(472)
					p.Match(SQLParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(473)
					p.Expression()
				}

				p.SetState(478)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(481)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewColumnReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(483)
			p.Identifier()
		}

	case 6:
		localctx = NewParenExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(484)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(485)
			p.Expression()
		}
		{
			p.SetState(486)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(495)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 52, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewDereferenceContext(p, NewPrimaryExpressionContext(p, _parentctx, _parentState))
			localctx.(*DereferenceContext).base = _prevctx

			p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_primaryExpression)
			p.SetState(490)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				goto errorExit
			}
			{
				p.SetState(491)
				p.Match(SQLParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(492)

				var _x = p.Identifier()

				localctx.(*DereferenceContext).fieldName = _x
			}

		}
		p.SetState(497)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 52, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPredicateContext is an interface to support dynamic dispatch.
type IPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPredicateContext differentiates from other interfaces.
	IsPredicateContext()
}

type PredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPredicateContext() *PredicateContext {
	var p = new(PredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_predicate
	return p
}

func InitEmptyPredicateContext(p *PredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_predicate
}

func (*PredicateContext) IsPredicateContext() {}

func NewPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PredicateContext {
	var p = new(PredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_predicate

	return p
}

func (s *PredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *PredicateContext) CopyAll(ctx *PredicateContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ValueExpressionPredicateContext struct {
	PredicateContext
}

func NewValueExpressionPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ValueExpressionPredicateContext {
	var p = new(ValueExpressionPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *ValueExpressionPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueExpressionPredicateContext) ValueExpression() IValueExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueExpressionContext)
}

func (s *ValueExpressionPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterValueExpressionPredicate(s)
	}
}

func (s *ValueExpressionPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitValueExpressionPredicate(s)
	}
}

func (s *ValueExpressionPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitValueExpressionPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

type BinaryComparisonPredicateContext struct {
	PredicateContext
	left     IValueExpressionContext
	operator IComparisonOperatorContext
	right    IValueExpressionContext
}

func NewBinaryComparisonPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BinaryComparisonPredicateContext {
	var p = new(BinaryComparisonPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *BinaryComparisonPredicateContext) GetLeft() IValueExpressionContext { return s.left }

func (s *BinaryComparisonPredicateContext) GetOperator() IComparisonOperatorContext {
	return s.operator
}

func (s *BinaryComparisonPredicateContext) GetRight() IValueExpressionContext { return s.right }

func (s *BinaryComparisonPredicateContext) SetLeft(v IValueExpressionContext) { s.left = v }

func (s *BinaryComparisonPredicateContext) SetOperator(v IComparisonOperatorContext) { s.operator = v }

func (s *BinaryComparisonPredicateContext) SetRight(v IValueExpressionContext) { s.right = v }

func (s *BinaryComparisonPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BinaryComparisonPredicateContext) AllValueExpression() []IValueExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueExpressionContext); ok {
			len++
		}
	}

	tst := make([]IValueExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueExpressionContext); ok {
			tst[i] = t.(IValueExpressionContext)
			i++
		}
	}

	return tst
}

func (s *BinaryComparisonPredicateContext) ValueExpression(i int) IValueExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueExpressionContext)
}

func (s *BinaryComparisonPredicateContext) ComparisonOperator() IComparisonOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonOperatorContext)
}

func (s *BinaryComparisonPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBinaryComparisonPredicate(s)
	}
}

func (s *BinaryComparisonPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBinaryComparisonPredicate(s)
	}
}

func (s *BinaryComparisonPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBinaryComparisonPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

type InPredicateContext struct {
	PredicateContext
	left IValueExpressionContext
}

func NewInPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InPredicateContext {
	var p = new(InPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *InPredicateContext) GetLeft() IValueExpressionContext { return s.left }

func (s *InPredicateContext) SetLeft(v IValueExpressionContext) { s.left = v }

func (s *InPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InPredicateContext) IN() antlr.TerminalNode {
	return s.GetToken(SQLParserIN, 0)
}

func (s *InPredicateContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *InPredicateContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *InPredicateContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *InPredicateContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *InPredicateContext) ValueExpression() IValueExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueExpressionContext)
}

func (s *InPredicateContext) NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserNOT, 0)
}

func (s *InPredicateContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *InPredicateContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *InPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterInPredicate(s)
	}
}

func (s *InPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitInPredicate(s)
	}
}

func (s *InPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitInPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

type LikePredicateContext struct {
	PredicateContext
	left    IValueExpressionContext
	pattern IValueExpressionContext
	escape  IValueExpressionContext
}

func NewLikePredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LikePredicateContext {
	var p = new(LikePredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *LikePredicateContext) GetLeft() IValueExpressionContext { return s.left }

func (s *LikePredicateContext) GetPattern() IValueExpressionContext { return s.pattern }

func (s *LikePredicateContext) GetEscape() IValueExpressionContext { return s.escape }

func (s *LikePredicateContext) SetLeft(v IValueExpressionContext) { s.left = v }

func (s *LikePredicateContext) SetPattern(v IValueExpressionContext) { s.pattern = v }

func (s *LikePredicateContext) SetEscape(v IValueExpressionContext) { s.escape = v }

func (s *LikePredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LikePredicateContext) LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserLIKE, 0)
}

func (s *LikePredicateContext) AllValueExpression() []IValueExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueExpressionContext); ok {
			len++
		}
	}

	tst := make([]IValueExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueExpressionContext); ok {
			tst[i] = t.(IValueExpressionContext)
			i++
		}
	}

	return tst
}

func (s *LikePredicateContext) ValueExpression(i int) IValueExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueExpressionContext)
}

func (s *LikePredicateContext) NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserNOT, 0)
}

func (s *LikePredicateContext) ESCAPE() antlr.TerminalNode {
	return s.GetToken(SQLParserESCAPE, 0)
}

func (s *LikePredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterLikePredicate(s)
	}
}

func (s *LikePredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitLikePredicate(s)
	}
}

func (s *LikePredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitLikePredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

type RegexpPredicateContext struct {
	PredicateContext
	left     IValueExpressionContext
	operator antlr.Token
	pattern  IValueExpressionContext
}

func NewRegexpPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RegexpPredicateContext {
	var p = new(RegexpPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *RegexpPredicateContext) GetOperator() antlr.Token { return s.operator }

func (s *RegexpPredicateContext) SetOperator(v antlr.Token) { s.operator = v }

func (s *RegexpPredicateContext) GetLeft() IValueExpressionContext { return s.left }

func (s *RegexpPredicateContext) GetPattern() IValueExpressionContext { return s.pattern }

func (s *RegexpPredicateContext) SetLeft(v IValueExpressionContext) { s.left = v }

func (s *RegexpPredicateContext) SetPattern(v IValueExpressionContext) { s.pattern = v }

func (s *RegexpPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RegexpPredicateContext) AllValueExpression() []IValueExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueExpressionContext); ok {
			len++
		}
	}

	tst := make([]IValueExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueExpressionContext); ok {
			tst[i] = t.(IValueExpressionContext)
			i++
		}
	}

	return tst
}

func (s *RegexpPredicateContext) ValueExpression(i int) IValueExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueExpressionContext)
}

func (s *RegexpPredicateContext) REGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserREGEXP, 0)
}

func (s *RegexpPredicateContext) NEQREGEXP() antlr.TerminalNode {
	return s.GetToken(SQLParserNEQREGEXP, 0)
}

func (s *RegexpPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterRegexpPredicate(s)
	}
}

func (s *RegexpPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitRegexpPredicate(s)
	}
}

func (s *RegexpPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitRegexpPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Predicate() (localctx IPredicateContext) {
	localctx = NewPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, SQLParserRULE_predicate)
	var _la int

	p.SetState(534)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext()) {
	case 1:
		localctx = NewBinaryComparisonPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(498)

			var _x = p.valueExpression(0)

			localctx.(*BinaryComparisonPredicateContext).left = _x
		}
		{
			p.SetState(499)

			var _x = p.ComparisonOperator()

			localctx.(*BinaryComparisonPredicateContext).operator = _x
		}
		{
			p.SetState(500)

			var _x = p.valueExpression(0)

			localctx.(*BinaryComparisonPredicateContext).right = _x
		}

	case 2:
		localctx = NewInPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(502)

			var _x = p.valueExpression(0)

			localctx.(*InPredicateContext).left = _x
		}
		p.SetState(504)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserNOT {
			{
				p.SetState(503)
				p.Match(SQLParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(506)
			p.Match(SQLParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(507)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(508)
			p.Expression()
		}
		p.SetState(513)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(509)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(510)
				p.Expression()
			}

			p.SetState(515)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(516)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewLikePredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(518)

			var _x = p.valueExpression(0)

			localctx.(*LikePredicateContext).left = _x
		}
		p.SetState(520)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserNOT {
			{
				p.SetState(519)
				p.Match(SQLParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(522)
			p.Match(SQLParserLIKE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(523)

			var _x = p.valueExpression(0)

			localctx.(*LikePredicateContext).pattern = _x
		}
		p.SetState(526)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 56, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(524)
				p.Match(SQLParserESCAPE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(525)

				var _x = p.valueExpression(0)

				localctx.(*LikePredicateContext).escape = _x
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 4:
		localctx = NewRegexpPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(528)

			var _x = p.valueExpression(0)

			localctx.(*RegexpPredicateContext).left = _x
		}
		{
			p.SetState(529)

			var _lt = p.GetTokenStream().LT(1)

			localctx.(*RegexpPredicateContext).operator = _lt

			_la = p.GetTokenStream().LA(1)

			if !(_la == SQLParserREGEXP || _la == SQLParserNEQREGEXP) {
				var _ri = p.GetErrorHandler().RecoverInline(p)

				localctx.(*RegexpPredicateContext).operator = _ri
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(531)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(530)

				var _x = p.valueExpression(0)

				localctx.(*RegexpPredicateContext).pattern = _x
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 5:
		localctx = NewValueExpressionPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(533)
			p.valueExpression(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonOperatorContext is an interface to support dynamic dispatch.
type IComparisonOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EQ() antlr.TerminalNode
	NEQ() antlr.TerminalNode
	LT() antlr.TerminalNode
	LTE() antlr.TerminalNode
	GT() antlr.TerminalNode
	GTE() antlr.TerminalNode

	// IsComparisonOperatorContext differentiates from other interfaces.
	IsComparisonOperatorContext()
}

type ComparisonOperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonOperatorContext() *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_comparisonOperator
	return p
}

func InitEmptyComparisonOperatorContext(p *ComparisonOperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_comparisonOperator
}

func (*ComparisonOperatorContext) IsComparisonOperatorContext() {}

func NewComparisonOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_comparisonOperator

	return p
}

func (s *ComparisonOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonOperatorContext) EQ() antlr.TerminalNode {
	return s.GetToken(SQLParserEQ, 0)
}

func (s *ComparisonOperatorContext) NEQ() antlr.TerminalNode {
	return s.GetToken(SQLParserNEQ, 0)
}

func (s *ComparisonOperatorContext) LT() antlr.TerminalNode {
	return s.GetToken(SQLParserLT, 0)
}

func (s *ComparisonOperatorContext) LTE() antlr.TerminalNode {
	return s.GetToken(SQLParserLTE, 0)
}

func (s *ComparisonOperatorContext) GT() antlr.TerminalNode {
	return s.GetToken(SQLParserGT, 0)
}

func (s *ComparisonOperatorContext) GTE() antlr.TerminalNode {
	return s.GetToken(SQLParserGTE, 0)
}

func (s *ComparisonOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterComparisonOperator(s)
	}
}

func (s *ComparisonOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitComparisonOperator(s)
	}
}

func (s *ComparisonOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitComparisonOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ComparisonOperator() (localctx IComparisonOperatorContext) {
	localctx = NewComparisonOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, SQLParserRULE_comparisonOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(536)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-73)) & ^0x3f) == 0 && ((int64(1)<<(_la-73))&63) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterContext is an interface to support dynamic dispatch.
type IFilterContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FILTER() antlr.TerminalNode
	LR_BRACKET() antlr.TerminalNode
	WHERE() antlr.TerminalNode
	BooleanExpression() IBooleanExpressionContext
	RR_BRACKET() antlr.TerminalNode

	// IsFilterContext differentiates from other interfaces.
	IsFilterContext()
}

type FilterContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterContext() *FilterContext {
	var p = new(FilterContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_filter
	return p
}

func InitEmptyFilterContext(p *FilterContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_filter
}

func (*FilterContext) IsFilterContext() {}

func NewFilterContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterContext {
	var p = new(FilterContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_filter

	return p
}

func (s *FilterContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterContext) FILTER() antlr.TerminalNode {
	return s.GetToken(SQLParserFILTER, 0)
}

func (s *FilterContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *FilterContext) WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserWHERE, 0)
}

func (s *FilterContext) BooleanExpression() IBooleanExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBooleanExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBooleanExpressionContext)
}

func (s *FilterContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *FilterContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterFilter(s)
	}
}

func (s *FilterContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitFilter(s)
	}
}

func (s *FilterContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitFilter(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Filter() (localctx IFilterContext) {
	localctx = NewFilterContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, SQLParserRULE_filter)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(538)
		p.Match(SQLParserFILTER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(539)
		p.Match(SQLParserLR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(540)
		p.Match(SQLParserWHERE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(541)
		p.booleanExpression(0)
	}
	{
		p.SetState(542)
		p.Match(SQLParserRR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQualifiedNameContext is an interface to support dynamic dispatch.
type IQualifiedNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsQualifiedNameContext differentiates from other interfaces.
	IsQualifiedNameContext()
}

type QualifiedNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQualifiedNameContext() *QualifiedNameContext {
	var p = new(QualifiedNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_qualifiedName
	return p
}

func InitEmptyQualifiedNameContext(p *QualifiedNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_qualifiedName
}

func (*QualifiedNameContext) IsQualifiedNameContext() {}

func NewQualifiedNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualifiedNameContext {
	var p = new(QualifiedNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_qualifiedName

	return p
}

func (s *QualifiedNameContext) GetParser() antlr.Parser { return s.parser }

func (s *QualifiedNameContext) AllIdentifier() []IIdentifierContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentifierContext); ok {
			len++
		}
	}

	tst := make([]IIdentifierContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentifierContext); ok {
			tst[i] = t.(IIdentifierContext)
			i++
		}
	}

	return tst
}

func (s *QualifiedNameContext) Identifier(i int) IIdentifierContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *QualifiedNameContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(SQLParserDOT)
}

func (s *QualifiedNameContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserDOT, i)
}

func (s *QualifiedNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QualifiedNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QualifiedNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQualifiedName(s)
	}
}

func (s *QualifiedNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQualifiedName(s)
	}
}

func (s *QualifiedNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQualifiedName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) QualifiedName() (localctx IQualifiedNameContext) {
	localctx = NewQualifiedNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, SQLParserRULE_qualifiedName)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(544)
		p.Identifier()
	}
	p.SetState(549)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(545)
				p.Match(SQLParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(546)
				p.Identifier()
			}

		}
		p.SetState(551)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 59, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPropertiesContext is an interface to support dynamic dispatch.
type IPropertiesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LR_BRACKET() antlr.TerminalNode
	PropertyAssignments() IPropertyAssignmentsContext
	RR_BRACKET() antlr.TerminalNode

	// IsPropertiesContext differentiates from other interfaces.
	IsPropertiesContext()
}

type PropertiesContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertiesContext() *PropertiesContext {
	var p = new(PropertiesContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_properties
	return p
}

func InitEmptyPropertiesContext(p *PropertiesContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_properties
}

func (*PropertiesContext) IsPropertiesContext() {}

func NewPropertiesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertiesContext {
	var p = new(PropertiesContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_properties

	return p
}

func (s *PropertiesContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertiesContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *PropertiesContext) PropertyAssignments() IPropertyAssignmentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyAssignmentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyAssignmentsContext)
}

func (s *PropertiesContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *PropertiesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertiesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertiesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterProperties(s)
	}
}

func (s *PropertiesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitProperties(s)
	}
}

func (s *PropertiesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitProperties(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Properties() (localctx IPropertiesContext) {
	localctx = NewPropertiesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, SQLParserRULE_properties)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(552)
		p.Match(SQLParserLR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(553)
		p.PropertyAssignments()
	}
	{
		p.SetState(554)
		p.Match(SQLParserRR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPropertyAssignmentsContext is an interface to support dynamic dispatch.
type IPropertyAssignmentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllProperty() []IPropertyContext
	Property(i int) IPropertyContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsPropertyAssignmentsContext differentiates from other interfaces.
	IsPropertyAssignmentsContext()
}

type PropertyAssignmentsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyAssignmentsContext() *PropertyAssignmentsContext {
	var p = new(PropertyAssignmentsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_propertyAssignments
	return p
}

func InitEmptyPropertyAssignmentsContext(p *PropertyAssignmentsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_propertyAssignments
}

func (*PropertyAssignmentsContext) IsPropertyAssignmentsContext() {}

func NewPropertyAssignmentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyAssignmentsContext {
	var p = new(PropertyAssignmentsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_propertyAssignments

	return p
}

func (s *PropertyAssignmentsContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyAssignmentsContext) AllProperty() []IPropertyContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPropertyContext); ok {
			len++
		}
	}

	tst := make([]IPropertyContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPropertyContext); ok {
			tst[i] = t.(IPropertyContext)
			i++
		}
	}

	return tst
}

func (s *PropertyAssignmentsContext) Property(i int) IPropertyContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyContext)
}

func (s *PropertyAssignmentsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *PropertyAssignmentsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *PropertyAssignmentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyAssignmentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyAssignmentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterPropertyAssignments(s)
	}
}

func (s *PropertyAssignmentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitPropertyAssignments(s)
	}
}

func (s *PropertyAssignmentsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitPropertyAssignments(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) PropertyAssignments() (localctx IPropertyAssignmentsContext) {
	localctx = NewPropertyAssignmentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, SQLParserRULE_propertyAssignments)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(556)
		p.Property()
	}
	p.SetState(561)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(557)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(558)
			p.Property()
		}

		p.SetState(563)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPropertyContext is an interface to support dynamic dispatch.
type IPropertyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetName returns the name rule contexts.
	GetName() IIdentifierContext

	// GetValue returns the value rule contexts.
	GetValue() IPropertyValueContext

	// SetName sets the name rule contexts.
	SetName(IIdentifierContext)

	// SetValue sets the value rule contexts.
	SetValue(IPropertyValueContext)

	// Getter signatures
	EQ() antlr.TerminalNode
	Identifier() IIdentifierContext
	PropertyValue() IPropertyValueContext

	// IsPropertyContext differentiates from other interfaces.
	IsPropertyContext()
}

type PropertyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	name   IIdentifierContext
	value  IPropertyValueContext
}

func NewEmptyPropertyContext() *PropertyContext {
	var p = new(PropertyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_property
	return p
}

func InitEmptyPropertyContext(p *PropertyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_property
}

func (*PropertyContext) IsPropertyContext() {}

func NewPropertyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyContext {
	var p = new(PropertyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_property

	return p
}

func (s *PropertyContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyContext) GetName() IIdentifierContext { return s.name }

func (s *PropertyContext) GetValue() IPropertyValueContext { return s.value }

func (s *PropertyContext) SetName(v IIdentifierContext) { s.name = v }

func (s *PropertyContext) SetValue(v IPropertyValueContext) { s.value = v }

func (s *PropertyContext) EQ() antlr.TerminalNode {
	return s.GetToken(SQLParserEQ, 0)
}

func (s *PropertyContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *PropertyContext) PropertyValue() IPropertyValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyValueContext)
}

func (s *PropertyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PropertyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterProperty(s)
	}
}

func (s *PropertyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitProperty(s)
	}
}

func (s *PropertyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitProperty(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Property() (localctx IPropertyContext) {
	localctx = NewPropertyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, SQLParserRULE_property)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(564)

		var _x = p.Identifier()

		localctx.(*PropertyContext).name = _x
	}
	{
		p.SetState(565)
		p.Match(SQLParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(566)

		var _x = p.PropertyValue()

		localctx.(*PropertyContext).value = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPropertyValueContext is an interface to support dynamic dispatch.
type IPropertyValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPropertyValueContext differentiates from other interfaces.
	IsPropertyValueContext()
}

type PropertyValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPropertyValueContext() *PropertyValueContext {
	var p = new(PropertyValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_propertyValue
	return p
}

func InitEmptyPropertyValueContext(p *PropertyValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_propertyValue
}

func (*PropertyValueContext) IsPropertyValueContext() {}

func NewPropertyValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PropertyValueContext {
	var p = new(PropertyValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_propertyValue

	return p
}

func (s *PropertyValueContext) GetParser() antlr.Parser { return s.parser }

func (s *PropertyValueContext) CopyAll(ctx *PropertyValueContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PropertyValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PropertyValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DefaultPropertyValueContext struct {
	PropertyValueContext
}

func NewDefaultPropertyValueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DefaultPropertyValueContext {
	var p = new(DefaultPropertyValueContext)

	InitEmptyPropertyValueContext(&p.PropertyValueContext)
	p.parser = parser
	p.CopyAll(ctx.(*PropertyValueContext))

	return p
}

func (s *DefaultPropertyValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DefaultPropertyValueContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(SQLParserDEFAULT, 0)
}

func (s *DefaultPropertyValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDefaultPropertyValue(s)
	}
}

func (s *DefaultPropertyValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDefaultPropertyValue(s)
	}
}

func (s *DefaultPropertyValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDefaultPropertyValue(s)

	default:
		return t.VisitChildren(s)
	}
}

type NonDefaultPropertyValueContext struct {
	PropertyValueContext
}

func NewNonDefaultPropertyValueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NonDefaultPropertyValueContext {
	var p = new(NonDefaultPropertyValueContext)

	InitEmptyPropertyValueContext(&p.PropertyValueContext)
	p.parser = parser
	p.CopyAll(ctx.(*PropertyValueContext))

	return p
}

func (s *NonDefaultPropertyValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NonDefaultPropertyValueContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *NonDefaultPropertyValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterNonDefaultPropertyValue(s)
	}
}

func (s *NonDefaultPropertyValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitNonDefaultPropertyValue(s)
	}
}

func (s *NonDefaultPropertyValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitNonDefaultPropertyValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) PropertyValue() (localctx IPropertyValueContext) {
	localctx = NewPropertyValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, SQLParserRULE_propertyValue)
	p.SetState(570)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDefaultPropertyValueContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(568)
			p.Match(SQLParserDEFAULT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewNonDefaultPropertyValueContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(569)
			p.Expression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBooleanValueContext is an interface to support dynamic dispatch.
type IBooleanValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode

	// IsBooleanValueContext differentiates from other interfaces.
	IsBooleanValueContext()
}

type BooleanValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBooleanValueContext() *BooleanValueContext {
	var p = new(BooleanValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_booleanValue
	return p
}

func InitEmptyBooleanValueContext(p *BooleanValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_booleanValue
}

func (*BooleanValueContext) IsBooleanValueContext() {}

func NewBooleanValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BooleanValueContext {
	var p = new(BooleanValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_booleanValue

	return p
}

func (s *BooleanValueContext) GetParser() antlr.Parser { return s.parser }

func (s *BooleanValueContext) TRUE() antlr.TerminalNode {
	return s.GetToken(SQLParserTRUE, 0)
}

func (s *BooleanValueContext) FALSE() antlr.TerminalNode {
	return s.GetToken(SQLParserFALSE, 0)
}

func (s *BooleanValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BooleanValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBooleanValue(s)
	}
}

func (s *BooleanValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBooleanValue(s)
	}
}

func (s *BooleanValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBooleanValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) BooleanValue() (localctx IBooleanValueContext) {
	localctx = NewBooleanValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, SQLParserRULE_booleanValue)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(572)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SQLParserFALSE || _la == SQLParserTRUE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStringContext is an interface to support dynamic dispatch.
type IStringContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsStringContext differentiates from other interfaces.
	IsStringContext()
}

type StringContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringContext() *StringContext {
	var p = new(StringContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_string
	return p
}

func InitEmptyStringContext(p *StringContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_string
}

func (*StringContext) IsStringContext() {}

func NewStringContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringContext {
	var p = new(StringContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_string

	return p
}

func (s *StringContext) GetParser() antlr.Parser { return s.parser }

func (s *StringContext) CopyAll(ctx *StringContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BasicStringLiteralContext struct {
	StringContext
}

func NewBasicStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BasicStringLiteralContext {
	var p = new(BasicStringLiteralContext)

	InitEmptyStringContext(&p.StringContext)
	p.parser = parser
	p.CopyAll(ctx.(*StringContext))

	return p
}

func (s *BasicStringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BasicStringLiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(SQLParserSTRING, 0)
}

func (s *BasicStringLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBasicStringLiteral(s)
	}
}

func (s *BasicStringLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBasicStringLiteral(s)
	}
}

func (s *BasicStringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBasicStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) String_() (localctx IStringContext) {
	localctx = NewStringContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, SQLParserRULE_string)
	localctx = NewBasicStringLiteralContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(574)
		p.Match(SQLParserSTRING)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIdentifierContext is an interface to support dynamic dispatch.
type IIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsIdentifierContext differentiates from other interfaces.
	IsIdentifierContext()
}

type IdentifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifierContext() *IdentifierContext {
	var p = new(IdentifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_identifier
	return p
}

func InitEmptyIdentifierContext(p *IdentifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_identifier
}

func (*IdentifierContext) IsIdentifierContext() {}

func NewIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierContext {
	var p = new(IdentifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_identifier

	return p
}

func (s *IdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifierContext) CopyAll(ctx *IdentifierContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BackQuotedIdentifierContext struct {
	IdentifierContext
}

func NewBackQuotedIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BackQuotedIdentifierContext {
	var p = new(BackQuotedIdentifierContext)

	InitEmptyIdentifierContext(&p.IdentifierContext)
	p.parser = parser
	p.CopyAll(ctx.(*IdentifierContext))

	return p
}

func (s *BackQuotedIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BackQuotedIdentifierContext) BACKQUOTED_IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SQLParserBACKQUOTED_IDENTIFIER, 0)
}

func (s *BackQuotedIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBackQuotedIdentifier(s)
	}
}

func (s *BackQuotedIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBackQuotedIdentifier(s)
	}
}

func (s *BackQuotedIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBackQuotedIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

type QuotedIdentifierContext struct {
	IdentifierContext
}

func NewQuotedIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *QuotedIdentifierContext {
	var p = new(QuotedIdentifierContext)

	InitEmptyIdentifierContext(&p.IdentifierContext)
	p.parser = parser
	p.CopyAll(ctx.(*IdentifierContext))

	return p
}

func (s *QuotedIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuotedIdentifierContext) QUOTED_IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SQLParserQUOTED_IDENTIFIER, 0)
}

func (s *QuotedIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterQuotedIdentifier(s)
	}
}

func (s *QuotedIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitQuotedIdentifier(s)
	}
}

func (s *QuotedIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitQuotedIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

type DigitIdentifierContext struct {
	IdentifierContext
}

func NewDigitIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DigitIdentifierContext {
	var p = new(DigitIdentifierContext)

	InitEmptyIdentifierContext(&p.IdentifierContext)
	p.parser = parser
	p.CopyAll(ctx.(*IdentifierContext))

	return p
}

func (s *DigitIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DigitIdentifierContext) DIGIT_IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SQLParserDIGIT_IDENTIFIER, 0)
}

func (s *DigitIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDigitIdentifier(s)
	}
}

func (s *DigitIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDigitIdentifier(s)
	}
}

func (s *DigitIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDigitIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

type UnquotedIdentifierContext struct {
	IdentifierContext
}

func NewUnquotedIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnquotedIdentifierContext {
	var p = new(UnquotedIdentifierContext)

	InitEmptyIdentifierContext(&p.IdentifierContext)
	p.parser = parser
	p.CopyAll(ctx.(*IdentifierContext))

	return p
}

func (s *UnquotedIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnquotedIdentifierContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SQLParserIDENTIFIER, 0)
}

func (s *UnquotedIdentifierContext) NonReserved() INonReservedContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INonReservedContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INonReservedContext)
}

func (s *UnquotedIdentifierContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterUnquotedIdentifier(s)
	}
}

func (s *UnquotedIdentifierContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitUnquotedIdentifier(s)
	}
}

func (s *UnquotedIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitUnquotedIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Identifier() (localctx IIdentifierContext) {
	localctx = NewIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, SQLParserRULE_identifier)
	p.SetState(581)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserIDENTIFIER:
		localctx = NewUnquotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(576)
			p.Match(SQLParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case SQLParserQUOTED_IDENTIFIER:
		localctx = NewQuotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(577)
			p.Match(SQLParserQUOTED_IDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case SQLParserALL, SQLParserALIVE, SQLParserAND, SQLParserAS, SQLParserASC, SQLParserBROKER, SQLParserBROKERS, SQLParserBY, SQLParserCOMPACT, SQLParserCREATE, SQLParserCROSS, SQLParserDATABASE, SQLParserDATABASES, SQLParserDEFAULT, SQLParserDESC, SQLParserDISTRIBUTED, SQLParserDROP, SQLParserENGINE, SQLParserESCAPE, SQLParserEXPLAIN, SQLParserEXISTS, SQLParserFALSE, SQLParserFIELDS, SQLParserFILTER, SQLParserFLUSH, SQLParserFROM, SQLParserGROUP, SQLParserHAVING, SQLParserIF, SQLParserIN, SQLParserJOIN, SQLParserKEYS, SQLParserLEFT, SQLParserLIKE, SQLParserLIMIT, SQLParserLOG, SQLParserLOGICAL, SQLParserMASTER, SQLParserMETRICS, SQLParserMETRIC, SQLParserMETADATA, SQLParserMETADATAS, SQLParserNAMESPACE, SQLParserNAMESPACES, SQLParserNOT, SQLParserON, SQLParserOR, SQLParserORDER, SQLParserPLAN, SQLParserREQUESTS, SQLParserREPLICATIONS, SQLParserRIGHT, SQLParserROLLUP, SQLParserSELECT, SQLParserSHOW, SQLParserSTATE, SQLParserSTORAGE, SQLParserTAG, SQLParserTRACE, SQLParserTRUE, SQLParserTYPE, SQLParserTYPES, SQLParserVALUES, SQLParserWHERE, SQLParserWITH, SQLParserWITHIN, SQLParserUSING, SQLParserUSE:
		localctx = NewUnquotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(578)
			p.NonReserved()
		}

	case SQLParserBACKQUOTED_IDENTIFIER:
		localctx = NewBackQuotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(579)
			p.Match(SQLParserBACKQUOTED_IDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case SQLParserDIGIT_IDENTIFIER:
		localctx = NewDigitIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(580)
			p.Match(SQLParserDIGIT_IDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INumberContext is an interface to support dynamic dispatch.
type INumberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsNumberContext differentiates from other interfaces.
	IsNumberContext()
}

type NumberContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumberContext() *NumberContext {
	var p = new(NumberContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_number
	return p
}

func InitEmptyNumberContext(p *NumberContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_number
}

func (*NumberContext) IsNumberContext() {}

func NewNumberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumberContext {
	var p = new(NumberContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_number

	return p
}

func (s *NumberContext) GetParser() antlr.Parser { return s.parser }

func (s *NumberContext) CopyAll(ctx *NumberContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *NumberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DecimalLiteralContext struct {
	NumberContext
}

func NewDecimalLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalLiteralContext {
	var p = new(DecimalLiteralContext)

	InitEmptyNumberContext(&p.NumberContext)
	p.parser = parser
	p.CopyAll(ctx.(*NumberContext))

	return p
}

func (s *DecimalLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalLiteralContext) DECIMAL_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserDECIMAL_VALUE, 0)
}

func (s *DecimalLiteralContext) MINUS() antlr.TerminalNode {
	return s.GetToken(SQLParserMINUS, 0)
}

func (s *DecimalLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDecimalLiteral(s)
	}
}

func (s *DecimalLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDecimalLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type DoubleLiteralContext struct {
	NumberContext
}

func NewDoubleLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DoubleLiteralContext {
	var p = new(DoubleLiteralContext)

	InitEmptyNumberContext(&p.NumberContext)
	p.parser = parser
	p.CopyAll(ctx.(*NumberContext))

	return p
}

func (s *DoubleLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DoubleLiteralContext) DOUBLE_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserDOUBLE_VALUE, 0)
}

func (s *DoubleLiteralContext) MINUS() antlr.TerminalNode {
	return s.GetToken(SQLParserMINUS, 0)
}

func (s *DoubleLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDoubleLiteral(s)
	}
}

func (s *DoubleLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDoubleLiteral(s)
	}
}

func (s *DoubleLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDoubleLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type IntegerLiteralContext struct {
	NumberContext
}

func NewIntegerLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntegerLiteralContext {
	var p = new(IntegerLiteralContext)

	InitEmptyNumberContext(&p.NumberContext)
	p.parser = parser
	p.CopyAll(ctx.(*NumberContext))

	return p
}

func (s *IntegerLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegerLiteralContext) INTEGER_VALUE() antlr.TerminalNode {
	return s.GetToken(SQLParserINTEGER_VALUE, 0)
}

func (s *IntegerLiteralContext) MINUS() antlr.TerminalNode {
	return s.GetToken(SQLParserMINUS, 0)
}

func (s *IntegerLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitIntegerLiteral(s)
	}
}

func (s *IntegerLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitIntegerLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Number() (localctx INumberContext) {
	localctx = NewNumberContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, SQLParserRULE_number)
	var _la int

	p.SetState(595)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 66, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(584)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserMINUS {
			{
				p.SetState(583)
				p.Match(SQLParserMINUS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(586)
			p.Match(SQLParserDECIMAL_VALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewDoubleLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(588)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserMINUS {
			{
				p.SetState(587)
				p.Match(SQLParserMINUS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(590)
			p.Match(SQLParserDOUBLE_VALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewIntegerLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(592)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserMINUS {
			{
				p.SetState(591)
				p.Match(SQLParserMINUS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(594)
			p.Match(SQLParserINTEGER_VALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INonReservedContext is an interface to support dynamic dispatch.
type INonReservedContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ALL() antlr.TerminalNode
	ALIVE() antlr.TerminalNode
	AND() antlr.TerminalNode
	AS() antlr.TerminalNode
	ASC() antlr.TerminalNode
	BROKER() antlr.TerminalNode
	BROKERS() antlr.TerminalNode
	BY() antlr.TerminalNode
	COMPACT() antlr.TerminalNode
	CREATE() antlr.TerminalNode
	CROSS() antlr.TerminalNode
	DATABASE() antlr.TerminalNode
	DATABASES() antlr.TerminalNode
	DEFAULT() antlr.TerminalNode
	DESC() antlr.TerminalNode
	DISTRIBUTED() antlr.TerminalNode
	DROP() antlr.TerminalNode
	ENGINE() antlr.TerminalNode
	ESCAPE() antlr.TerminalNode
	EXPLAIN() antlr.TerminalNode
	EXISTS() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	FIELDS() antlr.TerminalNode
	FILTER() antlr.TerminalNode
	FLUSH() antlr.TerminalNode
	FROM() antlr.TerminalNode
	GROUP() antlr.TerminalNode
	HAVING() antlr.TerminalNode
	IF() antlr.TerminalNode
	IN() antlr.TerminalNode
	JOIN() antlr.TerminalNode
	KEYS() antlr.TerminalNode
	LEFT() antlr.TerminalNode
	LIKE() antlr.TerminalNode
	LIMIT() antlr.TerminalNode
	LOG() antlr.TerminalNode
	LOGICAL() antlr.TerminalNode
	MASTER() antlr.TerminalNode
	METRIC() antlr.TerminalNode
	METRICS() antlr.TerminalNode
	METADATA() antlr.TerminalNode
	METADATAS() antlr.TerminalNode
	NAMESPACE() antlr.TerminalNode
	NAMESPACES() antlr.TerminalNode
	NOT() antlr.TerminalNode
	ON() antlr.TerminalNode
	OR() antlr.TerminalNode
	ORDER() antlr.TerminalNode
	PLAN() antlr.TerminalNode
	REQUESTS() antlr.TerminalNode
	REPLICATIONS() antlr.TerminalNode
	RIGHT() antlr.TerminalNode
	ROLLUP() antlr.TerminalNode
	SELECT() antlr.TerminalNode
	SHOW() antlr.TerminalNode
	STATE() antlr.TerminalNode
	STORAGE() antlr.TerminalNode
	TAG() antlr.TerminalNode
	TRACE() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	TYPE() antlr.TerminalNode
	TYPES() antlr.TerminalNode
	VALUES() antlr.TerminalNode
	WHERE() antlr.TerminalNode
	WITH() antlr.TerminalNode
	WITHIN() antlr.TerminalNode
	USING() antlr.TerminalNode
	USE() antlr.TerminalNode

	// IsNonReservedContext differentiates from other interfaces.
	IsNonReservedContext()
}

type NonReservedContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNonReservedContext() *NonReservedContext {
	var p = new(NonReservedContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_nonReserved
	return p
}

func InitEmptyNonReservedContext(p *NonReservedContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_nonReserved
}

func (*NonReservedContext) IsNonReservedContext() {}

func NewNonReservedContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NonReservedContext {
	var p = new(NonReservedContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_nonReserved

	return p
}

func (s *NonReservedContext) GetParser() antlr.Parser { return s.parser }

func (s *NonReservedContext) ALL() antlr.TerminalNode {
	return s.GetToken(SQLParserALL, 0)
}

func (s *NonReservedContext) ALIVE() antlr.TerminalNode {
	return s.GetToken(SQLParserALIVE, 0)
}

func (s *NonReservedContext) AND() antlr.TerminalNode {
	return s.GetToken(SQLParserAND, 0)
}

func (s *NonReservedContext) AS() antlr.TerminalNode {
	return s.GetToken(SQLParserAS, 0)
}

func (s *NonReservedContext) ASC() antlr.TerminalNode {
	return s.GetToken(SQLParserASC, 0)
}

func (s *NonReservedContext) BROKER() antlr.TerminalNode {
	return s.GetToken(SQLParserBROKER, 0)
}

func (s *NonReservedContext) BROKERS() antlr.TerminalNode {
	return s.GetToken(SQLParserBROKERS, 0)
}

func (s *NonReservedContext) BY() antlr.TerminalNode {
	return s.GetToken(SQLParserBY, 0)
}

func (s *NonReservedContext) COMPACT() antlr.TerminalNode {
	return s.GetToken(SQLParserCOMPACT, 0)
}

func (s *NonReservedContext) CREATE() antlr.TerminalNode {
	return s.GetToken(SQLParserCREATE, 0)
}

func (s *NonReservedContext) CROSS() antlr.TerminalNode {
	return s.GetToken(SQLParserCROSS, 0)
}

func (s *NonReservedContext) DATABASE() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASE, 0)
}

func (s *NonReservedContext) DATABASES() antlr.TerminalNode {
	return s.GetToken(SQLParserDATABASES, 0)
}

func (s *NonReservedContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(SQLParserDEFAULT, 0)
}

func (s *NonReservedContext) DESC() antlr.TerminalNode {
	return s.GetToken(SQLParserDESC, 0)
}

func (s *NonReservedContext) DISTRIBUTED() antlr.TerminalNode {
	return s.GetToken(SQLParserDISTRIBUTED, 0)
}

func (s *NonReservedContext) DROP() antlr.TerminalNode {
	return s.GetToken(SQLParserDROP, 0)
}

func (s *NonReservedContext) ENGINE() antlr.TerminalNode {
	return s.GetToken(SQLParserENGINE, 0)
}

func (s *NonReservedContext) ESCAPE() antlr.TerminalNode {
	return s.GetToken(SQLParserESCAPE, 0)
}

func (s *NonReservedContext) EXPLAIN() antlr.TerminalNode {
	return s.GetToken(SQLParserEXPLAIN, 0)
}

func (s *NonReservedContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(SQLParserEXISTS, 0)
}

func (s *NonReservedContext) FALSE() antlr.TerminalNode {
	return s.GetToken(SQLParserFALSE, 0)
}

func (s *NonReservedContext) FIELDS() antlr.TerminalNode {
	return s.GetToken(SQLParserFIELDS, 0)
}

func (s *NonReservedContext) FILTER() antlr.TerminalNode {
	return s.GetToken(SQLParserFILTER, 0)
}

func (s *NonReservedContext) FLUSH() antlr.TerminalNode {
	return s.GetToken(SQLParserFLUSH, 0)
}

func (s *NonReservedContext) FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserFROM, 0)
}

func (s *NonReservedContext) GROUP() antlr.TerminalNode {
	return s.GetToken(SQLParserGROUP, 0)
}

func (s *NonReservedContext) HAVING() antlr.TerminalNode {
	return s.GetToken(SQLParserHAVING, 0)
}

func (s *NonReservedContext) IF() antlr.TerminalNode {
	return s.GetToken(SQLParserIF, 0)
}

func (s *NonReservedContext) IN() antlr.TerminalNode {
	return s.GetToken(SQLParserIN, 0)
}

func (s *NonReservedContext) JOIN() antlr.TerminalNode {
	return s.GetToken(SQLParserJOIN, 0)
}

func (s *NonReservedContext) KEYS() antlr.TerminalNode {
	return s.GetToken(SQLParserKEYS, 0)
}

func (s *NonReservedContext) LEFT() antlr.TerminalNode {
	return s.GetToken(SQLParserLEFT, 0)
}

func (s *NonReservedContext) LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserLIKE, 0)
}

func (s *NonReservedContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(SQLParserLIMIT, 0)
}

func (s *NonReservedContext) LOG() antlr.TerminalNode {
	return s.GetToken(SQLParserLOG, 0)
}

func (s *NonReservedContext) LOGICAL() antlr.TerminalNode {
	return s.GetToken(SQLParserLOGICAL, 0)
}

func (s *NonReservedContext) MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserMASTER, 0)
}

func (s *NonReservedContext) METRIC() antlr.TerminalNode {
	return s.GetToken(SQLParserMETRIC, 0)
}

func (s *NonReservedContext) METRICS() antlr.TerminalNode {
	return s.GetToken(SQLParserMETRICS, 0)
}

func (s *NonReservedContext) METADATA() antlr.TerminalNode {
	return s.GetToken(SQLParserMETADATA, 0)
}

func (s *NonReservedContext) METADATAS() antlr.TerminalNode {
	return s.GetToken(SQLParserMETADATAS, 0)
}

func (s *NonReservedContext) NAMESPACE() antlr.TerminalNode {
	return s.GetToken(SQLParserNAMESPACE, 0)
}

func (s *NonReservedContext) NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserNAMESPACES, 0)
}

func (s *NonReservedContext) NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserNOT, 0)
}

func (s *NonReservedContext) ON() antlr.TerminalNode {
	return s.GetToken(SQLParserON, 0)
}

func (s *NonReservedContext) OR() antlr.TerminalNode {
	return s.GetToken(SQLParserOR, 0)
}

func (s *NonReservedContext) ORDER() antlr.TerminalNode {
	return s.GetToken(SQLParserORDER, 0)
}

func (s *NonReservedContext) PLAN() antlr.TerminalNode {
	return s.GetToken(SQLParserPLAN, 0)
}

func (s *NonReservedContext) REQUESTS() antlr.TerminalNode {
	return s.GetToken(SQLParserREQUESTS, 0)
}

func (s *NonReservedContext) REPLICATIONS() antlr.TerminalNode {
	return s.GetToken(SQLParserREPLICATIONS, 0)
}

func (s *NonReservedContext) RIGHT() antlr.TerminalNode {
	return s.GetToken(SQLParserRIGHT, 0)
}

func (s *NonReservedContext) ROLLUP() antlr.TerminalNode {
	return s.GetToken(SQLParserROLLUP, 0)
}

func (s *NonReservedContext) SELECT() antlr.TerminalNode {
	return s.GetToken(SQLParserSELECT, 0)
}

func (s *NonReservedContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *NonReservedContext) STATE() antlr.TerminalNode {
	return s.GetToken(SQLParserSTATE, 0)
}

func (s *NonReservedContext) STORAGE() antlr.TerminalNode {
	return s.GetToken(SQLParserSTORAGE, 0)
}

func (s *NonReservedContext) TAG() antlr.TerminalNode {
	return s.GetToken(SQLParserTAG, 0)
}

func (s *NonReservedContext) TRACE() antlr.TerminalNode {
	return s.GetToken(SQLParserTRACE, 0)
}

func (s *NonReservedContext) TRUE() antlr.TerminalNode {
	return s.GetToken(SQLParserTRUE, 0)
}

func (s *NonReservedContext) TYPE() antlr.TerminalNode {
	return s.GetToken(SQLParserTYPE, 0)
}

func (s *NonReservedContext) TYPES() antlr.TerminalNode {
	return s.GetToken(SQLParserTYPES, 0)
}

func (s *NonReservedContext) VALUES() antlr.TerminalNode {
	return s.GetToken(SQLParserVALUES, 0)
}

func (s *NonReservedContext) WHERE() antlr.TerminalNode {
	return s.GetToken(SQLParserWHERE, 0)
}

func (s *NonReservedContext) WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserWITH, 0)
}

func (s *NonReservedContext) WITHIN() antlr.TerminalNode {
	return s.GetToken(SQLParserWITHIN, 0)
}

func (s *NonReservedContext) USING() antlr.TerminalNode {
	return s.GetToken(SQLParserUSING, 0)
}

func (s *NonReservedContext) USE() antlr.TerminalNode {
	return s.GetToken(SQLParserUSE, 0)
}

func (s *NonReservedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NonReservedContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NonReservedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterNonReserved(s)
	}
}

func (s *NonReservedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitNonReserved(s)
	}
}

func (s *NonReservedContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitNonReserved(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) NonReserved() (localctx INonReservedContext) {
	localctx = NewNonReservedContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, SQLParserRULE_nonReserved)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(597)
		_la = p.GetTokenStream().LA(1)

		if !(((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-144) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&511) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *SQLParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 24:
		var t *RelationContext = nil
		if localctx != nil {
			t = localctx.(*RelationContext)
		}
		return p.Relation_Sempred(t, predIndex)

	case 37:
		var t *BooleanExpressionContext = nil
		if localctx != nil {
			t = localctx.(*BooleanExpressionContext)
		}
		return p.BooleanExpression_Sempred(t, predIndex)

	case 38:
		var t *ValueExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ValueExpressionContext)
		}
		return p.ValueExpression_Sempred(t, predIndex)

	case 39:
		var t *PrimaryExpressionContext = nil
		if localctx != nil {
			t = localctx.(*PrimaryExpressionContext)
		}
		return p.PrimaryExpression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *SQLParser) Relation_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) BooleanExpression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) ValueExpression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 3:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *SQLParser) PrimaryExpression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 5:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
