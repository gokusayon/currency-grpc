.PHONY: protos

# New
protos: 
	protoc --proto_path=protos --go_out=plugins=grpc:protos/currency protos/currency.proto

# Old
# protoc --proto_path=protos --go_out=protos/currency --go-grpc_out=protos/currency  protos/currency.proto
