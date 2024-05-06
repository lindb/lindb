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
	"encoding/binary"
	"io"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

// Schema represents metric scheam(tags/fields etc.)
type Schema struct {
	Fields  field.Metas
	TagKeys tag.Metas
}

// GetAllHistogramFields returns all histogram fields.
func (s *Schema) GetAllHistogramFields() (rs field.Metas) {
	// with format like __bucket_${boundary}
	for idx := range s.Fields {
		if s.Fields[idx].Type == field.HistogramField {
			rs = append(rs, s.Fields[idx])
		}
	}
	return
}

// MarkPersisted marks all fields/tags as persisted.
func (s *Schema) MarkPersisted() {
	for idx := range s.Fields {
		s.Fields[idx].Persisted = true
	}
	for idx := range s.TagKeys {
		s.TagKeys[idx].Persisted = true
	}
}

// NeedWrite checks whether the schema need persist.
func (s *Schema) NeedWrite() bool {
	for _, f := range s.Fields {
		if !f.Persisted {
			return true
		}
	}
	for _, t := range s.TagKeys {
		if !t.Persisted {
			return true
		}
	}
	return false
}

// Write writes the schema data.
func (s *Schema) Write(w io.Writer) error {
	// FIXME: check max len?
	var scratch [2]byte
	var size uint16
	for _, f := range s.Fields {
		if !f.Persisted {
			size++
		}
	}
	binary.LittleEndian.PutUint16(scratch[:], size)
	// write fields
	if _, err := w.Write(scratch[:]); err != nil {
		return err
	}
	for _, f := range s.Fields {
		if f.Persisted {
			continue
		}
		if err := f.Write(w); err != nil {
			return err
		}
	}
	size = 0
	for _, t := range s.TagKeys {
		if !t.Persisted {
			size++
		}
	}
	// write tag keys
	binary.LittleEndian.PutUint16(scratch[:], size)
	if _, err := w.Write(scratch[:]); err != nil {
		return err
	}
	for _, t := range s.TagKeys {
		if t.Persisted {
			continue
		}
		if err := t.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalFromPersist unmarshals schema from persist.
func (s *Schema) UnmarshalFromPersist(buf []byte) {
	s.unmarshal(buf, true)
}

// Unmarshal unmarshals schema from binary.
func (s *Schema) Unmarshal(buf []byte) {
	s.unmarshal(buf, false)
}

func (s *Schema) unmarshal(buf []byte, persist bool) {
	// read fields
	size := binary.LittleEndian.Uint16(buf[:2])
	buf = buf[2:] // remove size
	for i := uint16(0); i < size; i++ {
		f := &field.Meta{
			Persisted: persist,
		}
		buf = f.Unmarshal(buf)
		s.Fields = append(s.Fields, *f)
	}

	// read tags
	size = binary.LittleEndian.Uint16(buf[:2])
	buf = buf[2:] // remove size
	for i := uint16(0); i < size; i++ {
		t := &tag.Meta{
			Persisted: persist,
		}
		buf = t.Unmarshal(buf)
		s.TagKeys = append(s.TagKeys, *t)
	}
}
