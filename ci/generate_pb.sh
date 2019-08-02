#!/usr/bin/env bash

set -ex

# package name in xxx.proto should be xxx

echo "generate pb file"

GO_PREFIX_PATH=github.com/lindb/lindb/rpc/proto

function collect() {
    file=$(basename $1)
    base_name=$(basename $1 ".proto")
    mkdir -p ../proto/$base_name
    if [[ -z ${GO_OUT_M} ]];then
        GO_OUT_M="M$file=$GO_PREFIX_PATH/$base_name"
     else
        GO_OUT_M="$GO_OUT_M,M$file=$GO_PREFIX_PATH/$base_name"
     fi
}


function gen() {
    base_name=$(basename $1 ".proto")
#    mkdir -p ../pkg/$base_name
    protoc -I. --gofast_out=plugins=grpc,$GO_OUT_M:../proto/$base_name $1
}

cd rpc/pb

for file in `ls *.proto`
    do
    collect $file
done

for file in `ls *.proto`
     do
     gen $file
done





