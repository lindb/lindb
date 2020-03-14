module github.com/lindb/lindb

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/GeertJohan/go.rice v1.0.0
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/antlr/antlr4 v0.0.0-20190623224521-a770ff26ccc4
	github.com/benbjohnson/tmpl v1.0.0 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/coreos/bbolt v1.3.3
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/damnever/goctl v1.1.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.3.2
	github.com/golang/snappy v0.0.1
	github.com/google/btree v1.0.0 // indirect
	github.com/gorilla/mux v1.7.2
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hillbig/rsdic v0.0.0-20150805052524-6158e7a2d824
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.7
	github.com/lindb/roaring v0.0.0-00010101000000-000000000000
	github.com/mattn/go-isatty v0.0.8
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4
	github.com/prometheus/common v0.7.0
	github.com/prometheus/prometheus v2.15.2+incompatible
	github.com/shirou/gopsutil v2.19.9+incompatible
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/uber-go/tally v3.3.13+incompatible
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.5.0
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/sys v0.0.0-20191010194322-b09406accb47
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/grpc v1.22.1
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

// just redirect to local repo for local debug
// replace github.com/lindb/roaring => /Users/jie.huang/go/src/github.com/lindb/roaring
replace github.com/lindb/roaring => github.com/lindb/roaring v0.4.22-0.20200211075929-6661b4a242fa
