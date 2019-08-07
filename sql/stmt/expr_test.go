package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
)

func TestExpr_Rewrite(t *testing.T) {
	assert.Equal(t, "f", (&SelectItem{Expr: &FieldExpr{Name: "f"}}).Rewrite())
	assert.Equal(t, "f as f1", (&SelectItem{Expr: &FieldExpr{Name: "f"}, Alias: "f1"}).Rewrite())

	assert.Equal(t, "f", (&FieldExpr{Name: "f"}).Rewrite())

	assert.Equal(t, "sum(f)", (&CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}}).Rewrite())
	assert.Equal(t, "sum()", (&CallExpr{FuncType: function.Sum}).Rewrite())

	assert.Equal(t, "(sum(f)+a)", (&ParenExpr{
		Expr: &BinaryExpr{
			Left:     &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}},
			Operator: ADD,
			Right:    &FieldExpr{Name: "a"},
		}}).Rewrite())

	assert.Equal(t, "sum(f)+a",
		(&BinaryExpr{
			Left:     &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}},
			Operator: ADD,
			Right:    &FieldExpr{Name: "a"},
		}).Rewrite())

	assert.Equal(t, "not tagKey=tagValue",
		(&NotExpr{
			Expr: &EqualsExpr{Key: "tagKey", Value: "tagValue"},
		}).Rewrite())

	assert.Equal(t, "tagKey=tagValue", (&EqualsExpr{Key: "tagKey", Value: "tagValue"}).Rewrite())

	assert.Equal(t, "tagKey like tagValue", (&LikeExpr{Key: "tagKey", Value: "tagValue"}).Rewrite())

	assert.Equal(t, "tagKey in (a,b,c)", (&InExpr{Key: "tagKey", Values: []string{"a", "b", "c"}}).Rewrite())
	assert.Equal(t, "tagKey in ()", (&InExpr{Key: "tagKey"}).Rewrite())

	assert.Equal(t, "tagKey=~Regexp", (&RegexExpr{Key: "tagKey", Regexp: "Regexp"}).Rewrite())
}

func TestTagFilter(t *testing.T) {
	assert.Equal(t, "tagKey", (&EqualsExpr{Key: "tagKey", Value: "tagValue"}).TagKey())
	assert.Equal(t, "tagKey", (&LikeExpr{Key: "tagKey", Value: "tagValue"}).TagKey())
	assert.Equal(t, "tagKey", (&InExpr{Key: "tagKey", Values: []string{"a", "b", "c"}}).TagKey())
	assert.Equal(t, "tagKey", (&RegexExpr{Key: "tagKey", Regexp: "Regexp"}).TagKey())
}
