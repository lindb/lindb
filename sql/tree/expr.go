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
	"fmt"
	"strings"
)

// Expr represents a interface for all expression types.
type Expr interface {
	// Rewrite rewrites the expr after parse
	Rewrite() string
}

// ParenExpr represents a parenthesized expression
type ParenExpr struct {
	Expr Expr
}

// BinaryExpr represents an operations with two expressions
type BinaryExpr struct {
	Left, Right Expr
	Operator    BinaryOP
}

// EqualsExpr represents an equals expression
type EqualsExpr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// InExpr represents an in expression
type InExpr struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// LikeExpr represents a like expression
type LikeExpr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RegexExpr represents a regular expression
type RegexExpr struct {
	Name   string `json:"name"`
	Regexp string `json:"regexp"`
}

type NullExpr struct {
	Name string `json:"name"`
	Not  bool   `json:"not"`
}

// Rewrite rewrites the equals expr after parse
func (e *EqualsExpr) Rewrite() string {
	return fmt.Sprintf("%s=%s", e.Name, e.Value)
}

// Rewrite rewrites the in expr after parse
func (e *InExpr) Rewrite() string {
	return fmt.Sprintf("%s in (%s)", e.Name, strings.Join(e.Values, ","))
}

// Rewrite rewrites the like expr after parse
func (e *LikeExpr) Rewrite() string {
	return fmt.Sprintf("%s like %s", e.Name, e.Value)
}

// Rewrite rewrites the regex expr after parse
func (e *RegexExpr) Rewrite() string {
	return fmt.Sprintf("%s=~%s", e.Name, e.Regexp)
}

func (e *NullExpr) Rewrite() string {
	isNot := ""
	if e.Not {
		isNot = "NOT "
	}
	return fmt.Sprintf("%s IS %sNULL", e.Name, isNot)
}
