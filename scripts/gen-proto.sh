#!/bin/bash
echo "Running script with argument: $1"
CURRENT_DIR=$1
rm -rf "${CURRENT_DIR}/generated"
find "${CURRENT_DIR}/protos" -type f -name "*.proto" -print0 | while IFS= read -r -d '' file; do
  protoc -I=${CURRENT_DIR}/protos -I=/usr/local/go --go_out=${CURRENT_DIR} --go-grpc_out=${CURRENT_DIR} "${file}"
done  

