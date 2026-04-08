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

package ingest

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lindb/common/pkg/ltoml"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/app/broker/write"
	writerPkg "github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/internal/concurrent"
	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/internal/mock"
	"github.com/lindb/lindb/metrics"
)

// testPath is a fixed route used for BaseWriter tests.
const testPath = "/test-ingest"

// newTestDeps builds minimal HTTPDeps with a MockManager for BaseWriter tests.
func newTestDeps(ctrl *gomock.Controller) (*deps.HTTPDeps, *write.MockManager) {
	wm := write.NewMockManager(ctrl)
	d := &deps.HTTPDeps{
		BrokerCfg: &config.Broker{
			BrokerBase: config.BrokerBase{
				Ingestion: config.Ingestion{
					IngestTimeout: ltoml.Duration(time.Second * 2),
				},
			},
		},
		WriteManager: wm,
		IngestLimiter: concurrent.NewLimiter(
			context.TODO(),
			32,
			time.Second,
			metrics.NewLimitStatistics("writer_test", linmetric.BrokerRegistry)),
	}
	return d, wm
}

// registerHandler wires a BaseWriter with the given processors onto a fresh gin engine.
func registerHandler(d *deps.HTTPDeps, processors map[string]ProcessFunc) *gin.Engine {
	r := gin.New()
	bw := NewBaseWriter(d, processors)
	r.POST(testPath, bw.Write)
	r.PUT(testPath, bw.Write)
	return r
}

func TestBaseWriter_MissingDatabaseHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, _ := newTestDeps(ctrl)
	r := registerHandler(d, map[string]ProcessFunc{
		constants.ContentTypeOTelProto: func(_ context.Context, _ string, _ []byte) error {
			return nil
		},
	})

	// Neither POST nor PUT supplies X-LinDB-Database → 500
	header := make(http.Header)
	header.Set("Content-Type", constants.ContentTypeOTelProto)

	resp := mock.DoRequest(t, r, http.MethodPost, testPath, "data", header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	resp = mock.DoRequest(t, r, http.MethodPut, testPath, "data", header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestBaseWriter_UnsupportedContentType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, _ := newTestDeps(ctrl)
	r := registerHandler(d, map[string]ProcessFunc{
		constants.ContentTypeOTelProto: func(_ context.Context, _ string, _ []byte) error {
			return nil
		},
	})

	header := make(http.Header)
	header.Set("Content-Type", "text/plain")
	header.Set(constants.DatabaseHeader, "testdb")

	resp := mock.DoRequest(t, r, http.MethodPost, testPath, "data", header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestBaseWriter_ProcessError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, _ := newTestDeps(ctrl)
	r := registerHandler(d, map[string]ProcessFunc{
		constants.ContentTypeOTelProto: func(_ context.Context, _ string, _ []byte) error {
			return errors.New("process failed")
		},
	})

	header := make(http.Header)
	header.Set("Content-Type", constants.ContentTypeOTelProto)
	header.Set(constants.DatabaseHeader, "testdb")

	resp := mock.DoRequest(t, r, http.MethodPost, testPath, "data", header)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestBaseWriter_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, _ := newTestDeps(ctrl)
	r := registerHandler(d, map[string]ProcessFunc{
		constants.ContentTypeOTelProto: func(_ context.Context, _ string, _ []byte) error {
			return nil
		},
	})

	header := make(http.Header)
	header.Set("Content-Type", constants.ContentTypeOTelProto)
	header.Set(constants.DatabaseHeader, "testdb")

	resp := mock.DoRequest(t, r, http.MethodPost, testPath, "data", header)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	resp = mock.DoRequest(t, r, http.MethodPut, testPath, "data", header)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestBaseWriter_ArrowContentType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, _ := newTestDeps(ctrl)
	r := registerHandler(d, map[string]ProcessFunc{
		constants.ContentTypeArrow: func(_ context.Context, _ string, _ []byte) error {
			return nil
		},
	})

	header := make(http.Header)
	header.Set("Content-Type", constants.ContentTypeArrow)
	header.Set(constants.DatabaseHeader, "metricdb")

	resp := mock.DoRequest(t, r, http.MethodPost, testPath, "arrowdata", header)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestStandardProcess_DatabaseNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, wm := newTestDeps(ctrl)
	wm.EXPECT().GetWriter("unknowndb").Return(nil, false)

	proc := StandardProcess(d, constants.EncodingProto)
	err := proc(context.TODO(), "unknowndb", []byte("data"))
	assert.ErrorIs(t, err, constants.ErrDatabaseNotFound)
}

func TestStandardProcess_WriteSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, wm := newTestDeps(ctrl)
	mockWriter := writerPkg.NewMockWriter(ctrl)
	wm.EXPECT().GetWriter("mydb").Return(mockWriter, true)
	mockWriter.EXPECT().Write(gomock.Any(), []byte("payload"), constants.EncodingProto).Return(nil)

	proc := StandardProcess(d, constants.EncodingProto)
	err := proc(context.TODO(), "mydb", []byte("payload"))
	assert.NoError(t, err)
}

func TestStandardProcess_WriteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	d, wm := newTestDeps(ctrl)
	mockWriter := writerPkg.NewMockWriter(ctrl)
	wm.EXPECT().GetWriter("mydb").Return(mockWriter, true)
	mockWriter.EXPECT().Write(gomock.Any(), gomock.Any(), constants.EncodingProto).Return(errors.New("write error"))

	proc := StandardProcess(d, constants.EncodingProto)
	err := proc(context.TODO(), "mydb", []byte("data"))
	assert.Error(t, err)
}
