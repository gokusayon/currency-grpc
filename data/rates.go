package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRate struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewExchangeRate(log hclog.Logger) (*ExchangeRate, error) {
	er := &ExchangeRate{log: log, rates: map[string]float64{}}
	err := er.getRates()
	return er, err
}

func (e *ExchangeRate) GetRates(base, dest string) (float64, error) {
	br, ok := e.rates[base]

	if !ok {
		return 0, fmt.Errorf("Currency not found", "currency", base)
	}

	dr, ok := e.rates[dest]

	if !ok {
		return 0, fmt.Errorf("Currency not found", "currency", dest)
	}

	return dr / br, nil
}

func (e *ExchangeRate) getRates() error {

	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	md := &ExtractCubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	for _, val := range md.CubeData {
		r, err := strconv.ParseFloat(val.Rate, 64)

		if err != nil {
			return err
		}
		e.rates[val.Currency] = r
	}

	e.rates["EUR"] = 1

	return nil
}

type ExtractCubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"` // Must make camel case so they are public
	Rate     string `xml:"rate,attr"`
}
