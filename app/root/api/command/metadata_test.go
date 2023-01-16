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
	"testing"

	"github.com/golang/mock/gomock"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/coordinator/root"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ok := "ok"
	stateMgr := root.NewMockStateManager(ctrl)
	repo := state.NewMockRepository(ctrl)
	deps := &depspkg.HTTPDeps{
		Repo:     repo,
		StateMgr: stateMgr,
	}

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
			statement: &stmt.Metadata{MetadataType: stmt.RootMetadata, Source: stmt.SourceType(100)},
		},
		{
			name:      "state from state machine, live nodes state",
			statement: &stmt.Metadata{MetadataType: stmt.RootMetadata, Source: stmt.StateMachineSource},
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
			name:      "show live nodes metadata, but walk entry repo failure",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateRepoSource, Type: "LiveNode"},
			prepare: func() {
				repo.EXPECT().WalkEntry(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name:      "show live nodes metadata, but type not match",
			statement: &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Source: stmt.StateRepoSource, Type: "Master"},
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
		})
	}
}
