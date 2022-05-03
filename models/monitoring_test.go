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

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplicatorState_String(t *testing.T) {
	assert.Equal(t, "Init", ReplicatorInitState.String())
	assert.Equal(t, "Ready", ReplicatorReadyState.String())
	assert.Equal(t, "Failure", ReplicatorFailureState.String())
	assert.Equal(t, "Unknown", ReplicatorUnknownState.String())
}

func TestReplicatorState_MarshalJSON(t *testing.T) {
	json, err := ReplicatorInitState.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"Init"`, string(json))
}

func TestReplicatorState_UnmarshalJSON(t *testing.T) {
	rs := ReplicatorReadyState
	err := rs.UnmarshalJSON([]byte(`"Init"`))
	assert.NoError(t, err)
	assert.Equal(t, ReplicatorInitState, rs)
	err = rs.UnmarshalJSON([]byte(`"Ready"`))
	assert.NoError(t, err)
	assert.Equal(t, ReplicatorReadyState, rs)
	err = rs.UnmarshalJSON([]byte(`"Failure"`))
	assert.NoError(t, err)
	assert.Equal(t, ReplicatorFailureState, rs)
	err = rs.UnmarshalJSON([]byte(`"jjj"`))
	assert.NoError(t, err)
	assert.Equal(t, ReplicatorUnknownState, rs)
}
