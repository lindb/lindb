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

package monitoring

import (
	"bufio"
	"io"
	"os"
	"path"

	"github.com/gin-gonic/gin"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	httppkg "github.com/lindb/lindb/pkg/http"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	LogListPath = "/log/list"
	LogViewPath = "/log/view"
)

// LoggerAPI represents view log file rest api.
type LoggerAPI struct {
	logDir string
	logger *logger.Logger
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
func (d *LoggerAPI) List(c *gin.Context) {
	logFiles, err := fileutil.ListDir(d.logDir)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	httppkg.OK(c, logFiles)
}

// View tails the log file, return the last n lines.
func (d *LoggerAPI) View(c *gin.Context) {
	var param struct {
		FileName string `form:"file" binding:"required"`
		// default: read last 1MB data from log file
		Size int64 `form:"size,default=1048576"`
	}
	err := c.ShouldBind(&param)
	if err != nil {
		httppkg.Error(c, err)
		return
	}
	file, err := os.Open(path.Join(d.logDir, param.FileName))
	if err != nil {
		httppkg.Error(c, err)
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
