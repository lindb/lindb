package index

import (
	"encoding/binary"
	"errors"

	"github.com/lindb/common/pkg/fileutil"
	"github.com/linxGnu/grocksdb"
	"github.com/samber/lo"
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/encoding"
)

var (
	wo *grocksdb.WriteOptions
	ro *grocksdb.ReadOptions

	sequenceKey = []byte("sequence")
)

func init() {
	wo = grocksdb.NewDefaultWriteOptions()
	wo.DisableWAL(true)

	ro = grocksdb.NewDefaultReadOptions()
}

type Database interface {
	GetOrCreateNamespaceID(namespace []byte) (uint32, error)
	GetOrCreateFieldKeyID(ns uint32, key []byte) (uint32, error)
	GetOrCreateFieldValueID(key uint32, value []byte) (uint32, error)

	GetNamespaceID(namespace []byte) (uint32, error)
	GetFieldKeyID(ns uint32, key []byte) (uint32, error)
	GetFieldValueID(key uint32, value []byte) (uint32, error)

	Flush() error
	Close()
}

type database struct {
	db         *grocksdb.DB
	namespace  *grocksdb.ColumnFamilyHandle
	fieldName  *grocksdb.ColumnFamilyHandle
	fieldValue *grocksdb.ColumnFamilyHandle

	sequence atomic.Uint32
}

func NewDatabase(dbPath string) Database {
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	var (
		familyNames []string
		db          *grocksdb.DB
		families    []*grocksdb.ColumnFamilyHandle
		err         error
	)
	// opts.SetMergeOperator(&IntAddMergeOperator{})
	if fileutil.Exist(dbPath) {
		// open existing backend database
		familyNames, err = grocksdb.ListColumnFamilies(opts, dbPath)
		if err != nil {
			panic(err)
		}
		var cfOpts []*grocksdb.Options
		for range len(familyNames) {
			cfOpts = append(cfOpts, opts)
		}

		db, families, err = grocksdb.OpenDbColumnFamilies(opts, dbPath, familyNames, cfOpts)
		if err != nil {
			panic(err)
		}
	} else {
		// create new backend database
		db, err = grocksdb.OpenDb(opts, dbPath)
		if err != nil {
			panic(err)
		}
	}

	getOrCreate := func(cfName string) *grocksdb.ColumnFamilyHandle {
		cfh, ok := lo.Find(families, func(cf *grocksdb.ColumnFamilyHandle) bool {
			return cfName == cf.Name()
		})
		if ok {
			return cfh
		}
		cfh, err = db.CreateColumnFamily(opts, cfName)
		if err != nil {
			panic(err)
		}
		return cfh
	}
	namespace := getOrCreate("ns")
	fieldName := getOrCreate("fn")
	fieldValue := getOrCreate("fv")

	indexDB := &database{
		db:         db,
		namespace:  namespace,
		fieldName:  fieldName,
		fieldValue: fieldValue,
	}

	return indexDB
}

func (db *database) initialize() {
	// initialize global sequence
	v, err := db.db.Get(ro, sequenceKey)
	if err != nil {
		panic(err)
	}
	defer v.Free()
	if v.Exists() {
		db.sequence.Store(encoding.BytesToU32(v.Data()))
	}
}

func (db *database) GetOrCreateNamespaceID(namespace []byte) (uint32, error) {
	v, err := db.db.GetCF(ro, db.namespace, namespace)
	if err != nil {
		return 0, err
	}
	defer v.Free()
	if v.Exists() {
		return encoding.BytesToU32(v.Data()), nil
	}

	id := db.sequence.Inc()
	ns := make([]byte, len(namespace))
	copy(ns, namespace)
	db.db.PutCF(wo, db.namespace, ns, encoding.U32ToBytes(id))
	return id, nil
}

func (db *database) GetOrCreateFieldKeyID(ns uint32, key []byte) (uint32, error) {
	return db.getOrCreateID(db.fieldName, ns, key)
}

func (db *database) GetOrCreateFieldValueID(key uint32, value []byte) (uint32, error) {
	return db.getOrCreateID(db.fieldValue, key, value)
}

func (db *database) GetNamespaceID(namespace []byte) (uint32, error) {
	v, err := db.db.GetCF(ro, db.namespace, namespace)
	if err != nil {
		return 0, err
	}
	defer v.Free()
	if v.Exists() {
		return encoding.BytesToU32(v.Data()), nil
	}
	return 0, errors.New("not exist")
}

func (db *database) GetFieldKeyID(ns uint32, key []byte) (uint32, error) {
	return db.getID(db.fieldName, ns, key)
}

func (db *database) GetFieldValueID(key uint32, value []byte) (uint32, error) {
	return db.getID(db.fieldValue, key, value)
}

func (db *database) Flush() error {
	opt := grocksdb.NewDefaultFlushOptions()
	db.db.Put(wo, sequenceKey, encoding.U32ToBytes(db.sequence.Load()))
	return db.db.Flush(opt)
}

func (db *database) Close() {
	db.db.Close()
}

func (db *database) getOrCreateID(cf *grocksdb.ColumnFamilyHandle, parent uint32, key []byte) (uint32, error) {
	keyBytes := make([]byte, 4+len(key))
	binary.BigEndian.PutUint32(keyBytes[:4], parent)
	copy(keyBytes[4:], key)

	v, err := db.db.GetCF(ro, cf, keyBytes)
	if err != nil {
		return 0, err
	}
	defer v.Free()
	if v.Exists() {
		return encoding.BytesToU32(v.Data()), nil
	}

	id := db.sequence.Inc()
	db.db.PutCF(wo, cf, keyBytes, encoding.U32ToBytes(id))
	return id, nil
}

func (db *database) getID(cf *grocksdb.ColumnFamilyHandle, parent uint32, key []byte) (uint32, error) {
	keyBytes := make([]byte, 4+len(key))
	binary.BigEndian.PutUint32(keyBytes[:4], parent)
	copy(keyBytes[4:], key)

	v, err := db.db.GetCF(ro, cf, keyBytes)
	if err != nil {
		return 0, err
	}
	defer v.Free()
	if v.Exists() {
		return encoding.BytesToU32(v.Data()), nil
	}
	return 0, errors.New("not exist")
}
