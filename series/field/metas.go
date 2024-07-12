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

package field

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"
	"strings"

	"github.com/lindb/lindb/pkg/stream"
)

// Meta is the meta-data for field, which contains field-name, fieldID and field-type
type Meta struct {
	Name Name `json:"name"`
	Type Type `json:"type"` // query not use type
	ID   ID   `json:"id"`   // query not use id, don't get id in query phase
	// write: field index under memory database
	// read: field index of query fields
	Index     uint8
	Persisted bool // FIXME: can remove
}

// MarshalBinary marshals meta as binary.
func (m *Meta) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	writer := stream.NewBufferWriter(&buf)
	writer.PutByte(byte(m.ID))
	writer.PutByte(byte(m.Type))
	writer.PutInt16(int16(len(m.Name)))
	writer.PutBytes([]byte(m.Name))
	return buf.Bytes(), writer.Error()
}

// Write writes write field meta.
func (m *Meta) Write(w io.Writer) error {
	if _, err := w.Write([]byte{byte(m.ID), byte(m.Type)}); err != nil {
		return err
	}
	var scratch [2]byte
	binary.LittleEndian.PutUint16(scratch[:], uint16(len(m.Name)))
	if _, err := w.Write(scratch[:]); err != nil {
		return err
	}
	if _, err := w.Write([]byte(m.Name)); err != nil {
		return err
	}
	return nil
}

// Unmarshal unmarshals meta from binary.
func (m *Meta) Unmarshal(buf []byte) []byte {
	m.ID = ID(buf[0])
	m.Type = Type(buf[1])
	size := binary.LittleEndian.Uint16(buf[2:])
	end := 4 + size
	m.Name = Name(buf[4:end])
	return buf[end:]
}

// Metas implements sort.Interface, it's sorted by name
type Metas []Meta

func (fms Metas) Len() int { return len(fms) }

func (fms Metas) Less(i, j int) bool { return fms[i].Name < fms[j].Name }

func (fms Metas) Swap(i, j int) { fms[i], fms[j] = fms[j], fms[i] }

func UnmarshalBinary(data []byte) (Metas, ID, error) {
	reader := stream.NewReader(data)
	var max ID
	var fms Metas

	for !reader.Empty() && reader.Error() == nil {
		id := ID(reader.ReadByte())
		fType := Type(reader.ReadByte())
		nameLen := reader.ReadInt16()
		name := reader.ReadBytes(int(nameLen))
		fms = append(fms, Meta{ID: id, Type: fType, Name: Name(name)})
		if id > max {
			max = id
		}
	}
	return fms, max, reader.Error()
}

// Find returns Meta by given field name, if not exist returns false.
func (fms Metas) Find(fieldName Name) (Meta, bool) {
	for _, f := range fms {
		if f.Name == fieldName {
			return f, true
		}
	}
	return Meta{}, false
}

// GetFromName searches the meta by fieldName, return false when not exist
func (fms Metas) GetFromName(fieldName Name) (Meta, bool) {
	idx := sort.Search(len(fms), func(i int) bool { return fms[i].Name >= fieldName })
	if idx >= len(fms) || fms[idx].Name != fieldName {
		return Meta{}, false
	}
	return fms[idx], true
}

func (fms Metas) FindIndexByName(fieldName Name) (int, bool) {
	idx := sort.Search(len(fms), func(i int) bool { return fms[i].Name >= fieldName })
	if idx >= len(fms) || fms[idx].Name != fieldName {
		return -1, false
	}
	return idx, true
}

// GetFromID searches the meta by fieldID, returns false when not exist
func (fms Metas) GetFromID(fieldID ID) (Meta, bool) {
	for _, fm := range fms {
		if fm.ID == fieldID {
			return fm, true
		}
	}
	return Meta{}, false
}

// Clone clones a copy of fieldsMetas
func (fms Metas) Clone() (x2 Metas) {
	x2 = make([]Meta, fms.Len())
	copy(x2, fms)
	return x2
}

// Intersects checks whether each fieldID is in the list,
// and returns the new meta-list corresponding with the fieldID-list.
func (fms Metas) Intersects(fields Metas) (x2 Metas, isSubSet bool) {
	isSubSet = true
	for _, f := range fields {
		if fm, ok := fms.GetFromID(f.ID); ok {
			x2 = append(x2, fm)
		} else {
			isSubSet = false
		}
	}
	sort.Sort(x2)
	return x2, isSubSet
}

// Stringer returns the fields in string
func (fms Metas) String() string {
	switch len(fms) {
	case 0:
		return ""
	case 1:
		return string(fms[0].Name)
	case 2:
		return string(fms[0].Name) + "," + string(fms[1].Name)
	default:
		b := strings.Builder{}
		for i := range fms {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(string(fms[i].Name))
		}
		return b.String()
	}
}
