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

package config

import (
	"fmt"
	"reflect"
)

// PrintEnvFormat prints config as env format.
func PrintEnvFormat(v any) {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return
	}
	printStruct(ptrRef.Elem(), "")
}

// printStruct prints struct.
func printStruct(ref reflect.Value, prefix string) {
	refType := ref.Type()

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		refTypeField := refType.Field(i)

		printField(refField, &refTypeField, prefix)
	}
}

// printField prints field.
func printField(refField reflect.Value, refTypeField *reflect.StructField, parent string) {
	prefix := refTypeField.Tag.Get("envPrefix")
	if prefix != "" {
		if reflect.Struct == refField.Kind() {
			printStruct(refField, parent+prefix)
		}
	} else {
		envName := refTypeField.Tag.Get("env")
		if envName != "" {
			fmt.Printf("%s%s=%v\n", parent, envName, refField)
		}
	}
}
