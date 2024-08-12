// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation/function"
)

func TestExpr_Rewrite(t *testing.T) {
	assert.Equal(t, "f", (&SelectItem2{Expr: &FieldExpr{Name: "f"}}).Rewrite())
	assert.Equal(t, "1.90", (&SelectItem2{Expr: &NumberLiteral{Val: 1.9}}).Rewrite())
	assert.Equal(t, "f as f1", (&SelectItem2{Expr: &FieldExpr{Name: "f"}, Alias: "f1"}).Rewrite())

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
			Expr: &EqualsExpr{Name: "tagKey", Value: "tagValue"},
		}).Rewrite())

	assert.Equal(t, "tagKey=tagValue", (&EqualsExpr{Name: "tagKey", Value: "tagValue"}).Rewrite())

	assert.Equal(t, "tagKey like tagValue", (&LikeExpr{Name: "tagKey", Value: "tagValue"}).Rewrite())

	assert.Equal(t, "tagKey in (a,b,c)", (&InExpr{Name: "tagKey", Values: []string{"a", "b", "c"}}).Rewrite())
	assert.Equal(t, "tagKey in ()", (&InExpr{Name: "tagKey"}).Rewrite())

	assert.Equal(t, "tagKey=~Regexp", (&RegexExpr{Name: "tagKey", Regexp: "Regexp"}).Rewrite())

	assert.Equal(t, "f desc", (&OrderByExpr{Expr: &FieldExpr{Name: "f"}, Desc: true}).Rewrite())
	assert.Equal(t, "max(f) asc", (&OrderByExpr{Expr: &CallExpr{FuncType: function.Max, Params: []Expr{&FieldExpr{Name: "f"}}}}).Rewrite())
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
	_, err = unmarshalOrderByExpr([]byte("324"))
	assert.NotNil(t, err)
	_, err = unmarshalOrderByExpr([]byte("{\"type\":\"orderBy\",\"expr\":[\"213\"]}"))
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
	expr := &RegexExpr{Name: "tagKey", Regexp: "Regexp"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*RegexExpr)
	assert.Equal(t, *expr, *e)
}

func TestLikeExpr_Marshal(t *testing.T) {
	expr := &LikeExpr{Name: "tagKey", Value: "tagValue"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*LikeExpr)
	assert.Equal(t, *expr, *e)
}

func TestInExpr_Marshal(t *testing.T) {
	expr := &InExpr{Name: "tagKey"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*InExpr)
	assert.Equal(t, *expr, *e)

	expr = &InExpr{Name: "tagKey", Values: []string{"a", "b", "c"}}
	data = Marshal(expr)
	exprData, _ = Unmarshal(data)
	e = exprData.(*InExpr)
	assert.Equal(t, *expr, *e)
}

func TestEqualsExpr_Marshal(t *testing.T) {
	expr := &EqualsExpr{Name: "tagKey", Value: "tagValue"}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*EqualsExpr)
	assert.Equal(t, *expr, *e)
}

func TestNotExpr_Marshal(t *testing.T) {
	expr := &NotExpr{
		Expr: &EqualsExpr{Name: "tagKey", Value: "tagValue"},
	}
	data := Marshal(expr)
	exprData, _ := Unmarshal(data)
	e := exprData.(*NotExpr)
	assert.Equal(t, *expr, *e)
}

func TestNumberLiteral_Marshal(t *testing.T) {
	expr := &SelectItem2{Expr: &NumberLiteral{Val: 19.0}}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	assert.NoError(t, err)
	e := exprData.(*SelectItem2)
	assert.Equal(t, *expr, *e)
}

func TestSelectItem_Marshal(t *testing.T) {
	expr := &SelectItem2{Expr: &FieldExpr{Name: "f"}, Alias: "f1"}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	assert.NoError(t, err)
	e := exprData.(*SelectItem2)
	assert.Equal(t, *expr, *e)
}

func TestOrderByExpr_Marshal(t *testing.T) {
	t.Run("sample", func(t *testing.T) {
		expr := &OrderByExpr{Expr: &FieldExpr{Name: "f"}, Desc: true}
		data := Marshal(expr)
		exprData, err := Unmarshal(data)
		assert.NoError(t, err)
		e := exprData.(*OrderByExpr)
		assert.Equal(t, *expr, *e)
	})
	t.Run("order by with func", func(t *testing.T) {
		expr := &OrderByExpr{Expr: &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}}}
		data := Marshal(expr)
		exprData, err := Unmarshal(data)
		assert.NoError(t, err)
		e := exprData.(*OrderByExpr)
		assert.Equal(t, *expr, *e)
	})
}

func TestCallExpr_Marshal(t *testing.T) {
	expr := &CallExpr{FuncType: function.Sum, Params: []Expr{&FieldExpr{Name: "f"}}}
	data := Marshal(expr)
	exprData, err := Unmarshal(data)
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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