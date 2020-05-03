go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/GeertJohan/go.rice v1.0.0
	github.com/OneOfOne/xxhash v1.2.2
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/antlr/antlr4 v0.0.0-20200309161749-1284814c2112
	github.com/benbjohnson/tmpl v1.0.0 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/damnever/goctl v1.2.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/gogo/protobuf v1.2.1
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1
	github.com/gorilla/mux v1.7.4
	github.com/json-iterator/go v1.1.7
	github.com/lindb/roaring v0.0.0-00010101000000-000000000000
	github.com/m3db/prometheus_client_golang v0.8.1 // indirect
	github.com/m3db/prometheus_client_model v0.1.0 // indirect
	github.com/m3db/prometheus_common v0.1.0 // indirect
	github.com/m3db/prometheus_procfs v0.8.1 // indirect
	github.com/mattn/go-isatty v0.0.4
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/prometheus/common v0.4.1
	github.com/shirou/gopsutil v2.20.2+incompatible
	github.com/spf13/cobra v0.0.3
	github.com/stretchr/testify v1.4.0
	github.com/uber-go/tally v3.3.15+incompatible
	go.etcd.io/bbolt v1.3.4
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200320040136-0eee733220fc
	go.uber.org/atomic v1.6.0
	go.uber.org/zap v1.14.1
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5
	golang.org/x/tools v0.0.0-20191029190741-b9c20aec41a5
	google.golang.org/grpc v1.26.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

// just redirect to local repo for local debug
// replace github.com/lindb/roaring => /Users/jie.huang/go/src/github.com/lindb/roaring
replace github.com/lindb/roaring => github.com/lindb/roaring v0.4.22-0.20200211075929-6661b4a242fa

module github.com/lindb/lindb
