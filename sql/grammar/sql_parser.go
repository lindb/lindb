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
		"", "", "", "", "'ALL'", "'APP'", "'ALIVE'", "'AND'", "'ANALYZE'", "'AS'",
		"'ASC'", "'BETWEEN'", "'BROKER'", "'BROKERS'", "'BY'", "'COMPACT'",
		"'CREATE'", "'CROSS'", "'COLUMNS'", "'DATABASE'", "'DATABASES'", "'DEFAULT'",
		"'DESC'", "'DISTRIBUTED'", "'DROP'", "'ENGINE'", "'ESCAPE'", "'EXPLAIN'",
		"'EXISTS'", "'FALSE'", "'FIELDS'", "'FLUSH'", "'FROM'", "'GROUP'", "'HAVING'",
		"'IF'", "'IN'", "'INSERT'", "'INTO'", "'INTERVAL'", "'IS'", "'JOIN'",
		"'KEYS'", "'LEFT'", "'LIKE'", "'LIMIT'", "'LOG'", "'LOGICAL'", "'MASTER'",
		"'MEMORY_DATABASES'", "'METRICS'", "'METRIC'", "'METADATA'", "'METADATAS'",
		"'NAMESPACE'", "'NAMESPACES'", "'NULL'", "'NOT'", "'NOW'", "'ON'", "'OR'",
		"'ORDER'", "'REQUESTS'", "'REPLICATIONS'", "'RIGHT'", "'ROLLUP'", "'SELECT'",
		"'SHOW'", "'STATE'", "'STORAGE'", "'TABLE_NAMES'", "'TIMESTAMP'", "'TRACE'",
		"'TRUE'", "'TYPE'", "'TYPES'", "'VALUES'", "'WHERE'", "'WITH'", "'WITHIN'",
		"'USING'", "'USE'", "'SECOND'", "'MINUTE'", "'HOUR'", "'DAY'", "'MONTH'",
		"'YEAR'", "'='", "", "'<'", "'<='", "'>'", "'>='", "'+'", "'-'", "'*'",
		"'/'", "'%'", "'=~'", "'!~'", "'!'", "'.'", "'('", "')'", "','", "'@'",
		"';'",
	}
	staticData.SymbolicNames = []string{
		"", "SIMPLE_COMMENT", "BRACKETED_COMMENT", "WS", "ALL", "APP", "ALIVE",
		"AND", "ANALYZE", "AS", "ASC", "BETWEEN", "BROKER", "BROKERS", "BY",
		"COMPACT", "CREATE", "CROSS", "COLUMNS", "DATABASE", "DATABASES", "DEFAULT",
		"DESC", "DISTRIBUTED", "DROP", "ENGINE", "ESCAPE", "EXPLAIN", "EXISTS",
		"FALSE", "FIELDS", "FLUSH", "FROM", "GROUP", "HAVING", "IF", "IN", "INSERT",
		"INTO", "INTERVAL", "IS", "JOIN", "KEYS", "LEFT", "LIKE", "LIMIT", "LOG",
		"LOGICAL", "MASTER", "MEMORY_DATABASES", "METRICS", "METRIC", "METADATA",
		"METADATAS", "NAMESPACE", "NAMESPACES", "NULL", "NOT", "NOW", "ON",
		"OR", "ORDER", "REQUESTS", "REPLICATIONS", "RIGHT", "ROLLUP", "SELECT",
		"SHOW", "STATE", "STORAGE", "TABLE_NAMES", "TIMESTAMP", "TRACE", "TRUE",
		"TYPE", "TYPES", "VALUES", "WHERE", "WITH", "WITHIN", "USING", "USE",
		"SECOND", "MINUTE", "HOUR", "DAY", "MONTH", "YEAR", "EQ", "NEQ", "LT",
		"LTE", "GT", "GTE", "PLUS", "MINUS", "ASTERISK", "SLASH", "PERCENT",
		"REGEXP", "NEQREGEXP", "EXCLAMATION_SYMBOL", "DOT", "LR_BRACKET", "RR_BRACKET",
		"COMMA", "AT", "SEMICOLON", "STRING", "INTEGER_VALUE", "DECIMAL_VALUE",
		"DOUBLE_VALUE", "IDENTIFIER", "DIGIT_IDENTIFIER", "QUOTED_IDENTIFIER",
		"BACKQUOTED_IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"statement", "streamingApp", "streamingQuery", "ddlStatement", "dmlStatement",
		"adminStatement", "utilityStatement", "explainOption", "createDatabase",
		"databaseOptions", "createDatabaseOptions", "rollupOptions", "dropDatabase",
		"createBroker", "flushDatabase", "compactDatabase", "showStatement",
		"useStatement", "query", "with", "namedQuery", "queryNoWith", "queryTerm",
		"queryPrimary", "querySpecification", "selectItem", "relation", "joinType",
		"joinCriteria", "aliasedRelation", "relationPrimary", "groupBy", "groupingElement",
		"groupingSet", "having", "orderBy", "sortItem", "limitRowCount", "expression",
		"booleanExpression", "valueExpression", "primaryExpression", "predicate",
		"comparisonOperator", "qualifiedName", "columnAliases", "appAnnotation",
		"annotation", "annotation_element", "properties", "propertyAssignments",
		"property", "propertyValue", "booleanValue", "string", "identifier",
		"number", "interval", "intervalUnit", "nonReserved",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 115, 710, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
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
		2, 58, 7, 58, 2, 59, 7, 59, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 3,
		0, 128, 8, 0, 1, 1, 5, 1, 131, 8, 1, 10, 1, 12, 1, 134, 9, 1, 1, 1, 5,
		1, 137, 8, 1, 10, 1, 12, 1, 140, 9, 1, 1, 2, 5, 2, 143, 8, 2, 10, 2, 12,
		2, 146, 9, 2, 1, 2, 1, 2, 3, 2, 150, 8, 2, 1, 3, 1, 3, 1, 3, 3, 3, 155,
		8, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 163, 8, 4, 10, 4, 12, 4,
		166, 9, 4, 1, 4, 1, 4, 3, 4, 170, 8, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1,
		4, 1, 4, 1, 4, 3, 4, 180, 8, 4, 1, 4, 1, 4, 3, 4, 184, 8, 4, 1, 5, 1, 5,
		1, 5, 3, 5, 189, 8, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8,
		1, 8, 5, 8, 200, 8, 8, 10, 8, 12, 8, 203, 9, 8, 1, 9, 1, 9, 1, 9, 5, 9,
		208, 8, 9, 10, 9, 12, 9, 211, 9, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 9, 5, 9, 220, 8, 9, 10, 9, 12, 9, 223, 9, 9, 1, 9, 1, 9, 3, 9, 227,
		8, 9, 1, 10, 1, 10, 3, 10, 231, 8, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 12,
		1, 12, 1, 12, 1, 12, 3, 12, 241, 8, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1,
		13, 1, 13, 1, 13, 3, 13, 250, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15,
		1, 15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1,
		16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16,
		3, 16, 278, 8, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 284, 8, 16, 1, 16,
		1, 16, 3, 16, 288, 8, 16, 1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 294, 8, 16,
		1, 17, 1, 17, 1, 17, 1, 18, 3, 18, 300, 8, 18, 1, 18, 1, 18, 1, 19, 1,
		19, 1, 19, 1, 19, 5, 19, 308, 8, 19, 10, 19, 12, 19, 311, 9, 19, 1, 20,
		1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 1, 21, 3, 21, 323,
		8, 21, 1, 21, 1, 21, 3, 21, 327, 8, 21, 1, 22, 1, 22, 1, 23, 1, 23, 1,
		23, 1, 23, 1, 23, 3, 23, 336, 8, 23, 1, 24, 1, 24, 1, 24, 1, 24, 5, 24,
		342, 8, 24, 10, 24, 12, 24, 345, 9, 24, 1, 24, 1, 24, 1, 24, 1, 24, 5,
		24, 351, 8, 24, 10, 24, 12, 24, 354, 9, 24, 3, 24, 356, 8, 24, 1, 24, 1,
		24, 3, 24, 360, 8, 24, 1, 24, 1, 24, 1, 24, 3, 24, 365, 8, 24, 1, 24, 1,
		24, 3, 24, 369, 8, 24, 1, 25, 1, 25, 3, 25, 373, 8, 25, 1, 25, 3, 25, 376,
		8, 25, 1, 25, 1, 25, 1, 25, 1, 25, 1, 25, 3, 25, 383, 8, 25, 1, 26, 1,
		26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26,
		3, 26, 397, 8, 26, 5, 26, 399, 8, 26, 10, 26, 12, 26, 402, 9, 26, 1, 27,
		1, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 5, 28, 413, 8,
		28, 10, 28, 12, 28, 416, 9, 28, 1, 28, 1, 28, 3, 28, 420, 8, 28, 1, 29,
		1, 29, 3, 29, 424, 8, 29, 1, 29, 3, 29, 427, 8, 29, 1, 30, 1, 30, 1, 30,
		1, 30, 1, 30, 3, 30, 434, 8, 30, 1, 31, 1, 31, 1, 31, 5, 31, 439, 8, 31,
		10, 31, 12, 31, 442, 9, 31, 1, 32, 1, 32, 3, 32, 446, 8, 32, 1, 33, 1,
		33, 1, 33, 1, 33, 5, 33, 452, 8, 33, 10, 33, 12, 33, 455, 9, 33, 3, 33,
		457, 8, 33, 1, 33, 1, 33, 3, 33, 461, 8, 33, 1, 34, 1, 34, 1, 35, 1, 35,
		1, 35, 5, 35, 468, 8, 35, 10, 35, 12, 35, 471, 9, 35, 1, 36, 1, 36, 3,
		36, 475, 8, 36, 1, 37, 1, 37, 1, 38, 1, 38, 1, 39, 1, 39, 1, 39, 1, 39,
		3, 39, 485, 8, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 5, 39, 493,
		8, 39, 10, 39, 12, 39, 496, 9, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1,
		40, 1, 40, 1, 40, 1, 40, 5, 40, 507, 8, 40, 10, 40, 12, 40, 510, 9, 40,
		1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1,
		41, 5, 41, 523, 8, 41, 10, 41, 12, 41, 526, 9, 41, 3, 41, 528, 8, 41, 1,
		41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 3, 41, 537, 8, 41, 1, 41,
		1, 41, 1, 41, 5, 41, 542, 8, 41, 10, 41, 12, 41, 545, 9, 41, 1, 42, 1,
		42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42,
		1, 42, 1, 42, 1, 42, 1, 42, 3, 42, 563, 8, 42, 1, 42, 1, 42, 1, 42, 1,
		42, 1, 42, 5, 42, 570, 8, 42, 10, 42, 12, 42, 573, 9, 42, 1, 42, 1, 42,
		1, 42, 1, 42, 3, 42, 579, 8, 42, 1, 42, 1, 42, 1, 42, 1, 42, 3, 42, 585,
		8, 42, 1, 42, 1, 42, 1, 42, 3, 42, 590, 8, 42, 1, 42, 1, 42, 1, 42, 3,
		42, 595, 8, 42, 1, 42, 1, 42, 1, 42, 3, 42, 600, 8, 42, 1, 43, 1, 43, 1,
		44, 1, 44, 1, 44, 5, 44, 607, 8, 44, 10, 44, 12, 44, 610, 9, 44, 1, 45,
		1, 45, 1, 45, 1, 45, 5, 45, 616, 8, 45, 10, 45, 12, 45, 619, 9, 45, 1,
		45, 1, 45, 1, 46, 1, 46, 1, 46, 1, 46, 1, 46, 1, 46, 5, 46, 629, 8, 46,
		10, 46, 12, 46, 632, 9, 46, 1, 46, 1, 46, 3, 46, 636, 8, 46, 1, 47, 1,
		47, 1, 47, 1, 47, 1, 47, 1, 47, 5, 47, 644, 8, 47, 10, 47, 12, 47, 647,
		9, 47, 1, 47, 1, 47, 3, 47, 651, 8, 47, 1, 48, 1, 48, 3, 48, 655, 8, 48,
		1, 49, 1, 49, 1, 49, 1, 49, 1, 50, 1, 50, 1, 50, 5, 50, 664, 8, 50, 10,
		50, 12, 50, 667, 9, 50, 1, 51, 1, 51, 1, 51, 1, 51, 1, 52, 1, 52, 3, 52,
		675, 8, 52, 1, 53, 1, 53, 1, 54, 1, 54, 1, 55, 1, 55, 1, 55, 1, 55, 1,
		55, 3, 55, 686, 8, 55, 1, 56, 3, 56, 689, 8, 56, 1, 56, 1, 56, 3, 56, 693,
		8, 56, 1, 56, 1, 56, 3, 56, 697, 8, 56, 1, 56, 3, 56, 700, 8, 56, 1, 57,
		1, 57, 1, 57, 1, 57, 1, 58, 1, 58, 1, 59, 1, 59, 1, 59, 0, 4, 52, 78, 80,
		82, 60, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
		34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68,
		70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104,
		106, 108, 110, 112, 114, 116, 118, 0, 12, 2, 0, 23, 23, 47, 47, 3, 0, 46,
		46, 51, 51, 72, 72, 2, 0, 43, 43, 64, 64, 2, 0, 10, 10, 22, 22, 2, 0, 57,
		57, 101, 101, 1, 0, 96, 98, 1, 0, 94, 95, 1, 0, 99, 100, 1, 0, 88, 93,
		2, 0, 29, 29, 73, 73, 1, 0, 82, 87, 2, 0, 4, 7, 9, 87, 762, 0, 127, 1,
		0, 0, 0, 2, 132, 1, 0, 0, 0, 4, 144, 1, 0, 0, 0, 6, 154, 1, 0, 0, 0, 8,
		183, 1, 0, 0, 0, 10, 188, 1, 0, 0, 0, 12, 190, 1, 0, 0, 0, 14, 192, 1,
		0, 0, 0, 16, 195, 1, 0, 0, 0, 18, 226, 1, 0, 0, 0, 20, 228, 1, 0, 0, 0,
		22, 234, 1, 0, 0, 0, 24, 236, 1, 0, 0, 0, 26, 244, 1, 0, 0, 0, 28, 251,
		1, 0, 0, 0, 30, 255, 1, 0, 0, 0, 32, 293, 1, 0, 0, 0, 34, 295, 1, 0, 0,
		0, 36, 299, 1, 0, 0, 0, 38, 303, 1, 0, 0, 0, 40, 312, 1, 0, 0, 0, 42, 318,
		1, 0, 0, 0, 44, 328, 1, 0, 0, 0, 46, 335, 1, 0, 0, 0, 48, 337, 1, 0, 0,
		0, 50, 382, 1, 0, 0, 0, 52, 384, 1, 0, 0, 0, 54, 403, 1, 0, 0, 0, 56, 419,
		1, 0, 0, 0, 58, 421, 1, 0, 0, 0, 60, 433, 1, 0, 0, 0, 62, 435, 1, 0, 0,
		0, 64, 445, 1, 0, 0, 0, 66, 460, 1, 0, 0, 0, 68, 462, 1, 0, 0, 0, 70, 464,
		1, 0, 0, 0, 72, 472, 1, 0, 0, 0, 74, 476, 1, 0, 0, 0, 76, 478, 1, 0, 0,
		0, 78, 484, 1, 0, 0, 0, 80, 497, 1, 0, 0, 0, 82, 536, 1, 0, 0, 0, 84, 599,
		1, 0, 0, 0, 86, 601, 1, 0, 0, 0, 88, 603, 1, 0, 0, 0, 90, 611, 1, 0, 0,
		0, 92, 622, 1, 0, 0, 0, 94, 637, 1, 0, 0, 0, 96, 654, 1, 0, 0, 0, 98, 656,
		1, 0, 0, 0, 100, 660, 1, 0, 0, 0, 102, 668, 1, 0, 0, 0, 104, 674, 1, 0,
		0, 0, 106, 676, 1, 0, 0, 0, 108, 678, 1, 0, 0, 0, 110, 685, 1, 0, 0, 0,
		112, 699, 1, 0, 0, 0, 114, 701, 1, 0, 0, 0, 116, 705, 1, 0, 0, 0, 118,
		707, 1, 0, 0, 0, 120, 128, 3, 6, 3, 0, 121, 128, 3, 8, 4, 0, 122, 128,
		3, 10, 5, 0, 123, 128, 3, 12, 6, 0, 124, 125, 3, 2, 1, 0, 125, 126, 5,
		0, 0, 1, 126, 128, 1, 0, 0, 0, 127, 120, 1, 0, 0, 0, 127, 121, 1, 0, 0,
		0, 127, 122, 1, 0, 0, 0, 127, 123, 1, 0, 0, 0, 127, 124, 1, 0, 0, 0, 128,
		1, 1, 0, 0, 0, 129, 131, 3, 92, 46, 0, 130, 129, 1, 0, 0, 0, 131, 134,
		1, 0, 0, 0, 132, 130, 1, 0, 0, 0, 132, 133, 1, 0, 0, 0, 133, 138, 1, 0,
		0, 0, 134, 132, 1, 0, 0, 0, 135, 137, 3, 4, 2, 0, 136, 135, 1, 0, 0, 0,
		137, 140, 1, 0, 0, 0, 138, 136, 1, 0, 0, 0, 138, 139, 1, 0, 0, 0, 139,
		3, 1, 0, 0, 0, 140, 138, 1, 0, 0, 0, 141, 143, 3, 94, 47, 0, 142, 141,
		1, 0, 0, 0, 143, 146, 1, 0, 0, 0, 144, 142, 1, 0, 0, 0, 144, 145, 1, 0,
		0, 0, 145, 147, 1, 0, 0, 0, 146, 144, 1, 0, 0, 0, 147, 149, 3, 8, 4, 0,
		148, 150, 5, 107, 0, 0, 149, 148, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150,
		5, 1, 0, 0, 0, 151, 155, 3, 16, 8, 0, 152, 155, 3, 24, 12, 0, 153, 155,
		3, 26, 13, 0, 154, 151, 1, 0, 0, 0, 154, 152, 1, 0, 0, 0, 154, 153, 1,
		0, 0, 0, 155, 7, 1, 0, 0, 0, 156, 184, 3, 36, 18, 0, 157, 169, 5, 27, 0,
		0, 158, 159, 5, 103, 0, 0, 159, 164, 3, 14, 7, 0, 160, 161, 5, 105, 0,
		0, 161, 163, 3, 14, 7, 0, 162, 160, 1, 0, 0, 0, 163, 166, 1, 0, 0, 0, 164,
		162, 1, 0, 0, 0, 164, 165, 1, 0, 0, 0, 165, 167, 1, 0, 0, 0, 166, 164,
		1, 0, 0, 0, 167, 168, 5, 104, 0, 0, 168, 170, 1, 0, 0, 0, 169, 158, 1,
		0, 0, 0, 169, 170, 1, 0, 0, 0, 170, 171, 1, 0, 0, 0, 171, 184, 3, 8, 4,
		0, 172, 173, 5, 27, 0, 0, 173, 174, 5, 8, 0, 0, 174, 184, 3, 8, 4, 0, 175,
		176, 5, 37, 0, 0, 176, 177, 5, 38, 0, 0, 177, 179, 3, 88, 44, 0, 178, 180,
		3, 90, 45, 0, 179, 178, 1, 0, 0, 0, 179, 180, 1, 0, 0, 0, 180, 181, 1,
		0, 0, 0, 181, 182, 3, 36, 18, 0, 182, 184, 1, 0, 0, 0, 183, 156, 1, 0,
		0, 0, 183, 157, 1, 0, 0, 0, 183, 172, 1, 0, 0, 0, 183, 175, 1, 0, 0, 0,
		184, 9, 1, 0, 0, 0, 185, 189, 3, 28, 14, 0, 186, 189, 3, 30, 15, 0, 187,
		189, 3, 32, 16, 0, 188, 185, 1, 0, 0, 0, 188, 186, 1, 0, 0, 0, 188, 187,
		1, 0, 0, 0, 189, 11, 1, 0, 0, 0, 190, 191, 3, 34, 17, 0, 191, 13, 1, 0,
		0, 0, 192, 193, 5, 74, 0, 0, 193, 194, 7, 0, 0, 0, 194, 15, 1, 0, 0, 0,
		195, 196, 5, 16, 0, 0, 196, 197, 5, 19, 0, 0, 197, 201, 3, 88, 44, 0, 198,
		200, 3, 18, 9, 0, 199, 198, 1, 0, 0, 0, 200, 203, 1, 0, 0, 0, 201, 199,
		1, 0, 0, 0, 201, 202, 1, 0, 0, 0, 202, 17, 1, 0, 0, 0, 203, 201, 1, 0,
		0, 0, 204, 209, 3, 20, 10, 0, 205, 206, 5, 105, 0, 0, 206, 208, 3, 20,
		10, 0, 207, 205, 1, 0, 0, 0, 208, 211, 1, 0, 0, 0, 209, 207, 1, 0, 0, 0,
		209, 210, 1, 0, 0, 0, 210, 227, 1, 0, 0, 0, 211, 209, 1, 0, 0, 0, 212,
		213, 5, 78, 0, 0, 213, 227, 3, 98, 49, 0, 214, 215, 5, 65, 0, 0, 215, 216,
		5, 103, 0, 0, 216, 221, 3, 22, 11, 0, 217, 218, 5, 105, 0, 0, 218, 220,
		3, 22, 11, 0, 219, 217, 1, 0, 0, 0, 220, 223, 1, 0, 0, 0, 221, 219, 1,
		0, 0, 0, 221, 222, 1, 0, 0, 0, 222, 224, 1, 0, 0, 0, 223, 221, 1, 0, 0,
		0, 224, 225, 5, 104, 0, 0, 225, 227, 1, 0, 0, 0, 226, 204, 1, 0, 0, 0,
		226, 212, 1, 0, 0, 0, 226, 214, 1, 0, 0, 0, 227, 19, 1, 0, 0, 0, 228, 230,
		5, 25, 0, 0, 229, 231, 5, 88, 0, 0, 230, 229, 1, 0, 0, 0, 230, 231, 1,
		0, 0, 0, 231, 232, 1, 0, 0, 0, 232, 233, 7, 1, 0, 0, 233, 21, 1, 0, 0,
		0, 234, 235, 3, 98, 49, 0, 235, 23, 1, 0, 0, 0, 236, 237, 5, 24, 0, 0,
		237, 240, 5, 19, 0, 0, 238, 239, 5, 35, 0, 0, 239, 241, 5, 28, 0, 0, 240,
		238, 1, 0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 242, 1, 0, 0, 0, 242, 243,
		3, 88, 44, 0, 243, 25, 1, 0, 0, 0, 244, 245, 5, 16, 0, 0, 245, 246, 5,
		12, 0, 0, 246, 249, 3, 88, 44, 0, 247, 248, 5, 78, 0, 0, 248, 250, 3, 98,
		49, 0, 249, 247, 1, 0, 0, 0, 249, 250, 1, 0, 0, 0, 250, 27, 1, 0, 0, 0,
		251, 252, 5, 31, 0, 0, 252, 253, 5, 19, 0, 0, 253, 254, 3, 88, 44, 0, 254,
		29, 1, 0, 0, 0, 255, 256, 5, 15, 0, 0, 256, 257, 5, 19, 0, 0, 257, 258,
		3, 88, 44, 0, 258, 31, 1, 0, 0, 0, 259, 260, 5, 67, 0, 0, 260, 294, 5,
		48, 0, 0, 261, 262, 5, 67, 0, 0, 262, 294, 5, 13, 0, 0, 263, 264, 5, 67,
		0, 0, 264, 294, 5, 20, 0, 0, 265, 266, 5, 67, 0, 0, 266, 294, 5, 62, 0,
		0, 267, 268, 5, 67, 0, 0, 268, 294, 5, 45, 0, 0, 269, 270, 5, 67, 0, 0,
		270, 294, 5, 49, 0, 0, 271, 272, 5, 67, 0, 0, 272, 294, 5, 63, 0, 0, 273,
		274, 5, 67, 0, 0, 274, 277, 5, 55, 0, 0, 275, 276, 5, 44, 0, 0, 276, 278,
		3, 76, 38, 0, 277, 275, 1, 0, 0, 0, 277, 278, 1, 0, 0, 0, 278, 294, 1,
		0, 0, 0, 279, 280, 5, 67, 0, 0, 280, 283, 5, 70, 0, 0, 281, 282, 5, 32,
		0, 0, 282, 284, 3, 88, 44, 0, 283, 281, 1, 0, 0, 0, 283, 284, 1, 0, 0,
		0, 284, 287, 1, 0, 0, 0, 285, 286, 5, 44, 0, 0, 286, 288, 3, 76, 38, 0,
		287, 285, 1, 0, 0, 0, 287, 288, 1, 0, 0, 0, 288, 294, 1, 0, 0, 0, 289,
		290, 5, 67, 0, 0, 290, 291, 5, 18, 0, 0, 291, 292, 5, 32, 0, 0, 292, 294,
		3, 88, 44, 0, 293, 259, 1, 0, 0, 0, 293, 261, 1, 0, 0, 0, 293, 263, 1,
		0, 0, 0, 293, 265, 1, 0, 0, 0, 293, 267, 1, 0, 0, 0, 293, 269, 1, 0, 0,
		0, 293, 271, 1, 0, 0, 0, 293, 273, 1, 0, 0, 0, 293, 279, 1, 0, 0, 0, 293,
		289, 1, 0, 0, 0, 294, 33, 1, 0, 0, 0, 295, 296, 5, 81, 0, 0, 296, 297,
		3, 110, 55, 0, 297, 35, 1, 0, 0, 0, 298, 300, 3, 38, 19, 0, 299, 298, 1,
		0, 0, 0, 299, 300, 1, 0, 0, 0, 300, 301, 1, 0, 0, 0, 301, 302, 3, 42, 21,
		0, 302, 37, 1, 0, 0, 0, 303, 304, 5, 78, 0, 0, 304, 309, 3, 40, 20, 0,
		305, 306, 5, 105, 0, 0, 306, 308, 3, 40, 20, 0, 307, 305, 1, 0, 0, 0, 308,
		311, 1, 0, 0, 0, 309, 307, 1, 0, 0, 0, 309, 310, 1, 0, 0, 0, 310, 39, 1,
		0, 0, 0, 311, 309, 1, 0, 0, 0, 312, 313, 3, 110, 55, 0, 313, 314, 5, 9,
		0, 0, 314, 315, 5, 103, 0, 0, 315, 316, 3, 36, 18, 0, 316, 317, 5, 104,
		0, 0, 317, 41, 1, 0, 0, 0, 318, 322, 3, 44, 22, 0, 319, 320, 5, 61, 0,
		0, 320, 321, 5, 14, 0, 0, 321, 323, 3, 70, 35, 0, 322, 319, 1, 0, 0, 0,
		322, 323, 1, 0, 0, 0, 323, 326, 1, 0, 0, 0, 324, 325, 5, 45, 0, 0, 325,
		327, 3, 74, 37, 0, 326, 324, 1, 0, 0, 0, 326, 327, 1, 0, 0, 0, 327, 43,
		1, 0, 0, 0, 328, 329, 3, 46, 23, 0, 329, 45, 1, 0, 0, 0, 330, 336, 3, 48,
		24, 0, 331, 332, 5, 103, 0, 0, 332, 333, 3, 42, 21, 0, 333, 334, 5, 104,
		0, 0, 334, 336, 1, 0, 0, 0, 335, 330, 1, 0, 0, 0, 335, 331, 1, 0, 0, 0,
		336, 47, 1, 0, 0, 0, 337, 338, 5, 66, 0, 0, 338, 343, 3, 50, 25, 0, 339,
		340, 5, 105, 0, 0, 340, 342, 3, 50, 25, 0, 341, 339, 1, 0, 0, 0, 342, 345,
		1, 0, 0, 0, 343, 341, 1, 0, 0, 0, 343, 344, 1, 0, 0, 0, 344, 355, 1, 0,
		0, 0, 345, 343, 1, 0, 0, 0, 346, 347, 5, 32, 0, 0, 347, 352, 3, 52, 26,
		0, 348, 349, 5, 105, 0, 0, 349, 351, 3, 52, 26, 0, 350, 348, 1, 0, 0, 0,
		351, 354, 1, 0, 0, 0, 352, 350, 1, 0, 0, 0, 352, 353, 1, 0, 0, 0, 353,
		356, 1, 0, 0, 0, 354, 352, 1, 0, 0, 0, 355, 346, 1, 0, 0, 0, 355, 356,
		1, 0, 0, 0, 356, 359, 1, 0, 0, 0, 357, 358, 5, 77, 0, 0, 358, 360, 3, 78,
		39, 0, 359, 357, 1, 0, 0, 0, 359, 360, 1, 0, 0, 0, 360, 364, 1, 0, 0, 0,
		361, 362, 5, 33, 0, 0, 362, 363, 5, 14, 0, 0, 363, 365, 3, 62, 31, 0, 364,
		361, 1, 0, 0, 0, 364, 365, 1, 0, 0, 0, 365, 368, 1, 0, 0, 0, 366, 367,
		5, 34, 0, 0, 367, 369, 3, 68, 34, 0, 368, 366, 1, 0, 0, 0, 368, 369, 1,
		0, 0, 0, 369, 49, 1, 0, 0, 0, 370, 375, 3, 76, 38, 0, 371, 373, 5, 9, 0,
		0, 372, 371, 1, 0, 0, 0, 372, 373, 1, 0, 0, 0, 373, 374, 1, 0, 0, 0, 374,
		376, 3, 110, 55, 0, 375, 372, 1, 0, 0, 0, 375, 376, 1, 0, 0, 0, 376, 383,
		1, 0, 0, 0, 377, 378, 3, 82, 41, 0, 378, 379, 5, 102, 0, 0, 379, 380, 5,
		96, 0, 0, 380, 383, 1, 0, 0, 0, 381, 383, 5, 96, 0, 0, 382, 370, 1, 0,
		0, 0, 382, 377, 1, 0, 0, 0, 382, 381, 1, 0, 0, 0, 383, 51, 1, 0, 0, 0,
		384, 385, 6, 26, -1, 0, 385, 386, 3, 58, 29, 0, 386, 400, 1, 0, 0, 0, 387,
		396, 10, 2, 0, 0, 388, 389, 5, 17, 0, 0, 389, 390, 5, 41, 0, 0, 390, 397,
		3, 52, 26, 0, 391, 392, 3, 54, 27, 0, 392, 393, 5, 41, 0, 0, 393, 394,
		3, 52, 26, 0, 394, 395, 3, 56, 28, 0, 395, 397, 1, 0, 0, 0, 396, 388, 1,
		0, 0, 0, 396, 391, 1, 0, 0, 0, 397, 399, 1, 0, 0, 0, 398, 387, 1, 0, 0,
		0, 399, 402, 1, 0, 0, 0, 400, 398, 1, 0, 0, 0, 400, 401, 1, 0, 0, 0, 401,
		53, 1, 0, 0, 0, 402, 400, 1, 0, 0, 0, 403, 404, 7, 2, 0, 0, 404, 55, 1,
		0, 0, 0, 405, 406, 5, 59, 0, 0, 406, 420, 3, 78, 39, 0, 407, 408, 5, 80,
		0, 0, 408, 409, 5, 103, 0, 0, 409, 414, 3, 110, 55, 0, 410, 411, 5, 105,
		0, 0, 411, 413, 3, 110, 55, 0, 412, 410, 1, 0, 0, 0, 413, 416, 1, 0, 0,
		0, 414, 412, 1, 0, 0, 0, 414, 415, 1, 0, 0, 0, 415, 417, 1, 0, 0, 0, 416,
		414, 1, 0, 0, 0, 417, 418, 5, 104, 0, 0, 418, 420, 1, 0, 0, 0, 419, 405,
		1, 0, 0, 0, 419, 407, 1, 0, 0, 0, 420, 57, 1, 0, 0, 0, 421, 426, 3, 60,
		30, 0, 422, 424, 5, 9, 0, 0, 423, 422, 1, 0, 0, 0, 423, 424, 1, 0, 0, 0,
		424, 425, 1, 0, 0, 0, 425, 427, 3, 110, 55, 0, 426, 423, 1, 0, 0, 0, 426,
		427, 1, 0, 0, 0, 427, 59, 1, 0, 0, 0, 428, 434, 3, 88, 44, 0, 429, 430,
		5, 103, 0, 0, 430, 431, 3, 36, 18, 0, 431, 432, 5, 104, 0, 0, 432, 434,
		1, 0, 0, 0, 433, 428, 1, 0, 0, 0, 433, 429, 1, 0, 0, 0, 434, 61, 1, 0,
		0, 0, 435, 440, 3, 64, 32, 0, 436, 437, 5, 105, 0, 0, 437, 439, 3, 64,
		32, 0, 438, 436, 1, 0, 0, 0, 439, 442, 1, 0, 0, 0, 440, 438, 1, 0, 0, 0,
		440, 441, 1, 0, 0, 0, 441, 63, 1, 0, 0, 0, 442, 440, 1, 0, 0, 0, 443, 446,
		3, 66, 33, 0, 444, 446, 5, 4, 0, 0, 445, 443, 1, 0, 0, 0, 445, 444, 1,
		0, 0, 0, 446, 65, 1, 0, 0, 0, 447, 456, 5, 103, 0, 0, 448, 453, 3, 76,
		38, 0, 449, 450, 5, 105, 0, 0, 450, 452, 3, 76, 38, 0, 451, 449, 1, 0,
		0, 0, 452, 455, 1, 0, 0, 0, 453, 451, 1, 0, 0, 0, 453, 454, 1, 0, 0, 0,
		454, 457, 1, 0, 0, 0, 455, 453, 1, 0, 0, 0, 456, 448, 1, 0, 0, 0, 456,
		457, 1, 0, 0, 0, 457, 458, 1, 0, 0, 0, 458, 461, 5, 104, 0, 0, 459, 461,
		3, 76, 38, 0, 460, 447, 1, 0, 0, 0, 460, 459, 1, 0, 0, 0, 461, 67, 1, 0,
		0, 0, 462, 463, 3, 78, 39, 0, 463, 69, 1, 0, 0, 0, 464, 469, 3, 72, 36,
		0, 465, 466, 5, 105, 0, 0, 466, 468, 3, 72, 36, 0, 467, 465, 1, 0, 0, 0,
		468, 471, 1, 0, 0, 0, 469, 467, 1, 0, 0, 0, 469, 470, 1, 0, 0, 0, 470,
		71, 1, 0, 0, 0, 471, 469, 1, 0, 0, 0, 472, 474, 3, 76, 38, 0, 473, 475,
		7, 3, 0, 0, 474, 473, 1, 0, 0, 0, 474, 475, 1, 0, 0, 0, 475, 73, 1, 0,
		0, 0, 476, 477, 5, 109, 0, 0, 477, 75, 1, 0, 0, 0, 478, 479, 3, 78, 39,
		0, 479, 77, 1, 0, 0, 0, 480, 481, 6, 39, -1, 0, 481, 482, 7, 4, 0, 0, 482,
		485, 3, 78, 39, 4, 483, 485, 3, 84, 42, 0, 484, 480, 1, 0, 0, 0, 484, 483,
		1, 0, 0, 0, 485, 494, 1, 0, 0, 0, 486, 487, 10, 3, 0, 0, 487, 488, 5, 7,
		0, 0, 488, 493, 3, 78, 39, 4, 489, 490, 10, 2, 0, 0, 490, 491, 5, 60, 0,
		0, 491, 493, 3, 78, 39, 3, 492, 486, 1, 0, 0, 0, 492, 489, 1, 0, 0, 0,
		493, 496, 1, 0, 0, 0, 494, 492, 1, 0, 0, 0, 494, 495, 1, 0, 0, 0, 495,
		79, 1, 0, 0, 0, 496, 494, 1, 0, 0, 0, 497, 498, 6, 40, -1, 0, 498, 499,
		3, 82, 41, 0, 499, 508, 1, 0, 0, 0, 500, 501, 10, 2, 0, 0, 501, 502, 7,
		5, 0, 0, 502, 507, 3, 80, 40, 3, 503, 504, 10, 1, 0, 0, 504, 505, 7, 6,
		0, 0, 505, 507, 3, 80, 40, 2, 506, 500, 1, 0, 0, 0, 506, 503, 1, 0, 0,
		0, 507, 510, 1, 0, 0, 0, 508, 506, 1, 0, 0, 0, 508, 509, 1, 0, 0, 0, 509,
		81, 1, 0, 0, 0, 510, 508, 1, 0, 0, 0, 511, 512, 6, 41, -1, 0, 512, 537,
		3, 112, 56, 0, 513, 537, 3, 114, 57, 0, 514, 537, 3, 106, 53, 0, 515, 537,
		5, 56, 0, 0, 516, 537, 3, 108, 54, 0, 517, 518, 3, 88, 44, 0, 518, 527,
		5, 103, 0, 0, 519, 524, 3, 76, 38, 0, 520, 521, 5, 105, 0, 0, 521, 523,
		3, 76, 38, 0, 522, 520, 1, 0, 0, 0, 523, 526, 1, 0, 0, 0, 524, 522, 1,
		0, 0, 0, 524, 525, 1, 0, 0, 0, 525, 528, 1, 0, 0, 0, 526, 524, 1, 0, 0,
		0, 527, 519, 1, 0, 0, 0, 527, 528, 1, 0, 0, 0, 528, 529, 1, 0, 0, 0, 529,
		530, 5, 104, 0, 0, 530, 537, 1, 0, 0, 0, 531, 537, 3, 110, 55, 0, 532,
		533, 5, 103, 0, 0, 533, 534, 3, 76, 38, 0, 534, 535, 5, 104, 0, 0, 535,
		537, 1, 0, 0, 0, 536, 511, 1, 0, 0, 0, 536, 513, 1, 0, 0, 0, 536, 514,
		1, 0, 0, 0, 536, 515, 1, 0, 0, 0, 536, 516, 1, 0, 0, 0, 536, 517, 1, 0,
		0, 0, 536, 531, 1, 0, 0, 0, 536, 532, 1, 0, 0, 0, 537, 543, 1, 0, 0, 0,
		538, 539, 10, 2, 0, 0, 539, 540, 5, 102, 0, 0, 540, 542, 3, 110, 55, 0,
		541, 538, 1, 0, 0, 0, 542, 545, 1, 0, 0, 0, 543, 541, 1, 0, 0, 0, 543,
		544, 1, 0, 0, 0, 544, 83, 1, 0, 0, 0, 545, 543, 1, 0, 0, 0, 546, 547, 5,
		71, 0, 0, 547, 548, 3, 86, 43, 0, 548, 549, 3, 80, 40, 0, 549, 600, 1,
		0, 0, 0, 550, 551, 3, 80, 40, 0, 551, 552, 3, 86, 43, 0, 552, 553, 3, 80,
		40, 0, 553, 600, 1, 0, 0, 0, 554, 555, 5, 71, 0, 0, 555, 556, 5, 11, 0,
		0, 556, 557, 3, 80, 40, 0, 557, 558, 5, 7, 0, 0, 558, 559, 3, 80, 40, 0,
		559, 600, 1, 0, 0, 0, 560, 562, 3, 80, 40, 0, 561, 563, 5, 57, 0, 0, 562,
		561, 1, 0, 0, 0, 562, 563, 1, 0, 0, 0, 563, 564, 1, 0, 0, 0, 564, 565,
		5, 36, 0, 0, 565, 566, 5, 103, 0, 0, 566, 571, 3, 76, 38, 0, 567, 568,
		5, 105, 0, 0, 568, 570, 3, 76, 38, 0, 569, 567, 1, 0, 0, 0, 570, 573, 1,
		0, 0, 0, 571, 569, 1, 0, 0, 0, 571, 572, 1, 0, 0, 0, 572, 574, 1, 0, 0,
		0, 573, 571, 1, 0, 0, 0, 574, 575, 5, 104, 0, 0, 575, 600, 1, 0, 0, 0,
		576, 578, 3, 80, 40, 0, 577, 579, 5, 57, 0, 0, 578, 577, 1, 0, 0, 0, 578,
		579, 1, 0, 0, 0, 579, 580, 1, 0, 0, 0, 580, 581, 5, 44, 0, 0, 581, 584,
		3, 80, 40, 0, 582, 583, 5, 26, 0, 0, 583, 585, 3, 80, 40, 0, 584, 582,
		1, 0, 0, 0, 584, 585, 1, 0, 0, 0, 585, 600, 1, 0, 0, 0, 586, 587, 3, 80,
		40, 0, 587, 589, 7, 7, 0, 0, 588, 590, 3, 80, 40, 0, 589, 588, 1, 0, 0,
		0, 589, 590, 1, 0, 0, 0, 590, 600, 1, 0, 0, 0, 591, 592, 3, 80, 40, 0,
		592, 594, 5, 40, 0, 0, 593, 595, 5, 57, 0, 0, 594, 593, 1, 0, 0, 0, 594,
		595, 1, 0, 0, 0, 595, 596, 1, 0, 0, 0, 596, 597, 5, 56, 0, 0, 597, 600,
		1, 0, 0, 0, 598, 600, 3, 80, 40, 0, 599, 546, 1, 0, 0, 0, 599, 550, 1,
		0, 0, 0, 599, 554, 1, 0, 0, 0, 599, 560, 1, 0, 0, 0, 599, 576, 1, 0, 0,
		0, 599, 586, 1, 0, 0, 0, 599, 591, 1, 0, 0, 0, 599, 598, 1, 0, 0, 0, 600,
		85, 1, 0, 0, 0, 601, 602, 7, 8, 0, 0, 602, 87, 1, 0, 0, 0, 603, 608, 3,
		110, 55, 0, 604, 605, 5, 102, 0, 0, 605, 607, 3, 110, 55, 0, 606, 604,
		1, 0, 0, 0, 607, 610, 1, 0, 0, 0, 608, 606, 1, 0, 0, 0, 608, 609, 1, 0,
		0, 0, 609, 89, 1, 0, 0, 0, 610, 608, 1, 0, 0, 0, 611, 612, 5, 103, 0, 0,
		612, 617, 3, 110, 55, 0, 613, 614, 5, 105, 0, 0, 614, 616, 3, 110, 55,
		0, 615, 613, 1, 0, 0, 0, 616, 619, 1, 0, 0, 0, 617, 615, 1, 0, 0, 0, 617,
		618, 1, 0, 0, 0, 618, 620, 1, 0, 0, 0, 619, 617, 1, 0, 0, 0, 620, 621,
		5, 104, 0, 0, 621, 91, 1, 0, 0, 0, 622, 623, 5, 106, 0, 0, 623, 635, 5,
		5, 0, 0, 624, 625, 5, 103, 0, 0, 625, 630, 3, 96, 48, 0, 626, 627, 5, 105,
		0, 0, 627, 629, 3, 96, 48, 0, 628, 626, 1, 0, 0, 0, 629, 632, 1, 0, 0,
		0, 630, 628, 1, 0, 0, 0, 630, 631, 1, 0, 0, 0, 631, 633, 1, 0, 0, 0, 632,
		630, 1, 0, 0, 0, 633, 634, 5, 104, 0, 0, 634, 636, 1, 0, 0, 0, 635, 624,
		1, 0, 0, 0, 635, 636, 1, 0, 0, 0, 636, 93, 1, 0, 0, 0, 637, 638, 5, 106,
		0, 0, 638, 650, 3, 110, 55, 0, 639, 640, 5, 103, 0, 0, 640, 645, 3, 96,
		48, 0, 641, 642, 5, 105, 0, 0, 642, 644, 3, 96, 48, 0, 643, 641, 1, 0,
		0, 0, 644, 647, 1, 0, 0, 0, 645, 643, 1, 0, 0, 0, 645, 646, 1, 0, 0, 0,
		646, 648, 1, 0, 0, 0, 647, 645, 1, 0, 0, 0, 648, 649, 5, 104, 0, 0, 649,
		651, 1, 0, 0, 0, 650, 639, 1, 0, 0, 0, 650, 651, 1, 0, 0, 0, 651, 95, 1,
		0, 0, 0, 652, 655, 3, 102, 51, 0, 653, 655, 3, 94, 47, 0, 654, 652, 1,
		0, 0, 0, 654, 653, 1, 0, 0, 0, 655, 97, 1, 0, 0, 0, 656, 657, 5, 103, 0,
		0, 657, 658, 3, 100, 50, 0, 658, 659, 5, 104, 0, 0, 659, 99, 1, 0, 0, 0,
		660, 665, 3, 102, 51, 0, 661, 662, 5, 105, 0, 0, 662, 664, 3, 102, 51,
		0, 663, 661, 1, 0, 0, 0, 664, 667, 1, 0, 0, 0, 665, 663, 1, 0, 0, 0, 665,
		666, 1, 0, 0, 0, 666, 101, 1, 0, 0, 0, 667, 665, 1, 0, 0, 0, 668, 669,
		3, 110, 55, 0, 669, 670, 5, 88, 0, 0, 670, 671, 3, 104, 52, 0, 671, 103,
		1, 0, 0, 0, 672, 675, 5, 21, 0, 0, 673, 675, 3, 76, 38, 0, 674, 672, 1,
		0, 0, 0, 674, 673, 1, 0, 0, 0, 675, 105, 1, 0, 0, 0, 676, 677, 7, 9, 0,
		0, 677, 107, 1, 0, 0, 0, 678, 679, 5, 108, 0, 0, 679, 109, 1, 0, 0, 0,
		680, 686, 5, 112, 0, 0, 681, 686, 5, 114, 0, 0, 682, 686, 3, 118, 59, 0,
		683, 686, 5, 115, 0, 0, 684, 686, 5, 113, 0, 0, 685, 680, 1, 0, 0, 0, 685,
		681, 1, 0, 0, 0, 685, 682, 1, 0, 0, 0, 685, 683, 1, 0, 0, 0, 685, 684,
		1, 0, 0, 0, 686, 111, 1, 0, 0, 0, 687, 689, 5, 95, 0, 0, 688, 687, 1, 0,
		0, 0, 688, 689, 1, 0, 0, 0, 689, 690, 1, 0, 0, 0, 690, 700, 5, 110, 0,
		0, 691, 693, 5, 95, 0, 0, 692, 691, 1, 0, 0, 0, 692, 693, 1, 0, 0, 0, 693,
		694, 1, 0, 0, 0, 694, 700, 5, 111, 0, 0, 695, 697, 5, 95, 0, 0, 696, 695,
		1, 0, 0, 0, 696, 697, 1, 0, 0, 0, 697, 698, 1, 0, 0, 0, 698, 700, 5, 109,
		0, 0, 699, 688, 1, 0, 0, 0, 699, 692, 1, 0, 0, 0, 699, 696, 1, 0, 0, 0,
		700, 113, 1, 0, 0, 0, 701, 702, 5, 39, 0, 0, 702, 703, 3, 112, 56, 0, 703,
		704, 3, 116, 58, 0, 704, 115, 1, 0, 0, 0, 705, 706, 7, 10, 0, 0, 706, 117,
		1, 0, 0, 0, 707, 708, 7, 11, 0, 0, 708, 119, 1, 0, 0, 0, 80, 127, 132,
		138, 144, 149, 154, 164, 169, 179, 183, 188, 201, 209, 221, 226, 230, 240,
		249, 277, 283, 287, 293, 299, 309, 322, 326, 335, 343, 352, 355, 359, 364,
		368, 372, 375, 382, 396, 400, 414, 419, 423, 426, 433, 440, 445, 453, 456,
		460, 469, 474, 484, 492, 494, 506, 508, 524, 527, 536, 543, 562, 571, 578,
		584, 589, 594, 599, 608, 617, 630, 635, 645, 650, 654, 665, 674, 685, 688,
		692, 696, 699,
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
	SQLParserAPP                   = 5
	SQLParserALIVE                 = 6
	SQLParserAND                   = 7
	SQLParserANALYZE               = 8
	SQLParserAS                    = 9
	SQLParserASC                   = 10
	SQLParserBETWEEN               = 11
	SQLParserBROKER                = 12
	SQLParserBROKERS               = 13
	SQLParserBY                    = 14
	SQLParserCOMPACT               = 15
	SQLParserCREATE                = 16
	SQLParserCROSS                 = 17
	SQLParserCOLUMNS               = 18
	SQLParserDATABASE              = 19
	SQLParserDATABASES             = 20
	SQLParserDEFAULT               = 21
	SQLParserDESC                  = 22
	SQLParserDISTRIBUTED           = 23
	SQLParserDROP                  = 24
	SQLParserENGINE                = 25
	SQLParserESCAPE                = 26
	SQLParserEXPLAIN               = 27
	SQLParserEXISTS                = 28
	SQLParserFALSE                 = 29
	SQLParserFIELDS                = 30
	SQLParserFLUSH                 = 31
	SQLParserFROM                  = 32
	SQLParserGROUP                 = 33
	SQLParserHAVING                = 34
	SQLParserIF                    = 35
	SQLParserIN                    = 36
	SQLParserINSERT                = 37
	SQLParserINTO                  = 38
	SQLParserINTERVAL              = 39
	SQLParserIS                    = 40
	SQLParserJOIN                  = 41
	SQLParserKEYS                  = 42
	SQLParserLEFT                  = 43
	SQLParserLIKE                  = 44
	SQLParserLIMIT                 = 45
	SQLParserLOG                   = 46
	SQLParserLOGICAL               = 47
	SQLParserMASTER                = 48
	SQLParserMEMORY_DATABASES      = 49
	SQLParserMETRICS               = 50
	SQLParserMETRIC                = 51
	SQLParserMETADATA              = 52
	SQLParserMETADATAS             = 53
	SQLParserNAMESPACE             = 54
	SQLParserNAMESPACES            = 55
	SQLParserNULL                  = 56
	SQLParserNOT                   = 57
	SQLParserNOW                   = 58
	SQLParserON                    = 59
	SQLParserOR                    = 60
	SQLParserORDER                 = 61
	SQLParserREQUESTS              = 62
	SQLParserREPLICATIONS          = 63
	SQLParserRIGHT                 = 64
	SQLParserROLLUP                = 65
	SQLParserSELECT                = 66
	SQLParserSHOW                  = 67
	SQLParserSTATE                 = 68
	SQLParserSTORAGE               = 69
	SQLParserTABLE_NAMES           = 70
	SQLParserTIMESTAMP             = 71
	SQLParserTRACE                 = 72
	SQLParserTRUE                  = 73
	SQLParserTYPE                  = 74
	SQLParserTYPES                 = 75
	SQLParserVALUES                = 76
	SQLParserWHERE                 = 77
	SQLParserWITH                  = 78
	SQLParserWITHIN                = 79
	SQLParserUSING                 = 80
	SQLParserUSE                   = 81
	SQLParserSECOND                = 82
	SQLParserMINUTE                = 83
	SQLParserHOUR                  = 84
	SQLParserDAY                   = 85
	SQLParserMONTH                 = 86
	SQLParserYEAR                  = 87
	SQLParserEQ                    = 88
	SQLParserNEQ                   = 89
	SQLParserLT                    = 90
	SQLParserLTE                   = 91
	SQLParserGT                    = 92
	SQLParserGTE                   = 93
	SQLParserPLUS                  = 94
	SQLParserMINUS                 = 95
	SQLParserASTERISK              = 96
	SQLParserSLASH                 = 97
	SQLParserPERCENT               = 98
	SQLParserREGEXP                = 99
	SQLParserNEQREGEXP             = 100
	SQLParserEXCLAMATION_SYMBOL    = 101
	SQLParserDOT                   = 102
	SQLParserLR_BRACKET            = 103
	SQLParserRR_BRACKET            = 104
	SQLParserCOMMA                 = 105
	SQLParserAT                    = 106
	SQLParserSEMICOLON             = 107
	SQLParserSTRING                = 108
	SQLParserINTEGER_VALUE         = 109
	SQLParserDECIMAL_VALUE         = 110
	SQLParserDOUBLE_VALUE          = 111
	SQLParserIDENTIFIER            = 112
	SQLParserDIGIT_IDENTIFIER      = 113
	SQLParserQUOTED_IDENTIFIER     = 114
	SQLParserBACKQUOTED_IDENTIFIER = 115
)

