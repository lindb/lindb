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

package kv

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestInitStoreManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		once4StoreManager = sync.Once{}
		ctrl.Finish()
	}()
	t.Run("set init store manager", func(t *testing.T) {
		defer func() {
			InitStoreManager(nil)
		}()
		storeMgr := NewMockStoreManager(ctrl)
		InitStoreManager(storeMgr)
		assert.Equal(t, storeMgr, GetStoreManager())
	})
	t.Run("get singleton store manager", func(t *testing.T) {
		defer func() {
			once4StoreManager = sync.Once{}
		}()
		storeMgr1 := GetStoreManager()
		assert.NotNil(t, storeMgr1)
		storeMgr2 := GetStoreManager()
		assert.Equal(t, storeMgr1, storeMgr2)
	})
}

func TestStoreManager_CreateStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newStoreFunc = newStore
		ctrl.Finish()
	}()
	storeMgr := newStoreManager()
	cases := []struct {
		name      string
		storeName string
		prepare   func()
		wantErr   bool
	}{
		{
			name:      "create store successfully",
			storeName: "new-store",
			prepare: func() {
				newStoreFunc = func(name, path string, option StoreOption) (s Store, err error) {
					return NewMockStore(ctrl), nil
				}
			},
		},
		{
			name:      "create store successfully, but it exist",
			storeName: "exist-store",
			prepare: func() {
				newStoreFunc = func(name, path string, option StoreOption) (s Store, err error) {
					return NewMockStore(ctrl), nil
				}
				store, err := storeMgr.CreateStore("exist-store", StoreOption{})
				assert.NoError(t, err)
				assert.NotNil(t, store)
			},
		},
		{
			name:      "create store err",
			storeName: "err-store",
			prepare: func() {
				newStoreFunc = func(name, path string, option StoreOption) (s Store, err error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				newStoreFunc = func(name, path string, option StoreOption) (s Store, err error) {
					return NewMockStore(ctrl), nil
				}
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			store, err := storeMgr.CreateStore(tt.storeName, StoreOption{})
			if ((err != nil) != tt.wantErr && store == nil) || (!tt.wantErr && store == nil) {
				t.Errorf("CreateStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockStoreManager_CloseStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newStoreFunc = newStore
		ctrl.Finish()
	}()
	storeMgr := newStoreManager()
	store := NewMockStore(ctrl)
	newStoreFunc = func(name, path string, option StoreOption) (s Store, err error) {
		return store, nil
	}
	cases := []struct {
		name      string
		storeName string
		prepare   func()
		wantErr   bool
	}{
		{
			name:      "close store err",
			storeName: "test",
			prepare: func() {
				store1, err := storeMgr.CreateStore("test", StoreOption{})
				assert.NoError(t, err)
				assert.Equal(t, store, store1)
				store.EXPECT().close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "close store successfully",
			storeName: "test",
			prepare: func() {
				store1, err := storeMgr.CreateStore("test", StoreOption{})
				assert.NoError(t, err)
				assert.Equal(t, store, store1)
				store.EXPECT().close().Return(nil)
			},
		},
		{
			name:      "close store successfully",
			storeName: "not-exist",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			err := storeMgr.CloseStore(tt.storeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStore() error = %v, wantErr %v", err, tt.wantErr)
			}
			stores := storeMgr.GetStores()
			assert.Empty(t, stores)
		})
	}
}

func TestStoreManager_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newStoreFunc = newStore
		ctrl.Finish()
	}()
	storeMgr := newStoreManager()
	store := NewMockStore(ctrl)
	newStoreFunc = func(name, path string, option StoreOption) (s Store, err error) {
		return store, nil
	}

	stores := storeMgr.GetStores()
	assert.Empty(t, stores)

	store1, err := storeMgr.CreateStore("test", StoreOption{})
	assert.NoError(t, err)
	assert.Equal(t, store, store1)

	store1, ok := storeMgr.GetStoreByName("test")
	assert.True(t, ok)
	assert.Equal(t, store, store1)

	store1, ok = storeMgr.GetStoreByName("not-exist")
	assert.False(t, ok)
	assert.Nil(t, store1)

	stores = storeMgr.GetStores()
	assert.Len(t, stores, 1)
	assert.Equal(t, store, stores[0])
}
