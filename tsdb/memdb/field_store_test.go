package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newFieldStore(t *testing.T) {
	fStore := newFieldStore()
	assert.NotNil(t, fStore)
}

func Test_getSegmentStore(t *testing.T) {
	fStore := newFieldStore()
	sStore := fStore.getSegmentStore(11)
	assert.NotNil(t, sStore)
}
