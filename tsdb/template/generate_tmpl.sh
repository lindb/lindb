#!/usr/bin/env bash

set -ex

# install tmpl
if ! type tmpl >/dev/null;  then \
  go get github.com/benbjohnson/tmpl
fi

tmpl -data=@tag_store.data -o=../metadb/tag_store.gen.go int_map.tmpl
tmpl -data=@tag_store_test.data -o=../metadb/tag_store.gen_test.go int_map_test.tmpl
tmpl -data=@inverted_store.data -o=../indexdb/inverted_store.gen.go int_map.tmpl
tmpl -data=@inverted_store_test.data -o=../indexdb/inverted_store.gen_test.go int_map_test.tmpl
tmpl -data=@tag_index_store.data -o=../indexdb/tag_index_store.gen.go int_map.tmpl
tmpl -data=@tag_index_store_test.data -o=../indexdb/tag_index_store.gen_test.go int_map_test.tmpl
tmpl -data=@forward_store.data -o=../indexdb/forward_store.gen.go int_map.tmpl
tmpl -data=@forward_store_test.data -o=../indexdb/forward_store.gen_test.go int_map_test.tmpl
tmpl -data=@metric_bucket.data -o=../memdb/metric_bucket_store.gen.go int_map.tmpl
tmpl -data=@metric_bucket_test.data -o=../memdb/metric_bucket_store.gen_test.go int_map_test.tmpl
tmpl -data=@metric_store.data -o=../memdb/metric_store.gen.go int_map.tmpl
tmpl -data=@metric_store_test.data -o=../memdb/metric_store.gen_test.go int_map_test.tmpl