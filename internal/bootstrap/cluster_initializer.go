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

package bootstrap

import (
	"net/url"
	"path"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
)

//go:generate mockgen -source=./cluster_initializer.go -destination=./cluster_initializer_mock.go -package=bootstrap

// ClusterInitializer initializes cluster(storage/internal database)
type ClusterInitializer interface {
	// InitInternalDatabase initializes internal database
	InitInternalDatabase(sql string) error
}

// clusterInitializer implements ClusterInitializer interface.
type clusterInitializer struct {
	endpoint string
}

// NewClusterInitializer creates a initializer
func NewClusterInitializer(endpoint string) ClusterInitializer {
	u, _ := url.Parse(endpoint)
	u.Path = path.Join(u.Path, constants.APIVersion1CliPath)
	return &clusterInitializer{endpoint: u.String()}
}

// InitInternalDatabase initializes internal database
func (i *clusterInitializer) InitInternalDatabase(sql string) error {
	cli := client.NewExecuteCli(i.endpoint)
	if _, err := cli.Execute(models.ExecuteParam{
		SQL: sql,
	}); err != nil {
		return err
	}
	return nil
}
