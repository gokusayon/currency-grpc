package server

import (
	"context"
	"io"
	"time"

	data "github.com/gokusayon/currency/data"
	protos "github.com/gokusayon/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct {
	log         hclog.Logger
	er          *data.ExchangeRate
	subscribers map[protos.Currency_SubscribeServer][]*protos.RateRequest
}

// NewCurrency returns the handler for the @Currency
func NewCurrency(log hclog.Logger, er *data.ExchangeRate) *Currency {
	c := &Currency{log: log, er: er, subscribers: make(map[protos.Currency_SubscribeServer][]*protos.RateRequest)}
	go c.handleUpdates()
	return c
}

// GetRates returns the exchange rates for the given base and destination currency
func (c *Currency) GetRate(ctx context.Context, req *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())

	//  if base currency is same as the destionation currency then return rich client error
	if req.GetBase() == req.GetDestination() {
		c.log.Info("Unable to fetch rates as base and destionation currency are same", "base", req.GetBase(), "destination", req.GetDestination())

		st := status.Newf(
			codes.InvalidArgument,
			"Base curreny %s can not be same as the destination %s",
			req.GetBase().String(),
			req.GetDestination().String(),
		)

		err, wde := st.WithDetails(req)

		if wde != nil {
			return nil, wde
		}

		return nil, err.Err()
	}

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
			c.log.Error("Client has closed connection")
			break
		}

		if err != nil {
			c.log.Error("Unable to fetch rates as base and destionation currency are same", "err", err)
			return err
		}

		// Sanity check if key exists
		rrs, ok := c.subscribers[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		// check if sub is already added to the subscibers list
		var validationErr *status.Status
		for _, val := range rrs {

			// if this currency is already subscribed then return error
			if val.Base == rr.Base && val.Destination == rr.Destination {

				validationErr = status.Newf(codes.AlreadyExists, "Subscription active for the currenct destination")

				validationErr, err = validationErr.WithDetails(rr)

				if err != nil {
					c.log.Error("Unable to add metadata to error.", "err", err)
				}

				break
			}
		}

		// If validationErr is not nil then return error and continue
		if validationErr != nil {
			c.log.Error("Unable to subscibe", "err", validationErr.Message())

			rrs := &protos.StreamingRateResponse_Error{Error: validationErr.Proto()}
			err := &protos.StreamingRateResponse{Message: rrs}

			src.Send(err)

			continue
		}

		rrs = append(rrs, rr)
		c.subscribers[src] = rrs
	}

	return nil
}

func (c *Currency) handleUpdates() {
	c.log.Info("handling updates")
	ch := c.er.MonitorRates(5 * time.Second)

	// Iterate over channels
	// Monitor rates pushes the data to the channel. Its consumed using for loop over channel.
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

				resp := &protos.StreamingRateResponse{
					Message: &protos.StreamingRateResponse_RateResponse{
						RateResponse: &protos.RateResponse{Base: rr.GetBase(), Destination: rr.GetDestination(), Rate: rate},
					},
				}

				err = key.Send(resp)

				// if unable to publish to client then remove entry from the subsciption list
				if err != nil {
					c.log.Error("Error publishing response. Removing subscription", "base", rr.GetBase(), "Destination", rr.GetDestination())
					delete(c.subscribers, key)
				}
			}
		}
	}
}
