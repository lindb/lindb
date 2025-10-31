package context

import (
	"context"

	"github.com/lindb/lindb/streaming/stream/input"
)

type QueryContext struct {
	Context      context.Context
	Query        string
	InputManager input.InputManager
}
