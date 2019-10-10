package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
)

func TestExpr_Rewrite(t *testing.T) {
	assert.Equal(t, "f", (&SelectItem{Expr: &FieldExpr{Name: "f"}}).Rewrite())
	assert.Equal(t, "1.90", (&SelectItem{Expr: &NumberLiteral{Val: 1.9}}).Rewrite())
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

func TestExpr_Marshal_Fail(t *testing.T) {
	data := Marshal(nil)
	assert.Nil(t, data)
}

func TestExpr_Unmarshal_Fail(t *testing.T) {
	_, err := Unmarshal([]byte{1, 2, 3})
	assert.NotNil(t, err)
	_, err = Unmarshal([]byte("{\"type\":\"unknown\"}"))
	assert.NotNil(t, err)
	_, err = unmarshal(&exprData{Type: "test", Expr: []byte{1, 2, 3}}, &EqualsExpr{})
	assert.NotNil(t, err)
	_, err = unmarshalCall([]byte{1, 2, 3})
	assert.NotNil(t, err)
	_, err = unmarshalCall([]byte("{\"type\":\"call\",\"params\":[\"213\"]}"))
	assert.NotNil(t, err)
	_, err = Unmarshal([]byte("{\"type\":\"paren\",\"expr\":[\"213\"]}"))
	assert.NotNil(t, err)
	_, err = Unmarshal([]byte("{\"type\":\"number\",\"expr\":{\"val\":\"sf\"}}"))
	assert.NotNil(t, err)
	_, err = Unmarshal([]byte("{\"type\":\"not\",\"expr\":[\"213\"]}"))
	assert.NotNil(t, err)
	_, err = unmarshalSelectItem([]byte("324"))
	assert.NotNil(t, err)
	_, err = unmarshalSelectItem([]byte("{\"type\":\"selectItem\",\"expr\":[\"213\"]}"))
	assert.NotNil(t, err)
	_, err = unmarshalBinary([]byte("123"))
	assert.NotNil(t, err)
	_, err = unmarshalBinary([]byte("{\"type\":\"binary\",\"left\":\"123\"}"))
	assert.NotNil(t, err)
	_, err = unmarshalBinary([]byte("{\"type\":\"binary\",\"left\":{\"type\":\"field\",\"expr\":{\"name\":\"f\"}}," +
		"\"right\":\"123\"}"))
	assert.NotNil(t, err)
}

func TestRegexExpr_Marshal(t *testing.T) {
	expr := &RegexExpr{Key: "tagKey", Regexp: "Regexp"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*RegexExpr)
	assert.Equal(t, *expr, *e)
}

func TestLikeExpr_Marshal(t *testing.T) {
	expr := &LikeExpr{Key: "tagKey", Value: "tagValue"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*LikeExpr)
	assert.Equal(t, *expr, *e)
}

func TestInExpr_Marshal(t *testing.T) {
	expr := &InExpr{Key: "tagKey"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*InExpr)
	assert.Equal(t, *expr, *e)

	expr = &InExpr{Key: "tagKey", Values: []string{"a", "b", "c"}}
	data = Marshal(expr)
	exprData, _ = Unmarshal(data)
	e = exprData.(*InExpr)
	assert.Equal(t, *expr, *e)
}

func TestEqualsExpr_Marshal(t *testing.T) {
	expr := &EqualsExpr{Key: "tagKey", Value: "tagValue"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*EqualsExpr)
	assert.Equal(t, *expr, *e)
}

func TestNotExpr_Marshal(t *testing.T) {
	expr := &NotExpr{
		Expr: &EqualsExpr{Key: "tagKey", Value: "tagValue"},
	}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*NotExpr)
	assert.Equal(t, *expr, *e)
}

func TestNumberLiteral_Marshal(t *testing.T) {
	expr := &SelectItem{Expr: &NumberLiteral{Val: 19.0}}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	assert.NoError(t, err)
	e := exprData.(*SelectItem)
	assert.Equal(t, *expr, *e)
}

func TestSelectItem_Marshal(t *testing.T) {
	expr := &SelectItem{Expr: &FieldExpr{Name: "f"}, Alias: "f1"}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	e := exprData.(*SelectItem)
	assert.Equal(t, *expr, *e)
}

func TestCallExpr_Marshal(t *testing.T) {
	expr := &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	e := exprData.(*CallExpr)
	assert.Equal(t, *expr, *e)
}

func TestParenExpr_Marshal(t *testing.T) {
	expr := &ParenExpr{
		Expr: &BinaryExpr{
			Left:     &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}},
			Operator: ADD,
			Right:    &FieldExpr{Name: "a"},
		}}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	e := exprData.(*ParenExpr)
	assert.Equal(t, *expr, *e)
}

func TestBinaryExpr_Marshal(t *testing.T) {
	expr := &BinaryExpr{
		Left:     &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}},
		Operator: ADD,
		Right:    &FieldExpr{Name: "a"},
	}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*BinaryExpr)
	assert.Equal(t, *expr, *e)
}
