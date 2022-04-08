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

package version

import (
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/timeutil"
)

// rollup represents the rollup metadata for rollup job.
// source <=> target family reference
type rollup struct {
	// file number -> target interval type for raw family
	rollupFiles map[table.FileNumber][]timeutil.Interval // source family

	// family id -> file number for source family,
	// reference to raw family file number, reference add after rollup successfully
	referenceFiles map[FamilyID][]table.FileNumber // target family
}

// newRollup creates the rollup job metadata
func newRollup() *rollup {
	return &rollup{
		rollupFiles:    make(map[table.FileNumber][]timeutil.Interval),
		referenceFiles: make(map[FamilyID][]table.FileNumber),
	}
}

// addRollupFile adds need rollup file and target intervals
func (r *rollup) addRollupFile(fileNumber table.FileNumber, interval timeutil.Interval) {
	r.rollupFiles[fileNumber] = append(r.rollupFiles[fileNumber], interval)
}

// removeRollupFile removes rollup file and interval after rollup job complete successfully
func (r *rollup) removeRollupFile(fileNumber table.FileNumber, interval timeutil.Interval) {
	var rs []timeutil.Interval
	intervals, ok := r.rollupFiles[fileNumber]
	if !ok {
		// not found
		return
	}
	for idx := range intervals {
		if interval != intervals[idx] {
			// keep other, remove completed interval
			rs = append(rs, intervals[idx])
		}
	}
	if len(rs) == 0 {
		delete(r.rollupFiles, fileNumber)
		return
	}
	r.rollupFiles[fileNumber] = rs
}

// getRollupFiles returns all need rollup files
func (r *rollup) getRollupFiles() map[table.FileNumber][]timeutil.Interval {
	result := make(map[table.FileNumber][]timeutil.Interval)
	for k, v := range r.rollupFiles {
		rs := make([]timeutil.Interval, len(v))
		copy(rs, v)
		result[k] = rs
	}
	return result
}

// addReferenceFile adds rollup reference file under target family
func (r *rollup) addReferenceFile(familyID FamilyID, fileNumber table.FileNumber) {
	files, ok := r.referenceFiles[familyID]
	if !ok {
		r.referenceFiles[familyID] = []table.FileNumber{fileNumber}
		return
	}
	for _, file := range files {
		if file == fileNumber {
			return
		}
	}
	files = append(files, fileNumber)
	r.referenceFiles[familyID] = files
}

// removeReferenceFile removes rollup reference file under target family
func (r *rollup) removeReferenceFile(familyID FamilyID, fileNumber table.FileNumber) {
	files, ok := r.referenceFiles[familyID]
	if !ok {
		return
	}
	var newFiles []table.FileNumber
	for _, file := range files {
		if file != fileNumber {
			newFiles = append(newFiles, file)
		}
	}
	if len(newFiles) == 0 {
		// if source files is empty, remove family reference
		delete(r.referenceFiles, familyID)
		return
	}
	r.referenceFiles[familyID] = newFiles
}

// getReferenceFiles returns the reference files under target family
func (r *rollup) getReferenceFiles() map[FamilyID][]table.FileNumber {
	result := make(map[FamilyID][]table.FileNumber)
	for k, v := range r.referenceFiles {
		d := make([]table.FileNumber, len(v))
		copy(d, v)
		result[k] = d
	}
	return result
}
