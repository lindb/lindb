package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/pkg/field"
)

func Test_getSegmentStore(t *testing.T) {
	fStore := newFieldStore(field.SumField)
	sStore := fStore.getSegmentStore(11)
	assert.Nil(t, sStore)
}
