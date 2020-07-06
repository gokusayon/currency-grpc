package handlers

import (
	"net/http"

	dataimport "github.com/gokusayon/products-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: listProducts

// GetProducts returns the products from datastore
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.log.Debug("Handle GET All Products")

	cur := r.URL.Query().Get("currency")

	lp, err := p.productsDB.GetProducts(cur)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		dataimport.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// lp := dataimport.GetProducts()
	err = dataimport.ToJSON(lp, rw)
	if err != nil {
		p.log.Error("Unable to Encode Products")
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Returns a product from list of products
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle returns a products based on ID from database
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	p.log.Debug("Handle GET Product By ID")

	id := getProductID(r)
	cur := r.URL.Query().Get("currency")

	prod, err := p.productsDB.GetProductByID(id, cur)

	switch err {
	case nil:

	case dataimport.ErrorProductNotFound:
		p.log.Error("Unable to find product", "id", id)

		rw.WriteHeader(http.StatusNotFound)
		dataimport.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.log.Error("Unable to find product", "id", id)

		rw.WriteHeader(http.StatusInternalServerError)
		dataimport.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = dataimport.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.log.Error("[ERROR] serializing product", err)
	}
}
