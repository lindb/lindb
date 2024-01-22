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
)

// Prometheus represents a configuration for prometheus
type Prometheus struct {
	Namespace string `env:"NAMESPACE" toml:"namespace"`
	Database  string `env:"DATABASE" toml:"database"`
	Field     string `env:"FIELD" toml:"field"`
}

// TOML returns Monitor's toml config
func (p *Prometheus) TOML() string {
	return fmt.Sprintf(`
## Config for the Prometheus
[prometheus]
## namespace
## Default: %s
## Env: LINDB_PROMETHEUS_NAMESPACE
namespace = "%s"
# database
## Default: %s
## Env: LINDB_PROMETHEUS_DATABASE
database = "%s"
# field
## Default: %s
## Env: LINDB_PROMETHEUS_FIELD
field = "%s"`,
		p.Namespace,
		p.Namespace,
		p.Database,
		p.Database,
		p.Field,
		p.Field,
	)
}

// NewDefaultPrometheus returns a new default prometheus config
func NewDefaultPrometheus() *Prometheus {
	return &Prometheus{
		Namespace: "default-ns",
		Database:  "prometheus",
		Field:     "x",
	}
}
