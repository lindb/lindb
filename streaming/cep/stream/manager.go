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

package stream

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

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
	RegisterStreamBySchema(name string, schema *types.TableSchema) error
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

func (mgr *streamManager) RegisterStreamBySchema(name string, schema *types.TableSchema) error {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	mgr.streams[name] = schema
	return nil
}

func (mgr *streamManager) RegisterStreamByType(event any) error {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Pointer {
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
		schema.AddColumn(types.ColumnMetadata{Name: lo.SnakeCase(field.Name), DataType: fieldType(field)})
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
		Database: db,
		Stream:   table,
	}
}

func fieldType(field reflect.StructField) types.DataType {
	t := field.Type
	switch {
	case t == reflect.TypeFor[map[string]string]():
		return types.DTMap
	case t == reflect.TypeFor[string]():
		return types.DTString
	case t == reflect.TypeFor[int](), t == reflect.TypeFor[int32](), t == reflect.TypeFor[int64]():
		return types.DTInt
	case t == reflect.TypeFor[float32](), t == reflect.TypeFor[float64]():
		return types.DTFloat
	case t == reflect.TypeFor[time.Time]():
		return types.DTTimestamp
	case t == reflect.TypeFor[time.Duration]():
		return types.DTDuration
	default:
		return types.DTBinary
	}
	panic(fmt.Sprintf("unsupported field type:%v", field.Type))
}