// SQLParser rules.
const (
	SQLParserRULE_statement             = 0
	SQLParserRULE_streamingApp          = 1
	SQLParserRULE_streamingQuery        = 2
	SQLParserRULE_ddlStatement          = 3
	SQLParserRULE_dmlStatement          = 4
	SQLParserRULE_adminStatement        = 5
	SQLParserRULE_utilityStatement      = 6
	SQLParserRULE_explainOption         = 7
	SQLParserRULE_createDatabase        = 8
	SQLParserRULE_databaseOptions       = 9
	SQLParserRULE_createDatabaseOptions = 10
	SQLParserRULE_rollupOptions         = 11
	SQLParserRULE_dropDatabase          = 12
	SQLParserRULE_createBroker          = 13
	SQLParserRULE_flushDatabase         = 14
	SQLParserRULE_compactDatabase       = 15
	SQLParserRULE_showStatement         = 16
	SQLParserRULE_useStatement          = 17
	SQLParserRULE_query                 = 18
	SQLParserRULE_with                  = 19
	SQLParserRULE_namedQuery            = 20
	SQLParserRULE_queryNoWith           = 21
	SQLParserRULE_queryTerm             = 22
	SQLParserRULE_queryPrimary          = 23
	SQLParserRULE_querySpecification    = 24
	SQLParserRULE_selectItem            = 25
	SQLParserRULE_relation              = 26
	SQLParserRULE_joinType              = 27
	SQLParserRULE_joinCriteria          = 28
	SQLParserRULE_aliasedRelation       = 29
	SQLParserRULE_relationPrimary       = 30
	SQLParserRULE_groupBy               = 31
	SQLParserRULE_groupingElement       = 32
	SQLParserRULE_groupingSet           = 33
	SQLParserRULE_having                = 34
	SQLParserRULE_orderBy               = 35
	SQLParserRULE_sortItem              = 36
	SQLParserRULE_limitRowCount         = 37
	SQLParserRULE_expression            = 38
	SQLParserRULE_booleanExpression     = 39
	SQLParserRULE_valueExpression       = 40
	SQLParserRULE_primaryExpression     = 41
	SQLParserRULE_predicate             = 42
	SQLParserRULE_comparisonOperator    = 43
	SQLParserRULE_qualifiedName         = 44
	SQLParserRULE_columnAliases         = 45
	SQLParserRULE_appAnnotation         = 46
	SQLParserRULE_annotation            = 47
	SQLParserRULE_annotation_element    = 48
	SQLParserRULE_properties            = 49
	SQLParserRULE_propertyAssignments   = 50
	SQLParserRULE_property              = 51
	SQLParserRULE_propertyValue         = 52
	SQLParserRULE_booleanValue          = 53
	SQLParserRULE_string                = 54
	SQLParserRULE_identifier            = 55
	SQLParserRULE_number                = 56
	SQLParserRULE_interval              = 57
	SQLParserRULE_intervalUnit          = 58
	SQLParserRULE_nonReserved           = 59
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
	StreamingApp() IStreamingAppContext
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

func (s *StatementContext) StreamingApp() IStreamingAppContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStreamingAppContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStreamingAppContext)
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
	p.SetState(127)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(120)
			p.DdlStatement()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(121)
			p.DmlStatement()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(122)
			p.AdminStatement()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(123)
			p.UtilityStatement()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(124)
			p.StreamingApp()
		}
		{
			p.SetState(125)
			p.Match(SQLParserEOF)
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

// IStreamingAppContext is an interface to support dynamic dispatch.
type IStreamingAppContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAppAnnotation() []IAppAnnotationContext
	AppAnnotation(i int) IAppAnnotationContext
	AllStreamingQuery() []IStreamingQueryContext
	StreamingQuery(i int) IStreamingQueryContext

	// IsStreamingAppContext differentiates from other interfaces.
	IsStreamingAppContext()
}

type StreamingAppContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStreamingAppContext() *StreamingAppContext {
	var p = new(StreamingAppContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_streamingApp
	return p
}

func InitEmptyStreamingAppContext(p *StreamingAppContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_streamingApp
}

func (*StreamingAppContext) IsStreamingAppContext() {}

func NewStreamingAppContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StreamingAppContext {
	var p = new(StreamingAppContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_streamingApp

	return p
}

func (s *StreamingAppContext) GetParser() antlr.Parser { return s.parser }

func (s *StreamingAppContext) AllAppAnnotation() []IAppAnnotationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAppAnnotationContext); ok {
			len++
		}
	}

	tst := make([]IAppAnnotationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAppAnnotationContext); ok {
			tst[i] = t.(IAppAnnotationContext)
			i++
		}
	}

	return tst
}

func (s *StreamingAppContext) AppAnnotation(i int) IAppAnnotationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAppAnnotationContext); ok {
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

	return t.(IAppAnnotationContext)
}

func (s *StreamingAppContext) AllStreamingQuery() []IStreamingQueryContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStreamingQueryContext); ok {
			len++
		}
	}

	tst := make([]IStreamingQueryContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStreamingQueryContext); ok {
			tst[i] = t.(IStreamingQueryContext)
			i++
		}
	}

	return tst
}

