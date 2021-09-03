@echo off
protoc --proto_path=../api/protobuf --go_out=../ --go-grpc_out=.././ ../api/protobuf/internal/*.proto