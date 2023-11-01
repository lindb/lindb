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

package command

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/common/pkg/encoding"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	masterpkg "github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ok := "ok"
	stateMgr := broker.NewMockStateManager(ctrl)
	master := coordinator.NewMockMasterController(ctrl)
	masterStateMgr := masterpkg.NewMockStateManager(ctrl)
	master.EXPECT().GetStateManager().Return(masterStateMgr).AnyTimes()
	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo:     repo,
		StateMgr: stateMgr,
		Master:   master,
	}
	var backend *httptest.Server

	cases := []struct {
		name      string
		statement stmt.Statement
		prepare   func()
		wantErr   bool
	}{
		{
			name:      "show metadata path successfully",
			statement: &stmt.Metadata{MetadataType: stmt.MetadataTypes},
		},
		{
			name:      "state from state machine, but source not found",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.SourceType(100)},
		},
		{
			name:      "state from state machine, but type not found",
			statement: &stmt.Metadata{MetadataType: stmt.MetadataType(100), Source: stmt.StateMachineSource},
		},
		{
			name:      "state from state machine, broker state",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateMachineSource},
			prepare: func() {
				stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{}})
				cli := client.NewMockStateMachineCli(ctrl)
				NewStateMachineCliFn = func() client.StateMachineCli {
					return cli
				}
				cli.EXPECT().FetchStateByNodes(gomock.Any(), gomock.Any()).Return(&ok)
			},
		},
		{
			name:      "state from state machine, master state",
			statement: &stmt.Metadata{MetadataType: stmt.MasterMetadata, Source: stmt.StateMachineSource},
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{}})
				cli := client.NewMockStateMachineCli(ctrl)
				NewStateMachineCliFn = func() client.StateMachineCli {
					return cli
				}
				cli.EXPECT().FetchStateByNodes(gomock.Any(), gomock.Any()).Return(&ok)
			},
		},
		{
			name:      "state from state machine, storage state",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateMachineSource},
			prepare: func() {
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{}})
				cli := client.NewMockStateMachineCli(ctrl)
				NewStateMachineCliFn = func() client.StateMachineCli {
					return cli
				}
				cli.EXPECT().FetchStateByNode(gomock.Any(), gomock.Any()).Return(&ok, nil)
			},
		},
		{
			name:      "show broker metadata, but type not found",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateRepoSource},
		},
		{
			name:      "show broker metadata, but walk entry repo failure",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateRepoSource, Type: "LiveNode"},
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "show broker metadata, but walk entry unmarshal failure",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateRepoSource, Type: "LiveNode"},
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
						fn([]byte("key"), []byte("value"))
						return nil
					})
			},
		},
		{
			name:      "show broker metadata, get live node successfully",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateRepoSource, Type: "LiveNode"},
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
						fn([]byte("key"), encoding.JSONMarshal(&models.Database{Name: "1.1.1.1"}))
						return nil
					})
			},
		},
		{
			name:      "show master metadata, get master successfully",
			statement: &stmt.Metadata{MetadataType: stmt.MasterMetadata, Source: stmt.StateRepoSource, Type: "Master"},
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, fn func(key, value []byte)) error {
						fn([]byte("key"), encoding.JSONMarshal(&models.Master{ElectTime: 11}))
						return nil
					})
			},
		},
		{
			name:      "show storage metadata, but storage name empty",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateRepoSource, Type: "LiveNode", ClusterName: ""},
			wantErr:   true,
		},
		{
			name:      "show storage metadata, but type not found",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateRepoSource, Type: "LiveNode2", ClusterName: "abc"},
			prepare: func() {
				master.EXPECT().IsMaster().Return(true)
			},
		},
		{
			name:      "show storage metadata, but storage state not found",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateRepoSource, Type: "LiveNode", ClusterName: "test"},
			prepare: func() {
				master.EXPECT().IsMaster().Return(true)
				masterStateMgr.EXPECT().GetStorageCluster("test").Return(nil)
			},
		},
		{
			name:      "show storage metadata, no data",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateRepoSource, Type: "LiveNode", ClusterName: "test"},
			prepare: func() {
				master.EXPECT().IsMaster().Return(true)
				storage := masterpkg.NewMockStorageCluster(ctrl)
				masterStateMgr.EXPECT().GetStorageCluster("test").Return(storage)
				storage.EXPECT().GetRepo().Return(repo)
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:      "show storage metadata, forward request failure",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateRepoSource, Type: "LiveNode", ClusterName: "test"},
			prepare: func() {
				master.EXPECT().IsMaster().Return(false)
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: 8089}})
			},
			wantErr: true,
		},
		{
			name:      "show storage metadata, forward request successfully",
			statement: &stmt.Metadata{MetadataType: stmt.StorageMetadata, Source: stmt.StateRepoSource, Type: "LiveNode", ClusterName: "test"},
			prepare: func() {
				port := uint16(8789)
				master.EXPECT().IsMaster().Return(false)
				master.EXPECT().GetMaster().Return(&models.Master{Node: &models.StatelessNode{HostIP: "127.0.0.1", HTTPPort: port}})
				backend = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte("test"))
				}))
				// hack
				_ = backend.Listener.Close()
				l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				assert.NoError(t, err)
				backend.Listener = l
				// Start the server.
				backend.Start()
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				NewStateMachineCliFn = client.NewStateMachineCli
			}()
			if tt.prepare != nil {
				tt.prepare()
			}
			rs, err := MetadataCommand(context.TODO(), deps, nil, tt.statement)
			if (err != nil) != tt.wantErr && rs == nil {
				t.Errorf("MetadataCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if backend != nil {
				backend.Close()
				backend = nil
			}
		})
	}
}
