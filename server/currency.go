package server

import (
	"context"
	"io"
	"time"

	data "github.com/gokusayon/currency/data"
	protos "github.com/gokusayon/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

type Currency struct {
	log         hclog.Logger
	er          *data.ExchangeRate
	subscribers map[protos.Currency_SubscribeServer][]*protos.RateRequest
}

func NewCurrency(log hclog.Logger, er *data.ExchangeRate) *Currency {
	c := &Currency{log: log, er: er, subscribers: make(map[protos.Currency_SubscribeServer][]*protos.RateRequest)}
	go c.handleUpdates()
	return c
}

func (c *Currency) GetRate(ctx context.Context, req *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())
	rate, _ := c.er.GetRates(req.GetBase().String(), req.GetDestination().String())
	return &protos.RateResponse{Base: req.GetBase(), Destination: req.GetDestination(), Rate: rate}, nil
}

func (c *Currency) Subscribe(src protos.Currency_SubscribeServer) error {

	for {
		// gRPC blocking client untill a new request for subscription is recieved
		rr, err := src.Recv()
		c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("Client has closed connection")
			break
		}

		if err != nil {
			return err
		}

		// Sanity check if key exists
		rrs, ok := c.subscribers[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscribers[src] = rrs
	}

	return nil
}

func (c *Currency) handleUpdates() {
	//TODO: handle updates for subscibers
	c.log.Info("handling updates")
	ch := c.er.MonitorRates(5 * time.Second)

	// Iterate over channels
	// Monitor rates pushes the data to the channel. Its consumed using for loop.
	for range ch {
		c.log.Info("Data Updated")

		//  Iterate over subscriptions
		for key, vals := range c.subscribers {

			// Iterate over all rate requests
			for _, rr := range vals {

				rate, err := c.er.GetRates(rr.GetBase().String(), rr.GetDestination().String())

				if err != nil {
					c.log.Error("Error fetching rates", key, "base", rr.GetBase(), "Destination", rr.GetDestination())
				}

				err = key.Send(&protos.RateResponse{Base: rr.GetBase(), Destination: rr.GetDestination(), Rate: rate})

				if err != nil {
					c.log.Error("Error publishing response. Removing subscription", "base", rr.GetBase(), "Destination", rr.GetDestination())
					delete(c.subscribers, key)
				}
			}
		}
	}
}
