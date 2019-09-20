package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestQuery_Marshal(t *testing.T) {
	query := Query{
		MetricName: "test",
		SelectItems: []Expr{
			&SelectItem{Expr: &FieldExpr{Name: "a"}},
			&SelectItem{Expr: &FieldExpr{Name: "b"}},
			&SelectItem{
				Expr: &CallExpr{
					FuncType: function.Stddev,
					Params: []Expr{
						&CallExpr{
							FuncType: function.Max,
							Params: []Expr{&CallExpr{
								FuncType: function.Sum,
								Params:   []Expr{&FieldExpr{Name: "c"}}},
							},
						},
					},
				},
			},
		},
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
		TimeRange: timeutil.TimeRange{Start: 10, End: 30},
		Interval:  1000,
		GroupBy:   []string{"a", "b", "c"},
		Limit:     100,
	}

	data := encoding.JSONMarshal(&query)
	query1 := Query{}
	err := encoding.JSONUnmarshal(data, &query1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, query, query1)
	assert.True(t, query.HasGroupBy())
}

func TestQuery_Marshal_Fail(t *testing.T) {
	query := &Query{}
	err := query.UnmarshalJSON([]byte{1, 2, 3})
	assert.NotNil(t, err)
	err = query.UnmarshalJSON([]byte("{\"condition\":\"123\"}"))
	assert.NotNil(t, err)
	err = query.UnmarshalJSON([]byte("{\"selectItems\":[\"123\"]}"))
	assert.NotNil(t, err)
}
