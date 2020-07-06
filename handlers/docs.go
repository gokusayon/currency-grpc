// Package classification for Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//  - application/json
//
//	Produces:
//  - application/json
// swagger:meta
package handlers

import data "github.com/gokusayon/products-api/data"

// A list of products returns in response
// swagger:response listProducts
type listProductsWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// A product returns in response
// swagger:response productResponse
type productResponseWrapper struct {
	// All products in the system
	// in: body
	Body data.Product
}

// ProductsNoContent No content is returned by this API endpoint
// swagger:response noContent
type ProductsNoContent struct{}

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// swagger:parameters listSingleProduct listProducts
type productQueryParam struct {
	// Currency used when returning the price of the product.
	// when not specified it returns in GBP
	//	in: query
	//	required: false
	Currency string `json:"currency"`
}

// swagger:parameters deleteProduct listSingleProduct updateProduct
type productIDParam struct {
	// The ID of product to be deleted
	//	in: path
	//	required: true
	ID int `json:"id"`
}
