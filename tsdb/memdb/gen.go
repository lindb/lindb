package memdb

// need install tmpl first (go get github.com/benbjohnson/tmpl)
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@block_store.gen.go.tmpldata block_store.gen.go.tmpl
