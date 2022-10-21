#!/bin/bash
set -e

readonly service="$1"

if [ ! -d "$service" ]; then
	mkdir -p "internal/genproto/$service"
	echo "Creating folder $service"
fi

protoc \
  -I=api/protobuf "api/protobuf/$service.proto" \
  "--go_out=internal/genproto/$service" --go_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=false \
  "--go-grpc_out=internal/genproto/$service" --go-grpc_opt=paths=source_relative

echo "Finish generate proto"