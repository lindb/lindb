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

package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInt(ask string) (int, error) {
	answer, err := readString(ask)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(answer, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func readString(ask string) (string, error) {
	var answer string
	_, _ = fmt.Fprintf(os.Stdout, ask)
	if _, err := fmt.Scanln(&answer); err != nil {
		return answer, err
	}
	return strings.TrimSpace(answer), nil
}
