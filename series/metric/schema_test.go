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
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

func TestSchema_GetAllHistogramFields(t *testing.T) {
	schema := &Schema{
		Fields: field.Metas{{ID: 1, Name: "field1"}},
	}
	assert.Empty(t, schema.GetAllHistogramFields())
	schema.Fields = append(schema.Fields, field.Meta{Name: "f", Type: field.HistogramField})
	assert.Len(t, schema.GetAllHistogramFields(), 1)
}

func TestSchema_MarkPersisted(t *testing.T) {
	schema := &Schema{
		Fields:  field.Metas{{ID: 1, Name: "field1"}},
		TagKeys: tag.Metas{{ID: 2, Key: "key1"}},
	}
	assert.True(t, schema.NeedWrite())
	schema.Fields[0].Persisted = true
	assert.True(t, schema.NeedWrite())
	schema.MarkPersisted()
	assert.False(t, schema.NeedWrite())
}

func TestSchema_Marshal(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	schema := &Schema{
		Fields:  field.Metas{{ID: 1, Name: "field1"}},
		TagKeys: tag.Metas{{ID: 2, Key: "key1"}},
	}
	assert.NoError(t, schema.Write(buf))
	data := buf.Bytes()
	schema2 := &Schema{}
	schema2.Unmarshal(data)
	assert.Equal(t, schema, schema2)

	schema3 := &Schema{
		Fields:  field.Metas{{ID: 2, Name: "field2"}},
		TagKeys: tag.Metas{{ID: 4, Key: "key2"}},
	}
	buf.Reset()
	assert.NoError(t, schema3.Write(buf))
	data = buf.Bytes()

	schema2.Unmarshal(data)
	tagMeta, ok := schema2.TagKeys.Find("key1")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(2), tagMeta.ID)
	tagMeta, ok = schema2.TagKeys.Find("key2")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(4), tagMeta.ID)

	buf.Reset()
	assert.NoError(t, schema2.Write(buf))
	data = buf.Bytes()

	schema4 := &Schema{}
	schema4.Unmarshal(data)
	assert.True(t, schema4.NeedWrite())
	tagMeta, ok = schema4.TagKeys.Find("key1")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(2), tagMeta.ID)
	tagMeta, ok = schema4.TagKeys.Find("key2")
	assert.True(t, ok)
	assert.Equal(t, tag.KeyID(4), tagMeta.ID)
	schema4 = &Schema{}
	schema4.UnmarshalFromPersist(data)
	assert.False(t, schema4.NeedWrite())
}

func TestSchema_Write_Error(t *testing.T) {
	schema := &Schema{
		Fields: field.Metas{
			{ID: 2, Name: "field2"},
			{ID: 2, Name: "field2", Persisted: true},
		},
		TagKeys: tag.Metas{{ID: 4, Key: "key2", Persisted: true}, {ID: 5, Key: "key"}},
	}
	cases := []struct {
		name string
		w    io.Writer
	}{
		{
			name: "write fields size error",
			w:    mock.NewWriter(1),
		},
		{
			name: "write field error",
			w:    mock.NewWriter(2),
		},
		{
			name: "write tags size error",
			w:    mock.NewWriter(5),
		},
		{
			name: "write tag error",
			w:    mock.NewWriter(6),
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, schema.Write(tt.w))
		})
	}
}
