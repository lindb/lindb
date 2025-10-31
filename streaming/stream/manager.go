package stream

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

var (
	instance *Manager
	once     sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			streams: make(map[string]StreamManager),
		}
	})
	return instance
}

type Manager struct {
	streams map[string]StreamManager

	mutex sync.Mutex
}

func (mgr *Manager) GetStreamManager(app string) StreamManager {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	streamMgr, ok := mgr.streams[app]
	if !ok {
		streamMgr = newStreamManager()
		mgr.streams[app] = streamMgr
	}
	return streamMgr
}

type StreamManager interface {
	spi.MetadataManager

	RegisterStreamByType(event any) error
}

type streamManager struct {
	streams map[string]*types.TableSchema

	mutex sync.RWMutex
}

func newStreamManager() StreamManager {
	return &streamManager{
		streams: make(map[string]*types.TableSchema),
	}
}

func (mgr *streamManager) RegisterStreamByType(event any) error {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return errors.New("input is not a struct or pointer to struct")
	}
	schema := types.NewTableSchema()
	for i := range t.NumField() {
		field := t.Field(i)
		// only support public field
		if field.PkgPath != "" {
			continue
		}
		schema.AddColumn(types.ColumnMetadata{Name: lo.CamelCase(field.Name), DataType: fieldType(field)})
	}
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.streams[t.Name()] = schema
	return nil
}

func (mgr *streamManager) GetTableMetadata(db string, ns string, table string) (*types.TableMetadata, error) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()

	schema, ok := mgr.streams[table]
	fmt.Println(mgr.streams)
	if !ok {
		return nil, errors.New("table not exist")
	}
	return &types.TableMetadata{
		Schema:              schema,
		SupportDynamicField: false,
	}, nil
}

func (mgr *streamManager) GetTableHandle(db string, ns string, table string) spi.TableHandle {
	return &TableHandle{
		App:    db,
		Stream: table,
	}
}

func fieldType(field reflect.StructField) types.DataType {
	switch field.Type.Kind() {
	case reflect.Map:
		return types.DTMap
	case reflect.String:
		return types.DTString
	case reflect.Int, reflect.Int32, reflect.Int64:
		return types.DTInt
	case reflect.Float32:
		return types.DTFloat
	default:
		panic("unsupported field type:" + field.Type.Key().String())
	}
}
