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

package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/cockroachdb/pebble/v2"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/tree"
)

// Key-space prefixes replace RocksDB column families.
// Each prefix is a single byte prepended to every key within the logical namespace.
var (
	// metaPrefix reserves key space [0x00, ...] for internal metadata.
	// 0x00 sorts before all user-data prefixes (0x01–0x03), ensuring no collision.
	metaPrefix = []byte{0x00}
	// nsPrefix is the prefix for namespace ID mappings.
	nsPrefix = []byte{0x01}
	// fnPrefix is the prefix for field-name ID mappings (parent is namespace ID).
	fnPrefix = []byte{0x02}
	// fvPrefix is the prefix for field-value ID mappings (parent is field-name ID).
	fvPrefix = []byte{0x03}

	sequenceKey = []byte("sequence")

	// writeOpts disables fsync for every individual write; data is recovered from segment WALs.
	writeOpts = pebble.NoSync
)

type Database interface {
	GetOrCreateNamespaceID(namespace []byte) (uint32, error)
	GetOrCreateFieldKeyID(ns uint32, key []byte) (uint32, error)
	GetOrCreateFieldValueID(key uint32, value []byte) (uint32, error)

	GetNamespaceID(namespace []byte) (uint32, error)
	GetFieldKeyID(ns uint32, key []byte) (uint32, error)
	GetFieldValueID(key uint32, value []byte) (uint32, error)

	FindFieldValueIDs(key uint32, expr tree.Expr) ([]uint32, error)

	ScanField(fieldKey uint32, prefix []byte, callback func(key []byte, value uint32) bool)

	Flush() error
	Close()
}

type database struct {
	db *pebble.DB

	sequence atomic.Uint32
}

func NewDatabase(dbPath string) Database {
	opts := &pebble.Options{
		// Disable the WAL for the index DB — durability is guaranteed by the segment WAL.
		DisableWAL: true,
	}
	db, err := pebble.Open(dbPath, opts)
	if err != nil {
		panic(err)
	}

	indexDB := &database{
		db: db,
	}

	indexDB.initialize()

	return indexDB
}

func (db *database) initialize() {
	// initialize global sequence from persisted state
	v, closer, err := db.db.Get(metaKey(sequenceKey))
	if err == pebble.ErrNotFound {
		// fresh database — sequence starts at 0
		return
	}
	if err != nil {
		panic(err)
	}
	db.sequence.Store(encoding.BytesToU32(v))
	closer.Close() //nolint:errcheck
}

// prefixedKey builds a lookup key by prepending a single-byte CF prefix to key.
func prefixedKey(prefix, key []byte) []byte {
	pk := make([]byte, len(prefix)+len(key))
	copy(pk, prefix)
	copy(pk[len(prefix):], key)
	return pk
}

// metaKey builds an internal metadata key: [0x00] + name.
func metaKey(name []byte) []byte {
	return prefixedKey(metaPrefix, name)
}

func (db *database) GetOrCreateNamespaceID(namespace []byte) (uint32, error) {
	pk := prefixedKey(nsPrefix, namespace)
	v, closer, err := db.db.Get(pk)
	if err != nil && err != pebble.ErrNotFound {
		return 0, err
	}
	if err == nil {
		id := encoding.BytesToU32(v)
		closer.Close() //nolint:errcheck
		return id, nil
	}

	id := db.sequence.Inc()
	ns := make([]byte, len(pk))
	copy(ns, pk)
	if err := db.db.Set(ns, encoding.U32ToBytes(id), writeOpts); err != nil {
		return 0, err
	}
	return id, nil
}

func (db *database) GetOrCreateFieldKeyID(ns uint32, key []byte) (uint32, error) {
	return db.getOrCreateID(fnPrefix, ns, key)
}

func (db *database) GetOrCreateFieldValueID(key uint32, value []byte) (uint32, error) {
	return db.getOrCreateID(fvPrefix, key, value)
}

func (db *database) GetNamespaceID(namespace []byte) (uint32, error) {
	pk := prefixedKey(nsPrefix, namespace)
	v, closer, err := db.db.Get(pk)
	if err == pebble.ErrNotFound {
		return 0, errors.New("not exist")
	}
	if err != nil {
		return 0, err
	}
	id := encoding.BytesToU32(v)
	closer.Close() //nolint:errcheck
	return id, nil
}

func (db *database) GetFieldKeyID(ns uint32, key []byte) (uint32, error) {
	return db.getID(fnPrefix, ns, key)
}

func (db *database) GetFieldValueID(key uint32, value []byte) (uint32, error) {
	return db.getID(fvPrefix, key, value)
}

