run-client:
	@go run ./*.go

gen:
	@protoc \
		--proto_path=protobuf "protobuf/clipBoardSync.proto" \
		--go_out=genproto/clipboardSync --go_opt=paths=source_relative \
  	--go-grpc_out=genproto/clipboardSync --go-grpc_opt=paths=source_relative
