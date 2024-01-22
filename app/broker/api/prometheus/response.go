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

package prometheus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/lindb/common/pkg/logger"
	"github.com/munnerz/goautoneg"
	"github.com/prometheus/prometheus/util/annotations"
	"net/http"
)

// respondError returns error to client.
func (e *ExecuteAPI) respondError(w http.ResponseWriter, apiErr *apiError, data any) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(&Response{
		Status:    statusError,
		ErrorType: apiErr.typ,
		Error:     apiErr.err.Error(),
		Data:      data,
	})
	if err != nil {
		e.logger.Error("error marshaling json response", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var code int
	switch apiErr.typ {
	case errorBadData:
		code = http.StatusBadRequest
	case errorExec:
		code = http.StatusUnprocessableEntity
	case errorCanceled:
		code = statusClientClosedConnection
	case errorTimeout:
		code = http.StatusServiceUnavailable
	case errorInternal:
		code = http.StatusInternalServerError
	case errorNotFound:
		code = http.StatusNotFound
	case errorNotAcceptable:
		code = http.StatusNotAcceptable
	default:
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if n, err := w.Write(b); err != nil {
		e.logger.Error("error writing response", logger.Error(err), logger.Int("bytesWritten", n))
	}
}

// respond returns normal data to the client.
func (e *ExecuteAPI) respond(w http.ResponseWriter, req *http.Request, data any, warnings annotations.Annotations, query string) {
	statusMessage := statusSuccess

	resp := &Response{
		Status:   statusMessage,
		Data:     data,
		Warnings: warnings.AsStrings(query, 10),
	}

	codec, err := e.negotiateCodec(req, resp)
	if err != nil {
		e.respondError(w, &apiError{errorNotAcceptable, err}, nil)
		return
	}

	b, err := codec.Encode(resp)
	if err != nil {
		e.logger.Error("error marshaling response", logger.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", codec.ContentType().String())
	w.WriteHeader(http.StatusOK)
	if n, err := w.Write(b); err != nil {
		e.logger.Error("error marshaling response", logger.Error(err), logger.Int("bytesWritten", n))
	}
}

// negotiateCodec returns a decoder based on the accept of header.
// currently, only application/json is supported.
func (e *ExecuteAPI) negotiateCodec(req *http.Request, resp *Response) (Codec, error) {
	for _, clause := range goautoneg.ParseAccept(req.Header.Get("Accept")) {
		for _, codec := range e.codecs {
			if codec.ContentType().Satisfies(clause) && codec.CanEncode(resp) {
				return codec, nil
			}
		}
	}

	defaultCodec := e.codecs[0]
	if !defaultCodec.CanEncode(resp) {
		return nil, fmt.Errorf("cannot encode response as %s", defaultCodec.ContentType())
	}

	return defaultCodec, nil
}

// response returns error or normal data to client.
func (e *ExecuteAPI) response(c *gin.Context, result apiFuncResult) {
	r, w := c.Request, c.Writer
	if result.err != nil {
		e.respondError(w, result.err, result.data)
	} else if result.data != nil {
		e.respond(w, r, result.data, result.warnings, r.FormValue("query"))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
