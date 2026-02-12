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

package collections

// Properties represents a collection of key-value pairs where keys are strings.
type Properties struct {
	data map[string]any
}

// NewProperties creates and returns a new Properties instance.
func NewProperties() *Properties {
	return &Properties{
		data: make(map[string]any),
	}
}

// Set sets the value for a given key.
func (p *Properties) Set(key string, value any) {
	p.data[key] = value
}

func (p *Properties) Get(key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

func (p *Properties) GetStringDefault(key, defaultVal string) string {
	if val, ok := p.GetString(key); ok {
		return val
	}
	return defaultVal
}

func (p *Properties) GetString(key string) (string, bool) {
	val, ok := p.data[key]
	if !ok {
		return "", false
	}
	strVal, ok := val.(string)
	return strVal, ok
}

func (p *Properties) GetInt(key string) (int, bool) {
	val, ok := p.data[key]
	if !ok {
		return 0, false
	}
	intVal, ok := val.(int)
	return intVal, ok
}

func (p *Properties) GetBool(key string) (bool, bool) {
	val, ok := p.data[key]
	if !ok {
		return false, false
	}
	boolVal, ok := val.(bool)
	return boolVal, ok
}

func (p *Properties) GetStringSlice(key string) ([]string, bool) {
	val, ok := p.data[key]
	if !ok {
		return nil, false
	}
	switch v := val.(type) {
	case []string:
		return v, true
	case string:
		return []string{v}, true
	}
	return nil, false
}

func (p *Properties) Values() map[string]any {
	return p.data
}
