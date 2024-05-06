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

package conntrack

import (
	"fmt"
	"net"
	"os"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

func TestNewTrackedListener(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		newListenFn = net.Listen
		ctrl.Finish()
	}()

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "new tcp listen failure",
			prepare: func() {
				newListenFn = func(network, address string) (net.Listener, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create tcp listen successfully",
			prepare: func() {
				newListenFn = func(network, address string) (net.Listener, error) {
					return NewMockListener(ctrl), nil
				}
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			listener, err := NewTrackedListener("tcp", "1.1.1.1:8080", linmetric.BrokerRegistry)
			if ((err != nil) != tt.wantErr && listener == nil) || (!tt.wantErr && listener == nil) {
				t.Errorf("NewTrackedListener() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTrackedListener_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listener := NewMockListener(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "accept failure",
			prepare: func() {
				listener.EXPECT().Accept().Return(nil, os.ErrDeadlineExceeded)
				listener.EXPECT().Accept().Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "accept successfully",
			prepare: func() {
				listener.EXPECT().Accept().Return(NewMockConn(ctrl), nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tracker := &TrackedListener{
				Listener:   listener,
				statistics: metrics.NewConnStatistics(linmetric.BrokerRegistry, "1.1.1.1:8080"),
			}
			if tt.prepare != nil {
				tt.prepare()
			}

			conn, err := tracker.Accept()
			if ((err != nil) != tt.wantErr && conn == nil) || (!tt.wantErr && conn == nil) {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
