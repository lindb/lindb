#!/usr/bin/env bash

set -ex

# install go-licenser(https://github.com/elastic/go-licenser)
if ! type go-licenser >/dev/null;  then \
  go install github.com/elastic/go-licenser@v0.4.0
fi
go-licenser -license ASL2 -licensor LinDB ../