// ScanField iterates over all field-value entries whose key starts with the given prefix
// under the specified fieldKey namespace.
func (db *database) ScanField(fieldKey uint32, prefix []byte, callback func(key []byte, value uint32) bool) {
	// Build the seek key: fvPrefix + BE(fieldKey) + prefix
	seekKey := buildCompositeKey(fvPrefix, fieldKey, prefix)

	// Compute upper bound: first key of the next parent (fieldKey+1).
	// If fieldKey is MaxUint32, upper bound is the start of the next prefix space.
	var upperBound []byte
	if fieldKey == math.MaxUint32 {
		upperBound = []byte{fvPrefix[0] + 1}
	} else {
		upperBound = make([]byte, 5) // fvPrefix(1B) + BE(fieldKey+1)(4B)
		upperBound[0] = fvPrefix[0]
		binary.BigEndian.PutUint32(upperBound[1:], fieldKey+1)
	}

	iterOpts := &pebble.IterOptions{
		LowerBound: seekKey,
		UpperBound: upperBound,
	}
	it, err := db.db.NewIter(iterOpts)
	if err != nil {
		return
	}
	defer it.Close() //nolint:errcheck

	for valid := it.SeekGE(seekKey); valid; valid = it.Next() {
		k := it.Key()

		if !bytes.HasPrefix(k, seekKey) {
			break
		}
		// Strip the fvPrefix (1 byte) and the 4-byte parent ID to get the original field key.
		ok := callback(k[1+4:], encoding.BytesToU32(it.Value()))
		if !ok {
			break
		}
	}
}

func (db *database) FindFieldValueIDs(key uint32, expr tree.Expr) (ids []uint32, err error) {
	switch expression := expr.(type) {
	case *tree.EqualsExpr:
		id, err := db.GetFieldValueID(key, strutil.String2ByteSlice(expression.Value))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
		return ids, nil
	case *tree.InExpr:
		for _, value := range expression.Values {
			id, idErr := db.GetFieldValueID(key, strutil.String2ByteSlice(value))
			if idErr != nil {
				// Value not in index; skip — IN matches whatever exists.
				continue
			}
			ids = append(ids, id)
		}
		return ids, nil
	case *tree.LikeExpr:
		rp, err := regexp.Compile(likePatternToRegex(expression.Value))
		if err != nil {
			return nil, err
		}
		prefix, _ := rp.LiteralPrefix() // e.g. "^INFO.*$" → prefix="INFO"
		return db.scanFieldValuesByRegexp(key, strutil.String2ByteSlice(prefix), rp), nil
	case *tree.RegexExpr:
		rp, err := regexp.Compile(expression.Regexp)
		if err != nil {
			return nil, err
		}
		prefix, _ := rp.LiteralPrefix() // e.g. "^INFO.*" → prefix="INFO"
		return db.scanFieldValuesByRegexp(key, strutil.String2ByteSlice(prefix), rp), nil
	}
	return ids, nil
}

// scanFieldValuesByRegexp scans field values under fieldKey starting from prefix
// and collects IDs of those whose key matches rp.
// An empty prefix causes a full scan of all values for the field.
func (db *database) scanFieldValuesByRegexp(fieldKey uint32, prefix []byte, rp *regexp.Regexp) []uint32 {
	var ids []uint32
	db.ScanField(fieldKey, prefix, func(value []byte, id uint32) bool {
		if rp.Match(value) {
			ids = append(ids, id)
		}
		return true
	})
	return ids
}

// likePatternToRegex converts a SQL LIKE pattern to an anchored Go regex string.
// SQL '%' matches any sequence of characters; '_' matches exactly one character.
func likePatternToRegex(like string) string {
	var sb strings.Builder
	sb.WriteString("^")
	for _, ch := range like {
		switch ch {
		case '%':
			sb.WriteString(".*")
		case '_':
			sb.WriteString(".")
		default:
			sb.WriteString(regexp.QuoteMeta(string(ch)))
		}
	}
	sb.WriteString("$")
	return sb.String()
}

func (db *database) Flush() error {
	if err := db.db.Set(metaKey(sequenceKey), encoding.U32ToBytes(db.sequence.Load()), writeOpts); err != nil {
		return err
	}
	return db.db.Flush()
}

func (db *database) Close() {
	db.db.Close() //nolint:errcheck
}

// getOrCreateID looks up an ID for (prefix, parent, key). If not found, it allocates
// the next sequence value and stores it.
func (db *database) getOrCreateID(prefix []byte, parent uint32, key []byte) (uint32, error) {
	pk := buildCompositeKey(prefix, parent, key)

	v, closer, err := db.db.Get(pk)
	if err != nil && err != pebble.ErrNotFound {
		return 0, err
	}
	if err == nil {
		id := encoding.BytesToU32(v)
		closer.Close() //nolint:errcheck
		return id, nil
	}

	id := db.sequence.Inc()
	if err := db.db.Set(pk, encoding.U32ToBytes(id), writeOpts); err != nil {
		return 0, err
	}
	return id, nil
}

func (db *database) getID(prefix []byte, parent uint32, key []byte) (uint32, error) {
	pk := buildCompositeKey(prefix, parent, key)

	v, closer, err := db.db.Get(pk)
	if err == pebble.ErrNotFound {
		return 0, errors.New("not exist")
	}
	if err != nil {
		return 0, err
	}
	id := encoding.BytesToU32(v)
	closer.Close() //nolint:errcheck
	return id, nil
}

// buildCompositeKey constructs a key as: prefix(1B) + BE(parent)(4B) + key.
// This replicates the previous key layout used with RocksDB column families.
func buildCompositeKey(prefix []byte, parent uint32, key []byte) []byte {
	pk := make([]byte, len(prefix)+4+len(key))
	copy(pk, prefix)
	binary.BigEndian.PutUint32(pk[len(prefix):], parent)
	copy(pk[len(prefix)+4:], key)
	return pk
}