func (s *StreamingAppContext) StreamingQuery(i int) IStreamingQueryContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStreamingQueryContext); ok {
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

	return t.(IStreamingQueryContext)
}

func (s *StreamingAppContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StreamingAppContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StreamingAppContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterStreamingApp(s)
	}
}

func (s *StreamingAppContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitStreamingApp(s)
	}
}

func (s *StreamingAppContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitStreamingApp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) StreamingApp() (localctx IStreamingAppContext) {
	localctx = NewStreamingAppContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SQLParserRULE_streamingApp)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(132)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(129)
				p.AppAnnotation()
			}

		}
		p.SetState(134)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(138)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserEXPLAIN || _la == SQLParserINSERT || ((int64((_la-66)) & ^0x3f) == 0 && ((int64(1)<<(_la-66))&1236950585345) != 0) {
		{
			p.SetState(135)
			p.StreamingQuery()
		}

		p.SetState(140)
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

// IStreamingQueryContext is an interface to support dynamic dispatch.
type IStreamingQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DmlStatement() IDmlStatementContext
	AllAnnotation() []IAnnotationContext
	Annotation(i int) IAnnotationContext
	SEMICOLON() antlr.TerminalNode

	// IsStreamingQueryContext differentiates from other interfaces.
	IsStreamingQueryContext()
}

type StreamingQueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStreamingQueryContext() *StreamingQueryContext {
	var p = new(StreamingQueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_streamingQuery
	return p
}

func InitEmptyStreamingQueryContext(p *StreamingQueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_streamingQuery
}

func (*StreamingQueryContext) IsStreamingQueryContext() {}

func NewStreamingQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StreamingQueryContext {
	var p = new(StreamingQueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_streamingQuery

	return p
}

func (s *StreamingQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *StreamingQueryContext) DmlStatement() IDmlStatementContext {
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

func (s *StreamingQueryContext) AllAnnotation() []IAnnotationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAnnotationContext); ok {
			len++
		}
	}

	tst := make([]IAnnotationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAnnotationContext); ok {
			tst[i] = t.(IAnnotationContext)
			i++
		}
	}

	return tst
}

func (s *StreamingQueryContext) Annotation(i int) IAnnotationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnnotationContext); ok {
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

	return t.(IAnnotationContext)
}

func (s *StreamingQueryContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(SQLParserSEMICOLON, 0)
}

func (s *StreamingQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StreamingQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StreamingQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterStreamingQuery(s)
	}
}

func (s *StreamingQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitStreamingQuery(s)
	}
}

func (s *StreamingQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitStreamingQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) StreamingQuery() (localctx IStreamingQueryContext) {
	localctx = NewStreamingQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SQLParserRULE_streamingQuery)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserAT {
		{
			p.SetState(141)
			p.Annotation()
		}

		p.SetState(146)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(147)
		p.DmlStatement()
	}
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserSEMICOLON {
		{
			p.SetState(148)
			p.Match(SQLParserSEMICOLON)
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
	p.EnterRule(localctx, 6, SQLParserRULE_ddlStatement)
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(151)
			p.CreateDatabase()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(152)
			p.DropDatabase()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(153)
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

type InsertIntoContext struct {
	DmlStatementContext
}

func NewInsertIntoContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InsertIntoContext {
	var p = new(InsertIntoContext)

	InitEmptyDmlStatementContext(&p.DmlStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*DmlStatementContext))

	return p
}

func (s *InsertIntoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InsertIntoContext) INSERT() antlr.TerminalNode {
	return s.GetToken(SQLParserINSERT, 0)
}

func (s *InsertIntoContext) INTO() antlr.TerminalNode {
	return s.GetToken(SQLParserINTO, 0)
}

func (s *InsertIntoContext) QualifiedName() IQualifiedNameContext {
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

func (s *InsertIntoContext) Query() IQueryContext {
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

func (s *InsertIntoContext) ColumnAliases() IColumnAliasesContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnAliasesContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnAliasesContext)
}

func (s *InsertIntoContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterInsertInto(s)
	}
}

