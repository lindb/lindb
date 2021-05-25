#!/usr/bin/env bash

set -ex

# install go-licenser
if ! type go-licenser >/dev/null;  then \
  go get -u github.com/elastic/go-licenser
fi
go-licenser -license ASL2 -licensor LinDB