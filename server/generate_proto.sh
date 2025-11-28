#!/bin/bash

# 生成Go的protobuf代码
protoc --go_out=. --go_opt=paths=source_relative \
    pkg/proto/game.proto

echo "Protobuf code generated successfully!"