func (s *InsertIntoContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitInsertInto(s)
	}
}

func (s *InsertIntoContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitInsertInto(s)

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
	p.EnterRule(localctx, 8, SQLParserRULE_dmlStatement)
	var _la int

	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext()) {
	case 1:
		localctx = NewStatementDefaultContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(156)
			p.Query()
		}

	case 2:
		localctx = NewExplainContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(157)
			p.Match(SQLParserEXPLAIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(169)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(158)
				p.Match(SQLParserLR_BRACKET)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(159)
				p.ExplainOption()
			}
			p.SetState(164)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == SQLParserCOMMA {
				{
					p.SetState(160)
					p.Match(SQLParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(161)
					p.ExplainOption()
				}

				p.SetState(166)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(167)
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
			p.SetState(171)
			p.DmlStatement()
		}

	case 3:
		localctx = NewExplainAnalyzeContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(172)
			p.Match(SQLParserEXPLAIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(173)
			p.Match(SQLParserANALYZE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(174)
			p.DmlStatement()
		}

	case 4:
		localctx = NewInsertIntoContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(175)
			p.Match(SQLParserINSERT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(176)
			p.Match(SQLParserINTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(177)
			p.QualifiedName()
		}
		p.SetState(179)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 8, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(178)
				p.ColumnAliases()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}
		{
			p.SetState(181)
			p.Query()
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
	p.EnterRule(localctx, 10, SQLParserRULE_adminStatement)
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserFLUSH:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(185)
			p.FlushDatabase()
		}

	case SQLParserCOMPACT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(186)
			p.CompactDatabase()
		}

	case SQLParserSHOW:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(187)
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
	p.EnterRule(localctx, 12, SQLParserRULE_utilityStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(190)
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
	p.EnterRule(localctx, 14, SQLParserRULE_explainOption)
	var _la int

	localctx = NewExplainTypeContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(192)
		p.Match(SQLParserTYPE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(193)

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
	AllDatabaseOptions() []IDatabaseOptionsContext
	DatabaseOptions(i int) IDatabaseOptionsContext

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

func (s *CreateDatabaseContext) AllDatabaseOptions() []IDatabaseOptionsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDatabaseOptionsContext); ok {
			len++
		}
	}

	tst := make([]IDatabaseOptionsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDatabaseOptionsContext); ok {
			tst[i] = t.(IDatabaseOptionsContext)
			i++
		}
	}

	return tst
}

func (s *CreateDatabaseContext) DatabaseOptions(i int) IDatabaseOptionsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDatabaseOptionsContext); ok {
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

	return t.(IDatabaseOptionsContext)
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
	p.EnterRule(localctx, 16, SQLParserRULE_createDatabase)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(195)
		p.Match(SQLParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(196)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(197)

		var _x = p.QualifiedName()

		localctx.(*CreateDatabaseContext).name = _x
	}
	p.SetState(201)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64((_la-25)) & ^0x3f) == 0 && ((int64(1)<<(_la-25))&9008298766368769) != 0 {
		{
			p.SetState(198)
			p.DatabaseOptions()
		}

		p.SetState(203)
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

// IDatabaseOptionsContext is an interface to support dynamic dispatch.
type IDatabaseOptionsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsDatabaseOptionsContext differentiates from other interfaces.
	IsDatabaseOptionsContext()
}

type DatabaseOptionsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDatabaseOptionsContext() *DatabaseOptionsContext {
	var p = new(DatabaseOptionsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_databaseOptions
	return p
}

func InitEmptyDatabaseOptionsContext(p *DatabaseOptionsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_databaseOptions
}

func (*DatabaseOptionsContext) IsDatabaseOptionsContext() {}

func NewDatabaseOptionsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DatabaseOptionsContext {
	var p = new(DatabaseOptionsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_databaseOptions

	return p
}

func (s *DatabaseOptionsContext) GetParser() antlr.Parser { return s.parser }

func (s *DatabaseOptionsContext) CopyAll(ctx *DatabaseOptionsContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *DatabaseOptionsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DatabaseOptionsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DatabaseOptsContext struct {
	DatabaseOptionsContext
}

func NewDatabaseOptsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DatabaseOptsContext {
	var p = new(DatabaseOptsContext)

	InitEmptyDatabaseOptionsContext(&p.DatabaseOptionsContext)
	p.parser = parser
	p.CopyAll(ctx.(*DatabaseOptionsContext))

	return p
}

func (s *DatabaseOptsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DatabaseOptsContext) AllCreateDatabaseOptions() []ICreateDatabaseOptionsContext {
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

func (s *DatabaseOptsContext) CreateDatabaseOptions(i int) ICreateDatabaseOptionsContext {
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

func (s *DatabaseOptsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *DatabaseOptsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *DatabaseOptsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterDatabaseOpts(s)
	}
}

func (s *DatabaseOptsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitDatabaseOpts(s)
	}
}

func (s *DatabaseOptsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitDatabaseOpts(s)

	default:
		return t.VisitChildren(s)
	}
}

type WithPropsContext struct {
	DatabaseOptionsContext
}

func NewWithPropsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WithPropsContext {
	var p = new(WithPropsContext)

	InitEmptyDatabaseOptionsContext(&p.DatabaseOptionsContext)
	p.parser = parser
	p.CopyAll(ctx.(*DatabaseOptionsContext))

	return p
}

func (s *WithPropsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WithPropsContext) WITH() antlr.TerminalNode {
	return s.GetToken(SQLParserWITH, 0)
}

func (s *WithPropsContext) Properties() IPropertiesContext {
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

func (s *WithPropsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterWithProps(s)
	}
}

func (s *WithPropsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitWithProps(s)
	}
}

func (s *WithPropsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitWithProps(s)

	default:
		return t.VisitChildren(s)
	}
}

type RollupPropsContext struct {
	DatabaseOptionsContext
}

func NewRollupPropsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RollupPropsContext {
	var p = new(RollupPropsContext)

	InitEmptyDatabaseOptionsContext(&p.DatabaseOptionsContext)
	p.parser = parser
	p.CopyAll(ctx.(*DatabaseOptionsContext))

	return p
}

func (s *RollupPropsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RollupPropsContext) ROLLUP() antlr.TerminalNode {
	return s.GetToken(SQLParserROLLUP, 0)
}

func (s *RollupPropsContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *RollupPropsContext) AllRollupOptions() []IRollupOptionsContext {
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

func (s *RollupPropsContext) RollupOptions(i int) IRollupOptionsContext {
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

func (s *RollupPropsContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *RollupPropsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *RollupPropsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *RollupPropsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterRollupProps(s)
	}
}

func (s *RollupPropsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitRollupProps(s)
	}
}

