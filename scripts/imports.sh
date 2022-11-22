#!/bin/sh

# need install https://github.com/incu6us/goimports-reviser
for f in $(find . -name "*.go"); do echo ${f} && goimports-reviser -rm-unused -set-alias -format -file-path ${f}; done
