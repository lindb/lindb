package sql

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type errorListener struct {
	antlr.ErrorListener
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer,
	offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	panic(msg)
}
