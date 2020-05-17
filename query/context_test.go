package query

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestStorageExecuteContext(t *testing.T) {
	ctx := newStorageExecuteContext(nil, &stmt.Query{Explain: true})
	ctx.setTagFilterResult(nil)
	assert.NotNil(t, ctx.QueryStats())
}
