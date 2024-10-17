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

type Env struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Default string `json:"default"`
}

// PrintEnvFormat prints config as env format.
func PrintEnvFormat(v any) {
	envs := toEnvs(v)
	for k, v := range envs {
		fmt.Printf("%s=%s\n", k, v)
	}
}

func ToEnvs(cfg, defaultCfg any) (rs []Env) {
	envs := toEnvs(cfg)
	defaultEnvs := toEnvs(defaultCfg)
	for k, v := range envs {
		val := v
		defaultV := defaultEnvs[k]
		if v == defaultV {
			val = ""
		}
		rs = append(rs, Env{
			Key:     k,
			Value:   val,
			Default: defaultV,
		})
	}
	return
}

func toEnvs(cfg any) map[string]string {
	envsMap := make(map[string]string)
	ptrRef := reflect.ValueOf(cfg)
	if ptrRef.Kind() != reflect.Ptr {
		return nil
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return nil
	}
	structToEnvs(ptrRef.Elem(), "", envsMap)
	return envsMap
}

// printStruct prints struct.
func structToEnvs(ref reflect.Value, prefix string, envs map[string]string) {
	refType := ref.Type()

	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		refTypeField := refType.Field(i)

		fieldToEnvs(refField, &refTypeField, prefix, envs)
	}
}

// fieldToEnvs prints field.
func fieldToEnvs(refField reflect.Value, refTypeField *reflect.StructField, parent string, envs map[string]string) {
	prefix := refTypeField.Tag.Get("envPrefix")
	if prefix != "" {
		if reflect.Struct == refField.Kind() {
			structToEnvs(refField, parent+prefix, envs)
		}
	} else {
		envName := refTypeField.Tag.Get("env")
		if envName != "" {
			envs[fmt.Sprintf("%s%s", parent, envName)] = fmt.Sprintf("%v", refField)
		}
	}
}
