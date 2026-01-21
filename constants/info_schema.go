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

package constants

const (
	// InformationSchema represents information schema database.
	InformationSchema = "information_schema"

	TableEnv = "env"

	TableMaster  = "master"
	TableBrokers = "brokers"
	TableMetrics = "metrics"

	TableMetadatas     = "metadatas"
	TableMetadataTypes = "metadata_types"

	TableStorages        = "storages"
	TableEngines         = "engines"
	TableReplications    = "replications"
	TableMemoryDatabases = "memory_databases"
	TableSchemata        = "schemata"
	TableNamespaces      = "namespaces"
	TableTableNames      = "table_names"
	TableColumns         = "columns"

	TableFunctions = "functions"
	TableSnippets  = "snippets"

	TableStreamings    = "streamings"
	TableStreamingJobs = "jos"
)

const (
	// TimestampColumnName represents timestamp column name.
	TimestampColumnName = "timestamp"
)
