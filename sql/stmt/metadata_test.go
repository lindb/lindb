package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
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

func TestMetadata_MarshalJSON(t *testing.T) {
	query := Metadata{
		Namespace:  "ns",
		MetricName: "test",
		Type:       TagValue,
		Condition: &BinaryExpr{
			Left: &ParenExpr{Expr: &BinaryExpr{
				Left:     &InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}},
				Operator: AND,
				Right:    &EqualsExpr{Key: "region", Value: "sh"},
			}},
			Operator: AND,
			Right: &ParenExpr{Expr: &BinaryExpr{
				Left:     &EqualsExpr{Key: "path", Value: "/data"},
				Operator: OR,
				Right:    &EqualsExpr{Key: "path", Value: "/home"},
			}},
		},
		TagKey: "tagKey",
		Prefix: "prefix",
		Limit:  100,
	}

	data := encoding.JSONMarshal(&query)
	query1 := Metadata{}
	err := encoding.JSONUnmarshal(data, &query1)
	assert.NoError(t, err)
	assert.Equal(t, query, query1)
}

func TestMetadata_Marshal_Fail(t *testing.T) {
	query := &Metadata{}
	err := query.UnmarshalJSON([]byte{1, 2, 3})
	assert.NotNil(t, err)
	err = query.UnmarshalJSON([]byte("{\"condition\":\"123\"}"))
	assert.NotNil(t, err)
}
