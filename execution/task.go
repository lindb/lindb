package execution

import "context"

type DataDefinitionTask interface {
	Name() string
	Execute(ctx context.Context) error
}
