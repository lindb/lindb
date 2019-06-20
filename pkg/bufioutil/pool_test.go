package bufioutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuffer_PutBuffer(t *testing.T) {
	buf3 := make([]byte, 3)

	PutBuffer(&buf3)

	gotBuf3 := GetBuffer(3)
	assert.NotNil(t, gotBuf3)
	assert.Len(t, *gotBuf3, 3)

	assert.Len(t, *GetBuffer(1), 1)

	gotBuf5 := GetBuffer(5)
	assert.Len(t, *gotBuf5, 5)

	gotBuf4 := GetBuffer(4)
	assert.Len(t, *gotBuf4, 4)
	assert.Equal(t, 4, cap(*gotBuf4))
}
