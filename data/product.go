package dataimport

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"io"
	"regexp"
	"time"
)

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

func (p * Product) FromJson(reader io.Reader) error  {
	e := json.NewDecoder(reader)
	return e.Decode(p)
}

func (p *Product) Validate() error{
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

// fomat : abc-asdf-asdfs
func validateSKU(f validator.FieldLevel) bool{

	re := regexp.MustCompile(`[a-z+]-[a-z]+-[a-z]`)
	matches := re.FindAllString(f.Field().String(), -1)

	if len(matches) != 1{
		return false
	}

	return true
}

// Products is a collection of Product
type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() Products {
	return productList
}

func AddProduct(p *Product){
	p.ID = getNextID()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error{
	_, pos, err := findNextProduct(id)

	if err != nil {
		return err
	}

	productList[pos] = p
	p.ID = id
	
	return nil
}

var ErrorProductNotFound error = fmt.Errorf("Product Not Found!")

func findNextProduct(id int) (*Product , int, error){
	for i, p := range productList {
		if p.ID ==  id{
			return p, i ,nil
		}
	}

	return nil, -1,  ErrorProductNotFound
}

func getNextID() int{
	lp := productList[len(productList) -1]
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
