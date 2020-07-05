.PHONY: protos

protos:
	protoc --proto_path=protos --go_out=protos/currency --go-grpc_out=protos/currency protos/currency.proto
