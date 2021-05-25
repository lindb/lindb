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

package queue

const (
	dataPath  = "data"
	indexPath = "index"
	metaPath  = "meta"

	metaPageIndex = 0

	indexItemLength          = 8 + 4 + 4 // data page id(8bytes) + message offset in data page(4bytes) + message length(4bytes)
	indexItemsPerPage        = 1024 * 256
	indexPageSize            = indexItemsPerPage * indexItemLength
	dataPageSize             = 128 * 1024 * 1024 // 128MB
	metaPageSize             = 8 + 8 + 8 + 8     // headSeq(int64), tailSeq(int64), data expire page, index expire page
	queueHeadSeqOffset       = 0
	queueTailSeqOffset       = queueHeadSeqOffset + 8
	queueExpireDataOffset    = queueTailSeqOffset + 8
	queueExpireIndexOffset   = queueExpireDataOffset + 8
	queueDataPageIndexOffset = 0
	messageOffsetOffset      = 8
	messageLengthOffset      = 8 + 4

	defaultDataSizeLimit = 4 * dataPageSize

	fanOutDirName = "fan_out"
	// headSeq(int64), tailSeq(int64)
	fanOutMetaSize      = 8 + 8
	fanOutHeadSeqOffset = 0
	fanOutTailSeqOffset = 8
	// SeqNoNewMessageAvailable is the seqNum returned when no new message available
	SeqNoNewMessageAvailable = int64(-1)
)
