package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataType_String(t *testing.T) {
	assert.Equal(t, "database", Database.String())
	assert.Equal(t, "namespace", Namespace.String())
	assert.Equal(t, "measurement", Metric.String())
	assert.Equal(t, "field", Field.String())
	assert.Equal(t, "tagKey", TagKey.String())
	assert.Equal(t, "tagValue", TagValue.String())
	assert.Equal(t, "unknown", MetadataType(0).String())
}
