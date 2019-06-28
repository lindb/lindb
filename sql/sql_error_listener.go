package sql

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type sqlErrorListener struct {
	antlr.ErrorListener
}

func (this *sqlErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line int, column int, msg string, e antlr.RecognitionException) {
	panic(fmt.Sprintf("LinSQL parse error at line %d: %d %s", line, column, msg))
}
