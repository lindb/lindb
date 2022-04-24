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
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestGRPCUtil(t *testing.T) {
	s, m := splitMethodName("service/method")
	assert.Equal(t, "service", s)
	assert.Equal(t, "method", m)
	s, m = splitMethodName("service")
	assert.Equal(t, "unknown", s)
	assert.Equal(t, "unknown", m)

	assert.Equal(t, ClientStream, clientStreamType(&grpc.StreamDesc{ClientStreams: true, ServerStreams: false}))
	assert.Equal(t, ServerStream, clientStreamType(&grpc.StreamDesc{ClientStreams: false, ServerStreams: true}))
	assert.Equal(t, BidiStream, clientStreamType(&grpc.StreamDesc{ClientStreams: true, ServerStreams: true}))

	assert.Equal(t, ClientStream, streamRPCType(&grpc.StreamServerInfo{IsClientStream: true, IsServerStream: false}))
	assert.Equal(t, ServerStream, streamRPCType(&grpc.StreamServerInfo{IsClientStream: false, IsServerStream: true}))
	assert.Equal(t, BidiStream, streamRPCType(&grpc.StreamServerInfo{IsClientStream: true, IsServerStream: true}))
}
