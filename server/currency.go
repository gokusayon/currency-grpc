package server

import (
	"context"

	data "github.com/gokusayon/currency/data"
	protos "github.com/gokusayon/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	log hclog.Logger
	er  *data.ExchangeRate
}

func NewCurrency(log hclog.Logger, er *data.ExchangeRate) *Currency {
	return &Currency{log: log, er: er}
}

func (c *Currency) GetRate(ctx context.Context, req *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())
	rate, _ := c.er.GetRates(req.GetBase().String(), req.GetDestination().String())
	c.log.Info("Currency rate is: ", "currency", req.GetDestination(), "rate", rate)
	return &protos.RateResponse{Rate: rate}, nil
}
