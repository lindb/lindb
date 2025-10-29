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

package trace

type TraceMergeOperator struct{}

func (op *TraceMergeOperator) Name() string {
	return "TraceMergeOperator"
}

func (op *TraceMergeOperator) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	total := len(existingValue)
	for _, v := range operands {
		total += len(v)
	}
	dest := make([]byte, total)
	offset := copy(dest, existingValue)
	for _, operand := range operands {
		offset += copy(dest[offset:], operand)
	}
	return dest, true
}

func (op *TraceMergeOperator) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	dest := make([]byte, (len(rightOperand) + len(leftOperand)))
	copy(dest, leftOperand)
	copy(dest[len(leftOperand):], rightOperand)
	return dest, true
}
