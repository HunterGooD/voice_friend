
generate:
	protoc -I protos protos/user_service/*.proto \
	 --go_out=./gen/go --go_opt=paths=source_relative \
	 --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative \
	 --experimental_allow_proto3_optional
