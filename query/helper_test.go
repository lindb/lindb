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

package query

import (
	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/metadb"
)

///////////////////////////////////////////////////
//                mock interface				 //
///////////////////////////////////////////////////

func newMockDatabase(ctrl *gomock.Controller) *tsdb.MockDatabase {
	shard := tsdb.NewMockShard(ctrl)
	shard.EXPECT().GetDataFamilies(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	shard.EXPECT().IndexDatabase().Return(nil).AnyTimes()
	metadata := metadb.NewMockMetadata(ctrl)
	metadataIndex := metadb.NewMockMetadataDatabase(ctrl)
	metadata.EXPECT().MetadataDatabase().Return(metadataIndex).AnyTimes()
	metadataIndex.EXPECT().GetMetricID(gomock.Any(), gomock.Any()).Return(uint32(10), nil).AnyTimes()
	metadataIndex.EXPECT().GetField(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(field.Meta{ID: 10, Type: field.SumField}, nil).AnyTimes()

	mockedDatabase := tsdb.NewMockDatabase(ctrl)
	mockedDatabase.EXPECT().GetShard(gomock.Any()).Return(shard, true).AnyTimes()
	mockedDatabase.EXPECT().Metadata().Return(metadata).AnyTimes()
	mockedDatabase.EXPECT().NumOfShards().Return(3).AnyTimes()
	return mockedDatabase
}
