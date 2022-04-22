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
	"fmt"
	"strconv"
	"strings"

	"github.com/lindb/lindb/kv/table"
)

const sstSuffix = "sst"

const TmpSuffix = "tmp"

const Lock = "LOCK"

const Options = "OPTIONS"

const ManifestPrefix = "MANIFEST-"

// FileType represents a file type.
type FileType int

// File types.
const (
	TypeManifest FileType = iota
	TypeJournal
	TypeTable
	TypeTemp
	TypeInfo
)

// FileDesc represents file type and file number
type FileDesc struct {
	FileType   FileType
	FileNumber table.FileNumber
}

// current returns current file name for saving manifest file name
func current() string {
	return "CURRENT"
}

// Table returns the sst's file name
func Table(fileNumber table.FileNumber) string {
	return fmt.Sprintf("%06d.%s", fileNumber, sstSuffix)
}

// ManifestFileName returns manifest file name
func ManifestFileName(fileNumber table.FileNumber) string {
	return fmt.Sprintf("%s%06d", ManifestPrefix, fileNumber)
}

// ParseFileName parses file name.
// if the file name was successfully parsed, returns file desc instance, else return nil.
func ParseFileName(fileName string) *FileDesc {
	if strings.HasSuffix(fileName, ".sst") {
		n, err := strconv.ParseInt(removeSuffix(fileName, ".sst"), 10, 64)
		if err != nil {
			return nil
		}
		return &FileDesc{
			FileType:   TypeTable,
			FileNumber: table.FileNumber(n),
		}
	}
	return nil
}

// removeSuffix removes suffix, then returns new string
func removeSuffix(value, suffix string) string {
	return value[0 : len(value)-len(suffix)]
}
