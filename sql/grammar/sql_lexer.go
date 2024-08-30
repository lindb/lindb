// Code generated from ./sql/grammar/SQLLexer.g4 by ANTLR 4.13.2. DO NOT EDIT.

package grammar

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"sync"
	"unicode"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type SQLLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var SQLLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sqllexerLexerInit() {
	staticData := &SQLLexerLexerStaticData
	staticData.ChannelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN", "COMMENT",
	}
	staticData.ModeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.LiteralNames = []string{
		"", "", "", "", "'ALL'", "'ALIVE'", "'AND'", "'AS'", "'ASC'", "'BROKER'",
		"'BROKERS'", "'BY'", "'COMPACT'", "'CREATE'", "'CROSS'", "'DATABASE'",
		"'DATABASES'", "'DEFAULT'", "'DESC'", "'DROP'", "'ESCAPE'", "'EXPLAIN'",
		"'EXISTS'", "'FALSE'", "'FIELDS'", "'FILTER'", "'FLUSH'", "'FROM'",
		"'GROUP'", "'HAVING'", "'IF'", "'IN'", "'JOIN'", "'KEYS'", "'LEFT'",
		"'LIKE'", "'LIMIT'", "'MASTER'", "'METRICS'", "'METADATA'", "'METADATAS'",
		"'NAMESPACE'", "'NAMESPACES'", "'NOT'", "'ON'", "'OR'", "'ORDER'", "'PLAN'",
		"'REQUESTS'", "'REPLICATIONS'", "'RIGHT'", "'ROLLUP'", "'SELECT'", "'SHOW'",
		"'STATE'", "'STORAGE'", "'TAG'", "'TRUE'", "'TYPES'", "'VALUES'", "'WHERE'",
		"'WITH'", "'WITHIN'", "'USING'", "'USE'", "'='", "", "'<'", "'<='",
		"'>'", "'>='", "'+'", "'-'", "'*'", "'/'", "'%'", "'=~'", "'!~'", "'!'",
		"'.'", "'('", "')'", "','",
	}
	staticData.SymbolicNames = []string{
		"", "SIMPLE_COMMENT", "BRACKETED_COMMENT", "WS", "ALL", "ALIVE", "AND",
		"AS", "ASC", "BROKER", "BROKERS", "BY", "COMPACT", "CREATE", "CROSS",
		"DATABASE", "DATABASES", "DEFAULT", "DESC", "DROP", "ESCAPE", "EXPLAIN",
		"EXISTS", "FALSE", "FIELDS", "FILTER", "FLUSH", "FROM", "GROUP", "HAVING",
		"IF", "IN", "JOIN", "KEYS", "LEFT", "LIKE", "LIMIT", "MASTER", "METRICS",
		"METADATA", "METADATAS", "NAMESPACE", "NAMESPACES", "NOT", "ON", "OR",
		"ORDER", "PLAN", "REQUESTS", "REPLICATIONS", "RIGHT", "ROLLUP", "SELECT",
		"SHOW", "STATE", "STORAGE", "TAG", "TRUE", "TYPES", "VALUES", "WHERE",
		"WITH", "WITHIN", "USING", "USE", "EQ", "NEQ", "LT", "LTE", "GT", "GTE",
		"PLUS", "MINUS", "ASTERISK", "SLASH", "PERCENT", "REGEXP", "NEQREGEXP",
		"EXCLAMATION_SYMBOL", "DOT", "LR_BRACKET", "RR_BRACKET", "COMMA", "STRING",
		"INTEGER_VALUE", "DECIMAL_VALUE", "DOUBLE_VALUE", "IDENTIFIER", "DIGIT_IDENTIFIER",
		"QUOTED_IDENTIFIER", "BACKQUOTED_IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"SIMPLE_COMMENT", "BRACKETED_COMMENT", "WS", "ALL", "ALIVE", "AND",
		"AS", "ASC", "BROKER", "BROKERS", "BY", "COMPACT", "CREATE", "CROSS",
		"DATABASE", "DATABASES", "DEFAULT", "DESC", "DROP", "ESCAPE", "EXPLAIN",
		"EXISTS", "FALSE", "FIELDS", "FILTER", "FLUSH", "FROM", "GROUP", "HAVING",
		"IF", "IN", "JOIN", "KEYS", "LEFT", "LIKE", "LIMIT", "MASTER", "METRICS",
		"METADATA", "METADATAS", "NAMESPACE", "NAMESPACES", "NOT", "ON", "OR",
		"ORDER", "PLAN", "REQUESTS", "REPLICATIONS", "RIGHT", "ROLLUP", "SELECT",
		"SHOW", "STATE", "STORAGE", "TAG", "TRUE", "TYPES", "VALUES", "WHERE",
		"WITH", "WITHIN", "USING", "USE", "EQ", "NEQ", "LT", "LTE", "GT", "GTE",
		"PLUS", "MINUS", "ASTERISK", "SLASH", "PERCENT", "REGEXP", "NEQREGEXP",
		"EXCLAMATION_SYMBOL", "DOT", "LR_BRACKET", "RR_BRACKET", "COMMA", "STRING",
		"INTEGER_VALUE", "DECIMAL_VALUE", "DOUBLE_VALUE", "IDENTIFIER", "DIGIT_IDENTIFIER",
		"QUOTED_IDENTIFIER", "BACKQUOTED_IDENTIFIER", "DECIMAL_INTEGER", "EXPONENT",
		"DIGIT", "LETTER",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 90, 766, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15,
		7, 15, 2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7,
		20, 2, 21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25,
		2, 26, 7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2,
		31, 7, 31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36,
		7, 36, 2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7,
		41, 2, 42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46,
		2, 47, 7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2,
		52, 7, 52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 2, 57,
		7, 57, 2, 58, 7, 58, 2, 59, 7, 59, 2, 60, 7, 60, 2, 61, 7, 61, 2, 62, 7,
		62, 2, 63, 7, 63, 2, 64, 7, 64, 2, 65, 7, 65, 2, 66, 7, 66, 2, 67, 7, 67,
		2, 68, 7, 68, 2, 69, 7, 69, 2, 70, 7, 70, 2, 71, 7, 71, 2, 72, 7, 72, 2,
		73, 7, 73, 2, 74, 7, 74, 2, 75, 7, 75, 2, 76, 7, 76, 2, 77, 7, 77, 2, 78,
		7, 78, 2, 79, 7, 79, 2, 80, 7, 80, 2, 81, 7, 81, 2, 82, 7, 82, 2, 83, 7,
		83, 2, 84, 7, 84, 2, 85, 7, 85, 2, 86, 7, 86, 2, 87, 7, 87, 2, 88, 7, 88,
		2, 89, 7, 89, 2, 90, 7, 90, 2, 91, 7, 91, 2, 92, 7, 92, 2, 93, 7, 93, 1,
		0, 1, 0, 1, 0, 1, 0, 5, 0, 194, 8, 0, 10, 0, 12, 0, 197, 9, 0, 1, 0, 3,
		0, 200, 8, 0, 1, 0, 3, 0, 203, 8, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		5, 1, 211, 8, 1, 10, 1, 12, 1, 214, 9, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 2, 4, 2, 222, 8, 2, 11, 2, 12, 2, 223, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 6,
		1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8,
		1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 10, 1, 10, 1,
		10, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12,
		1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1,
		13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15,
		1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 16, 1,
		16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 17,
		1, 17, 1, 18, 1, 18, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 19, 1, 19, 1,
		19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20,
		1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 22, 1, 22, 1, 22, 1,
		22, 1, 22, 1, 22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 24,
		1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25, 1, 25, 1,
		25, 1, 25, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 27, 1, 27, 1, 27, 1, 27,
		1, 27, 1, 27, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 28, 1, 29, 1,
		29, 1, 29, 1, 30, 1, 30, 1, 30, 1, 31, 1, 31, 1, 31, 1, 31, 1, 31, 1, 32,
		1, 32, 1, 32, 1, 32, 1, 32, 1, 33, 1, 33, 1, 33, 1, 33, 1, 33, 1, 34, 1,
		34, 1, 34, 1, 34, 1, 34, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 36,
		1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 37, 1, 37, 1, 37, 1, 37, 1,
		37, 1, 37, 1, 37, 1, 37, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38,
		1, 38, 1, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1, 39, 1,
		39, 1, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40, 1, 40,
		1, 40, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1,
		41, 1, 41, 1, 42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43, 1, 44, 1, 44,
		1, 44, 1, 45, 1, 45, 1, 45, 1, 45, 1, 45, 1, 45, 1, 46, 1, 46, 1, 46, 1,
		46, 1, 46, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47, 1, 47,
		1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1,
		48, 1, 48, 1, 48, 1, 49, 1, 49, 1, 49, 1, 49, 1, 49, 1, 49, 1, 50, 1, 50,
		1, 50, 1, 50, 1, 50, 1, 50, 1, 50, 1, 51, 1, 51, 1, 51, 1, 51, 1, 51, 1,
		51, 1, 51, 1, 52, 1, 52, 1, 52, 1, 52, 1, 52, 1, 53, 1, 53, 1, 53, 1, 53,
		1, 53, 1, 53, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1,
		55, 1, 55, 1, 55, 1, 55, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56, 1, 57, 1, 57,
		1, 57, 1, 57, 1, 57, 1, 57, 1, 58, 1, 58, 1, 58, 1, 58, 1, 58, 1, 58, 1,
		58, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 59, 1, 60, 1, 60, 1, 60, 1, 60,
		1, 60, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 61, 1, 62, 1, 62, 1,
		62, 1, 62, 1, 62, 1, 62, 1, 63, 1, 63, 1, 63, 1, 63, 1, 64, 1, 64, 1, 65,
		1, 65, 1, 65, 1, 65, 3, 65, 616, 8, 65, 1, 66, 1, 66, 1, 67, 1, 67, 1,
		67, 1, 68, 1, 68, 1, 69, 1, 69, 1, 69, 1, 70, 1, 70, 1, 71, 1, 71, 1, 72,
		1, 72, 1, 73, 1, 73, 1, 74, 1, 74, 1, 75, 1, 75, 1, 75, 1, 76, 1, 76, 1,
		76, 1, 77, 1, 77, 1, 78, 1, 78, 1, 79, 1, 79, 1, 80, 1, 80, 1, 81, 1, 81,
		1, 82, 1, 82, 1, 82, 1, 82, 5, 82, 658, 8, 82, 10, 82, 12, 82, 661, 9,
		82, 1, 82, 1, 82, 1, 83, 1, 83, 1, 84, 1, 84, 1, 84, 3, 84, 670, 8, 84,
		1, 84, 1, 84, 3, 84, 674, 8, 84, 1, 85, 4, 85, 677, 8, 85, 11, 85, 12,
		85, 678, 1, 85, 1, 85, 5, 85, 683, 8, 85, 10, 85, 12, 85, 686, 9, 85, 3,
		85, 688, 8, 85, 1, 85, 1, 85, 1, 85, 1, 85, 4, 85, 694, 8, 85, 11, 85,
		12, 85, 695, 1, 85, 1, 85, 3, 85, 700, 8, 85, 1, 86, 1, 86, 3, 86, 704,
		8, 86, 1, 86, 1, 86, 1, 86, 5, 86, 709, 8, 86, 10, 86, 12, 86, 712, 9,
		86, 1, 87, 1, 87, 1, 87, 1, 87, 4, 87, 718, 8, 87, 11, 87, 12, 87, 719,
		1, 88, 1, 88, 1, 88, 1, 88, 5, 88, 726, 8, 88, 10, 88, 12, 88, 729, 9,
		88, 1, 88, 1, 88, 1, 89, 1, 89, 1, 89, 1, 89, 5, 89, 737, 8, 89, 10, 89,
		12, 89, 740, 9, 89, 1, 89, 1, 89, 1, 90, 1, 90, 3, 90, 746, 8, 90, 1, 90,
		5, 90, 749, 8, 90, 10, 90, 12, 90, 752, 9, 90, 1, 91, 1, 91, 3, 91, 756,
		8, 91, 1, 91, 4, 91, 759, 8, 91, 11, 91, 12, 91, 760, 1, 92, 1, 92, 1,
		93, 1, 93, 1, 212, 0, 94, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15,
		8, 17, 9, 19, 10, 21, 11, 23, 12, 25, 13, 27, 14, 29, 15, 31, 16, 33, 17,
		35, 18, 37, 19, 39, 20, 41, 21, 43, 22, 45, 23, 47, 24, 49, 25, 51, 26,
		53, 27, 55, 28, 57, 29, 59, 30, 61, 31, 63, 32, 65, 33, 67, 34, 69, 35,
		71, 36, 73, 37, 75, 38, 77, 39, 79, 40, 81, 41, 83, 42, 85, 43, 87, 44,
		89, 45, 91, 46, 93, 47, 95, 48, 97, 49, 99, 50, 101, 51, 103, 52, 105,
		53, 107, 54, 109, 55, 111, 56, 113, 57, 115, 58, 117, 59, 119, 60, 121,
		61, 123, 62, 125, 63, 127, 64, 129, 65, 131, 66, 133, 67, 135, 68, 137,
		69, 139, 70, 141, 71, 143, 72, 145, 73, 147, 74, 149, 75, 151, 76, 153,
		77, 155, 78, 157, 79, 159, 80, 161, 81, 163, 82, 165, 83, 167, 84, 169,
		85, 171, 86, 173, 87, 175, 88, 177, 89, 179, 90, 181, 0, 183, 0, 185, 0,
		187, 0, 1, 0, 33, 2, 0, 10, 10, 13, 13, 3, 0, 9, 10, 13, 13, 32, 32, 2,
		0, 65, 65, 97, 97, 2, 0, 76, 76, 108, 108, 2, 0, 73, 73, 105, 105, 2, 0,
		86, 86, 118, 118, 2, 0, 69, 69, 101, 101, 2, 0, 78, 78, 110, 110, 2, 0,
		68, 68, 100, 100, 2, 0, 83, 83, 115, 115, 2, 0, 67, 67, 99, 99, 2, 0, 66,
		66, 98, 98, 2, 0, 82, 82, 114, 114, 2, 0, 79, 79, 111, 111, 2, 0, 75, 75,
		107, 107, 2, 0, 89, 89, 121, 121, 2, 0, 77, 77, 109, 109, 2, 0, 80, 80,
		112, 112, 2, 0, 84, 84, 116, 116, 2, 0, 70, 70, 102, 102, 2, 0, 85, 85,
		117, 117, 2, 0, 88, 88, 120, 120, 2, 0, 72, 72, 104, 104, 2, 0, 71, 71,
		103, 103, 2, 0, 74, 74, 106, 106, 2, 0, 81, 81, 113, 113, 2, 0, 87, 87,
		119, 119, 1, 0, 39, 39, 1, 0, 34, 34, 1, 0, 96, 96, 2, 0, 43, 43, 45, 45,
		1, 0, 48, 57, 2, 0, 65, 90, 97, 122, 791, 0, 1, 1, 0, 0, 0, 0, 3, 1, 0,
		0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0,
		0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1,
		0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 0, 25, 1, 0, 0, 0, 0, 27,
		1, 0, 0, 0, 0, 29, 1, 0, 0, 0, 0, 31, 1, 0, 0, 0, 0, 33, 1, 0, 0, 0, 0,
		35, 1, 0, 0, 0, 0, 37, 1, 0, 0, 0, 0, 39, 1, 0, 0, 0, 0, 41, 1, 0, 0, 0,
		0, 43, 1, 0, 0, 0, 0, 45, 1, 0, 0, 0, 0, 47, 1, 0, 0, 0, 0, 49, 1, 0, 0,
		0, 0, 51, 1, 0, 0, 0, 0, 53, 1, 0, 0, 0, 0, 55, 1, 0, 0, 0, 0, 57, 1, 0,
		0, 0, 0, 59, 1, 0, 0, 0, 0, 61, 1, 0, 0, 0, 0, 63, 1, 0, 0, 0, 0, 65, 1,
		0, 0, 0, 0, 67, 1, 0, 0, 0, 0, 69, 1, 0, 0, 0, 0, 71, 1, 0, 0, 0, 0, 73,
		1, 0, 0, 0, 0, 75, 1, 0, 0, 0, 0, 77, 1, 0, 0, 0, 0, 79, 1, 0, 0, 0, 0,
		81, 1, 0, 0, 0, 0, 83, 1, 0, 0, 0, 0, 85, 1, 0, 0, 0, 0, 87, 1, 0, 0, 0,
		0, 89, 1, 0, 0, 0, 0, 91, 1, 0, 0, 0, 0, 93, 1, 0, 0, 0, 0, 95, 1, 0, 0,
		0, 0, 97, 1, 0, 0, 0, 0, 99, 1, 0, 0, 0, 0, 101, 1, 0, 0, 0, 0, 103, 1,
		0, 0, 0, 0, 105, 1, 0, 0, 0, 0, 107, 1, 0, 0, 0, 0, 109, 1, 0, 0, 0, 0,
		111, 1, 0, 0, 0, 0, 113, 1, 0, 0, 0, 0, 115, 1, 0, 0, 0, 0, 117, 1, 0,
		0, 0, 0, 119, 1, 0, 0, 0, 0, 121, 1, 0, 0, 0, 0, 123, 1, 0, 0, 0, 0, 125,
		1, 0, 0, 0, 0, 127, 1, 0, 0, 0, 0, 129, 1, 0, 0, 0, 0, 131, 1, 0, 0, 0,
		0, 133, 1, 0, 0, 0, 0, 135, 1, 0, 0, 0, 0, 137, 1, 0, 0, 0, 0, 139, 1,
		0, 0, 0, 0, 141, 1, 0, 0, 0, 0, 143, 1, 0, 0, 0, 0, 145, 1, 0, 0, 0, 0,
		147, 1, 0, 0, 0, 0, 149, 1, 0, 0, 0, 0, 151, 1, 0, 0, 0, 0, 153, 1, 0,
		0, 0, 0, 155, 1, 0, 0, 0, 0, 157, 1, 0, 0, 0, 0, 159, 1, 0, 0, 0, 0, 161,
		1, 0, 0, 0, 0, 163, 1, 0, 0, 0, 0, 165, 1, 0, 0, 0, 0, 167, 1, 0, 0, 0,
		0, 169, 1, 0, 0, 0, 0, 171, 1, 0, 0, 0, 0, 173, 1, 0, 0, 0, 0, 175, 1,
		0, 0, 0, 0, 177, 1, 0, 0, 0, 0, 179, 1, 0, 0, 0, 1, 189, 1, 0, 0, 0, 3,
		206, 1, 0, 0, 0, 5, 221, 1, 0, 0, 0, 7, 227, 1, 0, 0, 0, 9, 231, 1, 0,
		0, 0, 11, 237, 1, 0, 0, 0, 13, 241, 1, 0, 0, 0, 15, 244, 1, 0, 0, 0, 17,
		248, 1, 0, 0, 0, 19, 255, 1, 0, 0, 0, 21, 263, 1, 0, 0, 0, 23, 266, 1,
		0, 0, 0, 25, 274, 1, 0, 0, 0, 27, 281, 1, 0, 0, 0, 29, 287, 1, 0, 0, 0,
		31, 296, 1, 0, 0, 0, 33, 306, 1, 0, 0, 0, 35, 314, 1, 0, 0, 0, 37, 319,
		1, 0, 0, 0, 39, 324, 1, 0, 0, 0, 41, 331, 1, 0, 0, 0, 43, 339, 1, 0, 0,
		0, 45, 346, 1, 0, 0, 0, 47, 352, 1, 0, 0, 0, 49, 359, 1, 0, 0, 0, 51, 366,
		1, 0, 0, 0, 53, 372, 1, 0, 0, 0, 55, 377, 1, 0, 0, 0, 57, 383, 1, 0, 0,
		0, 59, 390, 1, 0, 0, 0, 61, 393, 1, 0, 0, 0, 63, 396, 1, 0, 0, 0, 65, 401,
		1, 0, 0, 0, 67, 406, 1, 0, 0, 0, 69, 411, 1, 0, 0, 0, 71, 416, 1, 0, 0,
		0, 73, 422, 1, 0, 0, 0, 75, 429, 1, 0, 0, 0, 77, 437, 1, 0, 0, 0, 79, 446,
		1, 0, 0, 0, 81, 456, 1, 0, 0, 0, 83, 466, 1, 0, 0, 0, 85, 477, 1, 0, 0,
		0, 87, 481, 1, 0, 0, 0, 89, 484, 1, 0, 0, 0, 91, 487, 1, 0, 0, 0, 93, 493,
		1, 0, 0, 0, 95, 498, 1, 0, 0, 0, 97, 507, 1, 0, 0, 0, 99, 520, 1, 0, 0,
		0, 101, 526, 1, 0, 0, 0, 103, 533, 1, 0, 0, 0, 105, 540, 1, 0, 0, 0, 107,
		545, 1, 0, 0, 0, 109, 551, 1, 0, 0, 0, 111, 559, 1, 0, 0, 0, 113, 563,
		1, 0, 0, 0, 115, 568, 1, 0, 0, 0, 117, 574, 1, 0, 0, 0, 119, 581, 1, 0,
		0, 0, 121, 587, 1, 0, 0, 0, 123, 592, 1, 0, 0, 0, 125, 599, 1, 0, 0, 0,
		127, 605, 1, 0, 0, 0, 129, 609, 1, 0, 0, 0, 131, 615, 1, 0, 0, 0, 133,
		617, 1, 0, 0, 0, 135, 619, 1, 0, 0, 0, 137, 622, 1, 0, 0, 0, 139, 624,
		1, 0, 0, 0, 141, 627, 1, 0, 0, 0, 143, 629, 1, 0, 0, 0, 145, 631, 1, 0,
		0, 0, 147, 633, 1, 0, 0, 0, 149, 635, 1, 0, 0, 0, 151, 637, 1, 0, 0, 0,
		153, 640, 1, 0, 0, 0, 155, 643, 1, 0, 0, 0, 157, 645, 1, 0, 0, 0, 159,
		647, 1, 0, 0, 0, 161, 649, 1, 0, 0, 0, 163, 651, 1, 0, 0, 0, 165, 653,
		1, 0, 0, 0, 167, 664, 1, 0, 0, 0, 169, 673, 1, 0, 0, 0, 171, 699, 1, 0,
		0, 0, 173, 703, 1, 0, 0, 0, 175, 713, 1, 0, 0, 0, 177, 721, 1, 0, 0, 0,
		179, 732, 1, 0, 0, 0, 181, 743, 1, 0, 0, 0, 183, 753, 1, 0, 0, 0, 185,
		762, 1, 0, 0, 0, 187, 764, 1, 0, 0, 0, 189, 190, 5, 45, 0, 0, 190, 191,
		5, 45, 0, 0, 191, 195, 1, 0, 0, 0, 192, 194, 8, 0, 0, 0, 193, 192, 1, 0,
		0, 0, 194, 197, 1, 0, 0, 0, 195, 193, 1, 0, 0, 0, 195, 196, 1, 0, 0, 0,
		196, 199, 1, 0, 0, 0, 197, 195, 1, 0, 0, 0, 198, 200, 5, 13, 0, 0, 199,
		198, 1, 0, 0, 0, 199, 200, 1, 0, 0, 0, 200, 202, 1, 0, 0, 0, 201, 203,
		5, 10, 0, 0, 202, 201, 1, 0, 0, 0, 202, 203, 1, 0, 0, 0, 203, 204, 1, 0,
		0, 0, 204, 205, 6, 0, 0, 0, 205, 2, 1, 0, 0, 0, 206, 207, 5, 47, 0, 0,
		207, 208, 5, 42, 0, 0, 208, 212, 1, 0, 0, 0, 209, 211, 9, 0, 0, 0, 210,
		209, 1, 0, 0, 0, 211, 214, 1, 0, 0, 0, 212, 213, 1, 0, 0, 0, 212, 210,
		1, 0, 0, 0, 213, 215, 1, 0, 0, 0, 214, 212, 1, 0, 0, 0, 215, 216, 5, 42,
		0, 0, 216, 217, 5, 47, 0, 0, 217, 218, 1, 0, 0, 0, 218, 219, 6, 1, 0, 0,
		219, 4, 1, 0, 0, 0, 220, 222, 7, 1, 0, 0, 221, 220, 1, 0, 0, 0, 222, 223,
		1, 0, 0, 0, 223, 221, 1, 0, 0, 0, 223, 224, 1, 0, 0, 0, 224, 225, 1, 0,
		0, 0, 225, 226, 6, 2, 1, 0, 226, 6, 1, 0, 0, 0, 227, 228, 7, 2, 0, 0, 228,
		229, 7, 3, 0, 0, 229, 230, 7, 3, 0, 0, 230, 8, 1, 0, 0, 0, 231, 232, 7,
		2, 0, 0, 232, 233, 7, 3, 0, 0, 233, 234, 7, 4, 0, 0, 234, 235, 7, 5, 0,
		0, 235, 236, 7, 6, 0, 0, 236, 10, 1, 0, 0, 0, 237, 238, 7, 2, 0, 0, 238,
		239, 7, 7, 0, 0, 239, 240, 7, 8, 0, 0, 240, 12, 1, 0, 0, 0, 241, 242, 7,
		2, 0, 0, 242, 243, 7, 9, 0, 0, 243, 14, 1, 0, 0, 0, 244, 245, 7, 2, 0,
		0, 245, 246, 7, 9, 0, 0, 246, 247, 7, 10, 0, 0, 247, 16, 1, 0, 0, 0, 248,
		249, 7, 11, 0, 0, 249, 250, 7, 12, 0, 0, 250, 251, 7, 13, 0, 0, 251, 252,
		7, 14, 0, 0, 252, 253, 7, 6, 0, 0, 253, 254, 7, 12, 0, 0, 254, 18, 1, 0,
		0, 0, 255, 256, 7, 11, 0, 0, 256, 257, 7, 12, 0, 0, 257, 258, 7, 13, 0,
		0, 258, 259, 7, 14, 0, 0, 259, 260, 7, 6, 0, 0, 260, 261, 7, 12, 0, 0,
		261, 262, 7, 9, 0, 0, 262, 20, 1, 0, 0, 0, 263, 264, 7, 11, 0, 0, 264,
		265, 7, 15, 0, 0, 265, 22, 1, 0, 0, 0, 266, 267, 7, 10, 0, 0, 267, 268,
		7, 13, 0, 0, 268, 269, 7, 16, 0, 0, 269, 270, 7, 17, 0, 0, 270, 271, 7,
		2, 0, 0, 271, 272, 7, 10, 0, 0, 272, 273, 7, 18, 0, 0, 273, 24, 1, 0, 0,
		0, 274, 275, 7, 10, 0, 0, 275, 276, 7, 12, 0, 0, 276, 277, 7, 6, 0, 0,
		277, 278, 7, 2, 0, 0, 278, 279, 7, 18, 0, 0, 279, 280, 7, 6, 0, 0, 280,
		26, 1, 0, 0, 0, 281, 282, 7, 10, 0, 0, 282, 283, 7, 12, 0, 0, 283, 284,
		7, 13, 0, 0, 284, 285, 7, 9, 0, 0, 285, 286, 7, 9, 0, 0, 286, 28, 1, 0,
		0, 0, 287, 288, 7, 8, 0, 0, 288, 289, 7, 2, 0, 0, 289, 290, 7, 18, 0, 0,
		290, 291, 7, 2, 0, 0, 291, 292, 7, 11, 0, 0, 292, 293, 7, 2, 0, 0, 293,
		294, 7, 9, 0, 0, 294, 295, 7, 6, 0, 0, 295, 30, 1, 0, 0, 0, 296, 297, 7,
		8, 0, 0, 297, 298, 7, 2, 0, 0, 298, 299, 7, 18, 0, 0, 299, 300, 7, 2, 0,
		0, 300, 301, 7, 11, 0, 0, 301, 302, 7, 2, 0, 0, 302, 303, 7, 9, 0, 0, 303,
		304, 7, 6, 0, 0, 304, 305, 7, 9, 0, 0, 305, 32, 1, 0, 0, 0, 306, 307, 7,
		8, 0, 0, 307, 308, 7, 6, 0, 0, 308, 309, 7, 19, 0, 0, 309, 310, 7, 2, 0,
		0, 310, 311, 7, 20, 0, 0, 311, 312, 7, 3, 0, 0, 312, 313, 7, 18, 0, 0,
		313, 34, 1, 0, 0, 0, 314, 315, 7, 8, 0, 0, 315, 316, 7, 6, 0, 0, 316, 317,
		7, 9, 0, 0, 317, 318, 7, 10, 0, 0, 318, 36, 1, 0, 0, 0, 319, 320, 7, 8,
		0, 0, 320, 321, 7, 12, 0, 0, 321, 322, 7, 13, 0, 0, 322, 323, 7, 17, 0,
		0, 323, 38, 1, 0, 0, 0, 324, 325, 7, 6, 0, 0, 325, 326, 7, 9, 0, 0, 326,
		327, 7, 10, 0, 0, 327, 328, 7, 2, 0, 0, 328, 329, 7, 17, 0, 0, 329, 330,
		7, 6, 0, 0, 330, 40, 1, 0, 0, 0, 331, 332, 7, 6, 0, 0, 332, 333, 7, 21,
		0, 0, 333, 334, 7, 17, 0, 0, 334, 335, 7, 3, 0, 0, 335, 336, 7, 2, 0, 0,
		336, 337, 7, 4, 0, 0, 337, 338, 7, 7, 0, 0, 338, 42, 1, 0, 0, 0, 339, 340,
		7, 6, 0, 0, 340, 341, 7, 21, 0, 0, 341, 342, 7, 4, 0, 0, 342, 343, 7, 9,
		0, 0, 343, 344, 7, 18, 0, 0, 344, 345, 7, 9, 0, 0, 345, 44, 1, 0, 0, 0,
		346, 347, 7, 19, 0, 0, 347, 348, 7, 2, 0, 0, 348, 349, 7, 3, 0, 0, 349,
		350, 7, 9, 0, 0, 350, 351, 7, 6, 0, 0, 351, 46, 1, 0, 0, 0, 352, 353, 7,
		19, 0, 0, 353, 354, 7, 4, 0, 0, 354, 355, 7, 6, 0, 0, 355, 356, 7, 3, 0,
		0, 356, 357, 7, 8, 0, 0, 357, 358, 7, 9, 0, 0, 358, 48, 1, 0, 0, 0, 359,
		360, 7, 19, 0, 0, 360, 361, 7, 4, 0, 0, 361, 362, 7, 3, 0, 0, 362, 363,
		7, 18, 0, 0, 363, 364, 7, 6, 0, 0, 364, 365, 7, 12, 0, 0, 365, 50, 1, 0,
		0, 0, 366, 367, 7, 19, 0, 0, 367, 368, 7, 3, 0, 0, 368, 369, 7, 20, 0,
		0, 369, 370, 7, 9, 0, 0, 370, 371, 7, 22, 0, 0, 371, 52, 1, 0, 0, 0, 372,
		373, 7, 19, 0, 0, 373, 374, 7, 12, 0, 0, 374, 375, 7, 13, 0, 0, 375, 376,
		7, 16, 0, 0, 376, 54, 1, 0, 0, 0, 377, 378, 7, 23, 0, 0, 378, 379, 7, 12,
		0, 0, 379, 380, 7, 13, 0, 0, 380, 381, 7, 20, 0, 0, 381, 382, 7, 17, 0,
		0, 382, 56, 1, 0, 0, 0, 383, 384, 7, 22, 0, 0, 384, 385, 7, 2, 0, 0, 385,
		386, 7, 5, 0, 0, 386, 387, 7, 4, 0, 0, 387, 388, 7, 7, 0, 0, 388, 389,
		7, 23, 0, 0, 389, 58, 1, 0, 0, 0, 390, 391, 7, 4, 0, 0, 391, 392, 7, 19,
		0, 0, 392, 60, 1, 0, 0, 0, 393, 394, 7, 4, 0, 0, 394, 395, 7, 7, 0, 0,
		395, 62, 1, 0, 0, 0, 396, 397, 7, 24, 0, 0, 397, 398, 7, 13, 0, 0, 398,
		399, 7, 4, 0, 0, 399, 400, 7, 7, 0, 0, 400, 64, 1, 0, 0, 0, 401, 402, 7,
		14, 0, 0, 402, 403, 7, 6, 0, 0, 403, 404, 7, 15, 0, 0, 404, 405, 7, 9,
		0, 0, 405, 66, 1, 0, 0, 0, 406, 407, 7, 3, 0, 0, 407, 408, 7, 6, 0, 0,
		408, 409, 7, 19, 0, 0, 409, 410, 7, 18, 0, 0, 410, 68, 1, 0, 0, 0, 411,
		412, 7, 3, 0, 0, 412, 413, 7, 4, 0, 0, 413, 414, 7, 14, 0, 0, 414, 415,
		7, 6, 0, 0, 415, 70, 1, 0, 0, 0, 416, 417, 7, 3, 0, 0, 417, 418, 7, 4,
		0, 0, 418, 419, 7, 16, 0, 0, 419, 420, 7, 4, 0, 0, 420, 421, 7, 18, 0,
		0, 421, 72, 1, 0, 0, 0, 422, 423, 7, 16, 0, 0, 423, 424, 7, 2, 0, 0, 424,
		425, 7, 9, 0, 0, 425, 426, 7, 18, 0, 0, 426, 427, 7, 6, 0, 0, 427, 428,
		7, 12, 0, 0, 428, 74, 1, 0, 0, 0, 429, 430, 7, 16, 0, 0, 430, 431, 7, 6,
		0, 0, 431, 432, 7, 18, 0, 0, 432, 433, 7, 12, 0, 0, 433, 434, 7, 4, 0,
		0, 434, 435, 7, 10, 0, 0, 435, 436, 7, 9, 0, 0, 436, 76, 1, 0, 0, 0, 437,
		438, 7, 16, 0, 0, 438, 439, 7, 6, 0, 0, 439, 440, 7, 18, 0, 0, 440, 441,
		7, 2, 0, 0, 441, 442, 7, 8, 0, 0, 442, 443, 7, 2, 0, 0, 443, 444, 7, 18,
		0, 0, 444, 445, 7, 2, 0, 0, 445, 78, 1, 0, 0, 0, 446, 447, 7, 16, 0, 0,
		447, 448, 7, 6, 0, 0, 448, 449, 7, 18, 0, 0, 449, 450, 7, 2, 0, 0, 450,
		451, 7, 8, 0, 0, 451, 452, 7, 2, 0, 0, 452, 453, 7, 18, 0, 0, 453, 454,
		7, 2, 0, 0, 454, 455, 7, 9, 0, 0, 455, 80, 1, 0, 0, 0, 456, 457, 7, 7,
		0, 0, 457, 458, 7, 2, 0, 0, 458, 459, 7, 16, 0, 0, 459, 460, 7, 6, 0, 0,
		460, 461, 7, 9, 0, 0, 461, 462, 7, 17, 0, 0, 462, 463, 7, 2, 0, 0, 463,
		464, 7, 10, 0, 0, 464, 465, 7, 6, 0, 0, 465, 82, 1, 0, 0, 0, 466, 467,
		7, 7, 0, 0, 467, 468, 7, 2, 0, 0, 468, 469, 7, 16, 0, 0, 469, 470, 7, 6,
		0, 0, 470, 471, 7, 9, 0, 0, 471, 472, 7, 17, 0, 0, 472, 473, 7, 2, 0, 0,
		473, 474, 7, 10, 0, 0, 474, 475, 7, 6, 0, 0, 475, 476, 7, 9, 0, 0, 476,
		84, 1, 0, 0, 0, 477, 478, 7, 7, 0, 0, 478, 479, 7, 13, 0, 0, 479, 480,
		7, 18, 0, 0, 480, 86, 1, 0, 0, 0, 481, 482, 7, 13, 0, 0, 482, 483, 7, 7,
		0, 0, 483, 88, 1, 0, 0, 0, 484, 485, 7, 13, 0, 0, 485, 486, 7, 12, 0, 0,
		486, 90, 1, 0, 0, 0, 487, 488, 7, 13, 0, 0, 488, 489, 7, 12, 0, 0, 489,
		490, 7, 8, 0, 0, 490, 491, 7, 6, 0, 0, 491, 492, 7, 12, 0, 0, 492, 92,
		1, 0, 0, 0, 493, 494, 7, 17, 0, 0, 494, 495, 7, 3, 0, 0, 495, 496, 7, 2,
		0, 0, 496, 497, 7, 7, 0, 0, 497, 94, 1, 0, 0, 0, 498, 499, 7, 12, 0, 0,
		499, 500, 7, 6, 0, 0, 500, 501, 7, 25, 0, 0, 501, 502, 7, 20, 0, 0, 502,
		503, 7, 6, 0, 0, 503, 504, 7, 9, 0, 0, 504, 505, 7, 18, 0, 0, 505, 506,
		7, 9, 0, 0, 506, 96, 1, 0, 0, 0, 507, 508, 7, 12, 0, 0, 508, 509, 7, 6,
		0, 0, 509, 510, 7, 17, 0, 0, 510, 511, 7, 3, 0, 0, 511, 512, 7, 4, 0, 0,
		512, 513, 7, 10, 0, 0, 513, 514, 7, 2, 0, 0, 514, 515, 7, 18, 0, 0, 515,
		516, 7, 4, 0, 0, 516, 517, 7, 13, 0, 0, 517, 518, 7, 7, 0, 0, 518, 519,
		7, 9, 0, 0, 519, 98, 1, 0, 0, 0, 520, 521, 7, 12, 0, 0, 521, 522, 7, 4,
		0, 0, 522, 523, 7, 23, 0, 0, 523, 524, 7, 22, 0, 0, 524, 525, 7, 18, 0,
		0, 525, 100, 1, 0, 0, 0, 526, 527, 7, 12, 0, 0, 527, 528, 7, 13, 0, 0,
		528, 529, 7, 3, 0, 0, 529, 530, 7, 3, 0, 0, 530, 531, 7, 20, 0, 0, 531,
		532, 7, 17, 0, 0, 532, 102, 1, 0, 0, 0, 533, 534, 7, 9, 0, 0, 534, 535,
		7, 6, 0, 0, 535, 536, 7, 3, 0, 0, 536, 537, 7, 6, 0, 0, 537, 538, 7, 10,
		0, 0, 538, 539, 7, 18, 0, 0, 539, 104, 1, 0, 0, 0, 540, 541, 7, 9, 0, 0,
		541, 542, 7, 22, 0, 0, 542, 543, 7, 13, 0, 0, 543, 544, 7, 26, 0, 0, 544,
		106, 1, 0, 0, 0, 545, 546, 7, 9, 0, 0, 546, 547, 7, 18, 0, 0, 547, 548,
		7, 2, 0, 0, 548, 549, 7, 18, 0, 0, 549, 550, 7, 6, 0, 0, 550, 108, 1, 0,
		0, 0, 551, 552, 7, 9, 0, 0, 552, 553, 7, 18, 0, 0, 553, 554, 7, 13, 0,
		0, 554, 555, 7, 12, 0, 0, 555, 556, 7, 2, 0, 0, 556, 557, 7, 23, 0, 0,
		557, 558, 7, 6, 0, 0, 558, 110, 1, 0, 0, 0, 559, 560, 7, 18, 0, 0, 560,
		561, 7, 2, 0, 0, 561, 562, 7, 23, 0, 0, 562, 112, 1, 0, 0, 0, 563, 564,
		7, 18, 0, 0, 564, 565, 7, 12, 0, 0, 565, 566, 7, 20, 0, 0, 566, 567, 7,
		6, 0, 0, 567, 114, 1, 0, 0, 0, 568, 569, 7, 18, 0, 0, 569, 570, 7, 15,
		0, 0, 570, 571, 7, 17, 0, 0, 571, 572, 7, 6, 0, 0, 572, 573, 7, 9, 0, 0,
		573, 116, 1, 0, 0, 0, 574, 575, 7, 5, 0, 0, 575, 576, 7, 2, 0, 0, 576,
		577, 7, 3, 0, 0, 577, 578, 7, 20, 0, 0, 578, 579, 7, 6, 0, 0, 579, 580,
		7, 9, 0, 0, 580, 118, 1, 0, 0, 0, 581, 582, 7, 26, 0, 0, 582, 583, 7, 22,
		0, 0, 583, 584, 7, 6, 0, 0, 584, 585, 7, 12, 0, 0, 585, 586, 7, 6, 0, 0,
		586, 120, 1, 0, 0, 0, 587, 588, 7, 26, 0, 0, 588, 589, 7, 4, 0, 0, 589,
		590, 7, 18, 0, 0, 590, 591, 7, 22, 0, 0, 591, 122, 1, 0, 0, 0, 592, 593,
		7, 26, 0, 0, 593, 594, 7, 4, 0, 0, 594, 595, 7, 18, 0, 0, 595, 596, 7,
		22, 0, 0, 596, 597, 7, 4, 0, 0, 597, 598, 7, 7, 0, 0, 598, 124, 1, 0, 0,
		0, 599, 600, 7, 20, 0, 0, 600, 601, 7, 9, 0, 0, 601, 602, 7, 4, 0, 0, 602,
		603, 7, 7, 0, 0, 603, 604, 7, 23, 0, 0, 604, 126, 1, 0, 0, 0, 605, 606,
		7, 20, 0, 0, 606, 607, 7, 9, 0, 0, 607, 608, 7, 6, 0, 0, 608, 128, 1, 0,
		0, 0, 609, 610, 5, 61, 0, 0, 610, 130, 1, 0, 0, 0, 611, 612, 5, 60, 0,
		0, 612, 616, 5, 62, 0, 0, 613, 614, 5, 33, 0, 0, 614, 616, 5, 61, 0, 0,
		615, 611, 1, 0, 0, 0, 615, 613, 1, 0, 0, 0, 616, 132, 1, 0, 0, 0, 617,
		618, 5, 60, 0, 0, 618, 134, 1, 0, 0, 0, 619, 620, 5, 60, 0, 0, 620, 621,
		5, 61, 0, 0, 621, 136, 1, 0, 0, 0, 622, 623, 5, 62, 0, 0, 623, 138, 1,
		0, 0, 0, 624, 625, 5, 62, 0, 0, 625, 626, 5, 61, 0, 0, 626, 140, 1, 0,
		0, 0, 627, 628, 5, 43, 0, 0, 628, 142, 1, 0, 0, 0, 629, 630, 5, 45, 0,
		0, 630, 144, 1, 0, 0, 0, 631, 632, 5, 42, 0, 0, 632, 146, 1, 0, 0, 0, 633,
		634, 5, 47, 0, 0, 634, 148, 1, 0, 0, 0, 635, 636, 5, 37, 0, 0, 636, 150,
		1, 0, 0, 0, 637, 638, 5, 61, 0, 0, 638, 639, 5, 126, 0, 0, 639, 152, 1,
		0, 0, 0, 640, 641, 5, 33, 0, 0, 641, 642, 5, 126, 0, 0, 642, 154, 1, 0,
		0, 0, 643, 644, 5, 33, 0, 0, 644, 156, 1, 0, 0, 0, 645, 646, 5, 46, 0,
		0, 646, 158, 1, 0, 0, 0, 647, 648, 5, 40, 0, 0, 648, 160, 1, 0, 0, 0, 649,
		650, 5, 41, 0, 0, 650, 162, 1, 0, 0, 0, 651, 652, 5, 44, 0, 0, 652, 164,
		1, 0, 0, 0, 653, 659, 5, 39, 0, 0, 654, 658, 8, 27, 0, 0, 655, 656, 5,
		39, 0, 0, 656, 658, 5, 39, 0, 0, 657, 654, 1, 0, 0, 0, 657, 655, 1, 0,
		0, 0, 658, 661, 1, 0, 0, 0, 659, 657, 1, 0, 0, 0, 659, 660, 1, 0, 0, 0,
		660, 662, 1, 0, 0, 0, 661, 659, 1, 0, 0, 0, 662, 663, 5, 39, 0, 0, 663,
		166, 1, 0, 0, 0, 664, 665, 3, 181, 90, 0, 665, 168, 1, 0, 0, 0, 666, 667,
		3, 181, 90, 0, 667, 669, 5, 46, 0, 0, 668, 670, 3, 181, 90, 0, 669, 668,
		1, 0, 0, 0, 669, 670, 1, 0, 0, 0, 670, 674, 1, 0, 0, 0, 671, 672, 5, 46,
		0, 0, 672, 674, 3, 181, 90, 0, 673, 666, 1, 0, 0, 0, 673, 671, 1, 0, 0,
		0, 674, 170, 1, 0, 0, 0, 675, 677, 3, 185, 92, 0, 676, 675, 1, 0, 0, 0,
		677, 678, 1, 0, 0, 0, 678, 676, 1, 0, 0, 0, 678, 679, 1, 0, 0, 0, 679,
		687, 1, 0, 0, 0, 680, 684, 5, 46, 0, 0, 681, 683, 3, 185, 92, 0, 682, 681,
		1, 0, 0, 0, 683, 686, 1, 0, 0, 0, 684, 682, 1, 0, 0, 0, 684, 685, 1, 0,
		0, 0, 685, 688, 1, 0, 0, 0, 686, 684, 1, 0, 0, 0, 687, 680, 1, 0, 0, 0,
		687, 688, 1, 0, 0, 0, 688, 689, 1, 0, 0, 0, 689, 690, 3, 183, 91, 0, 690,
		700, 1, 0, 0, 0, 691, 693, 5, 46, 0, 0, 692, 694, 3, 185, 92, 0, 693, 692,
		1, 0, 0, 0, 694, 695, 1, 0, 0, 0, 695, 693, 1, 0, 0, 0, 695, 696, 1, 0,
		0, 0, 696, 697, 1, 0, 0, 0, 697, 698, 3, 183, 91, 0, 698, 700, 1, 0, 0,
		0, 699, 676, 1, 0, 0, 0, 699, 691, 1, 0, 0, 0, 700, 172, 1, 0, 0, 0, 701,
		704, 3, 187, 93, 0, 702, 704, 5, 95, 0, 0, 703, 701, 1, 0, 0, 0, 703, 702,
		1, 0, 0, 0, 704, 710, 1, 0, 0, 0, 705, 709, 3, 187, 93, 0, 706, 709, 3,
		185, 92, 0, 707, 709, 5, 95, 0, 0, 708, 705, 1, 0, 0, 0, 708, 706, 1, 0,
		0, 0, 708, 707, 1, 0, 0, 0, 709, 712, 1, 0, 0, 0, 710, 708, 1, 0, 0, 0,
		710, 711, 1, 0, 0, 0, 711, 174, 1, 0, 0, 0, 712, 710, 1, 0, 0, 0, 713,
		717, 3, 185, 92, 0, 714, 718, 3, 187, 93, 0, 715, 718, 3, 185, 92, 0, 716,
		718, 5, 95, 0, 0, 717, 714, 1, 0, 0, 0, 717, 715, 1, 0, 0, 0, 717, 716,
		1, 0, 0, 0, 718, 719, 1, 0, 0, 0, 719, 717, 1, 0, 0, 0, 719, 720, 1, 0,
		0, 0, 720, 176, 1, 0, 0, 0, 721, 727, 5, 34, 0, 0, 722, 726, 8, 28, 0,
		0, 723, 724, 5, 34, 0, 0, 724, 726, 5, 34, 0, 0, 725, 722, 1, 0, 0, 0,
		725, 723, 1, 0, 0, 0, 726, 729, 1, 0, 0, 0, 727, 725, 1, 0, 0, 0, 727,
		728, 1, 0, 0, 0, 728, 730, 1, 0, 0, 0, 729, 727, 1, 0, 0, 0, 730, 731,
		5, 34, 0, 0, 731, 178, 1, 0, 0, 0, 732, 738, 5, 96, 0, 0, 733, 737, 8,
		29, 0, 0, 734, 735, 5, 96, 0, 0, 735, 737, 5, 96, 0, 0, 736, 733, 1, 0,
		0, 0, 736, 734, 1, 0, 0, 0, 737, 740, 1, 0, 0, 0, 738, 736, 1, 0, 0, 0,
		738, 739, 1, 0, 0, 0, 739, 741, 1, 0, 0, 0, 740, 738, 1, 0, 0, 0, 741,
		742, 5, 96, 0, 0, 742, 180, 1, 0, 0, 0, 743, 750, 3, 185, 92, 0, 744, 746,
		5, 95, 0, 0, 745, 744, 1, 0, 0, 0, 745, 746, 1, 0, 0, 0, 746, 747, 1, 0,
		0, 0, 747, 749, 3, 185, 92, 0, 748, 745, 1, 0, 0, 0, 749, 752, 1, 0, 0,
		0, 750, 748, 1, 0, 0, 0, 750, 751, 1, 0, 0, 0, 751, 182, 1, 0, 0, 0, 752,
		750, 1, 0, 0, 0, 753, 755, 7, 6, 0, 0, 754, 756, 7, 30, 0, 0, 755, 754,
		1, 0, 0, 0, 755, 756, 1, 0, 0, 0, 756, 758, 1, 0, 0, 0, 757, 759, 3, 185,
		92, 0, 758, 757, 1, 0, 0, 0, 759, 760, 1, 0, 0, 0, 760, 758, 1, 0, 0, 0,
		760, 761, 1, 0, 0, 0, 761, 184, 1, 0, 0, 0, 762, 763, 7, 31, 0, 0, 763,
		186, 1, 0, 0, 0, 764, 765, 7, 32, 0, 0, 765, 188, 1, 0, 0, 0, 29, 0, 195,
		199, 202, 212, 223, 615, 657, 659, 669, 673, 678, 684, 687, 695, 699, 703,
		708, 710, 717, 719, 725, 727, 736, 738, 745, 750, 755, 760, 2, 0, 2, 0,
		0, 1, 0,
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

// SQLLexerInit initializes any static state used to implement SQLLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewSQLLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func SQLLexerInit() {
	staticData := &SQLLexerLexerStaticData
	staticData.once.Do(sqllexerLexerInit)
}

// NewSQLLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewSQLLexer(input antlr.CharStream) *SQLLexer {
	SQLLexerInit()
	l := new(SQLLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &SQLLexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	l.channelNames = staticData.ChannelNames
	l.modeNames = staticData.ModeNames
	l.RuleNames = staticData.RuleNames
	l.LiteralNames = staticData.LiteralNames
	l.SymbolicNames = staticData.SymbolicNames
	l.GrammarFileName = "SQLLexer.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// SQLLexer tokens.
const (
	SQLLexerSIMPLE_COMMENT        = 1
	SQLLexerBRACKETED_COMMENT     = 2
	SQLLexerWS                    = 3
	SQLLexerALL                   = 4
	SQLLexerALIVE                 = 5
	SQLLexerAND                   = 6
	SQLLexerAS                    = 7
	SQLLexerASC                   = 8
	SQLLexerBROKER                = 9
	SQLLexerBROKERS               = 10
	SQLLexerBY                    = 11
	SQLLexerCOMPACT               = 12
	SQLLexerCREATE                = 13
	SQLLexerCROSS                 = 14
	SQLLexerDATABASE              = 15
	SQLLexerDATABASES             = 16
	SQLLexerDEFAULT               = 17
	SQLLexerDESC                  = 18
	SQLLexerDROP                  = 19
	SQLLexerESCAPE                = 20
	SQLLexerEXPLAIN               = 21
	SQLLexerEXISTS                = 22
	SQLLexerFALSE                 = 23
	SQLLexerFIELDS                = 24
	SQLLexerFILTER                = 25
	SQLLexerFLUSH                 = 26
	SQLLexerFROM                  = 27
	SQLLexerGROUP                 = 28
	SQLLexerHAVING                = 29
	SQLLexerIF                    = 30
	SQLLexerIN                    = 31
	SQLLexerJOIN                  = 32
	SQLLexerKEYS                  = 33
	SQLLexerLEFT                  = 34
	SQLLexerLIKE                  = 35
	SQLLexerLIMIT                 = 36
	SQLLexerMASTER                = 37
	SQLLexerMETRICS               = 38
	SQLLexerMETADATA              = 39
	SQLLexerMETADATAS             = 40
	SQLLexerNAMESPACE             = 41
	SQLLexerNAMESPACES            = 42
	SQLLexerNOT                   = 43
	SQLLexerON                    = 44
	SQLLexerOR                    = 45
	SQLLexerORDER                 = 46
	SQLLexerPLAN                  = 47
	SQLLexerREQUESTS              = 48
	SQLLexerREPLICATIONS          = 49
	SQLLexerRIGHT                 = 50
	SQLLexerROLLUP                = 51
	SQLLexerSELECT                = 52
	SQLLexerSHOW                  = 53
	SQLLexerSTATE                 = 54
	SQLLexerSTORAGE               = 55
	SQLLexerTAG                   = 56
	SQLLexerTRUE                  = 57
	SQLLexerTYPES                 = 58
	SQLLexerVALUES                = 59
	SQLLexerWHERE                 = 60
	SQLLexerWITH                  = 61
	SQLLexerWITHIN                = 62
	SQLLexerUSING                 = 63
	SQLLexerUSE                   = 64
	SQLLexerEQ                    = 65
	SQLLexerNEQ                   = 66
	SQLLexerLT                    = 67
	SQLLexerLTE                   = 68
	SQLLexerGT                    = 69
	SQLLexerGTE                   = 70
	SQLLexerPLUS                  = 71
	SQLLexerMINUS                 = 72
	SQLLexerASTERISK              = 73
	SQLLexerSLASH                 = 74
	SQLLexerPERCENT               = 75
	SQLLexerREGEXP                = 76
	SQLLexerNEQREGEXP             = 77
	SQLLexerEXCLAMATION_SYMBOL    = 78
	SQLLexerDOT                   = 79
	SQLLexerLR_BRACKET            = 80
	SQLLexerRR_BRACKET            = 81
	SQLLexerCOMMA                 = 82
	SQLLexerSTRING                = 83
	SQLLexerINTEGER_VALUE         = 84
	SQLLexerDECIMAL_VALUE         = 85
	SQLLexerDOUBLE_VALUE          = 86
	SQLLexerIDENTIFIER            = 87
	SQLLexerDIGIT_IDENTIFIER      = 88
	SQLLexerQUOTED_IDENTIFIER     = 89
	SQLLexerBACKQUOTED_IDENTIFIER = 90
)

// SQLLexerCOMMENT is the SQLLexer channel.
const SQLLexerCOMMENT = 2