func (s *RollupPropsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitRollupProps(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) DatabaseOptions() (localctx IDatabaseOptionsContext) {
	localctx = NewDatabaseOptionsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SQLParserRULE_databaseOptions)
	var _la int

	p.SetState(226)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserENGINE:
		localctx = NewDatabaseOptsContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(204)
			p.CreateDatabaseOptions()
		}
		p.SetState(209)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(205)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(206)
				p.CreateDatabaseOptions()
			}

			p.SetState(211)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	case SQLParserWITH:
		localctx = NewWithPropsContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(212)
			p.Match(SQLParserWITH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(213)
			p.Properties()
		}

	case SQLParserROLLUP:
		localctx = NewRollupPropsContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(214)
			p.Match(SQLParserROLLUP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(215)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(216)
			p.RollupOptions()
		}
		p.SetState(221)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(217)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(218)
				p.RollupOptions()
			}

			p.SetState(223)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(224)
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
	p.EnterRule(localctx, 20, SQLParserRULE_createDatabaseOptions)
	var _la int

	localctx = NewEngineOptionContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(228)
		p.Match(SQLParserENGINE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(230)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserEQ {
		{
			p.SetState(229)
			p.Match(SQLParserEQ)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(232)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*EngineOptionContext).value = _lt

		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-46)) & ^0x3f) == 0 && ((int64(1)<<(_la-46))&67108897) != 0) {
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
	p.EnterRule(localctx, 22, SQLParserRULE_rollupOptions)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(234)
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
	p.EnterRule(localctx, 24, SQLParserRULE_dropDatabase)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(236)
		p.Match(SQLParserDROP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(237)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(240)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 16, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(238)
			p.Match(SQLParserIF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(239)
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
		p.SetState(242)

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
	p.EnterRule(localctx, 26, SQLParserRULE_createBroker)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(244)
		p.Match(SQLParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(245)
		p.Match(SQLParserBROKER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(246)

		var _x = p.QualifiedName()

		localctx.(*CreateBrokerContext).name = _x
	}
	p.SetState(249)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWITH {
		{
			p.SetState(247)
			p.Match(SQLParserWITH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(248)
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
	p.EnterRule(localctx, 28, SQLParserRULE_flushDatabase)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(251)
		p.Match(SQLParserFLUSH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(252)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(253)

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
	p.EnterRule(localctx, 30, SQLParserRULE_compactDatabase)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(255)
		p.Match(SQLParserCOMPACT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(256)
		p.Match(SQLParserDATABASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(257)

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

type ShowMemoryDatabasesContext struct {
	ShowStatementContext
}

func NewShowMemoryDatabasesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowMemoryDatabasesContext {
	var p = new(ShowMemoryDatabasesContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowMemoryDatabasesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowMemoryDatabasesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowMemoryDatabasesContext) MEMORY_DATABASES() antlr.TerminalNode {
	return s.GetToken(SQLParserMEMORY_DATABASES, 0)
}

func (s *ShowMemoryDatabasesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowMemoryDatabases(s)
	}
}

func (s *ShowMemoryDatabasesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowMemoryDatabases(s)
	}
}

func (s *ShowMemoryDatabasesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowMemoryDatabases(s)

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

type ShowTableNamesContext struct {
	ShowStatementContext
	tableName IExpressionContext
}

func NewShowTableNamesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowTableNamesContext {
	var p = new(ShowTableNamesContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowTableNamesContext) GetTableName() IExpressionContext { return s.tableName }

func (s *ShowTableNamesContext) SetTableName(v IExpressionContext) { s.tableName = v }

func (s *ShowTableNamesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowTableNamesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowTableNamesContext) TABLE_NAMES() antlr.TerminalNode {
	return s.GetToken(SQLParserTABLE_NAMES, 0)
}

func (s *ShowTableNamesContext) FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserFROM, 0)
}

func (s *ShowTableNamesContext) QualifiedName() IQualifiedNameContext {
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

func (s *ShowTableNamesContext) LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserLIKE, 0)
}

func (s *ShowTableNamesContext) Expression() IExpressionContext {
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

func (s *ShowTableNamesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowTableNames(s)
	}
}

func (s *ShowTableNamesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowTableNames(s)
	}
}

func (s *ShowTableNamesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowTableNames(s)

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

type ShowNamespacesContext struct {
	ShowStatementContext
	namespace IExpressionContext
}

func NewShowNamespacesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowNamespacesContext {
	var p = new(ShowNamespacesContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowNamespacesContext) GetNamespace() IExpressionContext { return s.namespace }

func (s *ShowNamespacesContext) SetNamespace(v IExpressionContext) { s.namespace = v }

func (s *ShowNamespacesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowNamespacesContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowNamespacesContext) NAMESPACES() antlr.TerminalNode {
	return s.GetToken(SQLParserNAMESPACES, 0)
}

func (s *ShowNamespacesContext) LIKE() antlr.TerminalNode {
	return s.GetToken(SQLParserLIKE, 0)
}

func (s *ShowNamespacesContext) Expression() IExpressionContext {
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

type ShowColumnsContext struct {
	ShowStatementContext
}

func NewShowColumnsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowColumnsContext {
	var p = new(ShowColumnsContext)

	InitEmptyShowStatementContext(&p.ShowStatementContext)
	p.parser = parser
	p.CopyAll(ctx.(*ShowStatementContext))

	return p
}

func (s *ShowColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowColumnsContext) SHOW() antlr.TerminalNode {
	return s.GetToken(SQLParserSHOW, 0)
}

func (s *ShowColumnsContext) COLUMNS() antlr.TerminalNode {
	return s.GetToken(SQLParserCOLUMNS, 0)
}

func (s *ShowColumnsContext) FROM() antlr.TerminalNode {
	return s.GetToken(SQLParserFROM, 0)
}

func (s *ShowColumnsContext) QualifiedName() IQualifiedNameContext {
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

func (s *ShowColumnsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterShowColumns(s)
	}
}

func (s *ShowColumnsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitShowColumns(s)
	}
}

func (s *ShowColumnsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitShowColumns(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ShowStatement() (localctx IShowStatementContext) {
	localctx = NewShowStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, SQLParserRULE_showStatement)
	var _la int

	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewShowMasterContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(259)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(260)
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
			p.SetState(261)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(262)
			p.Match(SQLParserBROKERS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewShowDatabasesContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(263)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(264)
			p.Match(SQLParserDATABASES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewShowRequestsContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(265)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(266)
			p.Match(SQLParserREQUESTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewShowLimitContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(267)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(268)
			p.Match(SQLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		localctx = NewShowMemoryDatabasesContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(269)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(270)
			p.Match(SQLParserMEMORY_DATABASES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		localctx = NewShowReplicationsContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(271)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(272)
			p.Match(SQLParserREPLICATIONS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		localctx = NewShowNamespacesContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(273)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(274)
			p.Match(SQLParserNAMESPACES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(277)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserLIKE {
			{
				p.SetState(275)
				p.Match(SQLParserLIKE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(276)

				var _x = p.Expression()

				localctx.(*ShowNamespacesContext).namespace = _x
			}

		}

	case 9:
		localctx = NewShowTableNamesContext(p, localctx)
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(279)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(280)
			p.Match(SQLParserTABLE_NAMES)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(283)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserFROM {
			{
				p.SetState(281)
				p.Match(SQLParserFROM)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(282)
				p.QualifiedName()
			}

		}
		p.SetState(287)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserLIKE {
			{
				p.SetState(285)
				p.Match(SQLParserLIKE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(286)

				var _x = p.Expression()

				localctx.(*ShowTableNamesContext).tableName = _x
			}

		}

	case 10:
		localctx = NewShowColumnsContext(p, localctx)
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(289)
			p.Match(SQLParserSHOW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(290)
			p.Match(SQLParserCOLUMNS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(291)
			p.Match(SQLParserFROM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(292)
			p.QualifiedName()
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
	p.EnterRule(localctx, 34, SQLParserRULE_useStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(295)
		p.Match(SQLParserUSE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(296)

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
	p.EnterRule(localctx, 36, SQLParserRULE_query)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(299)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWITH {
		{
			p.SetState(298)
			p.With()
		}

	}
	{
		p.SetState(301)
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
	p.EnterRule(localctx, 38, SQLParserRULE_with)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(303)
		p.Match(SQLParserWITH)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(304)
		p.NamedQuery()
	}
	p.SetState(309)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(305)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(306)
			p.NamedQuery()
		}

		p.SetState(311)
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
	p.EnterRule(localctx, 40, SQLParserRULE_namedQuery)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(312)

		var _x = p.Identifier()

		localctx.(*NamedQueryContext).name = _x
	}
	{
		p.SetState(313)
		p.Match(SQLParserAS)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(314)
		p.Match(SQLParserLR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(315)
		p.Query()
	}
	{
		p.SetState(316)
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
	p.EnterRule(localctx, 42, SQLParserRULE_queryNoWith)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(318)
		p.QueryTerm()
	}
	p.SetState(322)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserORDER {
		{
			p.SetState(319)
			p.Match(SQLParserORDER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(320)
			p.Match(SQLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(321)
			p.OrderBy()
		}

	}
	p.SetState(326)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserLIMIT {
		{
			p.SetState(324)
			p.Match(SQLParserLIMIT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(325)

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
	p.EnterRule(localctx, 44, SQLParserRULE_queryTerm)
	localctx = NewQueryTermDefaultContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(328)
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
	p.EnterRule(localctx, 46, SQLParserRULE_queryPrimary)
	p.SetState(335)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserSELECT:
		localctx = NewQueryPrimaryDefaultContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(330)
			p.QuerySpecification()
		}

	case SQLParserLR_BRACKET:
		localctx = NewSubqueryContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(331)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(332)
			p.QueryNoWith()
		}
		{
			p.SetState(333)
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
	p.EnterRule(localctx, 48, SQLParserRULE_querySpecification)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(337)
		p.Match(SQLParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(338)
		p.SelectItem()
	}
	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(339)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(340)
			p.SelectItem()
		}

		p.SetState(345)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(355)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserFROM {
		{
			p.SetState(346)
			p.Match(SQLParserFROM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(347)
			p.relation(0)
		}
		p.SetState(352)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(348)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(349)
				p.relation(0)
			}

			p.SetState(354)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	p.SetState(359)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserWHERE {
		{
			p.SetState(357)
			p.Match(SQLParserWHERE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(358)

			var _x = p.booleanExpression(0)

			localctx.(*QuerySpecificationContext).where = _x
		}

	}
	p.SetState(364)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserGROUP {
		{
			p.SetState(361)
			p.Match(SQLParserGROUP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(362)
			p.Match(SQLParserBY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(363)
			p.GroupBy()
		}

	}
	p.SetState(368)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserHAVING {
		{
			p.SetState(366)
			p.Match(SQLParserHAVING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(367)
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
	p.EnterRule(localctx, 50, SQLParserRULE_selectItem)
	p.SetState(382)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 35, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSelectSingleContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(370)
			p.Expression()
		}
		p.SetState(375)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 34, p.GetParserRuleContext()) == 1 {
			p.SetState(372)
			p.GetErrorHandler().Sync(p)

			if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
				{
					p.SetState(371)
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
				p.SetState(374)
				p.Identifier()
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 2:
		localctx = NewSelectAllContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(377)
			p.primaryExpression(0)
		}
		{
			p.SetState(378)
			p.Match(SQLParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(379)
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
			p.SetState(381)
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
	_startState := 52
	p.EnterRecursionRule(localctx, 52, SQLParserRULE_relation, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewRelationDefaultContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(385)
		p.AliasedRelation()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(400)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
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
			p.SetState(387)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				goto errorExit
			}
			p.SetState(396)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetTokenStream().LA(1) {
			case SQLParserCROSS:
				{
					p.SetState(388)
					p.Match(SQLParserCROSS)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(389)
					p.Match(SQLParserJOIN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(390)

					var _x = p.relation(0)

					localctx.(*JoinRelationContext).right = _x
				}

			case SQLParserLEFT, SQLParserRIGHT:
				{
					p.SetState(391)
					p.JoinType()
				}
				{
					p.SetState(392)
					p.Match(SQLParserJOIN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(393)

					var _x = p.relation(0)

					localctx.(*JoinRelationContext).rightRelation = _x
				}
				{
					p.SetState(394)
					p.JoinCriteria()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

		}
		p.SetState(402)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 37, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 54, SQLParserRULE_joinType)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(403)
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
	p.EnterRule(localctx, 56, SQLParserRULE_joinCriteria)
	var _la int

	p.SetState(419)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserON:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(405)
			p.Match(SQLParserON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(406)
			p.booleanExpression(0)
		}

	case SQLParserUSING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(407)
			p.Match(SQLParserUSING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(408)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(409)
			p.Identifier()
		}
		p.SetState(414)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(410)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(411)
				p.Identifier()
			}

			p.SetState(416)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(417)
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
	p.EnterRule(localctx, 58, SQLParserRULE_aliasedRelation)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(421)
		p.RelationPrimary()
	}
	p.SetState(426)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext()) == 1 {
		p.SetState(423)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 40, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(422)
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
			p.SetState(425)
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
	p.EnterRule(localctx, 60, SQLParserRULE_relationPrimary)
	p.SetState(433)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserALL, SQLParserAPP, SQLParserALIVE, SQLParserAND, SQLParserAS, SQLParserASC, SQLParserBETWEEN, SQLParserBROKER, SQLParserBROKERS, SQLParserBY, SQLParserCOMPACT, SQLParserCREATE, SQLParserCROSS, SQLParserCOLUMNS, SQLParserDATABASE, SQLParserDATABASES, SQLParserDEFAULT, SQLParserDESC, SQLParserDISTRIBUTED, SQLParserDROP, SQLParserENGINE, SQLParserESCAPE, SQLParserEXPLAIN, SQLParserEXISTS, SQLParserFALSE, SQLParserFIELDS, SQLParserFLUSH, SQLParserFROM, SQLParserGROUP, SQLParserHAVING, SQLParserIF, SQLParserIN, SQLParserINSERT, SQLParserINTO, SQLParserINTERVAL, SQLParserIS, SQLParserJOIN, SQLParserKEYS, SQLParserLEFT, SQLParserLIKE, SQLParserLIMIT, SQLParserLOG, SQLParserLOGICAL, SQLParserMASTER, SQLParserMEMORY_DATABASES, SQLParserMETRICS, SQLParserMETRIC, SQLParserMETADATA, SQLParserMETADATAS, SQLParserNAMESPACE, SQLParserNAMESPACES, SQLParserNULL, SQLParserNOT, SQLParserNOW, SQLParserON, SQLParserOR, SQLParserORDER, SQLParserREQUESTS, SQLParserREPLICATIONS, SQLParserRIGHT, SQLParserROLLUP, SQLParserSELECT, SQLParserSHOW, SQLParserSTATE, SQLParserSTORAGE, SQLParserTABLE_NAMES, SQLParserTIMESTAMP, SQLParserTRACE, SQLParserTRUE, SQLParserTYPE, SQLParserTYPES, SQLParserVALUES, SQLParserWHERE, SQLParserWITH, SQLParserWITHIN, SQLParserUSING, SQLParserUSE, SQLParserSECOND, SQLParserMINUTE, SQLParserHOUR, SQLParserDAY, SQLParserMONTH, SQLParserYEAR, SQLParserIDENTIFIER, SQLParserDIGIT_IDENTIFIER, SQLParserQUOTED_IDENTIFIER, SQLParserBACKQUOTED_IDENTIFIER:
		localctx = NewTableNameContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(428)
			p.QualifiedName()
		}

	case SQLParserLR_BRACKET:
		localctx = NewSubQueryRelationContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(429)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(430)
			p.Query()
		}
		{
			p.SetState(431)
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
	p.EnterRule(localctx, 62, SQLParserRULE_groupBy)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(435)
		p.GroupingElement()
	}
	p.SetState(440)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(436)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(437)
			p.GroupingElement()
		}

		p.SetState(442)
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
	p.EnterRule(localctx, 64, SQLParserRULE_groupingElement)
	p.SetState(445)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSingleGroupingSetContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(443)
			p.GroupingSet()
		}

	case 2:
		localctx = NewGroupByAllColumnsContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(444)
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
	p.EnterRule(localctx, 66, SQLParserRULE_groupingSet)
	var _la int

	p.SetState(460)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 47, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(447)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(456)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-272) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&4486696800354303) != 0) {
			{
				p.SetState(448)
				p.Expression()
			}
			p.SetState(453)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == SQLParserCOMMA {
				{
					p.SetState(449)
					p.Match(SQLParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(450)
					p.Expression()
				}

				p.SetState(455)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(458)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(459)
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
	p.EnterRule(localctx, 68, SQLParserRULE_having)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(462)
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
	p.EnterRule(localctx, 70, SQLParserRULE_orderBy)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(464)
		p.SortItem()
	}
	p.SetState(469)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(465)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(466)
			p.SortItem()
		}

		p.SetState(471)
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
	p.EnterRule(localctx, 72, SQLParserRULE_sortItem)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(472)
		p.Expression()
	}
	p.SetState(474)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == SQLParserASC || _la == SQLParserDESC {
		{
			p.SetState(473)

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
	p.EnterRule(localctx, 74, SQLParserRULE_limitRowCount)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(476)
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
	p.EnterRule(localctx, 76, SQLParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(478)
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
	_startState := 78
	p.EnterRecursionRule(localctx, 78, SQLParserRULE_booleanExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(484)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) {
	case 1:
		localctx = NewLogicalNotContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(481)

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
			p.SetState(482)
			p.booleanExpression(4)
		}

	case 2:
		localctx = NewPredicatedExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(483)
			p.Predicate()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(494)
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
			p.SetState(492)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 51, p.GetParserRuleContext()) {
			case 1:
				localctx = NewAndContext(p, NewBooleanExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_booleanExpression)
				p.SetState(486)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(487)
					p.Match(SQLParserAND)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(488)
					p.booleanExpression(4)
				}

			case 2:
				localctx = NewOrContext(p, NewBooleanExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_booleanExpression)
				p.SetState(489)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(490)
					p.Match(SQLParserOR)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(491)
					p.booleanExpression(3)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(496)
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
	_startState := 80
	p.EnterRecursionRule(localctx, 80, SQLParserRULE_valueExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewValueExpressionDefaultContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(498)
		p.primaryExpression(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(508)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(506)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 53, p.GetParserRuleContext()) {
			case 1:
				localctx = NewArithmeticBinaryContext(p, NewValueExpressionContext(p, _parentctx, _parentState))
				localctx.(*ArithmeticBinaryContext).left = _prevctx

				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_valueExpression)
				p.SetState(500)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(501)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*ArithmeticBinaryContext).operator = _lt

					_la = p.GetTokenStream().LA(1)

					if !((int64((_la-96)) & ^0x3f) == 0 && ((int64(1)<<(_la-96))&7) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*ArithmeticBinaryContext).operator = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(502)

					var _x = p.valueExpression(3)

					localctx.(*ArithmeticBinaryContext).right = _x
				}

			case 2:
				localctx = NewArithmeticBinaryContext(p, NewValueExpressionContext(p, _parentctx, _parentState))
				localctx.(*ArithmeticBinaryContext).left = _prevctx

				p.PushNewRecursionContext(localctx, _startState, SQLParserRULE_valueExpression)
				p.SetState(503)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(504)

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
					p.SetState(505)

					var _x = p.valueExpression(2)

					localctx.(*ArithmeticBinaryContext).right = _x
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(510)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext())
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

type NullLiteralContext struct {
	PrimaryExpressionContext
}

func NewNullLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullLiteralContext {
	var p = new(NullLiteralContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *NullLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullLiteralContext) NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserNULL, 0)
}

func (s *NullLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterNullLiteral(s)
	}
}

func (s *NullLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitNullLiteral(s)
	}
}

func (s *NullLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitNullLiteral(s)

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

type IntervalLiteralContext struct {
	PrimaryExpressionContext
}

func NewIntervalLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntervalLiteralContext {
	var p = new(IntervalLiteralContext)

	InitEmptyPrimaryExpressionContext(&p.PrimaryExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryExpressionContext))

	return p
}

func (s *IntervalLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalLiteralContext) Interval() IIntervalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntervalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntervalContext)
}

func (s *IntervalLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterIntervalLiteral(s)
	}
}

func (s *IntervalLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitIntervalLiteral(s)
	}
}

func (s *IntervalLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitIntervalLiteral(s)

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
	_startState := 82
	p.EnterRecursionRule(localctx, 82, SQLParserRULE_primaryExpression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(536)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 57, p.GetParserRuleContext()) {
	case 1:
		localctx = NewNumericLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(512)
			p.Number()
		}

	case 2:
		localctx = NewIntervalLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(513)
			p.Interval()
		}

	case 3:
		localctx = NewBooleanLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(514)
			p.BooleanValue()
		}

	case 4:
		localctx = NewNullLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(515)
			p.Match(SQLParserNULL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewStringLiteralContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(516)
			p.String_()
		}

	case 6:
		localctx = NewFunctionCallContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(517)
			p.QualifiedName()
		}
		{
			p.SetState(518)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(527)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-272) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&4486696800354303) != 0) {
			{
				p.SetState(519)
				p.Expression()
			}
			p.SetState(524)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			for _la == SQLParserCOMMA {
				{
					p.SetState(520)
					p.Match(SQLParserCOMMA)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(521)
					p.Expression()
				}

				p.SetState(526)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)
			}

		}
		{
			p.SetState(529)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		localctx = NewColumnReferenceContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(531)
			p.Identifier()
		}

	case 8:
		localctx = NewParenExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(532)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(533)
			p.Expression()
		}
		{
			p.SetState(534)
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
	p.SetState(543)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext())
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
			p.SetState(538)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				goto errorExit
			}
			{
				p.SetState(539)
				p.Match(SQLParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(540)

				var _x = p.Identifier()

				localctx.(*DereferenceContext).fieldName = _x
			}

		}
		p.SetState(545)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 58, p.GetParserRuleContext())
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

type BetweenPredicateContext struct {
	PredicateContext
	lower IValueExpressionContext
	upper IValueExpressionContext
}

func NewBetweenPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BetweenPredicateContext {
	var p = new(BetweenPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *BetweenPredicateContext) GetLower() IValueExpressionContext { return s.lower }

func (s *BetweenPredicateContext) GetUpper() IValueExpressionContext { return s.upper }

func (s *BetweenPredicateContext) SetLower(v IValueExpressionContext) { s.lower = v }

func (s *BetweenPredicateContext) SetUpper(v IValueExpressionContext) { s.upper = v }

func (s *BetweenPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BetweenPredicateContext) TIMESTAMP() antlr.TerminalNode {
	return s.GetToken(SQLParserTIMESTAMP, 0)
}

func (s *BetweenPredicateContext) BETWEEN() antlr.TerminalNode {
	return s.GetToken(SQLParserBETWEEN, 0)
}

func (s *BetweenPredicateContext) AND() antlr.TerminalNode {
	return s.GetToken(SQLParserAND, 0)
}

func (s *BetweenPredicateContext) AllValueExpression() []IValueExpressionContext {
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

func (s *BetweenPredicateContext) ValueExpression(i int) IValueExpressionContext {
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

func (s *BetweenPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterBetweenPredicate(s)
	}
}

func (s *BetweenPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitBetweenPredicate(s)
	}
}

func (s *BetweenPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitBetweenPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

type TimestampPredicateContext struct {
	PredicateContext
	operator IComparisonOperatorContext
	right    IValueExpressionContext
}

func NewTimestampPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TimestampPredicateContext {
	var p = new(TimestampPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *TimestampPredicateContext) GetOperator() IComparisonOperatorContext { return s.operator }

func (s *TimestampPredicateContext) GetRight() IValueExpressionContext { return s.right }

func (s *TimestampPredicateContext) SetOperator(v IComparisonOperatorContext) { s.operator = v }

func (s *TimestampPredicateContext) SetRight(v IValueExpressionContext) { s.right = v }

func (s *TimestampPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimestampPredicateContext) TIMESTAMP() antlr.TerminalNode {
	return s.GetToken(SQLParserTIMESTAMP, 0)
}

func (s *TimestampPredicateContext) ComparisonOperator() IComparisonOperatorContext {
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

func (s *TimestampPredicateContext) ValueExpression() IValueExpressionContext {
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

func (s *TimestampPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterTimestampPredicate(s)
	}
}

func (s *TimestampPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitTimestampPredicate(s)
	}
}

func (s *TimestampPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitTimestampPredicate(s)

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

type NullPredicateContext struct {
	PredicateContext
	left IValueExpressionContext
}

func NewNullPredicateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullPredicateContext {
	var p = new(NullPredicateContext)

	InitEmptyPredicateContext(&p.PredicateContext)
	p.parser = parser
	p.CopyAll(ctx.(*PredicateContext))

	return p
}

func (s *NullPredicateContext) GetLeft() IValueExpressionContext { return s.left }

func (s *NullPredicateContext) SetLeft(v IValueExpressionContext) { s.left = v }

func (s *NullPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullPredicateContext) IS() antlr.TerminalNode {
	return s.GetToken(SQLParserIS, 0)
}

func (s *NullPredicateContext) NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserNULL, 0)
}

func (s *NullPredicateContext) ValueExpression() IValueExpressionContext {
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

func (s *NullPredicateContext) NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserNOT, 0)
}

func (s *NullPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterNullPredicate(s)
	}
}

func (s *NullPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitNullPredicate(s)
	}
}

func (s *NullPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitNullPredicate(s)

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
	p.EnterRule(localctx, 84, SQLParserRULE_predicate)
	var _la int

	p.SetState(599)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 65, p.GetParserRuleContext()) {
	case 1:
		localctx = NewTimestampPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(546)
			p.Match(SQLParserTIMESTAMP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(547)

			var _x = p.ComparisonOperator()

			localctx.(*TimestampPredicateContext).operator = _x
		}
		{
			p.SetState(548)

			var _x = p.valueExpression(0)

			localctx.(*TimestampPredicateContext).right = _x
		}

	case 2:
		localctx = NewBinaryComparisonPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(550)

			var _x = p.valueExpression(0)

			localctx.(*BinaryComparisonPredicateContext).left = _x
		}
		{
			p.SetState(551)

			var _x = p.ComparisonOperator()

			localctx.(*BinaryComparisonPredicateContext).operator = _x
		}
		{
			p.SetState(552)

			var _x = p.valueExpression(0)

			localctx.(*BinaryComparisonPredicateContext).right = _x
		}

	case 3:
		localctx = NewBetweenPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(554)
			p.Match(SQLParserTIMESTAMP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(555)
			p.Match(SQLParserBETWEEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(556)

			var _x = p.valueExpression(0)

			localctx.(*BetweenPredicateContext).lower = _x
		}
		{
			p.SetState(557)
			p.Match(SQLParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(558)

			var _x = p.valueExpression(0)

			localctx.(*BetweenPredicateContext).upper = _x
		}

	case 4:
		localctx = NewInPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(560)

			var _x = p.valueExpression(0)

			localctx.(*InPredicateContext).left = _x
		}
		p.SetState(562)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserNOT {
			{
				p.SetState(561)
				p.Match(SQLParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(564)
			p.Match(SQLParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(565)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(566)
			p.Expression()
		}
		p.SetState(571)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(567)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(568)
				p.Expression()
			}

			p.SetState(573)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(574)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewLikePredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(576)

			var _x = p.valueExpression(0)

			localctx.(*LikePredicateContext).left = _x
		}
		p.SetState(578)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserNOT {
			{
				p.SetState(577)
				p.Match(SQLParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(580)
			p.Match(SQLParserLIKE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(581)

			var _x = p.valueExpression(0)

			localctx.(*LikePredicateContext).pattern = _x
		}
		p.SetState(584)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 62, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(582)
				p.Match(SQLParserESCAPE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(583)

				var _x = p.valueExpression(0)

				localctx.(*LikePredicateContext).escape = _x
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 6:
		localctx = NewRegexpPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(586)

			var _x = p.valueExpression(0)

			localctx.(*RegexpPredicateContext).left = _x
		}
		{
			p.SetState(587)

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
		p.SetState(589)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 63, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(588)

				var _x = p.valueExpression(0)

				localctx.(*RegexpPredicateContext).pattern = _x
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 7:
		localctx = NewNullPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(591)

			var _x = p.valueExpression(0)

			localctx.(*NullPredicateContext).left = _x
		}
		{
			p.SetState(592)
			p.Match(SQLParserIS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(594)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserNOT {
			{
				p.SetState(593)
				p.Match(SQLParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(596)
			p.Match(SQLParserNULL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		localctx = NewValueExpressionPredicateContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(598)
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
	p.EnterRule(localctx, 86, SQLParserRULE_comparisonOperator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(601)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-88)) & ^0x3f) == 0 && ((int64(1)<<(_la-88))&63) != 0) {
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
	p.EnterRule(localctx, 88, SQLParserRULE_qualifiedName)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(603)
		p.Identifier()
	}
	p.SetState(608)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 66, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(604)
				p.Match(SQLParserDOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(605)
				p.Identifier()
			}

		}
		p.SetState(610)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 66, p.GetParserRuleContext())
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

// IColumnAliasesContext is an interface to support dynamic dispatch.
type IColumnAliasesContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LR_BRACKET() antlr.TerminalNode
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext
	RR_BRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsColumnAliasesContext differentiates from other interfaces.
	IsColumnAliasesContext()
}

type ColumnAliasesContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnAliasesContext() *ColumnAliasesContext {
	var p = new(ColumnAliasesContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnAliases
	return p
}

func InitEmptyColumnAliasesContext(p *ColumnAliasesContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_columnAliases
}

func (*ColumnAliasesContext) IsColumnAliasesContext() {}

func NewColumnAliasesContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnAliasesContext {
	var p = new(ColumnAliasesContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_columnAliases

	return p
}

func (s *ColumnAliasesContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnAliasesContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *ColumnAliasesContext) AllIdentifier() []IIdentifierContext {
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

func (s *ColumnAliasesContext) Identifier(i int) IIdentifierContext {
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

func (s *ColumnAliasesContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *ColumnAliasesContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *ColumnAliasesContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *ColumnAliasesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnAliasesContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnAliasesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterColumnAliases(s)
	}
}

func (s *ColumnAliasesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitColumnAliases(s)
	}
}

func (s *ColumnAliasesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitColumnAliases(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) ColumnAliases() (localctx IColumnAliasesContext) {
	localctx = NewColumnAliasesContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, SQLParserRULE_columnAliases)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(611)
		p.Match(SQLParserLR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(612)
		p.Identifier()
	}
	p.SetState(617)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(613)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(614)
			p.Identifier()
		}

		p.SetState(619)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(620)
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

// IAppAnnotationContext is an interface to support dynamic dispatch.
type IAppAnnotationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	APP() antlr.TerminalNode
	LR_BRACKET() antlr.TerminalNode
	AllAnnotation_element() []IAnnotation_elementContext
	Annotation_element(i int) IAnnotation_elementContext
	RR_BRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAppAnnotationContext differentiates from other interfaces.
	IsAppAnnotationContext()
}

type AppAnnotationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAppAnnotationContext() *AppAnnotationContext {
	var p = new(AppAnnotationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_appAnnotation
	return p
}

func InitEmptyAppAnnotationContext(p *AppAnnotationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_appAnnotation
}

func (*AppAnnotationContext) IsAppAnnotationContext() {}

func NewAppAnnotationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AppAnnotationContext {
	var p = new(AppAnnotationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_appAnnotation

	return p
}

func (s *AppAnnotationContext) GetParser() antlr.Parser { return s.parser }

func (s *AppAnnotationContext) AT() antlr.TerminalNode {
	return s.GetToken(SQLParserAT, 0)
}

func (s *AppAnnotationContext) APP() antlr.TerminalNode {
	return s.GetToken(SQLParserAPP, 0)
}

func (s *AppAnnotationContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *AppAnnotationContext) AllAnnotation_element() []IAnnotation_elementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAnnotation_elementContext); ok {
			len++
		}
	}

	tst := make([]IAnnotation_elementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAnnotation_elementContext); ok {
			tst[i] = t.(IAnnotation_elementContext)
			i++
		}
	}

	return tst
}

func (s *AppAnnotationContext) Annotation_element(i int) IAnnotation_elementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnnotation_elementContext); ok {
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

	return t.(IAnnotation_elementContext)
}

func (s *AppAnnotationContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *AppAnnotationContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *AppAnnotationContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *AppAnnotationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AppAnnotationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AppAnnotationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterAppAnnotation(s)
	}
}

func (s *AppAnnotationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitAppAnnotation(s)
	}
}

func (s *AppAnnotationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitAppAnnotation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) AppAnnotation() (localctx IAppAnnotationContext) {
	localctx = NewAppAnnotationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, SQLParserRULE_appAnnotation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(622)
		p.Match(SQLParserAT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(623)
		p.Match(SQLParserAPP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(635)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 69, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(624)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(625)
			p.Annotation_element()
		}
		p.SetState(630)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(626)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(627)
				p.Annotation_element()
			}

			p.SetState(632)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(633)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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

// IAnnotationContext is an interface to support dynamic dispatch.
type IAnnotationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	Identifier() IIdentifierContext
	LR_BRACKET() antlr.TerminalNode
	AllAnnotation_element() []IAnnotation_elementContext
	Annotation_element(i int) IAnnotation_elementContext
	RR_BRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsAnnotationContext differentiates from other interfaces.
	IsAnnotationContext()
}

type AnnotationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnnotationContext() *AnnotationContext {
	var p = new(AnnotationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_annotation
	return p
}

func InitEmptyAnnotationContext(p *AnnotationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_annotation
}

func (*AnnotationContext) IsAnnotationContext() {}

func NewAnnotationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnnotationContext {
	var p = new(AnnotationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_annotation

	return p
}

func (s *AnnotationContext) GetParser() antlr.Parser { return s.parser }

func (s *AnnotationContext) AT() antlr.TerminalNode {
	return s.GetToken(SQLParserAT, 0)
}

func (s *AnnotationContext) Identifier() IIdentifierContext {
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

func (s *AnnotationContext) LR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserLR_BRACKET, 0)
}

func (s *AnnotationContext) AllAnnotation_element() []IAnnotation_elementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAnnotation_elementContext); ok {
			len++
		}
	}

	tst := make([]IAnnotation_elementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAnnotation_elementContext); ok {
			tst[i] = t.(IAnnotation_elementContext)
			i++
		}
	}

	return tst
}

func (s *AnnotationContext) Annotation_element(i int) IAnnotation_elementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnnotation_elementContext); ok {
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

	return t.(IAnnotation_elementContext)
}

func (s *AnnotationContext) RR_BRACKET() antlr.TerminalNode {
	return s.GetToken(SQLParserRR_BRACKET, 0)
}

func (s *AnnotationContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SQLParserCOMMA)
}

func (s *AnnotationContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SQLParserCOMMA, i)
}

func (s *AnnotationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnnotationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnnotationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterAnnotation(s)
	}
}

func (s *AnnotationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitAnnotation(s)
	}
}

func (s *AnnotationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitAnnotation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Annotation() (localctx IAnnotationContext) {
	localctx = NewAnnotationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, SQLParserRULE_annotation)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(637)
		p.Match(SQLParserAT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(638)
		p.Identifier()
	}
	p.SetState(650)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 71, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(639)
			p.Match(SQLParserLR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(640)
			p.Annotation_element()
		}
		p.SetState(645)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == SQLParserCOMMA {
			{
				p.SetState(641)
				p.Match(SQLParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(642)
				p.Annotation_element()
			}

			p.SetState(647)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(648)
			p.Match(SQLParserRR_BRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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

// IAnnotation_elementContext is an interface to support dynamic dispatch.
type IAnnotation_elementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Property() IPropertyContext
	Annotation() IAnnotationContext

	// IsAnnotation_elementContext differentiates from other interfaces.
	IsAnnotation_elementContext()
}

type Annotation_elementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnnotation_elementContext() *Annotation_elementContext {
	var p = new(Annotation_elementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_annotation_element
	return p
}

func InitEmptyAnnotation_elementContext(p *Annotation_elementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_annotation_element
}

func (*Annotation_elementContext) IsAnnotation_elementContext() {}

func NewAnnotation_elementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Annotation_elementContext {
	var p = new(Annotation_elementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_annotation_element

	return p
}

func (s *Annotation_elementContext) GetParser() antlr.Parser { return s.parser }

func (s *Annotation_elementContext) Property() IPropertyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPropertyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPropertyContext)
}

func (s *Annotation_elementContext) Annotation() IAnnotationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAnnotationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAnnotationContext)
}

func (s *Annotation_elementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Annotation_elementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Annotation_elementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterAnnotation_element(s)
	}
}

func (s *Annotation_elementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitAnnotation_element(s)
	}
}

func (s *Annotation_elementContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitAnnotation_element(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Annotation_element() (localctx IAnnotation_elementContext) {
	localctx = NewAnnotation_elementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, SQLParserRULE_annotation_element)
	p.SetState(654)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserALL, SQLParserAPP, SQLParserALIVE, SQLParserAND, SQLParserAS, SQLParserASC, SQLParserBETWEEN, SQLParserBROKER, SQLParserBROKERS, SQLParserBY, SQLParserCOMPACT, SQLParserCREATE, SQLParserCROSS, SQLParserCOLUMNS, SQLParserDATABASE, SQLParserDATABASES, SQLParserDEFAULT, SQLParserDESC, SQLParserDISTRIBUTED, SQLParserDROP, SQLParserENGINE, SQLParserESCAPE, SQLParserEXPLAIN, SQLParserEXISTS, SQLParserFALSE, SQLParserFIELDS, SQLParserFLUSH, SQLParserFROM, SQLParserGROUP, SQLParserHAVING, SQLParserIF, SQLParserIN, SQLParserINSERT, SQLParserINTO, SQLParserINTERVAL, SQLParserIS, SQLParserJOIN, SQLParserKEYS, SQLParserLEFT, SQLParserLIKE, SQLParserLIMIT, SQLParserLOG, SQLParserLOGICAL, SQLParserMASTER, SQLParserMEMORY_DATABASES, SQLParserMETRICS, SQLParserMETRIC, SQLParserMETADATA, SQLParserMETADATAS, SQLParserNAMESPACE, SQLParserNAMESPACES, SQLParserNULL, SQLParserNOT, SQLParserNOW, SQLParserON, SQLParserOR, SQLParserORDER, SQLParserREQUESTS, SQLParserREPLICATIONS, SQLParserRIGHT, SQLParserROLLUP, SQLParserSELECT, SQLParserSHOW, SQLParserSTATE, SQLParserSTORAGE, SQLParserTABLE_NAMES, SQLParserTIMESTAMP, SQLParserTRACE, SQLParserTRUE, SQLParserTYPE, SQLParserTYPES, SQLParserVALUES, SQLParserWHERE, SQLParserWITH, SQLParserWITHIN, SQLParserUSING, SQLParserUSE, SQLParserSECOND, SQLParserMINUTE, SQLParserHOUR, SQLParserDAY, SQLParserMONTH, SQLParserYEAR, SQLParserIDENTIFIER, SQLParserDIGIT_IDENTIFIER, SQLParserQUOTED_IDENTIFIER, SQLParserBACKQUOTED_IDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(652)
			p.Property()
		}

	case SQLParserAT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(653)
			p.Annotation()
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
	p.EnterRule(localctx, 98, SQLParserRULE_properties)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(656)
		p.Match(SQLParserLR_BRACKET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(657)
		p.PropertyAssignments()
	}
	{
		p.SetState(658)
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
	p.EnterRule(localctx, 100, SQLParserRULE_propertyAssignments)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(660)
		p.Property()
	}
	p.SetState(665)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == SQLParserCOMMA {
		{
			p.SetState(661)
			p.Match(SQLParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(662)
			p.Property()
		}

		p.SetState(667)
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
	p.EnterRule(localctx, 102, SQLParserRULE_property)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(668)

		var _x = p.Identifier()

		localctx.(*PropertyContext).name = _x
	}
	{
		p.SetState(669)
		p.Match(SQLParserEQ)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(670)

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
	p.EnterRule(localctx, 104, SQLParserRULE_propertyValue)
	p.SetState(674)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 74, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDefaultPropertyValueContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(672)
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
			p.SetState(673)
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
	p.EnterRule(localctx, 106, SQLParserRULE_booleanValue)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(676)
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
	p.EnterRule(localctx, 108, SQLParserRULE_string)
	localctx = NewBasicStringLiteralContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(678)
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
	p.EnterRule(localctx, 110, SQLParserRULE_identifier)
	p.SetState(685)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case SQLParserIDENTIFIER:
		localctx = NewUnquotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(680)
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
			p.SetState(681)
			p.Match(SQLParserQUOTED_IDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case SQLParserALL, SQLParserAPP, SQLParserALIVE, SQLParserAND, SQLParserAS, SQLParserASC, SQLParserBETWEEN, SQLParserBROKER, SQLParserBROKERS, SQLParserBY, SQLParserCOMPACT, SQLParserCREATE, SQLParserCROSS, SQLParserCOLUMNS, SQLParserDATABASE, SQLParserDATABASES, SQLParserDEFAULT, SQLParserDESC, SQLParserDISTRIBUTED, SQLParserDROP, SQLParserENGINE, SQLParserESCAPE, SQLParserEXPLAIN, SQLParserEXISTS, SQLParserFALSE, SQLParserFIELDS, SQLParserFLUSH, SQLParserFROM, SQLParserGROUP, SQLParserHAVING, SQLParserIF, SQLParserIN, SQLParserINSERT, SQLParserINTO, SQLParserINTERVAL, SQLParserIS, SQLParserJOIN, SQLParserKEYS, SQLParserLEFT, SQLParserLIKE, SQLParserLIMIT, SQLParserLOG, SQLParserLOGICAL, SQLParserMASTER, SQLParserMEMORY_DATABASES, SQLParserMETRICS, SQLParserMETRIC, SQLParserMETADATA, SQLParserMETADATAS, SQLParserNAMESPACE, SQLParserNAMESPACES, SQLParserNULL, SQLParserNOT, SQLParserNOW, SQLParserON, SQLParserOR, SQLParserORDER, SQLParserREQUESTS, SQLParserREPLICATIONS, SQLParserRIGHT, SQLParserROLLUP, SQLParserSELECT, SQLParserSHOW, SQLParserSTATE, SQLParserSTORAGE, SQLParserTABLE_NAMES, SQLParserTIMESTAMP, SQLParserTRACE, SQLParserTRUE, SQLParserTYPE, SQLParserTYPES, SQLParserVALUES, SQLParserWHERE, SQLParserWITH, SQLParserWITHIN, SQLParserUSING, SQLParserUSE, SQLParserSECOND, SQLParserMINUTE, SQLParserHOUR, SQLParserDAY, SQLParserMONTH, SQLParserYEAR:
		localctx = NewUnquotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(682)
			p.NonReserved()
		}

	case SQLParserBACKQUOTED_IDENTIFIER:
		localctx = NewBackQuotedIdentifierContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(683)
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
			p.SetState(684)
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
	p.EnterRule(localctx, 112, SQLParserRULE_number)
	var _la int

	p.SetState(699)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 79, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDecimalLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(688)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserMINUS {
			{
				p.SetState(687)
				p.Match(SQLParserMINUS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(690)
			p.Match(SQLParserDECIMAL_VALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewDoubleLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(692)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserMINUS {
			{
				p.SetState(691)
				p.Match(SQLParserMINUS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(694)
			p.Match(SQLParserDOUBLE_VALUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewIntegerLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(696)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == SQLParserMINUS {
			{
				p.SetState(695)
				p.Match(SQLParserMINUS)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(698)
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

// IIntervalContext is an interface to support dynamic dispatch.
type IIntervalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetValue returns the value rule contexts.
	GetValue() INumberContext

	// GetUnit returns the unit rule contexts.
	GetUnit() IIntervalUnitContext

	// SetValue sets the value rule contexts.
	SetValue(INumberContext)

	// SetUnit sets the unit rule contexts.
	SetUnit(IIntervalUnitContext)

	// Getter signatures
	INTERVAL() antlr.TerminalNode
	Number() INumberContext
	IntervalUnit() IIntervalUnitContext

	// IsIntervalContext differentiates from other interfaces.
	IsIntervalContext()
}

type IntervalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	value  INumberContext
	unit   IIntervalUnitContext
}

func NewEmptyIntervalContext() *IntervalContext {
	var p = new(IntervalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_interval
	return p
}

func InitEmptyIntervalContext(p *IntervalContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_interval
}

func (*IntervalContext) IsIntervalContext() {}

func NewIntervalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalContext {
	var p = new(IntervalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_interval

	return p
}

func (s *IntervalContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalContext) GetValue() INumberContext { return s.value }

func (s *IntervalContext) GetUnit() IIntervalUnitContext { return s.unit }

func (s *IntervalContext) SetValue(v INumberContext) { s.value = v }

func (s *IntervalContext) SetUnit(v IIntervalUnitContext) { s.unit = v }

func (s *IntervalContext) INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserINTERVAL, 0)
}

func (s *IntervalContext) Number() INumberContext {
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

func (s *IntervalContext) IntervalUnit() IIntervalUnitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntervalUnitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntervalUnitContext)
}

func (s *IntervalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterInterval(s)
	}
}

func (s *IntervalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitInterval(s)
	}
}

func (s *IntervalContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitInterval(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) Interval() (localctx IIntervalContext) {
	localctx = NewIntervalContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 114, SQLParserRULE_interval)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(701)
		p.Match(SQLParserINTERVAL)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(702)

		var _x = p.Number()

		localctx.(*IntervalContext).value = _x
	}
	{
		p.SetState(703)

		var _x = p.IntervalUnit()

		localctx.(*IntervalContext).unit = _x
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

// IIntervalUnitContext is an interface to support dynamic dispatch.
type IIntervalUnitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SECOND() antlr.TerminalNode
	MINUTE() antlr.TerminalNode
	HOUR() antlr.TerminalNode
	DAY() antlr.TerminalNode
	MONTH() antlr.TerminalNode
	YEAR() antlr.TerminalNode

	// IsIntervalUnitContext differentiates from other interfaces.
	IsIntervalUnitContext()
}

type IntervalUnitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalUnitContext() *IntervalUnitContext {
	var p = new(IntervalUnitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_intervalUnit
	return p
}

func InitEmptyIntervalUnitContext(p *IntervalUnitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = SQLParserRULE_intervalUnit
}

func (*IntervalUnitContext) IsIntervalUnitContext() {}

func NewIntervalUnitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalUnitContext {
	var p = new(IntervalUnitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = SQLParserRULE_intervalUnit

	return p
}

func (s *IntervalUnitContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalUnitContext) SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserSECOND, 0)
}

func (s *IntervalUnitContext) MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserMINUTE, 0)
}

func (s *IntervalUnitContext) HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserHOUR, 0)
}

func (s *IntervalUnitContext) DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserDAY, 0)
}

func (s *IntervalUnitContext) MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserMONTH, 0)
}

func (s *IntervalUnitContext) YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserYEAR, 0)
}

func (s *IntervalUnitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalUnitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalUnitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.EnterIntervalUnit(s)
	}
}

func (s *IntervalUnitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SQLParserListener); ok {
		listenerT.ExitIntervalUnit(s)
	}
}

