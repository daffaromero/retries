run-purchases:
	@cd services/order-service && go run .

run-products:
	@cd services/products && go run .

gen-event:
	@protoc \
    --proto_path=protobuf "protobuf/event.proto" \
    --go_out=services/common/genproto/event --go_opt=paths=source_relative \
    --go-grpc_out=services/common/genproto/event --go-grpc_opt=paths=source_relative

gen-order:
	@protoc \
    --proto_path=protobuf "protobuf/orders.proto" \
    --go_out=services/common/genproto/orders --go_opt=paths=source_relative \
    --go-grpc_out=services/common/genproto/orders --go-grpc_opt=paths=source_relative