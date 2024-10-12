package expression

import (
	"context"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/utils"
)

// EvalContext is the context for evaluating expression.
type EvalContext interface {
	// CurrentTime returns the current time.
	CurrentTime() time.Time
}

// evalContext implements EvalContext interface.
type evalContext struct {
	ctx context.Context
}

// NewEvalContext creates an EvalContext.
func NewEvalContext(ctx context.Context) EvalContext {
	return &evalContext{
		ctx: ctx,
	}
}

// CurrentTime returns the current time.
func (e *evalContext) CurrentTime() time.Time {
	return utils.GetTimeFromContext(e.ctx, constants.ContextKeyCurrentTime)
}
