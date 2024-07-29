run-purchases:
	@cd services/order-service && go run .

run-products:
	@cd services/products && go run .

gen:
	@protoc \
    --proto_path=protobuf "protobuf/event.proto" \
    --go_out=services/common/genproto/event --go_opt=paths=source_relative \
    --go-grpc_out=services/common/genproto/event --go-grpc_opt=paths=source_relative