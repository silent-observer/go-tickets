#!/usr/bin/env bash

pushd proto
protoc --proto_path=. --go_out=../proto_go --go_opt=paths=source_relative *.proto
popd