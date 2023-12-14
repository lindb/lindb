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

package stmt

// BinaryOP represents binary operation type
type BinaryOP int

const (
	AND BinaryOP = iota + 1
	OR

	ADD
	SUB
	MUL
	DIV

	EQUAL
	NOTEQUAL
	GREATER
	GREATEREQUAL
	LESS
	LESSEQUAL
	LIKE

	UNKNOWN
)

// BinaryOPString returns the binary operator's string value
func BinaryOPString(op BinaryOP) string {
	switch op {
	case AND:
		return "and"
	case OR:
		return "or"
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case EQUAL:
		return "="
	case NOTEQUAL:
		return "!="
	case GREATER:
		return ">"
	case GREATEREQUAL:
		return ">="
	case LESS:
		return "<"
	case LESSEQUAL:
		return "<="
	case LIKE:
		return "like"
	default:
		return "unknown"
	}
}
