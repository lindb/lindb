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

//go:build integration
// +build integration

package standalone

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/app/standalone"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

func TestMain(m *testing.M) {
	cfg := config.NewDefaultStandalone()
	if err := logger.InitLogger(cfg.Logging, "standalone.log"); err != nil {
		panic(fmt.Errorf("init logging err: %s", err))
	}

	gin.SetMode(gin.ReleaseMode)
	// run cluster as standalone mode
	runtime := standalone.NewStandaloneRuntime(config.Version, &cfg)
	if err := runtime.Run(); err != nil {
		panic(fmt.Errorf("run standalone cluster err: %s", err))
	}
	fmt.Println("standalone cluster start successfully")
	defer func() {
		fmt.Println("all tests run completed, now stop standalone cluster")
		runtime.Stop()
	}()
	m.Run()
}

func Test_QueryMetric(t *testing.T) {
	time.Sleep(30 * time.Second)
	cli := client.NewExecuteCli("http://localhost:9000/api")
	rs, err := cli.ExecuteAsResult(models.ExecuteParam{
		Database: "_internal",
		SQL:      "select used_percent from lindb.monitor.system.mem_stat  where time>now()-1h group by node,role",
	}, &models.ResultSet{})
	assert.NoError(t, err)
	fmt.Println(rs)
}