func (s *IntervalUnitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SQLParserVisitor:
		return t.VisitIntervalUnit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SQLParser) IntervalUnit() (localctx IIntervalUnitContext) {
	localctx = NewIntervalUnitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 116, SQLParserRULE_intervalUnit)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(705)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-82)) & ^0x3f) == 0 && ((int64(1)<<(_la-82))&63) != 0) {
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

// INonReservedContext is an interface to support dynamic dispatch.
type INonReservedContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	APP() antlr.TerminalNode
	ALL() antlr.TerminalNode
	ALIVE() antlr.TerminalNode
	AND() antlr.TerminalNode
	AS() antlr.TerminalNode
	ASC() antlr.TerminalNode
	BETWEEN() antlr.TerminalNode
	BROKER() antlr.TerminalNode
	BROKERS() antlr.TerminalNode
	BY() antlr.TerminalNode
	COMPACT() antlr.TerminalNode
	CREATE() antlr.TerminalNode
	CROSS() antlr.TerminalNode
	COLUMNS() antlr.TerminalNode
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
	FLUSH() antlr.TerminalNode
	FROM() antlr.TerminalNode
	GROUP() antlr.TerminalNode
	HAVING() antlr.TerminalNode
	JOIN() antlr.TerminalNode
	KEYS() antlr.TerminalNode
	LEFT() antlr.TerminalNode
	LIKE() antlr.TerminalNode
	LIMIT() antlr.TerminalNode
	LOG() antlr.TerminalNode
	LOGICAL() antlr.TerminalNode
	IF() antlr.TerminalNode
	IN() antlr.TerminalNode
	INTERVAL() antlr.TerminalNode
	IS() antlr.TerminalNode
	INSERT() antlr.TerminalNode
	INTO() antlr.TerminalNode
	MASTER() antlr.TerminalNode
	MEMORY_DATABASES() antlr.TerminalNode
	METRIC() antlr.TerminalNode
	METRICS() antlr.TerminalNode
	METADATA() antlr.TerminalNode
	METADATAS() antlr.TerminalNode
	NAMESPACE() antlr.TerminalNode
	NAMESPACES() antlr.TerminalNode
	NULL() antlr.TerminalNode
	NOT() antlr.TerminalNode
	NOW() antlr.TerminalNode
	ON() antlr.TerminalNode
	OR() antlr.TerminalNode
	ORDER() antlr.TerminalNode
	REQUESTS() antlr.TerminalNode
	REPLICATIONS() antlr.TerminalNode
	RIGHT() antlr.TerminalNode
	ROLLUP() antlr.TerminalNode
	SELECT() antlr.TerminalNode
	SHOW() antlr.TerminalNode
	STATE() antlr.TerminalNode
	STORAGE() antlr.TerminalNode
	TABLE_NAMES() antlr.TerminalNode
	TIMESTAMP() antlr.TerminalNode
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
	SECOND() antlr.TerminalNode
	MINUTE() antlr.TerminalNode
	HOUR() antlr.TerminalNode
	DAY() antlr.TerminalNode
	MONTH() antlr.TerminalNode
	YEAR() antlr.TerminalNode

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

func (s *NonReservedContext) APP() antlr.TerminalNode {
	return s.GetToken(SQLParserAPP, 0)
}

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

func (s *NonReservedContext) BETWEEN() antlr.TerminalNode {
	return s.GetToken(SQLParserBETWEEN, 0)
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

func (s *NonReservedContext) COLUMNS() antlr.TerminalNode {
	return s.GetToken(SQLParserCOLUMNS, 0)
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

func (s *NonReservedContext) IF() antlr.TerminalNode {
	return s.GetToken(SQLParserIF, 0)
}

func (s *NonReservedContext) IN() antlr.TerminalNode {
	return s.GetToken(SQLParserIN, 0)
}

func (s *NonReservedContext) INTERVAL() antlr.TerminalNode {
	return s.GetToken(SQLParserINTERVAL, 0)
}

func (s *NonReservedContext) IS() antlr.TerminalNode {
	return s.GetToken(SQLParserIS, 0)
}

func (s *NonReservedContext) INSERT() antlr.TerminalNode {
	return s.GetToken(SQLParserINSERT, 0)
}

func (s *NonReservedContext) INTO() antlr.TerminalNode {
	return s.GetToken(SQLParserINTO, 0)
}

func (s *NonReservedContext) MASTER() antlr.TerminalNode {
	return s.GetToken(SQLParserMASTER, 0)
}

func (s *NonReservedContext) MEMORY_DATABASES() antlr.TerminalNode {
	return s.GetToken(SQLParserMEMORY_DATABASES, 0)
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

func (s *NonReservedContext) NULL() antlr.TerminalNode {
	return s.GetToken(SQLParserNULL, 0)
}

func (s *NonReservedContext) NOT() antlr.TerminalNode {
	return s.GetToken(SQLParserNOT, 0)
}

func (s *NonReservedContext) NOW() antlr.TerminalNode {
	return s.GetToken(SQLParserNOW, 0)
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

func (s *NonReservedContext) TABLE_NAMES() antlr.TerminalNode {
	return s.GetToken(SQLParserTABLE_NAMES, 0)
}

func (s *NonReservedContext) TIMESTAMP() antlr.TerminalNode {
	return s.GetToken(SQLParserTIMESTAMP, 0)
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

func (s *NonReservedContext) SECOND() antlr.TerminalNode {
	return s.GetToken(SQLParserSECOND, 0)
}

func (s *NonReservedContext) MINUTE() antlr.TerminalNode {
	return s.GetToken(SQLParserMINUTE, 0)
}

func (s *NonReservedContext) HOUR() antlr.TerminalNode {
	return s.GetToken(SQLParserHOUR, 0)
}

func (s *NonReservedContext) DAY() antlr.TerminalNode {
	return s.GetToken(SQLParserDAY, 0)
}

func (s *NonReservedContext) MONTH() antlr.TerminalNode {
	return s.GetToken(SQLParserMONTH, 0)
}

func (s *NonReservedContext) YEAR() antlr.TerminalNode {
	return s.GetToken(SQLParserYEAR, 0)
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
	p.EnterRule(localctx, 118, SQLParserRULE_nonReserved)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(707)
		_la = p.GetTokenStream().LA(1)

		if !(((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-272) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&16777215) != 0)) {
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
	case 26:
		var t *RelationContext = nil
		if localctx != nil {
			t = localctx.(*RelationContext)
		}
		return p.Relation_Sempred(t, predIndex)

	case 39:
		var t *BooleanExpressionContext = nil
		if localctx != nil {
			t = localctx.(*BooleanExpressionContext)
		}
		return p.BooleanExpression_Sempred(t, predIndex)

	case 40:
		var t *ValueExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ValueExpressionContext)
		}
		return p.ValueExpression_Sempred(t, predIndex)

	case 41:
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
