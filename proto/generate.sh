#!/usr/bin/env bash

set -ex

# package name in xxx.proto should be xxx

echo "generate pb file"

GO_PREFIX_PATH=github.com/lindb/lindb/proto

function collect() {
    # v1/replica.proto
    file=$1/$(basename $2)
    base_name=$1/$(basename $2 ".proto")
    mkdir -p ../gen/$base_name
    if [[ -z ${GO_OUT_M} ]];then
        GO_OUT_M="M$file=$GO_PREFIX_PATH/$base_name"
     else
        GO_OUT_M="$GO_OUT_M,M$file=$GO_PREFIX_PATH/$base_name"
     fi
}

function gen() {
    dir_name=$1
    base_name=$(basename $2 ".proto")
    protoc -I. --gofast_out=plugins=grpc,$GO_OUT_M:../gen/$dir_name/$base_name $2
}


# for dir in v1
for dir in v1; do

  cd proto/$dir

  for file in `ls *.proto`; do
    collect $dir $file
  done

  for file in `ls *.proto`; do
    gen $dir $file
  done
  cd ../../
  pwd
done


