.PHONY: swagger

check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger.yaml --scan-models

client:
	swagger generate client -f ./swagger.yaml -A products-api --principal