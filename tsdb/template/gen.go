package template

// need install tmpl first (go get github.com/benbjohnson/tmpl)
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@tag_store.data -o=../metadb/tag_store.gen.go int_map.tmpl
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@tag_store_test.data -o=../metadb/tag_store.gen_test.go int_map_test.tmpl
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@tag_index_map.data -o=../indexdb/tag_store.gen.go int_map.tmpl
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@tag_index_map_test.data -o=../indexdb/tag_store.gen_test.go int_map_test.tmpl
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@tag_index_store.data -o=../indexdb/tag_index_store.gen.go int_map.tmpl
//go:generate env GO111MODULE=on go run github.com/benbjohnson/tmpl -data=@tag_index_store_test.data -o=../indexdb/tag_index_store.gen_test.go int_map_test.tmpl
