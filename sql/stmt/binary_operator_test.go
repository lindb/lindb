package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryOPString(t *testing.T) {
	assert.Equal(t, "and", BinaryOPString(AND))
	assert.Equal(t, "or", BinaryOPString(OR))

	assert.Equal(t, "+", BinaryOPString(ADD))
	assert.Equal(t, "-", BinaryOPString(SUB))
	assert.Equal(t, "*", BinaryOPString(MUL))
	assert.Equal(t, "/", BinaryOPString(DIV))

	assert.Equal(t, "unknown", BinaryOPString(UNKNOWN))
}
