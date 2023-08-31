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

package api

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	httppkg "github.com/lindb/common/pkg/http"
	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/constants"
)

// for testing
var (
	readDirFn = os.ReadDir
	relFn     = filepath.Rel
	absFn     = filepath.Abs
	openFn    = os.Open
)

// FileInfo represents file info include name/size.
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

var (
	LogListPath = "/log/list"
	LogViewPath = "/log/view"
)

// LoggerAPI represents view log file rest api.
type LoggerAPI struct {
	logDir string
	logger logger.Logger
}

// NewLoggerAPI creates log view api instance.
func NewLoggerAPI(logDir string) *LoggerAPI {
	return &LoggerAPI{
		logDir: logDir,
		logger: logger.GetLogger("Monitoring", "ExploreAPI"),
	}
}

// Register adds explore url route.
func (d *LoggerAPI) Register(route gin.IRoutes) {
	route.GET(LogListPath, d.List)
	route.GET(LogViewPath, d.View)
}

// List returns all log files in log dir.

// @Summary list log files
// @Description return all log files in log dir.
// @Tags State
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal error"
// @Router /log/list [get]
func (d *LoggerAPI) List(c *gin.Context) {
	files, err := readDirFn(d.logDir)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	var logFiles []FileInfo
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".log") {
			fileInfo, err := file.Info()
			if err != nil {
				httppkg.Error(c, err)
				return
			}
			logFiles = append(logFiles, FileInfo{
				Name: name,
				Size: fileInfo.Size(),
			})
		}
	}
	httppkg.OK(c, logFiles)
}

// View tails the log file, return the last n lines.
// @Summary tail log file
// @Description return last N lines in log file.
// @Tags State
// @Accept json
// @Produce plain
// @Success 200 {string} string
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal error"
// @Router /log/view [get]
func (d *LoggerAPI) View(c *gin.Context) {
	var param struct {
		FileName string `form:"file" binding:"required"`
		// default: read last 1MB data from log file
		Size int64 `form:"size,default=1048576"`
	}
	err := c.ShouldBindQuery(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	// prepend slash for cleaning relative paths
	requestedFile := filepath.Clean(filepath.Join(string(os.PathSeparator), param.FileName))
	rel, err := relFn(string(os.PathSeparator), requestedFile)
	if err != nil {
		// slash is prepended above therefore this is not expected to fail
		httppkg.Error(c, fmt.Errorf("failed to get the relative path"))
		d.logger.Error("failed to get the relative path", logger.Error(err))
		return
	}
	absLogDir, err := absFn(d.logDir)
	if err != nil {
		httppkg.Error(c, fmt.Errorf("failed to get log absolute path"))
		d.logger.Error("failed to get log absolute path", logger.Error(err))
		return
	}

	logFilePath := filepath.Join(absLogDir, rel)
	file, err := openFn(logFilePath)
	if err != nil {
		httppkg.Error(c, fmt.Errorf("failed to open log file: %s", param.FileName))
		d.logger.Error("failed to open log file", logger.Error(err))
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			d.logger.Warn("close file err",
				logger.String("file", param.FileName),
				logger.Error(err))
		}
	}()
	stat, err := file.Stat()
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	if stat.Size() > param.Size {
		// if log file size > read size, need skip
		_, err = file.Seek(stat.Size()-param.Size, io.SeekStart)
		if err != nil {
			httppkg.Error(c, err)
			return
		}
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip first line
	c.Stream(func(w io.Writer) bool {
		for scanner.Scan() {
			if err := writeLine(w, [][]byte{scanner.Bytes(), constants.LBBytes}); err != nil {
				d.logger.Warn("write log data to response stream err",
					logger.String("file", param.FileName),
					logger.Error(err))
				return false
			}
		}
		return false
	})
}

// writeLine writes a line into stream.
func writeLine(w io.Writer, data [][]byte) error {
	for _, d := range data {
		_, err := w.Write(d)
		if err != nil {
			return err
		}
	}
	return nil
}
