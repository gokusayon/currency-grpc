package dataimport

import (
	"context"
	"fmt"
	"time"

	protos "github.com/gokusayon/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the ID of the product
	//
	// required: true
	// min: 1
	ID int `json:"id"`

	// the Name of the Product
	//
	// required: true
	// max length: 25
	Name string `json:"name" validate:"required"`

	// the Description of the Product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`

	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// Products is a collection of Product
type Products []*Product

// ProductsDB is object for managing list of products
type ProductsDB struct {
	log      hclog.Logger
	currency protos.CurrencyClient
	client   protos.Currency_SubscribeClient
	rates    map[string]float64
}

// NewProductsDB returns a ProductDB handler
func NewProductsDB(log hclog.Logger, currency protos.CurrencyClient) *ProductsDB {
	pdb := &ProductsDB{log: log, currency: currency, client: nil, rates: map[string]float64{}}
	pdb.handleUpdates()
	return pdb
}

func (p *ProductsDB) handleUpdates() {
	p.log.Info("Registering client")
	sub, err := p.currency.Subscribe(context.Background())

	if err != nil {
		p.log.Debug("Error setting up the subscription", "error", err)
	}

	p.client = sub

	go func() {

		p.log.Info("Listening for updates ...")
		for {
			rr, err := p.client.Recv()
			p.log.Info("Recived Update for", "destination", rr.GetDestination(), "rate", rr.GetRate())

			if err != nil {
				p.log.Error("Error while waiting for message", "error", err)
				return
			}

			p.rates[rr.GetDestination().String()] = rr.Rate

		}
	}()

}

// GetProducts GETS list of products from database
func (p *ProductsDB) GetProducts(destination string) (Products, error) {
	if destination == "" {
		return productList, nil
	}

	rate, err := p.getRate(destination)

	if err != nil {
		p.log.Error("Unable to get currency rates", "currency", destination, "err", err)
		return nil, err
	}

	newProductList := Products{}
	for _, prod := range productList {
		temp := *prod
		temp.Price = temp.Price * rate

		newProductList = append(newProductList, &temp)
	}

	return newProductList, nil
}

// GetProductByID GETS a product with given ID from database
func (p *ProductsDB) GetProductByID(id int, destination string) (*Product, error) {
	index := findIndexByID(id)

	if index == -1 {
		return nil, ErrorProductNotFound
	}

	// Return default values if destination string is empty
	if destination == "" {
		return productList[index], nil
	}

	rate, err := p.getRate(destination)

	if err != nil {
		p.log.Error("Unable to get currency rates", "currency", destination, "err", err)
		return nil, err
	}

	pr := *productList[index]
	pr.Price = pr.Price * rate

	return &pr, nil

}

// getRate returns rate value for destination.
func (p *ProductsDB) getRate(destination string) (float64, error) {

	// check cache
	if r, ok := p.rates[destination]; ok {
		return r, nil
	}

	// Get initial rate
	rr := protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := p.currency.GetRate(context.Background(), &rr)

	if err != nil {
		p.log.Error("Unable to get currency rates", "currency", destination, "err", err)
		return 0, err
	}

	// Subscribe for updates
	p.client.Send(&rr)

	return resp.Rate, nil
}

// DeleteProduct DELETES a product with given ID from database
func (p *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByID(id)

	if i == -1 {
		return ErrorProductNotFound
	}

	productList = append(productList[:i], productList[i+1:]...)
	return nil
}

// AddProduct PUT a product into the database
func (p *ProductsDB) AddProduct(pr Product) {
	pr.ID = getNextID()
	productList = append(productList, &pr)
}

// UpdateProduct POST a product into the database
func (p *ProductsDB) UpdateProduct(id int, pr Product) error {
	i := findIndexByID(id)

	if i == -1 {
		return ErrorProductNotFound
	}

	productList[i] = &pr

	return nil
}

// ErrorProductNotFound in case the product does not exist in the database
var ErrorProductNotFound error = fmt.Errorf("Product Not Found")

func findIndexByID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
