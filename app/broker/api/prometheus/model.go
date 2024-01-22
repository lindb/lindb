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
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/munnerz/goautoneg"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/model/textparse"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/annotations"
	"math"
	"time"
)

var (
	// MinTime is the default timestamp used for the begin of optional time ranges.
	// Exposed to let downstream projects to reference it.
	MinTime = time.Unix(math.MinInt64/1000+62135596801, 0).UTC()

	// MaxTime is the default timestamp used for the end of optional time ranges.
	// Exposed to let downstream projects to reference it.
	MaxTime = time.Unix(math.MaxInt64/1000-62135596801, 999999999).UTC()

	minTimeFormatted = MinTime.Format(time.RFC3339Nano)
	maxTimeFormatted = MaxTime.Format(time.RFC3339Nano)
)

type status string

const (
	statusSuccess status = "success"
	statusError   status = "error"

	// Non-standard status code (originally introduced by nginx) for the case when a client closes
	// the connection while the server is still processing the request.
	statusClientClosedConnection = 499
)

type errorType string

const (
	errorNone          errorType = ""
	errorTimeout       errorType = "timeout"
	errorCanceled      errorType = "canceled"
	errorExec          errorType = "execution"
	errorBadData       errorType = "bad_data"
	errorInternal      errorType = "internal"
	errorUnavailable   errorType = "unavailable"
	errorNotFound      errorType = "not_found"
	errorNotAcceptable errorType = "not_acceptable"
)

type apiError struct {
	typ errorType
	err error
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s: %s", e.typ, e.err)
}

func returnAPIError(err error) *apiError {
	if err == nil {
		return nil
	}

	cause := errors.Unwrap(err)
	if cause == nil {
		cause = err
	}

	switch cause.(type) {
	case promql.ErrQueryCanceled:
		return &apiError{errorCanceled, err}
	case promql.ErrQueryTimeout:
		return &apiError{errorTimeout, err}
	case promql.ErrStorage:
		return &apiError{errorInternal, err}
	}

	if errors.Is(err, context.Canceled) {
		return &apiError{errorCanceled, err}
	}

	return &apiError{errorExec, err}
}

// Response contains a response to a HTTP API request.
type Response struct {
	Status    status    `json:"status"`
	Data      any       `json:"data,omitempty"`
	ErrorType errorType `json:"errorType,omitempty"`
	Error     string    `json:"error,omitempty"`
	Warnings  []string  `json:"warnings,omitempty"`
}

// apiFuncResult is the return result of the Prometheus API.
type apiFuncResult struct {
	data      any
	err       *apiError
	warnings  annotations.Annotations
	finalizer func()
}

// metadata is the structure that the /metadata endpoint depends on.
type metadata struct {
	Type textparse.MetricType `json:"type"`
	Help string               `json:"help"`
	Unit string               `json:"unit"`
}

// seriesSet is implementation of storage.SeriesSet.
type seriesSet struct {
	series []storage.Series
	index  int
	err    error
}

func newSeriesSet() *seriesSet {
	return &seriesSet{index: -1}
}

func (s *seriesSet) Next() bool {
	if len(s.series) == 0 {
		return false
	}
	if s.index < len(s.series)-1 {
		s.index++
		return true
	}
	return false
}

func (s *seriesSet) At() storage.Series {
	if s.index < 0 {
		return nil
	}
	return s.series[s.index]
}

func (s *seriesSet) Err() error                        { return s.err }
func (s *seriesSet) Warnings() annotations.Annotations { return nil }

func (s *seriesSet) setErr(err error) {
	s.err = err
}

func (s *seriesSet) setSeries(series []storage.Series) {
	s.series = series
}

// A Codec performs encoding of API responses.
type Codec interface {
	// ContentType returns the MIME time that this Codec emits.
	ContentType() MIMEType

	// CanEncode determines if this Codec can encode resp.
	CanEncode(resp *Response) bool

	// Encode encodes resp, ready for transmission to an API consumer.
	Encode(resp *Response) ([]byte, error)
}

// JSONCodec is a Codec that encodes API responses as JSON.
type JSONCodec struct{}

func (j JSONCodec) ContentType() MIMEType {
	return MIMEType{Type: "application", SubType: "json"}
}

func (j JSONCodec) CanEncode(_ *Response) bool {
	return true
}

func (j JSONCodec) Encode(resp *Response) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Marshal(resp)
}

type MIMEType struct {
	Type    string
	SubType string
}

func (m MIMEType) String() string {
	return m.Type + "/" + m.SubType
}

func (m MIMEType) Satisfies(accept goautoneg.Accept) bool {
	if accept.Type == "*" && accept.SubType == "*" {
		return true
	}

	if accept.Type == m.Type && accept.SubType == "*" {
		return true
	}

	if accept.Type == m.Type && accept.SubType == m.SubType {
		return true
	}

	return false
}