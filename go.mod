module github.com/gokusayon/products-api

go 1.13

require (
	github.com/go-kit/kit v0.10.0
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/runtime v0.19.19
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gokusayon/currency v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/hashicorp/go-hclog v0.14.1
	github.com/leodido/go-urn v1.2.0 // indirect
	google.golang.org/grpc v1.30.0
)

replace github.com/gokusayon/currency => ../currency
