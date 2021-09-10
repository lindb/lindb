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

package metric

import (
	"fmt"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/lindb/lindb/series/tag"
)

type BrokerRowFlatDecoder struct {
	block      []byte
	rowBuilder RowBuilder
	originRow  readOnlyRow // used for unmarshal

	compoundValues []float64
	compoundBounds []float64

	namespace    []byte
	enrichedTags tag.Tags
}

var brokerRowFlatDecoderPool sync.Pool

func NewBrokerRowFlatDecoder(
	block []byte,
	namespace []byte,
	enrichedTags tag.Tags,
) (
	decoder *BrokerRowFlatDecoder,
	releaseFunc func(decoder *BrokerRowFlatDecoder),
) {
	releaseFunc = func(decoder *BrokerRowFlatDecoder) { brokerRowFlatDecoderPool.Put(decoder) }
	item := brokerRowFlatDecoderPool.Get()
	if item != nil {
		decoder = item.(*BrokerRowFlatDecoder)
	} else {
		decoder = &BrokerRowFlatDecoder{rowBuilder: *newRowBuilder()}
	}
	decoder.namespace = namespace
	decoder.block = block
	decoder.enrichedTags = enrichedTags
	return decoder, releaseFunc
}

// resetForNextDecode resets context for decoding next row
func (itr *BrokerRowFlatDecoder) resetForNextDecode() {
	itr.rowBuilder.Reset()

	itr.compoundValues = itr.compoundValues[:0]
	itr.compoundBounds = itr.compoundBounds[:0]
}

// HasNext checks if the raw block is fully decode
func (itr *BrokerRowFlatDecoder) HasNext() bool { return len(itr.block) > 0 }

// DecodeTo decodes next flat block into BrokerRow
func (itr *BrokerRowFlatDecoder) DecodeTo(row *BrokerRow) error {
	itr.resetForNextDecode()

	if len(itr.block) < 4 {
		return fmt.Errorf("flat block too short: %d", len(itr.block))
	}
	size := flatbuffers.GetSizePrefix(itr.block, 0)
	partition := itr.block[flatbuffers.SizeUOffsetT : flatbuffers.SizeUOffsetT+size]
	itr.originRow.m.Init(partition, flatbuffers.GetUOffsetT(partition))

	if err := itr.rebuild(); err != nil {
		return err
	}
	defer func() {
		itr.block = itr.block[flatbuffers.SizeUOffsetT+size:]
	}()
	return itr.rowBuilder.BuildTo(row)
}

func (itr *BrokerRowFlatDecoder) rebuild() error {
	var tagCount int
	kvItr := itr.originRow.NewKeyValueIterator()
	for kvItr.HasNext() {
		tagCount++
		if err := itr.rowBuilder.AddTag(kvItr.NextKey(), kvItr.NextValue()); err != nil {
			return err
		}
	}
	if len(itr.enrichedTags) > 0 {
		tagCount++
		for i := 0; i < len(itr.enrichedTags); i++ {
			if err := itr.rowBuilder.AddTag(itr.enrichedTags[i].Key, itr.enrichedTags[i].Value); err != nil {
				return err
			}
		}
	}

	simpleFieldItr := itr.originRow.NewSimpleFieldIterator()
	for simpleFieldItr.HasNext() {
		if err := itr.rowBuilder.AddSimpleField(
			simpleFieldItr.NextRawName(),
			simpleFieldItr.NextRawType(),
			simpleFieldItr.NextValue(),
		); err != nil {
			return err
		}
	}
	compoundFieldItr, ok := itr.originRow.NewCompoundFieldIterator()
	if !ok {
		goto End
	}
	for compoundFieldItr.HasNextBucket() {
		itr.compoundBounds = append(itr.compoundBounds, compoundFieldItr.NextExplicitBound())
		itr.compoundValues = append(itr.compoundValues, compoundFieldItr.NextValue())
	}
	if err := itr.rowBuilder.AddCompoundFieldData(itr.compoundValues, itr.compoundBounds); err != nil {
		return err
	}
	if err := itr.rowBuilder.AddCompoundFieldMMSC(
		compoundFieldItr.Min(),
		compoundFieldItr.Max(),
		compoundFieldItr.Sum(),
		compoundFieldItr.Count(),
	); err != nil {
		return err
	}

End:
	itr.rowBuilder.AddMetricName(itr.originRow.Name())
	if len(itr.namespace) > 0 {
		itr.rowBuilder.AddNameSpace(itr.namespace)
	} else {
		itr.rowBuilder.AddNameSpace(itr.originRow.NameSpace())
	}
	return nil
}
