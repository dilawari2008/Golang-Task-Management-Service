protoc --proto_path=api/proto \
  --go_out=paths=source_relative:api/proto \
  --go-grpc_out=paths=source_relative:api/proto \
  api/proto/task.proto